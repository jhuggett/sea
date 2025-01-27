import { CrewInformationReq } from "@sea/shared";
import { Copy } from "./ui/copy";
import { Tablet } from "./ui/tablet";
import { round } from "./player-ship";
import { upperCaseFirstLetter } from "./App";
import { useState } from "react";
import { Button } from "./ui/button";

export const CrewWidget = ({ crew }: { crew: CrewInformationReq }) => {
  const [listCrew, setListCrew] = useState(false);

  return (
    <Tablet>
      <div>
        <div className="flex flex-col gap-4">
          <Copy>Crew count: {crew.size} people</Copy>
          <Copy>
            Recommend manning: {crew.minimumSafeManning}-
            {crew.maximumSafeManning} people
          </Copy>
          <Copy>Morale {round(crew.morale * 100, 0)}%</Copy>
        </div>
        <Button onClick={() => setListCrew((prev) => !prev)}>
          {listCrew ? "Hide Members" : "Show Members"}
        </Button>
        {listCrew && (
          <div className="flex flex-col gap-2 max-h-96 overflow-auto">
            {crew.crew_members.map((member) => {
              console.log({ member });

              return (
                <div className="bg-orange-800 p-2">
                  <div className="flex gap-2">
                    <div>{upperCaseFirstLetter(member.person.first_name)}</div>
                    <div>
                      {" "}
                      "{upperCaseFirstLetter(member.person.nick_name)}"
                    </div>
                    <div> {upperCaseFirstLetter(member.person.last_name)}</div>
                  </div>
                  <div>{member.contract.title}</div>
                  <div>Age: {member.person.age}</div>
                  <div>Morale: {round(member.person.morale, 2)}%</div>
                  <div>
                    Started on {member.contract.start_date}, ends{" "}
                    {member.contract.end_date || "indefinitely"}
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </Tablet>
  );
};
