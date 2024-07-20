import { Ship } from "./ship";
import "./style.css";
import * as ex from "excalibur";
import { SetupRPC, GameContext, Continent } from "@sea/shared";
import { sortPoints } from "./utils/clockwise-sort";

const setGameContext = (ctx: GameContext) => {
  localStorage.setItem("game_context", JSON.stringify(ctx));
};

const getGameContext = (): GameContext | undefined => {
  const ctx = localStorage.getItem("game_context");
  if (!ctx) {
    return undefined;
  }
  return JSON.parse(ctx);
};

function drawContinent(continent: Continent) {
  var pointCenter = { x: 0, y: 0 };
  for (const point of continent.coastal_points) {
    pointCenter.x += point.x;
    pointCenter.y += point.y;
  }

  pointCenter.x /= continent.coastal_points.length;
  pointCenter.y /= continent.coastal_points.length;

  const multiplier = 10;

  const triangle = new ex.Polygon({
    points: sortPoints(continent.coastal_points, pointCenter).flatMap((p) => [
      new ex.Vector(p.x * multiplier, p.y * multiplier),
    ]),
    color: ex.Color.Green,
  });

  const actor = new ex.Actor({
    x: pointCenter.x * multiplier,
    y: pointCenter.y * multiplier,
  });

  actor.graphics.add(triangle);

  game.add(actor);
}

const conn = new WebSocket("ws://localhost:8080/ws");
const rpc = SetupRPC(conn);

const ship = new Ship();

conn.onopen = async () => {
  console.log("Connected to server");

  var ctx = getGameContext();

  if (!ctx) {
    const registerResp = await rpc.send("Register", {});

    ctx = registerResp.context;
    setGameContext(ctx);
  }

  const loginResp = await rpc.send("Login", {
    context: ctx,
  });

  ship.actor.pos.x = loginResp.ship.x;
  ship.actor.pos.y = loginResp.ship.y;

  console.log("Login response", loginResp);

  const worldMapResp = await rpc.send("GetWorldMap", {});

  console.log("World map response", worldMapResp);

  for (const continent of worldMapResp.continents) {
    console.log("Continent", continent);
    drawContinent(continent);
  }
};

const game = new ex.Engine({
  displayMode: ex.DisplayMode.FitScreenAndFill,
});

game.add(ship.actor);

game.input.pointers.primary.on("down", (evt) => {
  rpc.send("MoveShip", {
    x: evt.coordinates.worldPos.x,
    y: evt.coordinates.worldPos.y,
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

rpc.receive("ShipChangedTarget", ({ x, y }) => {
  console.log("Received ShipChangedTarget:", x, y);

  ship.setTarget(x, y);

  return Promise.resolve({
    result: {},
  });
});

game.start();
