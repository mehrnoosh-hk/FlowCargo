import { apiClient } from "@/lib/api";
import { createTenantService } from "./tenant.service";

/**
 * Default tenant service instance using the apiClient
 *
 * Usage in components:
 * ```ts
 * import { tenantService } from "@/services/tenant";
 *
 * const tenant = await tenantService.getById("123");
 * ```
 *
 * For testing, use createTenantService() with a mock:
 * ```ts
 * import { createTenantService } from "@/services/tenant.service";
 *
 * const mockClient = { get: vi.fn(), post: vi.fn(), ... };
 * const service = createTenantService(mockClient);
 * ```
 */
export const tenantService = createTenantService(apiClient);

// Re-export for convenience
export { createTenantService, formatTenantForDisplay } from "./tenant.service";
export type { HttpClient, TenantService } from "./tenant.service";
