import { Snapshot, Continent, SetupRPC, JSONRPC } from "@sea/shared";
import * as ex from "excalibur";
import "./App.css";
import { Ship } from "./ship";
import { useEffect, useState } from "react";
import { Button } from "./ui/button";
import { Tablet } from "./ui/tablet";

export const TILE_SIZE = 32;

let rpc: JSONRPC | undefined = undefined;

// async function 

// let isDocked = false;

const useStart = () => {

  const [isDocked, setIsDocked] = useState(false)

  useEffect(() => {
    (async () => {
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
          console.log("No context found, registering");
    
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
    
        // render continents
        const worldMapResp = await rpc!.send("GetWorldMap", {});
    
        console.log("World map response", worldMapResp);
    
        for (const continent of worldMapResp.continents) {
          console.log("Continent", continent);
          if (continent) {
            drawContinent(continent);
          }
        }
    
        // render ports
        const portsResp = await rpc!.send("GetPorts", {});
    
        console.log("Ports response", portsResp);
    
        for (const port of portsResp.ports) {
          const actor = new ex.Actor({
            x: port.point.x * TILE_SIZE,
            y: port.point.y * TILE_SIZE,
            radius: 10,
            color: ex.Color.Red,
          });
    
          game.add(actor);
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


      rpc?.receive("TimeChanged", ({ current_tick, ticks_per_second }) => {
        console.log("Time changed", current_tick, ticks_per_second);

        return Promise.resolve({
          result: {},
        });
      })
    
      rpc?.receive("ShipDocked", ({ undocked }) => {
    
        if (undocked) {
          console.log("Undocked");
          setIsDocked(false);
            ship.actor.graphics.visible = true;
        } else {
          console.log("Docked");

          ship.actor.graphics.visible = false
          
          setIsDocked(true);
        }
    
        return Promise.resolve({
          result: {},
        });
      })
    
      game.start();
    })()
      
    }, []);


    return {
      isDocked
    }
}


function App() {
  const { isDocked } = useStart();

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
        <div className="flex justify-between">
          <Tablet>
            <div className="flex flex-col">
              <span>Ticks per second: {ticksPerSecond}</span>
              <div>
                <Button onClick={() => setIsPaused((prev) => !prev)}>
                  {isPaused ? "Resume" : "Pause"}
                </Button>
                <Button onClick={() => setTicksPerSecond((prev) => prev + 1)}>
                  +
                </Button>
                <Button
                  onClick={() =>
                    setTicksPerSecond((prev) => Math.max(0, prev - 1))
                  }
                >
                  -
                </Button>
              </div>
            </div>
          </Tablet>
          <Tablet>
            <Button onClick={() => localStorage.removeItem("game_snapshot")}>
              Clear current context
            </Button>
          </Tablet>
        </div>
        {isDocked && (
          <Tablet>
            <div className="flex flex-col">
              <span>Port</span>
            </div>
            </Tablet>
          )}
      </div>
      <canvas id="game"></canvas>
    </>
  );
}

export default App;
