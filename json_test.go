package main

import (
	"net/http/httptest"
	"testing"
)

func TestRespondWithError(t *testing.T) {
	testCases := []struct {
		name     string
		code     int
		msg      string
		wantCode int
		wantBody string
	}{
		{
			name:     "400 Bad Request",
			code:     400,
			msg:      "Bad Request",
			wantCode: 400,
			wantBody: `{"error":"Bad Request"}`,
		},
		{
			name:     "500 Internal Server Error",
			code:     500,
			msg:      "Internal Server Error",
			wantCode: 500,
			wantBody: `{"error":"Internal Server Error"}`,
		},
		{
			name:     "200 OK, but error message",
			code:     200,
			msg:      "OK with error",
			wantCode: 200,
			wantBody: `{"error":"OK with error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			respondWithError(w, tc.code, tc.msg)

			if w.Code != tc.wantCode {
				t.Errorf("respondWithError(%d, %q) code = %d; want %d", tc.code, tc.msg, w.Code, tc.wantCode)
			}

			if body := w.Body.String(); body != tc.wantBody {
				t.Errorf("respondWithError(%d, %q) body = %q; want %q", tc.code, tc.msg, body, tc.wantBody)
			}
		})
	}
}

func TestRespondWithJSON(t *testing.T) {
	testCases := []struct {
		name     string
		code     int
		payload  interface{}
		wantCode int
		wantBody string
	}{
		{
			name:     "Success",
			code:     200,
			payload:  map[string]string{"message": "success"},
			wantCode: 200,
			wantBody: `{"message":"success"}`,
		},
		{
			name:     "Empty Payload",
			code:     200,
			payload:  nil,
			wantCode: 200,
			wantBody: `null`,
		},
		{
			name:     "Integer Payload",
			code:     200,
			payload:  123,
			wantCode: 200,
			wantBody: `123`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			respondWithJSON(w, tc.code, tc.payload)

			if w.Code != tc.wantCode {
				t.Errorf("respondWithJSON(%d, %v) code = %d; want %d", tc.code, tc.payload, w.Code, tc.wantCode)
			}

			if body := w.Body.String(); body != tc.wantBody {
				t.Errorf("respondWithJSON(%d, %v) body = %q; want %q", tc.code, tc.payload, body, tc.wantBody)
			}

			if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("respondWithJSON(%d, %v) Content-Type = %q; want %q", tc.code, tc.payload, contentType, "application/json")
			}
		})
	}

	t.Run("JSON Marshalling Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		// Create a payload that will cause JSON marshalling to fail.
		payload := make(chan int)
		respondWithJSON(w, 200, payload)

		if w.Code != 500 {
			t.Errorf("respondWithJSON() with unmarshallable payload code = %d; want 500", w.Code)
		}
	})
}
