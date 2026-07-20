package backlog

import (
	"strings"
	"testing"
)

// board builds the demo instance's board once for the assertions below.
func board(t *testing.T) Board {
	t.Helper()
	// The fixture instance is testdata/knowledge, so base = "testdata"/"knowledge".
	a, err := Resolve("testdata", "knowledge")
	if err != nil {
		t.Fatalf("resolve fixture: %v", err)
	}
	b, err := NewParser().Parse(a)
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	return b
}

func byID(b Board, raw string) *Card {
	for i := range b.Cards {
		if b.Cards[i].ID != nil && b.Cards[i].ID.Raw == raw {
			return &b.Cards[i]
		}
	}
	return nil
}

// AC1: card set + per-card {id, column, len(Acceptance), ParkingLot}.
func TestParseCardTable(t *testing.T) {
	b := board(t)

	type row struct {
		col     Column
		acLen   int
		parking bool
	}
	want := map[string]row{
		"DEMO-1": {ColInProgress, 2, false},
		"DEMO-2": {ColInProgress, 0, false},
		"DEMO-3": {ColBacklog, 0, false},
		"DEMO-4": {ColDone, 0, false},
		"DEMO-5": {ColDone, 0, false},
		"DEMO-6": {ColDone, 0, false},
	}

	idCards, parkingCards := 0, 0
	for _, c := range b.Cards {
		if c.ID == nil {
			parkingCards++
			if !c.ParkingLot {
				t.Errorf("id-less card %q: ParkingLot=false, want true", c.Title)
			}
			continue
		}
		idCards++
		w, ok := want[c.ID.Raw]
		if !ok {
			t.Errorf("unexpected card %s", c.ID.Raw)
			continue
		}
		if c.Column != w.col {
			t.Errorf("%s: column=%v, want %v", c.ID.Raw, c.Column, w.col)
		}
		if len(c.Acceptance) != w.acLen {
			t.Errorf("%s: len(Acceptance)=%d, want %d", c.ID.Raw, len(c.Acceptance), w.acLen)
		}
		if c.ParkingLot != w.parking {
			t.Errorf("%s: ParkingLot=%v, want %v", c.ID.Raw, c.ParkingLot, w.parking)
		}
	}
	if idCards != len(want) {
		t.Errorf("id cards=%d, want %d", idCards, len(want))
	}
	if parkingCards != 1 {
		t.Errorf("parking cards=%d, want 1", parkingCards)
	}
	if len(b.Cards) != 7 {
		t.Errorf("total cards=%d, want 7", len(b.Cards))
	}
}

// AC2: dedup — most-advanced-state wins, each id appears exactly once.
func TestParseDedup(t *testing.T) {
	b := board(t)
	count := map[string]int{}
	for _, c := range b.Cards {
		if c.ID != nil {
			count[c.ID.Raw]++
		}
	}
	for id, n := range count {
		if n != 1 {
			t.Errorf("%s appears %d times, want 1 (dedup)", id, n)
		}
	}
	// DEMO-1 is in CURRENT and unchecked in BACKLOG -> InProgress wins.
	if c := byID(b, "DEMO-1"); c == nil || c.Column != ColInProgress {
		t.Errorf("DEMO-1 should be a single In-Progress card, got %+v", c)
	}
	// DEMO-4 is [x] in BACKLOG and present in DONE -> Done wins.
	if c := byID(b, "DEMO-4"); c == nil || c.Column != ColDone {
		t.Errorf("DEMO-4 should be a single Done card, got %+v", c)
	}
}

// AC3: exactly the three reconciliation warnings on the seeded inconsistencies.
func TestParseWarnings(t *testing.T) {
	b := board(t)
	wantSub := []string{
		"CURRENT holds >1 active task",
		"DONE item not ticked in BACKLOG: DEMO-6",
		"shipped item missing from DONE.md: DEMO-5",
	}
	for _, sub := range wantSub {
		found := false
		for _, w := range b.Warnings {
			if strings.Contains(w, sub) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing warning containing %q; got %v", sub, b.Warnings)
		}
	}
	if len(b.Warnings) != 3 {
		t.Errorf("warnings=%d (%v), want exactly 3", len(b.Warnings), b.Warnings)
	}
}

// AC4: field extraction.
func TestParseFields(t *testing.T) {
	b := board(t)

	d1 := byID(b, "DEMO-1")
	if d1 == nil {
		t.Fatal("DEMO-1 missing")
	}
	if len(d1.Acceptance) != 2 || !d1.Acceptance[0].Checked || d1.Acceptance[1].Checked {
		t.Errorf("DEMO-1 acceptance = %+v, want [checked, unchecked]", d1.Acceptance)
	}
	if !strings.Contains(d1.Source, "BACKLOG") {
		t.Errorf("DEMO-1 Source = %q", d1.Source)
	}
	if !strings.Contains(d1.Notes, "Depends on the exporter") {
		t.Errorf("DEMO-1 Notes = %q", d1.Notes)
	}
	if !strings.Contains(d1.DeliveryOverride, "may commit directly") {
		t.Errorf("DEMO-1 DeliveryOverride = %q", d1.DeliveryOverride)
	}
	if d1.ID.Namespace != "DEMO" || d1.ID.Number != 1 {
		t.Errorf("DEMO-1 ID = %+v, want {DEMO 1}", *d1.ID)
	}

	d4 := byID(b, "DEMO-4")
	if d4 == nil || d4.Done == nil {
		t.Fatalf("DEMO-4 or its DoneRecord missing: %+v", d4)
	}
	if d4.Done.Date != "2026-07-18" {
		t.Errorf("DEMO-4 date = %q, want 2026-07-18", d4.Done.Date)
	}
	// Summary contains an embedded " — "; it must survive round-trip.
	if d4.Done.Summary != "Streaming exporter with a header row — split into read and write stages." {
		t.Errorf("DEMO-4 summary = %q", d4.Done.Summary)
	}
	if d4.Done.DeliveryRecord != "PR #12, merged `abc1234`" {
		t.Errorf("DEMO-4 delivery = %q", d4.Done.DeliveryRecord)
	}
	if !strings.Contains(d4.Done.JournalPointer, "See journal 2026-07-18 10:00") {
		t.Errorf("DEMO-4 journal ptr = %q", d4.Done.JournalPointer)
	}

	// DEMO-5 is Done (shipped) but has no DONE entry -> no DoneRecord.
	if d5 := byID(b, "DEMO-5"); d5 == nil || d5.Column != ColDone || d5.Done != nil {
		t.Errorf("DEMO-5 want Done column with nil DoneRecord, got %+v", d5)
	}

	// Parking card: no id, ParkingLot, title captured.
	var parking *Card
	for i := range b.Cards {
		if b.Cards[i].ID == nil {
			parking = &b.Cards[i]
		}
	}
	if parking == nil || !parking.ParkingLot || !strings.Contains(parking.Title, "colour palette") {
		t.Errorf("parking card = %+v", parking)
	}
}

// AC5: parse defensively — malformed input never panics and good entries survive.
func TestParseDefensive(t *testing.T) {
	md := strings.Join([]string{
		"## Completed",
		"- totally broken line with no fields",
		"- DEMO-9 — a good one — 2026-07-01 — fine — PR #1, merged `aaa`. See journal 2026-07-01",
	}, "\n")
	cards, warns := parseDone(md) // must not panic
	if len(cards) != 2 {
		t.Errorf("cards=%d, want 2 (best-effort keeps both)", len(cards))
	}
	if len(warns) == 0 {
		t.Errorf("expected a warning for the malformed line")
	}
	// The good entry still parses fully.
	var good *Card
	for i := range cards {
		if cards[i].ID != nil && cards[i].ID.Raw == "DEMO-9" {
			good = &cards[i]
		}
	}
	if good == nil || good.Done == nil || good.Done.Date != "2026-07-01" {
		t.Errorf("good entry not parsed: %+v", good)
	}

	// A whole board over a missing planning area must warn, not crash.
	b, err := NewParser().Parse(Areas{PlanningDir: "testdata/does-not-exist", ProgressDir: "testdata/nope"})
	if err != nil {
		t.Errorf("Parse should degrade, not error: %v", err)
	}
	if len(b.Warnings) == 0 {
		t.Errorf("expected read warnings for a missing planning area")
	}
}
