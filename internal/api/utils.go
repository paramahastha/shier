package api

import (
	"io"
	"net/http"
	"net/http/httptest"
)

func uniqueNum(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func DoRequest(handler http.Handler, method, path string, body io.Reader) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	return res, nil
}
