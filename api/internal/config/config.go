package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	envDatabasePath        = "DATABASE_PATH"
	envServerPort          = "SERVER_PORT"
	envClaudeCodePath      = "CLAUDE_CODE_PATH"
	envMaxResearchDepth    = "MAX_RESEARCH_DEPTH"
	envMaxParallelAgents   = "MAX_PARALLEL_AGENTS"
	envTopicSizeLimit      = "TOPIC_SIZE_LIMIT"
	envAutoExpandPriority  = "AUTO_EXPAND_PRIORITY"
	envCurriculumStaleDays = "CURRICULUM_STALE_DAYS"
	envMasteryThreshold    = "MASTERY_THRESHOLD_DAYS"
	envResearchWorkDir     = "RESEARCH_WORK_DIR"
	envLogLevel            = "LOG_LEVEL"
)

const (
	defaultDatabasePath       = "./data/apollo.db"
	defaultServerPort         = 8080
	defaultClaudeCodePath     = "claude"
	defaultMaxResearchDepth   = 3
	defaultMaxParallelAgents  = 3
	defaultTopicSizeLimit     = 8
	defaultAutoExpandPriority = "essential"
	defaultCurriculumStale    = 180
	defaultMasteryThreshold   = 90
	defaultResearchWorkDir    = "./data/research"
	defaultLogLevel           = "info"
)

// Config holds runtime configuration loaded from environment variables.
type Config struct {
	DatabasePath       string
	ServerPort         int
	ClaudeCodePath     string
	MaxResearchDepth   int
	MaxParallelAgents  int
	TopicSizeLimit     int
	AutoExpandPriority string
	CurriculumStale    int
	MasteryThreshold   int
	ResearchWorkDir    string
	LogLevel           string
}

// Load reads environment variables and returns an application configuration.
func Load() (Config, error) {
	serverPort, err := intEnv(envServerPort, defaultServerPort)
	if err != nil {
		return Config{}, err
	}

	maxResearchDepth, err := intEnv(envMaxResearchDepth, defaultMaxResearchDepth)
	if err != nil {
		return Config{}, err
	}

	maxParallelAgents, err := intEnv(envMaxParallelAgents, defaultMaxParallelAgents)
	if err != nil {
		return Config{}, err
	}

	topicSizeLimit, err := intEnv(envTopicSizeLimit, defaultTopicSizeLimit)
	if err != nil {
		return Config{}, err
	}

	curriculumStale, err := intEnv(envCurriculumStaleDays, defaultCurriculumStale)
	if err != nil {
		return Config{}, err
	}

	masteryThreshold, err := intEnv(envMasteryThreshold, defaultMasteryThreshold)
	if err != nil {
		return Config{}, err
	}

	return Config{
		DatabasePath:       stringEnv(envDatabasePath, defaultDatabasePath),
		ServerPort:         serverPort,
		ClaudeCodePath:     stringEnv(envClaudeCodePath, defaultClaudeCodePath),
		MaxResearchDepth:   maxResearchDepth,
		MaxParallelAgents:  maxParallelAgents,
		TopicSizeLimit:     topicSizeLimit,
		AutoExpandPriority: stringEnv(envAutoExpandPriority, defaultAutoExpandPriority),
		CurriculumStale:    curriculumStale,
		MasteryThreshold:   masteryThreshold,
		ResearchWorkDir:    stringEnv(envResearchWorkDir, defaultResearchWorkDir),
		LogLevel:           stringEnv(envLogLevel, defaultLogLevel),
	}, nil
}

func stringEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}

	return value
}

func intEnv(key string, fallback int) (int, error) {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback, nil
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}

	return parsedValue, nil
}
