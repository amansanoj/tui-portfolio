package main

import "github.com/charmbracelet/lipgloss"

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
			Background(lipgloss.Color(neutral900)).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(neutral600)).
			Foreground(lipgloss.Color(neutral100)).
			Padding(1, 2),

		mainPaneActiveStyle: r.NewStyle().
			Background(lipgloss.Color(neutral900)).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(primaryDefault)).
			Foreground(lipgloss.Color(neutral100)).
			Padding(1, 2),

		activeItemStyle: r.NewStyle().
			Foreground(lipgloss.Color(secondaryDefault)).
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
			Foreground(lipgloss.Color(secondaryDefault)),

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
