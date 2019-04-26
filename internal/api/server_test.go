package api

import (
	"encoding/json"
	"fmt"
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

// func TestStatusOk(t *testing.T) {
// 	// expected body
// 	body := gin.H{
// 		"status": map[string]string{"message": "Ok"},
// 	}

// 	// router
// 	router := SetupRouter()

// 	// Perform a GET request with that handler.
// 	w := performRequest(router, "GET", "/")

// 	// Assert we encoded correctly,
// 	// the request gives a 200
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Convert json response to a map
// 	var response map[string]map[string]string
// 	err := json.Unmarshal([]byte(w.Body.Bytes()), &response)

// 	// Grab the value & whether or not it exists
// 	value, exists := response["status"]["message"]

// 	// Make some assertions on the correctness of the response.
// 	assert.Nil(t, err)
// 	assert.True(t, exists)
// 	assert.Equal(t, body["status"].(map[string]string)["message"], value)

// }

func TestCreateUser(t *testing.T) {
	// expected body
	body := gin.H{
		"data": map[string]interface{}{
			"user": "Create user successfully",
		},
		"status": 200,
	}

	// router
	router := SetupRouter()

	// Perform a GET request with that handler.
	w := performRequest(router, "GET", "/")

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)

	type DataUser struct {
		User string `json:"user"`
	}

	type Response struct {
		Data   DataUser `json:"data"`
		Status int      `json:"status"`
	}

	// Convert json response to a map
	var response Response
	err := json.Unmarshal([]byte(w.Body.Bytes()), &response)

	// Grab the value & whether or not it exists
	fmt.Println("Log", body["data"].(map[string]interface{})["user"], response.Data)

	// Make some assertions on the correctness of the response.
	assert.Nil(t, err)
	// assert.True(t, exists)
	// assert.Equal(t, body["data"].(map[string]interface{})["user"], value)

}
