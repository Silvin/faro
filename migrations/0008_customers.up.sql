-- 0008_customers (up): clientes del programa de lealtad y su asociación a ventas.
CREATE TABLE customers (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  uuid NOT NULL REFERENCES tenants(id),
    phone      text NOT NULL,
    first_name text NOT NULL,
    last_name  text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT customers_phone_per_tenant UNIQUE (tenant_id, phone)
);

CREATE INDEX customers_tenant_id_idx ON customers (tenant_id);

-- Asociación opcional de la venta a un cliente.
ALTER TABLE sales ADD COLUMN customer_id uuid REFERENCES customers(id);
