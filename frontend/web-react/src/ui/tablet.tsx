// import { useDraggable } from "../App";

export const Tablet = ({
  children,
  classNames,
}: {
  children: React.ReactNode;
  classNames?: string;
}) => {
  // const dragging = useDraggable();

  return (
    <div
      className={classNames}
      // {...dragging}
      // style={{
      //   top: dragging.position.y,
      //   left: dragging.position.x,
      // }}
    >
      <div className="p-4 rounded-md border border-slate-400 gap-4 flex bg-opacity-95 bg-orange-900 text-orange-100 font-medium font-mono text-left pointer-events-auto">
        {children}
      </div>
    </div>
  );
};
