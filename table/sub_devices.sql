CREATE TABLE public.sub_devices (
	id bigserial NOT NULL,
	main_device_id int8 NOT NULL,
	platform varchar NOT NULL,
	"version" varchar NOT NULL,
	device_id varchar NOT NULL,
	create_time timestamptz NOT NULL,
	update_time timestamptz NOT NULL,
	CONSTRAINT sub_devices_pkey PRIMARY KEY (id),
	CONSTRAINT sub_devices_unique UNIQUE (platform, version, device_id)
);

