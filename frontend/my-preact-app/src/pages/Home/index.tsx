import * as ex from "excalibur";
import preactLogo from "../../assets/preact.svg";
import "./style.css";
import { Snapshot, Continent, SetupRPC } from "@sea/shared";
import { Ship } from "../../../../web-react/src/ship";

/*
<div class="home">
			<a href="https://preactjs.com" target="_blank">
				<img src={preactLogo} alt="Preact logo" height="160" width="160" />
			</a>
			<h1>Get Started building Vite-powered Preact Apps </h1>
			<section>
				<Resource
					title="Learn Preact"
					description="If you're new to Preact, try the interactive tutorial to learn important concepts"
					href="https://preactjs.com/tutorial"
				/>
				<Resource
					title="Differences to React"
					description="If you're coming from React, you may want to   check out our docs to see where Preact differs"
					href="https://preactjs.com/guide/v10/differences-to-react"
				/>
				<Resource
					title="Learn Vite"
					description="To learn more about Vite and how you can customize it to fit your needs, take a look at their excellent documentation"
					href="https://vitejs.dev"
				/>
			</section>
		</div>
*/
export const TILE_SIZE = 32;

// export function togglePause(button: HTMLButtonElement) {
//   return async () => {
//     console.log("Toggling pause");

//     if (isPaused) {
//       rpc.send("ControlTime", {
//         set_ticks_per_second_by: 5,
//       });
//     } else {
//       rpc.send("ControlTime", {
//         set_ticks_per_second_to: 0,
//       });
//     }

//     isPaused = !isPaused;
//     button.innerText = isPaused ? "Pause" : "Resume";
//   };
// }

function abc() {
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

  // var isPaused = false;

  // const pauseButton = document.getElementById("pauseButton");
  // if (pauseButton) {
  //   pauseButton.onclick = togglePause(pauseButton as HTMLButtonElement);
  // }

  // class TestButton extends HTMLElement {
  //   text = "";

  //   constructor() {
  //     super();
  //   }

  //   connectedCallback() {
  //     this.text = this.getAttribute("text") || "Not set";

  //     this.innerHTML = `<button>${this.text}</button>`;
  //   }

  // render() {
  //   return `<button>${this.text}</button>`;
  // }
  //}

  // customElements.define("test-button", TestButton);

  game.start();
}

export function Home() {
  abc();
  return <canvas id="game"></canvas>;
}

function Resource(props) {
  return (
    <a href={props.href} target="_blank" class="resource">
      <h2>{props.title}</h2>
      <p>{props.description}</p>
    </a>
  );
}
