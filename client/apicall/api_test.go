package apicall_test

import (
	"filestore/client/apicall"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func HExistsTest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameter
		assert.Equal(t, req.URL.String(), "/hexist")
		file := req.Header.Get("file")
		hash := req.Header.Get("hash")
		assert.Equal(t, "abcd", hash)
		assert.Equal(t, "temp1.txt", file)
		// Send response to be tested
		rw.WriteHeader(http.StatusOK)
	}))
	// Close the server when test finishes
	defer server.Close()
	apicall.ServerUrl = server.URL
	status1, status2 := apicall.Hexists("abcd", "temp1.txt")
	assert.Equal(t, true, status1)
	assert.Equal(t, true, status2)
}

func CopyCallTest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameter
		assert.Equal(t, req.URL.String(), "/copy")
		hash := req.Header.Get("hash")
		assert.Equal(t, "abcd", hash)
		// Send response to be tested
		rw.WriteHeader(http.StatusOK)
	}))
	// Close the server when test finishes
	defer server.Close()
	apicall.ServerUrl = server.URL
	status1 := apicall.CopyCall("temp1.txt", "abcd")
	assert.Equal(t, true, status1)
}

func AddCallTest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameter
		assert.Equal(t, req.URL.String(), "/store")
		hash := req.Header.Get("hash")
		assert.Equal(t, "abcd", hash)
		// Send response to be tested
		rw.WriteHeader(http.StatusOK)
	}))
	// Close the server when test finishes
	defer server.Close()
	apicall.ServerUrl = server.URL
	status1 := apicall.AddCall("temp1.txt", "abcd")
	assert.Equal(t, true, status1)
}

func RemoveTest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameter
		assert.Equal(t, req.Method, http.MethodDelete)
		assert.Equal(t, req.URL.String(), "/rm")
		// Send response to be tested
		rw.WriteHeader(http.StatusOK)
	}))
	// Close the server when test finishes
	defer server.Close()
	apicall.ServerUrl = server.URL
	apicall.Remove("temp1.txt")
}
