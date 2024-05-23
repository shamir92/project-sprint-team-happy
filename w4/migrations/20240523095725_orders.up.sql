CREATE TABLE IF NOT EXISTS public.orders(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    user_lon float not null,
    user_lat float not null,
    total_price int not null,
    estimated_delivery_time smallint not null,
    state char varying(10) not null DEFAULT 'estimated',
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT orders_pkey PRIMARY KEY (id),
    CONSTRAINT orders_user_id_fk
        FOREIGN KEY (user_id)
        REFERENCES public.users(id)
        ON DELETE CASCADE
);

CREATE INDEX orders_user_id_idx ON orders(user_id);
CREATE INDEX orders_state_idx ON orders(state);