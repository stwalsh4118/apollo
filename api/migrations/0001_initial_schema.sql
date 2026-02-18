CREATE TABLE IF NOT EXISTS topics (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  difficulty TEXT CHECK (difficulty IN ('foundational', 'intermediate', 'advanced')),
  estimated_hours REAL,
  tags TEXT CHECK (tags IS NULL OR json_valid(tags)),
  status TEXT NOT NULL CHECK (status IN ('researching', 'draft', 'published', 'outdated')),
  version INTEGER NOT NULL DEFAULT 1,
  source_urls TEXT CHECK (source_urls IS NULL OR json_valid(source_urls)),
  generated_at TEXT,
  generated_by TEXT,
  parent_topic_id TEXT REFERENCES topics(id),
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS modules (
  id TEXT PRIMARY KEY,
  topic_id TEXT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  description TEXT,
  learning_objectives TEXT CHECK (learning_objectives IS NULL OR json_valid(learning_objectives)),
  estimated_minutes INTEGER,
  sort_order INTEGER NOT NULL,
  assessment TEXT CHECK (assessment IS NULL OR json_valid(assessment))
);

CREATE TABLE IF NOT EXISTS lessons (
  id TEXT PRIMARY KEY,
  module_id TEXT NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  sort_order INTEGER NOT NULL,
  estimated_minutes INTEGER,
  content TEXT NOT NULL CHECK (json_valid(content)),
  examples TEXT CHECK (examples IS NULL OR json_valid(examples)),
  exercises TEXT CHECK (exercises IS NULL OR json_valid(exercises)),
  review_questions TEXT CHECK (review_questions IS NULL OR json_valid(review_questions))
);

CREATE TABLE IF NOT EXISTS concepts (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  definition TEXT NOT NULL,
  defined_in_lesson TEXT REFERENCES lessons(id) ON DELETE SET NULL,
  defined_in_topic TEXT REFERENCES topics(id) ON DELETE SET NULL,
  difficulty TEXT CHECK (difficulty IN ('foundational', 'intermediate', 'advanced')),
  flashcard_front TEXT,
  flashcard_back TEXT,
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'unresolved', 'conflict')),
  aliases TEXT CHECK (aliases IS NULL OR json_valid(aliases))
);

CREATE TABLE IF NOT EXISTS concept_references (
  concept_id TEXT NOT NULL REFERENCES concepts(id) ON DELETE CASCADE,
  lesson_id TEXT NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
  context TEXT,
  PRIMARY KEY (concept_id, lesson_id)
);

CREATE TABLE IF NOT EXISTS topic_prerequisites (
  topic_id TEXT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
  prerequisite_topic_id TEXT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
  priority TEXT NOT NULL CHECK (priority IN ('essential', 'helpful', 'deep_background')),
  reason TEXT,
  PRIMARY KEY (topic_id, prerequisite_topic_id),
  CHECK (topic_id <> prerequisite_topic_id)
);

CREATE TABLE IF NOT EXISTS topic_relations (
  topic_a TEXT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
  topic_b TEXT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
  relation_type TEXT NOT NULL CHECK (relation_type IN ('related', 'builds_on', 'contrasts_with', 'subset_of')),
  description TEXT,
  PRIMARY KEY (topic_a, topic_b),
  CHECK (topic_a <> topic_b)
);

CREATE TABLE IF NOT EXISTS expansion_queue (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  topic_id TEXT NOT NULL,
  requested_by_topic TEXT REFERENCES topics(id) ON DELETE SET NULL,
  priority TEXT NOT NULL CHECK (priority IN ('essential', 'helpful', 'deep_background')),
  reason TEXT,
  status TEXT NOT NULL DEFAULT 'available' CHECK (status IN ('available', 'queued', 'researching', 'completed', 'skipped')),
  depth_from_root INTEGER,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS research_jobs (
  id TEXT PRIMARY KEY,
  root_topic TEXT,
  current_topic TEXT,
  status TEXT NOT NULL CHECK (status IN ('queued', 'researching', 'resolving', 'published', 'failed', 'cancelled')),
  progress TEXT CHECK (progress IS NULL OR json_valid(progress)),
  error TEXT,
  started_at TEXT,
  completed_at TEXT
);

CREATE TABLE IF NOT EXISTS learning_progress (
  lesson_id TEXT PRIMARY KEY REFERENCES lessons(id) ON DELETE CASCADE,
  status TEXT NOT NULL DEFAULT 'not_started' CHECK (status IN ('not_started', 'in_progress', 'completed')),
  started_at TEXT,
  completed_at TEXT,
  notes TEXT
);

CREATE TABLE IF NOT EXISTS concept_retention (
  concept_id TEXT PRIMARY KEY REFERENCES concepts(id) ON DELETE CASCADE,
  status TEXT NOT NULL DEFAULT 'new' CHECK (status IN ('new', 'learning', 'reviewing', 'mastered')),
  next_review TEXT,
  review_count INTEGER NOT NULL DEFAULT 0,
  ease_factor REAL NOT NULL DEFAULT 2.5,
  interval_days INTEGER NOT NULL DEFAULT 0,
  last_reviewed TEXT,
  last_rating TEXT CHECK (last_rating IS NULL OR last_rating IN ('forgot', 'hard', 'good', 'easy'))
);

CREATE VIRTUAL TABLE IF NOT EXISTS search_index USING fts5(
  entity_type,
  entity_id UNINDEXED,
  title,
  body
);

CREATE INDEX IF NOT EXISTS idx_modules_topic_id ON modules(topic_id);
CREATE INDEX IF NOT EXISTS idx_lessons_module_id ON lessons(module_id);
CREATE INDEX IF NOT EXISTS idx_concepts_defined_in_topic ON concepts(defined_in_topic);
CREATE INDEX IF NOT EXISTS idx_concept_references_lesson_id ON concept_references(lesson_id);
CREATE INDEX IF NOT EXISTS idx_topic_prerequisites_prereq ON topic_prerequisites(prerequisite_topic_id);
CREATE INDEX IF NOT EXISTS idx_topic_relations_topic_b ON topic_relations(topic_b);
CREATE INDEX IF NOT EXISTS idx_expansion_queue_status ON expansion_queue(status);
CREATE INDEX IF NOT EXISTS idx_expansion_queue_topic_id ON expansion_queue(topic_id);
CREATE INDEX IF NOT EXISTS idx_research_jobs_status ON research_jobs(status);
CREATE INDEX IF NOT EXISTS idx_learning_progress_status ON learning_progress(status);
CREATE INDEX IF NOT EXISTS idx_concept_retention_next_review ON concept_retention(next_review);
