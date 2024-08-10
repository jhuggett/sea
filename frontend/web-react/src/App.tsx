import { Snapshot, Continent, SetupRPC, JSONRPC } from "@sea/shared";
import * as ex from "excalibur";
import "./App.css";
import { Ship } from "./ship";
import { useEffect, useState } from "react";

export const TILE_SIZE = 32;

let rpc: JSONRPC | undefined = undefined;

async function start() {
  console.log("STARTING");

  const setGameSnapshot = (ctx: Snapshot) => {
    localStorage.setItem("game_snapshot", JSON.stringify(ctx));
  };

  const getGameSnapshot = (): Snapshot | undefined => {
    const ctx = localStorage.getItem("game_snapshot");
    if (!ctx) {
      return undefined;
    }
    return JSON.parse(ctx);
  };

  function drawContinent(continent: Continent) {
    const triangle = new ex.Polygon({
      points: continent.coastal_points.flatMap((p) => [
        new ex.Vector(p.x * TILE_SIZE, p.y * TILE_SIZE),
      ]),
      color: ex.Color.Green,
    });

    const actor = new ex.Actor({
      x: continent.center.x * TILE_SIZE,
      y: continent.center.y * TILE_SIZE,
    });

    actor.graphics.add(triangle);

    game.add(actor);
  }

  console.log("Connecting to server");

  const conn = new WebSocket("ws://localhost:8080/ws");

  const ship = new Ship();

  let setIsReady: ((value: unknown) => void) | undefined = undefined;
  const ready = new Promise((resolve) => {
    setIsReady = resolve;
  });

  conn.onopen = async () => {
    console.log("Connected to server");

    rpc = SetupRPC(conn);

    if (!rpc) {
      throw new Error("RPC not set up");
    }

    setIsReady?.(null);

    let ctx = getGameSnapshot();

    if (!ctx) {
      const registerResp = await rpc!.send("Register", {});

      ctx = registerResp.snapshot;
      setGameSnapshot(ctx!);
    }

    const loginResp = await rpc!.send("Login", {
      snapshot: ctx,
    });

    ship.actor.pos.x = loginResp.ship.x;
    ship.actor.pos.y = loginResp.ship.y;

    console.log("Login response", loginResp);

    const worldMapResp = await rpc!.send("GetWorldMap", {});

    console.log("World map response", worldMapResp);

    for (const continent of worldMapResp.continents) {
      console.log("Continent", continent);
      if (continent) {
        drawContinent(continent);
      }
    }
  };

  await ready;

  const game = new ex.Engine({
    displayMode: ex.DisplayMode.FitScreenAndFill,
    canvasElementId: "game",
  });

  game.add(ship.actor);

  game.input.pointers.primary.on("down", (evt) => {
    rpc!.send("MoveShip", {
      x: evt.coordinates.worldPos.x / TILE_SIZE,
      y: evt.coordinates.worldPos.y / TILE_SIZE,
    });
  });

  game.input.pointers.primary.on("wheel", (evt) => {
    game.currentScene.camera.zoom += evt.deltaY / 1000;
  });

  let locked = false;

  game.input.keyboard.on("press", (evt) => {
    if (evt.key === ex.Input.Keys.Space) {
      if (locked) {
        game.currentScene.camera.clearAllStrategies();
      } else {
        game.currentScene.camera.strategy.lockToActor(ship.actor);
      }
      locked = !locked;
    }
  });

  rpc?.receive("ShipMoved", ({ location }) => {
    ship.actor.actions.clearActions();
    ship.setTarget(location.x * TILE_SIZE, location.y * TILE_SIZE);

    return Promise.resolve({
      result: {},
    });
  });

  game.start();
}

function App() {
  useEffect(() => {
    console.log("Starting game");

    start();
  }, []);

  const [ticksPerSecond, setTicksPerSecond] = useState(1);
  const [isPaused, setIsPaused] = useState(false);

  useEffect(() => {
    if (isPaused) {
      rpc?.send("ControlTime", {
        set_ticks_per_second_to: 0,
      });

      return;
    }

    rpc?.send("ControlTime", {
      set_ticks_per_second_to: ticksPerSecond,
    });
  }, [ticksPerSecond, isPaused]);

  return (
    <>
      <div className="absolute w-screen h-screen  top-0 left-0 z-50 pointer-events-none">
        <div className="flex">
          <div className="p-4 gap-4 flex bg-slate-300 pointer-events-auto">
            <span>Ticks per second: {ticksPerSecond}</span>
            <button onClick={() => setIsPaused((prev) => !prev)}>
              {isPaused ? "Resume" : "Pause"}
            </button>
            <button onClick={() => setTicksPerSecond((prev) => prev + 1)}>
              +
            </button>
            <button
              onClick={() => setTicksPerSecond((prev) => Math.max(0, prev - 1))}
            >
              -
            </button>
          </div>
        </div>
      </div>
      <canvas id="game"></canvas>
    </>
  );
}

export default App;
