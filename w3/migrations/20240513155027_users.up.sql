CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.users(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    nip character(13) NULL,
    name character varying(50) NOT NULL,
    password character(72) DEFAULT NULL,
    role character varying(20) NOT NULL,
    identity_card_scan_img text NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone DEFAULT NULL,

    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_nip_unique UNIQUE (nip)
);

CREATE INDEX users_nip_idx ON users (nip);