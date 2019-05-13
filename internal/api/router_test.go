package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusOk(t *testing.T) {
	var payload = []byte(`{"message": "hello world"}`)

	// expected body
	body := map[string]string{
		"message": "hello world",
	}

	// router
	router := SetupRouter()

	// Perform a GET request with that handler.
	w, _ := doRequest(router, "POST", "/test", bytes.NewBuffer(payload))

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)

	// Convert json response to a map
	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.Bytes()), &response)

	// Grab the value & whether or not it exists
	value, exists := response["message"]

	// Make some assertions on the correctness of the response.
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, body["message"], value)

}
