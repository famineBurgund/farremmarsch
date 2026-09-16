CREATE SCHEMA IF NOT EXISTS farrepmach;

CREATE TYPE farrepmach.user_role AS ENUM ('tenant', 'employee', 'admin');

CREATE TABLE farrepmach.users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    phone TEXT,
    password_hash TEXT NOT NULL,
    full_name TEXT NOT NULL,
    role farrepmach.user_role NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE farrepmach.tenant_organizations (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    inn TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE farrepmach.tenant_users (
    user_id BIGINT NOT NULL REFERENCES farrepmach.users(id) ON DELETE CASCADE,
    organization_id BIGINT NOT NULL REFERENCES farrepmach.tenant_organizations(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, organization_id)
);

CREATE TYPE farrepmach.employee_position AS ENUM (
    'chief_accountant',
    'chief_power_engineer',
    'leasing_manager'
);

CREATE TABLE farrepmach.employees (
    user_id BIGINT PRIMARY KEY REFERENCES farrepmach.users(id) ON DELETE CASCADE,
    position farrepmach.employee_position NOT NULL,
    bitrix_responsible_id TEXT
);
