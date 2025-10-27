package log_parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func ParseLog(input string) (ParsedProtobufMessage, error) {
	re := regexp.MustCompile(`\d\d/\d\d/\d\d\d\d - \d\d:\d\d:\d\d: `)
	log := re.ReplaceAllString(input, "")
	startSignal := "SIGNOUT: Job created, Protobuf:"
	startIdx := strings.Index(log, startSignal)
	if startIdx == -1 {
		return ParsedProtobufMessage{}, fmt.Errorf("start signal not found")
	}
	startIdx += len(startSignal) + 1
	endIdx := strings.Index(log, "\ncluster_id")
	if endIdx == -1 {
		endIdx = len(log)
	}
	raw := log[startIdx:endIdx]
	clean := regexp.MustCompile(`#.*$`).ReplaceAllString(raw, "")
	clean = strings.TrimSpace(clean)

	tokens := tokenize(clean)
	obj, _ := parseObject(tokens, 0)

	// Collapse repeated arrays
	collapseRepeated(obj, map[string]bool{
		"players": true, "teams": true, "items": true,
		"tower_status": true, "barracks_status": true, "ability_upgrades": true,
	})

	b, _ := json.Marshal(obj)
	var msg ParsedProtobufMessage
	err := json.Unmarshal(b, &msg)
	return msg, err
}

func collapseRepeated(obj map[string]any, noCollapse map[string]bool) {
	for k, v := range obj {
		switch vv := v.(type) {
		case []any:
			for _, inner := range vv {
				if innerMap, ok := inner.(map[string]any); ok {
					collapseRepeated(innerMap, noCollapse)
				}
			}
			if len(vv) == 1 && !noCollapse[k] {
				obj[k] = vv[0]
			}
		case map[string]any:
			collapseRepeated(vv, noCollapse)
		}
	}
}
