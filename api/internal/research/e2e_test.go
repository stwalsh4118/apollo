package research_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/sean/apollo/api/internal/config"
	"github.com/sean/apollo/api/internal/handler"
	"github.com/sean/apollo/api/internal/models"
	"github.com/sean/apollo/api/internal/repository"
	"github.com/sean/apollo/api/internal/research"
)

// e2eEnv bundles test dependencies for an end-to-end test.
type e2eEnv struct {
	server *httptest.Server
	orch   *research.Orchestrator
	repo   repository.ResearchJobRepository
	cli    *e2eMockCLI
}

// e2eMockCLI is a mock CLIRunner that records calls and returns canned responses.
type e2eMockCLI struct {
	mu        sync.Mutex
	calls     []e2eCall
	passNum   int
	responses map[int]*models.CLIResponse
	errors    map[int]int
	attempts  map[int]int
	blockPass int
	blocking  chan struct{}
	blockOnce sync.Once
}

type e2eCall struct {
	Method    string // "initial" or "resume"
	SessionID string
	HasSchema bool
}

func newE2ECLI() *e2eMockCLI {
	return &e2eMockCLI{
		responses: make(map[int]*models.CLIResponse),
		errors:    make(map[int]int),
		attempts:  make(map[int]int),
	}
}

func (m *e2eMockCLI) RunInitialPass(_ context.Context, opts research.InitialPassOpts) (*models.CLIResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.calls = append(m.calls, e2eCall{Method: "initial"})
	m.passNum = 1
	m.attempts[1]++

	if failures, ok := m.errors[1]; ok && m.attempts[1] <= failures {
		return nil, fmt.Errorf("mock CLI error pass 1")
	}

	if resp, ok := m.responses[1]; ok {
		return resp, nil
	}

	return &models.CLIResponse{SessionID: "e2e-session", Result: "survey done"}, nil
}

func (m *e2eMockCLI) RunResumePass(ctx context.Context, opts research.ResumePassOpts) (*models.CLIResponse, error) {
	m.mu.Lock()

	currentPass := m.passNum
	if m.attempts[currentPass+1] == 0 && !m.isLastFailed(currentPass) {
		currentPass = m.passNum + 1
	}
	m.passNum = currentPass
	m.attempts[currentPass]++

	call := e2eCall{
		Method:    "resume",
		SessionID: opts.SessionID,
		HasSchema: opts.JSONSchemaFile != "",
	}
	m.calls = append(m.calls, call)

	if m.blockPass > 0 && currentPass == m.blockPass {
		m.mu.Unlock()
		m.blockOnce.Do(func() { close(m.blocking) })
		<-ctx.Done()
		return nil, ctx.Err()
	}

	if failures, ok := m.errors[currentPass]; ok && m.attempts[currentPass] <= failures {
		m.mu.Unlock()
		return nil, fmt.Errorf("mock CLI error pass %d", currentPass)
	}

	if resp, ok := m.responses[currentPass]; ok {
		m.mu.Unlock()
		return resp, nil
	}

	m.mu.Unlock()
	return &models.CLIResponse{SessionID: "e2e-session", Result: fmt.Sprintf("pass %d done", currentPass)}, nil
}

func (m *e2eMockCLI) isLastFailed(passNum int) bool {
	failures, has := m.errors[passNum]
	if !has {
		return false
	}
	return m.attempts[passNum] > 0 && m.attempts[passNum] <= failures
}

func (m *e2eMockCLI) getCalls() []e2eCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]e2eCall, len(m.calls))
	copy(result, m.calls)
	return result
}

func setupE2E(t *testing.T, cli *e2eMockCLI) *e2eEnv {
	t.Helper()

	db := setupTestDB(t)
	researchRepo := repository.NewResearchJobRepository(db)
	topicRepo := repository.NewTopicRepository(db)
	conceptRepo := repository.NewConceptRepository(db)
	pool := research.NewPoolSummaryBuilder(db)
	ingest := research.NewCurriculumIngester(db)
	logger := zerolog.Nop()

	workDir := t.TempDir()
	cfg := config.Config{
		ResearchWorkDir: workDir,
	}

	orch := research.NewOrchestrator(cli, pool, ingest, researchRepo, logger, cfg)

	r := chi.NewRouter()

	researchHandler := handler.NewResearchHandler(researchRepo, orch.Cancel)
	researchHandler.RegisterRoutes(r)

	topicHandler := handler.NewTopicHandler(topicRepo)
	topicHandler.RegisterRoutes(r)

	conceptHandler := handler.NewConceptHandler(conceptRepo)
	conceptHandler.RegisterRoutes(r)

	ts := httptest.NewServer(r)
	t.Cleanup(ts.Close)

	return &e2eEnv{
		server: ts,
		orch:   orch,
		repo:   researchRepo,
		cli:    cli,
	}
}

// TestE2EHappyPath verifies AC1, AC2, AC3, AC4, AC5, AC6, AC7, AC10.
func TestE2EHappyPath(t *testing.T) {
	cli := newE2ECLI()

	// Pass 4 returns valid curriculum.
	cli.responses[4] = &models.CLIResponse{
		SessionID:        "e2e-session",
		Result:           "curriculum generated",
		StructuredOutput: json.RawMessage(sampleCurriculum),
	}

	env := setupE2E(t, cli)

	// AC1: POST /api/research creates a job.
	resp := doPost(t, env.server.URL+"/api/research", `{"topic":"Go Concurrency"}`)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("AC1: expected 201, got %d", resp.StatusCode)
	}

	var job models.ResearchJob
	decodeBody(t, resp, &job)

	if job.ID == "" {
		t.Fatal("AC1: expected job ID")
	}

	if job.Status != models.ResearchStatusQueued {
		t.Fatalf("AC1: expected status 'queued', got %q", job.Status)
	}

	// Trigger orchestrator to process the job.
	if err := env.orch.RunJob(context.Background(), job.ID); err != nil {
		t.Fatalf("orchestrator run: %v", err)
	}

	// AC2, AC3: Verify CLI calls — 4 total (1 initial + 3 resume).
	calls := cli.getCalls()
	if len(calls) != 4 {
		t.Fatalf("AC2/AC3: expected 4 CLI calls, got %d", len(calls))
	}

	if calls[0].Method != "initial" {
		t.Fatal("AC2: first call should be initial")
	}

	for i := 1; i < 4; i++ {
		if calls[i].Method != "resume" {
			t.Fatalf("AC3: call %d should be resume, got %q", i+1, calls[i].Method)
		}
		if calls[i].SessionID != "e2e-session" {
			t.Fatalf("AC3: call %d session_id mismatch: %q", i+1, calls[i].SessionID)
		}
	}

	// AC4: Final pass includes --json-schema.
	if !calls[3].HasSchema {
		t.Fatal("AC4: final pass should have --json-schema flag")
	}

	// AC7: GET /api/research/jobs/:id shows published with progress.
	resp = doGet(t, env.server.URL+"/api/research/jobs/"+job.ID)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("AC7: expected 200, got %d", resp.StatusCode)
	}

	var updatedJob models.ResearchJob
	decodeBody(t, resp, &updatedJob)

	if updatedJob.Status != models.ResearchStatusPublished {
		t.Fatalf("AC7: expected status 'published', got %q", updatedJob.Status)
	}

	// AC6: Progress shows 4/4.
	if len(updatedJob.Progress) == 0 {
		t.Fatal("AC6: expected progress to be set")
	}

	var progress models.ResearchProgress
	if err := json.Unmarshal(updatedJob.Progress, &progress); err != nil {
		t.Fatalf("AC6: unmarshal progress: %v", err)
	}

	if progress.CurrentPass != 4 || progress.TotalPasses != 4 {
		t.Fatalf("AC6: expected 4/4 progress, got %d/%d", progress.CurrentPass, progress.TotalPasses)
	}

	// AC5, AC10: Topic stored and browsable via API.
	// Verify via direct DB query that the topic was stored (the ingester created it).
	topicJob, err := env.repo.GetJobByID(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("AC5: get job: %v", err)
	}

	if topicJob.Status != models.ResearchStatusPublished {
		t.Fatalf("AC5: expected published, got %q", topicJob.Status)
	}

	// Verify topic is accessible via the API.
	resp = doGet(t, env.server.URL+"/api/topics/go-concurrency")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("AC10: expected 200 for topic, got %d", resp.StatusCode)
	}

	var topicDetail map[string]any
	decodeBody(t, resp, &topicDetail)

	if topicDetail["id"] != "go-concurrency" {
		t.Fatalf("AC5: expected topic id 'go-concurrency', got %v", topicDetail["id"])
	}

	if topicDetail["title"] != "Go Concurrency" {
		t.Fatalf("AC5: expected title 'Go Concurrency', got %v", topicDetail["title"])
	}

	// AC5: Verify full topic with modules and lessons via /full endpoint.
	resp = doGet(t, env.server.URL+"/api/topics/go-concurrency/full")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("AC5: expected 200 for topic full, got %d", resp.StatusCode)
	}

	var topicFull map[string]any
	decodeBody(t, resp, &topicFull)

	modules, ok := topicFull["modules"].([]any)
	if !ok || len(modules) == 0 {
		t.Fatal("AC5: expected at least 1 module in full topic")
	}

	firstMod, _ := modules[0].(map[string]any)
	lessons, ok := firstMod["lessons"].([]any)
	if !ok || len(lessons) == 0 {
		t.Fatal("AC5: expected at least 1 lesson in first module")
	}

	// AC5: Verify concepts accessible via API.
	resp = doGet(t, env.server.URL+"/api/concepts/goroutine")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("AC5: expected 200 for concept, got %d", resp.StatusCode)
	}

	var conceptDetail map[string]any
	decodeBody(t, resp, &conceptDetail)

	if conceptDetail["name"] != "Goroutine" {
		t.Fatalf("AC5: expected concept name 'Goroutine', got %v", conceptDetail["name"])
	}

	// Verify flashcard fields stored.
	if conceptDetail["flashcard_front"] == nil || conceptDetail["flashcard_front"] == "" {
		t.Fatal("AC5: expected flashcard_front to be set")
	}
}

// TestE2EJobFailure verifies AC8.
func TestE2EJobFailure(t *testing.T) {
	cli := newE2ECLI()

	// Pass 2 fails on all attempts.
	cli.errors[2] = 3

	env := setupE2E(t, cli)

	// Create job.
	resp := doPost(t, env.server.URL+"/api/research", `{"topic":"Failing Topic"}`)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var job models.ResearchJob
	decodeBody(t, resp, &job)

	// Run orchestrator (will fail).
	_ = env.orch.RunJob(context.Background(), job.ID)

	// AC8: GET shows failed with error.
	resp = doGet(t, env.server.URL+"/api/research/jobs/"+job.ID)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("AC8: expected 200, got %d", resp.StatusCode)
	}

	var updated models.ResearchJob
	decodeBody(t, resp, &updated)

	if updated.Status != models.ResearchStatusFailed {
		t.Fatalf("AC8: expected status 'failed', got %q", updated.Status)
	}

	if updated.Error == "" {
		t.Fatal("AC8: expected error message")
	}
}

// TestE2EJobCancellation verifies AC9.
func TestE2EJobCancellation(t *testing.T) {
	cli := newE2ECLI()
	cli.blockPass = 2
	cli.blocking = make(chan struct{})

	cli.responses[4] = &models.CLIResponse{
		SessionID:        "e2e-session",
		Result:           "done",
		StructuredOutput: json.RawMessage(sampleCurriculum),
	}

	env := setupE2E(t, cli)

	// Create job.
	resp := doPost(t, env.server.URL+"/api/research", `{"topic":"Cancellable Topic"}`)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var job models.ResearchJob
	decodeBody(t, resp, &job)

	// Start orchestrator in background.
	errCh := make(chan error, 1)
	go func() {
		errCh <- env.orch.RunJob(context.Background(), job.ID)
	}()

	// Wait for the CLI to block on pass 2.
	<-cli.blocking

	// AC9: Cancel the job via API.
	resp = doPost(t, env.server.URL+"/api/research/jobs/"+job.ID+"/cancel", "")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("AC9: expected 200 for cancel, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Wait for orchestrator to finish.
	<-errCh

	// Verify status is cancelled.
	resp = doGet(t, env.server.URL+"/api/research/jobs/"+job.ID)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("AC9: expected 200, got %d", resp.StatusCode)
	}

	var updated models.ResearchJob
	decodeBody(t, resp, &updated)

	if updated.Status != models.ResearchStatusCancelled {
		t.Fatalf("AC9: expected status 'cancelled', got %q", updated.Status)
	}
}

// TestE2EProgressTracking verifies AC6 and AC7 in more detail.
func TestE2EProgressTracking(t *testing.T) {
	cli := newE2ECLI()

	cli.responses[4] = &models.CLIResponse{
		SessionID:        "e2e-session",
		Result:           "done",
		StructuredOutput: json.RawMessage(sampleCurriculum),
	}

	env := setupE2E(t, cli)

	resp := doPost(t, env.server.URL+"/api/research", `{"topic":"Progress Topic"}`)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var job models.ResearchJob
	decodeBody(t, resp, &job)

	if err := env.orch.RunJob(context.Background(), job.ID); err != nil {
		t.Fatalf("run job: %v", err)
	}

	// Verify progress via API.
	resp = doGet(t, env.server.URL+"/api/research/jobs/"+job.ID)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var updated models.ResearchJob
	decodeBody(t, resp, &updated)

	var progress models.ResearchProgress
	if err := json.Unmarshal(updated.Progress, &progress); err != nil {
		t.Fatalf("unmarshal progress: %v", err)
	}

	// Verify pass descriptions are populated.
	if len(progress.PassDescriptions) != config.ResearchPassCount {
		t.Fatalf("expected %d pass descriptions, got %d", config.ResearchPassCount, len(progress.PassDescriptions))
	}

	if progress.PassDescriptions[1] == "" {
		t.Fatal("expected pass 1 description to be non-empty")
	}
}

// HTTP helpers.
func doPost(t *testing.T, url, body string) *http.Response {
	t.Helper()

	var reader *strings.Reader
	if body != "" {
		reader = strings.NewReader(body)
	} else {
		reader = strings.NewReader("{}")
	}

	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}

	return resp
}

func doGet(t *testing.T, url string) *http.Response {
	t.Helper()

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}

	return resp
}

func decodeBody(t *testing.T, resp *http.Response, v any) {
	t.Helper()
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("decode body: %v", err)
	}
}
