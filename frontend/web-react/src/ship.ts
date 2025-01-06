import * as ex from "excalibur";
import { TILE_SIZE } from "./game";
export class Ship {
  actor: ex.Actor;
  constructor() {
    this.actor = new ex.Actor({
    });


    // this.actor.graphics.use(new ex.Circle({
    //   radius: TILE_SIZE / 8,
    //   color: ex.Color.fromRGB(0, 0, 0, 1),
    //   quality: 2,
    // }))

    const x =  (new ex.ImageSource('http://localhost:5173/another-attempt-at-a-ship.png'))



   x.load().then(() => {

    if (!x.isLoaded()) {
      console.log("Image not loaded");
      return
    }


      this.actor.graphics.use(
        new ex.Sprite({
          image: x,
         scale: new ex.Vector(.015625, .015625),
        })    
      )
    })

   

    
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
