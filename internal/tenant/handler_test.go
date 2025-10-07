package tenant

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"flowcargo/db/testutils"
	ru "flowcargo/internal/shared/restutils"
)

func TestTenantHandler(t *testing.T) {
	l, err := testutils.NewTestLogger()
	require.NoError(t, err)
	service := newmockService()
	handler := NewTenantHandler(service, l)

	t.Run("Create tenant with valid input", func(t *testing.T) {
		w := httptest.NewRecorder()
		rawBody := CreateTenantParams{
			Name:   "New Tenant",
			Email:  "new_tenant@test.com",
			Domain: nil,
		}
		jsonBody, err := json.Marshal(rawBody)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/tenants", bytes.NewBuffer(jsonBody))
		handler.CreateTenant(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		var resp ru.APIResponse[Tenant]
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.True(t, resp.Success)
		require.Equal(t, "Tenant created successfully", resp.Message)
		require.Equal(t, rawBody.Name, resp.Data.Name)
		require.Equal(t, rawBody.Email, resp.Data.Email)
		require.NotNil(t, resp.Data.ID)
	})

	t.Run("Create tenant with existing email", func(t *testing.T) {
		w := httptest.NewRecorder()
		rawBody := CreateTenantParams{
			Name:   "Already Exists",
			Email:  "existing_tenant@test.com",
			Domain: nil,
		}
		jsonBody, err := json.Marshal(rawBody)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/tenants", bytes.NewBuffer(jsonBody))
		handler.CreateTenant(w, req)

		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		// TODO: Fix this test after updating error handling in service layer
		require.Equal(t, http.StatusInternalServerError, w.Code)
		var resp ru.APIResponse[Tenant]
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.False(t, resp.Success)
	})

	t.Run("Create tenant with invalid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		rawBody := map[string]string{
			"age":   "38",
			"email": "invalid_email_format",
		}
		jsonBody, err := json.Marshal(rawBody)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/tenants", bytes.NewBuffer(jsonBody))
		handler.CreateTenant(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		var resp ru.APIResponse[Tenant]
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
	})

	t.Run("Create tenant with correct format but invalid data", func(t *testing.T) {
		w := httptest.NewRecorder()

		body := CreateTenantParams{
			Name:  "a",
			Email: "b",
		}

		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/tenants", bytes.NewBuffer(jsonBody))
		handler.CreateTenant(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		var resp ru.APIResponse[Tenant]
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.False(t, resp.Success)
	})

	t.Run("Get tenant by ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		id := ("11111111-1111-1111-1111-111111111111")
		req := httptest.NewRequest(http.MethodGet, "/tenants/", nil)
		req.SetPathValue("id", id)
		handler.GetTenant(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp ru.APIResponse[Tenant]
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.True(t, resp.Success)
		require.Equal(t, "Tenant retrieved successfully", resp.Message)
		require.Equal(t, uuid.MustParse(id), resp.Data.ID)
	})

	t.Run("Get tenant with missing ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/tenants/", nil)
		handler.GetTenant(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		var resp ru.APIResponse[Tenant]
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.False(t, resp.Success)
		require.Nil(t, resp.Data)
		require.NotNil(t, resp.Error)
	})

	t.Run("Get tenant with that does not exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/tenants/", nil)
		req.SetPathValue("id", "22222222-2222-2222-2222-222222222222") // Invalid UUID
		handler.GetTenant(w, req)
		// TODO: Fix this test after updating error handling in service layer
		require.Equal(t, http.StatusInternalServerError, w.Code)
		var resp ru.APIResponse[Tenant]
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.False(t, resp.Success)
		require.Nil(t, resp.Data)
		require.NotNil(t, resp.Error)
	})

	t.Run("Update tenant with valid input", func(t *testing.T) {
		w := httptest.NewRecorder()
		id := "11111111-1111-1111-1111-111111111111"
		rawBody := UpdateTenantParams{
			ID:     uuid.MustParse(id),
			Name:   testutils.ToPointer("Updated Tenant"),
			Email:  testutils.ToPointer("new_email_updated@example.com"),
			Domain: testutils.ToPointer("http://updated.com"),
		}
		jsonBody, err := json.Marshal(rawBody)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPut, "/tenants/", bytes.NewBuffer(jsonBody))
		req.SetPathValue("id", id)
		handler.UpdateTenant(w, req)

		var resp ru.APIResponse[Tenant]
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		t.Log(resp)
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		require.Equal(t, resp.Data.Name, *rawBody.Name)
		require.Equal(t, resp.Data.Email, *rawBody.Email)
		require.Equal(t, resp.Data.Domain, rawBody.Domain)
		require.Equal(t, uuid.MustParse(id), resp.Data.ID)
		require.True(t, resp.Success)
		require.Equal(t, "Tenant updated successfully", resp.Message)
	})

	t.Run("Update tenant with invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		rawBody := UpdateTenantParams{
			Name:  testutils.ToPointer("Updated Tenant"),
			Email: testutils.ToPointer("updated_email@example.com"),
		}
		jsonBody, err := json.Marshal(rawBody)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPut, "/tenants/", bytes.NewBuffer(jsonBody))
		req.SetPathValue("id", "invalid-id")
		handler.UpdateTenant(w, req)

		var resp ru.APIResponse[Tenant]
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.False(t, resp.Success)
		require.Nil(t, resp.Data)
		require.NotNil(t, resp.Error)
	})

	t.Run("Delete a tenant with valid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		id := "11111111-1111-1111-1111-111111111111"

		req := httptest.NewRequest(http.MethodDelete, "/tenants/", nil)
		req.SetPathValue("id", id)
		handler.DeleteTenant(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		var resp ru.APIResponse[Tenant]
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.True(t, resp.Success)
		require.Equal(t, "Tenant deleted successfully", resp.Message)
		require.Equal(t, uuid.MustParse(id), resp.Data.ID)
		require.Equal(t, false, resp.Data.IsActive)
	})

	t.Run("Delete a tenant with invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodDelete, "/tenants/", nil)
		req.SetPathValue("id", "invalid-id")
		handler.DeleteTenant(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		var resp ru.APIResponse[Tenant]
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.False(t, resp.Success)
		require.Nil(t, resp.Data)
		require.NotNil(t, resp.Error)
	})
}