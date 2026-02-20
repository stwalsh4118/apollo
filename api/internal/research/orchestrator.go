package research

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/sean/apollo/api/internal/config"
	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/repository"
	"github.com/sean/apollo/api/internal/schema"
)

// maxRetries is the number of times a failed pass is retried before giving up.
const maxRetries = 1

// Filenames for embedded assets written to each job's work directory.
const (
	embeddedSystemPromptFile = "research.md"
	embeddedCurriculumSchema = "curriculum.json"
)

// pollInterval is the delay between checking for queued jobs.
const pollInterval = 2 * time.Second

// passPrompts contains the instruction sent to the CLI for each pass.
var passPrompts = map[int]string{
	1: "Survey this topic: identify the key areas, plan modules and lessons, and outline the curriculum structure. Focus on breadth — cover the full landscape before going deep.",
	2: "Deep dive: generate detailed lesson content for every module and lesson you planned. Include thorough explanations, key concepts with definitions, and flashcards for each concept.",
	3: "Generate exercises, practice problems, and review questions for every lesson. Include worked examples with explanations. Ensure exercises match the lesson difficulty.",
	4: "Final validation pass: review all content for accuracy, completeness, and consistency. Read through the file tree and fix any issues by rewriting individual files.",
}

// Orchestrator drives the research pipeline from queued job to published curriculum.
type Orchestrator struct {
	cli     CLIRunner
	pool    *PoolSummaryBuilder
	ingest  *CurriculumIngester
	repo    repository.ResearchJobRepository
	logger  zerolog.Logger
	cfg     config.Config
	mu      sync.Mutex
	cancels map[string]context.CancelFunc
}

// NewOrchestrator creates an Orchestrator with all required dependencies.
func NewOrchestrator(
	cli CLIRunner,
	pool *PoolSummaryBuilder,
	ingest *CurriculumIngester,
	repo repository.ResearchJobRepository,
	logger zerolog.Logger,
	cfg config.Config,
) *Orchestrator {
	return &Orchestrator{
		cli:     cli,
		pool:    pool,
		ingest:  ingest,
		repo:    repo,
		logger:  logger,
		cfg:     cfg,
		cancels: make(map[string]context.CancelFunc),
	}
}

// Start runs a background loop that picks up queued jobs and processes them.
// It blocks until ctx is cancelled. For PBI 6, jobs are processed one at a time.
func (o *Orchestrator) Start(ctx context.Context) {
	o.logger.Info().Msg("research orchestrator started")

	for {
		select {
		case <-ctx.Done():
			o.logger.Info().Msg("research orchestrator stopping")
			return
		default:
		}

		jobID, err := o.findQueuedJob(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			o.logger.Error().Err(err).Msg("find queued job failed")
			o.sleep(ctx)

			continue
		}

		if jobID == "" {
			o.sleep(ctx)
			continue
		}

		if err := o.RunJob(ctx, jobID); err != nil {
			o.logger.Error().Err(err).Str("job_id", jobID).Msg("job execution failed")
		}
	}
}

// findQueuedJob returns the ID of the oldest queued job, or "" if none.
func (o *Orchestrator) findQueuedJob(ctx context.Context) (string, error) {
	return o.repo.FindOldestByStatus(ctx, models.ResearchStatusQueued)
}

// sleep waits for the poll interval or until ctx is cancelled.
func (o *Orchestrator) sleep(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-time.After(pollInterval):
	}
}

// RunJob executes the full 4-pass research pipeline for a single job.
func (o *Orchestrator) RunJob(ctx context.Context, jobID string) error {
	log := o.logger.With().Str("job_id", jobID).Logger()

	job, err := o.repo.GetJobByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("get job: %w", err)
	}

	// Create a cancellable context for this job.
	jobCtx, cancel := context.WithCancel(ctx)
	o.mu.Lock()
	o.cancels[jobID] = cancel
	o.mu.Unlock()

	defer func() {
		cancel()
		o.mu.Lock()
		delete(o.cancels, jobID)
		o.mu.Unlock()
	}()

	// Transition to researching.
	if err := o.repo.UpdateJobStatus(jobCtx, jobID, models.ResearchStatusResearching, ""); err != nil {
		return fmt.Errorf("update status to researching: %w", err)
	}

	log.Info().Str("topic", job.RootTopic).Msg("starting research pipeline")

	// Prepare working directory.
	workDir, err := o.prepareWorkDir(jobCtx, jobID)
	if err != nil {
		return o.failJob(ctx, jobID, fmt.Errorf("prepare work dir: %w", err))
	}

	// Build the initial prompt. Note: job.CurrentTopic equals RootTopic for new jobs.
	// The Brief field from CreateResearchJobInput is not yet stored in the DB (see 6-2).
	topicPrompt := buildTopicPrompt(job.RootTopic, job.CurrentTopic)

	// Pass 1: Survey.
	sessionID, err := o.runPass(jobCtx, jobID, 1, topicPrompt, "", workDir, log)
	if err != nil {
		if jobCtx.Err() != nil {
			return o.handleCancellation(jobID, log)
		}

		return o.failJob(ctx, jobID, fmt.Errorf("pass 1: %w", err))
	}

	// Pass 2: Deep Dive.
	_, err = o.runPass(jobCtx, jobID, 2, passPrompts[2], sessionID, workDir, log)
	if err != nil {
		if jobCtx.Err() != nil {
			return o.handleCancellation(jobID, log)
		}

		return o.failJob(ctx, jobID, fmt.Errorf("pass 2: %w", err))
	}

	// Pass 3: Exercises.
	_, err = o.runPass(jobCtx, jobID, 3, passPrompts[3], sessionID, workDir, log)
	if err != nil {
		if jobCtx.Err() != nil {
			return o.handleCancellation(jobID, log)
		}

		return o.failJob(ctx, jobID, fmt.Errorf("pass 3: %w", err))
	}

	// Pass 4: Validation (quality review of the file tree).
	_, err = o.runPass(jobCtx, jobID, 4, passPrompts[4], sessionID, workDir, log)
	if err != nil {
		if jobCtx.Err() != nil {
			return o.handleCancellation(jobID, log)
		}

		return o.failJob(ctx, jobID, fmt.Errorf("pass 4: %w", err))
	}

	// Transition to resolving. Use parent ctx to avoid cancellation race after pass 4.
	if err := o.repo.UpdateJobStatus(ctx, jobID, models.ResearchStatusResolving, ""); err != nil {
		return o.failJob(ctx, jobID, fmt.Errorf("update status to resolving: %w", err))
	}

	// Assemble the file tree into a CurriculumOutput.
	curriculum, err := AssembleFromDir(workDir)
	if err != nil {
		return o.failJob(ctx, jobID, fmt.Errorf("assemble curriculum: %w", err))
	}

	assembledJSON, err := json.Marshal(curriculum)
	if err != nil {
		return o.failJob(ctx, jobID, fmt.Errorf("marshal assembled curriculum: %w", err))
	}

	// Ingest the assembled curriculum.
	if err := o.ingest.Ingest(jobCtx, json.RawMessage(assembledJSON)); err != nil {
		if jobCtx.Err() != nil {
			return o.handleCancellation(jobID, log)
		}

		return o.failJob(ctx, jobID, fmt.Errorf("ingest curriculum: %w", err))
	}

	// Transition to published.
	if err := o.repo.UpdateJobStatus(ctx, jobID, models.ResearchStatusPublished, ""); err != nil {
		return fmt.Errorf("update status to published: %w", err)
	}

	log.Info().Msg("research pipeline completed successfully")

	return nil
}

// Cancel cancels a running job by invoking its context cancel function.
func (o *Orchestrator) Cancel(jobID string) {
	o.mu.Lock()
	cancel, ok := o.cancels[jobID]
	o.mu.Unlock()

	if ok {
		cancel()
	}
}

// prepareWorkDir creates the working directory for a job and writes the
// knowledge pool summary, system prompt, and curriculum schema.
// The system prompt and schema are embedded at compile time and written
// to the work directory so the CLI can reference them via relative paths.
func (o *Orchestrator) prepareWorkDir(ctx context.Context, jobID string) (string, error) {
	workDir := filepath.Join(o.cfg.ResearchWorkDir, jobID)
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		return "", fmt.Errorf("create work dir: %w", err)
	}

	if err := o.pool.WriteToDir(ctx, workDir); err != nil {
		return "", fmt.Errorf("write pool summary: %w", err)
	}

	// Write the embedded system prompt.
	if err := os.WriteFile(filepath.Join(workDir, embeddedSystemPromptFile), systemPromptContent, 0o644); err != nil {
		return "", fmt.Errorf("write system prompt: %w", err)
	}

	// Write the embedded curriculum schema.
	schemaBytes, err := schema.CurriculumSchemaJSON()
	if err != nil {
		return "", fmt.Errorf("read embedded curriculum schema: %w", err)
	}

	if err := os.WriteFile(filepath.Join(workDir, embeddedCurriculumSchema), schemaBytes, 0o644); err != nil {
		return "", fmt.Errorf("write curriculum schema: %w", err)
	}

	return workDir, nil
}

// buildTopicPrompt constructs the initial prompt for Pass 1.
func buildTopicPrompt(rootTopic, brief string) string {
	prompt := fmt.Sprintf("Research the topic: %s", rootTopic)
	if brief != "" && brief != rootTopic {
		prompt += fmt.Sprintf("\n\nAdditional context: %s", brief)
	}

	prompt += "\n\n" + passPrompts[1]

	return prompt
}

// runPass executes a single CLI pass with retry logic.
// For pass 1, sessionID is empty and an initial pass is executed.
// For passes 2-4, sessionID is provided and a resume pass is executed.
// Returns the session ID from the response.
func (o *Orchestrator) runPass(ctx context.Context, jobID string, passNum int, prompt, sessionID, workDir string, log zerolog.Logger) (string, error) {
	log.Info().Int("pass", passNum).Msg("starting pass")

	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			log.Warn().Int("pass", passNum).Int("attempt", attempt+1).Msg("retrying pass")
		}

		var resp *models.CLIResponse
		var err error

		if sessionID == "" {
			resp, err = o.cli.RunInitialPass(ctx, InitialPassOpts{
				Prompt:           prompt,
				WorkDir:          workDir,
				SystemPromptFile: embeddedSystemPromptFile,
				Model:            config.DefaultResearchModel,
				AllowedTools:     config.ResearchAllowedTools(),
			})
		} else {
			resp, err = o.cli.RunResumePass(ctx, ResumePassOpts{
				Prompt:    prompt,
				SessionID: sessionID,
				WorkDir:   workDir,
			})
		}

		if err != nil {
			lastErr = err

			if ctx.Err() != nil {
				break
			}

			continue
		}

		// Update progress after successful pass.
		if err := o.updateProgress(ctx, jobID, passNum); err != nil {
			log.Warn().Err(err).Int("pass", passNum).Msg("failed to update progress")
		}

		log.Info().Int("pass", passNum).Str("session_id", resp.SessionID).Msg("pass completed")

		return resp.SessionID, nil
	}

	return "", fmt.Errorf("pass %d failed after %d attempts: %w", passNum, maxRetries+1, lastErr)
}

// updateProgress writes the current pipeline progress to the job record.
// Note: ModulesPlanned, ModulesCompleted, and ConceptsFound are not populated
// because the CLI response does not include structured intermediate counts.
// These will be populated after ingestion in a future enhancement.
func (o *Orchestrator) updateProgress(ctx context.Context, jobID string, passNum int) error {
	progress := models.ResearchProgress{
		CurrentPass:      passNum,
		TotalPasses:      config.ResearchPassCount,
		PassDescriptions: config.ResearchPassDescription,
	}

	return o.repo.UpdateJobProgress(ctx, jobID, progress)
}

// failJob marks a job as failed with the given error and returns the error.
func (o *Orchestrator) failJob(ctx context.Context, jobID string, err error) error {
	if updateErr := o.repo.UpdateJobStatus(ctx, jobID, models.ResearchStatusFailed, err.Error()); updateErr != nil {
		o.logger.Error().Err(updateErr).Str("job_id", jobID).Msg("failed to update job status to failed")
	}

	return err
}

// handleCancellation logs the cancellation and ensures the job status is set.
// The cancel endpoint already sets the status to cancelled, so this is a safety net.
// Uses context.Background() because the parent context may also be cancelled.
func (o *Orchestrator) handleCancellation(jobID string, log zerolog.Logger) error {
	log.Info().Msg("job cancelled")

	bgCtx := context.Background()

	job, err := o.repo.GetJobByID(bgCtx, jobID)
	if err != nil {
		return fmt.Errorf("check cancelled job status: %w", err)
	}

	if job.Status != models.ResearchStatusCancelled {
		if err := o.repo.UpdateJobStatus(bgCtx, jobID, models.ResearchStatusCancelled, ""); err != nil {
			return fmt.Errorf("update status to cancelled: %w", err)
		}
	}

	return nil
}
