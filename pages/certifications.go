package pages

import "fmt"

type CertData struct {
	Title        string
	Date         string
	Organization string
	URL          string
}

func BuildCertsBody(certs []CertData) (string, []int, []int, int) {
	if len(certs) == 0 {
		return "No certifications found.\nMake sure NOTION_API_KEY is set.", nil, nil, 2
	}

	bodyOffsets := make([]int, len(certs))
	renderedOffsets := make([]int, len(certs))
	bodyLine := 0
	renderedLine := 0
	body := ""

	for i, cert := range certs {
		bodyOffsets[i] = bodyLine
		renderedOffsets[i] = renderedLine

		if cert.Date != "" {
			body += fmt.Sprintf("CERT|||%s|||%s\n", cert.Title, cert.Date)
		} else {
			body += fmt.Sprintf("CERT|||%s|||\n", cert.Title)
		}
		bodyLine++
		renderedLine += 2

		body += fmt.Sprintf("ORG|||%s\n", cert.Organization)
		bodyLine++
		renderedLine++

		if i < len(certs)-1 {
			body += "\n"
			bodyLine++
			renderedLine++
		}
	}

	return body, bodyOffsets, renderedOffsets, renderedLine
}
