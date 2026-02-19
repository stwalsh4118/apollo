const DIFFICULTY_STYLES: Record<string, string> = {
  foundational:
    "bg-green-100 text-green-800",
  intermediate:
    "bg-yellow-100 text-yellow-800",
  advanced:
    "bg-red-100 text-red-800",
};

const DEFAULT_STYLE = "bg-gray-100 text-gray-800";

interface DifficultyBadgeProps {
  difficulty: string;
}

export default function DifficultyBadge({ difficulty }: DifficultyBadgeProps) {
  const style = DIFFICULTY_STYLES[difficulty] ?? DEFAULT_STYLE;

  return (
    <span
      className={`inline-block rounded-full px-2.5 py-0.5 text-xs font-medium capitalize ${style}`}
    >
      {difficulty}
    </span>
  );
}
