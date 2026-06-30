-- 0006_sale_payment_method (down)
ALTER TABLE sales DROP COLUMN IF EXISTS payment_method;
