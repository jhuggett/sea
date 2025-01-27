import { Port, Inventory, Item } from "@sea/shared";
import { ReactNode, useState } from "react";
import { Button } from "./ui/button";
import { Copy } from "./ui/copy";
import { Tablet } from "./ui/tablet";
import { rpc } from "./game";

const TradeSquareItem = ({
  item,
  onAdd,
  onRemove,
  normalizedValue,
}: {
  item: Item;
  onAdd?: (item: Item, amount: number) => void;
  onRemove?: (item: Item, amount: number) => void;
  normalizedValue?: number;
}) => {
  const [range, setRange] = useState(0);

  return (
    <div className="flex gap-2 justify-between bg-orange-800  p-2 rounded-lg">
      <div key={item.id}>
        <Copy>
          {item.name}: {item.amount} {`{${normalizedValue}}`}
        </Copy>
        <div>
          <input
            className="range pr-6 accent-red-500"
            type="range"
            value={range}
            min="0"
            max={item.amount}
            name={`range-${item.id}`}
            onChange={(e) => {
              setRange(parseInt(e.target.value));
            }}
          ></input>
          {range}
        </div>
      </div>
      {onAdd && (
        <Button
          variant="ghost"
          onClick={() => {
            onAdd(item, range);
            setRange(0);
          }}
        >
          +
        </Button>
      )}
      {onRemove && (
        <Button
          variant="ghost"
          onClick={() => {
            onRemove(item, range);
            setRange(0);
          }}
        >
          -
        </Button>
      )}
    </div>
  );
};

export const TradeSquare = ({
  items,
  title,
  onAdd,
  onRemove,
  port,
}: {
  title: ReactNode;
  items: Item[];
  onAdd?: (item: Item, amount: number) => void;
  onRemove?: (item: Item, amount: number) => void;
  port: Port;
}) => {
  return (
    <div className=" p-4 m-4 basis-[100%] flex flex-col text-left">
      <div>{title}</div>
      <div className="flex flex-col gap-1 p-2 ">
        {items.map((item) => (
          <TradeSquareItem
            key={item.id}
            item={item}
            onAdd={onAdd}
            onRemove={onRemove}
            normalizedValue={port?.item_valuation?.[item.name]}
          />
        ))}
      </div>
    </div>
  );
};

export const Trade = ({
  port,
  playerInventory,
}: {
  port: Port;
  playerInventory: Inventory;
}) => {
  const [buying, setBuying] = useState<Item[]>([]);
  const [selling, setSelling] = useState<Item[]>([]);

  const [buyables, setBuyables] = useState<Item[]>(port.inventory?.items || []);
  const [sellables, setSellables] = useState<Item[]>(
    playerInventory?.items || []
  );

  return (
    <Tablet>
      <div className="flex flex-col gap-2">
        <div className="flex flex-col min-w-[33vw]">
          <span>Trade</span>
          <div className="flex ">
            <TradeSquare
              port={port}
              title={<Copy>Your stuff</Copy>}
              items={sellables}
              onRemove={(item, amount) => {
                setSellables((prev) => {
                  const sellableItem = prev.find((i) => i.id === item.id);

                  if (!sellableItem) {
                    return prev;
                  }

                  if (sellableItem.amount <= amount) {
                    return prev.filter((i) => i.id !== item.id);
                  }

                  return prev.map((i) =>
                    i.id === item.id ? { ...i, amount: i.amount - amount } : i
                  );
                });

                setSelling((prev) => {
                  const sellingItem = prev.find((i) => i.id === item.id);

                  if (!sellingItem) {
                    const newItem = { ...item, amount: amount };
                    return [...prev, newItem];
                  }

                  return prev.map((i) =>
                    i.id === item.id ? { ...i, amount: i.amount + amount } : i
                  );
                });
              }}
            />
            <TradeSquare
              port={port}
              title={<Copy>Their stuff (ID: {port.inventory?.id})</Copy>}
              items={buyables}
              onAdd={(item, amount) => {
                setBuyables((prev) => {
                  const buyableItem = prev.find((i) => i.id === item.id);

                  if (!buyableItem) {
                    return prev;
                  }

                  if (buyableItem.amount <= amount) {
                    return prev.filter((i) => i.id !== item.id);
                  }

                  return prev.map((i) =>
                    i.id === item.id ? { ...i, amount: i.amount - amount } : i
                  );
                });

                setBuying((prev) => {
                  const buyingItem = prev.find((i) => i.id === item.id);

                  if (!buyingItem) {
                    const newItem = { ...item, amount: amount };
                    return [...prev, newItem];
                  }

                  return prev.map((i) =>
                    i.id === item.id ? { ...i, amount: i.amount + amount } : i
                  );
                });
              }}
            />
          </div>
          <div className="flex ">
            <TradeSquare
              port={port}
              title={
                <Copy>
                  Selling{" "}
                  {selling.reduce(
                    (acc, item) =>
                      acc + item.amount * port.item_valuation![item.name],
                    0
                  )}
                </Copy>
              }
              items={selling}
              onRemove={(item, amount) => {
                setSelling((prev) => {
                  const sellingItem = prev.find((i) => i.id === item.id);

                  if (!sellingItem) {
                    return prev;
                  }

                  if (sellingItem.amount <= amount) {
                    return prev.filter((i) => i.id !== item.id);
                  }

                  return prev.map((i) =>
                    i.id === item.id ? { ...i, amount: i.amount - amount } : i
                  );
                });

                setSellables((prev) => {
                  const sellableItem = prev.find((i) => i.id === item.id);

                  if (!sellableItem) {
                    const newItem = { ...item, amount: amount };
                    return [...prev, newItem];
                  }

                  return prev.map((i) =>
                    i.id === item.id ? { ...i, amount: i.amount + amount } : i
                  );
                });
              }}
            />
            <TradeSquare
              port={port}
              title={
                <Copy>
                  Buying{" "}
                  {buying.reduce(
                    (acc, item) =>
                      acc + item.amount * port.item_valuation![item.name],
                    0
                  )}
                </Copy>
              }
              items={buying}
              onRemove={(item, amount) => {
                setBuying((prev) => {
                  const buyingItem = prev.find((i) => i.id === item.id);

                  if (!buyingItem) {
                    return prev;
                  }

                  if (buyingItem.amount <= amount) {
                    return prev.filter((i) => i.id !== item.id);
                  }

                  return prev.map((i) =>
                    i.id === item.id ? { ...i, amount: i.amount - amount } : i
                  );
                });

                setBuyables((prev) => {
                  const buyableItem = prev.find((i) => i.id === item.id);

                  if (!buyableItem) {
                    const newItem = { ...item, amount: amount };
                    return [...prev, newItem];
                  }

                  return prev.map((i) =>
                    i.id === item.id ? { ...i, amount: i.amount + amount } : i
                  );
                });
              }}
            />
          </div>
        </div>
        <Button
          disabled={
            buying.reduce(
              (acc, item) =>
                acc + item.amount * port.item_valuation![item.name],
              0
            ) >
            selling.reduce(
              (acc, item) =>
                acc + item.amount * port.item_valuation![item.name],
              0
            )
          }
          onClick={() => {
            const portInventoryID = port.inventory?.id;
            if (!portInventoryID) {
              console.log("No port ID");

              return;
            }

            const playerInventoryID = playerInventory.id;
            if (!playerInventoryID) {
              console.log("No player inventory ID");
              return;
            }

            if (!rpc) {
              console.log("No RPC");
              return;
            }

            const actions = [
              ...buying.map((item) => ({
                item: {
                  name: item.name,
                  amount: item.amount,
                },
                from: portInventoryID,
                to: playerInventoryID,
              })),
              ...selling.map((item) => ({
                item: {
                  name: item.name,
                  amount: item.amount,
                },
                from: playerInventoryID,
                to: portInventoryID,
              })),
            ];

            console.log({
              actions,
              playerInventory,
              portInventoryID,
              buying,
              selling,
            });

            rpc.send("Trade", {
              actions,
            });
          }}
        >
          Trade
        </Button>
      </div>
    </Tablet>
  );
};
