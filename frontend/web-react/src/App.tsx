import { Port, Inventory, ShipChangedReq, Person } from "@sea/shared";
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

const HiringMenu = ({ port }: { port: Port }) => {
  const [recruits, setRecruits] = useState<
    {
      selected: boolean;
      person: Person;
    }[]
  >();

  useEffect(() => {
    if (!port?.id || !rpc) return;

    rpc
      ?.send("GetHirablePeopleAtPort", {
        port_id: port.id,
      })
      .then((resp) => {
        setRecruits(
          resp.people.map((p) => ({
            selected: false,
            person: p,
          }))
        );
      });
  }, []);

  return (
    <Tablet>
      {!recruits && <>Loading</>}
      {recruits && (
        <div>
          <div className="flex flex-col gap-2">
            {recruits.map((r, i) => {
              return (
                <div className="flex justify-between gap-4">
                  <div className="flex gap-2">
                    <div>{upperCaseFirstLetter(r.person.first_name)}</div>
                    <div>{upperCaseFirstLetter(r.person.last_name)}</div>
                    <div>the {upperCaseFirstLetter(r.person.nick_name)}</div>

                    <div>
                      of {upperCaseFirstLetter(r.person.place_of_residence)}
                    </div>

                    <div>(age {r.person.age})</div>
                  </div>
                  <Button
                    onClick={() => {
                      setRecruits((prev) => {
                        if (!prev) return prev;
                        prev[i].selected = !prev[i].selected;
                        return [...prev];
                      });
                    }}
                  >
                    {r.selected ? "-" : "+"}
                  </Button>
                </div>
              );
            })}
          </div>
          <Button
            onClick={() => {
              rpc
                ?.send("HireCrew", {
                  people: recruits
                    .filter((r) => r.selected)
                    .map((r) => r.person),
                })
                .then(() => {
                  setRecruits((prev) => prev?.filter((r) => !r.selected));
                });
            }}
          >
            Hire
          </Button>
        </div>
      )}
    </Tablet>
  );
};

export const upperCaseFirstLetter = (s: string) => {
  return s[0].toUpperCase() + s.substring(1);
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
  const [hiringMenuOpen, setHiringMenuOpen] = useState(false);

  return (
    <>
      <Tablet>
        <div className="flex flex-col gap-2">
          <span>Port #{port.id}</span>
          <div className="flex flex-col gap-2">
            <Button onClick={() => setTradeOpen(!tradeOpen)}>Trade</Button>
            <Button
              onClick={() => {
                setHiringMenuOpen((prev) => !prev);
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
      {hiringMenuOpen && <HiringMenu port={port} />}
    </>
  );
};

function App() {
  const [cursorLocation, setCursorLocation] = useState({ x: 0, y: 0 });

  const {
    dockedAtPort,
    inventory,
    crew,
    route,
    setSail,
    playerShipInfo,
    cameraRotation,
    timeInformation,
  } = useStart({
    cursorSquareChanged(x, y) {
      setCursorLocation({ x, y });
    },
  });

  return (
    <>
      <div className="absolute w-screen h-screen  top-0 left-0 z-50 pointer-events-none">
        <div className="text-3xl text-gray-800">Ships Colonies Commerce</div>
        <div className="flex justify-between">
          <div className="flex p-1 gap-1">
            {playerShipInfo && <PlayerShip ship={playerShipInfo} />}
            {crew && <CrewWidget crew={crew} />}
            <InventoryWidget inventory={inventory} />
          </div>
          <div className="flex flex-col p-1 gap-1">
            <Tablet>
              <div className="flex flex-col">
                <span>
                  Ticks per second: {timeInformation?.ticks_per_second}
                </span>
                <Copy>Current time: {timeInformation?.current_tick}</Copy>
                <Copy>Current Day: {timeInformation?.current_day}</Copy>
                <Copy>Current Year: {timeInformation?.current_year}</Copy>
                <Copy>
                  Cursor: {cursorLocation.x}, {cursorLocation.y}
                </Copy>
                <div className="flex gap-1">
                  <Button
                    onClick={() => {
                      rpc?.send("ControlTime", {
                        pause: !timeInformation?.is_paused,
                        resume: timeInformation?.is_paused,
                      });
                    }}
                  >
                    {timeInformation?.is_paused ? "|>" : "||"}
                  </Button>
                  {[1, 3, 6, 9].map((speed) => (
                    <div>
                      <Button
                        disabled={speed === timeInformation?.ticks_per_second}
                        onClick={() =>
                          rpc?.send("ControlTime", {
                            set_ticks_per_second_to: speed,
                          })
                        }
                      >
                        {speed}x
                      </Button>
                    </div>
                  ))}
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
