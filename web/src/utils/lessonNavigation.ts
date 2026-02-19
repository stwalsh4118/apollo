import type { ModuleFull } from "../api";

export interface LessonNavTarget {
  id: string;
  title: string;
}

export interface LessonNav {
  prev: LessonNavTarget | null;
  next: LessonNavTarget | null;
}

export function getLessonNavigation(
  modules: ModuleFull[],
  currentLessonId: string,
): LessonNav {
  const flatLessons: LessonNavTarget[] = [];

  for (const mod of modules) {
    for (const lesson of mod.lessons) {
      flatLessons.push({ id: lesson.id, title: lesson.title });
    }
  }

  const currentIndex = flatLessons.findIndex((l) => l.id === currentLessonId);

  if (currentIndex === -1) {
    return { prev: null, next: null };
  }

  return {
    prev: currentIndex > 0 ? flatLessons[currentIndex - 1] : null,
    next:
      currentIndex < flatLessons.length - 1
        ? flatLessons[currentIndex + 1]
        : null,
  };
}
