-- 0004_sales (up): ventas del punto de venta (módulo M4).
-- Ver: .arete/modules/pos/tech-spec.md

CREATE TABLE sales (
    id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         uuid NOT NULL REFERENCES tenants(id),
    total_cents       integer NOT NULL,
    amount_paid_cents integer NOT NULL,
    change_cents      integer NOT NULL,
    created_at        timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE sale_items (
    id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    sale_id          uuid NOT NULL REFERENCES sales(id) ON DELETE CASCADE,
    product_id       uuid REFERENCES products(id),  -- snapshot de name/price abajo
    name             text NOT NULL,
    unit_price_cents integer NOT NULL,
    quantity         integer NOT NULL CHECK (quantity > 0),
    line_total_cents integer NOT NULL
);

CREATE INDEX sales_tenant_created_idx ON sales (tenant_id, created_at DESC);
CREATE INDEX sale_items_sale_id_idx ON sale_items (sale_id);
