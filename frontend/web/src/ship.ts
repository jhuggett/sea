import * as ex from "excalibur";
import { TILE_SIZE } from "./main";

export class Ship {
  actor: ex.Actor;
  constructor() {
    this.actor = new ex.Actor({
      width: TILE_SIZE,
      height: TILE_SIZE,

      // Let's give it some color with one of the predefined
      // color constants
      color: ex.Color.Chartreuse,
    });

    this.actor.body.collisionType = ex.CollisionType.Fixed;
  }

  setTarget(x: number, y: number) {
    this.actor.actions.moveTo(x, y, 250);
  }
}
