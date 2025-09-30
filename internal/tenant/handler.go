package tenant

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"flowcargo/internal/shared/logger"

	"github.com/go-playground/validator/v10"
)

type TenantHandler struct {
	service TenantService
	l       logger.Logger
}

func NewTenantHandler(service TenantService, l logger.Logger) TenantHandler {
	return TenantHandler{
		service: service,
		l:       l,
	}
}

var validate = validator.New()

func (h TenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var ctp CreateTenantParams
	if err := json.NewDecoder(r.Body).Decode(&ctp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := validate.Struct(ctp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tenant, err := h.service.CreateTenant(r.Context(), ctp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(tenant); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h TenantHandler) GetTenant(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	if idString == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(idString)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}
	tenant, err := h.service.GetTenantByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tenant); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h TenantHandler) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	if idString == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(idString)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	var utp UpdateTenantParams
	if err := json.NewDecoder(r.Body).Decode(&utp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validate.Struct(utp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tenant, err := h.service.UpdateTenant(r.Context(), id, utp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tenant); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h TenantHandler) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	if idString == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(idString)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	tenant, err := h.service.SoftDeleteTenant(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tenant); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
