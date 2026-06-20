package xlog

import "strings"

func maskDefault(s string) string {
	if len(s) <= 4 {
		return "***"
	}
	return s[:2] + "****" + s[len(s)-2:]
}

func maskEmail(s string) string {
	parts := strings.Split(s, "@")
	if len(parts) != 2 {
		return "***"
	}
	name := parts[0]
	if len(name) <= 2 {
		return "***@" + parts[1]
	}
	return name[:2] + "***@" + parts[1]
}

func maskPhone(s string) string {
	if len(s) < 7 {
		return "***"
	}
	return s[:3] + "****" + s[len(s)-4:]
}

func maskPassword(string) string {
	return "******"
}
