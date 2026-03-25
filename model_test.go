package main

import "testing"

func TestParseMonthYear(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "valid", in: "2025-03-01", want: "Mar 2025"},
		{name: "trimmed", in: " 2025-12-20 ", want: "Dec 2025"},
		{name: "invalid month", in: "2025-13-01", want: ""},
		{name: "bad format", in: "2025", want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := parseMonthYear(tc.in); got != tc.want {
				t.Fatalf("parseMonthYear(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestFormatDateRange(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "range", in: "2025-01-01 – 2025-03-01", want: "Jan 2025 – Mar 2025"},
		{name: "single", in: "2025-06-15", want: "Jun 2025"},
		{name: "invalid", in: "bad-value", want: ""},
		{name: "empty", in: "", want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := formatDateRange(tc.in); got != tc.want {
				t.Fatalf("formatDateRange(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
