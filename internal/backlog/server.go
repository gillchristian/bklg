package backlog

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

// Server serves the board for one resolved instance. It re-parses on every
// request (spec §7 freshness) so the board tracks a live agent session.
type Server struct {
	areas  Areas
	parser Parser
}

// NewServer builds a Server for the resolved areas.
func NewServer(a Areas) *Server {
	return &Server{areas: a, parser: NewParser()}
}

// board parses fresh. Parse never returns a hard error — unreadable files and
// malformed entries surface as Board.Warnings (spec §2), so the board renders
// rather than blanking.
func (s *Server) board() Board {
	b, _ := s.parser.Parse(s.areas)
	return b
}

// Routes wires the mux. Literal segments beat the {id} wildcard, so /_diag
// (and /_v, added in TASK-007) win over /{id} without ordering tricks (spec §7).
// /{$} matches exactly "/".
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.handleBoard)
	mux.HandleFunc("GET /_v", s.handleVersion)
	mux.HandleFunc("GET /_diag", s.handleDiag)
	mux.HandleFunc("GET /{id}", s.handleTask)
	return mux
}

// handleVersion returns the max mtime across the parsed files as a bare integer
// (spec §7). Cheap — it stats the files rather than reparsing — so the page can
// poll it every few seconds and reload on change.
func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, versionString(areaMTime(s.areas)))
}

// handleBoard renders the three-column board. It renders into a buffer first so
// a template error becomes a clean 500 instead of a half-written page.
func (s *Server) handleBoard(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	if err := boardTmpl.ExecuteTemplate(&buf, "layout", viewModel(s.board())); err != nil {
		http.Error(w, "bklg: render error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

// handleTask renders one task's detail page (spec §7). The id is matched
// case-insensitively; an unknown id (or a parking/id-less card) is a 404.
func (s *Server) handleTask(w http.ResponseWriter, r *http.Request) {
	b := s.board()
	card, ok := b.CardByRawID(r.PathValue("id"))
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
	}
	var buf bytes.Buffer
	if err := taskTmpl.ExecuteTemplate(&buf, "layout", vm); err != nil {
		http.Error(w, "bklg: render error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

// handleDiag lists the parse warnings verbatim, one per line (spec §7). Plain
// text — repo content is never rendered as HTML.
func (s *Server) handleDiag(w http.ResponseWriter, r *http.Request) {
	b := s.board()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if len(b.Warnings) == 0 {
		fmt.Fprintln(w, "no warnings")
		return
	}
	for _, warn := range b.Warnings {
		fmt.Fprintln(w, warn)
	}
}
