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

// --- dashboard adapter resolution (TASK-013, ADR-0004) ----------------------

func TestResolveDashboardViaLocations(t *testing.T) { // AC1 + AC2
	// The fixture's only manifest is index.md (not README.md) and it carries a
	// dashboard: key; resolution must find the index.md manifest, select
	// dashboard mode, and not require a planning/ dir (there is none).
	root := "testdata/resolve/dashboard"
	a, err := Resolve(root, "knowledge")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join(root, "knowledge/work/index.md")
	if a.DashboardFile != want {
		t.Errorf("DashboardFile = %q, want %q", a.DashboardFile, want)
	}
	if a.PlanningDir != "" || a.ProgressDir != "" {
		t.Errorf("dashboard mode should not set planning/progress, got planning=%q progress=%q", a.PlanningDir, a.ProgressDir)
	}
	if a.KnowledgeDir != filepath.Join(root, "knowledge") {
		t.Errorf("KnowledgeDir = %q, want %q", a.KnowledgeDir, filepath.Join(root, "knowledge"))
	}
	if a.LinkBase != "https://linear.app/acme/issue/" {
		t.Errorf("LinkBase = %q, want the manifest's linear: value", a.LinkBase)
	}
}

func TestResolveDashboardDefaultLinkBase(t *testing.T) { // TASK-015
	// The flag path has no manifest to read a linear: key from, so it defaults.
	a, err := ResolveDashboard("testdata/resolve/dashboard", "knowledge", "knowledge/work/index.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.LinkBase != defaultLinearBase {
		t.Errorf("LinkBase = %q, want default %q", a.LinkBase, defaultLinearBase)
	}
}

func TestResolveDashboardFlag(t *testing.T) { // AC3
	root := "testdata/resolve/dashboard"
	a, err := ResolveDashboard(root, "knowledge", "knowledge/work/index.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join(root, "knowledge/work/index.md")
	if a.DashboardFile != want {
		t.Errorf("DashboardFile = %q, want %q", a.DashboardFile, want)
	}
}

func TestResolveDashboardMissingFile(t *testing.T) { // AC4
	// Via the flag:
	if _, err := ResolveDashboard("testdata/resolve/dashboard", "knowledge", "knowledge/nope.md"); err == nil || !strings.Contains(err.Error(), "no dashboard file at") {
		t.Errorf("flag: want 'no dashboard file at' error, got %v", err)
	}
	// A non-directory path is still rejected before the file check:
	if _, err := ResolveDashboard("testdata/resolve/afile.txt", "knowledge", "x.md"); err == nil || !strings.Contains(err.Error(), "not a directory") {
		t.Errorf("path: want 'not a directory' error, got %v", err)
	}
	// Via a manifest dashboard: key pointing at a nonexistent file — must fail
	// with the same error, not fall through to the planning-area check:
	if _, err := Resolve("testdata/resolve/dashboard-missing", "knowledge"); err == nil || !strings.Contains(err.Error(), "no dashboard file at") {
		t.Errorf("locations: want 'no dashboard file at' error, got %v", err)
	}
}

func TestResolveFrameworkModeUnaffected(t *testing.T) { // AC5
	// A repo with a normal README.md + planning area still resolves framework
	// mode (no DashboardFile), i.e. the manifest-lookup widening didn't regress.
	a, err := Resolve("testdata/resolve/withloc", "knowledge")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.DashboardFile != "" {
		t.Errorf("framework repo should not be dashboard mode, got DashboardFile=%q", a.DashboardFile)
	}
}

func TestParseLocationsDashboard(t *testing.T) {
	loc := parseLocations("## Locations\n\ndashboard: work/index.md\n\n## End\n")
	if loc["dashboard"] != "work/index.md" {
		t.Errorf("dashboard = %q, want work/index.md", loc["dashboard"])
	}
}
