package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/inbound"
	"github.com/jhuggett/sea/jsonrpc"
	"github.com/jhuggett/sea/models/ship"
	"github.com/jhuggett/sea/models/world_map"
	"github.com/jhuggett/sea/outbound"
	"github.com/jhuggett/sea/timeline"
)

var upgrader = websocket.Upgrader{}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	slog.Info("Upgrading connection")
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		slog.Error("Error upgrading connection: ", err)
		return
	}

	slog.Info("Connection upgraded")

	run(c)
}

func main() {
	dbConn := db.Conn()
	dbConn.AutoMigrate(&ship.Ship{})
	dbConn.AutoMigrate(&world_map.WorldMap{})
	dbConn.AutoMigrate(&world_map.CoastalPoint{})
	dbConn.AutoMigrate(&world_map.Continent{})
	defer db.Close()

	http.HandleFunc("/ws", wsHandler)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}

type ExamplePayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Connection struct {
	RPC     jsonrpc.JSONRPC
	gameCtx *game_context.GameContext
}

func (c *Connection) Context() *game_context.GameContext {
	return c.gameCtx
}

func (c *Connection) Sender() *outbound.Sender {
	return outbound.NewSender(c.RPC, c.gameCtx)
}

func run(conn *websocket.Conn) {
	rpc := jsonrpc.New(conn)

	connection := &Connection{
		RPC: rpc,
	}

	Timeline := timeline.New()

	Timeline.Start()

	receivers := []func(){
		rpc.Receive("Login", inbound.Login(func(snapshot game_context.Snapshot) *game_context.GameContext {
			slog.Info("Setting game context")
			connection.gameCtx = game_context.New(snapshot)
			connection.gameCtx.Timeline = Timeline

			return connection.gameCtx
		})),
		rpc.Receive("MoveShip", inbound.MoveShip(connection)),
		rpc.Receive("Register", inbound.Register()),
		rpc.Receive("GetWorldMap", inbound.GetWorldMap(connection)),
	}

	<-rpc.ClosedChan

	Timeline.Stop()

	for _, stopReceiving := range receivers {
		stopReceiving()
	}
}
