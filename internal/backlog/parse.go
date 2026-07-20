package backlog

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Parser turns a resolved instance into a Board. The line scanner below is the
// default; the interface is the seam decision D4 records, so a goldmark-backed
// parser can replace it later without touching callers.
type Parser interface {
	Parse(a Areas) (Board, error)
}

// NewParser returns the default line-oriented parser (spec §2).
func NewParser() Parser { return lineParser{} }

type lineParser struct{}

const emDash = " — " // U+2014 with surrounding spaces — the field separator (spec §2)

var (
	idRe       = regexp.MustCompile(`^([A-Z]+)-(\d+)`)
	checkboxRe = regexp.MustCompile(`^-\s+\[([ xX])\]\s+(.*)$`)
	bulletRe   = regexp.MustCompile(`^-\s+(.*)$`)
)

func (lineParser) Parse(a Areas) (Board, error) {
	b := Board{Meta: Meta{
		KnowledgeDir: a.KnowledgeDir,
		PlanningDir:  a.PlanningDir,
		ProgressDir:  a.ProgressDir,
	}}

	current := parseCurrent(readArea(a.PlanningDir, "CURRENT.md", &b.Warnings))
	backlog := parseBacklog(readArea(a.PlanningDir, "BACKLOG.md", &b.Warnings))
	done, doneWarns := parseDone(readArea(a.PlanningDir, "DONE.md", &b.Warnings))
	b.Warnings = append(b.Warnings, doneWarns...)

	cards, recWarns := reconcile(current, backlog, done)
	b.Cards = cards
	b.Warnings = append(b.Warnings, recWarns...)

	// Blockers + badge join land in TASK-004; progress area is read there.
	return b, nil
}

// readArea reads a planning/progress file; a read error is a warning, not a
// crash — the board degrades rather than blanking (spec §2, parse defensively).
func readArea(dir, name string, warnings *[]string) string {
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		*warnings = append(*warnings, "could not read "+name+": "+err.Error())
		return ""
	}
	return string(data)
}

// --- sectioning -------------------------------------------------------------

type section struct {
	name  string
	lines []string
}

// splitSections splits markdown on "## " headings. The preamble (before the
// first heading) is the section named "". "### " (h3) lines are content, so
// CURRENT task headers survive intact.
func splitSections(md string) []section {
	var secs []section
	cur := section{}
	for _, line := range strings.Split(md, "\n") {
		if strings.HasPrefix(line, "## ") {
			secs = append(secs, cur)
			cur = section{name: strings.TrimSpace(line[len("## "):])}
			continue
		}
		cur.lines = append(cur.lines, line)
	}
	return append(secs, cur)
}

// --- id / title -------------------------------------------------------------

func parseID(s string) *ID {
	m := idRe.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return nil
	}
	n, _ := strconv.Atoi(m[2]) // regex guarantees digits
	return &ID{Namespace: m[1], Number: n, Raw: m[0]}
}

// parseIDTitle splits "<ID> — <title>" (or a freeform title). If there is an id
// but no separator, the title is whatever trails the id token.
func parseIDTitle(s string) (*ID, string) {
	s = strings.TrimSpace(s)
	id := parseID(s)
	if i := strings.Index(s, emDash); i >= 0 {
		return id, strings.TrimSpace(s[i+len(emDash):])
	}
	if id != nil {
		return id, strings.TrimSpace(s[len(id.Raw):])
	}
	return id, s
}

// --- CURRENT ----------------------------------------------------------------

// parseCurrent reads the "## Active" section's "### <ID> — <title>" task blocks
// (spec §2). Prose and "## Entry template" are ignored.
func parseCurrent(md string) []Card {
	var cards []Card
	for _, sec := range splitSections(md) {
		if sec.name != "Active" {
			continue
		}
		for _, block := range splitBlocks(sec.lines, "### ") {
			cards = append(cards, parseCurrentTask(block))
		}
	}
	return cards
}

// splitBlocks groups lines into blocks each starting at a line with prefix;
// lines before the first such line are dropped.
func splitBlocks(lines []string, prefix string) [][]string {
	var blocks [][]string
	var cur []string
	for _, ln := range lines {
		if strings.HasPrefix(ln, prefix) {
			if cur != nil {
				blocks = append(blocks, cur)
			}
			cur = []string{ln}
			continue
		}
		if cur != nil {
			cur = append(cur, ln)
		}
	}
	if cur != nil {
		blocks = append(blocks, cur)
	}
	return blocks
}

func parseCurrentTask(lines []string) Card {
	id, title := parseIDTitle(strings.TrimPrefix(lines[0], "### "))
	card := Card{ID: id, Title: title, Column: ColInProgress, Raw: strings.Join(lines, "\n")}

	field := "" // which field the current line belongs to
	buf := map[string][]string{}
	for _, ln := range lines[1:] {
		t := strings.TrimSpace(ln)
		switch {
		case strings.HasPrefix(t, "**Source:**"):
			field = "source"
			buf[field] = append(buf[field], strings.TrimSpace(strings.TrimPrefix(t, "**Source:**")))
		case strings.HasPrefix(t, "**Acceptance criteria:**"):
			field = "ac"
		case strings.HasPrefix(t, "**Delivery override:**"):
			field = "override"
			buf[field] = append(buf[field], strings.TrimSpace(strings.TrimPrefix(t, "**Delivery override:**")))
		case strings.HasPrefix(t, "**Notes:**"):
			field = "notes"
			buf[field] = append(buf[field], strings.TrimSpace(strings.TrimPrefix(t, "**Notes:**")))
		default:
			if field != "" {
				buf[field] = append(buf[field], ln)
			}
		}
	}
	card.Source = joinTrim(buf["source"])
	card.Notes = joinTrim(buf["notes"])
	card.DeliveryOverride = joinTrim(buf["override"])
	for _, ln := range buf["ac"] {
		if m := checkboxRe.FindStringSubmatch(strings.TrimSpace(ln)); m != nil {
			card.Acceptance = append(card.Acceptance, Criterion{
				Text:    strings.TrimSpace(m[2]),
				Checked: strings.EqualFold(m[1], "x"),
			})
		}
	}
	return card
}

func joinTrim(lines []string) string {
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// --- BACKLOG ----------------------------------------------------------------

type backlogItem struct {
	ID      *ID
	Title   string
	Checked bool
	Parking bool
	Raw     string
}

// parseBacklog reads "## Active" (checkbox items) and "## Parking lot"
// (checkbox or freeform bullets, id optional). The conventions blockquote in
// the preamble is ignored (spec §2).
func parseBacklog(md string) []backlogItem {
	var items []backlogItem
	for _, sec := range splitSections(md) {
		parking := false
		switch sec.name {
		case "Active":
		case "Parking lot":
			parking = true
		default:
			continue
		}
		for _, ln := range sec.lines {
			t := strings.TrimSpace(ln)
			if m := checkboxRe.FindStringSubmatch(t); m != nil {
				id, title := parseIDTitle(m[2])
				items = append(items, backlogItem{ID: id, Title: title, Checked: strings.EqualFold(m[1], "x"), Parking: parking, Raw: ln})
			} else if m := bulletRe.FindStringSubmatch(t); m != nil {
				id, title := parseIDTitle(m[1])
				items = append(items, backlogItem{ID: id, Title: title, Parking: parking, Raw: ln})
			}
		}
	}
	return items
}

// --- DONE -------------------------------------------------------------------

// parseDone reads "## Completed" entries of the shape
// "- <ID> — <title> — <date> — <summary> — <delivery record>. See journal …".
// Date is the 3rd field and the delivery record is the last, so a summary that
// itself contains " — " is still split correctly (spec §2 robustness).
func parseDone(md string) ([]Card, []string) {
	var cards []Card
	var warnings []string
	for _, sec := range splitSections(md) {
		if sec.name != "Completed" {
			continue
		}
		for _, ln := range sec.lines {
			t := strings.TrimSpace(ln)
			if !strings.HasPrefix(t, "- ") {
				continue
			}
			card, warn := parseDoneEntry(t, ln)
			cards = append(cards, card)
			if warn != "" {
				warnings = append(warnings, warn)
			}
		}
	}
	return cards, warnings
}

func parseDoneEntry(trimmed, raw string) (Card, string) {
	content := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))

	rec := &DoneRecord{}
	if i := strings.Index(content, "See journal"); i >= 0 {
		rec.JournalPointer = strings.TrimSpace(content[i:])
		content = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(content[:i]), "."))
	}

	parts := strings.Split(content, emDash)
	card := Card{Column: ColDone, Raw: raw, Done: rec}
	var warn string
	if len(parts) >= 1 {
		card.ID = parseID(parts[0])
	}
	if len(parts) >= 2 {
		card.Title = strings.TrimSpace(parts[1])
	}
	if len(parts) >= 3 {
		rec.Date = strings.TrimSpace(parts[2])
	}
	switch {
	case len(parts) >= 5:
		rec.DeliveryRecord = strings.TrimSpace(parts[len(parts)-1])
		rec.Summary = strings.TrimSpace(strings.Join(parts[3:len(parts)-1], emDash))
	case len(parts) == 4:
		rec.Summary = strings.TrimSpace(parts[3])
	}
	if len(parts) < 3 {
		warn = "malformed DONE entry (too few '—' fields): " + raw
	}
	return card, warn
}

// --- reconciliation & dedup (spec §5) ---------------------------------------

// reconcile merges CURRENT/BACKLOG/DONE into one card per id in its furthest
// column (Done > InProgress > Backlog) and emits the three hygiene warnings.
func reconcile(current []Card, backlog []backlogItem, done []Card) ([]Card, []string) {
	var warns []string
	var cards []Card
	seen := map[string]bool{} // ids already placed in their furthest column

	doneIDs := map[string]bool{}
	for _, d := range done {
		if d.ID != nil {
			doneIDs[d.ID.Raw] = true
		}
	}
	backlogChecked := map[string]bool{}
	for _, b := range backlog {
		if b.Checked && b.ID != nil {
			backlogChecked[b.ID.Raw] = true
		}
	}

	if len(current) > 1 {
		warns = append(warns, "CURRENT holds >1 active task (framework one-task invariant)")
	}

	// Done column: real DONE entries first.
	for _, d := range done {
		if d.ID != nil {
			seen[d.ID.Raw] = true
			if !backlogChecked[d.ID.Raw] {
				warns = append(warns, "DONE item not ticked in BACKLOG: "+d.ID.Raw)
			}
		}
		cards = append(cards, d)
	}
	// A [x] backlog item with no DONE entry is still shown in Done, with a warning.
	for _, b := range backlog {
		if b.Parking || !b.Checked || b.ID == nil || doneIDs[b.ID.Raw] {
			continue
		}
		warns = append(warns, "shipped item missing from DONE.md: "+b.ID.Raw)
		if !seen[b.ID.Raw] {
			seen[b.ID.Raw] = true
			cards = append(cards, Card{ID: b.ID, Title: b.Title, Column: ColDone, Raw: b.Raw})
		}
	}
	// In Progress: current tasks not already claimed by Done.
	for _, c := range current {
		if c.ID != nil {
			if seen[c.ID.Raw] {
				continue
			}
			seen[c.ID.Raw] = true
		}
		cards = append(cards, c)
	}
	// Backlog: remaining active items (unchecked, not deduped away).
	for _, b := range backlog {
		if b.Parking {
			continue
		}
		if b.ID != nil {
			if seen[b.ID.Raw] {
				continue
			}
			seen[b.ID.Raw] = true
		}
		cards = append(cards, Card{ID: b.ID, Title: b.Title, Column: ColBacklog, Raw: b.Raw})
	}
	// Parking-lot items: Backlog column, muted; never deduped (no id).
	for _, b := range backlog {
		if !b.Parking {
			continue
		}
		cards = append(cards, Card{ID: b.ID, Title: b.Title, Column: ColBacklog, ParkingLot: true, Raw: b.Raw})
	}
	return cards, warns
}
