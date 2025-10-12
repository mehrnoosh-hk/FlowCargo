package tenant

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"flowcargo/internal/shared/logger"
	ru "flowcargo/internal/shared/restutils"
)

type TenantHandler interface {
	CreateTenant(w http.ResponseWriter, r *http.Request)
	GetTenant(w http.ResponseWriter, r *http.Request)
	UpdateTenant(w http.ResponseWriter, r *http.Request)
	DeleteTenant(w http.ResponseWriter, r *http.Request)
}

type tenantHandler struct {
	service TenantService
	l       logger.Logger
}

func NewTenantHandler(service TenantService, l logger.Logger) TenantHandler {
	return tenantHandler{
		service: service,
		l:       l,
	}
}

var validate = validator.New()

// CreateTenant godoc
// @Summary      Create a new tenant
// @Description  Create a new tenant with the provided information
// @Tags         tenants
// @Accept       json
// @Produce      json
// @Param        tenant  body      CreateTenantParams  true  "Tenant creation parameters"
// @Success      201     {object}  ru.APIResponse[Tenant]
// @Failure      400     {object}  ru.APIResponse[any]
// @Failure      500     {object}  ru.APIResponse[any]
// @Router       /tenants [post]
func (h tenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ru.HandleMethodNotAllowed(w, r.Method, r.URL, ru.ResourceTenant, h.l)
	}
	var ctp CreateTenantParams
	decoderErr := json.NewDecoder(r.Body).Decode(&ctp)
	if decoderErr != nil {
		ru.HandleBadRequest(w, decoderErr, h.l)
		return
	}

	validateErr := validate.Struct(ctp)
	if validateErr != nil {
		ru.HandleBadRequest(w, validateErr, h.l)
		return
	}

	tenant, err := h.service.CreateTenant(r.Context(), ctp)
	if err != nil {
		ru.HandleInternalServerError(w, err, ru.ResourceTenant, h.l)
		return
	}
	ru.WriteSuccessResponse(
		w,
		http.StatusCreated,
		tenant,
		"Tenant created successfully",
	)
}

// GetTenant godoc
// @Summary      Get a tenant by ID
// @Description  Retrieve a tenant's information by their unique identifier
// @Tags         tenants
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Tenant ID (UUID format)"
// @Success      200  {object}  ru.APIResponse[Tenant]
// @Failure      400  {object}  ru.APIResponse[any]
// @Failure      500  {object}  ru.APIResponse[any]
// @Router       /tenants/{id} [get]
func (h tenantHandler) GetTenant(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	if idString == "" {
		ru.HandleBadRequest(w, errors.New("missing id parameter"), h.l)
		return
	}
	id, err := uuid.Parse(idString)
	if err != nil {
		ru.HandleBadRequest(w, errors.New("invalid id parameter"), h.l)
		return
	}
	tenant, err := h.service.GetTenantByID(r.Context(), id)
	if err != nil {
		ru.HandleInternalServerError(w, err, ru.ResourceTenant, h.l)
		return
	}
	ru.WriteSuccessResponse(w, http.StatusOK, tenant, "Tenant retrieved successfully")
}

// UpdateTenant godoc
// @Summary      Update a tenant
// @Description  Update an existing tenant's information by their unique identifier
// @Tags         tenants
// @Accept       json
// @Produce      json
// @Param        id      path      string              true  "Tenant ID (UUID format)"
// @Param        tenant  body      UpdateTenantParams  true  "Tenant update parameters"
// @Success      200     {object}  ru.APIResponse[Tenant]
// @Failure      400     {object}  ru.APIResponse[any]
// @Failure      500     {object}  ru.APIResponse[any]
// @Router       /tenants/{id} [put]
func (h tenantHandler) UpdateTenant(w http.ResponseWriter, r *http.Request) {

	id, err := ru.RetrieveID(r)

	if err != nil {
		ru.HandleBadRequest(w, err, h.l)
		return
	}
	// Note: The request body should contain the tenant ID.
	// Improve: Get the ID from the URL path instead.
	utp, err := ru.RetrieveBody[UpdateTenantParams](r)
	if err != nil {
		ru.HandleBadRequest(w, err, h.l)
		return
	}

	tenant, err := h.service.UpdateTenant(r.Context(), id, utp)
	if err != nil {
		ru.HandleInternalServerError(w, err, ru.ResourceTenant, h.l)
		return
	}

	ru.WriteSuccessResponse(w, http.StatusOK, tenant, "Tenant updated successfully")
}

// DeleteTenant godoc
// @Summary      Delete a tenant
// @Description  Delete an existing tenant by their unique identifier
// @Tags         tenants
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Tenant ID (UUID format)"
// @Success      200  {object}  ru.APIResponse[Tenant]
// @Failure      400  {object}  ru.APIResponse[any]
// @Failure      500  {object}  ru.APIResponse[any]
// @Router       /tenants/{id} [delete]
func (h tenantHandler) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	id, err := ru.RetrieveID(r)
	if err != nil {
		ru.HandleBadRequest(w, err, h.l)
		return
	}

	tenant, err := h.service.DeleteTenant(r.Context(), id)
	if err != nil {
		ru.HandleInternalServerError(w, err, ru.ResourceTenant, h.l)
		return
	}

	ru.WriteSuccessResponse(w, http.StatusOK, tenant, "Tenant deleted successfully")
}
