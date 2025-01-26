-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Function to create enum if it doesn't exist
CREATE OR REPLACE FUNCTION create_enum_if_not_exists(enum_name TEXT, enum_values TEXT[])
RETURNS void AS $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = enum_name) THEN
        EXECUTE 'CREATE TYPE ' || enum_name || ' AS ENUM (' ||
            array_to_string(enum_values, ', ') || ')';
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Create enums if they don't exist
SELECT create_enum_if_not_exists('user_role', 
    ARRAY['''ADMIN''', '''STAFF''']
);

SELECT create_enum_if_not_exists('space_status', 
    ARRAY['''Occupied (Paid)''', '''Occupied (Payment Due)''', '''Occupied (Overdue)''', 
          '''Vacant''', '''Reserved''']
);

SELECT create_enum_if_not_exists('payment_type', 
    ARRAY['''Monthly''', '''Weekly''', '''Daily''']
);

SELECT create_enum_if_not_exists('payment_status', 
    ARRAY['''Paid''', '''Due''', '''Overdue''']
);

-- Create tokens table if it doesn't exist
CREATE TABLE IF NOT EXISTS tokens (
    token_hash TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    revoked BOOLEAN DEFAULT false,
    CONSTRAINT token_expiry_valid CHECK (expires_at > created_at)
);

-- Create users table if it doesn't exist
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    role user_role NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP WITH TIME ZONE,
    CONSTRAINT email_valid CHECK (email ~* '^[A-Za-z0-9._+%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$')
);

-- Create sections table if it doesn't exist
CREATE TABLE IF NOT EXISTS sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL
);

-- Create spaces table if it doesn't exist
CREATE TABLE IF NOT EXISTS spaces (
    id VARCHAR(20) PRIMARY KEY,
    section_id UUID NOT NULL,
    status space_status NOT NULL DEFAULT 'Vacant',
    tenant_id UUID,
    reserved BOOLEAN NOT NULL DEFAULT false,
    payment_type payment_type,
    next_payment TIMESTAMP WITH TIME ZONE,
    past_due_amount DECIMAL(10,2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create tenants table if it doesn't exist
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    move_in_date TIMESTAMP WITH TIME ZONE NOT NULL,
    space_id VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create payments table if it doesn't exist
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    due_date TIMESTAMP WITH TIME ZONE NOT NULL,
    paid_date TIMESTAMP WITH TIME ZONE,
    previous_payment_date TIMESTAMP WITH TIME ZONE,
    payment_type payment_type NOT NULL,
    status payment_status NOT NULL DEFAULT 'Due',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- Create indexes if they don't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_tokens_user_id') THEN
        CREATE INDEX idx_tokens_user_id ON tokens(user_id);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_tokens_expires_at') THEN
        CREATE INDEX idx_tokens_expires_at ON tokens(expires_at);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_spaces_section_id') THEN
        CREATE INDEX idx_spaces_section_id ON spaces(section_id);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_spaces_tenant_id') THEN
        CREATE INDEX idx_spaces_tenant_id ON spaces(tenant_id);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_tenants_space_id') THEN
        CREATE INDEX idx_tenants_space_id ON tenants(space_id);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_payments_tenant_id') THEN
        CREATE INDEX idx_payments_tenant_id ON payments(tenant_id);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_payments_due_date') THEN
        CREATE INDEX idx_payments_due_date ON payments(due_date);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_payments_status') THEN
        CREATE INDEX idx_payments_status ON payments(status);
    END IF;
END $$;

-- Create or replace the updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers if they don't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_spaces_updated_at') THEN
        CREATE TRIGGER update_spaces_updated_at
            BEFORE UPDATE ON spaces
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_tenants_updated_at') THEN
        CREATE TRIGGER update_tenants_updated_at
            BEFORE UPDATE ON tenants
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_payments_updated_at') THEN
        CREATE TRIGGER update_payments_updated_at
            BEFORE UPDATE ON payments
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

-- Add foreign key constraints if they don't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_spaces_section_id'
    ) THEN
        ALTER TABLE spaces
        ADD CONSTRAINT fk_spaces_section_id
        FOREIGN KEY (section_id)
        REFERENCES sections(id)
        ON DELETE CASCADE;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_spaces_tenant'
    ) THEN
        ALTER TABLE spaces
        ADD CONSTRAINT fk_spaces_tenant
        FOREIGN KEY (tenant_id)
        REFERENCES tenants(id)
        ON DELETE SET NULL;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_tenants_space'
    ) THEN
        ALTER TABLE tenants
        ADD CONSTRAINT fk_tenants_space
        FOREIGN KEY (space_id)
        REFERENCES spaces(id)
        ON DELETE SET NULL;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_payments_tenant'
    ) THEN
        ALTER TABLE payments
        ADD CONSTRAINT fk_payments_tenant
        FOREIGN KEY (tenant_id)
        REFERENCES tenants(id)
        ON DELETE CASCADE;
    END IF;
END $$;

-- Insert default sections if they don't exist
INSERT INTO sections (name)
SELECT 'Mane Street' WHERE NOT EXISTS (SELECT 1 FROM sections WHERE name = 'Mane Street');

INSERT INTO sections (name)
SELECT 'Grace Street' WHERE NOT EXISTS (SELECT 1 FROM sections WHERE name = 'Grace Street');

INSERT INTO sections (name)
SELECT 'Trae Street' WHERE NOT EXISTS (SELECT 1 FROM sections WHERE name = 'Trae Street');

INSERT INTO sections (name)
SELECT 'Summer Street' WHERE NOT EXISTS (SELECT 1 FROM sections WHERE name = 'Summer Street');

INSERT INTO sections (name)
SELECT 'Rock Street' WHERE NOT EXISTS (SELECT 1 FROM sections WHERE name = 'Rock Street');

INSERT INTO sections (name)
SELECT 'Cedar Street' WHERE NOT EXISTS (SELECT 1 FROM sections WHERE name = 'Cedar Street');


-- Create space initialization function if it doesn't exist
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
        SELECT 
            prefix || i::TEXT,
            section_id,
            'Vacant'
        WHERE NOT EXISTS (
            SELECT 1 FROM spaces WHERE id = prefix || i::TEXT
        );
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Initialize spaces if they don't exist
SELECT initialize_section_spaces('Mane Street', 'M', 24);
SELECT initialize_section_spaces('Grace Street', 'G', 32);
SELECT initialize_section_spaces('Trae Street', 'T', 32);
SELECT initialize_section_spaces('Summer Street', 'S', 34);
SELECT initialize_section_spaces('Rock Street', 'R', 18);
SELECT initialize_section_spaces('Cedar Street', 'C', 13);