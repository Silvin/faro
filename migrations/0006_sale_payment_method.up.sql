-- 0006_sale_payment_method (up): forma de pago de la venta (efectivo | tarjeta).
ALTER TABLE sales ADD COLUMN payment_method text NOT NULL DEFAULT 'cash';
