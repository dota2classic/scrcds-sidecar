package state

import (
	"fmt"
	"os"
	"sidecar/internal/models"
	"strconv"

	d2cmodels "github.com/dota2classic/d2c-go-models/models"
)

type MatchInfo struct {
	MatchID       int64
	LobbyType     d2cmodels.MatchmakingMode
	GameMode      d2cmodels.DotaGameMode
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

	fmt.Println("Initialized global state: ", GlobalMatchInfo)
}
