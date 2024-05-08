ALTER TABLE public.checkouts
    DROP CONSTRAINT IF EXISTS checkouts_fk_customer_id;

ALTER TABLE public.checkouts
    DROP CONSTRAINT IF EXISTS checkouts_fk_user_id;