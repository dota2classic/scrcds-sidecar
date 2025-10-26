package srcds

import (
	"fmt"
	"strconv"
	"strings"
)

// SrcdsServerMetrics represents server stats
type SrcdsServerMetrics struct {
	CPU     float64
	In      float64
	Out     float64
	Uptime  float64
	Users   int
	FPS     float64
	Players int
}

// ParseStatsResponse parses the raw stats string into SrcdsServerMetrics
func ParseStatsResponse(statsRaw string) (*SrcdsServerMetrics, error) {
	lines := strings.Split(statsRaw, "\n")
	if len(lines) < 3 {
		return nil, fmt.Errorf("invalid stats format")
	}

	// remove first and last line
	lines = lines[1 : len(lines)-1]
	fields := strings.Fields(lines[0])

	if len(fields) < 7 {
		return nil, fmt.Errorf("not enough fields in stats line")
	}

	cpu, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return nil, err
	}
	in, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, err
	}
	out, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, err
	}
	uptime, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, err
	}
	users, err := strconv.Atoi(fields[4])
	if err != nil {
		return nil, err
	}
	fps, err := strconv.ParseFloat(fields[5], 64)
	if err != nil {
		return nil, err
	}
	players, err := strconv.Atoi(fields[6])
	if err != nil {
		return nil, err
	}

	return &SrcdsServerMetrics{
		CPU:     cpu,
		In:      in,
		Out:     out,
		Uptime:  uptime,
		Users:   users,
		FPS:     fps,
		Players: players,
	}, nil
}
