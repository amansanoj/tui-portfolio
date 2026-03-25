package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/amansanoj/tui-portfolio/pages"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) pageTheme() pages.Theme {
	return pages.Theme{
		SectionHeader: func(s string) string { return m.styles.sectionHeaderStyle.Render(s) },
		Bullet:        func(s string) string { return m.styles.bulletStyle.Render(s) },
		Content:       func(s string) string { return m.styles.contentStyle.Render(s) },
		Muted:         func(s string) string { return m.styles.mutedStyle.Render(s) },
		Active:        func(s string) string { return m.styles.activeItemStyle.Render(s) },
		Dim:           func(s string) string { return m.styles.dimStyle.Render(s) },
		ProjectTitle:  func(s string) string { return m.styles.projectTitleStyle.Render(s) },
	}
}

func (m Model) View() string {
	if m.windowWidth < 90 || m.windowHeight < 15 {
		msg := "Terminal too small — please resize to at least 90×15"
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(neutral400)).
			Padding(1, 2).
			Render(msg)
	}

	sidebarWidth := 20
	mainWidth := m.windowWidth - sidebarWidth - 4
	paneHeight := m.windowHeight - 3

	sidebar := m.renderSidebar()
	mainContent := m.renderMainContent(mainWidth)

	sidePane := m.styles.sidePaneStyle.Width(sidebarWidth).Height(paneHeight).Render(sidebar)
	mainPane := m.styles.mainPaneActiveStyle.Width(mainWidth).Height(paneHeight).Render(mainContent)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, sidePane, mainPane)
	statusBar := m.renderStatusBar()
	base := lipgloss.JoinVertical(lipgloss.Left, topRow, statusBar)

	if m.showingURL != "" {
		return m.renderURLPopup()
	}
	return base
}

func (m Model) renderURLPopup() string {
	url := m.showingURL

	maxURLWidth := m.windowWidth - 16
	if maxURLWidth < 20 {
		maxURLWidth = 20
	}

	var urlLines []string
	for len(url) > maxURLWidth {
		urlLines = append(urlLines, url[:maxURLWidth])
		url = url[maxURLWidth:]
	}
	if len(url) > 0 {
		urlLines = append(urlLines, url)
	}

	var urlRendered string
	for _, l := range urlLines {
		urlRendered += m.styles.projectTitleStyle.Render(l) + "\n"
	}
	urlRendered = strings.TrimRight(urlRendered, "\n")

	popupWidth := maxURLWidth + 8
	if popupWidth > m.windowWidth-4 {
		popupWidth = m.windowWidth - 4
	}

	content := m.styles.mutedStyle.Render("Open this link in your browser:") + "\n\n" +
		urlRendered + "\n\n" +
		m.styles.dimStyle.Render("press any key to dismiss")

	popup := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(primaryDefault)).
		Background(lipgloss.Color(neutral800)).
		Padding(1, 3).
		Width(popupWidth).
		Render(content)

	return lipgloss.Place(
		m.windowWidth,
		m.windowHeight,
		lipgloss.Center,
		lipgloss.Center,
		popup,
		lipgloss.WithWhitespaceBackground(lipgloss.Color(neutral900)),
	)
}

func (m Model) renderSidebar() string {
	var sb strings.Builder

	sb.WriteString(m.styles.nameStyle.Render("Aman Sanoj") + "\n")
	sb.WriteString(m.styles.dimStyle.Render(strings.Repeat("─", 16)) + "\n\n")

	for i, item := range m.menuItems {
		num := fmt.Sprintf("%d", i+1)
		if i == m.selectedIndex {
			sb.WriteString(m.styles.activeItemStyle.Render(num+" "+item.title) + "\n")
		} else {
			sb.WriteString(m.styles.inactiveItemStyle.Render(num+" "+item.title) + "\n")
		}
	}

	return sb.String()
}

func (m Model) statusPageRawLines() []string {
	snapshot := appContentStore.Snapshot()
	refreshValue := os.Getenv(contentRefreshEnvVar)
	if refreshValue == "" {
		refreshValue = "300"
	}

	cacheState := "warm"
	cacheAge := "unknown"
	lastLoaded := "not loaded yet"
	if !snapshot.Ready {
		cacheState = "cold"
	} else {
		age := time.Since(snapshot.LoadedAt)
		if age < 0 {
			age = 0
		}
		cacheAge = age.Round(time.Second).String()
		lastLoaded = snapshot.LoadedAt.Local().Format(time.RFC1123)
	}

	return pages.BuildStatusLines(pages.StatusData{
		CacheState:          cacheState,
		Refreshing:          snapshot.Refreshing,
		CacheAge:            cacheAge,
		LastLoaded:          lastLoaded,
		ProjectsCount:       len(snapshot.Projects),
		CertificationsCount: len(snapshot.Certifications),
		RefreshInterval:     refreshValue + "s",
		ListenAddress:       envWithDefault("APP_ADDR", defaultSSHAddress),
		HostKeyPath:         envWithDefault("HOST_KEY_PATH", defaultHostKeyPath),
	})
}

func (m Model) statusPageLineCount(mainWidth int) int {
	return pages.StatusLineCount(m.statusPageRawLines(), mainWidth)
}

func envWithDefault(name, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	return value
}

func (m Model) renderMainContent(mainWidth int) string {
	if m.selectedIndex < 0 || m.selectedIndex >= len(m.pageContents) {
		return m.styles.contentStyle.Render("No content available")
	}

	page := m.pageContents[m.selectedIndex]
	title := m.styles.subtitleStyle.Render(m.menuItems[m.selectedIndex].title)
	dividerWidth := mainWidth - 4
	if dividerWidth < 1 {
		dividerWidth = 1
	}
	divider := m.styles.dimStyle.Render(strings.Repeat("─", dividerWidth))

	allLines := strings.Split(page.body, "\n")
	avail := m.availableContentHeight()
	scroll := m.contentScroll
	theme := m.pageTheme()

	var visibleContent string
	switch m.selectedIndex {
	case 1:
		visibleContent = pages.RenderAbout(pages.ClampVisibleLines(allLines, scroll, avail), mainWidth, theme)
	case 2:
		if len(m.projects) > 0 {
			visibleContent = pages.RenderProjects(allLines, scroll, avail, mainWidth, m.selectedProject, theme)
		} else {
			visibleContent = m.styles.contentStyle.Render(strings.Join(allLines, "\n"))
		}
	case 3:
		if len(m.certifications) > 0 {
			visibleContent = pages.RenderCerts(allLines, scroll, avail, m.selectedCert, toPageCerts(m.certifications), theme)
		} else {
			visibleContent = m.styles.contentStyle.Render(strings.Join(allLines, "\n"))
		}
	case 4:
		visibleContent = pages.RenderContact(allLines, m.selectedContact, theme)
	case 5:
		visibleContent = pages.RenderStatus(m.statusPageRawLines(), mainWidth, scroll, avail, theme)
	default:
		visibleContent = pages.RenderDefault(pages.ClampVisibleLines(allLines, scroll, avail), theme)
	}

	return title + "\n" + divider + "\n\n" + visibleContent
}

func (m Model) renderStatusBar() string {
	sep := m.styles.statusSepStyle.Render("  │  ")

	hint := func(key, desc string) string {
		return m.styles.statusKeyStyle.Render(" "+key+" ") +
			m.styles.statusBarStyle.Render(" "+desc)
	}

	var left string
	switch m.selectedIndex {
	case 2:
		if len(m.projects) > 0 {
			left = m.styles.statusBarStyle.Render(" ") +
				hint("1-6", "navigate") + sep +
				hint("↑/↓", "select") + sep +
				hint("enter", "open") + sep +
				hint("q", "quit")
		} else {
			left = m.styles.statusBarStyle.Render(" ") +
				hint("1-6", "navigate") + sep +
				hint("↑/↓", "scroll") + sep +
				hint("q", "quit")
		}
	case 3:
		if len(m.certifications) > 0 {
			left = m.styles.statusBarStyle.Render(" ") +
				hint("1-6", "navigate") + sep +
				hint("↑/↓", "select") + sep +
				hint("enter", "open") + sep +
				hint("q", "quit")
		} else {
			left = m.styles.statusBarStyle.Render(" ") +
				hint("1-6", "navigate") + sep +
				hint("↑/↓", "scroll") + sep +
				hint("q", "quit")
		}
	case 4:
		left = m.styles.statusBarStyle.Render(" ") +
			hint("1-6", "navigate") + sep +
			hint("↑/↓", "select") + sep +
			hint("enter", "open") + sep +
			hint("q", "quit")
	default:
		left = m.styles.statusBarStyle.Render(" ") +
			hint("1-6", "navigate") + sep +
			hint("↑/↓", "scroll") + sep +
			hint("q", "quit")
	}

	maxScroll := m.getMaxContentScroll()
	var right string
	if maxScroll > 0 {
		pct := int((float32(m.contentScroll) / float32(maxScroll)) * 100)
		filled := int((float32(m.contentScroll) / float32(maxScroll)) * 10)
		if pct > 100 {
			pct = 100
		}
		if filled > 10 {
			filled = 10
		}
		if filled < 0 {
			filled = 0
		}
		bar := m.styles.statusScrollBarStyle.Render(strings.Repeat("█", filled)) +
			m.styles.statusScrollDimStyle.Render(strings.Repeat("░", 10-filled))
		pctStr := m.styles.statusScrollPctStyle.Render(fmt.Sprintf(" %3d%% ", pct))
		right = bar + pctStr
	} else {
		right = m.styles.statusBarStyle.Render("  ")
	}

	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	gap := m.windowWidth - leftWidth - rightWidth
	if gap < 0 {
		gap = 0
	}
	padding := m.styles.statusBarStyle.Render(strings.Repeat(" ", gap))

	return left + padding + right
}
