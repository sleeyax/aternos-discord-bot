package worker

import (
	"context"
	"encoding/json"
	aternos "github.com/sleeyax/aternos-api"
	"log"
	"time"
)

// New creates a new Worker instance.
func New(id string, options *aternos.Options) *Worker {
	return &Worker{id: id, api: aternos.New(options)}
}

func (w *Worker) log(msg string) {
	log.Printf("worker %s: %s\n", w.id, msg)
}

// Reconfigure reconfigures the worker with given options.
func (w *Worker) Reconfigure(options *aternos.Options) {
	// TODO: check performance of this call
	w.api = aternos.New(options)
}

// Init initializes the worker.
func (w *Worker) Init() error {
	_, err := w.getWebsocketConnection()
	return err
}

// Start starts the minecraft server.
func (w *Worker) Start() error {
	return w.api.StartServer()
}

// Stop stops the minecraft server.
func (w *Worker) Stop() error {
	return w.api.StopServer()
}

// On can be used for event handling.
//
// Once Init and Start have been called on the worker, the event handler will be fired whenever an interesting websocket message is received.
func (w *Worker) On(ctx context.Context, event func(messageType string, info *aternos.ServerInfo)) {
	ctxHeartBeat, cancelHeartBeat := context.WithCancel(context.Background())
	ctxConfirm, cancelConfirm := context.WithCancel(context.Background())

	defer func() {
		cancelHeartBeat()
		cancelConfirm()
		w.wss.Close()
		w.wss = nil
		// TODO: log worker number (based on discord id?)
		w.log("Background routines stopped & connections closed")
	}()

	for {
		select {
		case msg, ok := <-w.wss.Message:
			if !ok {
				w.log("Message channel closed. Tying to reconnect...")
				w.Init()
			}

			switch msg.Type {
			case "ready":
				event(msg.Type, w.serverInfo)

				// Start sending keep-alive requests in the background (until the server is offline, see below).
				go w.sendHeartBeats(ctxHeartBeat)
			case "status":
				var info aternos.ServerInfo
				json.Unmarshal(msg.MessageBytes, &info)

				switch info.Status {
				case aternos.Online:
					if info.StatusLabelClass == "online" && w.serverInfo.StatusLabelClass != "online" {
						event(msg.Type, &info)
					}
				case aternos.Offline:
					event(msg.Type, &info)
					w.serverInfo = &info
					return
				case aternos.Preparing: // stuck in queue (only happens when traffic is high)
					if (info.StatusLabelClass == "queueing" && info.Queue.Status == "pending") && (w.serverInfo.Queue.Status != "pending") {
						event(msg.Type, &info)
						go w.api.ConfirmServer(ctxConfirm, 10*time.Second)
					}
				case aternos.Loading:
					event(msg.Type, &info)
					cancelConfirm()
				}

				w.serverInfo = &info

			}
		case <-ctx.Done():
			return
		}
	}
}

// GetServerInfo returns minecraft server information.
// If the info is not known yet, it wil be fetched and cached in memory for later use.
func (w *Worker) GetServerInfo() (*aternos.ServerInfo, error) {
	if w.serverInfo == nil {
		info, err := w.api.GetServerInfo()
		if err != nil {
			return nil, err
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
