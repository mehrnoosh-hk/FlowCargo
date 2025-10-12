package restutils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// RetrieveID extracts and validates a UUID from the "id" path parameter of the HTTP request.
func RetrieveID(r *http.Request) (uuid.UUID, error) {
	idString := r.PathValue("id")
	if idString == "" {
		return uuid.Nil, errors.New("missing id parameter")
	}
	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.Nil, errors.New("invalid id parameter")
	}
	return id, nil
}

// RetrieveBody decodes and validates the JSON body of the HTTP request into the specified struct type.
func RetrieveBody[T any](r *http.Request) (T, error) {
	var body T
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return body, err
	}
	var validate = validator.New()
	if err := validate.Struct(body); err != nil {
		return body, err
	}
	return body, nil
}
