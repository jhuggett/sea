import { CrewInformationReq } from "@sea/shared";
import { Copy } from "./ui/copy";
import { Tablet } from "./ui/tablet";
import { round } from "./player-ship";

export const CrewWidget = ({ crew }: { crew: CrewInformationReq }) => {
  return (
    <Tablet>
      <div>
        <div className="flex flex-col gap-4">
          <Copy>Crew count: {crew.size} people</Copy>
          <Copy>
            Recommend manning: {crew.minimumSafeManning}-
            {crew.maximumSafeManning} people
          </Copy>
          <Copy>
            Daily Wage: {crew.wage * crew.size} piece_of_eight per day
          </Copy>
          <Copy>Daily Rations: {crew.rations * crew.size}</Copy>
          <Copy>Morale {round(crew.morale * 100, 0)}%</Copy>
        </div>
      </div>
    </Tablet>
  );
};
