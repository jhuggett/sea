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
	"github.com/jhuggett/sea/start"
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

			BlockList: []string{"backend/timeline", "backend/utils/callback"},
			Allowlist: []string{},

			WriteToFile: "log.txt",
		})),
	)

	slog.Debug("Starting server")

	db.Conn()
	db.Migrate()
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

func run(conn *websocket.Conn) {
	rpc := jsonrpc.New(conn)

	Connection := &start.Connection{
		RPC: rpc,
	}

	Timeline := timeline.New()

	/*

		If the client doesn't have a session, they register then login.

	*/

	var cleanUpGame func() = nil

	receivers := []func(){
		rpc.Receive("Login", inbound.WSLogin(func(snapshot game_context.Snapshot) inbound.Connection {
			slog.Info("Setting game context")
			Connection.GameCtx = game_context.New(snapshot)
			Connection.GameCtx.Timeline = Timeline

			cleanUpGame = start.Game(Connection)

			return Connection
		})),
		rpc.Receive("MoveShip", inbound.WSMoveShip(Connection)),
		rpc.Receive("Register", inbound.WSRegister()),
		rpc.Receive("GetWorldMap", inbound.WSGetWorldMap(Connection)),
		rpc.Receive("GetPorts", inbound.WSGetPorts(Connection)),
		rpc.Receive("ControlTime", inbound.ControlTime(Connection)),
		rpc.Receive("Trade", inbound.Trade(Connection)),
		rpc.Receive("PlotRoute", inbound.WSPlotRoute(Connection)),
		rpc.Receive("HireCrew", inbound.HireCrew(Connection)),
		rpc.Receive("RepairShip", inbound.RepairShip(Connection)),
		rpc.Receive("ManageRoute", inbound.ManageRoute(Connection)),
		rpc.Receive("GetHirablePeopleAtPort", inbound.GetHirablePeopleAtPort(Connection)),
	}

	<-rpc.ClosedChan

	Timeline.Stop()

	if cleanUpGame != nil {
		cleanUpGame()
	}

	for _, stopReceiving := range receivers {
		stopReceiving()
	}
}
