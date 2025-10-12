import { UUID } from "crypto";

export interface Tenant {
  id: UUID;
  name: string;
  email: string;
  domain: string | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateTenantRequest {
  name: string;
  email: string;
  domain?: string | null;
}

export interface UpdateTenantRequest {
  name?: string;
  email?: string;
  domain?: string | null;
  is_active?: boolean;
}

export interface TenantResponse {
  tenant: Tenant;
}

export interface TenantsResponse {
  tenants: Tenant[];
  total?: number;
}
