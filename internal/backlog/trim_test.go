package backlog

import (
	"bytes"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestTruncate(t *testing.T) {
	// AC3: short strings pass through unchanged.
	if got := truncate("short title", 140); got != "short title" {
		t.Errorf("short: got %q", got)
	}
	if got := truncate("  trimmed  ", 140); got != "trimmed" {
		t.Errorf("whitespace: got %q", got)
	}

	// AC1: long strings are cut and get an ellipsis.
	long := strings.Repeat("word ", 60) // 300 chars
	got := truncate(long, 140)
	if !strings.HasSuffix(got, "…") {
		t.Errorf("long title should end with …, got %q", got)
	}
	if utf8.RuneCountInString(got) > 141 {
		t.Errorf("truncated length %d > 141", utf8.RuneCountInString(got))
	}
	if got == long {
		t.Error("long title was not truncated")
	}

	// AC4: rune-safe — never splits a multibyte rune.
	multibyte := strings.Repeat("é", 200)
	m := truncate(multibyte, 140)
	if !utf8.ValidString(m) {
		t.Errorf("truncation produced invalid UTF-8: %q", m)
	}
	if !strings.HasSuffix(m, "…") {
		t.Error("multibyte truncation should end with …")
	}
}

// AC1 + AC2: the board trims a long title; the detail page keeps it full.
func TestCardTitleBoardVsDetail(t *testing.T) {
	long := "Swift/iOS toolchain bootstrap and orientation — " + strings.Repeat("a long clause that keeps going and going ", 8)
	card := Card{ID: &ID{Namespace: "X", Number: 1, Raw: "X-1"}, Title: long, Column: ColBacklog}

	var board bytes.Buffer
	bvm := boardVM{PlanningDir: "x", Version: "0", Columns: []columnVM{{Title: "Backlog", Cards: []Card{card}}}}
	if err := boardTmpl.ExecuteTemplate(&board, "layout", bvm); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(board.String(), long) {
		t.Error("board shows the FULL long title; want it truncated")
	}
	if !strings.Contains(board.String(), "…") {
		t.Error("board card should show an ellipsis for the long title")
	}

	var detail bytes.Buffer
	tvm := taskVM{PlanningDir: "x", Version: "0", Card: card}
	if err := taskTmpl.ExecuteTemplate(&detail, "layout", tvm); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(detail.String(), long) {
		t.Error("detail page should show the FULL title, untrimmed")
	}
}
