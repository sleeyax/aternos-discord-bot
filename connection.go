package aternos_discord_bot

import (
	"context"
	aternos "github.com/sleeyax/aternos-api"
	"time"
)

// getServerInfo returns the current server info.
// If the status is not known yet, it wil be (re)fetched and cached in memory.
func (c *Connection) getServerInfo() (*aternos.ServerInfo, error) {
	if c.serverInfo == nil {
		info, err := c.api.GetServerInfo()
		if err != nil {
			return nil, err
		}
		c.serverInfo = &info
	}

	return c.serverInfo, nil
}

// getWSS connects to the Aternos websocket server and stores the active connection in memory for later use.
func (c *Connection) getWSS() (*aternos.Websocket, error) {
	if c.wss == nil || !c.wss.IsConnected() {
		wss, err := c.api.ConnectWebSocket()
		if err != nil {
			return nil, err
		}
		c.wss = wss
	}

	return c.wss, nil
}

func (c *Connection) sendHeartBeats(context context.Context) {
	ticker := time.NewTicker(time.Millisecond * 49000)

	for {
		select {
		case <-context.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			c.wss.SendHeartBeat()
		}
	}
}
