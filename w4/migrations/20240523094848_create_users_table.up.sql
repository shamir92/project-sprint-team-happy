CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.users(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    username character varying(30) NOT NULL,
    email character varying(100) NOT NULL,
    password character(72) DEFAULT NULL,
    role character varying(20) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_username_unique UNIQUE (username)
);

CREATE INDEX users_username_idx ON users (username);