import * as ex from "excalibur";
import { TILE_SIZE } from "./App";
export class Ship {
  actor: ex.Actor;
  constructor() {
    this.actor = new ex.Actor({
      width: TILE_SIZE / 8,
      height: TILE_SIZE / 8,

      // Let's give it some color with one of the predefined
      // color constants
      color: ex.Color.Chartreuse,
    });

    this.actor.body.collisionType = ex.CollisionType.Fixed;
  }

  setTarget(x: number, y: number) {
    const shrinkSpeed = 1;
    const moveSpeed = 500;
    const shrinkFactor = 1.1;

    this.actor.actions.scaleTo(
      new ex.Vector(shrinkFactor, shrinkFactor),
      new ex.Vector(shrinkSpeed, shrinkSpeed)
    );
    this.actor.actions.moveTo(x, y, moveSpeed);
    this.actor.actions.scaleTo(
      new ex.Vector(1, 1),
      new ex.Vector(shrinkSpeed, shrinkSpeed)
    );
  }
}
