#!/usr/bin/env bash
# Seed script for E2E testing of PBI 3 acceptance criteria.
# Seeds a topic with 2 modules, 3 lessons covering all 6 content types,
# exercises with hints, and review questions.

set -euo pipefail

API="http://localhost:8080/api"

echo "=== Seeding E2E test data ==="

# --- Topic ---
echo "Creating topic..."
curl -sf -X POST "$API/topics" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "e2e-go-basics",
    "title": "Go Fundamentals",
    "description": "A comprehensive introduction to Go programming covering syntax, data types, and control flow.",
    "difficulty": "foundational",
    "status": "published",
    "estimated_hours": 4,
    "tags": ["go", "programming", "fundamentals"]
  }'

# --- Module 1 ---
echo "Creating module 1..."
curl -sf -X POST "$API/modules" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "e2e-mod-1",
    "topic_id": "e2e-go-basics",
    "title": "Getting Started with Go",
    "description": "Set up your environment and write your first Go program.",
    "sort_order": 1,
    "estimated_minutes": 45
  }'

# --- Module 2 ---
echo "Creating module 2..."
curl -sf -X POST "$API/modules" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "e2e-mod-2",
    "topic_id": "e2e-go-basics",
    "title": "Data Types and Control Flow",
    "description": "Learn about variables, types, and flow control in Go.",
    "sort_order": 2,
    "estimated_minutes": 60
  }'

# --- Lesson 1 (Module 1): All text-based content types ---
echo "Creating lesson 1 (text, callout, table, image)..."
curl -sf -X POST "$API/lessons" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "e2e-les-1",
    "module_id": "e2e-mod-1",
    "title": "Introduction to Go",
    "sort_order": 1,
    "estimated_minutes": 15,
    "content": [
      {
        "type": "text",
        "body": "Go is a statically typed, compiled language designed at Google. It provides **excellent concurrency** support and produces fast, reliable software.\n\nGo was created by Robert Griesemer, Rob Pike, and Ken Thompson in 2009."
      },
      {
        "type": "callout",
        "variant": "info",
        "body": "Go is sometimes called Golang because of its domain name golang.org."
      },
      {
        "type": "callout",
        "variant": "tip",
        "body": "Install Go from the official website at go.dev/dl for the best experience."
      },
      {
        "type": "callout",
        "variant": "warning",
        "body": "Make sure your GOPATH is set correctly before running Go programs."
      },
      {
        "type": "callout",
        "variant": "prerequisite",
        "body": "Basic command-line familiarity is required for this lesson."
      },
      {
        "type": "table",
        "headers": ["Feature", "Go", "Python", "Java"],
        "rows": [
          ["Typing", "Static", "Dynamic", "Static"],
          ["Compilation", "Compiled", "Interpreted", "JIT Compiled"],
          ["Concurrency", "Goroutines", "Threading", "Threads"],
          ["Memory", "GC", "GC", "GC"]
        ]
      },
      {
        "type": "image",
        "url": "https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Blue.png",
        "alt": "Go programming language logo",
        "caption": "The official Go gopher logo"
      }
    ],
    "review_questions": [
      {
        "question": "Who created the Go programming language?",
        "answer": "Go was created by Robert Griesemer, Rob Pike, and Ken Thompson at Google in 2009.",
        "concepts_tested": ["go-history"]
      },
      {
        "question": "What type system does Go use?",
        "answer": "Go uses a static type system where types are checked at compile time.",
        "concepts_tested": ["type-system"]
      }
    ]
  }'

# --- Lesson 2 (Module 1): Code sections with various languages ---
echo "Creating lesson 2 (code sections, exercises)..."
curl -sf -X POST "$API/lessons" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "e2e-les-2",
    "module_id": "e2e-mod-1",
    "title": "Your First Go Program",
    "sort_order": 2,
    "estimated_minutes": 20,
    "content": [
      {
        "type": "text",
        "body": "Let us write your first Go program — the classic Hello World example."
      },
      {
        "type": "code",
        "language": "go",
        "title": "Hello World in Go",
        "code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}",
        "explanation": "Every Go program starts with a package declaration. The main package is the entry point."
      },
      {
        "type": "code",
        "language": "bash",
        "code": "go run main.go\n# Output: Hello, World!",
        "explanation": "Use go run to compile and execute your program in one step."
      },
      {
        "type": "code",
        "language": "json",
        "title": "Go module configuration",
        "code": "{\n  \"module\": \"example.com/hello\",\n  \"go\": \"1.21\"\n}"
      },
      {
        "type": "code",
        "language": "yaml",
        "title": "CI configuration",
        "code": "name: Go CI\non: [push]\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v4\n      - uses: actions/setup-go@v5\n        with:\n          go-version: \"1.21\""
      },
      {
        "type": "code",
        "language": "typescript",
        "title": "TypeScript equivalent",
        "code": "function main(): void {\n  console.log(\"Hello, World!\");\n}\n\nmain();"
      },
      {
        "type": "code",
        "language": "javascript",
        "title": "JavaScript equivalent",
        "code": "function main() {\n  console.log(\"Hello, World!\");\n}\n\nmain();"
      }
    ],
    "exercises": [
      {
        "type": "command",
        "title": "Run Hello World",
        "instructions": "Create a file called main.go with the Hello World program shown above, then run it using the go run command.",
        "environment": "Terminal with Go installed",
        "success_criteria": [
          "Program compiles without errors",
          "Output shows Hello, World!"
        ],
        "hints": [
          "Make sure you are in the directory containing main.go",
          "The command is: go run main.go",
          "If you get an error, check that Go is installed with: go version"
        ]
      },
      {
        "type": "exploration",
        "title": "Modify the Greeting",
        "instructions": "Change the Hello World program to print your name instead. Try using fmt.Printf for formatted output.",
        "success_criteria": [
          "Program prints a custom greeting",
          "Uses fmt.Printf instead of fmt.Println"
        ],
        "hints": [
          "fmt.Printf uses format verbs like %s for strings",
          "Example: fmt.Printf(\"Hello, %s!\\n\", \"Alice\")"
        ]
      }
    ]
  }'

# --- Lesson 3 (Module 2): Diagram section ---
echo "Creating lesson 3 (diagram, cross-module nav test)..."
curl -sf -X POST "$API/lessons" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "e2e-les-3",
    "module_id": "e2e-mod-2",
    "title": "Variables and Types",
    "sort_order": 1,
    "estimated_minutes": 25,
    "content": [
      {
        "type": "text",
        "body": "Go has a rich type system. Variables can be declared using `var` or the short declaration operator `:=`."
      },
      {
        "type": "diagram",
        "format": "mermaid",
        "source": "graph TD\n    A[Variable Declaration] --> B{Short or Long?}\n    B -->|Short| C[\":= operator\"]\n    B -->|Long| D[\"var keyword\"]\n    C --> E[Type Inferred]\n    D --> F[Type Explicit]\n    D --> G[Type Inferred]",
        "title": "Variable Declaration Flow"
      },
      {
        "type": "code",
        "language": "go",
        "title": "Variable declarations",
        "code": "// Long form\nvar name string = \"Alice\"\nvar age int = 30\n\n// Short form (type inferred)\ncity := \"New York\"\ncount := 42",
        "explanation": "The short declaration := can only be used inside functions."
      },
      {
        "type": "text",
        "body": "Go supports several basic types including `int`, `float64`, `string`, `bool`, and more."
      }
    ],
    "review_questions": [
      {
        "question": "What is the difference between var and := in Go?",
        "answer": "The var keyword can be used anywhere and requires explicit type or initialization. The := operator is a short declaration that can only be used inside functions and always infers the type from the value.",
        "concepts_tested": ["variables", "type-inference"]
      }
    ]
  }'

echo ""
echo "=== Seed data created successfully ==="
echo "Topic: Go Fundamentals (e2e-go-basics)"
echo "  Module 1: Getting Started with Go (2 lessons)"
echo "  Module 2: Data Types and Control Flow (1 lesson)"
echo "  Content types: text, code, callout, diagram, table, image"
echo "  Exercises: 2 (with hints)"
echo "  Review questions: 3"
