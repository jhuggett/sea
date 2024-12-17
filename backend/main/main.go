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
	"github.com/jhuggett/sea/models"
	"github.com/jhuggett/sea/models/inventory"
	"github.com/jhuggett/sea/models/ship"
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

			BlockList: []string{},
			Allowlist: []string{},
		})),
	)

	slog.Debug("Starting server")

	dbConn := db.Conn()
	dbConn.AutoMigrate(&models.Ship{})
	dbConn.AutoMigrate(&models.WorldMap{})
	dbConn.AutoMigrate(&models.Point{})
	dbConn.AutoMigrate(&models.Continent{})
	dbConn.AutoMigrate(&models.Port{})
	dbConn.AutoMigrate(&models.Crew{})
	dbConn.AutoMigrate(&models.Inventory{})
	dbConn.AutoMigrate(&models.Item{})
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

	/*

		If the client doesn't have a session, they register then login.

	*/

	receivers := []func(){
		rpc.Receive("Login", inbound.Login(func(snapshot game_context.Snapshot) inbound.Connection {
			slog.Info("Setting game context")
			Connection.gameCtx = game_context.New(snapshot)
			Connection.gameCtx.Timeline = Timeline

			startGame(Connection)

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

func startGame(conn *Connection) {
	Timeline := conn.Context().Timeline

	Timeline.OnTicksPerCycleChangedDo(func(data timeline.TicksPerCycleChangedEventData) {
		conn.Sender().TimeChanged(data.CurrentTick, data.NewTicksPerCycle)
	})

	Timeline.Start()

	s, err := conn.Context().Ship()
	if err != nil {
		slog.Error("Error getting ship", "err", err)
		panic(err)
	}

	Timeline.Do(func() uint64 {
		slog.Info("A day has passed")

		// Pay wages (probably should be payed in a different way later)

		crew, err := s.Crew()
		if err != nil {
			slog.Error("Error getting crew", "err", err)
			return timeline.Day
		}

		slog.Info("Got crew", "crew", crew)

		inventory, err := s.Inventory().Fetch()
		if err != nil {
			slog.Error("Error fetching inventory", "err", err)
			return timeline.Day
		}

		err = inventory.RemoveItem(models.Item{
			Name:   string(models.ItemTypePieceOfEight),
			Amount: float32(crew.Persistent.Wage),
		})

		if err != nil {
			slog.Error("Error removing item", "err", err)
		}

		return timeline.Day
	}, timeline.Day)

	Timeline.Do(func() uint64 {
		conn.Sender().TimeChanged(Timeline.CurrentTick(), Timeline.TicksPerCycle())

		return 1
	}, 1)

	s.Inventory().OnChangedDo(func(data inventory.OnChangedEventData) {
		slog.Info("Inventory changed", "id", s.Persistent.ID, "inventory", data.Inventory)
		conn.Sender().ShipInventoryChanged(s.Persistent.ID, data.Inventory)
	})

	var stopFishing func()

	s.OnAnchorRaisedDo(func(data ship.AnchorRaisedEventData) {
		if stopFishing != nil {
			stopFishing()
		}
	})

	s.OnAnchorLoweredDo(func(data ship.AnchorLoweredEventData) {
		if data.Location == ship.AnchorLoweredLocationOpenSea {
			stopFishing = Timeline.Do(func() uint64 {
				inventory, err := s.Inventory().Fetch()
				if err != nil {
					slog.Error("Error fetching inventory", "err", err)
					return timeline.Day

				}

				err = inventory.AddItem(models.Item{
					Name:   string(models.ItemTypeFish),
					Amount: 1,
				})

				if err != nil {
					slog.Error("Error adding item", "err", err)
				}

				return timeline.Day
			}, timeline.Day)
		}
	})
}
