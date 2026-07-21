package backlog

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Parser turns a resolved instance into a Board. The line scanner below is the
// default; the interface is the seam decision D4 records, so a goldmark-backed
// parser can replace it later without touching callers.
type Parser interface {
	Parse(a Areas) (Board, error)
}

// NewParser returns the default parser. It dispatches on the resolved Areas: a
// dashboard-mode instance (Areas.DashboardFile set, ADR-0004) uses the
// single-file dashboard adapter; otherwise the line-oriented framework parser
// (spec §2). The seam is decision D4.
func NewParser() Parser { return defaultParser{} }

type defaultParser struct{}

func (defaultParser) Parse(a Areas) (Board, error) {
	if a.DashboardFile != "" {
		return parseDashboard(a)
	}
	return lineParser{}.Parse(a)
}

var (
	leadingBoldRe = regexp.MustCompile(`^\*\*(.+?)\*\*`)
	ticketRe      = regexp.MustCompile(`[A-Z]+-\d+`) // inline Linear-style ids
	groupHeadRe   = regexp.MustCompile(`^\*\*(.+):\*\*$`)
	tableSepRe    = regexp.MustCompile(`^:?-+:?$`)
)

// parseDashboard reads a single-file Active/Backlog/Done dashboard into a Board
// (ADR-0004; contract in reference/specs/dashboard-format.md). Active/Done are
// pipe tables, Backlog is bullet groups. Parsing is defensive: a row it can't
// read becomes a warning and is skipped, never a crash (spec §2).
func parseDashboard(a Areas) (Board, error) {
	b := Board{Meta: Meta{
		KnowledgeDir: a.KnowledgeDir,
		PlanningDir:  a.DashboardFile, // shown in the header + startup echo
		LatestMTime:  areaMTime(a),
		LinearBase:   a.LinkBase,
	}}
	data, err := os.ReadFile(a.DashboardFile)
	if err != nil {
		b.Warnings = append(b.Warnings, Warning{Kind: "read-error", Message: "could not read dashboard " + a.DashboardFile + ": " + err.Error()})
		return b, nil
	}
	for _, sec := range splitSections(string(data)) {
		switch strings.ToLower(strings.TrimSpace(sec.name)) {
		case "active":
			b.Cards = append(b.Cards, parseDashTable(sec.lines, ColInProgress, &b.Warnings)...)
		case "done":
			b.Cards = append(b.Cards, parseDashTable(sec.lines, ColDone, &b.Warnings)...)
		case "backlog":
			b.Cards = append(b.Cards, parseDashBacklog(sec.lines)...)
		}
	}
	computeDashboardBadges(b.Cards)
	assignSlugs(b.Cards)
	return b, nil
}

// assignSlugs gives each card a unique url-safe slug from its title (the
// dashboard detail-route key). Collisions get a -2, -3 … suffix; the used-set
// guard keeps a title that literally slugifies to "foo-2" from clashing.
func assignSlugs(cards []Card) {
	used := map[string]bool{}
	for i := range cards {
		base := slugify(cards[i].Title)
		if base == "" {
			base = "card"
		}
		s := base
		for n := 2; used[s]; n++ {
			s = base + "-" + strconv.Itoa(n)
		}
		used[s] = true
		cards[i].Slug = s
	}
}

// slugify lowercases s and collapses each run of non-alphanumeric characters to
// a single hyphen, trimming leading/trailing hyphens.
func slugify(s string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(s) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			dash = false
		} else if b.Len() > 0 && !dash {
			b.WriteByte('-')
			dash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

// computeDashboardBadges attaches the dashboard-mode chips: a red blocked badge
// and a group chip. The framework computeBadges is deliberately NOT used — its
// no-ac hygiene and blocker join don't apply to dashboard cards (ADR-0004).
func computeDashboardBadges(cards []Card) {
	for i := range cards {
		c := &cards[i]
		var badges []Badge
		if c.Blocked {
			badges = append(badges, Badge{Kind: "blocked"})
		}
		if c.Group != "" {
			badges = append(badges, Badge{Kind: "group", Label: c.Group})
		}
		c.Badges = badges
	}
}

// parseDashTable parses a pipe-table section (Active/Done). The first pipe row
// is the header (mapping the title/material/status columns by name), the dashes
// row is skipped, and each remaining row is a card. Cells split on unescaped
// "|"; "\|" is a literal pipe.
func parseDashTable(lines []string, col Column, warnings *[]Warning) []Card {
	var cards []Card
	titleIdx, materialIdx, statusIdx := 0, -1, -1
	haveHeader := false
	for _, ln := range lines {
		t := strings.TrimSpace(ln)
		if !strings.HasPrefix(t, "|") {
			continue
		}
		cells := splitCells(t)
		if isSeparatorRow(cells) {
			continue
		}
		if !haveHeader {
			titleIdx, materialIdx, statusIdx = mapDashColumns(cells)
			haveHeader = true
			continue
		}
		titleCell := cellAt(cells, titleIdx)
		if strings.TrimSpace(titleCell) == "" {
			*warnings = append(*warnings, Warning{Kind: "dashboard-malformed", Message: "dashboard row with an empty title cell: " + t})
			continue
		}
		title, subtitle := splitDashTitle(titleCell)
		card := Card{
			Dashboard: true,
			Column:    col,
			Title:     title,
			Subtitle:  subtitle,
			Material:  cellAt(cells, materialIdx),
			Status:    cellAt(cells, statusIdx),
			Tickets:   findTickets(t),
			Raw:       t,
		}
		if col != ColDone {
			card.Blocked = hasBlockedMarker(card.Status)
		}
		cards = append(cards, card)
	}
	return cards
}

// parseDashBacklog parses the Backlog section's bullet groups: a "**Group:**"
// line sets the current group label; each "- " bullet under it is a card.
func parseDashBacklog(lines []string) []Card {
	var cards []Card
	group := ""
	for _, ln := range lines {
		t := strings.TrimSpace(ln)
		if m := groupHeadRe.FindStringSubmatch(t); m != nil {
			group = strings.TrimSpace(m[1])
			continue
		}
		if !strings.HasPrefix(t, "- ") {
			continue
		}
		content := strings.TrimSpace(strings.TrimPrefix(t, "- "))
		if content == "" {
			continue
		}
		blocked := hasBlockedMarker(content)
		content = strings.TrimSpace(strings.TrimPrefix(content, "⛔")) // clean the title source
		title, subtitle := splitDashTitle(content)
		cards = append(cards, Card{
			Dashboard: true,
			Column:    ColBacklog,
			Title:     title,
			Subtitle:  subtitle,
			Group:     group,
			Tickets:   findTickets(t),
			Blocked:   blocked,
			Raw:       t,
		})
	}
	return cards
}

// splitDashTitle splits a row/bullet's title cell into the leading **bold**
// title and the subtitle after the em-dash. With no leading bold, the whole
// cell is the title (still splitting off an em-dash subtitle if present).
func splitDashTitle(s string) (title, subtitle string) {
	s = strings.TrimSpace(s)
	// Search for the subtitle separator only in the text AFTER the bold title,
	// so an em-dash inside the bold phrase ("**Foo — bar** — sub") doesn't split
	// the subtitle in the wrong place.
	rest := s
	if m := leadingBoldRe.FindStringSubmatch(s); m != nil {
		title = strings.TrimSpace(m[1])
		rest = s[len(m[0]):]
	}
	if i := strings.Index(rest, emDash); i >= 0 {
		if title == "" {
			title = strings.TrimSpace(rest[:i])
		}
		return title, strings.TrimSpace(rest[i+len(emDash):])
	}
	if title == "" {
		title = strings.TrimSpace(rest)
	}
	return title, ""
}

// splitCells splits a markdown table row on unescaped "|", unescaping "\|" to a
// literal "|", and drops the empty cells produced by the bounding pipes.
func splitCells(row string) []string {
	var cells []string
	var b strings.Builder
	rs := []rune(row)
	for i := 0; i < len(rs); i++ {
		if rs[i] == '\\' && i+1 < len(rs) && rs[i+1] == '|' {
			b.WriteRune('|')
			i++
			continue
		}
		if rs[i] == '|' {
			cells = append(cells, strings.TrimSpace(b.String()))
			b.Reset()
			continue
		}
		b.WriteRune(rs[i])
	}
	cells = append(cells, strings.TrimSpace(b.String()))
	if len(cells) > 0 && cells[0] == "" {
		cells = cells[1:]
	}
	if len(cells) > 0 && cells[len(cells)-1] == "" {
		cells = cells[:len(cells)-1]
	}
	return cells
}

func isSeparatorRow(cells []string) bool {
	if len(cells) == 0 {
		return false
	}
	for _, c := range cells {
		if !tableSepRe.MatchString(strings.TrimSpace(c)) {
			return false
		}
	}
	return true
}

// mapDashColumns finds the title/material/status column indices by header name,
// defaulting title to the first column when no "Work"/"What" header is found.
func mapDashColumns(header []string) (titleIdx, materialIdx, statusIdx int) {
	titleIdx, materialIdx, statusIdx = 0, -1, -1
	for i, h := range header {
		switch hl := strings.ToLower(strings.TrimSpace(h)); {
		case hl == "work" || hl == "what":
			titleIdx = i
		case hl == "material" || hl == "record":
			materialIdx = i
		case strings.Contains(hl, "status"):
			statusIdx = i
		}
	}
	return
}

func cellAt(cells []string, i int) string {
	if i < 0 || i >= len(cells) {
		return ""
	}
	return cells[i]
}

// findTickets collects distinct inline Linear-style ids (e.g. PINATA-602) in
// first-seen order. GitHub PR refs (#123) never match the id pattern.
func findTickets(s string) []ID {
	var out []ID
	seen := map[string]bool{}
	for _, m := range ticketRe.FindAllString(s, -1) {
		if seen[m] {
			continue
		}
		seen[m] = true
		if id := parseID(m); id != nil {
			out = append(out, *id)
		}
	}
	return out
}

// hasBlockedMarker reports whether text leads with the dashboard blocked marker
// (a ⛔ or **Blocked); a mid-text ⛔ (a decision marker) does not count.
func hasBlockedMarker(s string) bool {
	s = strings.TrimSpace(s)
	return strings.HasPrefix(s, "⛔") || strings.HasPrefix(s, "**Blocked")
}

type lineParser struct{}

const emDash = " — " // U+2014 with surrounding spaces — the field separator (spec §2)

var (
	idRe          = regexp.MustCompile(`^([A-Z]+)-(\d+)`)
	checkboxRe    = regexp.MustCompile(`^-\s+\[([ xX])\]\s+(.*)$`)
	bulletRe      = regexp.MustCompile(`^-\s+(.*)$`)
	blockerHeadRe = regexp.MustCompile(`^BLOCKER-\d+`)
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
	b.Warnings = append(b.Warnings, recWarns...)

	b.Blockers = parseBlockers(readArea(a.ProgressDir, "blockers.md", &b.Warnings))
	computeBadges(cards, b.Blockers)
	b.Cards = cards
	b.Meta.LatestMTime = areaMTime(a)
	return b, nil
}

// areaMTime is the newest modification time across the parsed files — the
// freshness stamp behind /_v. Cheap (4 stats), so /_v can poll without a full
// parse. Missing files are skipped.
func areaMTime(a Areas) time.Time {
	// Dashboard mode watches its single file; framework mode the four planning
	// and progress files.
	files := []string{
		filepath.Join(a.PlanningDir, "CURRENT.md"),
		filepath.Join(a.PlanningDir, "BACKLOG.md"),
		filepath.Join(a.PlanningDir, "DONE.md"),
		filepath.Join(a.ProgressDir, "blockers.md"),
	}
	if a.DashboardFile != "" {
		files = []string{a.DashboardFile}
	}
	var latest time.Time
	for _, p := range files {
		if fi, err := os.Stat(p); err == nil {
			if m := fi.ModTime(); m.After(latest) {
				latest = m
			}
		}
	}
	return latest
}

// readArea reads a planning/progress file; a read error is a warning, not a
// crash — the board degrades rather than blanking (spec §2, parse defensively).
func readArea(dir, name string, warnings *[]Warning) string {
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		*warnings = append(*warnings, Warning{Kind: "read-error", Message: "could not read " + name + ": " + err.Error()})
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
	// Strip leading markdown decoration so an id written as **WI-8 — …** or
	// `TRACK-1` is still recognized (real instances decorate ids; spec §2's
	// skeleton does not).
	s = strings.TrimLeft(strings.TrimSpace(s), "*`_ ")
	m := idRe.FindStringSubmatch(s)
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
		// id.Raw may not sit at index 0 — parseID strips leading decoration — so
		// locate it in the stripped string rather than slicing from the start.
		rest := strings.TrimPrefix(strings.TrimLeft(s, "*`_ "), id.Raw)
		return id, strings.Trim(rest, "*`_ ")
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
func parseDone(md string) ([]Card, []Warning) {
	var cards []Card
	var warnings []Warning
	addWarn := func(msg string) {
		if msg != "" {
			warnings = append(warnings, Warning{Kind: "malformed-done", Message: msg})
		}
	}
	for _, sec := range splitSections(md) {
		if sec.name != "Completed" {
			continue
		}
		// Two real-world shapes: the skeleton's "- <ID> — … — …" bullets, and the
		// "### <ID> — <title>" + prose-body headings some instances use. A section
		// with any "### " line is heading-style (its body may contain "- " list
		// lines, which must not be mistaken for bullet entries).
		headingStyle := false
		for _, ln := range sec.lines {
			if strings.HasPrefix(ln, "### ") {
				headingStyle = true
				break
			}
		}
		if headingStyle {
			for _, block := range splitBlocks(sec.lines, "### ") {
				card, warn := parseDoneHeading(block)
				cards = append(cards, card)
				addWarn(warn)
			}
			continue
		}
		for _, ln := range sec.lines {
			t := strings.TrimSpace(ln)
			if !strings.HasPrefix(t, "- ") {
				continue
			}
			card, warn := parseDoneEntry(t, ln)
			cards = append(cards, card)
			addWarn(warn)
		}
	}
	return cards, warnings
}

// parseDoneHeading parses a "### <ID> — <title>" DONE entry whose body carries a
// "**Completed:** <date> · **PR:** #N · **Journal:** …" line plus a prose summary.
func parseDoneHeading(lines []string) (Card, string) {
	id, title := parseIDTitle(strings.TrimPrefix(lines[0], "### "))
	body := strings.Join(lines[1:], "\n")
	rec := &DoneRecord{
		Date:           fieldAfter(body, "**Completed:**"),
		DeliveryRecord: prPrefixed(fieldAfter(body, "**PR:**")),
		JournalPointer: fieldAfter(body, "**Journal:**"),
		Summary:        doneSummary(body),
	}
	card := Card{ID: id, Title: title, Column: ColDone, Raw: strings.Join(lines, "\n"), Done: rec}
	if id == nil {
		return card, "DONE heading without a parseable id: " + strings.TrimSpace(lines[0])
	}
	return card, ""
}

// fieldAfter returns the text after a "**Marker:**" up to the next " · " field
// separator or end of line (the metadata line packs several fields with " · ").
func fieldAfter(body, marker string) string {
	i := strings.Index(body, marker)
	if i < 0 {
		return ""
	}
	rest := body[i+len(marker):]
	if j := strings.IndexByte(rest, '\n'); j >= 0 {
		rest = rest[:j]
	}
	if j := strings.Index(rest, " · "); j >= 0 {
		rest = rest[:j]
	}
	return strings.TrimSpace(rest)
}

// doneSummary is the prose body with the metadata line(s) removed.
func doneSummary(body string) string {
	var out []string
	for _, ln := range strings.Split(body, "\n") {
		t := strings.TrimSpace(ln)
		if t == "" || strings.HasPrefix(t, "**Completed:**") || strings.HasPrefix(t, "**PR:**") || strings.HasPrefix(t, "**Journal:**") {
			continue
		}
		out = append(out, t)
	}
	return strings.TrimSpace(strings.Join(out, " "))
}

// prPrefixed normalizes a bare "#161 …" PR field to "PR #161 …".
func prPrefixed(s string) string {
	if strings.HasPrefix(s, "#") {
		return "PR " + s
	}
	return s
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
func reconcile(current []Card, backlog []backlogItem, done []Card) ([]Card, []Warning) {
	var warns []Warning
	var cards []Card
	seen := map[string]bool{} // ids already placed in their furthest column

	doneIDs := map[string]bool{}
	for _, d := range done {
		if d.ID != nil {
			doneIDs[d.ID.Raw] = true
		}
	}
	backlogChecked := map[string]bool{}
	currentByID := map[string]Card{}
	for _, b := range backlog {
		if b.Checked && b.ID != nil && !b.Parking { // parking rows aren't "shipped"
			backlogChecked[b.ID.Raw] = true
		}
	}
	for _, c := range current {
		if c.ID != nil {
			currentByID[c.ID.Raw] = c
		}
	}

	if len(current) > 1 {
		warns = append(warns, Warning{Kind: "current-multiple", Message: "CURRENT holds >1 active task (framework one-task invariant)"})
	}

	// Done column: real DONE entries first.
	for _, d := range done {
		if d.ID != nil {
			seen[d.ID.Raw] = true
			if !backlogChecked[d.ID.Raw] {
				warns = append(warns, Warning{Kind: "done-not-ticked", Message: "DONE item not ticked in BACKLOG: " + d.ID.Raw, TaskRaw: d.ID.Raw})
			}
		}
		cards = append(cards, d)
	}
	// A [x] backlog item with no DONE entry is still shown in Done, with a warning
	// (once per id, even if the backlog lists it twice). If that id is also the
	// CURRENT task, promote that card so its detail survives, rather than a stub.
	for _, b := range backlog {
		if b.Parking || !b.Checked || b.ID == nil || doneIDs[b.ID.Raw] || seen[b.ID.Raw] {
			continue
		}
		seen[b.ID.Raw] = true
		warns = append(warns, Warning{Kind: "shipped-missing-done", Message: "shipped item missing from DONE.md: " + b.ID.Raw, TaskRaw: b.ID.Raw})
		if cc, ok := currentByID[b.ID.Raw]; ok {
			cc.Column = ColDone
			cards = append(cards, cc)
		} else {
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
	// Parking-lot items: Backlog column, muted. An id-less parking item is always
	// its own card; one whose id already appears as a tracked card is deduped away
	// so the "one card per id" invariant holds.
	for _, b := range backlog {
		if !b.Parking {
			continue
		}
		if b.ID != nil {
			if seen[b.ID.Raw] {
				continue
			}
			seen[b.ID.Raw] = true
		}
		cards = append(cards, Card{ID: b.ID, Title: b.Title, Column: ColBacklog, ParkingLot: true, Raw: b.Raw})
	}
	return cards, warns
}

// --- blockers (spec §2) -----------------------------------------------------

// parseBlockers reads blockers.md. Blocker headings and section headings are
// both "## ": a heading matching ^BLOCKER-\d+ is a blocker, assigned to the last
// section seen; Format/Open/Resolved are sections; everything under "## Format"
// is skipped. Open = the blocker sits under "## Open".
func parseBlockers(md string) []Blocker {
	var blockers []Blocker
	var cur *Blocker
	section := ""
	flush := func() {
		if cur != nil {
			cur.Body = strings.TrimSpace(cur.Body)
			blockers = append(blockers, *cur)
			cur = nil
		}
	}
	for _, line := range strings.Split(md, "\n") {
		if strings.HasPrefix(line, "## ") {
			head := strings.TrimSpace(line[len("## "):])
			flush()
			if blockerHeadRe.MatchString(head) {
				if section != "Format" { // skip the Format example
					cur = parseBlockerHead(head, section)
				}
			} else {
				section = head
			}
			continue
		}
		if cur == nil {
			continue
		}
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "**Task affected:**") {
			cur.TaskRaw = strings.TrimSpace(strings.TrimPrefix(t, "**Task affected:**"))
		} else {
			cur.Body += line + "\n"
		}
	}
	flush()
	return blockers
}

func parseBlockerHead(head, section string) *Blocker {
	b := &Blocker{Open: section == "Open"}
	if id := parseID(head); id != nil {
		b.ID = *id
	}
	// Anchor on the trailing "opened <ts>" field (like parseDoneEntry anchors on
	// its last field) so a title that itself contains " — " still parses: the
	// title is everything between the id and the opened field.
	parts := strings.Split(head, emDash)
	openedIdx := -1
	for i := len(parts) - 1; i >= 1; i-- {
		if strings.HasPrefix(strings.TrimSpace(parts[i]), "opened ") {
			openedIdx = i
			break
		}
	}
	switch {
	case openedIdx >= 1:
		b.Opened = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(parts[openedIdx]), "opened "))
		if openedIdx >= 2 {
			b.Title = strings.TrimSpace(strings.Join(parts[1:openedIdx], emDash))
		}
	case len(parts) >= 2:
		b.Title = strings.TrimSpace(strings.Join(parts[1:], emDash)) // no "opened" field
	}
	return b
}

// --- badges (spec §6) -------------------------------------------------------

// computeBadges attaches the cross-cutting chips by joining planning + progress.
// blocked is the load-bearing join (an open blocker's affected task); the rest
// are planning-only. Order: blocked, parking, override, no-ac, then namespace.
func computeBadges(cards []Card, blockers []Blocker) {
	blocked := map[string]bool{}
	for _, bl := range blockers {
		// Normalize the affected-task id: upper-case (case-insensitive join, per
		// the D2 route-key decision) then take the leading id token, so trailing
		// text like "DEMO-2 (importer path)" still matches.
		if bl.Open {
			if id := parseID(strings.ToUpper(bl.TaskRaw)); id != nil {
				blocked[id.Raw] = true
			}
		}
	}
	for i := range cards {
		c := &cards[i]
		var badges []Badge
		if c.ID != nil && blocked[strings.ToUpper(c.ID.Raw)] {
			badges = append(badges, Badge{Kind: "blocked"})
		}
		if c.ParkingLot {
			badges = append(badges, Badge{Kind: "parking"})
		}
		if c.DeliveryOverride != "" {
			badges = append(badges, Badge{Kind: "override"})
		}
		if c.Column == ColInProgress && len(c.Acceptance) == 0 {
			badges = append(badges, Badge{Kind: "no-ac"})
		}
		if c.ID != nil {
			badges = append(badges, Badge{Kind: "namespace", Label: c.ID.Namespace})
		}
		c.Badges = badges
	}
}
