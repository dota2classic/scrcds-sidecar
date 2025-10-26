package srcds

import (
	"regexp"
	"strconv"
	"strings"
)

// PlayerMetric represents a player entry from `status` command output.
type PlayerMetric struct {
	UserID    int
	Name      string
	SteamID   string
	Connected string
	Ping      int
	Loss      int
	State     string
	Rate      int
	Adr       string
}

// parseStatusRow parses a single row of the `status` output.
// Returns nil if it's not a valid player row (e.g. a bot).
func parseStatusRow(row string) *PlayerMetric {
	sep := regexp.MustCompile(`(?:[^\s"]+|"[^"]*")+`)
	spaced := sep.FindAllString(row, -1)
	if len(spaced) < 7 {
		// likely a bot or malformed row
		return nil
	}

	// userid is the 3rd column (index 2)
	userID, err := strconv.Atoi(spaced[2])
	if err != nil {
		return nil
	}

	usernameStart := strings.Index(row, `"`)
	usernameEnd := strings.LastIndex(row, `"`)
	if usernameStart == -1 || usernameEnd == -1 || usernameEnd <= usernameStart {
		return nil
	}

	username := strings.TrimSpace(row[usernameStart+1 : usernameEnd])
	afterUsername := row[usernameEnd:]
	spaced = sep.FindAllString(afterUsername, -1)

	if len(spaced) < 7 {
		return nil
	}

	ping, _ := strconv.Atoi(spaced[2])
	loss, _ := strconv.Atoi(spaced[3])
	rate, _ := strconv.Atoi(spaced[5])

	return &PlayerMetric{
		UserID:    userID,
		Name:      username,
		SteamID:   strings.TrimSuffix(strings.TrimPrefix(spaced[0], "[U:1:"), `]`),
		Connected: spaced[1],
		Ping:      ping,
		Loss:      loss,
		State:     spaced[4],
		Rate:      rate,
		Adr:       spaced[6],
	}
}

// ParseStatusResponse parses the full `status` command output.
func ParseStatusResponse(raw string) []PlayerMetric {
	lines := strings.Split(raw, "\n")
	startIdx := -1

	for i, line := range lines {
		if strings.TrimSpace(line) == "# userid name uniqueid connected ping loss state rate adr" {
			startIdx = i + 1
			break
		}
	}
	if startIdx == -1 {
		return nil
	}

	var results []PlayerMetric
	for _, line := range lines[startIdx:] {
		if strings.TrimSpace(line) == "#end" {
			break
		}
		if metric := parseStatusRow(line); metric != nil {
			results = append(results, *metric)
		}
	}
	return results
}
