BEGIN;

-- Create tenants table for multi-tenant architecture
-- This table stores information about each tenant in the system

CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email CITEXT NOT NULL,
    domain VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Add case-insensitive unique constraint on email
ALTER TABLE tenants ADD CONSTRAINT tenants_email_unique UNIQUE (email);

-- Create indexes for better performance
CREATE INDEX idx_tenants_email ON tenants (email);
CREATE INDEX idx_tenants_is_active ON tenants (is_active);
CREATE INDEX idx_tenants_domain ON tenants (domain) WHERE domain IS NOT NULL;
CREATE INDEX idx_tenants_created_at ON tenants (created_at);

-- Add update trigger for updated_at timestamp
CREATE TRIGGER set_tenants_updated_at
    BEFORE UPDATE ON tenants
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- Grant permissions to service_account role (data manipulation only)
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE tenants TO service_account;

-- Add comments for documentation
COMMENT ON TABLE tenants IS 'Stores tenant information for multi-tenant architecture';
COMMENT ON COLUMN tenants.id IS 'Unique identifier for the tenant';
COMMENT ON COLUMN tenants.name IS 'Display name of the tenant';
COMMENT ON COLUMN tenants.email IS 'Contact email for the tenant (case-insensitive, unique)';
COMMENT ON COLUMN tenants.domain IS 'Optional domain for subdomain-based multi-tenancy';
COMMENT ON COLUMN tenants.is_active IS 'Whether the tenant is active and can access the system';
COMMENT ON COLUMN tenants.created_at IS 'Timestamp when the tenant was created';
COMMENT ON COLUMN tenants.updated_at IS 'Timestamp when the tenant was last updated';

COMMIT;