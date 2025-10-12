import {
  type Tenant,
  type CreateTenantRequest,
  type UpdateTenantRequest,
} from "@/types/tenant";

/**
 * HTTP Client interface - defines what the tenant service NEEDS
 * This is the abstraction that any HTTP implementation must satisfy
 */
export interface HttpClient {
  get<T>(url: string): Promise<T>;
  post<T>(url: string, data?: any): Promise<T>;
  put<T>(url: string, data?: any): Promise<T>;
  delete<T>(url: string): Promise<T>;
}

const TENANT_BASE_URL = "/tenants";

/**
 * Factory function that creates a tenant service with dependency injection
 * @param client - Any HTTP client that implements the HttpClient interface
 */
export const createTenantService = (client: HttpClient) => ({
  /**
   * Get a tenant by ID
   * GET /tenants/{id}
   */
  getById: async (id: string): Promise<Tenant> => {
    return client.get<Tenant>(`${TENANT_BASE_URL}/${id}`);
  },

  /**
   * Create a new tenant
   * POST /tenants
   */
  create: async (data: CreateTenantRequest): Promise<Tenant> => {
    return client.post<Tenant>(TENANT_BASE_URL, data);
  },

  /**
   * Update an existing tenant
   * PUT /tenants/{id}
   */
  update: async (id: string, data: UpdateTenantRequest): Promise<Tenant> => {
    return client.put<Tenant>(`${TENANT_BASE_URL}/${id}`, data);
  },

  /**
   * Delete a tenant
   * DELETE /tenants/{id}
   */
  delete: async (id: string): Promise<Tenant> => {
    return client.delete<Tenant>(`${TENANT_BASE_URL}/${id}`);
  },
});

/**
 * Helper function to format tenant data for display in UI
 */
export const formatTenantForDisplay = (tenant: Tenant) => ({
  ...tenant,
  createdAt: new Date(tenant.created_at).toLocaleDateString(),
  updatedAt: new Date(tenant.updated_at).toLocaleDateString(),
  status: tenant.is_active ? "Active" : "Inactive",
});

/**
 * Type of the tenant service (useful for typing in components)
 */
export type TenantService = ReturnType<typeof createTenantService>;
