package worker

import aternos "github.com/sleeyax/aternos-api"

// The maximum number of times to try reconnecting to the websocket server on error.
const maxWssRetries = 3

type WorkersMap map[string]*Worker

// Worker represents a single worker that can connect to aternos for a specific message guild.
type Worker struct {
	// Unique worker id.
	id string

	// Current server status.
	// This is used as a cache for the info command and periodically updated once we receive a status update message over wss.
	serverInfo *aternos.ServerInfo

	// Current active websocket connection.
	wss *aternos.Websocket

	// Aternos API instance.
	api *aternos.Api

	// Current amount of websocket connection wssRetries (maximum maxWssRetries).
	wssRetries int
}
