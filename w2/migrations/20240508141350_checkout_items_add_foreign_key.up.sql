ALTER TABLE checkout_items
    ADD CONSTRAINT checkout_items_fk_checkout_id 
    FOREIGN KEY (checkout_id) 
    REFERENCES public.checkouts(id)
    ON DELETE CASCADE;

ALTER TABLE checkout_items
    ADD CONSTRAINT checkout_items_fk_product_id 
    FOREIGN KEY (product_id) 
    REFERENCES public.products(id)
    ON DELETE CASCADE;