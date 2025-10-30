package util

import (
	"path"
	"sidecar/internal/state"
	"strconv"
)

const (
	LOG_FOLDER    = "./dota/logs"
	REPLAY_FOLDER = "./dota/logs"
)

func GetLogFilePath() string {
	return path.Join(LOG_FOLDER, strconv.FormatInt(state.GlobalMatchInfo.MatchID, 10)+".log")
}
