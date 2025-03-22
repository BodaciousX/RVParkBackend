-- to access logs: psql "postgres://dbuser:dbpassword@localhost:5433/RVParkDB"
-- to quit logs: \q

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Drop existing tables if they exist (in reverse order of dependencies)
DROP TABLE IF EXISTS payments CASCADE;
DROP TABLE IF EXISTS tenants CASCADE;
DROP TABLE IF EXISTS spaces CASCADE;
DROP TABLE IF EXISTS sections CASCADE;
DROP TABLE IF EXISTS tokens CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop existing types if they exist
DROP TYPE IF EXISTS payment_method CASCADE;
DROP TYPE IF EXISTS space_status CASCADE;
DROP TYPE IF EXISTS user_role CASCADE;

-- Create enums
CREATE TYPE space_status AS ENUM ('Occupied', 'Vacant', 'Reserved');
CREATE TYPE payment_method AS ENUM ('CREDIT', 'CHECK', 'CASH');

-- Create tokens table
CREATE TABLE tokens (
    token_hash TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT LOCALTIMESTAMP,
    revoked BOOLEAN DEFAULT false,
    CONSTRAINT token_expiry_valid CHECK (expires_at > created_at)
);

-- Create users table (without role column)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT LOCALTIMESTAMP,
    last_login TIMESTAMP,
    CONSTRAINT email_valid CHECK (email ~* '^[A-Za-z0-9._+%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$')
);

-- Create sections table
CREATE TABLE sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL
);

-- Create spaces table
CREATE TABLE spaces (
    id VARCHAR(20) PRIMARY KEY,
    section_id UUID NOT NULL,
    status space_status NOT NULL DEFAULT 'Vacant',
    tenant_id UUID,
    reserved BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT LOCALTIMESTAMP,
    updated_at TIMESTAMP DEFAULT LOCALTIMESTAMP
);

-- Create tenants table
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    move_in_date TIMESTAMP NOT NULL,
    space_id VARCHAR(20),
    created_at TIMESTAMP DEFAULT LOCALTIMESTAMP,
    updated_at TIMESTAMP DEFAULT LOCALTIMESTAMP
);

-- Create payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    amount_due DECIMAL(10,2) NOT NULL,
    due_date TIMESTAMP NOT NULL,
    paid_date TIMESTAMP,
    next_payment_date TIMESTAMP NOT NULL,
    payment_method payment_method,
    created_at TIMESTAMP DEFAULT LOCALTIMESTAMP,
    updated_at TIMESTAMP DEFAULT LOCALTIMESTAMP,
    CONSTRAINT amount_positive CHECK (amount_due > 0)
);

-- Create indexes
CREATE INDEX idx_tokens_user_id ON tokens(user_id);
CREATE INDEX idx_tokens_expires_at ON tokens(expires_at);
CREATE INDEX idx_spaces_section_id ON spaces(section_id);
CREATE INDEX idx_spaces_tenant_id ON spaces(tenant_id);
CREATE INDEX idx_tenants_space_id ON tenants(space_id);
CREATE INDEX idx_payments_tenant_id ON payments(tenant_id);
CREATE INDEX idx_payments_due_date ON payments(due_date);
CREATE INDEX idx_payments_next_payment_date ON payments(next_payment_date);
CREATE INDEX idx_payments_payment_method ON payments(payment_method);

-- Create the updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = LOCALTIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers
CREATE TRIGGER update_spaces_updated_at
    BEFORE UPDATE ON spaces
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tenants_updated_at
    BEFORE UPDATE ON tenants
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payments_updated_at
    BEFORE UPDATE ON payments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add foreign key constraints
ALTER TABLE spaces
ADD CONSTRAINT fk_spaces_section_id
FOREIGN KEY (section_id)
REFERENCES sections(id)
ON DELETE CASCADE;

ALTER TABLE spaces
ADD CONSTRAINT fk_spaces_tenant
FOREIGN KEY (tenant_id)
REFERENCES tenants(id)
ON DELETE SET NULL;

ALTER TABLE tenants
ADD CONSTRAINT fk_tenants_space
FOREIGN KEY (space_id)
REFERENCES spaces(id)
ON DELETE SET NULL;

ALTER TABLE payments
ADD CONSTRAINT fk_payments_tenant
FOREIGN KEY (tenant_id)
REFERENCES tenants(id)
ON DELETE CASCADE;

-- Insert default sections
INSERT INTO sections (name) VALUES ('Mane Street');
INSERT INTO sections (name) VALUES ('Grace Street');
INSERT INTO sections (name) VALUES ('Trae Street');
INSERT INTO sections (name) VALUES ('Summer Street');
INSERT INTO sections (name) VALUES ('Rock Street');
INSERT INTO sections (name) VALUES ('Cedar Street');

-- Create space initialization function
CREATE OR REPLACE FUNCTION initialize_section_spaces(
    section_name VARCHAR,
    prefix CHAR,
    space_count INTEGER
)
RETURNS VOID AS $$
DECLARE
    section_id UUID;
    i INTEGER;
BEGIN
    SELECT id INTO section_id FROM sections WHERE name = section_name;
    
    FOR i IN 1..space_count LOOP
        INSERT INTO spaces (id, section_id, status)
        VALUES (
            prefix || i::TEXT,
            section_id,
            'Vacant'
        );
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Initialize spaces
SELECT initialize_section_spaces('Mane Street', 'M', 24);
SELECT initialize_section_spaces('Grace Street', 'G', 32);
SELECT initialize_section_spaces('Trae Street', 'T', 32);
SELECT initialize_section_spaces('Summer Street', 'S', 34);
SELECT initialize_section_spaces('Rock Street', 'R', 18);
SELECT initialize_section_spaces('Cedar Street', 'C', 13);