package handlers_test

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/Joey-Boivin/cdisk/api/handlers"
)

func TestPingHandler(t *testing.T) {
	methodsExceptGet := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions}

	t.Run("GivenGetRequest_WhenHandling_ReturnHttpStatusOk", func(t *testing.T) {
		response := httptest.NewRecorder()
		getRequest := createRequest(http.MethodGet)

		handlers.PingHandler(response, getRequest)

		got := response.Code
		want := http.StatusOK
		assertStatus(t, got, want)
	})

	t.Run("GivenMethodNotGet_WhenHandling_ReturnHttpMethodNotAllowed", func(t *testing.T) {
		for _, method := range methodsExceptGet {
			response := httptest.NewRecorder()
			getRequest := createRequest(method)

			handlers.PingHandler(response, getRequest)

			got := response.Code
			want := http.StatusMethodNotAllowed
			assertStatus(t, got, want)
		}
	})
}

func createRequest(method string) *http.Request {
	req, _ := http.NewRequest(method, handlers.PingEndpoint, nil)
	return req
}
func assertStatus(t *testing.T, gotCode int, wantCode int) {
	t.Helper()

	if gotCode != wantCode {
		t.Fatalf("Got http status code %d. Should've been %d", gotCode, wantCode)
	}
}
