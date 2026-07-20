package backlog

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func get(t *testing.T, srv *Server, path string) (int, string) {
	t.Helper()
	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
	return rec.Code, rec.Body.String()
}

// AC1 (id/title/column/namespace/badges) + AC2 (In-Progress fields) + AC4 (raw).
func TestDetailInProgress(t *testing.T) {
	srv := fixtureServer(t)
	code, body := get(t, srv, "/DEMO-1")
	if code != 200 {
		t.Fatalf("status %d", code)
	}
	for _, want := range []string{"DEMO-1", "Wire up the widget pipeline", "In Progress", ">DEMO<"} {
		if !strings.Contains(body, want) {
			t.Errorf("missing %q", want)
		}
	}
	if !strings.Contains(body, "the pipeline builds") || !strings.Contains(body, "the widget renders") {
		t.Error("missing acceptance criteria text")
	}
	if !strings.Contains(body, "☑") || !strings.Contains(body, "☐") {
		t.Error("acceptance checklist should show a checked and an unchecked box")
	}
	if !strings.Contains(body, "may commit directly") {
		t.Error("missing delivery override")
	}
	if !strings.Contains(body, "Depends on the exporter") {
		t.Error("missing notes")
	}
	if !strings.Contains(body, "<details") || !strings.Contains(body, "### DEMO-1") {
		t.Error("missing collapsed raw source block")
	}
}

// AC2 Done-state fields.
func TestDetailDone(t *testing.T) {
	srv := fixtureServer(t)
	code, body := get(t, srv, "/DEMO-4")
	if code != 200 {
		t.Fatalf("status %d", code)
	}
	for _, want := range []string{"2026-07-18", "Streaming exporter", "PR #12, merged", "See journal 2026-07-18"} {
		if !strings.Contains(body, want) {
			t.Errorf("Done detail missing %q", want)
		}
	}
}

// AC3 referencing blockers (open first, then resolved).
func TestDetailBlockers(t *testing.T) {
	srv := fixtureServer(t)
	_, b2 := get(t, srv, "/DEMO-2")
	if !strings.Contains(b2, "BLOCKER-001") || !strings.Contains(b2, ">open<") {
		t.Error("DEMO-2 detail should show open BLOCKER-001")
	}
	_, b1 := get(t, srv, "/DEMO-1")
	if !strings.Contains(b1, "BLOCKER-002") || !strings.Contains(b1, ">resolved<") {
		t.Error("DEMO-1 detail should show resolved BLOCKER-002")
	}
}

// AC5 unknown id -> 404; case-insensitive lookup.
func TestDetail404AndCase(t *testing.T) {
	srv := fixtureServer(t)
	if code, _ := get(t, srv, "/NOPE-999"); code != 404 {
		t.Errorf("/NOPE-999 status %d, want 404", code)
	}
	if code, _ := get(t, srv, "/demo-1"); code != 200 {
		t.Errorf("/demo-1 (case-insensitive) status %d, want 200", code)
	}
}
