package pages

import "fmt"

type ProjectData struct {
	Name        string
	Description string
	Date        string
	TechStack   string
}

func BuildProjectsBody(projects []ProjectData) (string, []int, []int, int) {
	if len(projects) == 0 {
		return "No projects found.\nMake sure NOTION_API_KEY is set.", nil, nil, 2
	}

	bodyOffsets := make([]int, len(projects))
	renderedOffsets := make([]int, len(projects))
	bodyLine := 0
	renderedLine := 0
	body := ""

	for i, proj := range projects {
		bodyOffsets[i] = bodyLine
		renderedOffsets[i] = renderedLine

		if proj.Date != "" {
			body += fmt.Sprintf("PROJ|||%s|||%s\n", proj.Name, proj.Date)
		} else {
			body += fmt.Sprintf("PROJ|||%s|||\n", proj.Name)
		}
		bodyLine++
		renderedLine++

		body += fmt.Sprintf("%s\n", proj.Description)
		bodyLine++
		renderedLine++

		if proj.TechStack != "" {
			body += fmt.Sprintf("%s\n", proj.TechStack)
			bodyLine++
			renderedLine++
		}
		if i < len(projects)-1 {
			body += "\n"
			bodyLine++
			renderedLine++
		}
	}

	return body, bodyOffsets, renderedOffsets, renderedLine
}
