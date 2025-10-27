package log_parser

import (
	"regexp"
	"strconv"
	"strings"
)

func parseValue(token string) any {
	if len(token) == 0 {
		return token
	}

	// Remove surrounding quotes
	if (strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"")) ||
		(strings.HasPrefix(token, "'") && strings.HasSuffix(token, "'")) {
		return token[1 : len(token)-1]
	}

	// Check if token is numeric
	matched, _ := regexp.MatchString(`^-?\d+(\.\d+)?([eE][+-]?\d+)?$`, token)
	if matched {
		// Determine digits before decimal
		intPart := token
		if dotIdx := strings.Index(token, "."); dotIdx != -1 {
			intPart = token[:dotIdx]
		}
		intPart = strings.TrimLeft(intPart, "+-")
		const BigNumberThreshold = 15

		if len(intPart) > BigNumberThreshold {
			// Too big to safely convert: keep as string
			return token
		}

		// Try int64 first if no decimal
		if !strings.Contains(token, ".") {
			if i64, err := strconv.ParseInt(token, 10, 64); err == nil {
				return i64
			}
		}

		// Otherwise parse as float64
		if f64, err := strconv.ParseFloat(token, 64); err == nil {
			return f64
		}

		// fallback
		return token
	}

	// Boolean
	lower := strings.ToLower(token)
	if lower == "true" {
		return true
	}
	if lower == "false" {
		return false
	}

	return token
}
