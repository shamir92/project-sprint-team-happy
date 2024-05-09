CREATE TABLE IF NOT EXISTS public.checkouts (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    customer_id uuid NOT NULL,
    paid int NOT NULL,
    change int NOT NULL,
    created_by uuid NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone,

    CONSTRAINT checkouts_pkey PRIMARY KEY (id),
    CONSTRAINT checkouts_paid CHECK (paid >= 1),
    CONSTRAINT checkouts_change CHECK (change >= 0)
)