package main

import "github.com/charmbracelet/lipgloss"

const (
	primaryDefault = "#156bb3"
	primaryLight   = "#3d8fd1"
	primaryDark    = "#0e4f87"

	accentDefault = "#ea944c"
	accentLight   = "#f0b07a"
	accentDark    = "#c8712a"

	neutral0   = "#ffffff"
	neutral100 = "#e8e8e8"
	neutral200 = "#d0d0d0"
	neutral400 = "#9a9a9a"
	neutral600 = "#5a5a5a"
	neutral700 = "#3f3f3f"
	neutral800 = "#2e2e2e"
	neutral900 = "#191919"
)

type Styles struct {
	sidePaneStyle        lipgloss.Style
	mainPaneActiveStyle  lipgloss.Style
	activeItemStyle      lipgloss.Style
	inactiveItemStyle    lipgloss.Style
	nameStyle            lipgloss.Style
	taglineStyle         lipgloss.Style
	subtitleStyle        lipgloss.Style
	sectionHeaderStyle   lipgloss.Style
	contentStyle         lipgloss.Style
	mutedStyle           lipgloss.Style
	dimStyle             lipgloss.Style
	bulletStyle          lipgloss.Style
	projectTitleStyle    lipgloss.Style
	statusBarStyle       lipgloss.Style
	statusKeyStyle       lipgloss.Style
	statusSepStyle       lipgloss.Style
	statusScrollBarStyle lipgloss.Style
	statusScrollDimStyle lipgloss.Style
	statusScrollPctStyle lipgloss.Style
}

func makeStyles(r *lipgloss.Renderer) Styles {
	return Styles{
		sidePaneStyle: r.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(neutral600)).
			Foreground(lipgloss.Color(neutral100)).
			Padding(1, 2),

		mainPaneActiveStyle: r.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(primaryDefault)).
			Foreground(lipgloss.Color(neutral100)).
			Padding(1, 3),

		activeItemStyle: r.NewStyle().
			Foreground(lipgloss.Color(accentDefault)).
			Bold(true),

		inactiveItemStyle: r.NewStyle().
			Foreground(lipgloss.Color(neutral600)),

		nameStyle: r.NewStyle().
			Foreground(lipgloss.Color(neutral100)).
			Bold(true),

		taglineStyle: r.NewStyle().
			Foreground(lipgloss.Color(neutral400)),

		subtitleStyle: r.NewStyle().
			Foreground(lipgloss.Color(primaryLight)).
			Bold(true),

		sectionHeaderStyle: r.NewStyle().
			Foreground(lipgloss.Color(primaryLight)).
			Bold(true),

		contentStyle: r.NewStyle().
			Foreground(lipgloss.Color(neutral100)),

		mutedStyle: r.NewStyle().
			Foreground(lipgloss.Color(neutral600)),

		dimStyle: r.NewStyle().
			Foreground(lipgloss.Color(neutral700)),

		bulletStyle: r.NewStyle().
			Foreground(lipgloss.Color(accentDefault)),

		projectTitleStyle: r.NewStyle().
			Foreground(lipgloss.Color(primaryLight)).
			Bold(true),

		statusBarStyle: r.NewStyle().
			Background(lipgloss.Color(neutral800)).
			Foreground(lipgloss.Color(neutral400)),

		statusKeyStyle: r.NewStyle().
			Background(lipgloss.Color(neutral700)).
			Foreground(lipgloss.Color(neutral100)).
			Padding(0, 1).
			Bold(true),

		statusSepStyle: r.NewStyle().
			Background(lipgloss.Color(neutral800)).
			Foreground(lipgloss.Color(neutral700)),

		statusScrollBarStyle: r.NewStyle().
			Background(lipgloss.Color(neutral800)).
			Foreground(lipgloss.Color(primaryLight)),

		statusScrollDimStyle: r.NewStyle().
			Background(lipgloss.Color(neutral800)).
			Foreground(lipgloss.Color(neutral700)),

		statusScrollPctStyle: r.NewStyle().
			Background(lipgloss.Color(neutral800)).
			Foreground(lipgloss.Color(neutral400)),
	}
}
