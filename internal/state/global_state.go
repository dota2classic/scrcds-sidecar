package state

import (
	"fmt"
	"os"
	"sidecar/internal/models"
	"strconv"

	d2cmodels "github.com/dota2classic/d2c-go-models/models"
)

type MatchInfo struct {
	MatchID        int64
	LobbyType      d2cmodels.MatchmakingMode
	GameMode       d2cmodels.DotaGameMode
	Host           string
	GameServerPort int
	SourceTvPort   int
	ServerAddress  string
}

var GlobalMatchInfo = &MatchInfo{}

func InitGlobalState() {

	host := os.Getenv("NODE_IP")
	gamePort, _ := strconv.Atoi(os.Getenv("HOST_PORT"))
	tvPort, _ := strconv.Atoi(os.Getenv("HOST_TV_PORT"))

	serverAddress := fmt.Sprintf("%s:%d", host, gamePort)

	matchId, _ := strconv.ParseInt(os.Getenv("MATCH_ID"), 10, 64)

	GlobalMatchInfo.MatchID = matchId
	GlobalMatchInfo.LobbyType = models.ParseLobbyType(os.Getenv("LOBBY_TYPE"))
	GlobalMatchInfo.GameMode = models.ParseGameMode(os.Getenv("GAME_MODE"))
	GlobalMatchInfo.Host = host
	GlobalMatchInfo.GameServerPort = gamePort
	GlobalMatchInfo.SourceTvPort = tvPort
	GlobalMatchInfo.ServerAddress = serverAddress

	fmt.Println("Initialized global state: ", GlobalMatchInfo)
}
