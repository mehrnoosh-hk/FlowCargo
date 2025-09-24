-- Create all functions, types, and extensions for the initial database setup.

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable case-insensitive text
CREATE EXTENSION IF NOT EXISTS "citext";



-- Function to automatically update the updated_at timestamp
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create function to get current tenant ID from session, raising a clear error if not set.
CREATE OR REPLACE FUNCTION current_tenant_id() RETURNS UUID AS $$
DECLARE
    tenant_id_text TEXT;
BEGIN
    tenant_id_text := current_setting('app.tenant_id', true);

    IF tenant_id_text IS NULL OR tenant_id_text = '' THEN
        RAISE EXCEPTION 'app.tenant_id is not set. All queries on tenant-scoped tables must be wrapped in With Tenant Context.';
    END IF;

    RETURN tenant_id_text::UUID;
EXCEPTION
    -- Catch potential invalid UUID format
    WHEN invalid_text_representation THEN
        RAISE EXCEPTION 'Invalid UUID format for app.tenant_id: "%"', tenant_id_text;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;


-- Function to check if the current user is a service_account
CREATE OR REPLACE FUNCTION is_service_account() RETURNS BOOLEAN AS $$
DECLARE
    setting_value TEXT;
BEGIN
    setting_value := current_setting('app.is_service_account', true);

    -- Safely check the value. If it's not 'true', it's false.
    RETURN setting_value IS NOT NULL AND setting_value = 'true';
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Create the service_account role
DO
$$
BEGIN
    IF NOT EXISTS (
       SELECT FROM pg_catalog.pg_roles
       WHERE  rolname = 'service_account') THEN

       CREATE ROLE service_account;
    END IF;
END
$$;

-- Grant privileges to the service_account role (data manipulation only)
GRANT USAGE ON SCHEMA public TO service_account;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO service_account;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO service_account;


