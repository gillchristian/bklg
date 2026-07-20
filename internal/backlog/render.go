package backlog

import (
	"embed"
	"html/template"
)

//go:embed templates
var templatesFS embed.FS

// funcs are the template helpers. Styling lives here (near rendering), not on
// the model.
var funcs = template.FuncMap{
	"badgeClass": badgeClass,
	"badgeText":  badgeText,
}

// One template set per page (each defines its own "content"), sharing the
// layout — separate sets avoid a "content" redefinition collision.
var (
	boardTmpl = template.Must(
		template.New("board").Funcs(funcs).ParseFS(templatesFS, "templates/layout.html", "templates/board.html"),
	)
	taskTmpl = template.Must(
		template.New("task").Funcs(funcs).ParseFS(templatesFS, "templates/layout.html", "templates/task.html"),
	)
)

// boardVM is the board's view model: the three columns in flow order, plus the
// header/banner inputs.
type boardVM struct {
	PlanningDir string
	Warnings    []string
	Columns     []columnVM
}

type columnVM struct {
	Title string
	Cards []Card
}

// taskVM is the detail page's view model. It carries the shared header inputs
// (PlanningDir/Warnings, so the layout banner renders) plus the card and the
// blockers referencing it (open first, then resolved).
type taskVM struct {
	PlanningDir string
	Warnings    []string
	Card        Card
	Blockers    []Blocker
}

func viewModel(b Board) boardVM {
	cols := []columnVM{{Title: "Backlog"}, {Title: "In Progress"}, {Title: "Done"}}
	for _, c := range b.Cards {
		switch c.Column {
		case ColBacklog:
			cols[0].Cards = append(cols[0].Cards, c)
		case ColInProgress:
			cols[1].Cards = append(cols[1].Cards, c)
		case ColDone:
			cols[2].Cards = append(cols[2].Cards, c)
		}
	}
	return boardVM{
		PlanningDir: DisplayPath(b.Meta.PlanningDir),
		Warnings:    b.Warnings,
		Columns:     cols,
	}
}

// badgeClass maps a badge kind to Tailwind chip classes (spec §6 suggested
// styles). blocked is the load-bearing red.
func badgeClass(kind string) string {
	switch kind {
	case "blocked":
		return "bg-red-100 text-red-800 ring-red-600/20"
	case "parking":
		return "bg-slate-100 text-slate-600 ring-slate-500/20"
	case "override":
		return "bg-amber-100 text-amber-800 ring-amber-600/20"
	case "no-ac":
		return "bg-yellow-100 text-yellow-800 ring-yellow-600/20"
	default: // namespace + anything else
		return "bg-gray-100 text-gray-700 ring-gray-500/20"
	}
}

// badgeText is the chip label: the namespace value for the namespace chip, the
// kind otherwise.
func badgeText(b Badge) string {
	if b.Kind == "namespace" {
		return b.Label
	}
	return b.Kind
}
