CREATE TABLE public.main_devices (
	id bigserial NOT NULL,
	user_id int8 NOT NULL,
	platform varchar NOT NULL,
	"version" varchar NOT NULL,
	device_id varchar NOT NULL,
	create_time timestamptz NOT NULL,
	update_time timestamptz NOT NULL,
	CONSTRAINT main_devices_pkey PRIMARY KEY (id),
	CONSTRAINT main_devices_unique UNIQUE (platform, version, device_id)
);