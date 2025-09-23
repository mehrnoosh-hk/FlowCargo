-- Revert all changes made in the up migration

-- Drop the app_admin role
DROP ROLE IF EXISTS '''app_admin''';

-- Drop the utility functions
DROP FUNCTION IF EXISTS is_app_admin();
DROP FUNCTION IF EXISTS current_tenant_id();
DROP FUNCTION IF EXISTS set_updated_at();

-- Disable extensions
DROP EXTENSION IF EXISTS "citext";
DROP EXTENSION IF EXISTS "uuid-ossp";
