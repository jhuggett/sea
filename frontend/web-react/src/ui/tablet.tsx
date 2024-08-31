export const Tablet = ({ children }: { children: React.ReactNode }) => {
  return (
    <div className="flex p-2">
      <div className="p-4 rounded-md border border-slate-400 gap-4 flex bg-opacity-75 bg-slate-200 pointer-events-auto">
        {children}
      </div>
    </div>
  );
};
