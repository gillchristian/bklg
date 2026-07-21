// Package backlog resolves, parses, and models a knowledge-framework instance's
// planning and progress areas for the bklg board viewer.
package backlog

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Areas are the resolved directories bklg reads. Paths are in filepath form
// (suitable for os operations); use DisplayPath for the startup echo.
type Areas struct {
	KnowledgeDir string // base = path/dir
	PlanningDir  string
	ProgressDir  string
	// DashboardFile, when non-empty, selects the single-file dashboard adapter
	// (ADR-0004): the parser reads this one Active/Backlog/Done file instead of
	// the planning/progress skeleton. Contract: reference/specs/dashboard-format.md.
	DashboardFile string
}

// RootManifestError signals that the resolved planning area is absent while the
// knowledge manifest looks like a multi-system root (it lists systems/<name>).
// The caller lists the per-system invocations rather than erroring blankly.
type RootManifestError struct {
	Path         string   // the path argument (for building invocations)
	Dir          string   // the --dir value
	ManifestPath string   // base/README.md that carried the system index
	PlanningDir  string   // the planning dir that was missing
	Systems      []string // distinct "systems/<name>", first-seen order
}

func (e *RootManifestError) Error() string {
	return fmt.Sprintf("no planning area at %s (looks like a multi-system root manifest listing %d system(s))",
		e.PlanningDir, len(e.Systems))
}

var systemRe = regexp.MustCompile(`systems/[A-Za-z0-9._-]+`)

// Resolve locates the areas for the repo at path, using the knowledge dir named
// by dir. Resolution order matches spec §3, extended by ADR-0004 for dashboards:
//
//  1. the manifest (base/README.md, else base/index.md) "## Locations" block
//     (values are repo-root-relative, so resolved against path);
//  2. if that block carries a dashboard: key, return dashboard-mode Areas
//     (single-file adapter) and stop — no planning area is required;
//  3. else the planning/progress keys, or default base/planning + base/progress;
//  4. if the planning dir is absent, a manifest listing systems/<name> is
//     treated as a root manifest (RootManifestError); otherwise the error is
//     "no planning area at <planningDir>".
func Resolve(path, dir string) (Areas, error) {
	if fi, err := os.Stat(path); err != nil || !fi.IsDir() {
		return Areas{}, fmt.Errorf("path is not a directory: %s", path)
	}

	base := filepath.Join(path, dir)
	planning := filepath.Join(base, "planning")
	progress := filepath.Join(base, "progress")

	// The manifest is README.md by convention; dashboard-shaped KBs (ADR-0004)
	// name theirs index.md, so try that as a fallback. First hit wins.
	var manifest string
	var manifestBytes []byte
	manifestErr := os.ErrNotExist
	for _, name := range []string{"README.md", "index.md"} {
		p := filepath.Join(base, name)
		if data, err := os.ReadFile(p); err == nil {
			manifest, manifestBytes, manifestErr = p, data, nil
			break
		}
	}
	if manifestErr == nil {
		loc := parseLocations(string(manifestBytes))
		// A dashboard: key opts into the single-file adapter and short-circuits
		// planning/progress resolution (ADR-0004).
		if v, ok := loc["dashboard"]; ok {
			return dashboardAreas(base, path, v)
		}
		if v, ok := loc["planning"]; ok {
			planning = filepath.Join(path, v)
		}
		if v, ok := loc["progress"]; ok {
			progress = filepath.Join(path, v)
		}
	}

	if fi, err := os.Stat(planning); err != nil || !fi.IsDir() {
		if manifestErr == nil {
			if sys := parseSystems(string(manifestBytes)); len(sys) > 0 {
				return Areas{}, &RootManifestError{
					Path:         path,
					Dir:          dir,
					ManifestPath: manifest,
					PlanningDir:  planning,
					Systems:      sys,
				}
			}
		}
		return Areas{}, fmt.Errorf("no planning area at %s", planning)
	}

	return Areas{KnowledgeDir: base, PlanningDir: planning, ProgressDir: progress}, nil
}

// ResolveDashboard builds dashboard-mode Areas for an explicit dashboard file
// (the --dashboard flag), bypassing manifest/Locations detection. rel is
// resolved against path (repo-root-relative, matching Locations values), and
// only the path-is-a-directory and file-exists checks apply — a dashboard KB
// need not have a planning area at all (ADR-0004).
func ResolveDashboard(path, dir, rel string) (Areas, error) {
	if fi, err := os.Stat(path); err != nil || !fi.IsDir() {
		return Areas{}, fmt.Errorf("path is not a directory: %s", path)
	}
	return dashboardAreas(filepath.Join(path, dir), path, rel)
}

// dashboardAreas builds Areas for the single-file dashboard adapter. rel is the
// dashboard file's path (repo-root-relative, like other Locations values), so it
// is resolved against path, not base.
func dashboardAreas(base, path, rel string) (Areas, error) {
	file := filepath.Join(path, rel)
	if fi, err := os.Stat(file); err != nil || fi.IsDir() {
		return Areas{}, fmt.Errorf("no dashboard file at %s", file)
	}
	return Areas{KnowledgeDir: base, DashboardFile: file}, nil
}

// parseLocations extracts the planning/progress entries of a "## Locations"
// block: enter on the "## Locations" heading, leave at the next "## " heading,
// and in between split each non-empty line on its first ":" (spec §3). Only the
// planning and progress keys are kept, so prose lines are ignored.
func parseLocations(md string) map[string]string {
	out := map[string]string{}
	inBlock := false
	for _, line := range strings.Split(md, "\n") {
		if strings.HasPrefix(line, "## ") {
			inBlock = strings.EqualFold(strings.TrimSpace(line[len("## "):]), "Locations")
			continue
		}
		if !inBlock {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		i := strings.Index(trimmed, ":")
		if i < 0 {
			continue
		}
		// Keys are matched exactly (lowercase) — the framework skeleton always
		// writes "planning:" / "progress:", and exact matching keeps a prose
		// line that happens to contain a colon from being mistaken for a key.
		// (The heading match above is intentionally case-insensitive; the block
		// name is cosmetic, the keys are the contract.) A block that lists only
		// one of the two keys overrides just that area; the other keeps its
		// default in Resolve.
		key := strings.TrimSpace(trimmed[:i])
		if key == "planning" || key == "progress" || key == "dashboard" {
			out[key] = strings.TrimSpace(trimmed[i+1:])
		}
	}
	return out
}

// parseSystems collects distinct systems/<name> directory names from the cells
// of |-delimited table rows (spec §3 system-index parse).
func parseSystems(md string) []string {
	var out []string
	seen := map[string]bool{}
	for _, line := range strings.Split(md, "\n") {
		if !strings.Contains(line, "|") {
			continue
		}
		for _, m := range systemRe.FindAllString(line, -1) {
			if !seen[m] {
				seen[m] = true
				out = append(out, m)
			}
		}
	}
	return out
}

// DisplayPath renders a resolved path for the startup echo, restoring the "./"
// prefix the spec's startup block shows for plain relative paths (filepath.Join
// strips it). Absolute paths and paths already starting with "." pass through.
func DisplayPath(p string) string {
	if filepath.IsAbs(p) || strings.HasPrefix(p, ".") {
		return p
	}
	return "./" + p
}
