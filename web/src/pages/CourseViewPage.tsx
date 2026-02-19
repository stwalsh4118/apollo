import { useCallback, useMemo, useReducer, useRef, useState } from "react";
import { useParams } from "react-router";
import { useTopicFull, useTopicProgress } from "../api";
import type { ConceptSummary, ModuleFull, ProgressStatus } from "../api";
import ModuleSidebar from "../components/ModuleSidebar";
import LessonContent from "../components/LessonContent";
import LessonNavigation from "../components/LessonNavigation";
import { getLessonNavigation } from "../utils/lessonNavigation";

function getFirstLessonId(modules: ModuleFull[]): string | null {
  for (const mod of modules) {
    if (mod.lessons.length > 0) {
      return mod.lessons[0].id;
    }
  }
  return null;
}

type LessonAction =
  | { type: "select"; lessonId: string }
  | { type: "init"; lessonId: string };

function lessonReducer(
  state: string | null,
  action: LessonAction,
): string | null {
  switch (action.type) {
    case "select":
      return action.lessonId;
    case "init":
      return state ?? action.lessonId;
  }
}

export default function CourseViewPage() {
  const { id } = useParams<{ id: string }>();
  const topicId = id ?? "";
  const { data: topic, isLoading, isError, error } = useTopicFull(topicId);
  const { data: topicProgress } = useTopicProgress(topicId);
  const [activeLessonId, dispatch] = useReducer(lessonReducer, null);
  const [sidebarOpen, setSidebarOpen] = useState(false);

  const progressMap = useMemo(() => {
    const map = new Map<string, ProgressStatus>();
    if (topicProgress) {
      for (const lp of topicProgress.lessons) {
        map.set(lp.lesson_id, lp.status);
      }
    }
    return map;
  }, [topicProgress]);

  const notesMap = useMemo(() => {
    const map = new Map<string, string>();
    if (topicProgress) {
      for (const lp of topicProgress.lessons) {
        if (lp.notes) {
          map.set(lp.lesson_id, lp.notes);
        }
      }
    }
    return map;
  }, [topicProgress]);

  const conceptsMap = useMemo(() => {
    const map = new Map<string, ConceptSummary[]>();
    if (topic) {
      for (const mod of topic.modules) {
        for (const lesson of mod.lessons) {
          if (lesson.concepts && lesson.concepts.length > 0) {
            map.set(lesson.id, lesson.concepts);
          }
        }
      }
    }
    return map;
  }, [topic]);

  const sortedModules = useMemo(() => {
    if (!topic) return [];
    return [...topic.modules].sort((a, b) => a.sort_order - b.sort_order).map(
      (mod) => ({
        ...mod,
        lessons: [...mod.lessons].sort((a, b) => a.sort_order - b.sort_order),
      }),
    );
  }, [topic]);

  // Set initial lesson when data loads (idempotent — only sets if null)
  const firstLessonId = useMemo(
    () => getFirstLessonId(sortedModules),
    [sortedModules],
  );
  if (firstLessonId && activeLessonId === null) {
    dispatch({ type: "init", lessonId: firstLessonId });
  }

  const contentRef = useRef<HTMLDivElement>(null);

  const handleSelectLesson = useCallback((lessonId: string) => {
    dispatch({ type: "select", lessonId });
    setSidebarOpen(false);
  }, []);

  const handleNavigate = useCallback((lessonId: string) => {
    dispatch({ type: "select", lessonId });
    contentRef.current?.scrollTo({ top: 0, behavior: "smooth" });
  }, []);

  const lessonNav = useMemo(
    () =>
      activeLessonId
        ? getLessonNavigation(sortedModules, activeLessonId)
        : { prev: null, next: null },
    [sortedModules, activeLessonId],
  );

  if (isLoading) {
    return (
      <div className="py-12 text-center text-gray-500">
        Loading course...
      </div>
    );
  }

  if (isError) {
    return (
      <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="rounded-md bg-red-50 p-4 text-sm text-red-700">
          Failed to load topic: {error.message}
        </div>
      </div>
    );
  }

  if (!topic) {
    return null;
  }

  return (
    <div className="flex h-[calc(100vh-3.5rem)]">
      {/* Mobile sidebar toggle */}
      <button
        onClick={() => setSidebarOpen(true)}
        className="fixed bottom-4 left-4 z-30 rounded-full bg-gray-900 p-3 text-white shadow-lg lg:hidden"
        aria-label="Open module sidebar"
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
            d="M4 6h16M4 12h16M4 18h16"
          />
        </svg>
      </button>

      {/* Mobile sidebar overlay */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/30 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside
        className={`fixed top-14 bottom-0 left-0 z-50 w-72 border-r border-gray-200 bg-white transition-transform lg:relative lg:inset-auto lg:z-auto lg:translate-x-0 ${
          sidebarOpen ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        <ModuleSidebar
          modules={sortedModules}
          activeLessonId={activeLessonId ?? ""}
          onSelectLesson={handleSelectLesson}
          onClose={() => setSidebarOpen(false)}
          progressMap={progressMap}
        />
      </aside>

      {/* Main content */}
      <div ref={contentRef} className="flex-1 overflow-y-auto">
        <div className="mx-auto max-w-4xl px-6 py-8 lg:px-12">
          {activeLessonId ? (
            <>
              <LessonContent
                key={activeLessonId}
                lessonId={activeLessonId}
                topicId={topicId}
                lessonStatus={progressMap.get(activeLessonId)}
                lessonNotes={notesMap.get(activeLessonId)}
                concepts={conceptsMap.get(activeLessonId)}
              />
              <LessonNavigation
                prev={lessonNav.prev}
                next={lessonNav.next}
                onNavigate={handleNavigate}
              />
            </>
          ) : (
            <div className="py-12 text-center text-gray-500">
              Select a lesson from the sidebar to begin.
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
