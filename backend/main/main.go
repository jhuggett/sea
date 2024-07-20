package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/inbound"
	"github.com/jhuggett/sea/jsonrpc"
	"github.com/jhuggett/sea/models/ship"
	"github.com/jhuggett/sea/models/world_map"
	"github.com/jhuggett/sea/outbound"
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
	gameCtx inbound.GameContext
}

func (c *Connection) Context() inbound.GameContext {
	return c.gameCtx
}

func (c *Connection) Sender() *outbound.Sender {
	return outbound.NewSender(c.RPC)
}

func run(conn *websocket.Conn) {
	rpc := jsonrpc.New(conn)

	connection := &Connection{
		RPC: rpc,
	}

	receivers := []func(){
		rpc.Receive("Login", inbound.Login(func(gameCtx inbound.GameContext) {
			slog.Info("Setting game context", "gameCtx", gameCtx)
			connection.gameCtx = gameCtx
		})),
		rpc.Receive("MoveShip", inbound.MoveShip(connection)),
		rpc.Receive("Register", inbound.Register()),
		rpc.Receive("GetWorldMap", inbound.GetWorldMap(connection)),
	}

	<-rpc.ClosedChan

	for _, stopReceiving := range receivers {
		stopReceiving()
	}
}
