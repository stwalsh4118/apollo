import { useState } from "react";
import type { ModuleFull } from "../api";
import LessonItem from "./LessonItem";

interface ModuleItemProps {
  module: ModuleFull;
  activeLessonId: string;
  onSelectLesson: (lessonId: string) => void;
  defaultExpanded?: boolean;
}

export default function ModuleItem({
  module,
  activeLessonId,
  onSelectLesson,
  defaultExpanded = false,
}: ModuleItemProps) {
  const [expanded, setExpanded] = useState(defaultExpanded);

  const hasActiveLesson = module.lessons.some((l) => l.id === activeLessonId);

  return (
    <div>
      <button
        onClick={() => setExpanded((prev) => !prev)}
        className="flex w-full items-center justify-between gap-2 rounded px-3 py-2 text-left text-sm font-medium text-gray-900 hover:bg-gray-100"
      >
        <span className="truncate">{module.title}</span>
        <svg
          className={`size-4 shrink-0 text-gray-400 transition-transform ${expanded ? "rotate-90" : ""}`}
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={2}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M9 5l7 7-7 7"
          />
        </svg>
      </button>

      {(expanded || hasActiveLesson) && (
        <div className="ml-2 space-y-0.5 pb-2">
          {module.lessons.map((lesson) => (
            <LessonItem
              key={lesson.id}
              lesson={lesson}
              isActive={lesson.id === activeLessonId}
              onSelect={onSelectLesson}
            />
          ))}
        </div>
      )}
    </div>
  );
}
