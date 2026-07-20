package backlog

import (
	"strings"
	"testing"
)

// AC1: heading-style DONE entries parse into full DoneRecords.
func TestParseDoneHeadingFormat(t *testing.T) {
	md := strings.Join([]string{
		"# Done",
		"## Completed",
		"### TRACK-000 — Swift/iOS toolchain bootstrap + orientation",
		`**Completed:** 2026-06-25 · **PR:** #161 (squash-merged) · **Journal:** 2026-06-25 "TRACK-000 COMPLETE".`,
		"The build owner's first Swift/iOS setup. Installed Xcode 15.3.",
		"Deployment-target pinning deferred to TRACK-001.",
		"### TRACK-001 — WI-1: project skeleton",
		`**Completed:** 2026-06-26 · **PR:** #162 · **Journal:** later.`,
		"Xcode/SwiftUI app skeleton.",
	}, "\n")

	cards, warns := parseDone(md)
	if len(cards) != 2 {
		t.Fatalf("cards=%d, want 2", len(cards))
	}
	if len(warns) != 0 {
		t.Errorf("unexpected warnings: %v", warns)
	}

	c0 := cards[0]
	if c0.ID == nil || c0.ID.Raw != "TRACK-000" {
		t.Fatalf("id=%v, want TRACK-000", c0.ID)
	}
	if c0.Title != "Swift/iOS toolchain bootstrap + orientation" {
		t.Errorf("title=%q", c0.Title)
	}
	if c0.Column != ColDone || c0.Done == nil {
		t.Fatalf("want a Done card with a DoneRecord")
	}
	if c0.Done.Date != "2026-06-25" {
		t.Errorf("date=%q, want 2026-06-25", c0.Done.Date)
	}
	if c0.Done.DeliveryRecord != "PR #161 (squash-merged)" {
		t.Errorf("delivery=%q, want PR #161 (squash-merged)", c0.Done.DeliveryRecord)
	}
	if !strings.Contains(c0.Done.JournalPointer, "TRACK-000 COMPLETE") {
		t.Errorf("journal=%q", c0.Done.JournalPointer)
	}
	if !strings.Contains(c0.Done.Summary, "first Swift/iOS setup") {
		t.Errorf("summary should contain the prose, got %q", c0.Done.Summary)
	}
	if strings.Contains(c0.Done.Summary, "**Completed:**") {
		t.Errorf("summary must exclude the metadata line, got %q", c0.Done.Summary)
	}
}

// AC3: ids wrapped in markdown decoration are recognized.
func TestParseIDEmphasis(t *testing.T) {
	cases := map[string]string{
		"**WI-8 — .trace export**": "WI-8",
		"`TRACK-1` — x":            "TRACK-1",
		"TRACK-000 — normal":       "TRACK-000",
		"*GW-12*":                  "GW-12",
	}
	for in, want := range cases {
		if id := parseID(in); id == nil || id.Raw != want {
			t.Errorf("parseID(%q) = %v, want %s", in, id, want)
		}
	}
	// A checkbox line with a bold id yields an id-bearing backlog item.
	items := parseBacklog("## Active\n- [ ] **WI-8 — .trace export**\n")
	if len(items) != 1 || items[0].ID == nil || items[0].ID.Raw != "WI-8" {
		t.Errorf("bold-id backlog item: %+v", items)
	}
}

// AC2: a "### " line in the body doesn't create phantom bullet entries, and the
// bullet format still works when there are no headings.
func TestParseDoneFormatExclusivity(t *testing.T) {
	heading := "## Completed\n### A-1 — t\nbody with a - dash list line\n- not an entry\n"
	cards, _ := parseDone(heading)
	if len(cards) != 1 || cards[0].ID.Raw != "A-1" {
		t.Errorf("heading-style: want 1 card A-1, got %d %+v", len(cards), cards)
	}
	bullet := "## Completed\n- B-2 — t — 2026-01-01 — s — PR #1. See journal x\n"
	cards2, _ := parseDone(bullet)
	if len(cards2) != 1 || cards2[0].ID.Raw != "B-2" {
		t.Errorf("bullet-style: want 1 card B-2, got %d %+v", len(cards2), cards2)
	}
}
