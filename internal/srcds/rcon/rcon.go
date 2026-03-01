package rcon

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sidecar/internal/state"
	"sync"
	"time"

	"github.com/gorcon/rcon"
)

const connMaxAge = 5 * time.Minute

var (
	cachedConn    *rcon.Conn
	connCreatedAt time.Time
	connMu        sync.Mutex
)

// RunRconCommand executes an RCON command, reusing a cached connection where possible.
// If the connection is broken or stale, it is transparently replaced.
func RunRconCommand(cmd string) (string, error) {
	connMu.Lock()
	conn, err := getOrDial()
	connMu.Unlock()

	if err != nil {
		return "", err
	}

	result, err := conn.Execute(cmd)
	if err != nil {
		log.Printf("Failed to execute RCON command: %v", err)
		invalidate()
		return "", err
	}

	return result, nil
}

// InvalidateRconConnection closes and clears the cached connection.
func InvalidateRconConnection() {
	invalidate()
}

// getOrDial returns the cached connection or dials a new one.
// Must be called with connMu held.
func getOrDial() (*rcon.Conn, error) {
	if cachedConn != nil && time.Since(connCreatedAt) < connMaxAge {
		return cachedConn, nil
	}

	addr := fmt.Sprintf("127.0.0.1:%d", state.GlobalMatchInfo.GameServerPort)
	rconPassword := os.Getenv("RCON_PASSWORD")

	hasher := md5.New()
	hasher.Write([]byte(rconPassword))
	passwordHash := hex.EncodeToString(hasher.Sum(nil))

	log.Printf("Connecting RCON to %s with password hash %s", addr, passwordHash)
	conn, err := rcon.Dial(addr, rconPassword)
	if err != nil {
		return nil, err
	}

	cachedConn = conn
	connCreatedAt = time.Now()
	return conn, nil
}

func invalidate() {
	connMu.Lock()
	defer connMu.Unlock()

	if cachedConn != nil {
		_ = cachedConn.Close()
		cachedConn = nil
	}
}
