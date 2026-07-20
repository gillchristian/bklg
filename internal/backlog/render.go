package backlog

import (
	"embed"
	"html/template"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

//go:embed templates
var templatesFS embed.FS

// funcs are the template helpers. Styling lives here (near rendering), not on
// the model.
var funcs = template.FuncMap{
	"badgeClass": badgeClass,
	"badgeText":  badgeText,
	"truncate":   truncate,
	"md":         renderMarkdown,
}

// truncate shortens a card title for the board so a long, paragraph-sized title
// (real instances put whole records on a line) doesn't become a wall of text.
// Rune-based (never splits a multibyte rune), with a small word-boundary backoff,
// and appends an ellipsis. The full title stays on the detail page.
func truncate(s string, max int) string {
	s = strings.TrimSpace(s)
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	cut := max
	for i := cut; i > max-24 && i > 0; i-- { // back up to a space for a clean break
		if runes[i] == ' ' {
			cut = i
			break
		}
	}
	return strings.TrimSpace(string(runes[:cut])) + "…"
}

// One template set per page (each defines its own "content"), sharing the
// layout — separate sets avoid a "content" redefinition collision.
var (
	boardTmpl = template.Must(
		template.New("board").Funcs(funcs).ParseFS(templatesFS, "templates/layout.html", "templates/board.html"),
	)
	taskTmpl = template.Must(
		template.New("task").Funcs(funcs).ParseFS(templatesFS, "templates/layout.html", "templates/task.html"),
	)
	diagTmpl = template.Must(
		template.New("diag").Funcs(funcs).ParseFS(templatesFS, "templates/layout.html", "templates/diag.html"),
	)
)

// versionString renders a freshness stamp for /_v and the poll script: the max
// mtime as UnixNano, or "0" when there are no files. The page bakes this in and
// /_v returns the same computation, so equality means "unchanged".
func versionString(t time.Time) string {
	if t.IsZero() {
		return "0"
	}
	return strconv.FormatInt(t.UnixNano(), 10)
}

// boardVM is the board's view model: the three columns in flow order, plus the
// header/banner inputs and the live-reload version.
type boardVM struct {
	PlanningDir string
	Warnings    []Warning
	Version     string
	Columns     []columnVM
	Systems     []systemChip // multi-system filter bar; empty in single-system mode
}

// systemChip is one entry in the multi-system filter bar.
type systemChip struct {
	Name   string
	Count  int
	Active bool
	Href   string
}

type columnVM struct {
	Title string
	Cards []Card
}

// taskVM is the detail page's view model. It carries the shared header inputs
// (PlanningDir/Warnings, so the layout banner renders) plus the card and the
// blockers referencing it (open first, then resolved).
type taskVM struct {
	PlanningDir string
	Warnings    []Warning
	Version     string
	Card        Card
	Blockers    []Blocker
}

// viewModel builds the board view. systemFilter (from ?system=) restricts the
// cards to one system in aggregate mode; it is a no-op in single-system mode.
// allSystems is the full set of discovered systems (empty in single mode) — the
// filter bar lists every one, including those with zero cards, so any project is
// filterable and the bar matches the startup's system count.
func viewModel(b Board, systemFilter string, allSystems []string) boardVM {
	counts := map[string]int{}
	for _, c := range b.Cards {
		if c.System != "" {
			counts[c.System]++
		}
	}

	cols := []columnVM{{Title: "Backlog"}, {Title: "In Progress"}, {Title: "Done"}}
	for _, c := range b.Cards {
		if systemFilter != "" && c.System != systemFilter {
			continue
		}
		switch c.Column {
		case ColBacklog:
			cols[0].Cards = append(cols[0].Cards, c)
		case ColInProgress:
			cols[1].Cards = append(cols[1].Cards, c)
		case ColDone:
			cols[2].Cards = append(cols[2].Cards, c)
		}
	}

	var chips []systemChip
	if len(allSystems) > 0 { // aggregate mode
		chips = append(chips, systemChip{Name: "All", Count: len(b.Cards), Active: systemFilter == "", Href: "/"})
		for _, name := range allSystems {
			chips = append(chips, systemChip{Name: name, Count: counts[name], Active: systemFilter == name, Href: "/?system=" + url.QueryEscape(name)})
		}
	}

	return boardVM{
		PlanningDir: DisplayPath(b.Meta.PlanningDir),
		Warnings:    b.Warnings,
		Version:     versionString(b.Meta.LatestMTime),
		Columns:     cols,
		Systems:     chips,
	}
}

// --- /_diag view model ------------------------------------------------------

type diagGroup struct {
	Kind        string
	Title       string
	Explanation string
	Warnings    []Warning
}

type diagVM struct {
	PlanningDir string
	Version     string
	Warnings    []Warning // for the shared layout banner (count)
	Groups      []diagGroup
}

// diagKindOrder fixes the display order of known warning kinds (most
// invariant-violating first); explanations say what each means and what to do.
var diagKindOrder = []string{"current-multiple", "shipped-missing-done", "done-not-ticked", "malformed-done", "read-error"}

var diagExplain = map[string][2]string{
	"current-multiple":     {"More than one active task", "CURRENT holds more than one task under ## Active. The framework's one-task-at-a-time invariant expects exactly one — finish or move the extras."},
	"shipped-missing-done": {"Shipped, but not in DONE.md", "Checked [x] in BACKLOG but with no DONE.md entry. Add a DONE entry — or, if this instance keeps the full record inline on the [x] line, this is expected and purely informational."},
	"done-not-ticked":      {"In DONE.md, not ticked in BACKLOG", "A DONE.md entry whose id isn't checked [x] in BACKLOG. Tick it so the backlog stays an accurate shipping index."},
	"malformed-done":       {"Malformed DONE entry", "A DONE.md entry didn't match either supported shape; it was parsed best-effort."},
	"read-error":           {"Couldn't read a file", "A planning/progress file couldn't be read; the board renders best-effort without it."},
}

func buildDiagVM(b Board) diagVM {
	byKind := map[string][]Warning{}
	for _, w := range b.Warnings {
		byKind[w.Kind] = append(byKind[w.Kind], w)
	}
	var groups []diagGroup
	seen := map[string]bool{}
	for _, k := range diagKindOrder {
		if ws := byKind[k]; len(ws) > 0 {
			groups = append(groups, diagGroup{Kind: k, Title: diagExplain[k][0], Explanation: diagExplain[k][1], Warnings: ws})
			seen[k] = true
		}
	}
	var unknown []string // any future kinds not in diagKindOrder
	for k := range byKind {
		if !seen[k] {
			unknown = append(unknown, k)
		}
	}
	sort.Strings(unknown) // deterministic order (map iteration is randomized)
	for _, k := range unknown {
		groups = append(groups, diagGroup{Kind: k, Title: k, Warnings: byKind[k]})
	}
	return diagVM{
		PlanningDir: DisplayPath(b.Meta.PlanningDir),
		Version:     versionString(b.Meta.LatestMTime),
		Warnings:    b.Warnings,
		Groups:      groups,
	}
}

// badgeClass maps a badge kind to Tailwind chip classes (spec §6 suggested
// styles). blocked is the load-bearing red.
func badgeClass(kind string) string {
	switch kind {
	case "blocked":
		return "bg-red-100 text-red-800 ring-red-600/20"
	case "parking":
		return "bg-slate-100 text-slate-600 ring-slate-500/20"
	case "override":
		return "bg-amber-100 text-amber-800 ring-amber-600/20"
	case "no-ac":
		return "bg-yellow-100 text-yellow-800 ring-yellow-600/20"
	default: // namespace + anything else
		return "bg-gray-100 text-gray-700 ring-gray-500/20"
	}
}

// badgeText is the chip label: the namespace value for the namespace chip, the
// kind otherwise.
func badgeText(b Badge) string {
	if b.Kind == "namespace" {
		return b.Label
	}
	return b.Kind
}
