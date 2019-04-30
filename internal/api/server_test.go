package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestStatusOk(t *testing.T) {
	// expected body
	body := gin.H{
		"status": map[string]string{"message": "Ok"},
	}

	// router
	router := SetupRouter()

	// Perform a GET request with that handler.
	w := performRequest(router, "GET", "/")

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)

	// Convert json response to a map
	var response map[string]map[string]string
	err := json.Unmarshal([]byte(w.Body.Bytes()), &response)

	// Grab the value & whether or not it exists
	value, exists := response["status"]["message"]

	// Make some assertions on the correctness of the response.
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, body["status"].(map[string]string)["message"], value)

}
