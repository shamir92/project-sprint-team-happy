CREATE TABLE IF NOT EXISTS public.order_items (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    order_id uuid NOT NULL,
    order_item_id uuid NOT NULL,
    price int NOT NULL,
    quantity int NOT NULL,
    amount int NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone,
    CONSTRAINT order_items_pkey PRIMARY KEY (id),
    CONSTRAINT quantity CHECK (quantity >= 1),
    CONSTRAINT price CHECK (price >= 1),
    CONSTRAINT amount CHECK (amount >= 1),

    CONSTRAINT order_items_order_id_fk
        FOREIGN KEY (order_id)
        REFERENCES public.orders(id)
        ON DELETE CASCADE,

    CONSTRAINT order_items_order_item_id_fk
        FOREIGN KEY (order_item_id)
        REFERENCES public.merchant_items(id)
        ON DELETE CASCADE
);