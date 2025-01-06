import { Port, Inventory, TimeChangedReq, ShipChangedReq } from "@sea/shared";
import "./App.css";
import React, { useEffect, useState } from "react";
import { Button } from "./ui/button";
import { Tablet } from "./ui/tablet";
import { Copy } from "./ui/copy";
import { Trade } from "./trade";
import { rpc, useStart } from "./game";
import { CrewWidget } from "./crew";
import { Plotter } from "./plotter";
import { PlayerShip } from "./player-ship";

export const useDraggable = () => {
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [dragging, setDragging] = useState<{
    offsetX: number;
    offsetY: number;
  }>();

  return {
    onMouseMove: (evt: React.MouseEvent) => {
      if (dragging) {
        setPosition({
          x: evt.clientX + dragging.offsetX,
          y: evt.clientY + dragging.offsetY,
        });
      }
    },
    onMouseDown: (e: React.MouseEvent<HTMLDivElement, MouseEvent>) =>
      setDragging({
        offsetX: e.currentTarget.offsetLeft - e.clientX,
        offsetY: e.currentTarget.offsetTop - e.clientY,
      }),
    onMouseUp: () => setDragging(undefined),
    onMouseLeave: () => setDragging(undefined),
    position,
  };
};

const InventoryWidget = ({ inventory }: { inventory?: Inventory }) => {
  return (
    <Tablet>
      <Copy>Inventory</Copy>
      <ul>
        {inventory?.items?.map((item) => (
          <li key={item.id}>
            {item.name}: {item.amount}
          </li>
        ))}
      </ul>
      <Copy>ID: {inventory?.id}</Copy>
    </Tablet>
  );
};

const PortMenu = ({
  port,
  playerInventory,
  playerShipInfo,
}: {
  port: Port;
  playerInventory: Inventory;
  playerShipInfo: ShipChangedReq;
}) => {
  const [tradeOpen, setTradeOpen] = useState(false);

  return (
    <>
      <Tablet>
        <div className="flex flex-col gap-2">
          <span>Port #{port.id}</span>
          <div className="flex flex-col gap-2">
            <Button onClick={() => setTradeOpen(!tradeOpen)}>Trade</Button>
            <Button
              onClick={() => {
                rpc?.send("HireCrew", {
                  size: 1,
                });
              }}
            >
              Hire Crew
            </Button>
            <Button
              onClick={() =>
                rpc?.send("RepairShip", { ship_id: playerShipInfo.id })
              }
            >
              Repair Ship
            </Button>
          </div>
        </div>
      </Tablet>
      {tradeOpen && <Trade port={port} playerInventory={playerInventory} />}
    </>
  );
};

function App() {
  const [currentTime, setCurrentTime] = useState<TimeChangedReq>();

  const [cursorLocation, setCursorLocation] = useState({ x: 0, y: 0 });

  const {
    dockedAtPort,
    inventory,
    crew,
    route,
    setSail,
    playerShipInfo,
    cameraRotation,
  } = useStart({
    timeChanged: (req) => {
      setCurrentTime(req);
    },
    cursorSquareChanged(x, y) {
      setCursorLocation({ x, y });
    },
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
        <div className="flex justify-between">
          <div className="flex p-1 gap-1">
            {playerShipInfo && <PlayerShip ship={playerShipInfo} />}
            {crew && <CrewWidget crew={crew} />}
            <InventoryWidget inventory={inventory} />
          </div>
          <div className="flex flex-col p-1 gap-1">
            <Tablet>
              <div className="flex flex-col">
                <span>Ticks per second: {currentTime?.ticks_per_second}</span>
                <Copy>Current time: {currentTime?.current_tick}</Copy>
                <Copy>Current Day: {currentTime?.current_day}</Copy>
                <Copy>Current Year: {currentTime?.current_year}</Copy>
                <Copy>
                  Cursor: {cursorLocation.x}, {cursorLocation.y}
                </Copy>
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
              <Copy>Debug stuffs</Copy>
              <Button onClick={() => localStorage.removeItem("game_snapshot")}>
                Clear current context
              </Button>
              {playerShipInfo && (
                <Button
                  onClick={() => {
                    rpc?.send("RepairShip", { ship_id: playerShipInfo?.id });
                  }}
                >
                  Repair ship
                </Button>
              )}
            </Tablet>
          </div>
        </div>
        <div className="flex p-1 gap-1">
          {dockedAtPort && playerShipInfo && (
            <PortMenu
              port={dockedAtPort}
              playerInventory={inventory!}
              playerShipInfo={playerShipInfo}
            />
          )}
        </div>
        <div className="bottom-0 fixed flex justify-center w-full">
          {route && setSail && (
            <>
              <Plotter route={route} setSail={setSail} />
            </>
          )}
        </div>
        <img
          className="w-64 fixed right-0 bottom-0"
          style={{
            transform: `rotate(${cameraRotation}rad)`,
          }}
          src="/another-compass-rose.png"
        />
      </div>

      <canvas id="game"></canvas>
    </>
  );
}

export default App;
