package backlog

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// Server serves the board for either one resolved instance (single mode) or a
// monorepo root aggregating every systems/<name> instance (multi mode). It
// re-parses on every request (spec §7 freshness).
type Server struct {
	parser Parser

	areas Areas // single mode

	multi       bool     // aggregate mode
	rootPath    string   // multi: the repo-root argument
	systemPaths []string // multi: repo-root-relative, e.g. ["systems/track", …]
}

// NewServer builds a single-instance Server.
func NewServer(a Areas) *Server {
	return &Server{parser: NewParser(), areas: a}
}

// NewMultiServer builds an aggregate Server over the systems discovered in a
// root manifest. systemPaths are repo-root-relative (e.g. "systems/track").
func NewMultiServer(rootPath string, systemPaths []string) *Server {
	return &Server{parser: NewParser(), multi: true, rootPath: rootPath, systemPaths: systemPaths}
}

// Systems returns the aggregated systems' display names (multi mode).
func (s *Server) Systems() []string {
	names := make([]string, len(s.systemPaths))
	for i, sp := range s.systemPaths {
		names[i] = systemName(sp)
	}
	return names
}

func systemName(systemPath string) string { return filepath.Base(systemPath) }

// board parses fresh. In multi mode it resolves + parses each system, tags each
// card with its system, and concatenates — no cross-system dedup (namespaces
// differ per system and each system's board is already deduped). A system that
// fails to resolve is skipped with a warning, not a crash.
func (s *Server) board() Board {
	if !s.multi {
		b, _ := s.parser.Parse(s.areas)
		return b
	}
	agg := Board{}
	agg.Meta.PlanningDir = s.rootPath
	for _, sp := range s.systemPaths {
		name := systemName(sp)
		a, err := Resolve(s.rootPath, sp+"/knowledge")
		if err != nil {
			agg.Warnings = append(agg.Warnings, Warning{Kind: "read-error", Message: "system " + name + " did not resolve: " + err.Error()})
			continue
		}
		b, _ := s.parser.Parse(a)
		for i := range b.Cards {
			b.Cards[i].System = name
		}
		agg.Cards = append(agg.Cards, b.Cards...)
		agg.Blockers = append(agg.Blockers, b.Blockers...)
		agg.Warnings = append(agg.Warnings, b.Warnings...)
		if b.Meta.LatestMTime.After(agg.Meta.LatestMTime) {
			agg.Meta.LatestMTime = b.Meta.LatestMTime
		}
	}
	return agg
}

// latestMTime is the /_v freshness stamp — max mtime across the parsed files,
// statted (not reparsed) so polling stays cheap. Aggregates in multi mode.
func (s *Server) latestMTime() time.Time {
	if !s.multi {
		return areaMTime(s.areas)
	}
	var latest time.Time
	for _, sp := range s.systemPaths {
		if a, err := Resolve(s.rootPath, sp+"/knowledge"); err == nil {
			if m := areaMTime(a); m.After(latest) {
				latest = m
			}
		}
	}
	return latest
}

// Routes wires the mux. Literal segments beat the {id} wildcard, so /_v and
// /_diag win over /{id} without ordering tricks (spec §7). /{$} matches "/".
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.handleBoard)
	mux.HandleFunc("GET /_v", s.handleVersion)
	mux.HandleFunc("GET /_diag", s.handleDiag)
	mux.HandleFunc("GET /{id}", s.handleTask)
	return mux
}

// handleVersion returns the max mtime across the parsed files as a bare integer
// (spec §7) — cheap, so the page can poll it and reload on change.
func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, versionString(s.latestMTime()))
}

// handleBoard renders the board. The optional ?system=<name> query filters the
// aggregate to one system (ignored in single mode). Buffered so a template
// error becomes a clean 500, not a half-written page.
func (s *Server) handleBoard(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	if err := boardTmpl.ExecuteTemplate(&buf, "layout", viewModel(s.board(), r.URL.Query().Get("system"), s.Systems())); err != nil {
		http.Error(w, "bklg: render error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

// handleTask renders one task's detail page (spec §7). The id is matched
// case-insensitively across all (aggregated) cards; an unknown id (or a
// parking/id-less card) is a 404.
func (s *Server) handleTask(w http.ResponseWriter, r *http.Request) {
	b := s.board()
	card, ok := b.CardByRoute(r.PathValue("id"))
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Blockers referencing this card, open first then resolved.
	var refs []Blocker
	for _, open := range []bool{true, false} {
		for _, bl := range b.Blockers {
			if bl.Open != open {
				continue
			}
			if pid := parseID(strings.ToUpper(bl.TaskRaw)); pid != nil && card.ID != nil && pid.Raw == strings.ToUpper(card.ID.Raw) {
				refs = append(refs, bl)
			}
		}
	}

	vm := taskVM{
		PlanningDir: DisplayPath(b.Meta.PlanningDir),
		Warnings:    b.Warnings,
		Version:     versionString(b.Meta.LatestMTime),
		Card:        *card,
		Blockers:    refs,
		LinearBase:  b.Meta.LinearBase,
	}
	var buf bytes.Buffer
	if err := taskTmpl.ExecuteTemplate(&buf, "layout", vm); err != nil {
		http.Error(w, "bklg: render error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

// handleDiag renders the diagnostics page: warnings grouped by kind, each group
// with a count, an explanation of what it means + how to fix it, and links from
// id-bearing warnings to their task detail page (spec §7, made actionable).
func (s *Server) handleDiag(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	if err := diagTmpl.ExecuteTemplate(&buf, "layout", buildDiagVM(s.board())); err != nil {
		http.Error(w, "bklg: render error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}
