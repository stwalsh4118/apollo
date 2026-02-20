package research_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/rs/zerolog"

	"github.com/sean/apollo/api/internal/config"
	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/repository"
	"github.com/sean/apollo/api/internal/research"
)

// mockCLIRunner is a test double for CLIRunner that returns canned responses.
// It tracks a logical pass counter: pass 1 for RunInitialPass, then passes 2, 3, 4
// for sequential RunResumePass calls. The counter only advances after a successful call.
type mockCLIRunner struct {
	mu         sync.Mutex
	totalCalls int
	passNum    int                         // current logical pass being attempted
	responses  map[int]*models.CLIResponse // pass number -> response
	errors     map[int]int                 // pass number -> number of failures before success
	attempts   map[int]int                 // pass number -> attempts so far
}

func newMockCLI() *mockCLIRunner {
	return &mockCLIRunner{
		passNum:   0, // will be set to 1 on first call
		responses: make(map[int]*models.CLIResponse),
		errors:    make(map[int]int),
		attempts:  make(map[int]int),
	}
}

func (m *mockCLIRunner) setResponse(passNum int, resp *models.CLIResponse) {
	m.responses[passNum] = resp
}

func (m *mockCLIRunner) setFailCount(passNum int, count int) {
	m.errors[passNum] = count
}

func (m *mockCLIRunner) RunInitialPass(_ context.Context, _ research.InitialPassOpts) (*models.CLIResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalCalls++
	m.passNum = 1
	m.attempts[1]++

	if failures, ok := m.errors[1]; ok && m.attempts[1] <= failures {
		return nil, fmt.Errorf("mock CLI error for pass 1 (attempt %d)", m.attempts[1])
	}

	if resp, ok := m.responses[1]; ok {
		return resp, nil
	}

	return &models.CLIResponse{SessionID: "session-abc", Result: "pass 1 done"}, nil
}

func (m *mockCLIRunner) RunResumePass(ctx context.Context, _ research.ResumePassOpts) (*models.CLIResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalCalls++

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// If we haven't tried this pass yet, or the previous attempt succeeded
	// (meaning passNum was already advanced), advance to next pass.
	// If the previous attempt for the current pass failed, this is a retry.
	currentPass := m.passNum
	if m.attempts[currentPass+1] == 0 && !m.lastCallFailed(currentPass) {
		currentPass = m.passNum + 1
	}

	m.passNum = currentPass
	m.attempts[currentPass]++

	if failures, ok := m.errors[currentPass]; ok && m.attempts[currentPass] <= failures {
		return nil, fmt.Errorf("mock CLI error for pass %d (attempt %d)", currentPass, m.attempts[currentPass])
	}

	if resp, ok := m.responses[currentPass]; ok {
		return resp, nil
	}

	return &models.CLIResponse{SessionID: "session-abc", Result: fmt.Sprintf("pass %d done", currentPass)}, nil
}

// lastCallFailed returns true if the last attempt for the given pass failed.
func (m *mockCLIRunner) lastCallFailed(passNum int) bool {
	failures, hasFailures := m.errors[passNum]
	if !hasFailures {
		return false
	}

	return m.attempts[passNum] > 0 && m.attempts[passNum] <= failures
}

func (m *mockCLIRunner) callCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.totalCalls
}

// Helper to create the test orchestrator with real DB repo and mocked CLI.
func setupOrchestrator(t *testing.T, cli research.CLIRunner) (*research.Orchestrator, *sql.DB, repository.ResearchJobRepository) {
	t.Helper()

	db := setupTestDB(t)
	repo := repository.NewResearchJobRepository(db)
	pool := research.NewPoolSummaryBuilder(db)
	ingest := research.NewCurriculumIngester(db)
	logger := zerolog.New(os.Stderr).Level(zerolog.Disabled)

	workDir := t.TempDir()

	cfg := config.Config{
		ResearchWorkDir: workDir,
		ClaudeCodePath:  "claude",
	}

	orch := research.NewOrchestrator(cli, pool, ingest, repo, logger, cfg)

	return orch, db, repo
}

func TestOrchestratorHappyPath(t *testing.T) {
	cli := newMockCLI()

	// Pass 4 returns structured output that passes schema validation.
	cli.setResponse(4, &models.CLIResponse{
		SessionID:        "session-abc",
		Result:           "done",
		StructuredOutput: json.RawMessage(sampleCurriculum),
	})

	orch, _, repo := setupOrchestrator(t, cli)

	// Create a queued job.
	job, err := repo.CreateJob(context.Background(), models.CreateResearchJobInput{
		Topic: "Go Concurrency",
	})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	// Run the job.
	if err := orch.RunJob(context.Background(), job.ID); err != nil {
		t.Fatalf("run job: %v", err)
	}

	// Verify final status is published.
	updated, err := repo.GetJobByID(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	if updated.Status != models.ResearchStatusPublished {
		t.Fatalf("expected status 'published', got %q", updated.Status)
	}

	// Verify all 4 passes were executed (pass 1 initial + passes 2-4 resume).
	if cli.callCount() != 4 {
		t.Fatalf("expected 4 CLI calls, got %d", cli.callCount())
	}

	// Verify progress was updated.
	if len(updated.Progress) == 0 {
		t.Fatal("expected progress to be set")
	}

	var progress models.ResearchProgress
	if err := json.Unmarshal(updated.Progress, &progress); err != nil {
		t.Fatalf("unmarshal progress: %v", err)
	}

	if progress.CurrentPass != 4 {
		t.Fatalf("expected current pass 4, got %d", progress.CurrentPass)
	}

	if progress.TotalPasses != config.ResearchPassCount {
		t.Fatalf("expected total passes %d, got %d", config.ResearchPassCount, progress.TotalPasses)
	}
}

func TestOrchestratorPassRetrySuccess(t *testing.T) {
	cli := newMockCLI()

	// Pass 2 fails once then succeeds on retry.
	cli.setFailCount(2, 1)

	// Pass 4 returns structured output.
	cli.setResponse(4, &models.CLIResponse{
		SessionID:        "session-abc",
		Result:           "done",
		StructuredOutput: json.RawMessage(sampleCurriculum),
	})

	orch, _, repo := setupOrchestrator(t, cli)

	job, err := repo.CreateJob(context.Background(), models.CreateResearchJobInput{
		Topic: "Go Concurrency",
	})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	if err := orch.RunJob(context.Background(), job.ID); err != nil {
		t.Fatalf("run job: %v", err)
	}

	updated, err := repo.GetJobByID(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	if updated.Status != models.ResearchStatusPublished {
		t.Fatalf("expected status 'published', got %q", updated.Status)
	}

	// Pass 2 was called twice (fail + retry) + passes 1, 3, 4 = 5 total calls.
	if cli.callCount() != 5 {
		t.Fatalf("expected 5 CLI calls (1 retry), got %d", cli.callCount())
	}
}

func TestOrchestratorUnrecoverableFailure(t *testing.T) {
	cli := newMockCLI()

	// Pass 2 fails on all attempts (more failures than retries allow).
	cli.setFailCount(2, 3)

	orch, _, repo := setupOrchestrator(t, cli)

	job, err := repo.CreateJob(context.Background(), models.CreateResearchJobInput{
		Topic: "Go Concurrency",
	})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	err = orch.RunJob(context.Background(), job.ID)
	if err == nil {
		t.Fatal("expected error for unrecoverable failure")
	}

	updated, err := repo.GetJobByID(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	if updated.Status != models.ResearchStatusFailed {
		t.Fatalf("expected status 'failed', got %q", updated.Status)
	}

	if updated.Error == "" {
		t.Fatal("expected error message to be set")
	}
}

func TestOrchestratorCancellation(t *testing.T) {
	cli := newMockCLI()

	// Make pass 2 block until context is cancelled.
	originalRun := cli.RunResumePass
	_ = originalRun // avoid unused warning

	// We'll use a slow mock that checks context.
	slowCLI := &slowMockCLI{
		blockPass:  2,
		sessionID:  "session-abc",
		curriculum: json.RawMessage(sampleCurriculum),
	}

	orch, _, repo := setupOrchestrator(t, slowCLI)

	job, err := repo.CreateJob(context.Background(), models.CreateResearchJobInput{
		Topic: "Go Concurrency",
	})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	// Set the status to cancelled in DB (as the cancel handler would).
	if err := repo.UpdateJobStatus(context.Background(), job.ID, models.ResearchStatusCancelled, ""); err != nil {
		t.Fatalf("update status: %v", err)
	}

	// Run the job with a context that will be cancelled.
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- orch.RunJob(ctx, job.ID)
	}()

	// Wait for the slow CLI to signal it's blocking, then cancel.
	slowCLI.waitUntilBlocking()
	cancel()

	err = <-errCh

	// Cancellation returns nil (handled gracefully).
	if err != nil {
		t.Fatalf("expected nil error for cancellation, got: %v", err)
	}

	updated, err := repo.GetJobByID(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	if updated.Status != models.ResearchStatusCancelled {
		t.Fatalf("expected status 'cancelled', got %q", updated.Status)
	}
}

// slowMockCLI blocks on a specified pass until context is cancelled.
type slowMockCLI struct {
	blockPass  int
	sessionID  string
	curriculum json.RawMessage
	mu         sync.Mutex
	callNum    int
	blocking   chan struct{}
	once       sync.Once
}

func (s *slowMockCLI) init() {
	s.once.Do(func() {
		s.blocking = make(chan struct{})
	})
}

func (s *slowMockCLI) waitUntilBlocking() {
	s.init()
	<-s.blocking
}

func (s *slowMockCLI) RunInitialPass(_ context.Context, _ research.InitialPassOpts) (*models.CLIResponse, error) {
	s.init()
	s.mu.Lock()
	s.callNum++
	s.mu.Unlock()

	return &models.CLIResponse{SessionID: s.sessionID, Result: "pass 1 done"}, nil
}

func (s *slowMockCLI) RunResumePass(ctx context.Context, _ research.ResumePassOpts) (*models.CLIResponse, error) {
	s.init()
	s.mu.Lock()
	s.callNum++
	num := s.callNum
	s.mu.Unlock()

	if num == s.blockPass {
		// Signal that we're blocking.
		close(s.blocking)
		// Wait for cancellation.
		<-ctx.Done()
		return nil, ctx.Err()
	}

	if num == 4 {
		return &models.CLIResponse{
			SessionID:        s.sessionID,
			Result:           "done",
			StructuredOutput: s.curriculum,
		}, nil
	}

	return &models.CLIResponse{SessionID: s.sessionID, Result: fmt.Sprintf("pass %d done", num)}, nil
}

func TestOrchestratorSchemaValidationFailure(t *testing.T) {
	cli := newMockCLI()

	// Pass 4 returns invalid structured output (not matching curriculum schema).
	cli.setResponse(4, &models.CLIResponse{
		SessionID:        "session-abc",
		Result:           "done",
		StructuredOutput: json.RawMessage(`{"bad": "data"}`),
	})

	orch, _, repo := setupOrchestrator(t, cli)

	job, err := repo.CreateJob(context.Background(), models.CreateResearchJobInput{
		Topic: "Go Concurrency",
	})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	err = orch.RunJob(context.Background(), job.ID)
	if err == nil {
		t.Fatal("expected error for invalid structured output")
	}

	updated, err := repo.GetJobByID(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	if updated.Status != models.ResearchStatusFailed {
		t.Fatalf("expected status 'failed', got %q", updated.Status)
	}

	if updated.Error == "" {
		t.Fatal("expected error message for schema validation failure")
	}
}

func TestOrchestratorProgressTracking(t *testing.T) {
	cli := newMockCLI()

	cli.setResponse(4, &models.CLIResponse{
		SessionID:        "session-abc",
		Result:           "done",
		StructuredOutput: json.RawMessage(sampleCurriculum),
	})

	orch, _, repo := setupOrchestrator(t, cli)

	job, err := repo.CreateJob(context.Background(), models.CreateResearchJobInput{
		Topic: "Go Concurrency",
	})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	if err := orch.RunJob(context.Background(), job.ID); err != nil {
		t.Fatalf("run job: %v", err)
	}

	// The final progress should reflect pass 4.
	updated, err := repo.GetJobByID(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	var progress models.ResearchProgress
	if err := json.Unmarshal(updated.Progress, &progress); err != nil {
		t.Fatalf("unmarshal progress: %v", err)
	}

	if progress.CurrentPass != 4 {
		t.Fatalf("expected current pass 4, got %d", progress.CurrentPass)
	}

	if progress.TotalPasses != 4 {
		t.Fatalf("expected total passes 4, got %d", progress.TotalPasses)
	}

	if len(progress.PassDescriptions) != 4 {
		t.Fatalf("expected 4 pass descriptions, got %d", len(progress.PassDescriptions))
	}
}

func TestOrchestratorCancelMethod(t *testing.T) {
	cli := newMockCLI()
	orch, _, _ := setupOrchestrator(t, cli)

	// Cancel a non-existent job should not panic.
	orch.Cancel("nonexistent-id")
}

func TestOrchestratorFinalPassNoStructuredOutput(t *testing.T) {
	cli := newMockCLI()

	// Pass 4 returns response with empty structured output.
	cli.setResponse(4, &models.CLIResponse{
		SessionID: "session-abc",
		Result:    "done",
	})

	orch, _, repo := setupOrchestrator(t, cli)

	job, err := repo.CreateJob(context.Background(), models.CreateResearchJobInput{
		Topic: "Go Concurrency",
	})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	err = orch.RunJob(context.Background(), job.ID)
	if err == nil {
		t.Fatal("expected error for missing structured output")
	}

	updated, err := repo.GetJobByID(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}

	if updated.Status != models.ResearchStatusFailed {
		t.Fatalf("expected status 'failed', got %q", updated.Status)
	}
}
