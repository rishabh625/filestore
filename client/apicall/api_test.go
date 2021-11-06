package apicall_test

import (
	"filestore/client/apicall"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHExists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameter
		assert.Equal(t, req.URL.String(), "/hexist")
		//file := req.Header.Get("file")
		hash := req.Header.Get("hash")
		assert.Equal(t, "abcd", hash)
		//assert.Equal(t, "temp1.txt", file)
		// Send response to be tested
		rw.WriteHeader(http.StatusOK)
	}))
	// Close the server when test finishes
	defer server.Close()
	apicall.ServerUrl = server.URL
	status := apicall.Hexists("abcd")
	assert.Equal(t, true, status)
}

func TestCopyCall(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameter
		assert.Equal(t, req.URL.String(), "/copy")
		// Send response to be tested
		rw.WriteHeader(http.StatusOK)
	}))
	// Close the server when test finishes
	defer server.Close()
	apicall.ServerUrl = server.URL
	status1 := apicall.CopyCall("temp1.txt", "abcd")
	assert.Equal(t, true, status1)
}

func TestAddCall(t *testing.T) {
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

func TestRemove(t *testing.T) {
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

func TestList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameter
		assert.Equal(t, req.Method, http.MethodGet)
		assert.Equal(t, req.URL.String(), "/ls")
		// Send response to be tested
		rw.WriteHeader(http.StatusOK)
	}))
	// Close the server when test finishes
	defer server.Close()
	apicall.ServerUrl = server.URL
	apicall.List()
}

func TestWc(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameter
		assert.Equal(t, req.Method, http.MethodGet)
		assert.Equal(t, req.URL.String(), "/wc")
		// Send response to be tested
		rw.WriteHeader(http.StatusOK)
	}))
	// Close the server when test finishes
	defer server.Close()
	apicall.ServerUrl = server.URL
	apicall.WC()
}

func TestFreqWords(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameter
		assert.Equal(t, req.Method, http.MethodGet)
		assert.Equal(t, "/freqwords?order=asc&limit=10", req.URL.String())
		// Send response to be tested
		rw.WriteHeader(http.StatusOK)
	}))
	// Close the server when test finishes
	defer server.Close()
	apicall.ServerUrl = server.URL
	apicall.FreqWords("10", "asc")
}
