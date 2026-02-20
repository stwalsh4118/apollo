package research

import (
	_ "embed"
)

// systemPromptContent holds the research skill prompt, embedded at compile time.
// The canonical copy lives at skills/research.md (project root). This copy is
// kept under prompts/ for go:embed access from within the api module.
//
//go:embed prompts/research.md
var systemPromptContent []byte
