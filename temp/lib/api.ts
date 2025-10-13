import { HttpClient } from "@/services/tenant.service";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
const API_TIMEOUT = Number(process.env.NEXT_PUBLIC_API_TIMEOUT) || 10000;

export class ApiError extends Error {
  constructor(
    public status: number,
    public statusText: string,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

const applyRequestInterceptor = (
  url: string,
  options: RequestInit,
): RequestInit => {
  // TODO: Implement auth token
  // For now just loging in console
  if (process.env.NODE_ENV === "development") {
    console.log(`[API Request] ${options.method || "GET"} ${url}`);
  }
  return options;
};

const handleResponse = async (response: Response): Promise<Response> => {
  if (process.env.NODE_ENV === "development") {
    console.log(`[API Response] ${response.status} ${response.statusText}`);
  }

  if (!response.ok) {
    let errorMessage = response.statusText;
    try {
      const errData = await response.json();
      errorMessage = errData || errorMessage;
    } catch {
      // Just return errorMessage
    }

    switch (response.status) {
      case 401:
        console.warn("[API] Unauthorized - please log in");
        break;
      case 403:
        console.warn("[API] Forbidden - insufficient permissions");
        break;
      case 404:
        console.warn("[API] Resource not found");
        break;
      case 500:
        console.error("[API] Internal server error");
        break;
    }
    throw new ApiError(response.status, response.statusText, errorMessage);
  }
  return response;
};

const fetchWithTimeout = async (
  url: string,
  options: RequestInit = {},
): Promise<Response> => {
  const controller: AbortController = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), API_TIMEOUT);
  try {
    const fullUrl: string = `${API_BASE_URL}${url}`;
    const requestOptions: RequestInit = applyRequestInterceptor(fullUrl, {
      ...options,
      signal: controller.signal,
      headers: {
        "Content-Type": "application/json",
      },
    });
    const response: Response = await fetch(fullUrl, requestOptions);
    return handleResponse(response);
  } catch (error) {
    if (error instanceof Error && error.name === "AbortError") {
      throw new Error("Request timeout");
    }
    throw error;
  } finally {
    clearTimeout(timeoutId);
  }
};

/**
 * API Client implementation that satisfies the HttpClient interface
 * This is the concrete implementation that tenant service will use
 */
export const apiClient: HttpClient = {
  get: async <T>(url: string): Promise<T> => {
    const response: Response = await fetchWithTimeout(url, { method: "GET" });
    return response.json();
  },

  post: async <T>(url: string, data?: any): Promise<T> => {
    const response: Response = await fetchWithTimeout(url, {
      method: "POST",
      body: data ? JSON.stringify(data) : undefined,
    });
    return response.json();
  },

  put: async <T>(url: string, data?: any): Promise<T> => {
    const response: Response = await fetchWithTimeout(url, {
      method: "PUT",
      body: data ? JSON.stringify(data) : undefined,
    });
    return response.json();
  },

  delete: async <T>(url: string): Promise<T> => {
    const response: Response = await fetchWithTimeout(url, {
      method: "DELETE",
    });
    return response.json();
  },
};

export const isApiError = (err: Error): boolean => {
  return err instanceof ApiError;
};

export const getErrorMessage = (error: Error): string => {
  if (isApiError(error)) {
    return error.message;
  }
  if (error instanceof Error) {
    return error.message;
  }
  return "Unknown error";
};
