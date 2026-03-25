package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	notionAPIVersion               = "2022-06-28"
	notionRequestTimeout           = 10 * time.Second
	defaultNotionProjectsDatabase  = "32acb49d4dc9804ab1b5f3ccf42c375c"
	defaultNotionCertsDatabase     = "32bcb49d4dc9806e82aae4f172dbf8cd"
	notionProjectsDatabaseIDEnvVar = "NOTION_PROJECTS_DB_ID"
	notionCertsDatabaseIDEnvVar    = "NOTION_CERTS_DB_ID"
)

var notionHTTPClient = &http.Client{Timeout: notionRequestTimeout}

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
	apiKey := strings.TrimSpace(os.Getenv("NOTION_API_KEY"))
	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "Error: NOTION_API_KEY environment variable not set\n")
		return []Project{}
	}

	databaseID := notionDatabaseID(notionProjectsDatabaseIDEnvVar, defaultNotionProjectsDatabase)
	notionResp, err := queryNotionDatabase(apiKey, databaseID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "notion projects request failed: %v\n", err)
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

func notionDatabaseID(envName, fallback string) string {
	value := strings.TrimSpace(os.Getenv(envName))
	if value == "" {
		return fallback
	}
	return value
}

func queryNotionDatabase(apiKey, databaseID string) (NotionResponse, error) {
	var notionResp NotionResponse

	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", databaseID)
	ctx, cancel := context.WithTimeout(context.Background(), notionRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(`{}`))
	if err != nil {
		return notionResp, fmt.Errorf("creating request failed: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Notion-Version", notionAPIVersion)
	req.Header.Add("Content-Type", "application/json")

	resp, err := notionHTTPClient.Do(req)
	if err != nil {
		return notionResp, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return notionResp, fmt.Errorf("reading response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return notionResp, fmt.Errorf("notion API returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if err := json.Unmarshal(body, &notionResp); err != nil {
		return notionResp, fmt.Errorf("parsing response failed: %w", err)
	}

	return notionResp, nil
}

func extractPlainText(arr []interface{}) string {
	if len(arr) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, v := range arr {
		part, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		text, ok := part["plain_text"].(string)
		if !ok {
			continue
		}
		sb.WriteString(text)
	}

	return strings.TrimSpace(sb.String())
}

func extractStringProperty(props map[string]interface{}, propName string) string {
	if prop, exists := props[propName]; exists {
		propMap, ok := prop.(map[string]interface{})
		if !ok {
			return ""
		}

		propType, ok := propMap["type"].(string)
		if !ok {
			return ""
		}

		switch propType {
		case "title":
			if titleArr, ok := propMap["title"].([]interface{}); ok {
				return extractPlainText(titleArr)
			}
		case "rich_text":
			if richArr, ok := propMap["rich_text"].([]interface{}); ok {
				return extractPlainText(richArr)
			}
		}
	}
	return ""
}

func extractDateProperty(props map[string]interface{}, propName string) string {
	if prop, exists := props[propName]; exists {
		propMap, ok := prop.(map[string]interface{})
		if !ok {
			return ""
		}
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
		propMap, ok := prop.(map[string]interface{})
		if !ok {
			return ""
		}
		if url, ok := propMap["url"].(string); ok {
			return url
		}
	}
	return ""
}

type Certification struct {
	Title        string
	Date         string
	Organization string
	URL          string
}

func fetchCertificationsFromNotion() []Certification {
	apiKey := strings.TrimSpace(os.Getenv("NOTION_API_KEY"))
	if apiKey == "" {
		return []Certification{}
	}

	databaseID := notionDatabaseID(notionCertsDatabaseIDEnvVar, defaultNotionCertsDatabase)
	notionResp, err := queryNotionDatabase(apiKey, databaseID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "notion certs request failed: %v\n", err)
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
