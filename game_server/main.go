package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	WS "golang.org/x/net/websocket"
	"log"
	"net/http"
	"strings"
)

const (
	MessageTypeSystem    = "system"
	MessageTypeUserClick = "click"
)

type MessageType struct {
	Type string `json:"type"`
}

type SystemMessage struct {
	Message string `json:"message"`
	Channel int    `json:"channel"`
	Version string `json:"version"`
}

type UserClickMessage struct {
	Id string `json:"id"`
	X  string `json:"x"`
	Y  string `json:"y"`
}

var wsServer *websocket.Conn

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	setupAPI()

	heatmapConn := createHeatmapConnection()

	defer heatmapConn.Close()

	go handleMessage(heatmapConn)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func setupAPI() {
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		wsServer = websocketHandler(writer, request)
	})
}

func createHeatmapConnection() *WS.Conn {
	channelId := 44420434
	origin := "http://localhost/"
	url := fmt.Sprintf("wss://heat-api.j38.net/channel/%d", channelId)
	heatmapConn, err := WS.Dial(url, "", origin)

	if err != nil {
		log.Fatal(err)
	}

	return heatmapConn
}

func handleMessage(ws *WS.Conn) {
	for {
		msg := make([]byte, 512)

		if _, err := ws.Read(msg); err != nil {
			log.Fatal(err)
		}

		if len(msg) == 0 {
			continue
		}

		msg = []byte(strings.Trim(string(msg), "\x00"))

		var message MessageType

		if err := json.Unmarshal(msg, &message); err != nil {
			log.Printf("error decoding message: %v", err)
			continue
		}

		if wsServer != nil {
			if err := wsServer.WriteMessage(1, msg); err != nil {
				log.Println("Error writing message:", err)
				return
			}
		}

		switch message.Type {
		case MessageTypeSystem:
			var systemMessage SystemMessage

			if err := json.Unmarshal(msg, &systemMessage); err != nil {
				log.Printf("error decoding message: %v", err)
				continue
			}

			fmt.Printf("System message: %s\n", systemMessage.Message)
		case MessageTypeUserClick:
			var userClickMessage UserClickMessage

			if err := json.Unmarshal(msg, &userClickMessage); err != nil {
				log.Printf("error decoding message: %v", err)
				continue
			}

			fmt.Printf("User click message: %s - x:%s y:%s\n", userClickMessage.Id, userClickMessage.X, userClickMessage.Y)
		}
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return nil
	}

	return conn
}
