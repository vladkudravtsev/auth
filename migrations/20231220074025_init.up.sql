CREATE TABLE IF NOT EXISTS users (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	"name" text NULL,
	email text NULL,
	password_hash text NULL,
	CONSTRAINT idx_users_email UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users USING btree (deleted_at);

CREATE TABLE IF NOT EXISTS apps (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	"name" text NULL,
	secret text NULL,
	CONSTRAINT apps_pkey PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_apps_deleted_at ON apps USING btree (deleted_at);
