package handlers_test

import (
	"reflect"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/Joey-Boivin/sdisk-api/api/handlers"
)

func TestPingHandler(t *testing.T) {

	pingHandler := handlers.NewPingHandler()

	t.Run("ReturnHttpStatusOk", func(t *testing.T) {
		response := httptest.NewRecorder()
		getRequest, _ := http.NewRequest(http.MethodGet, handlers.PingEndpoint, nil)

		pingHandler.Ping(response, getRequest)

		got := response.Code
		want := http.StatusOK
		assertStatus(t, got, want)
	})

	t.Run("SendPong", func(t *testing.T) {
		response := httptest.NewRecorder()
		getRequest, _ := http.NewRequest(http.MethodGet, handlers.PingEndpoint, nil)

		pingHandler.Ping(response, getRequest)

		got := response.Body.Bytes()
		want := []byte("pong")
		assertEquals(t, got, want)
	})
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
