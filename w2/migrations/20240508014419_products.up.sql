CREATE TABLE IF NOT EXISTS public.products (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name character varying(30) NOT NULL,
    sku character varying(30) NOT NULL,
    category character varying(20) NOT NULL,
    image_url character varying(2048) NOT NULL,
    notes character varying(200) NOT NULL,
    price int NOT NULL,
    stock smallint NOT NULL,
    location character varying(200) NOT NULL,
    is_available boolean not null,
    created_by uuid NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone DEFAULT NULL,

    CONSTRAINT products_pkey PRIMARY KEY (id),
    CONSTRAINT products_stock CHECK (stock >= 0 AND stock <= 10000),
    CONSTRAINT products_price CHECK (price >= 1)
);