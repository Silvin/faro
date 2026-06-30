-- 0008_customers (down)
ALTER TABLE sales DROP COLUMN IF EXISTS customer_id;
DROP TABLE IF EXISTS customers;
