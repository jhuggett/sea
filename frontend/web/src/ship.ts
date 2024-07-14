import * as ex from "excalibur";

export class Ship {
  actor: ex.Actor;
  constructor() {
    this.actor = new ex.Actor({
      width: 50,
      height: 20,

      // Let's give it some color with one of the predefined
      // color constants
      color: ex.Color.Chartreuse,
    });

    this.actor.body.collisionType = ex.CollisionType.Fixed;
  }

  setTarget(x: number, y: number) {
    this.actor.actions.clearActions();
    this.actor.actions.rotateTo(
      Math.atan2(y - this.actor.pos.y, x - this.actor.pos.x),
      10
    );
    this.actor.actions.moveTo(x, y, 250);
  }
}
