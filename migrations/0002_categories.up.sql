-- 0002_categories (up): categorías de producto por negocio (módulo M2).
-- Ver: .arete/modules/category-management/tech-spec.md

CREATE TABLE categories (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  uuid NOT NULL REFERENCES tenants(id),
    name       text NOT NULL,
    status     text NOT NULL DEFAULT 'active', -- active | inactive
    sort_order integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT categories_name_per_tenant UNIQUE (tenant_id, name)
);

CREATE INDEX categories_tenant_id_idx ON categories (tenant_id);
