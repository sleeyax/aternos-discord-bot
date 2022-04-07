package worker

import (
	"context"
	aternos "github.com/sleeyax/aternos-api"
	"time"
)

// New creates a new Worker instance.
func New(options *aternos.Options) *Worker {
	return &Worker{api: aternos.New(options)}
}

func (w *Worker) Reconfigure(options *aternos.Options) {
	// TODO: check performance of this call
	w.api = aternos.New(options)
}

func (w *Worker) Start() error {
	/*conn, err := w.getWebsocketConnection()
	if err != nil {
		return err
	}*/

	return nil
}

// GetServerInfo returns minecraft server information.
// If the info is not known yet, it wil be fetched and cached in memory for later use.
func (w *Worker) GetServerInfo() (*aternos.ServerInfo, error) {
	if w.serverInfo == nil {
		info, err := w.api.GetServerInfo()
		if err != nil {
			return nil, err
		}

		if info.DynIP == "" {
			info.DynIP = "unavailable"
		}

		w.serverInfo = &info
	}

	return w.serverInfo, nil
}

// getWebsocketConnection returns the current active connection.
// If no connection has been instantiated yet, it will connect to the Aternos websocket server and cache the connection in memory for later use.
func (w *Worker) getWebsocketConnection() (*aternos.Websocket, error) {
	if w.wss == nil || !w.wss.IsConnected() {
		wss, err := w.api.ConnectWebSocket()
		if err != nil {
			return nil, err
		}
		w.wss = wss
	}

	return w.wss, nil
}

// sendHeartBeats sends keep alive packets over the websocket connection by the default interval.
func (w *Worker) sendHeartBeats(context context.Context) {
	ticker := time.NewTicker(time.Millisecond * 49000)

	for {
		select {
		case <-context.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			w.wss.SendHeartBeat()
		}
	}
}
