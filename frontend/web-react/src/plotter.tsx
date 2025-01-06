import { PlotRouteResp, ShipMovedReq } from "@sea/shared";
import { Tablet } from "./ui/tablet";
import { Copy } from "./ui/copy";
import { Button } from "./ui/button";
import { round } from "./player-ship";
import { rpc } from "./game";

export const Plotter = ({
  route,
  setSail,
}: {
  route: {
    plottedRoute: PlotRouteResp;
    shipMovedUpdate?: ShipMovedReq;
    active: boolean;
  };
  setSail: () => void;
}) => {
  if (route.shipMovedUpdate?.route_info.reached_destination) {
    return <></>;
  }

  if (route.shipMovedUpdate?.route_info.is_cancelled) {
    return <></>;
  }

  if (route.active) {
    return (
      <Tablet>
        <Copy>
          On route to{" "}
          {
            route.plottedRoute.coordinates[
              route.plottedRoute.coordinates.length - 1
            ].x
          }
          ,
          {
            route.plottedRoute.coordinates[
              route.plottedRoute.coordinates.length - 1
            ].y
          }
        </Copy>
        {route.shipMovedUpdate && (
          <>
            <Copy>
              Progress{" "}
              {round(
                (route.shipMovedUpdate?.route_info.total_tiles_moved /
                  route.shipMovedUpdate?.route_info.tiles_in_route) *
                  100,
                0
              )}
              %
            </Copy>
            <Copy>
              ETA: {route.shipMovedUpdate?.route_info.estimated_time_left} days
            </Copy>
          </>
        )}
        {route.shipMovedUpdate &&
          route.shipMovedUpdate.route_info.is_paused && (
            <Button
              onClick={() => {
                rpc?.send("ManageRoute", {
                  ship_id: route.shipMovedUpdate!.ship_id,
                  action: "start",
                });
              }}
            >
              Resume
            </Button>
          )}
        {route.shipMovedUpdate &&
          !route.shipMovedUpdate.route_info.is_paused && (
            <Button
              onClick={() => {
                rpc?.send("ManageRoute", {
                  ship_id: route.shipMovedUpdate!.ship_id,
                  action: "pause",
                });
              }}
            >
              Pause
            </Button>
          )}
        {route.shipMovedUpdate && (
          <Button
            onClick={() => {
              rpc?.send("ManageRoute", {
                ship_id: route.shipMovedUpdate!.ship_id,
                action: "cancel",
              });
            }}
          >
            Cancel
          </Button>
        )}
      </Tablet>
    );
  }

  return (
    <Tablet>
      <div className="flex gap-2">
        <Copy>{route.plottedRoute.coordinates.length} squares</Copy>
        <Copy>by ~{round(route.plottedRoute.speed, 2)} squares per day</Copy>
        <Copy>
          will likely take {round(route.plottedRoute.duration, 2)} days
        </Copy>
      </div>
      <Button
        onClick={() => {
          setSail();
        }}
      >
        Set Sail
      </Button>
    </Tablet>
  );
};
