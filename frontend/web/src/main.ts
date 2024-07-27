import { Ship } from "./ship";
import "./style.css";
import * as ex from "excalibur";
import { SetupRPC, Snapshot, Continent } from "@sea/shared";

export const TILE_SIZE = 32;

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

const conn = new WebSocket("ws://localhost:8080/ws");
const rpc = SetupRPC(conn);

const ship = new Ship();

conn.onopen = async () => {
  console.log("Connected to server");

  var ctx = getGameSnapshot();

  if (!ctx) {
    const registerResp = await rpc.send("Register", {});

    ctx = registerResp.snapshot;
    setGameSnapshot(ctx!);
  }

  const loginResp = await rpc.send("Login", {
    snapshot: ctx,
  });

  ship.actor.pos.x = loginResp.ship.x;
  ship.actor.pos.y = loginResp.ship.y;

  console.log("Login response", loginResp);

  const worldMapResp = await rpc.send("GetWorldMap", {});

  console.log("World map response", worldMapResp);

  for (const continent of worldMapResp.continents) {
    console.log("Continent", continent);
    if (continent) {
      drawContinent(continent);
    }
  }
};

const game = new ex.Engine({
  displayMode: ex.DisplayMode.FitScreenAndFill,
});

game.add(ship.actor);

game.input.pointers.primary.on("down", (evt) => {
  rpc.send("MoveShip", {
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

rpc.receive("ShipMoved", ({ location }) => {
  ship.actor.actions.clearActions();
  ship.setTarget(location.x * TILE_SIZE, location.y * TILE_SIZE);

  return Promise.resolve({
    result: {},
  });
});

game.start();
