CREATE TABLE IF NOT EXISTS public.checkout_items (
    checkout_id uuid NOT NULL,
    product_id uuid NOT NULL,
    quantity int NOT NULL,
    amount int NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone,
    CONSTRAINT checkout_items_pkey PRIMARY KEY (checkout_id, product_id),
    CONSTRAINT quantity CHECK (quantity >= 1)
)