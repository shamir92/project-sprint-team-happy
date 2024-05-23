CREATE TABLE IF NOT EXISTS public.order_items (
    order_id uuid NOT NULL,
    order_item_id uuid NOT NULL,
    price int NOT NULL,
    quantity int NOT NULL,
    amount int NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone,
    CONSTRAINT order_items_pkey PRIMARY KEY (order_id, order_item_id),
    CONSTRAINT quantity CHECK (quantity >= 1),
    CONSTRAINT price CHECK (price >= 1),
    CONSTRAINT amount CHECK (amount >= 1)
);