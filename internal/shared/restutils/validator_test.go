package restutils

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRetriveID(t *testing.T) {
	t.Run("It returns an error if the id parameter is missing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/resource/", nil)
		_, err := RetrieveID(req)
		require.Equal(t, err.Error(), "missing id parameter")
	})

	t.Run("It returns an error if the id parameter is invalid", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/resource", nil)
		req.SetPathValue("id", "invalid-uuid")
		_, err := RetrieveID(req)
		require.Equal(t, err.Error(), "invalid id parameter")
	})

	t.Run("It successfully retrieves a valid UUID from the id parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/resource", nil)
		id := uuid.New()
		req.SetPathValue("id", id.String())
		returnedID, err := RetrieveID(req)
		require.NoError(t, err)
		require.Equal(t, id, returnedID)
	})
}


func TestRetriveBody(t *testing.T) {
	type TestBody struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"required,gte=0"`
	}

	t.Run("It returns an error if the body is not valid JSON", func(t *testing.T) {
		body := json.RawMessage(`{"name": 25}`)
		req := httptest.NewRequest("POST", "/resource", bytes.NewReader(body))
		_, err := RetrieveBody[TestBody](req)
		require.Error(t, err)
	})

	t.Run("It returns an error if the body fails validation", func(t *testing.T) {
		body := json.RawMessage(`{"name": "", "age": -5}`)
		req := httptest.NewRequest("POST", "/resource", bytes.NewReader(body))
		_, err := RetrieveBody[TestBody](req)
		require.Error(t, err)
	})

	t.Run("It successfully retrieves and validates the body", func(t *testing.T) {
		body := json.RawMessage(`{"name": "John", "age": 30}`)
		req := httptest.NewRequest("POST", "/resource", bytes.NewReader(body))
		returnedBody, err := RetrieveBody[TestBody](req)
		require.NoError(t, err)
		require.Equal(t, "John", returnedBody.Name)
		require.Equal(t, 30, returnedBody.Age)
	})
}