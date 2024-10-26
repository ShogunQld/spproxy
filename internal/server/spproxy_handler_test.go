package server_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"spproxy/internal/assert"
	"spproxy/internal/configs"
	"spproxy/internal/server"
	"strings"
	"testing"
)

type TestEndpoint struct {
	mockServer *httptest.Server
	resource   configs.Resource
	url        *url.URL
	handler    func(http.ResponseWriter, *http.Request)
	response   string
}

func NewTestEndpoint(proxyCache *server.ProxyCache, name, endpoint string, sticky bool, response string) *TestEndpoint {
	mock := mockServer(http.StatusOK, response)
	url, _ := url.Parse(mock.URL)
	port := ""
	if sticky {
		port = url.Port()
	}
	testEndpoint := &TestEndpoint{
		mockServer: mockServer(200, response),
		resource: configs.Resource{
			Name:            name,
			Endpoint:        endpoint,
			Port:            port,
			Destination_URL: mock.URL,
		},
		url:      url,
		response: response,
	}
	testEndpoint.handler = server.ProxyRequestHandler(testEndpoint.resource, proxyCache)
	fmt.Printf("Created Test Endpoint %s from %s to %s with sticky port: %s\n", name, endpoint, mock.URL, port)
	return testEndpoint
}

func TestProxyRequestHandler(t *testing.T) {
	proxyCache := server.NewProxyCache()

	endpointRoot := NewTestEndpoint(proxyCache, "Root", "/", true, "Response from Root")
	endpoint1 := NewTestEndpoint(proxyCache, "R1", "/app1/", true, "Response from endpoint 1")
	endpoint2 := NewTestEndpoint(proxyCache, "R2", "/app2/", true, "Response from endpoint 2")
	endpoint3 := NewTestEndpoint(proxyCache, "R3", "/static/", false, "Response from endpoint 3")

	t.Log("Hit endpoint one - expect redirection")
	req := httptest.NewRequest(http.MethodGet, "/app1/hello", nil)
	w := httptest.NewRecorder()
	endpoint1.handler(w, req)

	assert.Equal[int](t, w.Code, http.StatusSeeOther)
	assert.Equal[string](t, w.Header().Get("Location"), "/hello")
	assert.Equal[string](t, proxyCache.CurrentPort, endpoint1.resource.Port)

	t.Log("Hit root endpoint with redirect url - expect response from endpoint1 server")
	req = httptest.NewRequest(http.MethodGet, w.Header().Get("Location"), nil)
	w = httptest.NewRecorder()
	endpointRoot.handler(w, req)

	assert.Equal[int](t, w.Code, http.StatusOK)
	assert.Equal[string](t, proxyCache.CurrentPort, endpoint1.resource.Port)
	assert.Equal[string](t, strings.TrimSpace(w.Body.String()), endpoint1.response)

	t.Log("Hit endpoint /new - expect response from endpoint1 server")
	req = httptest.NewRequest(http.MethodGet, "/new", nil)
	w = httptest.NewRecorder()
	endpointRoot.handler(w, req)

	assert.Equal[int](t, w.Code, http.StatusOK)
	assert.Equal[string](t, proxyCache.CurrentPort, endpoint1.resource.Port)
	assert.Equal[string](t, strings.TrimSpace(w.Body.String()), endpoint1.response)

	t.Log("Hit non sicky endpoint3 - expect response from endpoint3 server")
	req = httptest.NewRequest(http.MethodGet, "/static/image", nil)
	w = httptest.NewRecorder()
	endpoint3.handler(w, req)

	assert.Equal[int](t, w.Code, http.StatusOK)
	assert.Equal[string](t, proxyCache.CurrentPort, endpoint1.resource.Port)
	assert.Equal[string](t, strings.TrimSpace(w.Body.String()), endpoint3.response)

	t.Log("Hit endpoint /new again - expect response from endpoint1 server")
	req = httptest.NewRequest(http.MethodGet, "/new", nil)
	w = httptest.NewRecorder()
	endpointRoot.handler(w, req)

	assert.Equal[int](t, w.Code, http.StatusOK)
	assert.Equal[string](t, proxyCache.CurrentPort, endpoint1.resource.Port)
	assert.Equal[string](t, strings.TrimSpace(w.Body.String()), endpoint1.response)

	t.Log("Hit endpoint two - expect redirection")
	req = httptest.NewRequest(http.MethodGet, "/app2/boo", nil)
	w = httptest.NewRecorder()
	endpoint2.handler(w, req)

	assert.Equal[int](t, w.Code, http.StatusSeeOther)
	assert.Equal[string](t, w.Header().Get("Location"), "/boo")
	assert.Equal[string](t, proxyCache.CurrentPort, endpoint2.resource.Port)

	t.Log("Hit root endpoint with redirect url - expect response from endpoint2 server")
	req = httptest.NewRequest(http.MethodGet, w.Header().Get("Location"), nil)
	w = httptest.NewRecorder()
	endpointRoot.handler(w, req)

	assert.Equal[int](t, w.Code, http.StatusOK)
	assert.Equal[string](t, proxyCache.CurrentPort, endpoint2.resource.Port)
	assert.Equal[string](t, strings.TrimSpace(w.Body.String()), endpoint2.response)
}

// Create a mock HTTP Server that will return a response with HTTP code and body.
func mockServer(code int, body string) *httptest.Server {

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprintln(w, body)
	}))
}
