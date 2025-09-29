BEGIN;

-- Drop tenants table and all associated objects
-- This migration reverses the tenants table creation

-- Drop the table (this will also drop associated indexes and constraints)
DROP TABLE IF EXISTS tenants;

-- Note: The update trigger is automatically dropped when the table is dropped
-- The service_account role permissions are not explicitly revoked here as they may be needed for other tables

COMMIT;
