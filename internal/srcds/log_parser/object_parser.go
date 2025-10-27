package log_parser

func parseObject(tokens []string, start int) (map[string]any, int) {
	obj := map[string]any{}
	i := start
	for i < len(tokens) {
		token := tokens[i]
		if token == "}" {
			return obj, i + 1
		}
		if token == "{" {
			i++
			continue
		}
		key := token
		i++
		if i < len(tokens) && tokens[i] == ":" {
			i++
			if tokens[i] == "{" {
				i++
				nested, next := parseObject(tokens, i)
				addToResult(obj, key, nested)
				i = next
			} else {
				val := parseValue(tokens[i])
				addToResult(obj, key, val)
				i++
			}
		} else if i < len(tokens) && tokens[i] == "{" {
			i++
			nested, next := parseObject(tokens, i)
			addToResult(obj, key, nested)
			i = next
		} else {
			val := parseValue(tokens[i])
			addToResult(obj, key, val)
			i++
		}
	}
	return obj, i
}

func addToResult(obj map[string]any, key string, val any) {
	if _, exists := obj[key]; !exists {
		obj[key] = []any{}
	}
	obj[key] = append(obj[key].([]any), val)
}
