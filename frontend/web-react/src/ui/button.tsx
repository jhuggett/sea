export const Button = ({
  onClick,
  children,
}: {
  onClick: () => void;
  children: React.ReactNode;
}) => {
  return (
    <button
      className="bg-slate-800 grow-0 w-fit px-3 py-1 border-2 text-sm border-slate-400 hover:border-orange-600 text-slate-50 rounded-lg"
      onClick={onClick}
    >
      {children}
    </button>
  );
};
