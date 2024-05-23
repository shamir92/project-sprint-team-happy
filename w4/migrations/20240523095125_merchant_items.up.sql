CREATE TABLE IF NOT EXISTS public.merchant_items(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    merchant_id uuid NOT NULL,
    name char varying(30) NOT NULL,
    category char varying(20) NOT NULL,
    price int not null,
    image_url text not null,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT merchant_items_pkey PRIMARY KEY (id),
    CONSTRAINT merchant_items_price CHECK (price >= 1),
    CONSTRAINT merchant_items_merchant_id_fk
        FOREIGN KEY (merchant_id)
        REFERENCES public.merchants(id)
        ON DELETE CASCADE
);

CREATE INDEX merchant_items_merchant_id_idx ON merchant_items(merchant_id);
CREATE INDEX merchant_items_category_idx ON merchant_items(category);