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
