package log_parser

import "regexp"

func tokenize(input string) []string {
	re := regexp.MustCompile(`(\{|\}|:)|([^\s\{\}:]+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	tokens := []string{}
	for _, m := range matches {
		if m[1] != "" {
			tokens = append(tokens, m[1])
		} else if m[2] != "" {
			tokens = append(tokens, m[2])
		}
	}
	return tokens
}
