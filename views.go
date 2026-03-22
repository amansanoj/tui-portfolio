package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	sidebarWidth := 20
	mainWidth := m.windowWidth - sidebarWidth - 6
	paneHeight := m.windowHeight - 3

	sidebar := m.renderSidebar()
	mainContent := m.renderMainContent(mainWidth)

	sidePane := sidePaneStyle.Width(sidebarWidth).Height(paneHeight).Render(sidebar)
	mainPane := mainPaneActiveStyle.Width(mainWidth).Height(paneHeight).Render(mainContent)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, sidePane, mainPane)
	statusBar := m.renderStatusBar()

	return lipgloss.JoinVertical(lipgloss.Left, topRow, statusBar)
}

func (m Model) renderSidebar() string {
	var sb strings.Builder

	sb.WriteString(nameStyle.Render("Aman Sanoj") + "\n")
	sb.WriteString(taglineStyle.Render("dev & student") + "\n\n")

	for i, item := range m.menuItems {
		num := fmt.Sprintf("%d", i+1)
		if i == m.selectedIndex {
			sb.WriteString(activeItemStyle.Render(num+" "+item.title) + "\n")
		} else {
			sb.WriteString(inactiveItemStyle.Render(num+" "+item.title) + "\n")
		}
	}

	return sb.String()
}

func (m Model) renderMainContent(mainWidth int) string {
	if m.selectedIndex < 0 || m.selectedIndex >= len(m.pageContents) {
		return contentStyle.Render("No content available")
	}

	page := m.pageContents[m.selectedIndex]
	title := subtitleStyle.Render(m.menuItems[m.selectedIndex].title)
	divider := dimStyle.Render(strings.Repeat("─", 44))

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
			visibleContent = contentStyle.Render(strings.Join(allLines, "\n"))
		}

	case 3:
		if len(m.certifications) > 0 {
			visibleContent = m.buildStyledCertContent(allLines, scroll, avail)
		} else {
			visibleContent = contentStyle.Render(strings.Join(allLines, "\n"))
		}

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
		visibleContent = contentStyle.Render(strings.Join(visible, "\n"))
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
			result += sectionHeaderStyle.Render(trimmed) + "\n"
			continue
		}

		if strings.HasPrefix(trimmed, "▸") {
			rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "▸"))
			if strings.Contains(rest, "(") && strings.Contains(rest, "–") {
				parenIdx := strings.Index(rest, " (")
				if parenIdx != -1 {
					name := rest[:parenIdx]
					date := rest[parenIdx:]
					result += bulletStyle.Render("▸ ") + contentStyle.Render(name) + mutedStyle.Render(date) + "\n"
					continue
				}
			}
			result += bulletStyle.Render("▸ ") + contentStyle.Render(rest) + "\n"
			continue
		}

		if strings.HasPrefix(line, "  ") {
			if strings.Contains(trimmed, "–") && strings.Contains(trimmed, "(") {
				result += mutedStyle.Render("  "+trimmed) + "\n"
			} else if strings.Contains(trimmed, "Team") || strings.Contains(trimmed, "Cadet") {
				result += "  " + activeItemStyle.Render(trimmed) + "\n"
			} else {
				result += "  " + contentStyle.Render(trimmed) + "\n"
			}
			continue
		}

		if strings.ContainsAny(line, "█░") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				name := parts[0]
				bar := parts[1]
				label := strings.Join(parts[2:], " ")
				result += contentStyle.Render(fmt.Sprintf("%-10s", name)) +
					contentStyle.Render(bar) +
					mutedStyle.Render("  "+label) + "\n"
				continue
			}
		}

		wrapped := wordWrap(trimmed, wrapWidth)
		for _, wline := range wrapped {
			result += contentStyle.Render(wline) + "\n"
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
					indicator = activeItemStyle.Render("● ")
				} else {
					indicator = dimStyle.Render("○ ")
				}
				result += indicator + projectTitleStyle.Render(parts[0]) + "\n"
				linesEmitted++
			}
			if dateRL >= scroll && linesEmitted < avail {
				parts := strings.SplitN(line, " (", 2)
				if len(parts) == 2 {
					date := strings.TrimSuffix(parts[1], ")")
					result += "  " + mutedStyle.Render(date) + "\n"
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
					styled = append(styled, activeItemStyle.Render(strings.TrimSpace(tag)))
				}
				result += "  " + strings.Join(styled, dimStyle.Render(" · ")) + "\n"
				linesEmitted++
			}
			renderedLine++
			continue
		}

		if renderedLine >= scroll && linesEmitted < avail {
			result += contentStyle.Render("  "+trimmed) + "\n"
			linesEmitted++
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
				indicator = activeItemStyle.Render("● ")
			} else {
				indicator = dimStyle.Render("○ ")
			}

			linkIcon := dimStyle.Render(" ·")
			if idx < len(m.certifications) && m.certifications[idx].URL != "" {
				linkIcon = mutedStyle.Render(" ↗")
			}

			titleRL := renderedLine
			dateRL := renderedLine + 1
			renderedLine += 2

			if titleRL >= scroll && linesEmitted < avail {
				result += indicator + projectTitleStyle.Render(title) + linkIcon + "\n"
				linesEmitted++
			}
			if dateRL >= scroll && linesEmitted < avail {
				if date != "" {
					result += "  " + mutedStyle.Render(date) + "\n"
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
				result += "  " + contentStyle.Render(org) + "\n"
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

		// anything else
		if renderedLine >= scroll && linesEmitted < avail {
			result += "  " + contentStyle.Render(trimmed) + "\n"
			linesEmitted++
		}
		renderedLine++
	}

	return strings.TrimRight(result, "\n")
}

func (m Model) renderStatusBar() string {
	sep := statusSepStyle.Render("  │  ")

	hint := func(key, desc string) string {
		return statusKeyStyle.Render(" "+key+" ") +
			statusBarStyle.Render(" "+desc)
	}

	var left string
	switch m.selectedIndex {
	case 2:
		if len(m.projects) > 0 {
			left = statusBarStyle.Render(" ") +
				hint("1-5", "navigate") + sep +
				hint("↑/↓", "select") + sep +
				hint("enter", "open") + sep +
				hint("q", "quit")
		} else {
			left = statusBarStyle.Render(" ") +
				hint("1-5", "navigate") + sep +
				hint("↑/↓", "scroll") + sep +
				hint("q", "quit")
		}
	case 3:
		if len(m.certifications) > 0 {
			left = statusBarStyle.Render(" ") +
				hint("1-5", "navigate") + sep +
				hint("↑/↓", "select") + sep +
				hint("enter", "open") + sep +
				hint("q", "quit")
		} else {
			left = statusBarStyle.Render(" ") +
				hint("1-5", "navigate") + sep +
				hint("↑/↓", "scroll") + sep +
				hint("q", "quit")
		}
	default:
		left = statusBarStyle.Render(" ") +
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
		bar := statusScrollBarStyle.Render(strings.Repeat("█", filled)) +
			statusScrollDimStyle.Render(strings.Repeat("░", 10-filled))
		pctStr := statusScrollPctStyle.Render(fmt.Sprintf(" %3d%% ", pct))
		right = bar + pctStr
	} else {
		right = statusBarStyle.Render("  ")
	}

	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	gap := m.windowWidth - leftWidth - rightWidth
	if gap < 0 {
		gap = 0
	}
	padding := statusBarStyle.Render(strings.Repeat(" ", gap))

	return left + padding + right
}

func openURL(url string) {
	var cmd *exec.Cmd
	switch {
	case os.Getenv("BROWSER") != "":
		cmd = exec.Command(os.Getenv("BROWSER"), url)
	case isCommandAvailable("xdg-open"):
		cmd = exec.Command("xdg-open", url)
	case isCommandAvailable("open"):
		cmd = exec.Command("open", url)
	}
	if cmd != nil {
		_ = cmd.Start()
	}
}

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
