import { ShipChangedReq } from "@sea/shared";
import { Tablet } from "./ui/tablet";
import { Copy } from "./ui/copy";

export const PlayerShip = ({ ship }: { ship: ShipChangedReq }) => {
  return (
    <Tablet>
      <div className="flex flex-col gap-1">
        <Copy>Ship: {ship.id}</Copy>
        <Copy>
          Location: {ship.x}, {ship.y}
        </Copy>
        <Copy>
          Speed: {round(ship.estimatedSailingSpeed, 2)} squares per day
        </Copy>
        <Copy>State of repair: {round(ship.stateOfRepair * 100, 2)}%</Copy>
        <Copy>
          {round(ship.currentCargoSpace, 2)} cargo space used of{" "}
          {ship.maxCargoSpaceCapacity}
        </Copy>
        <Copy>
          {round(ship.currentCargoWeight, 2)} cargo weight of recommended max of{" "}
          {ship.recommendedMaxCargoWeightCapacity}
        </Copy>
        <Copy>{ship.isDocked ? "Docked" : "At sea"}</Copy>
      </div>
    </Tablet>
  );
};

export const round = (num: number, places: number) => {
  return Math.round(num * Math.pow(10, places)) / Math.pow(10, places);
};
