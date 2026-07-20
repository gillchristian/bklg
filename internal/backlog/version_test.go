package backlog

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// AC2: Meta.LatestMTime is the newest of the four parsed files.
func TestMetaLatestMTime(t *testing.T) {
	a, err := Resolve("testdata", "knowledge")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := NewParser().Parse(a)
	if b.Meta.LatestMTime.IsZero() {
		t.Fatal("LatestMTime is zero; want the newest parsed-file mtime")
	}
	var want time.Time
	for _, p := range []string{
		filepath.Join(a.PlanningDir, "CURRENT.md"),
		filepath.Join(a.PlanningDir, "BACKLOG.md"),
		filepath.Join(a.PlanningDir, "DONE.md"),
		filepath.Join(a.ProgressDir, "blockers.md"),
	} {
		fi, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		if fi.ModTime().After(want) {
			want = fi.ModTime()
		}
	}
	if !b.Meta.LatestMTime.Equal(want) {
		t.Errorf("LatestMTime = %v, want %v (max of the four files)", b.Meta.LatestMTime, want)
	}
}

// AC1 shape + AC5: /_v is a 200 text/plain bare integer.
func TestVersionRoute(t *testing.T) {
	srv := fixtureServer(t)
	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, httptest.NewRequest("GET", "/_v", nil))
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("/_v Content-Type = %q, want text/plain", ct)
	}
	body := strings.TrimSpace(rec.Body.String())
	if !regexp.MustCompile(`^-?\d+$`).MatchString(body) {
		t.Errorf("/_v body = %q, want a bare integer", body)
	}
	if body == "0" {
		t.Error("/_v = 0 on a real instance; expected a real mtime")
	}
}

// AC3: the board embeds the poll script, and its baked-in version matches /_v.
func TestPollScript(t *testing.T) {
	srv := fixtureServer(t)
	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	body := rec.Body.String()
	for _, want := range []string{`fetch("/_v"`, "location.reload"} {
		if !strings.Contains(body, want) {
			t.Errorf("board missing poll-script piece %q", want)
		}
	}
	rec2 := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec2, httptest.NewRequest("GET", "/_v", nil))
	v := strings.TrimSpace(rec2.Body.String())
	if !strings.Contains(body, `var v = "`+v+`"`) {
		t.Errorf("board's baked version does not match /_v (%q)", v)
	}
}

// AC1/AC4: the version strictly changes when a parsed file is modified. Chtimes
// to a future time makes this deterministic regardless of filesystem mtime
// resolution.
func TestVersionChangesOnModify(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "knowledge/planning"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "knowledge/progress"), 0o755); err != nil {
		t.Fatal(err)
	}
	cur := filepath.Join(dir, "knowledge/planning/CURRENT.md")
	if err := os.WriteFile(cur, []byte("## Active\n### X-1 — a task\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	a, err := Resolve(dir, "knowledge")
	if err != nil {
		t.Fatal(err)
	}
	v1 := versionString(areaMTime(a))
	future := time.Now().Add(2 * time.Second)
	if err := os.Chtimes(cur, future, future); err != nil {
		t.Fatal(err)
	}
	v2 := versionString(areaMTime(a))
	if v1 == v2 {
		t.Errorf("version unchanged after modifying a parsed file (both %s)", v1)
	}
}
