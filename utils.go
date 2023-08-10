package main

import "regexp"

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func getStatusIcon(status string) string {
	switch status {
	case "Отвечен":
		return "🟢"
	case "Не отвечен":
		return "🔴"
	default:
		return "?"
	}
}
