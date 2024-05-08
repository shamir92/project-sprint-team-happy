CREATE TABLE IF NOT EXISTS public.customers (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name character varying(50) NOT NULL,
    phone_number character varying(16) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT customers_pkey PRIMARY KEY (id),
    CONSTRAINT customers_phone_number_unique UNIQUE (phone_number)
);

CREATE INDEX idx_customers_phone_number ON customers (phone_number);