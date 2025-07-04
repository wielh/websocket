package service

import (
	"context"
	"device-communication/src/dto"
	"device-communication/src/dtoError"
	"device-communication/src/repository"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type CommunicationSerivice interface {
	MainDeviceConnection(ctx context.Context, req *dto.MainDeviceConnectionRequest, w http.ResponseWriter, r *http.Request) *dtoError.ServiceError
	SubDeviceConnection(ctx context.Context, req *dto.SubDeviceConnectionRequest, w http.ResponseWriter, r *http.Request) *dtoError.ServiceError
}

type communicationSeriviceImpl struct {
	deviceRepo             repository.DeviceRepository
	errWarpper             dtoError.ServiceErrorWarpper
	socket                 websocket.Upgrader
	rooms                  webSocketRoomArray
	mainDeviceIdleDuration time.Duration
}

type webSocketRoom struct {
	MainConnection *websocket.Conn
	SubConnections map[uint64]*websocket.Conn
	mu             sync.Mutex
}

type webSocketRoomArray struct {
	rooms                 map[string]*webSocketRoom
	MAX_ROOM_NUMBER       int64
	MAX_SUB_DEVICE_NUMBER int64
	mu                    sync.RWMutex
}

func (w *webSocketRoom) SendMessage(message []byte) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for key, conn := range w.SubConnections {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			conn.Close()
			delete(w.SubConnections, key)
		}
	}
}

func (w *webSocketRoomArray) GetRoomKey(userId uint64, mainDeviceId uint64) string {
	return fmt.Sprintf("%d::%d", userId, mainDeviceId)
}

func (w *webSocketRoomArray) GetOrCreateRoom(userId uint64, mainDeviceId uint64, mainConnection *websocket.Conn) (*webSocketRoom, bool) {
	key := w.GetRoomKey(userId, mainDeviceId)
	w.mu.Lock()
	defer w.mu.Unlock()

	room, exists := w.rooms[key]
	if exists {
		return room, true
	}

	newRoom := &webSocketRoom{
		MainConnection: mainConnection,
		SubConnections: make(map[uint64]*websocket.Conn),
	}
	w.rooms[key] = newRoom
	return newRoom, false
}

func (w *webSocketRoomArray) JoinRoom(userId uint64, mainDeviceId uint64, subDeviceId uint64, subConnection *websocket.Conn) string {
	key := w.GetRoomKey(userId, mainDeviceId)
	w.mu.Lock()
	defer w.mu.Unlock()
	room, ok := w.rooms[key]
	if !ok {
		return "room not exist"
	}

	room.mu.Lock()
	defer room.mu.Unlock()
	if len(room.SubConnections) >= int(w.MAX_SUB_DEVICE_NUMBER) {
		return fmt.Sprintf("number of sub_device should <= %d", w.MAX_SUB_DEVICE_NUMBER)
	}

	_, ok = room.SubConnections[subDeviceId]
	if ok {
		return "this subdevice has already joined"
	}

	room.SubConnections[subDeviceId] = subConnection
	return ""
}

func (w *webSocketRoomArray) LeaveRoom(userId, mainDeviceId, subDeviceId uint64) {
	key := w.GetRoomKey(userId, mainDeviceId)
	w.mu.Lock()
	defer w.mu.Unlock()
	if room, ok := w.rooms[key]; ok {
		room.mu.Lock()
		defer room.mu.Unlock()
		delete(room.SubConnections, subDeviceId)
	}
}

func (w *webSocketRoomArray) RemoveRoom(userId uint64, mainDeviceId uint64) {
	key := w.GetRoomKey(userId, mainDeviceId)
	w.mu.Lock()
	room, ok := w.rooms[key]
	if ok {
		delete(w.rooms, key)
	}
	w.mu.Unlock()
	if !ok {
		return
	}

	room.mu.Lock()
	if room.MainConnection != nil {
		_ = room.MainConnection.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "room closed by server"))
		_ = room.MainConnection.Close()
	}
	for subId, conn := range room.SubConnections {
		_ = conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "main device offline"))
		_ = conn.Close()
		delete(room.SubConnections, subId)
	}
	room.mu.Unlock()
}

func (c *communicationSeriviceImpl) MainDeviceConnection(
	ctx context.Context, req *dto.MainDeviceConnectionRequest, writer http.ResponseWriter, httpRequest *http.Request) *dtoError.ServiceError {
	ok, err := c.deviceRepo.CheckMainDeviceBinding(ctx, req.UserId, req.MainDeviceId)
	if err != nil {
		return c.errWarpper.NewDBServiceError(err)
	} else if !ok {
		return c.errWarpper.NewMainDeviceNotBindingError()
	}

	conn, err := c.socket.Upgrade(writer, httpRequest, nil)
	if err != nil {
		return c.errWarpper.NewWebsocketUpgradeFailedError(err)
	}
	defer conn.Close()

	room, exists := c.rooms.GetOrCreateRoom(req.UserId, req.MainDeviceId, conn)
	if exists {
		conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "This main device already has a websocket connection"))
		return nil
	}
	defer c.rooms.RemoveRoom(req.UserId, req.MainDeviceId)
	conn.SetReadDeadline(time.Now().Add(c.mainDeviceIdleDuration))

	for {
		msgType, msg, err := room.MainConnection.ReadMessage()
		if err != nil {
			return nil
		}

		conn.SetReadDeadline(time.Now().Add(c.mainDeviceIdleDuration))
		switch msgType {
		case websocket.TextMessage:
			room.SendMessage(msg)
		case websocket.CloseMessage:
			return nil
		}
	}
}

func (c *communicationSeriviceImpl) SubDeviceConnection(ctx context.Context, req *dto.SubDeviceConnectionRequest, writer http.ResponseWriter, httpRequest *http.Request) *dtoError.ServiceError {
	ok, err := c.deviceRepo.CheckSubDeviceBinding(ctx, req.UserId, req.MainDeviceId, req.SubDeviceId)
	if err != nil {
		return c.errWarpper.NewDBServiceError(err)
	} else if !ok {
		return c.errWarpper.NewSubDeviceNotBindingError()
	}

	conn, err := c.socket.Upgrade(writer, httpRequest, nil)
	if err != nil {
		return c.errWarpper.NewWebsocketUpgradeFailedError(err)
	}
	defer conn.Close()

	errMessage := c.rooms.JoinRoom(req.UserId, req.MainDeviceId, req.SubDeviceId, conn)
	defer c.rooms.LeaveRoom(req.UserId, req.MainDeviceId, req.SubDeviceId)
	if errMessage != "" {
		conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, errMessage))
		return nil
	}

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return nil
		}
	}
}

var communication CommunicationSerivice

func init() {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	communication = &communicationSeriviceImpl{
		errWarpper: dtoError.GetServiceErrorWarpper(),
		deviceRepo: repository.GetDeviceRepository(),
		socket:     upgrader,
		rooms: webSocketRoomArray{
			rooms:                 make(map[string]*webSocketRoom),
			MAX_ROOM_NUMBER:       100,
			MAX_SUB_DEVICE_NUMBER: 1,
		},
		mainDeviceIdleDuration: 2 * time.Hour,
	}
}

func GetCommunicationSerivice() CommunicationSerivice {
	return communication
}
