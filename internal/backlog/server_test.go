package backlog

import (
	"bytes"
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
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("/_diag Content-Type = %q, want text/html", ct)
	}
	body := rec.Body.String()

	// AC2: grouped by kind with titles + explanations.
	for _, sub := range []string{
		"More than one active task",         // current-multiple group title
		"Shipped, but not in DONE.md",       // shipped-missing-done group title
		"In DONE.md, not ticked in BACKLOG", // done-not-ticked group title
		"CURRENT holds &gt;1 active task",   // current-multiple message (no id), > escaped
	} {
		if !strings.Contains(body, sub) {
			t.Errorf("/_diag missing %q", sub)
		}
	}
	// AC3: id-bearing warnings link to their detail page.
	for _, id := range []string{"DEMO-5", "DEMO-6"} {
		if !strings.Contains(body, `href="/`+id+`"`) {
			t.Errorf("/_diag missing actionable link to /%s", id)
		}
	}
	// AC4: the shipped-missing explanation reframes it as informational.
	if !strings.Contains(body, "keeps the full record inline") {
		t.Error("/_diag missing the 'inline record is expected' explanation")
	}
}

func TestBoardRender(t *testing.T) {
	srv := fixtureServer(t)
	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != 200 {
		t.Fatalf("board status %d", rec.Code)
	}
	body := rec.Body.String()

	// AC1: three column headings.
	for _, h := range []string{"Backlog", "In Progress", "Done"} {
		if !strings.Contains(body, h) {
			t.Errorf("missing column heading %q", h)
		}
	}
	// AC2: a card per parsed task — id, title, and a link to its detail page.
	for _, id := range []string{"DEMO-1", "DEMO-2", "DEMO-3", "DEMO-4", "DEMO-5", "DEMO-6"} {
		if !strings.Contains(body, id) {
			t.Errorf("missing card id %q", id)
		}
		if !strings.Contains(body, `href="/`+id+`"`) {
			t.Errorf("missing detail link to /%s", id)
		}
	}
	if !strings.Contains(body, "Investigate the flaky importer") {
		t.Error("missing DEMO-2 title text")
	}
	// AC3: blocked badge + namespace chip.
	if !strings.Contains(body, "blocked") {
		t.Error("missing blocked badge (DEMO-2)")
	}
	if !strings.Contains(body, ">DEMO<") {
		t.Error("missing namespace chip DEMO")
	}
	// AC4: header planning path + diagnostics banner.
	if !strings.Contains(body, "testdata/knowledge/planning") {
		t.Error("missing resolved planning path in header")
	}
	if !strings.Contains(body, "/_diag") {
		t.Error("missing /_diag banner link (warnings present)")
	}
	// AC5: Tailwind Play CDN.
	if !strings.Contains(body, "cdn.tailwindcss.com") {
		t.Error("missing Tailwind Play CDN script")
	}
}

// AC5: repo text is auto-escaped by html/template — no HTML injection.
func TestBoardEscaping(t *testing.T) {
	var buf bytes.Buffer
	vm := boardVM{
		PlanningDir: "x",
		Columns:     []columnVM{{Title: "Backlog", Cards: []Card{{Title: "<script>alert(1)</script>", Column: ColBacklog}}}},
	}
	if err := boardTmpl.ExecuteTemplate(&buf, "layout", vm); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, "<script>alert(1)</script>") {
		t.Error("repo text was NOT escaped — injection risk")
	}
	if !strings.Contains(out, "&lt;script&gt;") {
		t.Error("expected escaped form &lt;script&gt;")
	}
}
