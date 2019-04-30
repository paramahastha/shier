package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	// form := &struct {
	// 	FirstName string `form:"first_name" json:"first_name"`
	// 	LastName  string `form:"last_name" json:"last_name"`
	// 	Email     string `form:"email" json:"email"`
	// 	Password  string `form:"password" json:"password"`
	// 	Confirm   string `form:"confirm" json:"confirm"`
	// 	Roles     []int  `form:"roles" json:"roles"`
	// }{
	// 	"foo",
	// 	"bar",
	// 	"foo.bar@mail.com",
	// 	"secret",
	// 	"secret",
	// 	[]int{1},
	// }
	form := map[string]interface{}{
		"first_name": "foo",
		"last_name":  "bar",
		"email":      "foo.bar@mail.com",
		"password":   "secret",
		"confirm":    "secret",
		"roles":      []int{1},
	}

	// expected body response
	expectedResponse := map[string]interface{}{
		"data": map[string]interface{}{
			"user": "Create user successfully",
		},
		"status": 200,
	}

	// router
	router := SetupRouter()

	jsonUser, _ := json.Marshal(form)
	// Perform a POST request with that handler.
	request, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonUser))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	// Assert we encoded correctly,
	// the request gives a 200
	fmt.Println("LOG", w.Code)

	// // // Convert json response to a map
	var response struct {
		Data struct {
			User string `json:"user"`
		} `json:"data"`
		Status int `json:"status"`
	}
	err := json.Unmarshal([]byte(w.Body.String()), &response)

	// // Grab the value & whether or not it exists
	fmt.Println("Log", expectedResponse["data"], response.Data, err)

	// // Make some assertions on the correctness of the response.
	// assert.Nil(t, err)
	// assert.True(t, exists)
	// assert.Equal(t, body["data"].(map[string]interface{})["user"], value)

}
