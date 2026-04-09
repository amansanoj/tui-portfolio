package main

import "github.com/charmbracelet/lipgloss"

const (
	primaryDefault = "rgb(61, 143, 209)"
	primaryLight   = "rgb(90, 159, 216)"
	primaryDeep    = "rgb(39, 108, 165)"

	secondaryDefault = "rgb(240, 176, 122)"
	secondaryStrong  = "rgb(230, 117, 25)"

	accentDefault = "rgb(168, 158, 105)"
	accentLight   = "rgb(181, 171, 125)"

	neutral100 = "rgb(230, 230, 230)"
	neutral400 = "rgb(153, 153, 153)"
	neutral600 = "rgb(102, 102, 102)"
	neutral700 = "rgb(77, 77, 77)"
	neutral800 = "rgb(51, 51, 51)"
	neutral900 = "rgb(26, 26, 26)"
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
