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
        RAISE EXCEPTION 'app.tenant_id is not set.';
    END IF;

    RETURN tenant_id_text::UUID;
EXCEPTION
    -- Catch potential invalid UUID format
    WHEN invalid_text_representation THEN
        RAISE EXCEPTION 'Invalid UUID format for app.tenant_id: "%''', tenant_id_text;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;


-- Function to check if the current user is an app_admin
CREATE OR REPLACE FUNCTION is_app_admin()
RETURNS BOOLEAN AS $$
BEGIN
    RETURN current_user = '''app_admin''';
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Create the app_admin role
DO
$$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE  rolname = '''app_admin''') THEN

      CREATE ROLE '''app_admin''';
   END IF;
END
$$;

-- Grant privileges to the app_admin role
GRANT USAGE ON SCHEMA public TO '''app_admin''';
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO '''app_admin''';
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO '''app_admin''';