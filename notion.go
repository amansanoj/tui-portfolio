package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	notionAPIVersion               = "2022-06-28"
	notionRequestTimeout           = 10 * time.Second
	notionMaxRetries               = 3
	notionRetryBaseDelay           = 250 * time.Millisecond
	notionUserAgent                = "tui-portfolio/1.0"
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

	var lastErr error
	for attempt := 1; attempt <= notionMaxRetries; attempt++ {
		status, body, retryAfter, err := executeNotionQuery(ctx, url, apiKey)
		if err != nil {
			lastErr = err
			if ctx.Err() != nil {
				break
			}
			if attempt < notionMaxRetries {
				if !sleepWithContext(ctx, notionRetryBaseDelay*time.Duration(1<<(attempt-1))) {
					break
				}
				continue
			}
			break
		}

		if status == http.StatusOK {
			if err := json.Unmarshal(body, &notionResp); err != nil {
				return notionResp, fmt.Errorf("parsing response failed: %w", err)
			}
			return notionResp, nil
		}

		msg := strings.TrimSpace(string(body))
		lastErr = fmt.Errorf("notion API returned %d: %s", status, msg)
		if !isRetryableStatus(status) || attempt == notionMaxRetries {
			break
		}

		delay := parseRetryAfter(retryAfter)
		if delay <= 0 {
			delay = notionRetryBaseDelay * time.Duration(1<<(attempt-1))
		}
		if !sleepWithContext(ctx, delay) {
			break
		}
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("unknown Notion request failure")
	}
	return notionResp, lastErr
}

func executeNotionQuery(ctx context.Context, url, apiKey string) (int, []byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(`{}`))
	if err != nil {
		return 0, nil, "", fmt.Errorf("creating request failed: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Notion-Version", notionAPIVersion)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", notionUserAgent)

	resp, err := notionHTTPClient.Do(req)
	if err != nil {
		return 0, nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, "", fmt.Errorf("reading response failed: %w", err)
	}

	return resp.StatusCode, body, resp.Header.Get("Retry-After"), nil
}

func isRetryableStatus(status int) bool {
	if status == http.StatusTooManyRequests {
		return true
	}
	return status == http.StatusInternalServerError ||
		status == http.StatusBadGateway ||
		status == http.StatusServiceUnavailable ||
		status == http.StatusGatewayTimeout
}

func parseRetryAfter(header string) time.Duration {
	header = strings.TrimSpace(header)
	if header == "" {
		return 0
	}
	seconds, err := strconv.Atoi(header)
	if err != nil || seconds <= 0 {
		return 0
	}
	return time.Duration(seconds) * time.Second
}

func sleepWithContext(ctx context.Context, delay time.Duration) bool {
	t := time.NewTimer(delay)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
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
