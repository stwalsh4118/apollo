import type { LessonSummary } from "../api";

interface LessonItemProps {
  lesson: LessonSummary;
  isActive: boolean;
  onSelect: (lessonId: string) => void;
}

export default function LessonItem({
  lesson,
  isActive,
  onSelect,
}: LessonItemProps) {
  return (
    <button
      onClick={() => onSelect(lesson.id)}
      className={`flex w-full items-center gap-2 rounded px-3 py-1.5 text-left text-sm transition-colors ${
        isActive
          ? "bg-blue-50 font-medium text-blue-700"
          : "text-gray-600 hover:bg-gray-100 hover:text-gray-900"
      }`}
    >
      <span className="size-1.5 shrink-0 rounded-full bg-current opacity-40" />
      <span className="truncate">{lesson.title}</span>
    </button>
  );
}
