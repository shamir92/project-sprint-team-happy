ALTER TABLE public.products
    ADD CONSTRAINT products_fk_created_by
        FOREIGN KEY (created_by) 
        REFERENCES public.users(user_id)
        ON DELETE CASCADE;