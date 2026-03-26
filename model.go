package main

import (
	"fmt"
	"strings"

	"github.com/amansanoj/tui-portfolio/pages"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuItem struct {
	title string
}

type pageContent struct {
	body string
}

type ContactItem struct {
	Label  string
	Handle string
	URL    string
}

type Model struct {
	styles                 Styles
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
	contactItems           []ContactItem
	selectedContact        int
	showingURL             string // non-empty = show URL popup
}

var defaultContactItems = []ContactItem{
	{Label: "Email", Handle: "aman@falak.me", URL: "mailto:aman@falak.me"},
	{Label: "GitHub", Handle: "amansanoj", URL: "https://github.com/amansanoj"},
	{Label: "LinkedIn", Handle: "amansanoj", URL: "https://linkedin.com/in/amansanoj"},
}

func toPageProjects(projectsData []Project) []pages.ProjectData {
	out := make([]pages.ProjectData, 0, len(projectsData))
	for _, proj := range projectsData {
		out = append(out, pages.ProjectData{
			Name:        proj.Name,
			Description: proj.Description,
			Date:        formatDateRange(proj.Date),
			TechStack:   proj.TechStack,
		})
	}
	return out
}

func toPageCerts(certsData []Certification) []pages.CertData {
	out := make([]pages.CertData, 0, len(certsData))
	for _, cert := range certsData {
		out = append(out, pages.CertData{
			Title:        cert.Title,
			Date:         formatDateRange(cert.Date),
			Organization: cert.Organization,
			URL:          cert.URL,
		})
	}
	return out
}

func toPageContacts(items []ContactItem) []pages.ContactData {
	out := make([]pages.ContactData, 0, len(items))
	for _, item := range items {
		out = append(out, pages.ContactData{
			Label:  item.Label,
			Handle: item.Handle,
			URL:    item.URL,
		})
	}
	return out
}

func NewModel(renderer *lipgloss.Renderer) Model {
	snapshot := appContentStore.Snapshot()
	notionProjects := snapshot.Projects
	notionCerts := snapshot.Certifications
	pageProjects := toPageProjects(notionProjects)
	pageCerts := toPageCerts(notionCerts)

	var projectsContent string
	var projBodyOffsets, projRenderedOffsets []int
	var projRenderedTotal int
	if snapshot.Ready {
		projectsContent, projBodyOffsets, projRenderedOffsets, projRenderedTotal = pages.BuildProjectsBody(pageProjects)
	} else {
		projectsContent = "Loading projects..."
		projRenderedTotal = 1
		go appContentStore.Refresh()
	}

	var certsContent string
	var certBodyOffsets, certRenderedOffsets []int
	var certRenderedTotal int
	if snapshot.Ready {
		certsContent, certBodyOffsets, certRenderedOffsets, certRenderedTotal = pages.BuildCertsBody(pageCerts)
	} else {
		certsContent = "Loading certifications..."
		certRenderedTotal = 1
	}

	contactItems := defaultContactItems
	contactBody := pages.BuildContactBody(toPageContacts(contactItems))

	return Model{
		styles:                 makeStyles(renderer),
		selectedIndex:          0,
		selectedProject:        0,
		selectedCert:           0,
		selectedContact:        0,
		projects:               notionProjects,
		projectBodyOffsets:     projBodyOffsets,
		projectRenderedOffsets: projRenderedOffsets,
		projectRenderedTotal:   projRenderedTotal,
		certifications:         notionCerts,
		certBodyOffsets:        certBodyOffsets,
		certRenderedOffsets:    certRenderedOffsets,
		certRenderedTotal:      certRenderedTotal,
		contactItems:           contactItems,
		menuItems: []MenuItem{
			{title: "Home"},
			{title: "About"},
			{title: "Projects"},
			{title: "Certs"},
			{title: "Contact"},
			{title: "Status"},
		},
		pageContents: []pageContent{
			{
				body: pages.HomeBody(),
			},
			{
				body: pages.AboutBody(),
			},
			{body: projectsContent},
			{body: certsContent},
			{body: contactBody},
			{body: ""},
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
		// Any key dismisses the popup
		if m.showingURL != "" {
			m.showingURL = ""
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			url := ""
			switch m.selectedIndex {
			case 2:
				if len(m.projects) > m.selectedProject {
					url = m.projects[m.selectedProject].URL
				}
			case 3:
				if len(m.certifications) > m.selectedCert {
					url = m.certifications[m.selectedCert].URL
				}
			case 4:
				if len(m.contactItems) > m.selectedContact {
					url = m.contactItems[m.selectedContact].URL
				}
			}
			if url != "" {
				m.showingURL = url
			}

		case "1":
			m.selectedIndex = 0
			m.contentScroll = 0
			m.selectedProject = 0
			m.selectedCert = 0
			m.selectedContact = 0
		case "2":
			m.selectedIndex = 1
			m.contentScroll = 0
			m.selectedProject = 0
			m.selectedCert = 0
			m.selectedContact = 0
		case "3":
			m.selectedIndex = 2
			m.contentScroll = 0
			m.selectedProject = 0
			m.selectedCert = 0
			m.selectedContact = 0
		case "4":
			m.selectedIndex = 3
			m.contentScroll = 0
			m.selectedProject = 0
			m.selectedCert = 0
			m.selectedContact = 0
		case "5":
			m.selectedIndex = 4
			m.contentScroll = 0
			m.selectedProject = 0
			m.selectedCert = 0
			m.selectedContact = 0
		case "6":
			m.selectedIndex = 5
			m.contentScroll = 0
			m.selectedProject = 0
			m.selectedCert = 0
			m.selectedContact = 0

		case "up":
			switch m.selectedIndex {
			case 2:
				if len(m.projects) > 0 && m.selectedProject > 0 {
					m.selectedProject--
					projectOffsets, _ := m.projectRenderedMetrics()
					m.contentScroll = m.scrollRenderedToShow(projectOffsets[m.selectedProject], 0)
				}
			case 3:
				if len(m.certifications) > 0 && m.selectedCert > 0 {
					m.selectedCert--
					m.contentScroll = m.scrollRenderedToShow(m.certRenderedOffsets[m.selectedCert], 0)
				}
			case 4:
				if m.selectedContact > 0 {
					m.selectedContact--
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
					projectOffsets, _ := m.projectRenderedMetrics()
					m.contentScroll = m.scrollRenderedToShow(
						projectOffsets[m.selectedProject], m.projectItemHeight(m.selectedProject))
				}
			case 3:
				if len(m.certifications) > 0 && m.selectedCert < len(m.certifications)-1 {
					m.selectedCert++
					m.contentScroll = m.scrollRenderedToShow(
						m.certRenderedOffsets[m.selectedCert], m.certItemHeight(m.selectedCert))
				}
			case 4:
				if m.selectedContact < len(m.contactItems)-1 {
					m.selectedContact++
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
				projectOffsets, _ := m.projectRenderedMetrics()
				m.contentScroll = m.scrollRenderedToShow(
					projectOffsets[m.selectedProject], m.projectItemHeight(m.selectedProject))
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
				projectOffsets, _ := m.projectRenderedMetrics()
				m.contentScroll = m.scrollRenderedToShow(projectOffsets[m.selectedProject], 0)
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

func (m Model) certItemHeight(_ int) int {
	return 3
}

func (m Model) projectItemHeight(idx int) int {
	if idx < 0 || idx >= len(m.projects) {
		return 0
	}

	wrapWidth := m.projectWrapWidth()
	h := 1

	desc := strings.TrimSpace(m.projects[idx].Description)
	if desc == "" {
		h++
	} else {
		h += len(wordWrap(desc, wrapWidth))
	}

	if m.projects[idx].TechStack != "" {
		if strings.Contains(m.projects[idx].TechStack, ",") && !strings.Contains(m.projects[idx].TechStack, ", ") {
			h++
		} else {
			tech := strings.TrimSpace(m.projects[idx].TechStack)
			if tech == "" {
				h++
			} else {
				h += len(wordWrap(tech, wrapWidth))
			}
		}
	}

	return h
}

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
	// Inner pane height minus: border(2) + padding(2) + title block spacer(1 line).
	inner := paneHeight - 2 - 2 - 1
	if inner < 1 {
		inner = 1
	}
	return inner
}

func (m Model) getMaxContentScroll() int {
	avail := m.availableContentHeight()
	switch m.selectedIndex {
	case 1:
		mainWidth := m.windowWidth - 20 - 4
		lines := m.aboutPageLineCount(mainWidth)
		max := lines - avail
		if max < 0 {
			max = 0
		}
		return max
	case 2:
		_, total := m.projectRenderedMetrics()
		max := total - avail
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
	case 5:
		mainWidth := m.windowWidth - 20 - 4
		lines := m.statusPageLineCount(mainWidth)
		max := lines - avail
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

func (m Model) aboutPageLineCount(mainWidth int) int {
	if len(m.pageContents) <= 1 {
		return 0
	}

	allLines := strings.Split(m.pageContents[1].body, "\n")
	plain := func(s string) string { return s }
	theme := pages.Theme{
		SectionHeader: plain,
		Bullet:        plain,
		Content:       plain,
		Muted:         plain,
		Active:        plain,
		Dim:           plain,
		ProjectTitle:  plain,
	}

	rendered := pages.RenderAbout(allLines, mainWidth, theme)
	if rendered == "" {
		return 0
	}

	return len(strings.Split(rendered, "\n"))
}

func (m Model) projectWrapWidth() int {
	mainWidth := m.windowWidth - 20 - 4
	wrapWidth := mainWidth - 2 - 2 - 2
	if wrapWidth < 40 {
		wrapWidth = 40
	}
	return wrapWidth
}

func (m Model) projectRenderedMetrics() ([]int, int) {
	if len(m.projects) == 0 {
		return nil, 0
	}

	wrapWidth := m.projectWrapWidth()
	offsets := make([]int, len(m.projects))
	renderedLine := 0

	for i, proj := range m.projects {
		offsets[i] = renderedLine

		// Title and date are rendered together on one line.
		renderedLine++

		desc := strings.TrimSpace(proj.Description)
		if desc == "" {
			renderedLine++
		} else {
			renderedLine += len(wordWrap(desc, wrapWidth))
		}

		if proj.TechStack != "" {
			if strings.Contains(proj.TechStack, ",") && !strings.Contains(proj.TechStack, ", ") {
				renderedLine++
			} else {
				tech := strings.TrimSpace(proj.TechStack)
				if tech == "" {
					renderedLine++
				} else {
					renderedLine += len(wordWrap(tech, wrapWidth))
				}
			}
		}

		if i < len(m.projects)-1 {
			renderedLine++
		}
	}

	return offsets, renderedLine
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
