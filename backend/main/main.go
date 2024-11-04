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
	"github.com/jhuggett/sea/log"
	"github.com/jhuggett/sea/models/port"
	"github.com/jhuggett/sea/models/ship"
	"github.com/jhuggett/sea/models/world_map"
	"github.com/jhuggett/sea/outbound"
	"github.com/jhuggett/sea/timeline"
)

var upgrader = websocket.Upgrader{}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	slog.Info("Upgrading Connection")
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		slog.Error("Error upgrading Connection: ", err)
		return
	}

	slog.Info("Connection upgraded")

	run(c)
}

func main() {
	// set global logger with custom options
	slog.SetDefault(
		slog.New(log.NewHandler(&log.HandlerOptions{
			HandlerOptions: slog.HandlerOptions{
				AddSource: true,
				Level:     log.OptInDebug,
			},
			UseColor: true,

			BlockList: []string{
				"utils/callback",
			},
			Allowlist: []string{},
		})),
	)

	slog.Debug("Starting server")

	dbConn := db.Conn()
	dbConn.AutoMigrate(&ship.Ship{})
	dbConn.AutoMigrate(&world_map.WorldMap{})
	dbConn.AutoMigrate(&world_map.CoastalPoint{})
	dbConn.AutoMigrate(&world_map.Continent{})
	dbConn.AutoMigrate(&port.Port{})
	defer db.Close()

	http.HandleFunc("/ws", wsHandler)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Error starting server: ", err)
	}

	slog.Debug("All done")
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

	Connection := &Connection{
		RPC: rpc,
	}

	Timeline := timeline.New()

	timeline.OnTicksPerSecondChanged.Register(func(s struct {
		Old       uint64
		New       uint64
		TickCount uint64
	}) {
		Connection.Sender().TimeChanged(s.New, s.TickCount)
	})

	Timeline.Start()

	receivers := []func(){
		rpc.Receive("Login", inbound.Login(func(snapshot game_context.Snapshot) inbound.Connection {
			slog.Info("Setting game context")
			Connection.gameCtx = game_context.New(snapshot)
			Connection.gameCtx.Timeline = Timeline

			return Connection
		})),
		rpc.Receive("MoveShip", inbound.MoveShip(Connection)),
		rpc.Receive("Register", inbound.Register()),
		rpc.Receive("GetWorldMap", inbound.GetWorldMap(Connection)),
		rpc.Receive("GetPorts", inbound.GetPorts(Connection)),
		rpc.Receive("ControlTime", inbound.ControlTime(Connection)),
	}

	<-rpc.ClosedChan

	Timeline.Stop()

	for _, stopReceiving := range receivers {
		stopReceiving()
	}
}
