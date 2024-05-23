CREATE TABLE IF NOT EXISTS public.merchants(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name char varying(30) NOT NULL,
    category char varying(50) NOT NULL,
    image_url text not null,
    lon float not null,
    lat float not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,

    CONSTRAINT merchants_pkey PRIMARY KEY (id)
);

CREATE INDEX merchants_category_index ON merchants (category);