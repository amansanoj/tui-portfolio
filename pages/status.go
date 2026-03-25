package pages

import "fmt"

type StatusData struct {
	CacheState          string
	Refreshing          bool
	CacheAge            string
	LastLoaded          string
	ProjectsCount       int
	CertificationsCount int
	RefreshInterval     string
	ListenAddress       string
	HostKeyPath         string
}

func BuildStatusLines(data StatusData) []string {
	refreshing := "no"
	if data.Refreshing {
		refreshing = "yes"
	}

	return []string{
		fmt.Sprintf("Cache state        : %s", data.CacheState),
		fmt.Sprintf("Cache refreshing   : %s", refreshing),
		fmt.Sprintf("Cache age          : %s", data.CacheAge),
		fmt.Sprintf("Last loaded        : %s", data.LastLoaded),
		fmt.Sprintf("Projects cached    : %d", data.ProjectsCount),
		fmt.Sprintf("Certifications     : %d", data.CertificationsCount),
		"",
		fmt.Sprintf("Refresh interval   : %s", data.RefreshInterval),
		fmt.Sprintf("Listen address     : %s", data.ListenAddress),
		fmt.Sprintf("Host key path      : %s", data.HostKeyPath),
	}
}
