package worker

import (
	"context"
	"encoding/json"
	"fmt"
	aternos "github.com/sleeyax/aternos-api"
	"log"
	"time"
)

// New creates a new Worker instance.
func New(id string, options *aternos.Options) *Worker {
	w := &Worker{id: id, api: aternos.New(options)}
	w.Log("Created")
	return w
}

func (w *Worker) Log(msg string) {
	log.Printf("worker %s: %s\n", w.id, msg)
}

// Reconfigure reconfigures the worker with given options.
func (w *Worker) Reconfigure(options *aternos.Options) {
	w.api.Options = options
	w.serverInfo = nil
}

// Init initializes the worker.
func (w *Worker) Init() error {
	_, err := w.getWebsocketConnection()
	return err
}

// Start starts the minecraft server.
func (w *Worker) Start() error {
	w.Log("Starting server")
	return w.api.StartServer()
}

// Stop stops the minecraft server.
func (w *Worker) Stop() error {
	w.Log("Stopping server")
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
		w.Log("Background routines stopped & connections closed")
	}()

	for {
		select {
		case msg, ok := <-w.wss.Message:
			if !ok {
				if w.wssRetries == maxWssRetries {
					w.Log("Max number of websocket connection retries reached.")
					event("connection_error", nil) // TODO: find better way to communicate errors
					return
				}

				w.Log(fmt.Sprintf("Message channel closed. Trying to reconnect %d more time(s)...", maxWssRetries-w.wssRetries))
				time.Sleep(time.Second * 3)
				w.Init()
				w.wssRetries++
			}

			switch msg.Type {
			case "ready":
				event(msg.Type, w.serverInfo)

				// Start sending keep-alive requests in the background (until the server is offline, see below).
				go w.sendHeartBeats(ctxHeartBeat)

				w.Log("Connection ready")
			case "status":
				var info aternos.ServerInfo
				json.Unmarshal(msg.MessageBytes, &info)

				w.Log(fmt.Sprintf("Server is %s", info.StatusLabel))

				switch info.Status {
				case aternos.Online:
					if info.StatusLabelClass == "online" && info.Countdown != 0 {
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
