package handlers_test

import (
	"reflect"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/Joey-Boivin/cdisk/api/handlers"
)

func TestPingHandler(t *testing.T) {

	pingHandler := handlers.PingHandler{}

	t.Run("ReturnHttpStatusOk", func(t *testing.T) {
		response := httptest.NewRecorder()
		getRequest := createRequest(http.MethodGet)

		pingHandler.Get(response, getRequest)

		got := response.Code
		want := http.StatusOK
		assertStatus(t, got, want)
	})

	t.Run("SendPong", func(t *testing.T) {
		response := httptest.NewRecorder()
		getRequest := createRequest(http.MethodGet)

		pingHandler.Get(response, getRequest)

		got := response.Body.Bytes()
		want := []byte("pong")
		assertEquals(t, got, want)
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

func assertEquals(t *testing.T, got []byte, want []byte) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got response %s. Should've been %s", string(got), string(want))
	}
}
