export const Button = ({
  onClick,
  children,
  variant,
}: {
  onClick: () => void;
  children: React.ReactNode;
  variant?: "primary" | "secondary" | "ghost";
}) => {
  switch (variant) {
    case "primary":
      return (
        <button
          className="bg-orange-600 grow-0 w-fit px-3 py-1 border-2 text-sm border-orange-600 text-slate-50 rounded-lg"
          onClick={onClick}
        >
          {children}
        </button>
      );
    case "secondary":
      return (
        <button
          className="bg-slate-800 grow-0 w-fit px-3 py-1 border-2 text-sm border-slate-400 text-slate-50 rounded-lg"
          onClick={onClick}
        >
          {children}
        </button>
      );
    case "ghost":
      return (
        <button
          className="bg-transparent grow-0 w-fit px-3 py-1 border-2 border-transparent text-orange-100 rounded-lg"
          onClick={onClick}
        >
          {children}
        </button>
      );
  }

  return (
    <button
      className="bg-orange-700 grow-0 w-fit px-3 py-1 border-2 text-sm border-orange-600 hover:border-orange-600 text-orange-50 rounded-lg"
      onClick={onClick}
    >
      {children}
    </button>
  );
};
