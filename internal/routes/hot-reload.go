package routes

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (app *Application) hotReload(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	for {
		app.logger.Debug("hot reload connection successful")
		msgType, _, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if err = conn.WriteMessage(msgType, []byte("connected")); err != nil {
			return
		}
	}
}

func (app *Application) testAlive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
