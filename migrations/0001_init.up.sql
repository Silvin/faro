-- 0001_init (up): módulo auth/login — negocios (tenants) y usuarios.
-- Ver: .arete/modules/login/tech-spec.md y ADR-001 (multi-tenancy).

CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto; -- gen_random_uuid()

CREATE TABLE tenants (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name       text NOT NULL,
    status     text NOT NULL DEFAULT 'active', -- active | suspended
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE users (
    id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      uuid REFERENCES tenants(id),       -- NULL solo para super admin global
    email          citext NOT NULL UNIQUE,            -- único global
    password_hash  text NOT NULL,                     -- bcrypt
    name           text NOT NULL,
    is_super_admin boolean NOT NULL DEFAULT false,
    status         text NOT NULL DEFAULT 'active',    -- active | disabled
    created_at     timestamptz NOT NULL DEFAULT now(),
    -- Un usuario de negocio siempre pertenece a un tenant; solo el super admin global puede no tenerlo.
    CONSTRAINT users_tenant_required CHECK (is_super_admin = true OR tenant_id IS NOT NULL)
);

CREATE INDEX users_tenant_id_idx ON users (tenant_id);
