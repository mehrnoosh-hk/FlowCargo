-- Revert all changes made in the up migration

-- Revoke privileges from the app_admin role
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM app_admin;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM app_admin;
REVOKE USAGE ON SCHEMA public FROM app_admin;

-- Drop the app_admin role
DROP ROLE IF EXISTS app_admin;

-- Drop the pgaudit_role
DROP ROLE IF EXISTS pgaudit_role;

-- Drop the utility functions
DROP FUNCTION IF EXISTS is_app_admin();
DROP FUNCTION IF EXISTS current_tenant_id();
DROP FUNCTION IF EXISTS set_updated_at();

-- Disable extensions
DROP EXTENSION IF EXISTS "citext";
DROP EXTENSION IF EXISTS "uuid-ossp";

-- Disable pgaudit extension
DROP EXTENSION IF EXISTS "pgaudit";