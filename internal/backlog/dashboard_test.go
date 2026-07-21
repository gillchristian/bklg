package backlog

import (
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

func hasWarnKind(ws []Warning, kind string) bool {
	for _, w := range ws {
		if w.Kind == kind {
			return true
		}
	}
	return false
}
