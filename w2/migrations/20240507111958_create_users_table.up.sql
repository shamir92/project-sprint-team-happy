CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.users (
    user_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name character varying(50) NOT NULL,
    phone_number character varying(16) NOT NULL,
    password character(72) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone,
    CONSTRAINT users_pkey PRIMARY KEY (user_id),
    CONSTRAINT users_phone_number_unique UNIQUE (phone_number)
);

CREATE INDEX idx_users_phone_number ON users (phone_number);