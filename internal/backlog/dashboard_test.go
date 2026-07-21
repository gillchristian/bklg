package backlog

import (
	"net/http/httptest"
	"strings"
	"testing"
)

// TestParseDashboard exercises the dashboard adapter (TASK-014, ADR-0004) on the
// testdata/dashboard/work.md fixture: pipe-table Active/Done, bullet-group
// Backlog, escaped pipes, a leading-⛔ blocked row, multi/zero-ticket rows, and a
// malformed (empty-title) row.
func TestParseDashboard(t *testing.T) {
	b, err := NewParser().Parse(Areas{KnowledgeDir: "testdata/dashboard", DashboardFile: "testdata/dashboard/work.md"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	byCol := map[Column]int{}
	for _, c := range b.Cards {
		if !c.Dashboard {
			t.Errorf("card %q not flagged Dashboard", c.Title)
		}
		byCol[c.Column]++
	}
	if byCol[ColInProgress] != 3 || byCol[ColBacklog] != 4 || byCol[ColDone] != 2 {
		t.Errorf("column counts = InProgress %d / Backlog %d / Done %d, want 3 / 4 / 2",
			byCol[ColInProgress], byCol[ColBacklog], byCol[ColDone])
	}

	find := func(title string) *Card {
		for i := range b.Cards {
			if b.Cards[i].Title == title {
				return &b.Cards[i]
			}
		}
		t.Fatalf("no card titled %q", title)
		return nil
	}
	tickets := func(c *Card) []string {
		var r []string
		for _, id := range c.Tickets {
			r = append(r, id.Raw)
		}
		return r
	}
	eq := func(got, want []string) bool {
		if len(got) != len(want) {
			return false
		}
		for i := range got {
			if got[i] != want[i] {
				return false
			}
		}
		return true
	}

	// Alpha: leading-⛔ status ⇒ blocked; two tickets from the status prose.
	alpha := find("Alpha task")
	if !alpha.Blocked {
		t.Errorf("Alpha: Blocked = false, want true (leading ⛔ in status)")
	}
	if alpha.Column != ColInProgress {
		t.Errorf("Alpha: Column = %v, want In Progress", alpha.Column)
	}
	if alpha.Subtitle != "first active thing" {
		t.Errorf("Alpha: Subtitle = %q", alpha.Subtitle)
	}
	if alpha.Material != "[alpha/](./alpha/index.md)" {
		t.Errorf("Alpha: Material = %q", alpha.Material)
	}
	if !eq(tickets(alpha), []string{"PINATA-100", "PINATA-101"}) {
		t.Errorf("Alpha: tickets = %v, want [PINATA-100 PINATA-101]", tickets(alpha))
	}

	// Beta: ticket in the title; the #11531 PR ref must NOT become a ticket.
	beta := find("Beta task (PINATA-200)")
	if !eq(tickets(beta), []string{"PINATA-200"}) {
		t.Errorf("Beta: tickets = %v, want [PINATA-200] (PR #11531 must be ignored)", tickets(beta))
	}
	if beta.Blocked {
		t.Errorf("Beta: Blocked = true, want false")
	}

	// Gamma: escaped pipes stay literal in the status; zero tickets.
	gamma := find("Gamma task")
	if !strings.Contains(gamma.Status, "cond A || cond B || cond C") {
		t.Errorf("Gamma: Status = %q, want the escaped pipes preserved as ||", gamma.Status)
	}
	if len(gamma.Tickets) != 0 {
		t.Errorf("Gamma: tickets = %v, want none", tickets(gamma))
	}

	// Backlog groups.
	if g := find("Ungrouped item").Group; g != "" {
		t.Errorf("Ungrouped item: Group = %q, want empty", g)
	}
	delta := find("Delta feature (PINATA-300, PINATA-301)")
	if delta.Group != "Product / code" {
		t.Errorf("Delta: Group = %q, want 'Product / code'", delta.Group)
	}
	if !eq(tickets(delta), []string{"PINATA-300", "PINATA-301"}) {
		t.Errorf("Delta: tickets = %v, want two", tickets(delta))
	}
	if find("Epsilon").Group != "Product / code" {
		t.Errorf("Epsilon should inherit the 'Product / code' group")
	}
	zeta := find("Zeta doc")
	if zeta.Group != "Knowledge base" || !zeta.Blocked {
		t.Errorf("Zeta: Group = %q Blocked = %v, want 'Knowledge base' + blocked", zeta.Group, zeta.Blocked)
	}

	// Done: title from the What column, material from Record, ticket parsed; the
	// PR-link row yields no ticket.
	eta := find("Eta shipped (PINATA-400)")
	if eta.Column != ColDone || eta.Material != "[eta.md](./eta.md)" || !eq(tickets(eta), []string{"PINATA-400"}) {
		t.Errorf("Eta: col=%v material=%q tickets=%v", eta.Column, eta.Material, tickets(eta))
	}
	if len(find("Theta shipped").Tickets) != 0 {
		t.Errorf("Theta: want no tickets (a #99 PR link is not a ticket)")
	}

	// The empty-title row is skipped with a warning, but the rest of the board
	// still parsed (9 good cards above).
	if !hasWarnKind(b.Warnings, "dashboard-malformed") {
		t.Errorf("want a dashboard-malformed warning for the empty-title row; warnings = %v", b.Warnings)
	}
}

func TestSplitDashTitle(t *testing.T) {
	cases := []struct{ in, title, subtitle string }{
		{"**Alpha task** — first active thing", "Alpha task", "first active thing"},
		{"**Just bold**", "Just bold", ""},
		{"Plain — with subtitle", "Plain", "with subtitle"},
		{"Plain only", "Plain only", ""},
		// Em-dash INSIDE the bold phrase must not split the subtitle early.
		{"**Foo — bar** — realsubtitle", "Foo — bar", "realsubtitle"},
	}
	for _, c := range cases {
		title, subtitle := splitDashTitle(c.in)
		if title != c.title || subtitle != c.subtitle {
			t.Errorf("splitDashTitle(%q) = (%q, %q), want (%q, %q)", c.in, title, subtitle, c.title, c.subtitle)
		}
	}
}

func TestTicketURL(t *testing.T) {
	cases := []struct{ base, id, want string }{
		{"https://linear.app/gopinata/issue/", "PINATA-1", "https://linear.app/gopinata/issue/PINATA-1"},
		{"https://linear.app/gopinata/issue", "PINATA-1", "https://linear.app/gopinata/issue/PINATA-1"}, // no trailing slash
		{"", "PINATA-1", "PINATA-1"}, // no base configured
	}
	for _, c := range cases {
		if got := ticketURL(c.base, c.id); got != c.want {
			t.Errorf("ticketURL(%q,%q) = %q, want %q", c.base, c.id, got, c.want)
		}
	}
}

// TestDashboardBoardRender renders the board in dashboard mode (TASK-015) and
// checks the blocked badge, group chip, ticket-chip links, and the absence of a
// spurious no-ac badge / internal id link.
func TestDashboardBoardRender(t *testing.T) {
	srv := NewServer(Areas{
		KnowledgeDir:  "testdata/dashboard",
		DashboardFile: "testdata/dashboard/work.md",
		LinkBase:      "https://linear.app/acme/issue/",
	})
	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != 200 {
		t.Fatalf("board status %d", rec.Code)
	}
	body := rec.Body.String()

	if !strings.Contains(body, ">blocked<") {
		t.Error("board missing a blocked badge (Alpha/Zeta are blocked)")
	}
	if !strings.Contains(body, "Product / code") {
		t.Error("board missing the 'Product / code' group chip")
	}
	if !strings.Contains(body, `href="https://linear.app/acme/issue/PINATA-100"`) {
		t.Error("board missing a ticket chip linking to the configured Linear base")
	}
	// Dashboard In-Progress cards are AC-less + id-less: no no-ac badge, and no
	// internal /<id> detail anchor (only external ticket links).
	if strings.Contains(body, ">no-ac<") {
		t.Error("dashboard cards must not carry a no-ac badge")
	}
	if strings.Contains(body, `href="/PINATA-100"`) {
		t.Error("dashboard cards must not emit an internal /<id> link")
	}
	// Instead, the card title links to its slug detail page (TASK-016).
	if !strings.Contains(body, `href="/alpha-task"`) {
		t.Error("board should link a dashboard card title to its /<slug> detail page")
	}
}

// TestDashboardBadges asserts the positive AND negative of the badge rules:
// blocked cards get the badge, non-blocked ones don't, and dashboard cards
// never get the framework-only no-ac badge.
func TestDashboardBadges(t *testing.T) {
	b, err := NewParser().Parse(Areas{KnowledgeDir: "testdata/dashboard", DashboardFile: "testdata/dashboard/work.md"})
	if err != nil {
		t.Fatal(err)
	}
	kinds := func(title string) []string {
		for _, c := range b.Cards {
			if c.Title == title {
				var k []string
				for _, bd := range c.Badges {
					k = append(k, bd.Kind)
				}
				return k
			}
		}
		t.Fatalf("no card %q", title)
		return nil
	}
	has := func(ks []string, want string) bool {
		for _, k := range ks {
			if k == want {
				return true
			}
		}
		return false
	}
	if !has(kinds("Alpha task"), "blocked") {
		t.Error("Alpha (leading ⛔) should have a blocked badge")
	}
	if has(kinds("Beta task (PINATA-200)"), "blocked") {
		t.Error("Beta (not blocked) must not have a blocked badge")
	}
	if has(kinds("Alpha task"), "no-ac") {
		t.Error("dashboard cards must never get the framework no-ac badge")
	}
}

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"Alpha task":              "alpha-task",
		"Form Building UI (C4)!":  "form-building-ui-c4",
		"  spaced  out  ":         "spaced-out",
		"PINATA-599 & PINATA-601": "pinata-599-pinata-601",
		"!!!":                     "",
	}
	for in, want := range cases {
		if got := slugify(in); got != want {
			t.Errorf("slugify(%q) = %q, want %q", in, got, want)
		}
	}
	// Uniqueness: duplicate titles get distinct slugs; empty falls back to card.
	cards := []Card{{Title: "Same Title"}, {Title: "Same Title"}, {Title: "!!!"}}
	assignSlugs(cards)
	if cards[0].Slug == cards[1].Slug {
		t.Errorf("duplicate titles got the same slug %q", cards[0].Slug)
	}
	if cards[2].Slug == "" {
		t.Error("empty-slug title should fall back to a non-empty slug")
	}
}

// TestDashboardDetail renders a dashboard card's detail page via its slug route
// (TASK-016): title, ticket links (to the configured base), material, status;
// and an unknown slug 404s.
func TestDashboardDetail(t *testing.T) {
	srv := NewServer(Areas{
		KnowledgeDir:  "testdata/dashboard",
		DashboardFile: "testdata/dashboard/work.md",
		LinkBase:      "https://linear.app/acme/issue/",
	})
	get := func(path string) (int, string) {
		rec := httptest.NewRecorder()
		srv.Routes().ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		return rec.Code, rec.Body.String()
	}

	code, body := get("/alpha-task") // slugify("Alpha task")
	if code != 200 {
		t.Fatalf("/alpha-task status %d, want 200", code)
	}
	for _, want := range []string{
		"Alpha task", // title
		`href="https://linear.app/acme/issue/PINATA-100"`, // ticket chip on detail (taskVM.LinearBase)
		"Status / next step",                              // status section label
		"Blocked on PINATA-100",                           // status prose
	} {
		if !strings.Contains(body, want) {
			t.Errorf("/alpha-task detail missing %q", want)
		}
	}

	if code, _ := get("/no-such-slug"); code != 404 {
		t.Errorf("/no-such-slug status %d, want 404", code)
	}
}

func hasWarnKind(ws []Warning, kind string) bool {
	for _, w := range ws {
		if w.Kind == kind {
			return true
		}
	}
	return false
}
