package backlog

import (
	"os"
	"strings"
	"testing"
)

func TestParseBlockers(t *testing.T) {
	data, err := os.ReadFile("testdata/knowledge/progress/blockers.md")
	if err != nil {
		t.Fatal(err)
	}
	bl := parseBlockers(string(data))
	if len(bl) != 2 {
		t.Fatalf("got %d blockers, want 2 (the ## Format example must be skipped): %+v", len(bl), bl)
	}
	byRaw := map[string]Blocker{}
	for _, b := range bl {
		byRaw[b.ID.Raw] = b
	}

	b1, ok := byRaw["BLOCKER-001"]
	if !ok || !b1.Open || b1.TaskRaw != "DEMO-2" {
		t.Errorf("BLOCKER-001 = %+v, want Open affecting DEMO-2", b1)
	}
	if !strings.Contains(b1.Title, "Importer flakiness") {
		t.Errorf("BLOCKER-001 title = %q", b1.Title)
	}
	if b1.Opened != "2026-07-19 15:00" {
		t.Errorf("BLOCKER-001 opened = %q, want 2026-07-19 15:00", b1.Opened)
	}
	if !strings.Contains(b1.Body, "429s") {
		t.Errorf("BLOCKER-001 body not captured: %q", b1.Body)
	}

	b2, ok := byRaw["BLOCKER-002"]
	if !ok || b2.Open || b2.TaskRaw != "DEMO-1" {
		t.Errorf("BLOCKER-002 = %+v, want Resolved affecting DEMO-1", b2)
	}
	if _, bad := byRaw["BLOCKER-000"]; bad {
		t.Errorf("BLOCKER-000 (under ## Format) must be skipped")
	}
}

func TestParseBlockerHeadEmDashTitle(t *testing.T) {
	b := parseBlockerHead("BLOCKER-020 — Fix the thing — and another clause — opened 2026-07-20 12:34", "Open")
	if b.Title != "Fix the thing — and another clause" {
		t.Errorf("title = %q, want the full em-dashed title", b.Title)
	}
	if b.Opened != "2026-07-20 12:34" {
		t.Errorf("opened = %q, want 2026-07-20 12:34", b.Opened)
	}
	if !b.Open {
		t.Error("want Open")
	}
	// No "opened" field: whole tail is the title.
	if b2 := parseBlockerHead("BLOCKER-021 — just a title", "Resolved"); b2.Title != "just a title" || b2.Opened != "" {
		t.Errorf("no-opened head: title=%q opened=%q", b2.Title, b2.Opened)
	}
}

func TestBadgeJoinTolerant(t *testing.T) {
	cards := []Card{{ID: &ID{Namespace: "DEMO", Number: 2, Raw: "DEMO-2"}, Column: ColInProgress}}
	// lower-case + trailing text after the id must still join.
	computeBadges(cards, []Blocker{{TaskRaw: "demo-2 (importer path)", Open: true}})
	if !hasBadge(&cards[0], "blocked") {
		t.Errorf("case-insensitive + trailing-text join should mark DEMO-2 blocked; badges=%+v", cards[0].Badges)
	}
}

func hasBadge(c *Card, kind string) bool {
	for _, bg := range c.Badges {
		if bg.Kind == kind {
			return true
		}
	}
	return false
}

func namespaceBadge(c *Card) string {
	for _, bg := range c.Badges {
		if bg.Kind == "namespace" {
			return bg.Label
		}
	}
	return ""
}

func TestBadges(t *testing.T) {
	b := board(t)

	d2 := byID(b, "DEMO-2")
	if d2 == nil {
		t.Fatal("DEMO-2 missing")
	}
	if !hasBadge(d2, "blocked") {
		t.Error("DEMO-2 should be blocked (open BLOCKER-001 affects it)")
	}
	if !hasBadge(d2, "no-ac") {
		t.Error("DEMO-2 should carry no-ac (In-Progress with zero criteria)")
	}
	if namespaceBadge(d2) != "DEMO" {
		t.Errorf("DEMO-2 namespace badge = %q, want DEMO", namespaceBadge(d2))
	}

	d1 := byID(b, "DEMO-1")
	if d1 == nil {
		t.Fatal("DEMO-1 missing")
	}
	if hasBadge(d1, "blocked") {
		t.Error("DEMO-1 must NOT be blocked (BLOCKER-002 affecting it is resolved)")
	}
	if !hasBadge(d1, "override") {
		t.Error("DEMO-1 should carry override (has a Delivery override)")
	}
	if hasBadge(d1, "no-ac") {
		t.Error("DEMO-1 has 2 criteria; must not carry no-ac")
	}

	var parking *Card
	for i := range b.Cards {
		if b.Cards[i].ID == nil {
			parking = &b.Cards[i]
		}
	}
	if parking == nil || !hasBadge(parking, "parking") {
		t.Errorf("parking card should carry the parking badge: %+v", parking)
	}
}
