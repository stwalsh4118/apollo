package config

// Exported default constants for downstream consumers (orchestrator, CLI spawner).
const (
	// DefaultMaxResearchDepth is the maximum recursive prerequisite depth.
	DefaultMaxResearchDepth = 3

	// DefaultMaxParallelAgents is the default number of concurrent research agents.
	DefaultMaxParallelAgents = 3

	// DefaultResearchWorkDir is the default directory for research session files.
	DefaultResearchWorkDir = "./data/research"

	// DefaultResearchModel is the Claude model used for research sessions.
	DefaultResearchModel = "opus"
)

// Research pipeline constants that are not user-configurable.
const (
	// ResearchPassCount is the number of passes in the research pipeline.
	ResearchPassCount = 4

	// ResearchOutputFormat is the output format flag for the CLI.
	ResearchOutputFormat = "json"
)

// ResearchAllowedTools returns the tools available to the research CLI session.
// Returns a new slice each call to prevent accidental mutation.
// The agent needs file I/O, web access, search, and the ability to spawn
// sub-agents (Task) for parallel content generation.
func ResearchAllowedTools() []string {
	return []string{
		"WebSearch", "WebFetch",
		"Read", "Write",
		"Bash", "Glob", "Grep",
		"Task", "TodoWrite",
	}
}

// ResearchPassDescription maps pass number to a human-readable description.
var ResearchPassDescription = map[int]string{
	1: "Survey — topic landscape and module planning",
	2: "Deep Dive — detailed lesson content generation",
	3: "Exercises — practice problems and review questions",
	4: "Validation — structured output and quality checks",
}
