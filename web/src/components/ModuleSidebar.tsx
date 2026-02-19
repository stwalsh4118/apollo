import type { ModuleFull } from "../api";
import ModuleItem from "./ModuleItem";

interface ModuleSidebarProps {
  modules: ModuleFull[];
  activeLessonId: string;
  onSelectLesson: (lessonId: string) => void;
  onClose?: () => void;
}

export default function ModuleSidebar({
  modules,
  activeLessonId,
  onSelectLesson,
  onClose,
}: ModuleSidebarProps) {
  return (
    <nav className="flex h-full flex-col">
      <div className="flex items-center justify-between border-b border-gray-200 px-4 py-3">
        <h2 className="text-sm font-semibold text-gray-900">Modules</h2>
        {onClose && (
          <button
            onClick={onClose}
            className="rounded p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600 lg:hidden"
            aria-label="Close sidebar"
          >
            <svg
              className="size-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        )}
      </div>
      <div className="flex-1 overflow-y-auto p-2">
        <div className="space-y-1">
          {modules.map((mod, index) => (
            <ModuleItem
              key={mod.id}
              module={mod}
              activeLessonId={activeLessonId}
              onSelectLesson={onSelectLesson}
              defaultExpanded={index === 0}
            />
          ))}
        </div>
      </div>
    </nav>
  );
}
