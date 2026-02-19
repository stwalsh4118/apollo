import type {
  TopicSummary,
  TopicDetail,
  TopicFull,
  ModuleDetail,
  LessonDetail,
  TopicProgress,
  ProgressSummary,
  LessonProgress,
  UpdateProgressInput,
} from "./types";

export class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

async function fetchJson<T>(url: string, init?: RequestInit): Promise<T> {
  const res = await fetch(url, init);
  if (!res.ok) {
    throw new ApiError(res.status, `HTTP ${res.status}: ${url}`);
  }
  return res.json();
}

async function mutateJson<T>(
  url: string,
  method: string,
  body: unknown,
): Promise<T> {
  return fetchJson(url, {
    method,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
}

export function fetchTopics(): Promise<TopicSummary[]> {
  return fetchJson("/api/topics");
}

export function fetchTopicDetail(id: string): Promise<TopicDetail> {
  return fetchJson(`/api/topics/${encodeURIComponent(id)}`);
}

export function fetchTopicFull(id: string): Promise<TopicFull> {
  return fetchJson(`/api/topics/${encodeURIComponent(id)}/full`);
}

export function fetchModuleDetail(id: string): Promise<ModuleDetail> {
  return fetchJson(`/api/modules/${encodeURIComponent(id)}`);
}

export function fetchLessonDetail(id: string): Promise<LessonDetail> {
  return fetchJson(`/api/lessons/${encodeURIComponent(id)}`);
}

export function fetchTopicProgress(topicId: string): Promise<TopicProgress> {
  return fetchJson(`/api/progress/topics/${encodeURIComponent(topicId)}`);
}

export function fetchProgressSummary(): Promise<ProgressSummary> {
  return fetchJson("/api/progress/summary");
}

export function updateLessonProgress(
  lessonId: string,
  input: UpdateProgressInput,
): Promise<LessonProgress> {
  return mutateJson(
    `/api/progress/lessons/${encodeURIComponent(lessonId)}`,
    "PUT",
    input,
  );
}
