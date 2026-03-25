package main

import "strings"

func (m Model) clampVisibleLines(allLines []string, scroll, avail int) []string {
	start := scroll
	end := scroll + avail
	if start < 0 {
		start = 0
	}
	if end > len(allLines) {
		end = len(allLines)
	}
	if start > len(allLines) {
		start = len(allLines)
	}

	visible := allLines[start:end]
	for len(visible) < avail {
		visible = append(visible, "")
	}
	return visible
}

func (m Model) renderAboutPage(allLines []string, scroll, avail, mainWidth int) string {
	visible := m.clampVisibleLines(allLines, scroll, avail)
	return m.buildStyledAboutContent(visible, mainWidth)
}

func (m Model) renderProjectsPage(allLines []string, scroll, avail, mainWidth int) string {
	if len(m.projects) > 0 {
		return m.buildStyledProjectContent(allLines, scroll, avail, mainWidth)
	}
	return m.styles.contentStyle.Render(strings.Join(allLines, "\n"))
}

func (m Model) renderCertsPage(allLines []string, scroll, avail int) string {
	if len(m.certifications) > 0 {
		return m.buildStyledCertContent(allLines, scroll, avail)
	}
	return m.styles.contentStyle.Render(strings.Join(allLines, "\n"))
}

func (m Model) renderContactPage(allLines []string) string {
	return m.buildStyledContactContent(allLines)
}

func (m Model) renderDefaultPage(allLines []string, scroll, avail int) string {
	visible := m.clampVisibleLines(allLines, scroll, avail)
	var styledLines []string
	for _, line := range visible {
		if strings.Contains(line, "▸") {
			parts := strings.Split(line, "▸")
			styledLine := m.styles.contentStyle.Render(parts[0])
			for i := 1; i < len(parts); i++ {
				styledLine += m.styles.bulletStyle.Render("▸") + m.styles.contentStyle.Render(parts[i])
			}
			styledLines = append(styledLines, styledLine)
		} else {
			styledLines = append(styledLines, m.styles.contentStyle.Render(line))
		}
	}
	return strings.Join(styledLines, "\n")
}
