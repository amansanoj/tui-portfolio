package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type MenuItem struct {
	title string
}

type pageContent struct {
	body string
}

type Model struct {
	selectedIndex          int
	menuItems              []MenuItem
	pageContents           []pageContent
	windowWidth            int
	windowHeight           int
	contentScroll          int
	projects               []Project
	selectedProject        int
	projectBodyOffsets     []int
	projectRenderedOffsets []int
	projectRenderedTotal   int
	certifications         []Certification
	selectedCert           int
	certBodyOffsets        []int
	certRenderedOffsets    []int
	certRenderedTotal      int
}

func NewModel() Model {
	notionProjects := fetchProjectsFromNotion()
	notionCerts := fetchCertificationsFromNotion()

	projectsContent, projBodyOffsets, projRenderedOffsets, projRenderedTotal := buildProjectsBody(notionProjects)
	certsContent, certBodyOffsets, certRenderedOffsets, certRenderedTotal := buildCertsBody(notionCerts)

	return Model{
		selectedIndex:          0,
		selectedProject:        0,
		selectedCert:           0,
		projects:               notionProjects,
		projectBodyOffsets:     projBodyOffsets,
		projectRenderedOffsets: projRenderedOffsets,
		projectRenderedTotal:   projRenderedTotal,
		certifications:         notionCerts,
		certBodyOffsets:        certBodyOffsets,
		certRenderedOffsets:    certRenderedOffsets,
		certRenderedTotal:      certRenderedTotal,
		menuItems: []MenuItem{
			{title: "Home"},
			{title: "About"},
			{title: "Projects"},
			{title: "Certs"},
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
 ▸ Certs     — certifications I've earned
 ▸ Contact   — get in touch

 Use 1 / 2 / 3 / 4 / 5 to navigate`,
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

Languages
English   ████████████  Bilingual
Malayalam ████████████  Native
Hindi     █████████░░░  Proficient
French    ████░░░░░░░░  Basic`,
			},
			{body: projectsContent},
			{body: certsContent},
			{
				body: `✉  aman@falak.me
⚙  github.com/amansanoj
💼  linkedin.com/in/amansanoj`,
			},
		},
		windowWidth:  120,
		windowHeight: 30,
	}
}

func buildProjectsBody(projects []Project) (string, []int, []int, int) {
	if len(projects) == 0 {
		return "No projects found.\nMake sure NOTION_API_KEY is set.", nil, nil, 2
	}
	var sb strings.Builder
	bodyOffsets := make([]int, len(projects))
	renderedOffsets := make([]int, len(projects))
	bodyLine := 0
	renderedLine := 0
	for i, proj := range projects {
		bodyOffsets[i] = bodyLine
		renderedOffsets[i] = renderedLine

		dateRange := formatDateRange(proj.Date)
		if dateRange != "" {
			sb.WriteString(fmt.Sprintf("%s (%s)\n", proj.Name, dateRange))
		} else {
			sb.WriteString(fmt.Sprintf("%s\n", proj.Name))
		}
		bodyLine++
		renderedLine += 2 // title + date

		sb.WriteString(fmt.Sprintf("%s\n", proj.Description))
		bodyLine++
		renderedLine++

		if proj.TechStack != "" {
			sb.WriteString(fmt.Sprintf("%s\n", proj.TechStack))
			bodyLine++
			renderedLine++
		}
		if i < len(projects)-1 {
			sb.WriteString("\n")
			bodyLine++
			renderedLine++
		}
	}
	return sb.String(), bodyOffsets, renderedOffsets, renderedLine
}

func buildCertsBody(certs []Certification) (string, []int, []int, int) {
	if len(certs) == 0 {
		return "No certifications found.\nMake sure NOTION_API_KEY is set.", nil, nil, 2
	}
	var sb strings.Builder
	bodyOffsets := make([]int, len(certs))
	renderedOffsets := make([]int, len(certs))
	bodyLine := 0
	renderedLine := 0
	for i, cert := range certs {
		bodyOffsets[i] = bodyLine
		renderedOffsets[i] = renderedLine

		date := formatDateRange(cert.Date)
		if date != "" {
			sb.WriteString(fmt.Sprintf("CERT|||%s|||%s\n", cert.Title, date))
		} else {
			sb.WriteString(fmt.Sprintf("CERT|||%s|||\n", cert.Title))
		}
		bodyLine++
		renderedLine += 2 // title + date

		sb.WriteString(fmt.Sprintf("ORG|||%s\n", cert.Organization))
		bodyLine++
		renderedLine++

		// No URL line in body — ↗ icon on title handles that visually

		if i < len(certs)-1 {
			sb.WriteString("\n")
			bodyLine++
			renderedLine++
		}
	}
	return sb.String(), bodyOffsets, renderedOffsets, renderedLine
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
				if url := m.projects[m.selectedProject].URL; url != "" {
					openURL(url)
				}
			}
			if m.selectedIndex == 3 && len(m.certifications) > m.selectedCert {
				if url := m.certifications[m.selectedCert].URL; url != "" {
					openURL(url)
				}
			}

		case "1":
			m.selectedIndex = 0
			m.contentScroll = 0
			m.selectedProject = 0
			m.selectedCert = 0
		case "2":
			m.selectedIndex = 1
			m.contentScroll = 0
			m.selectedProject = 0
			m.selectedCert = 0
		case "3":
			m.selectedIndex = 2
			m.contentScroll = 0
			m.selectedProject = 0
			m.selectedCert = 0
		case "4":
			m.selectedIndex = 3
			m.contentScroll = 0
			m.selectedProject = 0
			m.selectedCert = 0
		case "5":
			m.selectedIndex = 4
			m.contentScroll = 0
			m.selectedProject = 0
			m.selectedCert = 0

		case "up":
			switch m.selectedIndex {
			case 2:
				if len(m.projects) > 0 && m.selectedProject > 0 {
					m.selectedProject--
					m.contentScroll = m.scrollRenderedToShow(m.projectRenderedOffsets[m.selectedProject], 0)
				}
			case 3:
				if len(m.certifications) > 0 && m.selectedCert > 0 {
					m.selectedCert--
					m.contentScroll = m.scrollRenderedToShow(m.certRenderedOffsets[m.selectedCert], 0)
				}
			default:
				if m.contentScroll > 0 {
					m.contentScroll--
				}
			}

		case "down":
			switch m.selectedIndex {
			case 2:
				if len(m.projects) > 0 && m.selectedProject < len(m.projects)-1 {
					m.selectedProject++
					m.contentScroll = m.scrollRenderedToShow(
						m.projectRenderedOffsets[m.selectedProject], m.projectItemHeight(m.selectedProject))
				}
			case 3:
				if len(m.certifications) > 0 && m.selectedCert < len(m.certifications)-1 {
					m.selectedCert++
					m.contentScroll = m.scrollRenderedToShow(
						m.certRenderedOffsets[m.selectedCert], m.certItemHeight(m.selectedCert))
				}
			default:
				maxScroll := m.getMaxContentScroll()
				if m.contentScroll < maxScroll {
					m.contentScroll++
				}
			}

		case "pgdn":
			switch m.selectedIndex {
			case 2:
				for i := 0; i < 5 && m.selectedProject < len(m.projects)-1; i++ {
					m.selectedProject++
				}
				m.contentScroll = m.scrollRenderedToShow(
					m.projectRenderedOffsets[m.selectedProject], m.projectItemHeight(m.selectedProject))
			case 3:
				for i := 0; i < 5 && m.selectedCert < len(m.certifications)-1; i++ {
					m.selectedCert++
				}
				m.contentScroll = m.scrollRenderedToShow(
					m.certRenderedOffsets[m.selectedCert], m.certItemHeight(m.selectedCert))
			default:
				maxScroll := m.getMaxContentScroll()
				m.contentScroll += 5
				if m.contentScroll > maxScroll {
					m.contentScroll = maxScroll
				}
			}

		case "pgup":
			switch m.selectedIndex {
			case 2:
				for i := 0; i < 5 && m.selectedProject > 0; i++ {
					m.selectedProject--
				}
				m.contentScroll = m.scrollRenderedToShow(m.projectRenderedOffsets[m.selectedProject], 0)
			case 3:
				for i := 0; i < 5 && m.selectedCert > 0; i++ {
					m.selectedCert--
				}
				m.contentScroll = m.scrollRenderedToShow(m.certRenderedOffsets[m.selectedCert], 0)
			default:
				m.contentScroll -= 5
				if m.contentScroll < 0 {
					m.contentScroll = 0
				}
			}
		}

	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	}
	return m, nil
}

func (m Model) certItemHeight(idx int) int {
	// title(1) + date(1) + org(1) = 3, no URL line in body anymore
	return 3
}

func (m Model) projectItemHeight(idx int) int {
	h := 3 // title(rendered as 2) + desc = 3
	if idx >= 0 && idx < len(m.projects) && m.projects[idx].TechStack != "" {
		h++
	}
	return h
}

// scrollRenderedToShow keeps contentScroll in rendered-line space.
// itemHeight=0 → only ensure top of item is visible (used when going up).
// itemHeight>0 → ensure full item block is visible (used when going down).
func (m Model) scrollRenderedToShow(renderedStart, itemHeight int) int {
	avail := m.availableContentHeight()
	if itemHeight == 0 {
		if renderedStart < m.contentScroll {
			return renderedStart
		}
		return m.contentScroll
	}
	bottom := renderedStart + itemHeight - 1
	if bottom >= m.contentScroll+avail {
		newScroll := bottom - avail + 1
		if newScroll < 0 {
			newScroll = 0
		}
		return newScroll
	}
	return m.contentScroll
}

func (m Model) availableContentHeight() int {
	paneHeight := m.windowHeight - 3
	inner := paneHeight - 2 - 2 - 1
	if inner < 1 {
		inner = 1
	}
	return inner
}

// getMaxContentScroll returns max scroll in the correct unit space for each page.
// Pages 2 and 3 use rendered-line space; others use body-line space.
func (m Model) getMaxContentScroll() int {
	avail := m.availableContentHeight()
	switch m.selectedIndex {
	case 2:
		max := m.projectRenderedTotal - avail
		if max < 0 {
			max = 0
		}
		return max
	case 3:
		max := m.certRenderedTotal - avail
		if max < 0 {
			max = 0
		}
		return max
	default:
		if m.selectedIndex < 0 || m.selectedIndex >= len(m.pageContents) {
			return 0
		}
		lines := strings.Split(m.pageContents[m.selectedIndex].body, "\n")
		max := len(lines) - avail
		if max < 0 {
			max = 0
		}
		return max
	}
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
