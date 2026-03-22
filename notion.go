package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const notionDatabaseID = "32acb49d4dc9804ab1b5f3ccf42c375c"
const notionAPIVersion = "2022-06-28"

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

const notionCertsDatabaseID = "32bcb49d4dc9806e82aae4f172dbf8cd"

type Certification struct {
	Title        string
	Date         string
	Organization string
	URL          string
}

func fetchCertificationsFromNotion() []Certification {
	apiKey := os.Getenv("NOTION_API_KEY")
	if apiKey == "" {
		return []Certification{}
	}

	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", notionCertsDatabaseID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		return []Certification{}
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Notion-Version", notionAPIVersion)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []Certification{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Certification{}
	}

	var notionResp NotionResponse
	if err := json.Unmarshal(body, &notionResp); err != nil {
		return []Certification{}
	}

	var certs []Certification

	for _, page := range notionResp.Results {
		cert := Certification{
			Title:        extractStringProperty(page.Properties, "Name"),
			Date:         extractDateProperty(page.Properties, "Validity Period"),
			Organization: extractStringProperty(page.Properties, "Issuing Organization"),
			URL:          extractURLProperty(page.Properties, "Certification URL"),
		}
		if cert.Title != "" {
			certs = append(certs, cert)
		}
	}
	return certs
}
