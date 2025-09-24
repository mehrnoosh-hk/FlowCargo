-- Revert all changes made in the up migration

-- Revoke privileges from the service_account role
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM service_account;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM service_account;
REVOKE USAGE ON SCHEMA public FROM service_account;

-- Drop the service_account role
DROP ROLE IF EXISTS service_account;



-- Drop the utility functions
DROP FUNCTION IF EXISTS is_service_account();
DROP FUNCTION IF EXISTS current_tenant_id();
DROP FUNCTION IF EXISTS set_updated_at();

-- Disable extensions
DROP EXTENSION IF EXISTS "citext";
DROP EXTENSION IF EXISTS "uuid-ossp";
