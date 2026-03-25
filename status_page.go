package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func (m Model) renderRuntimeStatus(mainWidth int) string {
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

	lines := []string{
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

	wrapWidth := mainWidth - 2 - 2 - 2
	if wrapWidth < 30 {
		wrapWidth = 30
	}

	var out []string
	for i, line := range lines {
		if line == "" {
			out = append(out, "")
			continue
		}
		wrapped := wordWrap(line, wrapWidth)
		for _, w := range wrapped {
			if i == 0 {
				out = append(out, m.styles.sectionHeaderStyle.Render(w))
			} else {
				out = append(out, m.styles.contentStyle.Render(w))
			}
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
