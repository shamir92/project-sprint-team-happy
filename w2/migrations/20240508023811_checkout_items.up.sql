CREATE TABLE IF NOT EXISTS public.checkout_items (
    checkout_id uuid NOT NULL,
    product_id uuid NOT NULL,
    quantity int NOT NULL,
    amount int NOT NULL,

    CONSTRAINT checkout_items_pkey PRIMARY KEY (checkout_id, product_id),
    CONSTRAINT quantity CHECK (quantity >= 1)
)