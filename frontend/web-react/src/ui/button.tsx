export const Button = ({
  onClick,
  children,
  variant,
  disabled,
}: {
  onClick: () => void;
  children: React.ReactNode;
  variant?: "primary" | "secondary" | "ghost";
  disabled?: boolean;
}) => {
  if (disabled) {
    return (
      <button
        disabled={disabled}
        className="bg-gray-500 grow-0 w-fit px-3 py-1 border-2 text-sm border-gray-600 hover:border-gray-600 text-gray-50 rounded-lg"
        onClick={onClick}
      >
        {children}
      </button>
    );
  }

  switch (variant) {
    case "primary":
      return (
        <button
          disabled={disabled}
          className="bg-orange-600 grow-0 w-fit px-3 py-1 border-2 text-sm border-orange-600 text-slate-50 rounded-lg"
          onClick={onClick}
        >
          {children}
        </button>
      );
    case "secondary":
      return (
        <button
          disabled={disabled}
          className="bg-slate-800 grow-0 w-fit px-3 py-1 border-2 text-sm border-slate-400 text-slate-50 rounded-lg"
          onClick={onClick}
        >
          {children}
        </button>
      );
    case "ghost":
      return (
        <button
          disabled={disabled}
          className="bg-transparent grow-0 w-fit px-3 py-1 border-2 border-transparent text-orange-100 rounded-lg"
          onClick={onClick}
        >
          {children}
        </button>
      );
  }

  return (
    <button
      disabled={disabled}
      className="bg-orange-700 grow-0 w-fit px-3 py-1 border-2 text-sm border-orange-600 hover:border-orange-600 text-orange-50 rounded-lg"
      onClick={onClick}
    >
      {children}
    </button>
  );
};
