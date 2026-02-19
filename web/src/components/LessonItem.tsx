import type { LessonSummary, ProgressStatus } from "../api";

interface LessonItemProps {
  lesson: LessonSummary;
  isActive: boolean;
  onSelect: (lessonId: string) => void;
  status?: ProgressStatus;
}

function StatusIndicator({ status }: { status?: ProgressStatus }) {
  if (status === "completed") {
    return (
      <svg
        className="size-4 shrink-0 text-green-500"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={2.5}
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M5 13l4 4L19 7"
        />
      </svg>
    );
  }

  if (status === "in_progress") {
    return (
      <span className="size-2 shrink-0 rounded-full bg-blue-400" />
    );
  }

  return (
    <span className="size-2 shrink-0 rounded-full border border-gray-300" />
  );
}

export default function LessonItem({
  lesson,
  isActive,
  onSelect,
  status,
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
      <StatusIndicator status={status} />
      <span className="truncate">{lesson.title}</span>
    </button>
  );
}
