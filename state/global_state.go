package state

import (
	"fmt"
	"os"
	"sidecar/models"
	"strconv"
)

type MatchInfo struct {
	MatchID       int64
	LobbyType     models.MatchmakingMode
	GameMode      models.DotaGameMode
	Host          string
	ServerAddress string
}

var GlobalMatchInfo = &MatchInfo{}

func InitGlobalState() {

	host := os.Getenv("NODE_IP")
	hostPort, _ := strconv.Atoi(os.Getenv("HOST_PORT"))

	serverAddress := fmt.Sprintf("%s:%d", host, hostPort)

	matchId, _ := strconv.ParseInt(os.Getenv("MATCH_ID"), 10, 64)

	GlobalMatchInfo.MatchID = matchId
	GlobalMatchInfo.LobbyType = models.ParseLobbyType(os.Getenv("LOBBY_TYPE"))
	GlobalMatchInfo.GameMode = models.ParseGameMode(os.Getenv("GAME_MODE"))
	GlobalMatchInfo.Host = host
	GlobalMatchInfo.ServerAddress = serverAddress
}
