// --- Pagination ---

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  per_page: number;
}

// --- Topics ---

export interface TopicSummary {
  id: string;
  title: string;
  description?: string;
  difficulty?: string;
  estimated_hours?: number;
  tags?: string[];
  status: string;
  module_count: number;
}

interface TopicBase {
  id: string;
  title: string;
  description?: string;
  difficulty?: string;
  estimated_hours?: number;
  tags?: string[];
  status: string;
  version: number;
  source_urls?: string[];
  generated_at?: string;
  generated_by?: string;
  parent_topic_id?: string;
  created_at: string;
  updated_at: string;
}

export interface TopicDetail extends TopicBase {
  modules: ModuleSummary[];
}

export interface TopicFull extends TopicBase {
  modules: ModuleFull[];
}

// --- Modules ---

export interface ModuleSummary {
  id: string;
  title: string;
  description?: string;
  estimated_minutes?: number;
  sort_order: number;
}

interface ModuleBase {
  id: string;
  topic_id: string;
  title: string;
  description?: string;
  learning_objectives?: string[];
  estimated_minutes?: number;
  sort_order: number;
  assessment?: unknown;
}

export interface ModuleDetail extends ModuleBase {
  lessons: LessonSummary[];
}

export interface ModuleFull extends ModuleBase {
  lessons: LessonFull[];
}

// --- Lessons ---

export interface LessonSummary {
  id: string;
  title: string;
  sort_order: number;
  estimated_minutes?: number;
}

export interface LessonContent {
  sections: ContentSection[];
}

interface LessonBase {
  id: string;
  module_id: string;
  title: string;
  sort_order: number;
  estimated_minutes?: number;
  content: LessonContent;
  examples?: Example[];
  exercises?: Exercise[];
  review_questions?: ReviewQuestion[];
}

export type LessonDetail = LessonBase;

export interface LessonFull extends LessonBase {
  concepts?: ConceptSummary[];
}

// --- Content Sections (discriminated union) ---

export interface TextSection {
  type: "text";
  body: string;
}

export interface CodeSection {
  type: "code";
  language: string;
  code: string;
  title?: string;
  explanation?: string;
}

export type CalloutVariant = "prerequisite" | "warning" | "tip" | "info";

export interface CalloutSection {
  type: "callout";
  variant: CalloutVariant;
  body: string;
  concept_ref?: string;
}

export interface DiagramSection {
  type: "diagram";
  format: "mermaid" | "image";
  source: string;
  title?: string;
}

export interface TableSection {
  type: "table";
  headers: string[];
  rows: string[][];
}

export interface ImageSection {
  type: "image";
  url: string;
  alt: string;
  caption?: string;
}

export type ContentSection =
  | TextSection
  | CodeSection
  | CalloutSection
  | DiagramSection
  | TableSection
  | ImageSection;

// --- Examples ---

export interface Example {
  title: string;
  description?: string;
  code?: string;
  explanation?: string;
}

// --- Exercises ---

export type ExerciseType =
  | "command"
  | "configuration"
  | "exploration"
  | "build"
  | "troubleshooting"
  | "scenario"
  | "thought_experiment"
  | "hands_on";

export interface Exercise {
  type: ExerciseType;
  title: string;
  instructions: string;
  environment?: string;
  success_criteria?: string[];
  hints?: string[];
}

// --- Review Questions ---

export interface ReviewQuestion {
  question: string;
  answer: string;
  concepts_tested?: string[];
}

// --- Concepts ---

export interface ConceptSummary {
  id: string;
  name: string;
  definition: string;
  difficulty?: string;
  status: string;
  defined_in_topic?: string;
  aliases?: string[];
}

export interface ConceptDetail {
  id: string;
  name: string;
  definition: string;
  defined_in_lesson?: string;
  defined_in_topic?: string;
  difficulty?: string;
  flashcard_front?: string;
  flashcard_back?: string;
  status: string;
  aliases?: string[];
  references?: ConceptReference[];
}

export interface ConceptReference {
  lesson_id: string;
  lesson_title: string;
  context?: string;
}

// --- Learning Progress ---

export type ProgressStatus = "not_started" | "in_progress" | "completed";

export interface LessonProgress {
  lesson_id: string;
  lesson_title?: string;
  status: ProgressStatus;
  started_at?: string;
  completed_at?: string;
  notes?: string;
}

export interface TopicProgress {
  topic_id: string;
  lessons: LessonProgress[];
}

export interface ProgressSummary {
  total_lessons: number;
  completed_lessons: number;
  completion_percentage: number;
  active_topics: number;
}

export interface UpdateProgressInput {
  status: ProgressStatus;
  notes?: string;
}

// --- Search ---

export interface SearchResult {
  entity_type: string;
  entity_id: string;
  title: string;
  snippet: string;
}

// --- Graph ---

export interface GraphNode {
  id: string;
  label: string;
  type: string;
}

export interface GraphEdge {
  source: string;
  target: string;
  type: string;
}

export interface GraphData {
  nodes: GraphNode[];
  edges: GraphEdge[];
}
