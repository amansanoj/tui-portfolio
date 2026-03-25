package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func (m Model) statusPageWrapWidth(mainWidth int) int {
	wrapWidth := mainWidth - 2 - 2 - 2
	if wrapWidth < 30 {
		wrapWidth = 30
	}
	return wrapWidth
}

func (m Model) statusPageRawLines() []string {
	snapshot := appContentStore.Snapshot()
	refreshValue := os.Getenv(contentRefreshEnvVar)
	if refreshValue == "" {
		refreshValue = "300"
	}

	cacheState := "warm"
	cacheAge := "unknown"
	lastLoaded := "not loaded yet"
	if !snapshot.Ready {
		cacheState = "cold"
	} else {
		age := time.Since(snapshot.LoadedAt)
		if age < 0 {
			age = 0
		}
		cacheAge = age.Round(time.Second).String()
		lastLoaded = snapshot.LoadedAt.Local().Format(time.RFC1123)
	}

	refreshing := "no"
	if snapshot.Refreshing {
		refreshing = "yes"
	}

	return []string{
		"Runtime Status",
		"",
		fmt.Sprintf("Cache state        : %s", cacheState),
		fmt.Sprintf("Cache refreshing   : %s", refreshing),
		fmt.Sprintf("Cache age          : %s", cacheAge),
		fmt.Sprintf("Last loaded        : %s", lastLoaded),
		fmt.Sprintf("Projects cached    : %d", len(snapshot.Projects)),
		fmt.Sprintf("Certifications     : %d", len(snapshot.Certifications)),
		"",
		fmt.Sprintf("Refresh interval   : %ss", refreshValue),
		fmt.Sprintf("Listen address     : %s", envWithDefault("APP_ADDR", defaultSSHAddress)),
		fmt.Sprintf("Host key path      : %s", envWithDefault("HOST_KEY_PATH", defaultHostKeyPath)),
	}
}

func (m Model) statusPageLineCount(mainWidth int) int {
	raw := m.statusPageRawLines()
	wrapWidth := m.statusPageWrapWidth(mainWidth)
	total := 0

	for _, line := range raw {
		if line == "" {
			total++
			continue
		}
		wrapped := wordWrap(line, wrapWidth)
		if len(wrapped) == 0 {
			total++
			continue
		}
		total += len(wrapped)
	}

	return total
}

func (m Model) renderRuntimeStatus(mainWidth, scroll, avail int) string {
	raw := m.statusPageRawLines()
	wrapWidth := m.statusPageWrapWidth(mainWidth)

	var expanded []string
	for _, line := range raw {
		if line == "" {
			expanded = append(expanded, "")
			continue
		}
		wrapped := wordWrap(line, wrapWidth)
		if len(wrapped) == 0 {
			expanded = append(expanded, "")
			continue
		}
		expanded = append(expanded, wrapped...)
	}

	visible := m.clampVisibleLines(expanded, scroll, avail)

	var out []string
	for i, line := range visible {
		if line == "" {
			out = append(out, "")
			continue
		}
		if i == 0 && scroll == 0 {
			out = append(out, m.styles.sectionHeaderStyle.Render(line))
		} else {
			out = append(out, m.styles.contentStyle.Render(line))
		}
	}

	return strings.Join(out, "\n")
}

func envWithDefault(name, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	return value
}
