package pages

import (
	"fmt"
	"strings"
	"unicode"
)

var monthTokens = []string{
	"jan", "feb", "mar", "apr", "may", "jun",
	"jul", "aug", "sep", "sept", "oct", "nov", "dec",
	"january", "february", "march", "april", "june", "july",
	"august", "september", "october", "november", "december",
}

type Theme struct {
	SectionHeader func(string) string
	Bullet        func(string) string
	Content       func(string) string
	Muted         func(string) string
	Active        func(string) string
	Dim           func(string) string
	ProjectTitle  func(string) string
}

func Wrap(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}
	if len(text) <= maxWidth {
		return []string{text}
	}
	var lines []string
	words := strings.Fields(text)
	current := ""
	for _, word := range words {
		if current == "" {
			current = word
		} else if len(current)+1+len(word) <= maxWidth {
			current += " " + word
		} else {
			lines = append(lines, current)
			current = word
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func ClampVisibleLines(allLines []string, scroll, avail int) []string {
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

func renderBulletsInline(line string, t Theme) string {
	if strings.Contains(line, "▸") {
		parts := strings.Split(line, "▸")
		styledLine := renderMutedDateBrackets(parts[0], t)
		for i := 1; i < len(parts); i++ {
			styledLine += t.Bullet("▸") + renderMutedDateBrackets(parts[i], t)
		}
		return styledLine
	}

	if leading, digit, rest, ok := splitLeadingDigitLine(line); ok {
		return t.Content(leading) + t.Active(digit) + renderMutedDateBrackets(rest, t)
	}

	return renderMutedDateBrackets(line, t)
}

func renderMutedDateBrackets(line string, t Theme) string {
	if !strings.Contains(line, "(") || !strings.Contains(line, ")") {
		return t.Content(line)
	}

	var out strings.Builder
	idx := 0

	for idx < len(line) {
		open := strings.Index(line[idx:], "(")
		if open == -1 {
			out.WriteString(t.Content(line[idx:]))
			break
		}

		open += idx
		if open > idx {
			out.WriteString(t.Content(line[idx:open]))
		}

		close := strings.Index(line[open:], ")")
		if close == -1 {
			out.WriteString(t.Content(line[open:]))
			break
		}

		close += open
		segment := line[open : close+1]
		inner := line[open+1 : close]

		if isDateLikeBracket(inner) {
			out.WriteString(t.Muted(segment))
		} else {
			out.WriteString(t.Content(segment))
		}

		idx = close + 1
	}

	return out.String()
}

func isDateLikeBracket(value string) bool {
	v := strings.ToLower(strings.TrimSpace(value))
	if v == "" {
		return false
	}

	hasDigit := false
	for _, r := range v {
		if unicode.IsDigit(r) {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return false
	}

	hasMonth := false
	for _, m := range monthTokens {
		if strings.Contains(v, m) {
			hasMonth = true
			break
		}
	}

	hasRangeMarker := strings.Contains(v, "present") ||
		strings.Contains(v, "current") ||
		strings.Contains(v, "-") ||
		strings.Contains(v, "–")

	return hasMonth || hasRangeMarker
}

func splitLeadingDigitLine(line string) (leading, digit, rest string, ok bool) {
	i := 0
	for i < len(line) {
		r := rune(line[i])
		if r != ' ' && r != '\t' {
			break
		}
		i++
	}

	if i >= len(line) {
		return "", "", "", false
	}

	r := rune(line[i])
	if !unicode.IsDigit(r) {
		return "", "", "", false
	}

	if i+1 >= len(line) || line[i+1] != ' ' {
		return "", "", "", false
	}

	return line[:i], line[i : i+1], line[i+1:], true
}

func RenderDefault(visible []string, t Theme) string {
	var styledLines []string
	for _, line := range visible {
		styledLines = append(styledLines, renderBulletsInline(line, t))
	}
	return strings.Join(styledLines, "\n")
}

func isSectionHeader(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return false
	}
	if strings.ContainsAny(trimmed, "█░") {
		return false
	}
	if strings.HasPrefix(trimmed, "▸") || strings.HasPrefix(trimmed, "●") ||
		strings.HasPrefix(trimmed, "○") || strings.HasPrefix(trimmed, "█") ||
		strings.HasPrefix(trimmed, " ") || strings.HasPrefix(trimmed, "│") ||
		strings.HasPrefix(trimmed, "└") || strings.HasPrefix(trimmed, "┌") {
		return false
	}
	words := strings.Fields(trimmed)
	if len(words) > 4 {
		return false
	}
	if strings.ContainsAny(trimmed, ".,@•·—") {
		return false
	}
	upper := strings.ToUpper(trimmed)
	if trimmed == upper {
		return true
	}
	if len(trimmed) > 0 && trimmed[0] >= 'A' && trimmed[0] <= 'Z' && len(words) <= 3 {
		return true
	}
	return false
}

func RenderAbout(lines []string, mainWidth int, t Theme) string {
	var result string
	wrapWidth := mainWidth - 2 - 2 - 2
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
			result += t.SectionHeader(trimmed) + "\n"
			continue
		}

		if strings.HasPrefix(trimmed, "▸") {
			rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "▸"))
			result += t.Bullet("▸ ") + renderMutedDateBrackets(rest, t) + "\n"
			continue
		}

		if strings.HasPrefix(line, "  ") {
			if strings.Contains(trimmed, "–") && strings.Contains(trimmed, "(") {
				result += t.Muted("  "+trimmed) + "\n"
			} else if strings.Contains(trimmed, "Team") || strings.Contains(trimmed, "Cadet") {
				result += "  " + t.Active(trimmed) + "\n"
			} else {
				result += "  " + renderMutedDateBrackets(trimmed, t) + "\n"
			}
			continue
		}

		if strings.ContainsAny(line, "█░") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				name := parts[0]
				bar := parts[1]
				label := strings.Join(parts[2:], " ")
				result += t.Content(fmt.Sprintf("%-10s", name)) +
					t.Content(bar) +
					t.Muted("  "+label) + "\n"
				continue
			}
		}

		wrapped := Wrap(trimmed, wrapWidth)
		for _, wline := range wrapped {
			result += renderBulletsInline(wline, t) + "\n"
		}
	}

	return strings.TrimRight(result, "\n")
}

func RenderProjects(allLines []string, scroll, avail, mainWidth, selectedProject int, t Theme) string {
	var result string
	linesEmitted := 0
	projectIndex := 0
	renderedLine := 0

	wrapWidth := mainWidth - 2 - 2 - 2
	if wrapWidth < 40 {
		wrapWidth = 40
	}

	for _, line := range allLines {
		if linesEmitted >= avail {
			break
		}

		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "PROJ|||") {
			idx := projectIndex
			projectIndex++
			titleRL := renderedLine
			renderedLine++

			parts := strings.SplitN(trimmed, "|||", 3)
			title := ""
			date := ""
			if len(parts) >= 2 {
				title = parts[1]
			}
			if len(parts) >= 3 {
				date = strings.TrimSpace(parts[2])
			}

			if titleRL >= scroll && linesEmitted < avail {
				indicator := t.Dim("○ ")
				if idx == selectedProject {
					indicator = t.Active("● ")
				}

				if date != "" {
					result += indicator + t.ProjectTitle(title) + " " + t.Muted("("+date+")") + "\n"
				} else {
					result += indicator + t.ProjectTitle(title) + "\n"
				}
				linesEmitted++
			}
			continue
		}

		isTitleLine := strings.Contains(line, " (") && strings.HasSuffix(trimmed, ")")

		if isTitleLine {
			idx := projectIndex
			projectIndex++
			titleRL := renderedLine
			renderedLine++

			if titleRL >= scroll && linesEmitted < avail {
				parts := strings.SplitN(line, " (", 2)
				indicator := t.Dim("○ ")
				if idx == selectedProject {
					indicator = t.Active("● ")
				}
				if len(parts) == 2 {
					date := "(" + parts[1]
					result += indicator + t.ProjectTitle(parts[0]) + " " + t.Muted(date) + "\n"
				} else {
					result += indicator + t.ProjectTitle(line) + "\n"
				}
				linesEmitted++
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
					styled = append(styled, t.Active(strings.TrimSpace(tag)))
				}
				result += "  " + strings.Join(styled, t.Dim(" · ")) + "\n"
				linesEmitted++
			}
			renderedLine++
			continue
		}

		wrapped := Wrap(trimmed, wrapWidth)
		for i, wline := range wrapped {
			lineRL := renderedLine + i
			if lineRL >= scroll && linesEmitted < avail {
				result += t.Content("  "+wline) + "\n"
				linesEmitted++
			}
		}
		renderedLine += len(wrapped)
	}

	return strings.TrimRight(result, "\n")
}

func RenderCerts(allLines []string, scroll, avail, selectedCert int, certs []CertData, t Theme) string {
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

			indicator := t.Dim("○ ")
			if idx == selectedCert {
				indicator = t.Active("● ")
			}

			linkIcon := t.Dim(" ·")
			if idx < len(certs) && certs[idx].URL != "" {
				linkIcon = t.Muted(" ↗")
			}

			titleRL := renderedLine
			dateRL := renderedLine + 1
			renderedLine += 2

			if titleRL >= scroll && linesEmitted < avail {
				result += indicator + t.ProjectTitle(title) + linkIcon + "\n"
				linesEmitted++
			}
			if dateRL >= scroll && linesEmitted < avail {
				if date != "" {
					result += "  " + t.Muted(date) + "\n"
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
				result += "  " + t.Content(org) + "\n"
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
			result += "  " + t.Content(trimmed) + "\n"
			linesEmitted++
		}
		renderedLine++
	}

	return strings.TrimRight(result, "\n")
}

func RenderContact(allLines []string, selectedContact int, t Theme) string {
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

			indicator := t.Dim("○ ")
			if idx == selectedContact {
				indicator = t.Active("● ")
			}

			result += indicator +
				t.Muted(fmt.Sprintf("%-10s", label)) +
				t.Content(handle) + "\n"
			idx++
			continue
		}

		result += t.Content(trimmed) + "\n"
	}

	return strings.TrimRight(result, "\n")
}

func StatusWrapWidth(mainWidth int) int {
	wrapWidth := mainWidth - 2 - 2 - 2
	if wrapWidth < 30 {
		wrapWidth = 30
	}
	return wrapWidth
}

func StatusLineCount(lines []string, mainWidth int) int {
	wrapWidth := StatusWrapWidth(mainWidth)
	total := 0

	for _, line := range lines {
		if line == "" {
			total++
			continue
		}
		wrapped := Wrap(line, wrapWidth)
		if len(wrapped) == 0 {
			total++
			continue
		}
		total += len(wrapped)
	}

	return total
}

func RenderStatus(lines []string, mainWidth, scroll, avail int, t Theme) string {
	wrapWidth := StatusWrapWidth(mainWidth)

	var expanded []string
	for _, line := range lines {
		if line == "" {
			expanded = append(expanded, "")
			continue
		}
		wrapped := Wrap(line, wrapWidth)
		if len(wrapped) == 0 {
			expanded = append(expanded, "")
			continue
		}
		expanded = append(expanded, wrapped...)
	}

	visible := ClampVisibleLines(expanded, scroll, avail)

	var out []string
	for _, line := range visible {
		if line == "" {
			out = append(out, "")
			continue
		}
		out = append(out, t.Content(line))
	}

	return strings.Join(out, "\n")
}
