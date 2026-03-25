package pages

import (
	"fmt"
	"strings"
)

type ContactData struct {
	Label  string
	Handle string
	URL    string
}

func BuildContactBody(items []ContactData) string {
	var sb strings.Builder
	for i, item := range items {
		sb.WriteString(fmt.Sprintf("CONTACT|||%s|||%s|||%s\n", item.Label, item.Handle, item.URL))
		if i < len(items)-1 {
			sb.WriteString("\n")
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
