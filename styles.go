package main

import "github.com/charmbracelet/lipgloss"

const (
	primaryDefault = "#3d8fd1"
	primaryLight   = "#5a9fd8"
	primaryDeep    = "#276ca5"

	secondaryDefault = "#f0b07a"
	secondaryStrong  = "#e67519"

	accentDefault = "#a89e69"
	accentLight   = "#b5ab7d"

	neutral100 = "#e6e6e6"
	neutral400 = "#999999"
	neutral600 = "#666666"
	neutral700 = "#4d4d4d"
	neutral800 = "#333333"
	neutral900 = "#1a1a1a"
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
	accentStyle          lipgloss.Style
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
			Foreground(lipgloss.Color(secondaryStrong)).
			Bold(true),

		inactiveItemStyle: r.NewStyle().
			Foreground(lipgloss.Color(secondaryDefault)),

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

		accentStyle: r.NewStyle().
			Foreground(lipgloss.Color(accentDefault)).
			Bold(true),

		bulletStyle: r.NewStyle().
			Foreground(lipgloss.Color(accentDefault)),

		projectTitleStyle: r.NewStyle().
			Foreground(lipgloss.Color(primaryLight)).
			Bold(true),

		statusBarStyle: r.NewStyle().
			Background(lipgloss.Color(neutral800)).
			Foreground(lipgloss.Color(neutral400)),

		statusKeyStyle: r.NewStyle().
			Background(lipgloss.Color(primaryDeep)).
			Foreground(lipgloss.Color(neutral100)).
			Padding(0, 1).
			Bold(true),

		statusSepStyle: r.NewStyle().
			Background(lipgloss.Color(neutral800)).
			Foreground(lipgloss.Color(neutral700)),

		statusScrollBarStyle: r.NewStyle().
			Background(lipgloss.Color(neutral800)).
			Foreground(lipgloss.Color(primaryDefault)),

		statusScrollDimStyle: r.NewStyle().
			Background(lipgloss.Color(neutral800)).
			Foreground(lipgloss.Color(neutral700)),

		statusScrollPctStyle: r.NewStyle().
			Background(lipgloss.Color(neutral800)).
			Foreground(lipgloss.Color(accentLight)),
	}
}
