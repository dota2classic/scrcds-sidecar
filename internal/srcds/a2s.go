package srcds

import (
	"errors"
	"fmt"
	"net"
	"sidecar/internal/state"
	"time"
)

// a2sInfoQuery is the A2S_INFO request payload per the Valve server query protocol.
var a2sInfoQuery = []byte("\xFF\xFF\xFF\xFF\x54Source Engine Query\x00")

// isServerReady sends an A2S_INFO UDP query and returns true if the server responds.
// Handles the Valve anti-flood challenge (0x41) transparently.
// This works before RCON is ready and does not trigger RCON ban protection.
func isServerReady() bool {
	addr := fmt.Sprintf("127.0.0.1:%d", state.GlobalMatchInfo.GameServerPort)

	conn, err := net.DialTimeout("udp", addr, 2*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	udpQuery := func(payload []byte) ([]byte, error) {
		_ = conn.SetDeadline(time.Now().Add(2 * time.Second))
		if _, err := conn.Write(payload); err != nil {
			return nil, err
		}
		buf := make([]byte, 1400)
		n, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}
		return buf[:n], nil
	}

	resp, err := udpQuery(a2sInfoQuery)
	if err != nil {
		return false
	}

	// Valve anti-flood: server sends a challenge (0x41) that must be echoed back.
	if len(resp) >= 9 && resp[4] == 0x41 {
		withChallenge := append(a2sInfoQuery, resp[5:9]...)
		resp, err = udpQuery(withChallenge)
		if err != nil {
			return false
		}
	}

	return len(resp) > 5
}

// AwaitServerReady polls via A2S_INFO until the server responds or the timeout elapses.
// Unlike RCON polling, this never triggers srcds's rcon-ban protection.
func AwaitServerReady(maxWait time.Duration) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	deadline := time.Now().Add(maxWait)

	for range ticker.C {
		if time.Now().After(deadline) {
			return errors.New("timed out waiting for server to start")
		}
		if isServerReady() {
			return nil
		}
	}
	return nil
}
