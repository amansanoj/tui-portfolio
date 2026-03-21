package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const notionDatabaseID = "32acb49d4dc9804ab1b5f3ccf42c375c"
const notionAPIVersion = "2022-06-28"

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

type MenuItem struct {
	title string
}

type NotionPage struct {
	ID         string                 `json:"id"`
	Properties map[string]interface{} `json:"properties"`
}

type NotionResponse struct {
	Results []NotionPage `json:"results"`
}

type Project struct {
	Name        string
	Description string
	Date        string
	URL         string
	TechStack   string
}

func fetchProjectsFromNotion() []Project {
	apiKey := os.Getenv("NOTION_API_KEY")
	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "Error: NOTION_API_KEY environment variable not set\n")
		return []Project{}
	}

	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", notionDatabaseID)
	payload := []byte(`{}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return []Project{}
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Notion-Version", notionAPIVersion)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []Project{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "Notion API error: %s\n", string(body))
		return []Project{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Project{}
	}

	var notionResp NotionResponse
	if err := json.Unmarshal(body, &notionResp); err != nil {
		return []Project{}
	}

	var projects []Project
	for _, page := range notionResp.Results {
		project := Project{
			Name:        extractStringProperty(page.Properties, "Name"),
			Description: extractStringProperty(page.Properties, "Description"),
			Date:        extractDateProperty(page.Properties, "Date"),
			URL:         extractURLProperty(page.Properties, "Project URL"),
			TechStack:   extractStringProperty(page.Properties, "Tech Stack"),
		}
		if project.Name != "" {
			projects = append(projects, project)
		}
	}
	return projects
}

func extractStringProperty(props map[string]interface{}, propName string) string {
	if prop, exists := props[propName]; exists {
		propMap := prop.(map[string]interface{})
		propType := propMap["type"].(string)
		switch propType {
		case "title":
			if titleArr, ok := propMap["title"].([]interface{}); ok && len(titleArr) > 0 {
				if titleObj, ok := titleArr[0].(map[string]interface{}); ok {
					if text, ok := titleObj["plain_text"].(string); ok {
						return text
					}
				}
			}
		case "rich_text":
			if richArr, ok := propMap["rich_text"].([]interface{}); ok && len(richArr) > 0 {
				if richObj, ok := richArr[0].(map[string]interface{}); ok {
					if text, ok := richObj["plain_text"].(string); ok {
						return text
					}
				}
			}
		}
	}
	return ""
}

func extractDateProperty(props map[string]interface{}, propName string) string {
	if prop, exists := props[propName]; exists {
		propMap := prop.(map[string]interface{})
		if dateProp, ok := propMap["date"].(map[string]interface{}); ok {
			if start, ok := dateProp["start"].(string); ok {
				if end, ok := dateProp["end"].(string); ok {
					return fmt.Sprintf("%s – %s", start, end)
				}
				return start
			}
		}
	}
	return ""
}

func extractURLProperty(props map[string]interface{}, propName string) string {
	if prop, exists := props[propName]; exists {
		propMap := prop.(map[string]interface{})
		if url, ok := propMap["url"].(string); ok {
			return url
		}
	}
	return ""
}

func formatDateRange(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	parts := strings.Split(dateStr, " – ")
	if len(parts) == 2 {
		startDate := parseMonthYear(parts[0])
		endDate := parseMonthYear(parts[1])
		if startDate != "" && endDate != "" {
			return fmt.Sprintf("%s – %s", startDate, endDate)
		}
		if startDate != "" {
			return startDate
		}
	} else if len(parts) == 1 {
		date := parseMonthYear(parts[0])
		if date != "" {
			return date
		}
	}
	return ""
}

func parseMonthYear(dateStr string) string {
	dateStr = strings.TrimSpace(dateStr)
	months := []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	parts := strings.Split(dateStr, "-")
	if len(parts) >= 2 {
		year := parts[0]
		month := parts[1]
		monthNum := 0
		fmt.Sscanf(month, "%d", &monthNum)
		if monthNum > 0 && monthNum < len(months) {
			return fmt.Sprintf("%s %s", months[monthNum], year)
		}
	}
	return ""
}

func wordWrap(text string, maxWidth int) []string {
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

// isSectionHeader returns true for lines that are short, title-cased or all-caps
// section headers like "Experience", "SKILLS", "Education" etc.
func isSectionHeader(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return false
	}
	// Must not start with a bullet or special char
	if strings.HasPrefix(trimmed, "▸") || strings.HasPrefix(trimmed, "●") ||
		strings.HasPrefix(trimmed, "○") || strings.HasPrefix(trimmed, "█") ||
		strings.HasPrefix(trimmed, " ") || strings.HasPrefix(trimmed, "│") ||
		strings.HasPrefix(trimmed, "└") || strings.HasPrefix(trimmed, "┌") {
		return false
	}
	// Must be short (single word or two words, no punctuation mid-line)
	words := strings.Fields(trimmed)
	if len(words) > 4 {
		return false
	}
	// Must not contain sentence-like punctuation
	if strings.ContainsAny(trimmed, ".,@•·—") {
		return false
	}
	// All-caps OR title-cased first word
	upper := strings.ToUpper(trimmed)
	if trimmed == upper {
		return true
	}
	// Title case: first letter of first word is uppercase
	if len(trimmed) > 0 && trimmed[0] >= 'A' && trimmed[0] <= 'Z' && len(words) <= 3 {
		return true
	}
	return false
}

type pageContent struct {
	body string
}

type Model struct {
	selectedIndex   int
	menuItems       []MenuItem
	pageContents    []pageContent
	windowWidth     int
	windowHeight    int
	contentScroll   int
	projects        []Project
	selectedProject int
}

func NewModel() Model {
	notionProjects := fetchProjectsFromNotion()

	projectsContent := ""
	if len(notionProjects) > 0 {
		for i, proj := range notionProjects {
			dateRange := formatDateRange(proj.Date)
			if dateRange != "" {
				projectsContent += fmt.Sprintf("%s (%s)\n", proj.Name, dateRange)
			} else {
				projectsContent += fmt.Sprintf("%s\n", proj.Name)
			}
			projectsContent += fmt.Sprintf("%s\n", proj.Description)
			if proj.TechStack != "" {
				projectsContent += fmt.Sprintf("%s\n", proj.TechStack)
			}
			if i < len(notionProjects)-1 {
				projectsContent += "\n"
			}
		}
	} else {
		projectsContent = "No projects found.\nMake sure NOTION_API_KEY is set."
	}

	return Model{
		selectedIndex:   0,
		selectedProject: 0,
		projects:        notionProjects,
		menuItems: []MenuItem{
			{title: "Home"},
			{title: "About"},
			{title: "Projects"},
			{title: "Contact"},
		},
		pageContents: []pageContent{
			{
				body: ` █████╗ ███╗   ███╗ █████╗ ███╗  ██╗
██╔══██╗████╗ ████║██╔══██╗████╗ ██║
███████║██╔████╔██║███████║██╔██╗██║
██╔══██║██║╚██╔╝██║██╔══██║██║╚████║
██║  ██║██║ ╚═╝ ██║██║  ██║██║ ╚███║
╚═╝  ╚═╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚══╝

 Developer & Student  ·  Dubai, UAE

 ┌─────────────────────────────────────┐
 │  Currently                          │
 │  ▸ Co-Founder @ Falak.me            │
 │  ▸ Senior year, High School         │
 │  ▸ Building things                  │
 └─────────────────────────────────────┘

 Navigate
 ▸ About     — my story, skills & education
 ▸ Projects  — things I've shipped
 ▸ Contact   — get in touch

 Use 1 / 2 / 3 / 4 to navigate`,
			},
			{
				body: `Hey, I'm Aman — a developer and student who enjoys
building things for the web. I'm currently studying and
working on personal projects that keep me learning.

I care about writing clean, readable code and I'm always
looking to pick up new skills.

Experience
▸ Co-Founder – Falak.me
  Jan 2025 – Present | Dubai, UAE

Skills
Python • JavaScript • React.js • Node.js
Next.js • Supabase • Responsive Web Design
Machine Learning • Data Science • AI
Database Development • Data Visualization
Brand Management • Ethical AI Governance

Education
▸ GEMS Our Own Indian School (Apr 2024–Mar 2026)
  Senior School Certificate (CBSE)
  Physics, Chemistry, Math, CS, English
  Grade: 83.4% (Predicted) | Dubai
  Innovation and Coding Team

▸ Bhavans Pearl Wisdom School (Jan 2023–Mar 2024)
  Secondary School Certificate (CBSE)
  Science, Math, Social Science, French, English
  Grade: 94.2% | Al Ain

Certifications
▸ IS-66: Space Weather Events (FEMA) — Feb 2026
▸ Work Smarter with AI (Canva) — Feb 2026
▸ Canva for Work (Canva) — Feb 2026
▸ High-Performance Leadership: F1 (Santander) — Aug 2025
▸ Data Science and AI (IIT Madras) — Jan 2025
▸ Python (Basic) (HackerRank) — Jun 2025
▸ Node (Basic) (HackerRank) — Aug 2023
▸ JavaScript (Basic) (HackerRank) — Aug 2023

Languages
English   ████████████  Bilingual
Malayalam ████████████  Native
Hindi     █████████░░░  Proficient
French    ████░░░░░░░░  Basic`,
			},
			{body: projectsContent},
			{
				body: `✉  aman@falak.me
⚙  github.com/amansanoj
💼  linkedin.com/in/amansanoj

Certifications
▸ IS-66: Space Weather Events (FEMA) — Feb 2026
▸ Work Smarter with AI (Canva) — Feb 2026
▸ Data Science and AI (IIT Madras) — Jan 2025
▸ Python (Basic) (HackerRank) — Jun 2025
▸ Node (Basic) (HackerRank) — Aug 2023
▸ JavaScript (Basic) (HackerRank) — Aug 2023`,
			},
		},
		windowWidth:  120,
		windowHeight: 30,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			if m.selectedIndex == 2 && len(m.projects) > m.selectedProject {
				project := m.projects[m.selectedProject]
				if project.URL != "" {
					openURL(project.URL)
				}
			}

		case "1":
			m.selectedIndex = 0
			m.contentScroll = 0
			m.selectedProject = 0
		case "2":
			m.selectedIndex = 1
			m.contentScroll = 0
			m.selectedProject = 0
		case "3":
			m.selectedIndex = 2
			m.contentScroll = 0
			m.selectedProject = 0
		case "4":
			m.selectedIndex = 3
			m.contentScroll = 0
			m.selectedProject = 0

		case "up":
			if m.selectedIndex == 2 && len(m.projects) > 0 {
				if m.selectedProject > 0 {
					m.selectedProject--
				}
			} else {
				if m.contentScroll > 0 {
					m.contentScroll--
				}
			}

		case "down":
			if m.selectedIndex == 2 && len(m.projects) > 0 {
				if m.selectedProject < len(m.projects)-1 {
					m.selectedProject++
				}
			} else {
				maxScroll := m.getMaxContentScroll()
				if m.contentScroll < maxScroll {
					m.contentScroll++
				}
			}

		case "pgdn":
			maxScroll := m.getMaxContentScroll()
			m.contentScroll += 5
			if m.contentScroll > maxScroll {
				m.contentScroll = maxScroll
			}

		case "pgup":
			m.contentScroll -= 5
			if m.contentScroll < 0 {
				m.contentScroll = 0
			}
		}

	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	}
	return m, nil
}

func (m Model) availableContentHeight() int {
	paneHeight := m.windowHeight - 3
	inner := paneHeight - 2 - 2 - 1
	if inner < 1 {
		inner = 1
	}
	return inner
}

func (m Model) getMaxContentScroll() int {
	if m.selectedIndex < 0 || m.selectedIndex >= len(m.pageContents) {
		return 0
	}
	contentLines := strings.Split(m.pageContents[m.selectedIndex].body, "\n")
	avail := m.availableContentHeight()
	maxScroll := len(contentLines) - avail
	if maxScroll < 0 {
		maxScroll = 0
	}
	return maxScroll
}

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

	contentLines := strings.Split(page.body, "\n")
	avail := m.availableContentHeight()

	var visibleLines []string
	for i := m.contentScroll; i < len(contentLines) && i < m.contentScroll+avail; i++ {
		visibleLines = append(visibleLines, contentLines[i])
	}
	for len(visibleLines) < avail {
		visibleLines = append(visibleLines, "")
	}

	var visibleContent string
	switch m.selectedIndex {
	case 1: // About
		visibleContent = m.buildStyledAboutContent(visibleLines, mainWidth)
	case 2: // Projects
		if len(m.projects) > 0 {
			visibleContent = m.buildStyledProjectContent(visibleLines, m.selectedProject, mainWidth)
		} else {
			visibleContent = contentStyle.Render(strings.Join(visibleLines, "\n"))
		}
	default:
		visibleContent = contentStyle.Render(strings.Join(visibleLines, "\n"))
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

		// Empty lines
		if trimmed == "" {
			result += "\n"
			continue
		}

		// Section headers: "Experience", "Skills", "Education", etc.
		if isSectionHeader(line) {
			result += sectionHeaderStyle.Render(trimmed) + "\n"
			continue
		}

		// Bullet lines starting with ▸
		if strings.HasPrefix(trimmed, "▸") {
			rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "▸"))
			// If line contains a date in parens, split and grey the date
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

		// Indented sub-lines (start with spaces in original)
		if strings.HasPrefix(line, "  ") {
		// Date range lines → grey
		if strings.Contains(trimmed, "–") && strings.Contains(trimmed, "(") {
			result += mutedStyle.Render("  "+trimmed) + "\n"
		// Extracurricular/team lines → orange
		} else if strings.Contains(trimmed, "Team") || strings.Contains(trimmed, "Cadet") {
			result += "  " + activeItemStyle.Render(trimmed) + "\n"
		// Everything else → white
		} else {
			result += "  " + contentStyle.Render(trimmed) + "\n"
		}
		continue
	}

		// Language bars (contain █ or ░)
		if strings.ContainsAny(line, "█░") {
			// Split into name, bar, label
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				name := parts[0]
				bar := parts[1]
				label := strings.Join(parts[2:], " ")
				result += contentStyle.Render(fmt.Sprintf("%-10s", name)) +
					activeItemStyle.Render(bar) +
					mutedStyle.Render("  "+label) + "\n"
				continue
			}
		}

		// Default: wrap long lines
		wrapped := wordWrap(trimmed, wrapWidth)
		for _, wline := range wrapped {
			result += contentStyle.Render(wline) + "\n"
		}
	}

	return strings.TrimRight(result, "\n")
}

func (m Model) buildStyledProjectContent(lines []string, selectedProject int, mainWidth int) string {
	var result string
	projectIndex := 0

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

		// Project title lines: "Name (Date – Date)"
		if strings.Contains(line, "(") && strings.Contains(line, "–") && strings.Contains(line, ")") {
			parts := strings.SplitN(line, " (", 2)
			if len(parts) == 2 {
				var indicator string
				if projectIndex == selectedProject {
					indicator = activeItemStyle.Render("● ")
				} else {
					indicator = dimStyle.Render("○ ")
				}
				date := strings.TrimSuffix(parts[1], ")")
				result += indicator + projectTitleStyle.Render(parts[0]) + "\n"
				result += "  " + mutedStyle.Render(date) + "\n"
				projectIndex++
				continue
			}
		}

		// Tech stack lines: comma-separated, no spaces after commas
		if strings.Contains(line, ",") && !strings.Contains(line, ", ") {
			tags := strings.Split(line, ",")
			var styled []string
			for _, tag := range tags {
				styled = append(styled, activeItemStyle.Render(strings.TrimSpace(tag)))
			}
			result += "  " + strings.Join(styled, dimStyle.Render(" · ")) + "\n"
			continue
		}

		// Description: word-wrap
		wrapped := wordWrap(trimmed, wrapWidth)
		for _, wline := range wrapped {
			result += contentStyle.Render("  "+wline) + "\n"
		}
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
	if m.selectedIndex == 2 && len(m.projects) > 0 {
		left = statusBarStyle.Render(" ") +
			hint("1-4", "navigate") + sep +
			hint("↑/↓", "select") + sep +
			hint("enter", "open") + sep +
			hint("q", "quit")
	} else {
		left = statusBarStyle.Render(" ") +
			hint("1-4", "navigate") + sep +
			hint("↑/↓", "scroll") + sep +
			hint("q", "quit")
	}

	maxScroll := m.getMaxContentScroll()
	var right string
	if maxScroll > 0 {
		pct := int((float32(m.contentScroll) / float32(maxScroll)) * 100)
		filled := int((float32(m.contentScroll) / float32(maxScroll)) * 10)
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
	if gap < 1 {
		gap = 1
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

func main() {
	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}