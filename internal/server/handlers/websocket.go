package handlers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	conns    = make(map[*websocket.Conn]bool)
	m        = sync.Mutex{}
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func GetPassedSessionTime(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("failed to connect: ", err)
		http.Error(w, "failed to connecti", http.StatusInternalServerError)
		return
	}
	defer func() {
		m.Lock()
		defer m.Unlock()
		deleteConnection(conn)
	}()

	m.Lock()
	conns[conn] = true
	m.Unlock()

	for {
		_, _, err = conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func SendTime(time string) {
	m.Lock()
	defer m.Unlock()

	for conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(time)); err != nil {
			deleteConnection(conn)
			continue
		}
	}
}

func deleteConnection(conn *websocket.Conn) (err error) {
	if err = conn.Close(); err != nil {
		fmt.Println("failed to close connection: ", err)
		return
	}

	delete(conns, conn)
	return
}
