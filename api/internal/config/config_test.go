package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv(envDatabasePath, "")
	t.Setenv(envServerPort, "")
	t.Setenv(envClaudeCodePath, "")
	t.Setenv(envMaxResearchDepth, "")
	t.Setenv(envMaxParallelAgents, "")
	t.Setenv(envTopicSizeLimit, "")
	t.Setenv(envAutoExpandPriority, "")
	t.Setenv(envCurriculumStaleDays, "")
	t.Setenv(envMasteryThreshold, "")
	t.Setenv(envResearchWorkDir, "")
	t.Setenv(envLogLevel, "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.DatabasePath != defaultDatabasePath {
		t.Fatalf("expected DatabasePath %q, got %q", defaultDatabasePath, cfg.DatabasePath)
	}

	if cfg.ServerPort != defaultServerPort {
		t.Fatalf("expected ServerPort %d, got %d", defaultServerPort, cfg.ServerPort)
	}

	if cfg.ClaudeCodePath != defaultClaudeCodePath {
		t.Fatalf("expected ClaudeCodePath %q, got %q", defaultClaudeCodePath, cfg.ClaudeCodePath)
	}

	if cfg.MaxResearchDepth != defaultMaxResearchDepth {
		t.Fatalf("expected MaxResearchDepth %d, got %d", defaultMaxResearchDepth, cfg.MaxResearchDepth)
	}

	if cfg.MaxParallelAgents != defaultMaxParallelAgents {
		t.Fatalf("expected MaxParallelAgents %d, got %d", defaultMaxParallelAgents, cfg.MaxParallelAgents)
	}

	if cfg.TopicSizeLimit != defaultTopicSizeLimit {
		t.Fatalf("expected TopicSizeLimit %d, got %d", defaultTopicSizeLimit, cfg.TopicSizeLimit)
	}

	if cfg.AutoExpandPriority != defaultAutoExpandPriority {
		t.Fatalf("expected AutoExpandPriority %q, got %q", defaultAutoExpandPriority, cfg.AutoExpandPriority)
	}

	if cfg.CurriculumStale != defaultCurriculumStale {
		t.Fatalf("expected CurriculumStale %d, got %d", defaultCurriculumStale, cfg.CurriculumStale)
	}

	if cfg.MasteryThreshold != defaultMasteryThreshold {
		t.Fatalf("expected MasteryThreshold %d, got %d", defaultMasteryThreshold, cfg.MasteryThreshold)
	}

	if cfg.ResearchWorkDir != defaultResearchWorkDir {
		t.Fatalf("expected ResearchWorkDir %q, got %q", defaultResearchWorkDir, cfg.ResearchWorkDir)
	}

	if cfg.LogLevel != defaultLogLevel {
		t.Fatalf("expected LogLevel %q, got %q", defaultLogLevel, cfg.LogLevel)
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv(envDatabasePath, "/tmp/test.db")
	t.Setenv(envServerPort, "18080")
	t.Setenv(envClaudeCodePath, "/usr/local/bin/claude")
	t.Setenv(envMaxResearchDepth, "5")
	t.Setenv(envMaxParallelAgents, "7")
	t.Setenv(envTopicSizeLimit, "11")
	t.Setenv(envAutoExpandPriority, "helpful")
	t.Setenv(envCurriculumStaleDays, "120")
	t.Setenv(envMasteryThreshold, "30")
	t.Setenv(envResearchWorkDir, "/tmp/research")
	t.Setenv(envLogLevel, "debug")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.DatabasePath != "/tmp/test.db" {
		t.Fatalf("expected DatabasePath override, got %q", cfg.DatabasePath)
	}

	if cfg.ServerPort != 18080 {
		t.Fatalf("expected ServerPort override, got %d", cfg.ServerPort)
	}

	if cfg.ClaudeCodePath != "/usr/local/bin/claude" {
		t.Fatalf("expected ClaudeCodePath override, got %q", cfg.ClaudeCodePath)
	}

	if cfg.MaxResearchDepth != 5 {
		t.Fatalf("expected MaxResearchDepth override, got %d", cfg.MaxResearchDepth)
	}

	if cfg.MaxParallelAgents != 7 {
		t.Fatalf("expected MaxParallelAgents override, got %d", cfg.MaxParallelAgents)
	}

	if cfg.TopicSizeLimit != 11 {
		t.Fatalf("expected TopicSizeLimit override, got %d", cfg.TopicSizeLimit)
	}

	if cfg.AutoExpandPriority != "helpful" {
		t.Fatalf("expected AutoExpandPriority override, got %q", cfg.AutoExpandPriority)
	}

	if cfg.CurriculumStale != 120 {
		t.Fatalf("expected CurriculumStale override, got %d", cfg.CurriculumStale)
	}

	if cfg.MasteryThreshold != 30 {
		t.Fatalf("expected MasteryThreshold override, got %d", cfg.MasteryThreshold)
	}

	if cfg.ResearchWorkDir != "/tmp/research" {
		t.Fatalf("expected ResearchWorkDir override, got %q", cfg.ResearchWorkDir)
	}

	if cfg.LogLevel != "debug" {
		t.Fatalf("expected LogLevel override, got %q", cfg.LogLevel)
	}
}

func TestLoadInvalidInteger(t *testing.T) {
	t.Setenv(envMaxParallelAgents, "not-a-number")

	_, err := Load()
	if err == nil {
		t.Fatalf("expected Load() to fail for invalid integer environment value")
	}
}
