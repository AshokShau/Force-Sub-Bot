package modules

import (
	"fmt"
	"strings"
	"time"
)

// getFormattedDuration returns a formatted string representing the duration
func getFormattedDuration(diff time.Duration) string {
	seconds := int(diff.Seconds())
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24
	weeks := days / 7
	months := days / 30

	// Calculate remaining values after each unit is accounted for
	remainingSeconds := seconds % 60
	remainingMinutes := minutes % 60
	remainingHours := hours % 24
	remainingDays := days % 7

	var text string

	// Format months
	if months != 0 {
		text += fmt.Sprintf("%d months ", months)
	}

	// Format weeks
	if weeks != 0 {
		text += fmt.Sprintf("%d weeks ", weeks)
	}

	// Format days
	if remainingDays != 0 {
		text += fmt.Sprintf("%d days ", remainingDays)
	}

	// Format hours
	if remainingHours != 0 {
		text += fmt.Sprintf("%d hours ", remainingHours)
	}

	// Format minutes
	if remainingMinutes != 0 {
		text += fmt.Sprintf("%d minutes ", remainingMinutes)
	}

	// Format seconds
	if remainingSeconds != 0 || text == "" { // Include seconds if there's no larger unit or if there are remaining seconds
		text += fmt.Sprintf("%d seconds", remainingSeconds)
	}

	// Trim any trailing space
	text = trimSuffix(text, " ")

	return text
}

// trimSuffix removes the suffix from the string if it exists
func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}
