package research

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sean/apollo/api/internal/config"
	"github.com/sean/apollo/api/internal/models"
)

// CLIRunner is the interface for spawning Claude Code CLI sessions.
// This allows the orchestrator to be tested with a mock implementation.
type CLIRunner interface {
	RunInitialPass(ctx context.Context, opts InitialPassOpts) (*models.CLIResponse, error)
	RunResumePass(ctx context.Context, opts ResumePassOpts) (*models.CLIResponse, error)
}

// InitialPassOpts configures the initial (Pass 1) CLI invocation.
type InitialPassOpts struct {
	Prompt           string
	WorkDir          string
	SystemPromptFile string
	Model            string
	AllowedTools     []string
}

// ResumePassOpts configures a resume (Pass 2-4) CLI invocation.
type ResumePassOpts struct {
	Prompt         string
	SessionID      string
	WorkDir        string
	JSONSchemaFile string // set only for the final pass
}

// CLISession implements CLIRunner using os/exec to spawn the Claude Code CLI.
type CLISession struct {
	binaryPath string
}

// NewCLISession creates a CLISession with the given binary path.
func NewCLISession(binaryPath string) *CLISession {
	return &CLISession{binaryPath: binaryPath}
}

// RunInitialPass spawns the CLI for Pass 1 with the given options.
func (s *CLISession) RunInitialPass(ctx context.Context, opts InitialPassOpts) (*models.CLIResponse, error) {
	args := []string{
		"-p", opts.Prompt,
		"--output-format", config.ResearchOutputFormat,
	}

	if opts.SystemPromptFile != "" {
		args = append(args, "--system-prompt-file", opts.SystemPromptFile)
	}

	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}

	if len(opts.AllowedTools) > 0 {
		args = append(args, "--allowedTools", strings.Join(opts.AllowedTools, ","))
	}

	return s.run(ctx, opts.WorkDir, args)
}

// RunResumePass spawns the CLI for Passes 2-4 with resume.
func (s *CLISession) RunResumePass(ctx context.Context, opts ResumePassOpts) (*models.CLIResponse, error) {
	args := []string{
		"-p", opts.Prompt,
		"--resume", opts.SessionID,
		"--output-format", config.ResearchOutputFormat,
	}

	if opts.JSONSchemaFile != "" {
		args = append(args, "--json-schema", opts.JSONSchemaFile)
	}

	return s.run(ctx, opts.WorkDir, args)
}

func (s *CLISession) run(ctx context.Context, workDir string, args []string) (*models.CLIResponse, error) {
	// The CLI runs headless with no TTY to approve permission prompts.
	// --dangerously-skip-permissions is required for non-interactive use.
	args = append(args, "--dangerously-skip-permissions")

	cmd := exec.CommandContext(ctx, s.binaryPath, args...)
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), "CLAUDE_CODE_MAX_OUTPUT_TOKENS=65536")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		stderrStr := stderr.String()
		if stderrStr == "" {
			stderrStr = "(no stderr)"
		}

		return nil, fmt.Errorf("cli exited with error: %w; stderr: %s", err, stderrStr)
	}

	resp, err := parseCLIResponse(stdout.Bytes())
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// parseCLIResponse parses the JSON output from the CLI.
func parseCLIResponse(data []byte) (*models.CLIResponse, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("cli returned empty output")
	}

	var resp models.CLIResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		snippet := string(data)
		if len(snippet) > 200 {
			snippet = snippet[:200] + "..."
		}

		return nil, fmt.Errorf("parse cli response: %w; raw output: %s", err, snippet)
	}

	if resp.IsError {
		return nil, fmt.Errorf("cli reported error: %s", resp.ErrorMessage)
	}

	return &resp, nil
}

// BuildInitialArgs returns the argument list for an initial pass (useful for testing).
func BuildInitialArgs(opts InitialPassOpts) []string {
	args := []string{
		"-p", opts.Prompt,
		"--output-format", config.ResearchOutputFormat,
	}

	if opts.SystemPromptFile != "" {
		args = append(args, "--system-prompt-file", opts.SystemPromptFile)
	}

	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}

	if len(opts.AllowedTools) > 0 {
		args = append(args, "--allowedTools", strings.Join(opts.AllowedTools, ","))
	}

	return args
}

// BuildResumeArgs returns the argument list for a resume pass (useful for testing).
func BuildResumeArgs(opts ResumePassOpts) []string {
	args := []string{
		"-p", opts.Prompt,
		"--resume", opts.SessionID,
		"--output-format", config.ResearchOutputFormat,
	}

	if opts.JSONSchemaFile != "" {
		args = append(args, "--json-schema", opts.JSONSchemaFile)
	}

	return args
}

// Verify interface compliance at compile time.
var _ CLIRunner = (*CLISession)(nil)
