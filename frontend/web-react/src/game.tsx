import {
  JSONRPC,
  Port,
  Inventory,
  Snapshot,
  Continent,
  SetupRPC,
  Coordinate,
  CrewInformationReq,
  TimeChangedReq,
  PlotRouteResp,
  ShipChangedReq,
  ShipMovedReq,
} from "@sea/shared";
import * as ex from "excalibur";
import { useState, useEffect } from "react";
import { Ship } from "./ship";

export const TILE_SIZE = 32;

export let rpc: JSONRPC | undefined = undefined;

export const useStart = ({
  cursorSquareChanged,
}: {
  cursorSquareChanged?: (x: number, y: number) => void;
}) => {
  // const [isDocked, setIsDocked] = useState(false)

  const [dockedAtPort, setDockedAtPort] = useState<Port>();
  const [inventory, setInventory] = useState<Inventory>();
  const [crew, setCrew] = useState<CrewInformationReq>();
  const [route, setRoute] = useState<{
    plottedRoute: PlotRouteResp;
    shipMovedUpdate?: ShipMovedReq;
    active: boolean;
  }>();
  const [setSail, setSetSail] = useState<() => void>();
  const [playerShipInfo, setPlayerShipInfo] = useState<ShipChangedReq>();
  const [cameraRotation, setCameraRotation] = useState(0);
  const [timeInformation, setTimeInformation] = useState<TimeChangedReq>();

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
          members: continent.points.map((p) => {
            const dark = ex.Color.fromHex("#563010");
            const light = ex.Color.fromHex("#B39253");

            const c = ex.Color.fromRGB(
              dark.r + (light.r - dark.r) * p.elevation,
              dark.g + (light.g - dark.g) * p.elevation,
              dark.b + (light.b - dark.b) * p.elevation
            );

            return {
              graphic: new ex.Rectangle({
                width: TILE_SIZE,
                height: TILE_SIZE,
                color: c, // ex.Color.fromHex("#B39253"). // ex.Color.fromHex("#563010") // ex.Color.fromRGB(0, 180 - 100 * p.elevation, 0),
              }),
              offset: new ex.Vector(
                p.x * TILE_SIZE - TILE_SIZE / 2,
                p.y * TILE_SIZE - TILE_SIZE / 2
              ),
            };
          }),
        });

        const actor = new ex.Actor({
          z: 0,
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

        ship.actor.pos.x = loginResp.ship.x * TILE_SIZE;
        ship.actor.pos.y = loginResp.ship.y * TILE_SIZE;

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

          actor.graphics.use(
            new ex.Circle({
              radius: TILE_SIZE / 2,
              color: ex.Color.fromRGB(200, 50, 70, 1),
            })
          );

          game.add(actor);

          const portName = new ex.Label({
            text: port.name,
            x: port.point.X * TILE_SIZE - TILE_SIZE,
            y: port.point.Y * TILE_SIZE - TILE_SIZE,
            z: 2,
            color: ex.Color.White,
            font: new ex.Font({
              size: 20,
              family: "Arial",
            }),
          });

          game.add(portName);
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

      const plotter = new ex.Actor({
        z: 3,
      });

      game.add(plotter);

      //   const lastXY = { x: 0, y: 0 };

      //   // need to debounce plotting
      //   let last = Date.now();

      game.input.pointers.primary.on("move", (evt) => {
        // const x = Math.round(evt.worldPos.x / TILE_SIZE);
        // const y = Math.round(evt.worldPos.y / TILE_SIZE);

        // if (lastXY.x === x && lastXY.y === y) {
        //   return;
        // }

        // lastXY.x = x;
        // lastXY.y = y;

        targetActor.pos.x = Math.round(evt.worldPos.x / TILE_SIZE) * TILE_SIZE;
        targetActor.pos.y = Math.round(evt.worldPos.y / TILE_SIZE) * TILE_SIZE;

        // last = Date.now();

        // setTimeout(() => {
        //   if (last + 50 > Date.now()) {
        //     return;
        //   }
        // }, 51);

        cursorSquareChanged?.(
          Math.round(evt.coordinates.worldPos.x / TILE_SIZE),
          Math.round(evt.coordinates.worldPos.y / TILE_SIZE)
        );
      });

      let potentialRoute:
        | {
            path: Coordinate[];
            selectedXY: Coordinate;
          }
        | undefined = undefined;
      // let currentRoute

      game.input.pointers.primary.on("down", (evt) => {
        const x = Math.round(evt.worldPos.x / TILE_SIZE);
        const y = Math.round(evt.worldPos.y / TILE_SIZE);

        if (
          !potentialRoute ||
          x !== potentialRoute.selectedXY.x ||
          y !== potentialRoute.selectedXY.y
        ) {
          rpc
            ?.send("PlotRoute", {
              coordinate: {
                x,
                y,
              },
            })
            .then((resp) => {
              const { coordinates } = resp;
              if (!coordinates || coordinates.length === 0) {
                plotter.graphics.use(new ex.GraphicsGroup({ members: [] }));
                return;
              }

              const graphicsGroup = new ex.GraphicsGroup({
                useAnchor: false,
                members: coordinates.map((p) => {
                  return {
                    graphic: new ex.Rectangle({
                      width: TILE_SIZE / 4,
                      height: TILE_SIZE / 4,
                      color: ex.Color.Black, // ex.Color.fromHex("#B39253"). // ex.Color.fromHex("#563010") // ex.Color.fromRGB(0, 180 - 100 * p.elevation, 0),
                    }),
                    offset: new ex.Vector(p.x * TILE_SIZE, p.y * TILE_SIZE),
                  };
                }),
              });
              plotter.graphics.use(graphicsGroup);

              potentialRoute = {
                path: coordinates,
                selectedXY: { x, y },
              };

              setRoute({
                plottedRoute: resp,
                shipMovedUpdate: undefined,
                active: false,
              });

              setSetSail(() => {
                return () => {
                  potentialRoute = undefined;
                  rpc!.send("MoveShip", {
                    x,
                    y,
                  });
                  // plotter.graphics.use(new ex.GraphicsGroup({ members: [] }));
                  setRoute((r) => {
                    if (!r) {
                      return r;
                    }
                    return {
                      ...r,
                      active: true,
                    };
                  });
                };
              });
            });
        }
      });

      game.input.pointers.primary.on("wheel", (evt) => {
        game.currentScene.camera.zoom += evt.deltaY / 5000;
      });

      let paused = false;

      game.input.keyboard.on("press", (evt) => {
        if (evt.key === ex.Input.Keys.Space) {
          if (paused) {
            rpc?.send("ControlTime", {
              resume: true,
            });
          } else {
            rpc?.send("ControlTime", {
              pause: true,
            });
          }

          paused = !paused;
        } else if (evt.key === ex.Keys.Q) {
          game.currentScene.camera.rotation += Math.PI / 8;
          setCameraRotation(game.currentScene.camera.rotation);
        } else if (evt.key === ex.Keys.E) {
          game.currentScene.camera.rotation -= Math.PI / 8;
          setCameraRotation(game.currentScene.camera.rotation);
        }
      });

      rpc?.receive("ShipMoved", (req) => {
        const { location } = req;
        ship.actor.actions.clearActions();
        ship.setTarget(location.x * TILE_SIZE, location.y * TILE_SIZE);

        setRoute((r) => {
          if (!r) {
            return r;
          }

          return {
            ...r,
            shipMovedUpdate: req,
          };
        });

        if (req.route_info.is_cancelled) {
          plotter.graphics.use(new ex.GraphicsGroup({ members: [] }));
        } else {
          const graphicsGroup = new ex.GraphicsGroup({
            useAnchor: false,
            members: req.route_info.trajectory.map((p) => {
              return {
                graphic: new ex.Rectangle({
                  width: TILE_SIZE / 4,
                  height: TILE_SIZE / 4,
                  color: ex.Color.Blue, // ex.Color.fromHex("#B39253"). // ex.Color.fromHex("#563010") // ex.Color.fromRGB(0, 180 - 100 * p.elevation, 0),
                }),
                offset: new ex.Vector(p.x * TILE_SIZE, p.y * TILE_SIZE),
              };
            }),
          });
          plotter.graphics.use(graphicsGroup);
        }

        return Promise.resolve({
          result: {},
        });
      });

      rpc?.receive("TimeChanged", (req) => {
        //timeChanged?.(req);

        setTimeInformation(req);

        return Promise.resolve({
          result: {},
        });
      });

      rpc?.receive("ShipDocked", ({ undocked, port }) => {
        console.log("Docked", undocked, port);

        if (undocked) {
          setDockedAtPort(undefined);

          ship.actor.graphics.visible = true;
        } else {
          console.log("Docked");

          ship.actor.graphics.visible = false;

          setDockedAtPort(port);
        }

        return Promise.resolve({
          result: {},
        });
      });

      rpc?.receive("ShipInventoryChanged", ({ inventory }) => {
        setInventory(inventory);

        return Promise.resolve({
          result: {},
        });
      });

      game.backgroundColor = ex.Color.fromHex("#E0D08A");

      game.start();

      game.currentScene.camera.strategy.lockToActor(ship.actor);
    })();
  }, []);

  rpc?.receive("CrewInformation", (crewInfo) => {
    console.log("Crew info", crewInfo);

    setCrew(crewInfo);

    return Promise.resolve({
      result: {},
    });
  });

  rpc?.receive("ShipChanged", (req) => {
    setPlayerShipInfo(req);

    return Promise.resolve({
      result: {},
    });
  });

  return {
    dockedAtPort,
    inventory,
    crew,
    route,
    setSail,
    playerShipInfo,
    cameraRotation,
    timeInformation,
  };
};
