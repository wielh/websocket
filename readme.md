# 專案介紹

## 簡介

這是一個範例專案。用 websocket 通訊，將 main_device 的訊息通過此 api 進行轉發。

## 架構

+ 分層
    + controller: 負責接收/轉發 http/websocket 請求
    + service: 業務核心邏輯
    + repository: db 相關程式
    + log: 日誌
    + dto: controller 與 service 參數定義
    + model: service 與 repositroy 之間的參數定義

+  服務
    + user: 負責一般的註冊，登入等用戶邏輯
    + device: 將裝置登記到 user 底下
    + commuication: 驗證是登記的 device 後，建立 webscoket 連線。

+ 用到的工具: go, gin, postgreSQL, redis

## 測試
+ 工具: postman, ngrok
+ postman websocket 測試工具說明:
    + postman 左上角的 New 按鍵中按下後選擇 websocket 即可進行測試。
    + 測試網址為 wws://(host)/(path)
    + postman websocket 暫時沒有自動將 cookie 附上的功能，需要手動貼到 cookie
+ ngrok 說明
    + ngrok 是一個反向代理（reverse proxy）+ 隧道工具，能把你本地的 Web 伺服器透過一個公共 URL 暴露到外網。
    + 安裝好 ngrok 後，執行 ngrok http (port)，即可將在代理在 (port) 運行的api，terminal上會顯示代理網址。

## TODO
+ 如果此API並非單體架構，處理 main_device 與 sub_device 連接問題。
    

