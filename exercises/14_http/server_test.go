package todoapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func do(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestCreateAndGet(t *testing.T) {
	s := NewServer()

	rec := do(t, s, "POST", "/todos", `{"title":"learn go"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST status = %d; want 201, body=%s", rec.Code, rec.Body)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q; want application/json", ct)
	}
	var created Todo
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	if created.ID != 1 || created.Title != "learn go" || created.Done {
		t.Fatalf("created = %+v", created)
	}

	rec = do(t, s, "GET", "/todos/1", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET status = %d; want 200", rec.Code)
	}
	var got Todo
	_ = json.NewDecoder(rec.Body).Decode(&got)
	if got != created {
		t.Errorf("GET = %+v; want %+v", got, created)
	}
}

func TestList(t *testing.T) {
	s := NewServer()
	rec := do(t, s, "GET", "/todos", "")
	if rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Fatalf("empty list = %d %q; want 200 []", rec.Code, rec.Body)
	}
	for _, title := range []string{"a", "b", "c"} {
		do(t, s, "POST", "/todos", `{"title":"`+title+`"}`)
	}
	rec = do(t, s, "GET", "/todos", "")
	var list []Todo
	_ = json.NewDecoder(rec.Body).Decode(&list)
	if len(list) != 3 || list[0].Title != "a" || list[2].Title != "c" {
		t.Errorf("list = %+v; want 3 todos sorted by id", list)
	}
}

func TestErrors(t *testing.T) {
	s := NewServer()
	tests := []struct {
		name, method, path, body string
		want                     int
	}{
		{"invalid json", "POST", "/todos", `{bad`, http.StatusBadRequest},
		{"empty title", "POST", "/todos", `{"title":"  "}`, http.StatusBadRequest},
		{"not found", "GET", "/todos/99", "", http.StatusNotFound},
		{"bad id", "GET", "/todos/abc", "", http.StatusBadRequest},
		{"delete missing", "DELETE", "/todos/99", "", http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := do(t, s, tt.method, tt.path, tt.body)
			if rec.Code != tt.want {
				t.Fatalf("status = %d; want %d", rec.Code, tt.want)
			}
			var body map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil || body["error"] == "" {
				t.Errorf("error body = %q; want {\"error\": \"...\"}", rec.Body)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	s := NewServer()
	do(t, s, "POST", "/todos", `{"title":"x"}`)
	if rec := do(t, s, "DELETE", "/todos/1", ""); rec.Code != http.StatusNoContent {
		t.Fatalf("DELETE status = %d; want 204", rec.Code)
	}
	if rec := do(t, s, "GET", "/todos/1", ""); rec.Code != http.StatusNotFound {
		t.Errorf("GET after delete = %d; want 404", rec.Code)
	}
}
