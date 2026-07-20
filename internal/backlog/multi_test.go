package backlog

import (
	"errors"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

// The mono fixture's root manifest resolves to a RootManifestError listing its
// systems (the signal main uses to build a multi-server).
func TestResolveMonoRoot(t *testing.T) {
	_, err := Resolve("testdata/mono", "knowledge")
	var rme *RootManifestError
	if !errors.As(err, &rme) {
		t.Fatalf("want RootManifestError, got %v", err)
	}
	if !slices.Equal(rme.Systems, []string{"systems/alpha", "systems/beta"}) {
		t.Errorf("Systems = %v", rme.Systems)
	}
}

// AC1/AC2: aggregate board combines both systems, each card tagged with its system.
func TestMultiServerAggregates(t *testing.T) {
	srv := NewMultiServer("testdata/mono", []string{"systems/alpha", "systems/beta"})
	b := srv.board()
	sysOf := map[string]string{}
	for _, c := range b.Cards {
		if c.ID != nil {
			sysOf[c.ID.Raw] = c.System
		}
	}
	if sysOf["X-1"] != "alpha" || sysOf["X-2"] != "alpha" {
		t.Errorf("alpha cards mis-tagged: %v", sysOf)
	}
	if sysOf["Y-1"] != "beta" || sysOf["Y-2"] != "beta" {
		t.Errorf("beta cards mis-tagged: %v", sysOf)
	}
	if len(b.Cards) != 4 {
		t.Errorf("aggregate cards = %d, want 4", len(b.Cards))
	}
}

// AC3: ?system= filters; the bar lists all systems with counts.
func TestMultiBoardFilter(t *testing.T) {
	srv := NewMultiServer("testdata/mono", []string{"systems/alpha", "systems/beta"})
	b := srv.board()

	sys := []string{"alpha", "beta"}
	all := viewModel(b, "", sys)
	if len(all.Systems) != 3 { // All + alpha + beta
		t.Errorf("filter bar = %d chips, want 3 (All, alpha, beta)", len(all.Systems))
	}

	alpha := viewModel(b, "alpha", sys)
	for _, col := range alpha.Columns {
		for _, c := range col.Cards {
			if c.System != "alpha" {
				t.Errorf("filter=alpha leaked a %s card", c.System)
			}
		}
	}
	// the alpha chip is active
	var activeAlpha bool
	for _, chip := range alpha.Systems {
		if chip.Name == "alpha" {
			activeAlpha = chip.Active
		}
	}
	if !activeAlpha {
		t.Error("alpha chip should be active when filtered to alpha")
	}

	// A discovered system with zero cards still appears (filter to any project).
	withEmpty := viewModel(b, "", []string{"alpha", "beta", "empty"})
	var sawEmpty bool
	for _, chip := range withEmpty.Systems {
		if chip.Name == "empty" && chip.Count == 0 {
			sawEmpty = true
		}
	}
	if !sawEmpty {
		t.Error("an empty system should still show in the filter bar with count 0")
	}
}

// AC1/AC3/AC4 through the HTTP surface.
func TestMultiBoardRoutes(t *testing.T) {
	srv := NewMultiServer("testdata/mono", []string{"systems/alpha", "systems/beta"})
	h := srv.Routes()

	get := func(path string) (int, string) {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		return rec.Code, rec.Body.String()
	}

	code, body := get("/")
	if code != 200 {
		t.Fatalf("/ status %d", code)
	}
	for _, want := range []string{"X-1", "Y-1", `href="/?system=alpha"`, `href="/?system=beta"`} {
		if !strings.Contains(body, want) {
			t.Errorf("/ missing %q", want)
		}
	}

	// AC3: filter to alpha hides beta cards.
	_, bodyA := get("/?system=alpha")
	if !strings.Contains(bodyA, "X-1") || strings.Contains(bodyA, "Y-1") {
		t.Error("/?system=alpha should show X-1 and hide Y-1")
	}

	// AC4: detail across systems + /_v.
	if c, _ := get("/Y-1"); c != 200 {
		t.Errorf("/Y-1 detail status %d, want 200", c)
	}
	if c, v := get("/_v"); c != 200 || strings.TrimSpace(v) == "0" {
		t.Errorf("/_v = %q (status %d)", v, c)
	}
}

// AC5: a system in the index that fails to resolve is skipped with a warning,
// not a crash; the other systems still render.
func TestMultiServerSkipsUnresolvable(t *testing.T) {
	srv := NewMultiServer("testdata/mono", []string{"systems/alpha", "systems/nope"})
	b := srv.board() // must not panic
	var readErr bool
	for _, w := range b.Warnings {
		if w.Kind == "read-error" && strings.Contains(w.Message, "nope") {
			readErr = true
		}
	}
	if !readErr {
		t.Errorf("unresolvable system should add a read-error warning; warnings=%v", b.Warnings)
	}
	// alpha still contributed its cards.
	var sawAlpha bool
	for _, c := range b.Cards {
		if c.System == "alpha" {
			sawAlpha = true
		}
	}
	if !sawAlpha {
		t.Error("alpha cards should still render when a sibling system fails")
	}
}
