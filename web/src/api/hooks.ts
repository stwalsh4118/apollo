import { useQuery } from "@tanstack/react-query";
import {
  fetchTopics,
  fetchTopicDetail,
  fetchTopicFull,
  fetchModuleDetail,
  fetchLessonDetail,
} from "./client";

export const topicKeys = {
  all: ["topics"] as const,
  detail: (id: string) => [...topicKeys.all, "detail", id] as const,
  full: (id: string) => [...topicKeys.all, "full", id] as const,
};

export const moduleKeys = {
  all: ["modules"] as const,
  detail: (id: string) => [...moduleKeys.all, "detail", id] as const,
};

export const lessonKeys = {
  all: ["lessons"] as const,
  detail: (id: string) => [...lessonKeys.all, "detail", id] as const,
};

export function useTopics() {
  return useQuery({
    queryKey: topicKeys.all,
    queryFn: fetchTopics,
  });
}

export function useTopicDetail(id: string) {
  return useQuery({
    queryKey: topicKeys.detail(id),
    queryFn: () => fetchTopicDetail(id),
    enabled: id.length > 0,
  });
}

export function useTopicFull(id: string) {
  return useQuery({
    queryKey: topicKeys.full(id),
    queryFn: () => fetchTopicFull(id),
    enabled: id.length > 0,
  });
}

export function useModuleDetail(id: string) {
  return useQuery({
    queryKey: moduleKeys.detail(id),
    queryFn: () => fetchModuleDetail(id),
    enabled: id.length > 0,
  });
}

export function useLessonDetail(id: string) {
  return useQuery({
    queryKey: lessonKeys.detail(id),
    queryFn: () => fetchLessonDetail(id),
    enabled: id.length > 0,
  });
}
