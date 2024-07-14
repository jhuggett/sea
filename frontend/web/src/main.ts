import { Ship } from "./ship";
import "./style.css";
import * as ex from "excalibur";
import { SetupRPC, logger } from "@sea/shared";

// logger.log = console.info;

const getShipID = () => {
  const id = localStorage.getItem("ship_id");
  if (id === null) {
    return 0;
  }
  return parseInt(id);
};

const setShipID = (id: number) => {
  localStorage.setItem("ship_id", id.toString());
};

const conn = new WebSocket("ws://localhost:8080/ws");
const rpc = SetupRPC(conn);

conn.onopen = () => {
  console.log("Connected to server");

  rpc
    .send("Login", {
      ship_id: getShipID(),
    })
    .then(({ result, error }) => {
      if (error || result === undefined) {
        console.error("Failed to login:", error);
        return;
      }
      const { ship_id } = result;
      if (ship_id !== undefined) {
        setShipID(ship_id);
      }
    });
};

const game = new ex.Engine({});

const ship = new Ship();
game.add(ship.actor);

game.input.pointers.primary.on("down", (evt) => {
  rpc.send("MoveShip", {
    ship_id: getShipID(),
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
