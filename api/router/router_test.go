package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Joey-Boivin/cdisk/api/router"
)

type HandlerMock struct {
	Called bool
}

func (h *HandlerMock) Get(w http.ResponseWriter, r *http.Request) {
	h.Called = true
}

type HandlerWithMultipleEndpointsMock struct {
	CalledPost bool
}

func (h *HandlerWithMultipleEndpointsMock) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.CalledPost = true
	}
}

func (h *HandlerWithMultipleEndpointsMock) AnyOtherMethod(w http.ResponseWriter, r *http.Request) {
	h.CalledPost = false
}

func TestRouter(t *testing.T) {

	router := router.NewRouter()

	t.Run("GivenARoute_WhenRoutingWithMatchingRouteParameters_ThenCallHandler", func(t *testing.T) {
		response := httptest.NewRecorder()
		handler := HandlerMock{}
		method := http.MethodGet
		endpoint := "/mock"
		router.AddRoute(handler.Get, method, endpoint)
		getRequest := createRequest(method, endpoint)

		router.ServeHTTP(response, getRequest)

		assertTrue(t, handler.Called)
	})

	t.Run("GivenMultipleRoutes_WhenRoutingWithMatchingRouteParemeters_ThenCallCorrectHandler", func(t *testing.T) {
		response := httptest.NewRecorder()
		handler := HandlerWithMultipleEndpointsMock{}
		method := http.MethodPost
		anyOtherMethod := http.MethodPut
		endpoint := "/mock"
		router.AddRoute(handler.Post, method, endpoint)
		router.AddRoute(handler.AnyOtherMethod, anyOtherMethod, endpoint)
		postRequest := createRequest(method, endpoint)

		router.ServeHTTP(response, postRequest)

		assertTrue(t, handler.CalledPost)
	})

	t.Run("WhenRoutingWithNoMatchingEndpoints_ThenReturnHttpNotFound", func(t *testing.T) {
		response := httptest.NewRecorder()
		anyMethod := http.MethodConnect
		anyEndpoint := "/any"
		anyRequestWithNoMatchingEndpoint := createRequest(anyMethod, anyEndpoint)

		router.ServeHTTP(response, anyRequestWithNoMatchingEndpoint)

		got := response.Code
		want := http.StatusNotFound
		assertStatus(t, got, want)
	})

	t.Run("WhenRoutingWithNoMatchingMethod_ThenReturnHttpMethodNotAllowed", func(t *testing.T) {
		response := httptest.NewRecorder()
		handler := HandlerMock{}
		method := http.MethodGet
		endpoint := "/mock"
		router.AddRoute(handler.Get, method, endpoint)
		anyRequestWithNoMatchingMethod := createRequest(http.MethodDelete, endpoint)

		router.ServeHTTP(response, anyRequestWithNoMatchingMethod)

		got := response.Code
		want := http.StatusMethodNotAllowed
		assertStatus(t, got, want)
	})
}

func createRequest(method string, endpoint string) *http.Request {
	req, _ := http.NewRequest(method, endpoint, nil)
	return req
}

func assertTrue(t *testing.T, expression bool) {
	t.Helper()

	if !expression {
		t.Fatalf("Expected true, got false")
	}
}

func assertStatus(t *testing.T, gotCode int, wantCode int) {
	t.Helper()

	if gotCode != wantCode {
		t.Fatalf("Got http status code %d. Should've been %d", gotCode, wantCode)
	}
}
