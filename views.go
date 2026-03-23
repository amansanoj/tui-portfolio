package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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
		return m.renderURLPopup(base)
	}
	return base
}

// renderURLPopup overlays a centered box with the URL on top of the base view.
func (m Model) renderURLPopup(base string) string {
	url := m.showingURL

	maxURLWidth := m.windowWidth - 16
	if maxURLWidth < 20 {
		maxURLWidth = 20
	}

	// Wrap URL across lines so it stays copyable
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
	sb.WriteString(m.styles.taglineStyle.Render("dev & student") + "\n\n")

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

func (m Model) renderMainContent(mainWidth int) string {
	if m.selectedIndex < 0 || m.selectedIndex >= len(m.pageContents) {
		return m.styles.contentStyle.Render("No content available")
	}

	page := m.pageContents[m.selectedIndex]
	title := m.styles.subtitleStyle.Render(m.menuItems[m.selectedIndex].title)
	divider := m.styles.dimStyle.Render(strings.Repeat("─", 44))

	allLines := strings.Split(page.body, "\n")
	avail := m.availableContentHeight()
	scroll := m.contentScroll

	var visibleContent string
	switch m.selectedIndex {
	case 1:
		start := scroll
		end := scroll + avail
		if end > len(allLines) {
			end = len(allLines)
		}
		visible := allLines[start:end]
		for len(visible) < avail {
			visible = append(visible, "")
		}
		visibleContent = m.buildStyledAboutContent(visible, mainWidth)

	case 2:
		if len(m.projects) > 0 {
			visibleContent = m.buildStyledProjectContent(allLines, scroll, avail, mainWidth)
		} else {
			visibleContent = m.styles.contentStyle.Render(strings.Join(allLines, "\n"))
		}

	case 3:
		if len(m.certifications) > 0 {
			visibleContent = m.buildStyledCertContent(allLines, scroll, avail)
		} else {
			visibleContent = m.styles.contentStyle.Render(strings.Join(allLines, "\n"))
		}

	case 4:
		visibleContent = m.buildStyledContactContent(allLines)

	default:
		start := scroll
		end := scroll + avail
		if end > len(allLines) {
			end = len(allLines)
		}
		visible := allLines[start:end]
		for len(visible) < avail {
			visible = append(visible, "")
		}
		visibleContent = m.styles.contentStyle.Render(strings.Join(visible, "\n"))
	}

	return title + "\n" + divider + "\n\n" + visibleContent
}

func (m Model) buildStyledAboutContent(lines []string, mainWidth int) string {
	var result string
	wrapWidth := mainWidth - 3 - 3 - 2
	if wrapWidth < 40 {
		wrapWidth = 40
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			result += "\n"
			continue
		}

		if isSectionHeader(line) {
			result += m.styles.sectionHeaderStyle.Render(trimmed) + "\n"
			continue
		}

		if strings.HasPrefix(trimmed, "▸") {
			rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "▸"))
			if strings.Contains(rest, "(") && strings.Contains(rest, "–") {
				parenIdx := strings.Index(rest, " (")
				if parenIdx != -1 {
					name := rest[:parenIdx]
					date := rest[parenIdx:]
					result += m.styles.bulletStyle.Render("▸ ") + m.styles.contentStyle.Render(name) + m.styles.mutedStyle.Render(date) + "\n"
					continue
				}
			}
			result += m.styles.bulletStyle.Render("▸ ") + m.styles.contentStyle.Render(rest) + "\n"
			continue
		}

		if strings.HasPrefix(line, "  ") {
			if strings.Contains(trimmed, "–") && strings.Contains(trimmed, "(") {
				result += m.styles.mutedStyle.Render("  "+trimmed) + "\n"
			} else if strings.Contains(trimmed, "Team") || strings.Contains(trimmed, "Cadet") {
				result += "  " + m.styles.activeItemStyle.Render(trimmed) + "\n"
			} else {
				result += "  " + m.styles.contentStyle.Render(trimmed) + "\n"
			}
			continue
		}

		if strings.ContainsAny(line, "█░") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				name := parts[0]
				bar := parts[1]
				label := strings.Join(parts[2:], " ")
				result += m.styles.contentStyle.Render(fmt.Sprintf("%-10s", name)) +
					m.styles.contentStyle.Render(bar) +
					m.styles.mutedStyle.Render("  "+label) + "\n"
				continue
			}
		}

		wrapped := wordWrap(trimmed, wrapWidth)
		for _, wline := range wrapped {
			result += m.styles.contentStyle.Render(wline) + "\n"
		}
	}

	return strings.TrimRight(result, "\n")
}

func (m Model) buildStyledProjectContent(allLines []string, scroll, avail, mainWidth int) string {
	var result string
	linesEmitted := 0
	projectIndex := 0
	renderedLine := 0

	wrapWidth := mainWidth - 3 - 3 - 2
	if wrapWidth < 40 {
		wrapWidth = 40
	}

	for _, line := range allLines {
		if linesEmitted >= avail {
			break
		}

		trimmed := strings.TrimSpace(line)

		isTitleLine := strings.Contains(line, "(") &&
			strings.Contains(line, "–") &&
			strings.Contains(line, ")")

		if isTitleLine {
			idx := projectIndex
			projectIndex++
			titleRL := renderedLine
			dateRL := renderedLine + 1
			renderedLine += 2

			if titleRL >= scroll && linesEmitted < avail {
				parts := strings.SplitN(line, " (", 2)
				var indicator string
				if idx == m.selectedProject {
					indicator = m.styles.activeItemStyle.Render("● ")
				} else {
					indicator = m.styles.dimStyle.Render("○ ")
				}
				result += indicator + m.styles.projectTitleStyle.Render(parts[0]) + "\n"
				linesEmitted++
			}
			if dateRL >= scroll && linesEmitted < avail {
				parts := strings.SplitN(line, " (", 2)
				if len(parts) == 2 {
					date := strings.TrimSuffix(parts[1], ")")
					result += "  " + m.styles.mutedStyle.Render(date) + "\n"
					linesEmitted++
				}
			}
			continue
		}

		if trimmed == "" {
			if renderedLine >= scroll && linesEmitted < avail {
				result += "\n"
				linesEmitted++
			}
			renderedLine++
			continue
		}

		if strings.Contains(line, ",") && !strings.Contains(line, ", ") {
			if renderedLine >= scroll && linesEmitted < avail {
				tags := strings.Split(line, ",")
				var styled []string
				for _, tag := range tags {
					styled = append(styled, m.styles.activeItemStyle.Render(strings.TrimSpace(tag)))
				}
				result += "  " + strings.Join(styled, m.styles.dimStyle.Render(" · ")) + "\n"
				linesEmitted++
			}
			renderedLine++
			continue
		}

		if renderedLine >= scroll && linesEmitted < avail {
			wrapped := wordWrap(trimmed, wrapWidth)
			for _, wline := range wrapped {
				if linesEmitted >= avail {
					break
				}
				result += m.styles.contentStyle.Render("  "+wline) + "\n"
				linesEmitted++
			}
		}
		renderedLine++
	}

	return strings.TrimRight(result, "\n")
}

func (m Model) buildStyledCertContent(allLines []string, scroll, avail int) string {
	var result string
	linesEmitted := 0
	certIndex := 0
	renderedLine := 0

	for _, line := range allLines {
		if linesEmitted >= avail {
			break
		}

		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "CERT|||") {
			idx := certIndex
			certIndex++

			parts := strings.SplitN(trimmed, "|||", 3)
			title := ""
			date := ""
			if len(parts) >= 2 {
				title = parts[1]
			}
			if len(parts) >= 3 {
				date = parts[2]
			}

			var indicator string
			if idx == m.selectedCert {
				indicator = m.styles.activeItemStyle.Render("● ")
			} else {
				indicator = m.styles.dimStyle.Render("○ ")
			}

			linkIcon := m.styles.dimStyle.Render(" ·")
			if idx < len(m.certifications) && m.certifications[idx].URL != "" {
				linkIcon = m.styles.mutedStyle.Render(" ↗")
			}

			titleRL := renderedLine
			dateRL := renderedLine + 1
			renderedLine += 2

			if titleRL >= scroll && linesEmitted < avail {
				result += indicator + m.styles.projectTitleStyle.Render(title) + linkIcon + "\n"
				linesEmitted++
			}
			if dateRL >= scroll && linesEmitted < avail {
				if date != "" {
					result += "  " + m.styles.mutedStyle.Render(date) + "\n"
				} else {
					result += "\n"
				}
				linesEmitted++
			}
			continue
		}

		if strings.HasPrefix(trimmed, "ORG|||") {
			if renderedLine >= scroll && linesEmitted < avail {
				org := strings.TrimPrefix(trimmed, "ORG|||")
				result += "  " + m.styles.contentStyle.Render(org) + "\n"
				linesEmitted++
			}
			renderedLine++
			continue
		}

		if trimmed == "" {
			if renderedLine >= scroll && linesEmitted < avail {
				result += "\n"
				linesEmitted++
			}
			renderedLine++
			continue
		}

		if renderedLine >= scroll && linesEmitted < avail {
			result += "  " + m.styles.contentStyle.Render(trimmed) + "\n"
			linesEmitted++
		}
		renderedLine++
	}

	return strings.TrimRight(result, "\n")
}

func (m Model) buildStyledContactContent(allLines []string) string {
	var result string
	idx := 0

	for _, line := range allLines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			result += "\n"
			continue
		}

		if strings.HasPrefix(trimmed, "CONTACT|||") {
			parts := strings.SplitN(trimmed, "|||", 4)
			label := ""
			handle := ""
			if len(parts) >= 2 {
				label = parts[1]
			}
			if len(parts) >= 3 {
				handle = parts[2]
			}

			var indicator string
			if idx == m.selectedContact {
				indicator = m.styles.activeItemStyle.Render("● ")
			} else {
				indicator = m.styles.dimStyle.Render("○ ")
			}

			result += indicator +
				m.styles.mutedStyle.Render(fmt.Sprintf("%-10s", label)) +
				m.styles.contentStyle.Render(handle) + "\n"
			idx++
			continue
		}

		result += m.styles.contentStyle.Render(trimmed) + "\n"
	}

	return strings.TrimRight(result, "\n")
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
				hint("1-5", "navigate") + sep +
				hint("↑/↓", "select") + sep +
				hint("enter", "open") + sep +
				hint("q", "quit")
		} else {
			left = m.styles.statusBarStyle.Render(" ") +
				hint("1-5", "navigate") + sep +
				hint("↑/↓", "scroll") + sep +
				hint("q", "quit")
		}
	case 3:
		if len(m.certifications) > 0 {
			left = m.styles.statusBarStyle.Render(" ") +
				hint("1-5", "navigate") + sep +
				hint("↑/↓", "select") + sep +
				hint("enter", "open") + sep +
				hint("q", "quit")
		} else {
			left = m.styles.statusBarStyle.Render(" ") +
				hint("1-5", "navigate") + sep +
				hint("↑/↓", "scroll") + sep +
				hint("q", "quit")
		}
	case 4:
		left = m.styles.statusBarStyle.Render(" ") +
			hint("1-5", "navigate") + sep +
			hint("↑/↓", "select") + sep +
			hint("enter", "open") + sep +
			hint("q", "quit")
	default:
		left = m.styles.statusBarStyle.Render(" ") +
			hint("1-5", "navigate") + sep +
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
