package main

import (
	"testing"
	"time"
)

func TestExtractPlainText(t *testing.T) {
	arr := []interface{}{
		map[string]interface{}{"plain_text": "Hello"},
		map[string]interface{}{"plain_text": " "},
		map[string]interface{}{"plain_text": "World"},
	}

	got := extractPlainText(arr)
	want := "Hello World"
	if got != want {
		t.Fatalf("extractPlainText() = %q, want %q", got, want)
	}
}

func TestExtractStringProperty(t *testing.T) {
	props := map[string]interface{}{
		"Title": map[string]interface{}{
			"type": "title",
			"title": []interface{}{
				map[string]interface{}{"plain_text": "Aman"},
				map[string]interface{}{"plain_text": " Sanoj"},
			},
		},
		"Desc": map[string]interface{}{
			"type": "rich_text",
			"rich_text": []interface{}{
				map[string]interface{}{"plain_text": "Developer"},
				map[string]interface{}{"plain_text": " and Student"},
			},
		},
	}

	if got, want := extractStringProperty(props, "Title"), "Aman Sanoj"; got != want {
		t.Fatalf("title = %q, want %q", got, want)
	}
	if got, want := extractStringProperty(props, "Desc"), "Developer and Student"; got != want {
		t.Fatalf("desc = %q, want %q", got, want)
	}
}

func TestExtractDateAndURLProperty(t *testing.T) {
	props := map[string]interface{}{
		"Date": map[string]interface{}{
			"date": map[string]interface{}{
				"start": "2025-01-01",
				"end":   "2025-03-01",
			},
		},
		"URL": map[string]interface{}{
			"url": "https://example.com",
		},
	}

	if got, want := extractDateProperty(props, "Date"), "2025-01-01 – 2025-03-01"; got != want {
		t.Fatalf("date = %q, want %q", got, want)
	}
	if got, want := extractURLProperty(props, "URL"), "https://example.com"; got != want {
		t.Fatalf("url = %q, want %q", got, want)
	}
}

func TestParseRetryAfter(t *testing.T) {
	if got, want := parseRetryAfter("2"), 2*time.Second; got != want {
		t.Fatalf("parseRetryAfter(2) = %v, want %v", got, want)
	}
	if got := parseRetryAfter("bad"); got != 0 {
		t.Fatalf("parseRetryAfter(bad) = %v, want 0", got)
	}
}

func TestIsRetryableStatus(t *testing.T) {
	retryable := []int{429, 500, 502, 503, 504}
	for _, status := range retryable {
		if !isRetryableStatus(status) {
			t.Fatalf("status %d should be retryable", status)
		}
	}

	notRetryable := []int{400, 401, 403, 404}
	for _, status := range notRetryable {
		if isRetryableStatus(status) {
			t.Fatalf("status %d should not be retryable", status)
		}
	}
}
