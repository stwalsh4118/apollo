package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sean/apollo/api/internal/config"
	"github.com/sean/apollo/api/internal/database"
	"github.com/sean/apollo/api/internal/logging"
)

const exitCodeFailure = 1

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "apollo startup failed: %v\n", err)
		os.Exit(exitCodeFailure)
	}
}

func run() error {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger, err := logging.New(os.Stdout, cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("create logger: %w", err)
	}

	handle, err := database.Open(ctx, cfg.DatabasePath, logger)
	if err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}
	defer func() {
		_ = handle.Close()
	}()

	logger.Info().
		Str("database_path", cfg.DatabasePath).
		Int("server_port", cfg.ServerPort).
		Int("max_research_depth", cfg.MaxResearchDepth).
		Int("max_parallel_agents", cfg.MaxParallelAgents).
		Int("topic_size_limit", cfg.TopicSizeLimit).
		Str("auto_expand_priority", cfg.AutoExpandPriority).
		Int("curriculum_stale_days", cfg.CurriculumStale).
		Int("mastery_threshold_days", cfg.MasteryThreshold).
		Str("research_work_dir", cfg.ResearchWorkDir).
		Msg("apollo foundation initialized with database migrations")

	return nil
}
