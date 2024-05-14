CREATE TABLE IF NOT EXISTS public.patients (
    id character(16) NOT NULL,
    name character varying(50) NOT NULL,
    phone_number character varying(16) NOT NULL,
    birth_date date NOT NULL,
    gender character varying(6) NOT NULL,
    identity_card_scan_img text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT patients_pkey PRIMARY KEY (id)
);