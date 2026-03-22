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

var (
	sidePaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(neutral600)).
			Foreground(lipgloss.Color(neutral100)).
			Padding(1, 2)

	mainPaneActiveStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(primaryDefault)).
				Foreground(lipgloss.Color(neutral100)).
				Padding(1, 3)

	activeItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(accentDefault)).
			Bold(true)

	inactiveItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(neutral600))

	nameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(neutral100)).
			Bold(true)

	taglineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(neutral400))

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(primaryLight)).
			Bold(true)

	sectionHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(primaryLight)).
				Bold(true)

	contentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(neutral100))

	mutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(neutral600))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(neutral700))

	bulletStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(accentDefault))

	projectTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(primaryLight)).
				Bold(true)

	statusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(neutral800)).
			Foreground(lipgloss.Color(neutral400))

	statusKeyStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(neutral700)).
			Foreground(lipgloss.Color(neutral100)).
			Padding(0, 1).
			Bold(true)

	statusSepStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(neutral800)).
			Foreground(lipgloss.Color(neutral700))

	statusScrollBarStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(neutral800)).
				Foreground(lipgloss.Color(primaryLight))

	statusScrollDimStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(neutral800)).
				Foreground(lipgloss.Color(neutral700))

	statusScrollPctStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(neutral800)).
				Foreground(lipgloss.Color(neutral400))
)
