package tenant

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	
	"flowcargo/internal/shared/logger"
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

func (h TenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var ctr CreateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&ctr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tenant, err := h.service.CreateTenant(r.Context(), CreateTenantParams{
		ctr.Name,
		ctr.Email,
		ctr.Domain,
	})
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

	var utr UpdateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&utr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := UpdateTenantParams{
		Name:     utr.Name,
		Email:    utr.Email,
		Domain:   utr.Domain,
		IsActive: utr.IsActive,
	}
	
	h.l.Debugf("Update params: %+v", params)

	tenant, err := h.service.UpdateTenant(r.Context(), id, params)
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
	h.l.Info("Delete request received")
	idString := r.PathValue("id")
	if idString == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	h.l.Infof("ID: %v", idString)
	id, err := uuid.Parse(idString)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}
	h.l.Infof("ID: %v", id)

	tenant, err := h.service.DeleteTenant(r.Context(), id)
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

