-- 0003_products (up): catálogo de productos por negocio (módulo M3).
-- Ver: .arete/modules/product-management/tech-spec.md

CREATE TABLE products (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id),
    category_id uuid REFERENCES categories(id),       -- opcional; del mismo negocio
    name        text NOT NULL,
    price_cents integer NOT NULL,                      -- dinero como entero (centavos)
    status      text NOT NULL DEFAULT 'active',        -- active | inactive
    created_at  timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT products_name_per_tenant UNIQUE (tenant_id, name),
    CONSTRAINT products_price_positive CHECK (price_cents > 0)
);

CREATE INDEX products_tenant_id_idx ON products (tenant_id);
CREATE INDEX products_category_id_idx ON products (category_id);
