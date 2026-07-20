package backlog

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func fixtureServer(t *testing.T) *Server {
	t.Helper()
	a, err := Resolve("testdata", "knowledge")
	if err != nil {
		t.Fatal(err)
	}
	return NewServer(a)
}

func TestDiagRoute(t *testing.T) {
	srv := fixtureServer(t)
	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, httptest.NewRequest("GET", "/_diag", nil))
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	body := rec.Body.String()

	// AC5: exactly the three seeded reconciliation warnings, nothing unexpected.
	lines := strings.Split(strings.TrimSpace(body), "\n")
	if len(lines) != 3 {
		t.Errorf("/_diag has %d lines, want exactly 3:\n%s", len(lines), body)
	}
	for _, sub := range []string{
		"CURRENT holds >1 active task",
		"DONE item not ticked in BACKLOG: DEMO-6",
		"shipped item missing from DONE.md: DEMO-5",
	} {
		if !strings.Contains(body, sub) {
			t.Errorf("/_diag missing %q", sub)
		}
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("/_diag Content-Type = %q, want text/plain", ct)
	}
}

func TestBoardRoute(t *testing.T) {
	srv := fixtureServer(t)
	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != 200 {
		t.Fatalf("board status %d", rec.Code)
	}
}
