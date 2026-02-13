package rcon

import (
	"fmt"
	"log"
	"os"
	"sidecar/internal/state"

	"github.com/gorcon/rcon"
)

func RunRconCommand(cmd string) (string, error) {
	conn, err := GetRconConnection()
	if err != nil {
		return "", err
	}

	defer conn.Close()

	return RunRconCommandOnConnection(conn, cmd)
}

func GetRconConnection() (*rcon.Conn, error) {
	// We use localhost to not trigger rcon banning
	addr := fmt.Sprintf("127.0.0.1:%d", state.GlobalMatchInfo.GameServerPort)

	rconPassword := os.Getenv("RCON_PASSWORD")

	conn, err := rcon.Dial(addr, rconPassword)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func RunRconCommandOnConnection(conn *rcon.Conn, cmd string) (string, error) {
	status, err := conn.Execute(cmd)
	if err != nil {
		log.Printf("Failed to execute RCON command: %v", err)
		return "", err
	}

	return status, nil
}
