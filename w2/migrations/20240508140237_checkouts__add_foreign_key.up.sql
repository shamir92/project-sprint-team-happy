ALTER TABLE public.checkouts
    ADD CONSTRAINT checkouts_fk_customer_id 
        FOREIGN KEY (customer_id) 
        REFERENCES public.customers(id)
        ON DELETE CASCADE;

ALTER TABLE public.checkouts
    ADD CONSTRAINT checkouts_fk_user_id 
        FOREIGN KEY (created_by) 
        REFERENCES public.users(user_id)
        ON DELETE CASCADE;