import { Snapshot, Continent, SetupRPC, JSONRPC, ShipInventoryChangedReq } from "@sea/shared";
import * as ex from "excalibur";
import "./App.css";
import { Ship } from "./ship";
import { useEffect, useState } from "react";
import { Button } from "./ui/button";
import { Tablet } from "./ui/tablet";
import { Copy } from "./ui/copy";

type InventoryItem = ShipInventoryChangedReq["items"][0];

export const TILE_SIZE = 32;

let rpc: JSONRPC | undefined = undefined;

const useStart = ({
  timeChanged,
  inventoryChanged
} : {
  timeChanged?: (current_tick: number, ticks_per_second: number) => void;
  inventoryChanged?: (items: InventoryItem[]) => void;
}) => {

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

        const graphicsGroup = new ex.GraphicsGroup({
          useAnchor: false,
          members: continent.points.map(p => {

            const dark = ex.Color.fromHex("#563010")
            const light = ex.Color.fromHex("#B39253")

            const c = ex.Color.fromRGB(
              dark.r + (light.r - dark.r) * p.elevation,
              dark.g + (light.g - dark.g) * p.elevation,
              dark.b + (light.b - dark.b) * p.elevation
            )
            
            return {
              graphic: new ex.Rectangle({
                width: TILE_SIZE,
                height: TILE_SIZE,
                color: c // ex.Color.fromHex("#B39253"). // ex.Color.fromHex("#563010") // ex.Color.fromRGB(0, 180 - 100 * p.elevation, 0),
              }),
              offset: new ex.Vector((p.x) * TILE_SIZE - TILE_SIZE / 2, (p.y) * TILE_SIZE - TILE_SIZE / 2)
            }
          })
        })

          
        const actor = new ex.Actor({
          z: 0
        });

        actor.graphics.use(graphicsGroup);

        // game.toggleDebug();

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
            x: port.point.X * TILE_SIZE,
            y: port.point.Y * TILE_SIZE,
            color: ex.Color.Red,
            z: 1,
          });

          actor.graphics.use(new ex.Circle({
            radius: TILE_SIZE / 2,
            color: ex.Color.fromRGB(200, 50, 70, 1)
          }))
    
          game.add(actor);
        }
      };
    
      await ready;
    
      const game = new ex.Engine({
        displayMode: ex.DisplayMode.FitScreenAndFill,
        canvasElementId: "game",
      });
    
      game.add(ship.actor);

      const targetActor = new ex.Actor({
        x: 0,
        y: 0,
        z: 2,
      });


      // square outline
      targetActor.graphics.use(
        new ex.Rectangle({
          width: TILE_SIZE,
          height: TILE_SIZE,
          color: ex.Color.Transparent,
          lineWidth: 3,
          strokeColor: ex.Color.Black,
        })
      );

      game.add(targetActor);
    
      game.input.pointers.primary.on("move", (evt) => {
        targetActor.pos.x = Math.round(evt.worldPos.x / TILE_SIZE) * TILE_SIZE;
        targetActor.pos.y = Math.round(evt.worldPos.y / TILE_SIZE) * TILE_SIZE;
      })
      

      game.input.pointers.primary.on("down", (evt) => {
        rpc!.send("MoveShip", {
          x: Math.round(evt.coordinates.worldPos.x / TILE_SIZE),
          y: Math.round(evt.coordinates.worldPos.y / TILE_SIZE),
        });
      });
    
      game.input.pointers.primary.on("wheel", (evt) => {
        game.currentScene.camera.zoom += evt.deltaY / 5000;
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

        timeChanged?.(current_tick, ticks_per_second);

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

      rpc?.receive("ShipInventoryChanged", ({items}) => {
        inventoryChanged?.(items);

        return Promise.resolve({
          result: {}
        })
      } )

      

      game.backgroundColor = ex.Color.fromHex("#E0D08A");

      game.start();
    })()
      
    }, []);


    return {
      isDocked
    }
}

export const useDraggable = () => {
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [dragging, setDragging] = useState<{ offsetX: number; offsetY: number }>();

  return {
    onMouseMove: (evt: React.MouseEvent) => {
      if (dragging) {
        setPosition({
          x: evt.clientX + dragging.offsetX,
          y: evt.clientY + dragging.offsetY,
        });
      }
    },
    onMouseDown: (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => setDragging({
      offsetX: e.currentTarget.offsetLeft - e.clientX,
      offsetY: e.currentTarget.offsetTop - e.clientY,
    }),
    onMouseUp: () => setDragging(undefined),
    onMouseLeave: () => setDragging(undefined),
    position
  }
}

const useInventory = () => {
  const [open, setOpen] = useState(true);

  const [items, setItems] = useState<InventoryItem[]>([]);

  return {
    inventoryWidget: open ? 
      <Tablet classNames="fixed" > 
      <Copy>Inventory</Copy>
      <ul>
        {items.map((item) => (
          <li key={item.id}>
            {item.name}: {item.amount}
          </li>
        ))}
      </ul>
      </Tablet>: null,
      setOpen,
      setItems,
  };
}


function App() {

  const [currentTime, setCurrentTime] = useState(0);

  const { inventoryWidget, setItems} = useInventory();

  const { isDocked } = useStart({
    timeChanged: (current_tick, ticks_per_second) => {
      console.log("Time changed", current_tick, ticks_per_second);
      setCurrentTime(current_tick);
    },
    inventoryChanged: (items) => {
      console.log("Inventory changed", items);
      setItems(items);
    }
  });

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
        {inventoryWidget}
        <div className="flex justify-between">
          <Tablet>
            <div className="flex flex-col">
              <span>Ticks per second: {ticksPerSecond}</span>
              <Copy>Current time: {currentTime}</Copy>
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
          <Tablet classNames="fixed">
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
