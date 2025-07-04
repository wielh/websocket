CREATE TABLE public.users (
	id bigserial NOT NULL,
	username varchar NOT NULL,
	"password" varchar NOT NULL,
	"name" varchar NOT NULL,
	email varchar NOT NULL,
	create_time timestamptz NOT NULL,
	update_time timestamptz NOT NULL,
	delete_time timestamptz NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT users_unique UNIQUE (username)
);
CREATE INDEX idx_users_deleted_at ON public.users USING btree (delete_time);