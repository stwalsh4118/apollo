import type { LessonNavTarget } from "../utils/lessonNavigation";

interface LessonNavigationProps {
  prev: LessonNavTarget | null;
  next: LessonNavTarget | null;
  onNavigate: (lessonId: string) => void;
}

export default function LessonNavigation({
  prev,
  next,
  onNavigate,
}: LessonNavigationProps) {
  if (!prev && !next) {
    return null;
  }

  return (
    <nav
      className="mt-10 flex items-stretch gap-4 border-t border-gray-200 pt-6"
      aria-label="Lesson navigation"
    >
      {prev ? (
        <button
          type="button"
          onClick={() => onNavigate(prev.id)}
          className="flex flex-1 flex-col items-start rounded-lg border border-gray-200 px-4 py-3 text-left transition-colors hover:bg-gray-50"
        >
          <span className="text-xs font-medium text-gray-500">Previous</span>
          <span className="mt-1 text-sm font-medium text-gray-900">
            {prev.title}
          </span>
        </button>
      ) : (
        <div className="flex-1" />
      )}
      {next ? (
        <button
          type="button"
          onClick={() => onNavigate(next.id)}
          className="flex flex-1 flex-col items-end rounded-lg border border-gray-200 px-4 py-3 text-right transition-colors hover:bg-gray-50"
        >
          <span className="text-xs font-medium text-gray-500">Next</span>
          <span className="mt-1 text-sm font-medium text-gray-900">
            {next.title}
          </span>
        </button>
      ) : (
        <div className="flex-1" />
      )}
    </nav>
  );
}
