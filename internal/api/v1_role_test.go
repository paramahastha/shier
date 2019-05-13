package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {

	var payload = []byte(`{
		"name": "xyz"}`)

	// expected body
	body := map[string]interface{}{
		"data":    nil,
		"code":    201,
		"message": "Create user successfully",
	}

	// router
	router := SetupRouter()

	// Perform a GET request with that handler.
	w, _ := doRequest(router, "POST", "/test", bytes.NewBuffer(payload))

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusCreated, w.Code)

	// Convert json response to a map
	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.Bytes()), &response)

	// Grab the value & whether or not it exists
	value, exists := response["message"]

	// Make some assertions on the correctness of the response.
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, body["message"], value)

}
