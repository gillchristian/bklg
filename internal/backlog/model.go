package backlog

import (
	"strings"
	"time"
)

// ID is a task identifier like {"MONO", 6, "MONO-006"}. Raw is the source text,
// Namespace/Number the parsed halves. IDs match [A-Z]+-\d+ (spec §2/§4).
type ID struct {
	Namespace string
	Number    int
	Raw       string
}

// Column is a board column; a card lives in exactly one (spec §5).
type Column int

const (
	ColBacklog Column = iota
	ColInProgress
	ColDone
)

func (c Column) String() string {
	switch c {
	case ColBacklog:
		return "Backlog"
	case ColInProgress:
		return "In Progress"
	case ColDone:
		return "Done"
	default:
		return "?"
	}
}

// Criterion is one acceptance-criteria checklist item from a CURRENT task.
type Criterion struct {
	Text    string
	Checked bool
}

// DoneRecord is the parsed shape of a DONE.md entry's tail fields (spec §4).
type DoneRecord struct {
	Date           string
	Summary        string
	DeliveryRecord string
	JournalPointer string
}

// Badge is a cross-cutting state chip on a card (spec §6). Kind is one of
// blocked | parking | override | no-ac | namespace; Label is the display text
// (for the namespace chip it is the namespace itself). Populated when the board
// is assembled (blocked needs the blockers join — TASK-004).
type Badge struct {
	Kind  string
	Label string
}

// Card is one task as rendered on the board, derived by joining an id across the
// parsed files (spec §4). ID is nil for parking-lot items without an id.
type Card struct {
	ID     *ID
	Title  string
	Column Column
	Badges []Badge

	Source           string      // from CURRENT
	Notes            string      // from CURRENT
	DeliveryOverride string      // from CURRENT
	Acceptance       []Criterion // from CURRENT
	Done             *DoneRecord // from DONE
	ParkingLot       bool
	System           string // owning system in aggregate/multi-system mode; "" single-system
	Raw              string // the source block, shown verbatim on the detail page

	// Dashboard-adapter fields (ADR-0004); zero-valued for framework-mode cards.
	Dashboard bool   // parsed from a single-file dashboard, not the planning skeleton
	Subtitle  string // the row/bullet text after the leading **bold** title (post em-dash)
	Tickets   []ID   // inline Linear-style ids found in the row (0..N), first-seen order
	Group     string // backlog group label (e.g. "Product / code"); "" if ungrouped
	Material  string // the row's Material/Record cell (raw markdown, e.g. "[x](y)")
	Status    string // the row's Status/next-step cell (Active), unescaped; "" for Done/Backlog
	Blocked   bool   // a leading ⛔ / **Blocked** marker (dashboard's blocked signal)
}

// Blocker is a parsed blockers.md entry (spec §4). Populated in TASK-004.
type Blocker struct {
	ID      ID
	Title   string
	Opened  string
	Body    string
	TaskRaw string // affected task id, raw ("TRAIL-007") — the join key
	Open    bool
}

// Meta carries the resolved locations and the freshness stamp for /_v.
type Meta struct {
	KnowledgeDir string
	PlanningDir  string
	ProgressDir  string
	LatestMTime  time.Time
	LinearBase   string // ticket-chip URL prefix in dashboard mode (ADR-0004)
}

// Warning is a parse/reconciliation diagnostic surfaced on /_diag. Kind groups
// like warnings; TaskRaw (when set) links the warning to a task's detail page.
type Warning struct {
	Kind    string // shipped-missing-done | done-not-ticked | current-multiple | read-error | malformed-done | dashboard-malformed
	Message string
	TaskRaw string // id involved, for linking (optional)
}

// Board is the whole parsed model behind every page (spec §4).
type Board struct {
	Cards    []Card
	Blockers []Blocker
	Warnings []Warning // parse + reconciliation diagnostics -> /_diag
	Meta     Meta
}

// CardByRawID returns the card whose id matches raw (case-insensitively, per the
// D2 route-key decision). Parking/id-less cards have no id and never match, so
// they have no detail page.
func (b *Board) CardByRawID(raw string) (*Card, bool) {
	for i := range b.Cards {
		if b.Cards[i].ID != nil && strings.EqualFold(b.Cards[i].ID.Raw, raw) {
			return &b.Cards[i], true
		}
	}
	return nil, false
}
