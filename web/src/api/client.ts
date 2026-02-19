import type {
  TopicSummary,
  TopicDetail,
  TopicFull,
  ModuleDetail,
  LessonDetail,
} from "./types";

export class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

async function fetchJson<T>(url: string): Promise<T> {
  const res = await fetch(url);
  if (!res.ok) {
    throw new ApiError(res.status, `HTTP ${res.status}: ${url}`);
  }
  return res.json();
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
