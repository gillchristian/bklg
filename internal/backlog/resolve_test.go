package backlog

import (
	"errors"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestResolveLocationsDereference(t *testing.T) { // AC1
	a, err := Resolve("testdata/resolve/withloc", "knowledge")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantP := filepath.Join("testdata/resolve/withloc", "alt/planning")
	wantR := filepath.Join("testdata/resolve/withloc", "alt/progress")
	if a.PlanningDir != wantP {
		t.Errorf("PlanningDir = %q, want %q (Locations dereference, resolved against path)", a.PlanningDir, wantP)
	}
	if a.ProgressDir != wantR {
		t.Errorf("ProgressDir = %q, want %q (Locations dereference, resolved against path)", a.ProgressDir, wantR)
	}
}

func TestResolveDefaultFallback(t *testing.T) { // AC2
	a, err := Resolve("testdata/resolve/default", "knowledge")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantP := filepath.Join("testdata/resolve/default/knowledge", "planning")
	wantR := filepath.Join("testdata/resolve/default/knowledge", "progress")
	if a.PlanningDir != wantP || a.ProgressDir != wantR {
		t.Errorf("got planning=%q progress=%q, want %q / %q (default fallback)", a.PlanningDir, a.ProgressDir, wantP, wantR)
	}
}

func TestResolveRootManifest(t *testing.T) { // AC3
	_, err := Resolve("testdata/resolve/rootmanifest", "knowledge")
	var rme *RootManifestError
	if !errors.As(err, &rme) {
		t.Fatalf("want *RootManifestError, got %v", err)
	}
	want := []string{"systems/alpha", "systems/beta"}
	if !slices.Equal(rme.Systems, want) {
		t.Errorf("Systems = %v, want %v (distinct, first-seen order)", rme.Systems, want)
	}
}

func TestResolveNoPlanningArea(t *testing.T) { // AC4
	_, err := Resolve("testdata/resolve/empty", "knowledge")
	if err == nil || !strings.Contains(err.Error(), "no planning area at") {
		t.Fatalf("want 'no planning area at' error, got %v", err)
	}
	var rme *RootManifestError
	if errors.As(err, &rme) {
		t.Errorf("empty dir (no systems table) must not be a RootManifestError")
	}
}

func TestResolvePathNotDir(t *testing.T) { // AC5
	for _, p := range []string{"testdata/resolve/afile.txt", "testdata/resolve/does-not-exist"} {
		if _, err := Resolve(p, "knowledge"); err == nil || !strings.Contains(err.Error(), "not a directory") {
			t.Errorf("Resolve(%q): want 'not a directory' error, got %v", p, err)
		}
	}
}

func TestParseLocations(t *testing.T) {
	md := strings.Join([]string{
		"# title",
		"",
		"## Locations",
		"",
		"some prose: even with a colon is ignored",
		"planning:   a/plan",
		"progress: a/prog",
		"framework: fw",
		"",
		"## Next section",
		"planning: leaked/after/block",
	}, "\n")
	loc := parseLocations(md)
	if loc["planning"] != "a/plan" || loc["progress"] != "a/prog" {
		t.Errorf("got %v, want planning=a/plan progress=a/prog", loc)
	}
	if _, ok := loc["framework"]; ok {
		t.Errorf("non-planning/progress keys must not be kept")
	}
	if loc["planning"] == "leaked/after/block" {
		t.Errorf("a key line after the block leaked into the map")
	}
}

func TestParseLocationsPartial(t *testing.T) {
	// Only planning present: it is captured; progress is absent from the map so
	// Resolve keeps progress at its base/progress default.
	md := "## Locations\n\nplanning: only/plan\n\n## End\n"
	loc := parseLocations(md)
	if loc["planning"] != "only/plan" {
		t.Errorf("planning = %q, want only/plan", loc["planning"])
	}
	if _, ok := loc["progress"]; ok {
		t.Errorf("progress must be absent when the block omits it, got %q", loc["progress"])
	}
}

func TestParseSystems(t *testing.T) {
	md := strings.Join([]string{
		"## Systems",
		"| a | systems/alpha  | note |",
		"| b | systems/beta/knowledge | note |",
		"| a | systems/alpha  | dup  |",
		"prose systems/ignored has no pipe so is skipped",
	}, "\n")
	got := parseSystems(md)
	want := []string{"systems/alpha", "systems/beta"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestDisplayPath(t *testing.T) {
	cases := map[string]string{
		"knowledge/planning":           "./knowledge/planning",
		"./already":                    "./already",
		"/abs/path":                    "/abs/path",
		"systems/x/knowledge/planning": "./systems/x/knowledge/planning",
	}
	for in, want := range cases {
		if got := DisplayPath(in); got != want {
			t.Errorf("DisplayPath(%q) = %q, want %q", in, got, want)
		}
	}
}
