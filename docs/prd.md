# Apollo — Product Requirements Document

**AI-Powered Curriculum Builder & Learning System**

| | |
|---|---|
| **Author** | Sean |
| **Date** | February 14, 2026 |
| **Version** | 1.3 |
| **Status** | PRD |

---

## 1. Overview

Apollo is a self-hosted service that takes a topic you want to learn, deploys AI agents in a recursive research loop to build structured curricula, and serves them through an interactive course-style frontend with built-in retention tools. The system maintains a global knowledge pool — topics are never duplicated, only connected. Concepts are the atomic unit of knowledge: defined once, referenced everywhere, tested via spaced repetition.

"Give me a topic. I'll build you a course. And connect it to everything else you're learning."

This is a single-user, self-hosted application. There are no multi-tenancy, billing, or public-facing concerns.

---

## 2. Problem Statement

Sean reads a lot, bookmarks things, constantly tries new stacks and tools — but doesn't retain or connect what he learns. Bookmarks pile up unread. Tutorials are skimmed and forgotten. Adjacent knowledge domains (networking, Linux, virtualization) stay siloed in his head instead of connected.

Existing learning platforms (Coursera, Udemy, etc.) are designed for mass audiences, not personalized deep dives. They can't take "I just set up a Proxmox server, teach me everything" and recursively build out a curriculum that includes the Proxmox-specific material AND the prerequisite networking, Linux administration, and virtualization knowledge — all connected into a single knowledge graph.

There is no system that:
1. Goes deep on a topic automatically via multi-round agent research
2. Identifies prerequisites and related topics, then researches those too
3. Structures everything into learnable curricula with clear topic boundaries
4. Connects all curricula into a unified, deduplicated knowledge pool
5. Presents it as an interactive course with hands-on exercises and retention tools

---

## 3. Goals & Non-Goals

### 3.1 Goals

- **Recursive curriculum generation:** Given a topic, research it deeply and produce a structured, multi-module curriculum. Automatically identify and research essential prerequisites.
- **Global knowledge pool:** All topics, concepts, and connections live in one deduplicated pool. A concept exists once and is referenced everywhere.
- **Interactive course frontend:** Render curricula as navigable courses with typed content sections, inline examples, exercises, and review questions.
- **Retention system:** Automatically generate flashcards from concepts. SM-2 spaced repetition with a dedicated review session UI.
- **Knowledge graph:** Maintain a concept-level connection graph across all curricula. Visualize as an interactive concept map.
- **Incremental growth:** The knowledge pool grows over time. New topics connect to existing ones. Prerequisites already in the pool are linked, not re-researched.
- **Reference wiki:** All curriculum content browsable as a searchable knowledge wiki.
- **Quality control:** Validate research agent output before it enters the pool.

### 3.2 Non-Goals

- Multi-user support or collaboration
- Mobile app (web UI is sufficient)
- Real-time/streaming curriculum generation (async is fine — research takes minutes)
- Replacing external courses or certifications
- AI tutoring / conversational learning (this is structured curriculum, not a chatbot)
- Video content generation
- Integration with Mnemos (future consideration, not v1)

---

## 4. User Stories

| ID | As a... | I want to... | So that... |
|---|---|---|---|
| US-1 | Learner | Submit a topic and have a full curriculum generated | I don't have to manually find and organize learning material |
| US-2 | Learner | See prerequisite topics auto-generated alongside my main topic | I fill knowledge gaps I didn't know I had |
| US-3 | Learner | Work through lessons in order with progress tracking | I learn methodically instead of jumping around |
| US-4 | Learner | Do hands-on exercises within lessons | I practice, not just read |
| US-5 | Learner | Review flashcards generated from concepts I've studied | I actually retain what I learn via spaced repetition |
| US-6 | Learner | See how concepts connect across different curricula | I understand how topics relate and build on each other |
| US-7 | Learner | Search across my entire knowledge pool | I can find anything I've studied by keyword |
| US-8 | Learner | Expand "helpful" or "deep background" topics on demand | My knowledge pool grows based on my actual curiosity |
| US-9 | Learner | See a visual map of all my topics and their connections | I get a bird's-eye view of what I know and where the gaps are |
| US-10 | Learner | Add personal notes to any lesson | I capture my own insights alongside the generated content |
| US-11 | Learner | Refresh an outdated curriculum | My knowledge stays current when tools and technologies evolve |
| US-12 | Learner | See research progress while agents are working | I know what's happening during long research runs |

---

## 5. Core Architecture

### 5.1 The Concept Entity: The Atom of Knowledge

This is the design decision that makes everything else work. **Concepts are the atomic unit of knowledge**, not lessons or modules.

A concept is a discrete, named, self-contained piece of knowledge:

| Concept | Defined In | Referenced By |
|---------|-----------|---------------|
| VLAN tagging | Networking Fundamentals → Module 3 | Proxmox → Module 5, Docker Networking → Module 2 |
| Bridge interface | Networking Fundamentals → Module 4 | Proxmox → Module 5, KVM → Module 3 |
| Copy-on-write (CoW) | File Systems → Module 2 | ZFS → Module 1, Btrfs → Module 1, Docker → Module 4 |
| Systemd unit files | Linux Administration → Module 6 | Proxmox → Module 2, Docker → Module 1 |

**How concepts flow through the system:**

1. **Research agent** outputs lessons with inline concept markers — "this lesson teaches these concepts" and "this lesson references these concepts."
2. **Connection resolver** deduplicates. If two curricula both define "VLAN tagging," the resolver picks the more foundational one as the canonical definition and turns the other into a reference.
3. **Frontend** renders concept references as interactive links. Click "VLAN tagging" in a Proxmox lesson → jumps to the full definition in the networking curriculum.
4. **Flashcard engine** generates cards from concept definitions automatically. One concept = one flashcard (at minimum).
5. **Knowledge wiki** has a concept index. Every concept is browsable, showing its definition and every place it appears across all curricula.
6. **Concept map** visualizes concepts as nodes, with edges for "referenced by," "prerequisite of," and "related to."

Without the concept layer, you have isolated courses. With it, you have a **knowledge graph**.

### 5.2 The Research Loop

Research is **recursive and bounded**. Not one agent going deep on everything — a loop of focused agents, each responsible for one topic.

**Flow:**

1. **Initial request:** "I want to learn everything about Proxmox."
2. **Orchestrator** checks the knowledge pool. Proxmox doesn't exist. Adds it to the research queue.
3. **Research agent** goes deep on Proxmox. Produces a structured curriculum. As part of output, classifies prerequisites:
   - **Essential** (auto-recurse): Linux administration, networking fundamentals
   - **Helpful** (user opt-in): ZFS, virtualization concepts
   - **Deep background** (curious only): OS internals, kernel namespaces
4. **Connection resolver** integrates the curriculum into the pool. Deduplicates concepts, wires up cross-references.
5. **Orchestrator** checks the pool for each essential prerequisite. "Linux administration" doesn't exist → add to queue. "Networking fundamentals" doesn't exist → add to queue.
6. **Research agents** (potentially in parallel) go deep on each missing essential. They output their own curricula with their own prerequisites.
7. **Loop continues** until all essential prerequisites exist in the pool or max depth (3) is reached.
8. **Helpful and deep background** topics are stored as "available for expansion" — the user sees them in the UI and can trigger research on demand.

**Why this works:**
- **No duplication:** A topic exists once. If you later request "Docker" and it needs "networking fundamentals," that's already in the pool from the Proxmox research.
- **Clear boundaries:** Each agent is responsible for one topic. It knows where its topic ends and another begins.
- **Context management:** Each agent has a focused context window. The orchestrator maintains the global view.
- **Taste:** The essential/helpful/deep classification gives the system judgment about what to auto-expand vs. what's a rabbit hole.

### 5.3 Prerequisite Classification

The research agent classifies each prerequisite it identifies:

| Priority | Behavior | Example (for Proxmox) |
|----------|----------|----------------------|
| **Essential** | Auto-researched. You can't meaningfully learn the parent topic without this. | Linux administration, networking fundamentals |
| **Helpful** | Stored but not auto-researched. Improves understanding but not blocking. User can expand on demand. | ZFS, virtualization concepts, Ceph |
| **Deep Background** | Stored but not auto-researched. Academic/theoretical foundation. Only for the genuinely curious. | OS kernel internals, CPU virtualization extensions |

Only **essential** prerequisites auto-recurse. This prevents the loop from chasing every thread to infinity (Proxmox → networking → OSI model → physics).

### 5.4 Large Topic Splitting

If a research agent determines a topic would exceed ~8 modules, it should split it into coherent sub-topics instead of producing one massive curriculum:

**Example: "Kubernetes"**
- The research agent's first pass determines Kubernetes is too broad for a single curriculum.
- It proposes a split: "Kubernetes Core" (pods, deployments, services), "Kubernetes Networking" (CNI, ingress, service mesh), "Kubernetes Storage" (PVs, CSI, StatefulSets), "Kubernetes Operations" (monitoring, logging, upgrades).
- Each sub-topic becomes its own entry in the research queue.
- The orchestrator creates `subset_of` relationships between them and a parent "Kubernetes" topic that serves as an index.

**Implementation:** The research agent receives a `topic_size_limit` constraint (default: 8 modules). If it determines the topic exceeds this during its initial survey pass, it returns a split proposal instead of a full curriculum. The orchestrator processes the split and queues each sub-topic.

---

## 6. Research Agent: Detailed Design

### 6.1 Research Pipeline (Per Topic)

The research agent doesn't do a single search and dump results. It follows a structured multi-pass pipeline:

**Pass 1: Survey (Broad)**
- Web search for the topic: official docs, Wikipedia, introductory guides
- Identify the scope of the topic: what are the major areas/subtopics?
- Determine if splitting is needed (see 5.4)
- Output: Proposed module structure (titles + descriptions), preliminary prerequisite list

**Pass 2: Deep Dive (Per Module)**
- For each proposed module, focused research: official documentation, tutorials, community guides, Stack Overflow/forum answers, best practices
- Multiple search queries per module to cover different angles
- Identify key concepts within each module
- Output: Draft lessons with content sections, concepts, examples

**Pass 3: Exercises & Assessment**
- For each lesson, generate appropriate exercises based on topic type (see 6.2)
- Generate review questions that test understanding, not just recall
- Generate module-level assessments
- Cross-reference concepts with the existing knowledge pool
- Output: Complete curriculum matching the schema

**Pass 4: Validation (Self-Review)**
- The research agent reviews its own output against a checklist:
  - Are all learning objectives covered by lessons?
  - Does every lesson teach or reference at least one concept?
  - Are flashcard questions testing understanding, not just terminology?
  - Are exercises actionable (not vague "try this out")?
  - Do prerequisite classifications make sense?
  - Is the difficulty progression reasonable (foundational → advanced)?
- Fix issues found during self-review before returning

### 6.2 Exercise Type Spectrum

Different topics require different exercise types. The schema supports all of these, and the research agent picks the appropriate type based on the content:

| Exercise Type | Description | When to Use | Schema `type` |
|---------------|-------------|-------------|---------------|
| **Command** | "Run this command and observe the output" | CLI tools, server admin, DevOps | `command` |
| **Configuration** | "Edit this config file to achieve X" | Server setup, infrastructure, tool config | `configuration` |
| **Exploration** | "Open the Proxmox UI, navigate to X, find Y" | GUI-based tools, dashboards | `exploration` |
| **Build** | "Create a small X that does Y" | Programming, architecture, design | `build` |
| **Troubleshooting** | "Given this error/symptom, diagnose and fix" | Debugging, operations | `troubleshooting` |
| **Scenario** | "You have requirement X. Design an approach." | Architecture, system design, planning | `scenario` |
| **Thought Experiment** | "Consider what would happen if X. Why?" | Theory, fundamentals, conceptual understanding | `thought_experiment` |

Each exercise includes:
- `title`: What you're doing
- `instructions`: Step-by-step or open-ended, depending on type
- `success_criteria`: How you know you did it right
- `hints`: Progressive hints (don't give away the answer immediately)
- `type`: From the table above
- `environment`: What you need (e.g., "a Proxmox server," "a terminal," "none — just think")

### 6.3 Research Agent Tooling

Each research agent runs as a **Claude Code CLI session** (`claude -p`) invoked by the Go orchestrator via `exec.Command`. This leverages the Max plan ($200/month) — all Claude usage, web search, and tool use are included in the subscription. **Zero marginal cost per topic.** The CLI is language-agnostic — the Go orchestrator spawns a process, passes flags, and reads structured JSON from stdout.

**Key CLI flags per research session:**

```bash
claude -p "Research topic brief here..." \
  --system-prompt-file ./skills/research.md \
  --json-schema ./schemas/curriculum.json \
  --output-format json \
  --model opus \
  --allowedTools "WebSearch,WebFetch,Read,Write" \
  --permission-mode acceptEdits
```

| Flag | Purpose |
|------|---------|
| `-p` | Non-interactive mode. Runs the query, outputs result, exits. |
| `--system-prompt-file` | Load the research skill prompt (4-pass pipeline, schema spec, self-review checklist). |
| `--json-schema` | Enforce structured output matching the curriculum schema. Claude returns valid JSON in `structured_output`. |
| `--output-format json` | Full JSON response including `structured_output`, `session_id`, usage, cost. |
| `--model opus` | Use Opus for deep research reasoning. Quality is everything. |
| `--allowedTools` | Pre-approve tools. No interactive permission prompts. |
| `--mcp-config` | Optional. Add Exa or other search MCP servers if search quality needs improving. |
| `--resume SESSION_ID` | Multi-turn: continue a previous session with full context preserved. |

**Available tools in each session:**

| Tool | Purpose |
|------|---------|
| `WebSearch` | Search the web for documentation, guides, tutorials. Built-in, included in Max plan. |
| `WebFetch` | Fetch and read a specific URL's full content. Converts HTML to markdown. |
| `Read` | Read context files prepared by the orchestrator (knowledge pool summary). |
| `Write` | Write intermediate files if needed during research. |

**Why Claude Code CLI over alternatives:**

| Option | Cost | Decision |
|--------|------|---------|
| **Claude Code CLI (`claude -p`)** | $0 marginal (Max plan) | **Chosen.** Zero cost, `--json-schema` for structured output, `--resume` for multi-turn, `--mcp-config` for search upgrades. |
| **Agent SDK (`@anthropic-ai/claude-code`)** | Pay per token (~$2-5/topic) | Rejected. Same capabilities as CLI but uses API billing, not the Max plan. Revisit if CLI has quality issues. |
| **Raw Anthropic API + Tavily/Exa** | Pay per token + per search | Rejected. Most expensive, most complex, no benefit. |

**Tradeoffs accepted:**
- Claude's WebSearch is keyword-based, not semantic/neural. Sufficient for technical docs. Exa can be added as an MCP server via `--mcp-config` if needed — zero architecture change.
- WebSearch content is encrypted — Claude can reason about it but raw results aren't extractable. Doesn't matter — we only care about the synthesized curriculum output.

### 6.4 Multi-Turn Research Pipeline

Rather than cramming all 4 research passes into a single prompt, the orchestrator runs them as a multi-turn conversation using `--resume`. Each pass gets focused instructions, context accumulates across turns, and only the final pass uses `--json-schema` for structured output.

```
Go Orchestrator                       Claude Code CLI
    │                                        │
    │  1. Prepare context:                   │
    │     - Write knowledge_pool_summary.json│
    │       to working directory             │
    │                                        │
    │  Pass 1 (Survey):                      │
    │  claude -p "Survey Proxmox VE.         │
    │    Read knowledge_pool_summary.json    │
    │    for existing topics/concepts.        │
    │    Identify modules and prereqs."       │
    │    --system-prompt-file research.md     │
    │    --output-format json                 │
    │    --model opus                         │
    │    --allowedTools "WebSearch,WebFetch,  │
    │      Read"                              │
    │──────────────────────────────────────>  │
    │                                        │  WebSearch (broad survey)
    │                                        │  WebFetch (official docs)
    │                                        │  Read (knowledge pool context)
    │  <─────────────────────────────────────│
    │  Extract session_id from JSON response │
    │                                        │
    │  Pass 2 (Deep Dive):                   │
    │  claude -p "Deep dive: flesh out each  │
    │    module with full lessons, content    │
    │    sections, and concepts."             │
    │    --resume $SESSION_ID                 │
    │    --output-format json                 │
    │──────────────────────────────────────>  │
    │                                        │  WebSearch (per module)
    │                                        │  WebFetch (tutorials, guides)
    │  <─────────────────────────────────────│
    │                                        │
    │  Pass 3 (Exercises & Assessment):      │
    │  claude -p "Generate exercises, review │
    │    questions, and module assessments."  │
    │    --resume $SESSION_ID                 │
    │    --output-format json                 │
    │──────────────────────────────────────>  │
    │                                        │  Generate exercises, questions
    │  <─────────────────────────────────────│
    │                                        │
    │  Pass 4 (Self-Review + Output):        │
    │  claude -p "Self-review against the    │
    │    checklist, fix issues, output final  │
    │    curriculum."                         │
    │    --resume $SESSION_ID                 │
    │    --json-schema curriculum_schema.json │
    │    --output-format json                 │
    │──────────────────────────────────────>  │
    │                                        │  Validate, fix issues
    │                                        │  Output structured JSON
    │  <─────────────────────────────────────│
    │                                        │
    │  Parse .structured_output from JSON    │
    │  Run connection resolver               │
    │  Store in SQLite                   │
    │  Process prerequisites → add to queue  │
    │                                        │
```

**Why multi-turn over single-shot:**
- Each pass gets a focused instruction instead of a monolithic prompt
- Context accumulates naturally — pass 2 has full context of pass 1's research
- If a pass fails or produces poor results, the orchestrator can retry just that pass
- Only the final pass needs `--json-schema`, reducing structured output failures
- Easier to add/remove/reorder passes as we iterate on research quality

### 6.5 Research Session Input

The orchestrator writes a `knowledge_pool_summary.json` file to the working directory before spawning the session. The research skill prompt (loaded via `--system-prompt-file`) tells the agent to read this file during the survey pass.

**`knowledge_pool_summary.json`:**
```json
{
  "existing_topics": [
    { "id": "linux-administration", "modules": ["filesystem", "users-permissions", "systemd", "package-management", "shell-scripting", "networking-config"] },
    { "id": "networking-fundamentals", "modules": ["osi-model", "tcp-ip", "dns", "vlans-and-trunking", "bridging", "firewalls"] }
  ],
  "existing_concepts": ["linux-bridge", "vlan-tagging", "systemd-unit-files", "iptables", "subnet-mask"]
}
```

**The prompt** (passed as the `-p` argument) contains the topic brief and constraints:
```
Research topic: proxmox-ve
Brief: I just installed Proxmox on my home server. I want to learn everything about managing it.
Depth from root: 0
Topic size limit: 8 modules
Boundary guidance: Stay within Proxmox-specific knowledge. Reference existing topics
for shared concepts. Define Proxmox-specific applications as new concepts.
```

### 6.6 Research Session Output

The orchestrator reads the JSON response from stdout and extracts `structured_output` — a complete Topic object matching the curriculum schema (section 7), including all modules, lessons, concepts, exercises, review questions, and the prerequisite list with priority classification. Schema validation is enforced by `--json-schema`.

The orchestrator validates this file against the schema before ingesting it.

### 6.6 Connection Resolver

Runs after each research agent returns. Integrates new content into the existing pool.

**Step 1: Concept Deduplication**
- Compare each `concepts_taught` entry against existing concepts in the pool.
- Exact slug match → merge. Pick the more complete definition as canonical.
- Fuzzy match (e.g., agent outputs "linux-bridge-interface" but pool has "linux-bridge") → flag for review. Use an LLM call to determine if they're the same concept. If yes, merge and alias the slug.
- No match → create new concept.

**Step 2: Cross-Reference Injection**
- For each `concepts_referenced` entry, verify the referenced concept exists.
- If it exists: create the bidirectional link (lesson → concept, concept.referenced_by → lesson).
- If it doesn't exist (agent referenced a concept from a topic not yet researched): create a placeholder concept marked as `unresolved`. It gets resolved when that topic is eventually researched.

**Step 3: Prerequisite Validation**
- Lightweight LLM check: "Given topic X, are these prerequisites reasonable?"
- Flags obviously wrong prerequisites (e.g., "Proxmox requires advanced calculus").
- Does NOT deeply validate content accuracy — that's the research agent's self-review responsibility.

**Step 4: Conflict Detection**
- If two curricula define the same concept with materially different definitions, flag for user review rather than silently picking one.
- Store both definitions with a `conflict` status. Surface in the UI for the user to resolve.

---

## 7. Curriculum Schema

The structured output format that research agents produce and the frontend consumes.

### 7.1 Topic (Top-Level)

```json
{
  "id": "proxmox-ve",
  "title": "Proxmox Virtual Environment",
  "description": "Complete guide to Proxmox VE — installation, virtual machines, containers, storage, networking, clustering, and operations.",
  "difficulty": "intermediate",
  "estimated_hours": 20,
  "tags": ["infrastructure", "virtualization", "homelab", "linux"],
  "prerequisites": {
    "essential": [
      { "topic_id": "linux-administration", "reason": "Proxmox runs on Debian Linux. You need to be comfortable with the command line, package management, and system services." },
      { "topic_id": "networking-fundamentals", "reason": "VM and container networking requires understanding of bridges, VLANs, subnets, and DNS." }
    ],
    "helpful": [
      { "topic_id": "zfs", "reason": "ZFS is the recommended storage backend for Proxmox. Not required but significantly improves storage management." },
      { "topic_id": "virtualization-concepts", "reason": "Understanding hypervisors, KVM, and QEMU helps you make better VM configuration decisions." }
    ],
    "deep_background": [
      { "topic_id": "operating-system-internals", "reason": "Understanding kernel namespaces and cgroups explains how LXC containers work under the hood." }
    ]
  },
  "related_topics": ["docker", "kubernetes", "terraform", "ansible"],
  "modules": ["..."],
  "source_urls": ["https://pve.proxmox.com/wiki/", "..."],
  "generated_at": "2026-02-14T...",
  "version": 1
}
```

### 7.2 Module

```json
{
  "id": "proxmox-ve/networking",
  "title": "Networking in Proxmox",
  "description": "Configure virtual networks, bridges, VLANs, and firewall rules for VMs and containers.",
  "learning_objectives": [
    "Configure Linux bridges for VM connectivity",
    "Set up VLAN-aware networking",
    "Manage the Proxmox firewall",
    "Troubleshoot common networking issues"
  ],
  "estimated_minutes": 90,
  "order": 5,
  "lessons": ["..."],
  "assessment": {
    "questions": [
      {
        "type": "conceptual",
        "question": "Why does Proxmox use Linux bridges instead of virtual switches by default?",
        "answer": "Linux bridges are native to the kernel, require no additional software, and integrate directly with the host's network stack. They're simpler to configure and debug than OVS for most use cases.",
        "concepts_tested": ["linux-bridge", "virtual-switch"]
      },
      {
        "type": "practical",
        "question": "You have two VMs that need to communicate with each other but NOT with the external network. Describe how you'd configure this.",
        "answer": "Create an internal-only bridge (vmbr1) with no physical interface attached. Assign both VMs to this bridge. They can communicate via the bridge but have no route to the outside.",
        "concepts_tested": ["linux-bridge", "network-isolation"]
      }
    ]
  }
}
```

### 7.3 Lesson

```json
{
  "id": "proxmox-ve/networking/vlan-configuration",
  "title": "VLAN Configuration in Proxmox",
  "order": 3,
  "estimated_minutes": 15,
  "content": {
    "sections": [
      {
        "type": "text",
        "body": "Proxmox supports VLAN-aware bridges, allowing you to tag VM traffic with VLAN IDs without creating separate bridges for each VLAN..."
      },
      {
        "type": "callout",
        "variant": "prerequisite",
        "body": "This lesson assumes you understand VLAN tagging (802.1Q). If not, review Networking Fundamentals → VLANs and Trunking.",
        "concept_ref": "vlan-tagging"
      },
      {
        "type": "code",
        "language": "bash",
        "title": "Enable VLAN-aware bridge",
        "code": "auto vmbr0\niface vmbr0 inet static\n    address 10.0.0.1/24\n    bridge-ports enp0s3\n    bridge-stp off\n    bridge-fd 0\n    bridge-vlan-aware yes\n    bridge-vids 2-4094",
        "explanation": "The bridge-vlan-aware yes flag enables 802.1Q VLAN support on this bridge."
      },
      {
        "type": "diagram",
        "format": "mermaid",
        "title": "VLAN-Aware Bridge Architecture",
        "source": "graph LR\n    VM1[VM 1 - VLAN 10] --> vmbr0\n    VM2[VM 2 - VLAN 20] --> vmbr0\n    VM3[VM 3 - VLAN 10] --> vmbr0\n    vmbr0[vmbr0 - VLAN-Aware Bridge] --> NIC[Physical NIC - Trunk]\n    NIC --> Switch[Managed Switch]"
      }
    ]
  },
  "concepts_taught": [
    {
      "id": "proxmox-vlan-aware-bridge",
      "name": "VLAN-aware bridge (Proxmox)",
      "definition": "A Linux bridge with 802.1Q support enabled, allowing per-VM VLAN tagging without creating separate bridges for each VLAN.",
      "flashcard": {
        "front": "What does bridge-vlan-aware yes do in a Proxmox network config?",
        "back": "Enables 802.1Q VLAN support on the bridge, so individual VMs can be assigned VLAN tags without needing separate bridges per VLAN."
      }
    }
  ],
  "concepts_referenced": [
    { "id": "vlan-tagging", "defined_in": "networking-fundamentals/vlans-and-trunking" },
    { "id": "linux-bridge", "defined_in": "networking-fundamentals/bridging" }
  ],
  "examples": [
    {
      "title": "Assign a VM to VLAN 100",
      "description": "Set the VLAN tag on a VM's network device via CLI.",
      "code": "qm set 100 -net0 virtio,bridge=vmbr0,tag=100",
      "explanation": "The tag=100 parameter tags all traffic from this VM with VLAN ID 100."
    }
  ],
  "exercises": [
    {
      "type": "hands_on",
      "title": "Create an isolated management VLAN",
      "instructions": "Create a VLAN-aware bridge and configure two VMs: one on VLAN 10 (management) and one on VLAN 20 (application). Verify they cannot communicate with each other.",
      "environment": "A running Proxmox server with at least two VMs",
      "success_criteria": [
        "Both VMs can reach the Proxmox host",
        "VM on VLAN 10 cannot ping VM on VLAN 20",
        "bridge vlan show confirms VLAN assignments"
      ],
      "hints": [
        "You'll need a managed switch or virtual router for inter-VLAN routing",
        "Use bridge vlan show on the Proxmox host to verify VLAN assignments"
      ]
    }
  ],
  "review_questions": [
    {
      "question": "What's the advantage of a VLAN-aware bridge over creating separate bridges for each VLAN?",
      "answer": "A single VLAN-aware bridge handles all VLANs, reducing configuration complexity and resource usage. Without it, you'd need one bridge per VLAN.",
      "concepts_tested": ["proxmox-vlan-aware-bridge", "vlan-tagging"]
    }
  ]
}
```

### 7.4 Content Section Types

| Type | Purpose | Rendering | Required Fields |
|------|---------|-----------|-----------------|
| `text` | Explanatory prose | Markdown paragraphs | `body` |
| `code` | Commands, config files, code | Syntax-highlighted block | `language`, `code`, `explanation` (optional: `title`) |
| `callout` | Prerequisites, warnings, tips | Styled callout box | `variant` (prerequisite/warning/tip/info), `body` (optional: `concept_ref`) |
| `diagram` | Visual architecture/flow | Mermaid renderer or static image | `format` (mermaid/image), `source`, `title` |
| `table` | Comparison, reference data | HTML table | `headers`, `rows` |
| `image` | Screenshots, diagrams | Image with caption | `url`, `alt`, `caption` |

This covers the vast majority of technical content. The frontend knows exactly how to render each type — no ambiguous blobs.

---

## 8. System Design

### 8.1 Architecture Overview

```
                                    ┌─────────────────┐
                                    │   React Frontend │
                                    │   (Course View,  │
                                    │    Wiki, Map,    │
                                    │    Review)       │
                                    └────────┬────────┘
                                             │ HTTP/JSON
                                             │
┌────────────────────────────────────────────┴──────────────────────────────────┐
│                              Go API Server                                    │
│                                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐ │
│  │ Curriculum    │  │ Learning     │  │ Search       │  │ Research         │ │
│  │ API          │  │ Progress API │  │ API          │  │ Orchestrator     │ │
│  │              │  │              │  │              │  │                  │ │
│  │ GET topics   │  │ GET/PUT      │  │ GET search   │  │ POST /research   │ │
│  │ GET modules  │  │ progress     │  │ concepts     │  │ GET /research    │ │
│  │ GET lessons  │  │ GET/PUT      │  │ topics       │  │   /status        │ │
│  │ GET concepts │  │ retention    │  │ fulltext     │  │                  │ │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  └──────┬───────────┘ │
│         │                  │                  │                  │             │
│         └──────────────────┴──────────────────┴──────────────────┘             │
│                                    │                                          │
│                              ┌─────┴─────┐                                   │
│                              │  SQLite    │                                   │
│                              └───────────┘                                   │
│                                                                               │
│  Research Orchestrator internals:                                             │
│  ┌─────────────┐    ┌──────────────┐    ┌──────────────────┐                │
│  │ Research     │───>│ Connection   │───>│ Knowledge Pool   │                │
│  │ Queue        │    │ Resolver     │    │ (SQLite)         │                │
│  └──────┬──────┘    └──────────────┘    └──────────────────┘                │
│         │                                                                     │
│   ┌─────┴──────────────────────────────────────────────┐                    │
│   │  Claude Code Sessions (spawned per topic)          │                    │
│   │                                                     │                    │
│   │  Session A: proxmox-ve                             │                    │
│   │    Tools: WebSearch, WebFetch, Read, Write          │                    │
│   │    Input:  context files + topic brief              │                    │
│   │    Output: curriculum.json                          │                    │
│   │                                                     │                    │
│   │  Session B: linux-administration (parallel)        │                    │
│   │  Session C: networking-fundamentals (parallel)     │                    │
│   └─────────────────────────────────────────────────────┘                    │
└───────────────────────────────────────────────────────────────────────────────┘
```

### 8.2 Key Technical Decisions

**Go for the backend.** The orchestrator needs to manage concurrent research agents, process queues, and serve the frontend API — all things Go excels at. Consistent with Sean's primary backend language.

**SQLite for storage.** Single-user, self-hosted app — no concurrent write contention. SQLite handles the relational model (foreign keys), JSON columns (`json_extract()`), and full-text search (FTS5) with zero infrastructure. The database is a file in `./data/apollo.db`. Backup = copy the file. No separate container, no connection string, no pg_dump cron. Pure Go driver (`modernc.org/sqlite`) means no CGO dependency.

**Claude Code sessions for research agents.** Each topic research runs as a dedicated Claude Code session, spawned by the Go orchestrator via CLI. This leverages the Max plan ($200/month) — all token usage, web search, and tool use are included. Zero marginal cost per topic. The orchestrator communicates with sessions via the filesystem (context files in, curriculum JSON out). No Claude API client library needed in the Go code.

**No external search API.** Tavily, Exa, and other search providers were evaluated and rejected. On the Max plan, Claude Code's built-in WebSearch and WebFetch are included at no additional cost. The cost difference per topic ($0.10-0.30) is negligible, and eliminating a dependency + API key is worth more than semantic search capabilities we don't need for technical documentation research.

**SM-2 for spaced repetition.** Well-understood algorithm, simple to implement, proven effective. Same as Mnemos — if both are eventually built, the implementation patterns are directly portable.

**React + TypeScript for frontend.** Standard choice. Rich ecosystem for the specific rendering needs: syntax highlighting (Shiki), diagrams (Mermaid.js), graph visualization (D3.js), markdown rendering.

### 8.3 Cost Model

All research runs on Claude Code sessions using the Max plan ($200/month). This plan includes:
- Unlimited Claude API token usage (Opus, Sonnet, Haiku)
- Unlimited WebSearch and WebFetch calls
- No per-search or per-token billing

**Marginal cost per topic: $0.** The only cost is the flat monthly subscription, which Sean is already paying for other use (development, this ideas repo, etc.). Apollo research is an incremental use of existing capacity.

**Rate limits are a non-concern.** The Max plan has an extremely generous usage window. Research sessions are a fraction of available capacity.

**Comparison to alternatives (rejected):**
| Approach | Cost per Topic | Annual Cost (50 topics) | Dependencies |
|----------|---------------|------------------------|-------------|
| Claude Code sessions (Max plan) | $0 marginal | $0 marginal (already paying $200/mo) | None |
| Claude API + Tavily | ~$2.25-4.25 | ~$112-212 | Tavily API key |
| Claude API + Exa | ~$2.19-4.19 | ~$110-210 | Exa API key |
| Claude API + Claude WebSearch | ~$2.30-4.30 | ~$115-215 | None (but per-token + per-search billing) |

The choice is obvious. Zero marginal cost, zero extra dependencies.

### 8.4 Infrastructure

| Concern | Approach |
|---------|----------|
| Hosting | Self-hosted (single Docker container on homelab or VPS, or just a binary). |
| Deployment | Single binary or `docker run`. No multi-container orchestration needed. |
| CI/CD | GitHub Actions: lint, test, build Docker image. Manual deploy for self-hosted. |
| Monitoring | Structured logging (zerolog). Research job status tracked in DB and exposed via API. No need for Prometheus/Grafana for a single-user app. |
| Backups | Copy `apollo.db`. The knowledge pool is the valuable data — back it up. A cron running `sqlite3 apollo.db ".backup backup.db"` handles online backups safely. |

---

## 9. API Design

### 9.1 Research

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/research` | Start a new research job. Body: `{ topic: string, brief?: string }` |
| GET | `/api/research/jobs` | List all research jobs (active and completed) |
| GET | `/api/research/jobs/:id` | Get status and progress of a specific research job |
| POST | `/api/research/jobs/:id/cancel` | Cancel a running research job |
| POST | `/api/research/expand/:topicId` | Expand a "helpful" or "deep_background" prerequisite into a full curriculum |
| POST | `/api/research/refresh/:topicId` | Refresh an existing curriculum (re-research with current version as context) |

**Research job status flow:** `queued` → `researching` → `resolving` → `published` (or `failed`)

**Progress reporting:** The research job record includes a `progress` field that updates as the agent works:

```json
{
  "id": "job-123",
  "status": "researching",
  "topic": "proxmox-ve",
  "progress": {
    "current_pass": "deep_dive",
    "modules_planned": 7,
    "modules_completed": 3,
    "current_module": "Storage Management",
    "prerequisites_found": {
      "essential": ["linux-administration", "networking-fundamentals"],
      "helpful": ["zfs", "virtualization-concepts"],
      "deep_background": ["operating-system-internals"]
    },
    "concepts_identified": 42,
    "elapsed_seconds": 180
  },
  "started_at": "2026-02-14T10:30:00Z"
}
```

### 9.2 Curriculum

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/topics` | List all topics in the knowledge pool |
| GET | `/api/topics/:id` | Get a topic with its modules (no lesson content) |
| GET | `/api/topics/:id/full` | Get a topic with all modules, lessons, concepts (full tree) |
| GET | `/api/modules/:id` | Get a module with its lessons |
| GET | `/api/lessons/:id` | Get a single lesson with full content |
| GET | `/api/concepts` | List all concepts (paginated, filterable by topic) |
| GET | `/api/concepts/:id` | Get a concept with all its references |
| GET | `/api/concepts/:id/references` | Get all lessons that teach or reference this concept |

### 9.3 Learning Progress

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/progress/topics/:id` | Get learning progress for a topic (per-lesson status) |
| PUT | `/api/progress/lessons/:id` | Update lesson progress (status, notes) |
| GET | `/api/progress/summary` | Dashboard data: completion %, active topics, review stats |

### 9.4 Spaced Repetition

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/review/due` | Get concepts due for review today |
| GET | `/api/review/stats` | Review stats: due today, upcoming, total mastered |
| POST | `/api/review/:conceptId` | Submit a review rating. Body: `{ rating: "forgot" | "hard" | "good" | "easy" }` |

### 9.5 Search

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/search?q=...` | Full-text search across topics, lessons, and concepts |

### 9.6 Knowledge Graph

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/graph` | Get the full topic + concept graph for visualization (nodes and edges) |
| GET | `/api/graph/topic/:id` | Get the subgraph for a single topic (its concepts and their connections) |

---

## 10. Data Model (SQLite)

### topics

| Column | Type | Description |
|--------|------|-------------|
| id | TEXT PK | Slug, e.g. `proxmox-ve` |
| title | TEXT NOT NULL | Display title |
| description | TEXT | 1-2 paragraph overview |
| difficulty | TEXT | `foundational`, `intermediate`, `advanced` |
| estimated_hours | REAL | Total estimated learning time |
| tags | TEXT (JSON) | Taxonomy tags |
| status | TEXT NOT NULL | `researching`, `draft`, `published`, `outdated` |
| version | INTEGER DEFAULT 1 | Curriculum version |
| source_urls | TEXT (JSON) | URLs used during research |
| generated_at | TEXT (ISO 8601) | When generated |
| generated_by | TEXT | Agent/model identifier |
| parent_topic_id | TEXT FK → topics | For sub-topics created by splitting |
| created_at | TEXT (ISO 8601) DEFAULT NOW() | |
| updated_at | TEXT (ISO 8601) DEFAULT NOW() | |

### modules

| Column | Type | Description |
|--------|------|-------------|
| id | TEXT PK | `topic-slug/module-slug` |
| topic_id | TEXT FK → topics NOT NULL | |
| title | TEXT NOT NULL | |
| description | TEXT | |
| learning_objectives | TEXT (JSON) | Array of strings |
| estimated_minutes | INTEGER | |
| sort_order | INTEGER NOT NULL | Position within topic |
| assessment | TEXT (JSON) | Module assessment questions |

### lessons

| Column | Type | Description |
|--------|------|-------------|
| id | TEXT PK | `topic/module/lesson-slug` |
| module_id | TEXT FK → modules NOT NULL | |
| title | TEXT NOT NULL | |
| sort_order | INTEGER NOT NULL | Position within module |
| estimated_minutes | INTEGER | |
| content | TEXT (JSON) NOT NULL | Array of typed content sections |
| examples | TEXT (JSON) | Worked examples |
| exercises | TEXT (JSON) | Hands-on exercises |
| review_questions | TEXT (JSON) | End-of-lesson review |

### concepts

| Column | Type | Description |
|--------|------|-------------|
| id | TEXT PK | Slug, globally unique |
| name | TEXT NOT NULL | Display name |
| definition | TEXT NOT NULL | 1-3 sentence explanation |
| defined_in_lesson | TEXT FK → lessons | Canonical lesson where taught |
| defined_in_topic | TEXT FK → topics | Shortcut for "which topic owns this" |
| difficulty | TEXT | `foundational`, `intermediate`, `advanced` |
| flashcard_front | TEXT | Spaced repetition question |
| flashcard_back | TEXT | Spaced repetition answer |
| status | TEXT DEFAULT 'active' | `active`, `unresolved`, `conflict` |
| aliases | TEXT (JSON) | Alternative slugs that map to this concept |

### concept_references

| Column | Type | Description |
|--------|------|-------------|
| concept_id | TEXT FK → concepts | |
| lesson_id | TEXT FK → lessons | |
| context | TEXT | How it's used in this lesson |
| PRIMARY KEY | (concept_id, lesson_id) | |

### topic_prerequisites

| Column | Type | Description |
|--------|------|-------------|
| topic_id | TEXT FK → topics | |
| prerequisite_topic_id | TEXT FK → topics | |
| priority | TEXT NOT NULL | `essential`, `helpful`, `deep_background` |
| reason | TEXT | Why this is a prerequisite |
| PRIMARY KEY | (topic_id, prerequisite_topic_id) | |

### topic_relations

| Column | Type | Description |
|--------|------|-------------|
| topic_a | TEXT FK → topics | |
| topic_b | TEXT FK → topics | |
| relation_type | TEXT NOT NULL | `related`, `builds_on`, `contrasts_with`, `subset_of` |
| description | TEXT | |
| PRIMARY KEY | (topic_a, topic_b) | |

### expansion_queue

| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER PRIMARY KEY AUTOINCREMENT | |
| topic_id | TEXT NOT NULL | Topic slug to research |
| requested_by_topic | TEXT FK → topics | Which topic identified this as a prerequisite |
| priority | TEXT NOT NULL | `essential`, `helpful`, `deep_background` |
| reason | TEXT | |
| status | TEXT DEFAULT 'available' | `available`, `queued`, `researching`, `completed`, `skipped` |
| depth_from_root | INTEGER | How many levels from the original request |

### research_jobs

| Column | Type | Description |
|--------|------|-------------|
| id | TEXT PK | Job identifier |
| root_topic | TEXT | The original user-requested topic |
| current_topic | TEXT | The topic currently being researched |
| status | TEXT NOT NULL | `queued`, `researching`, `resolving`, `published`, `failed`, `cancelled` |
| progress | TEXT (JSON) | Structured progress data |
| error | TEXT | Error message if failed |
| started_at | TEXT (ISO 8601) | |
| completed_at | TEXT (ISO 8601) | |

### learning_progress

| Column | Type | Description |
|--------|------|-------------|
| lesson_id | TEXT PK FK → lessons | |
| status | TEXT DEFAULT 'not_started' | `not_started`, `in_progress`, `completed` |
| started_at | TEXT (ISO 8601) | |
| completed_at | TEXT (ISO 8601) | |
| notes | TEXT | User's personal notes |

### concept_retention

| Column | Type | Description |
|--------|------|-------------|
| concept_id | TEXT PK FK → concepts | |
| status | TEXT DEFAULT 'new' | `new`, `learning`, `reviewing`, `mastered` |
| next_review | TEXT (ISO 8601) | |
| review_count | INTEGER DEFAULT 0 | |
| ease_factor | REAL DEFAULT 2.5 | SM-2 ease factor |
| interval_days | INTEGER DEFAULT 0 | Current review interval |
| last_reviewed | TEXT (ISO 8601) | |
| last_rating | TEXT | `forgot`, `hard`, `good`, `easy` |

---

## 11. Spaced Repetition (SM-2)

### 11.1 Algorithm

When a concept enters `learning` status (user completes the lesson that teaches it), `next_review` is set to tomorrow. After each review:

```
Input: rating (forgot=0, hard=1, good=2, easy=3)

if rating == forgot:
    interval = 1 day
    ease_factor = max(1.3, ease_factor - 0.2)
else if rating == hard:
    interval = max(1, interval * 1.2)
    ease_factor = max(1.3, ease_factor - 0.15)
else if rating == good:
    if review_count == 0: interval = 1
    else if review_count == 1: interval = 3
    else: interval = round(interval * ease_factor)
    ease_factor = ease_factor + 0.0  (no change)
else if rating == easy:
    if review_count == 0: interval = 4
    else: interval = round(interval * ease_factor * 1.3)
    ease_factor = ease_factor + 0.15

next_review = now + interval days
review_count += 1
```

### 11.2 Concept Lifecycle

```
new → learning → reviewing → mastered
       ↑              │
       └──────────────┘ (forgot resets to learning)
```

- **New:** Concept exists but user hasn't studied the lesson yet.
- **Learning:** User completed the lesson. First review due tomorrow.
- **Reviewing:** Active in the review queue. Interval grows with successful recalls.
- **Mastered:** Interval exceeds 90 days. Exits the active review queue (still reviewable on demand).

### 11.3 Review Session UX

1. User opens Review Session. Sees count of concepts due today.
2. Flashcard appears: front only (question).
3. User thinks about the answer (no time limit).
4. User clicks "Show Answer." Back is revealed.
5. User self-rates: Forgot / Hard / Good / Easy.
6. Next flashcard appears. After the last one: summary (reviewed X, next due dates).
7. If the user rated "Forgot," a link to the relevant lesson is shown for re-study.

---

## 12. Curriculum Updates

### 12.1 When to Update

Curricula can become outdated when:
- A new version of the software is released (Proxmox 8 → 9)
- Best practices change
- New features are added to the topic
- The user explicitly requests a refresh

### 12.2 Update Flow

1. User clicks "Refresh" on a topic, or the system flags a topic as potentially outdated (based on age — configurable threshold, default 6 months).
2. The orchestrator triggers a refresh research job.
3. The research agent receives the **existing curriculum as context** alongside its research brief. It's told: "Here is the current curriculum for Proxmox (version 1). Research what has changed. Produce an updated version."
4. The agent outputs a new version of the curriculum.
5. The orchestrator diffs old vs. new:
   - New lessons are added.
   - Removed lessons are archived (not deleted — learning progress preserved).
   - Modified lessons are updated. If a concept's definition changed, the concept is flagged for re-review in spaced repetition.
6. The topic's `version` field increments. A changelog is stored.

### 12.3 Impact on Retention

When a concept's definition changes during a curriculum update:
- The concept's `status` is set to `updated`.
- Its spaced repetition state is reset to `learning` (next review tomorrow).
- The flashcard is regenerated from the new definition.
- The user sees a notification: "3 concepts in Proxmox were updated. They've been added back to your review queue."

---

## 13. Frontend Architecture

### 13.1 Views

**Dashboard (Home)**
- Topic cards: each topic in the pool with title, difficulty badge, completion %, module count
- "Currently Studying" section: topics with `in_progress` lessons
- Review Queue widget: "X concepts due for review today" with a button to start a session
- "Available for Expansion" section: helpful/deep prerequisites not yet researched
- Research Jobs widget: status of any active research jobs
- Concept Map thumbnail: click to expand

**Course View (Primary Learning)**
- Left sidebar: collapsible module list with completion indicators per module
- Main area: lesson content rendered from typed sections
- Concept chips: key concepts for this lesson shown as clickable chips, linking to their canonical definition
- Prerequisite callouts: when a lesson references a concept from another topic, rendered as a styled callout with a direct link
- Progress bar per module
- "Mark Complete" button per lesson
- Review questions: collapsible section at lesson end
- Personal notes: inline text area per lesson

**Knowledge Wiki**
- Topic index: browse all topics
- Concept index: alphabetical list of all concepts across all curricula
- Concept detail page: definition, flashcard, every lesson that teaches or references it, related concepts
- Full-text search across everything
- Breadcrumb navigation: Topic → Module → Lesson

**Concept Map**
- Force-directed graph (D3.js)
- Topics as large nodes, concepts as small nodes
- Edge types: prerequisite (directional arrow), reference (dotted line), related (solid line)
- Color-coded by: topic membership, difficulty, or study progress (user toggle)
- Click a node to navigate to that topic/concept
- Zoom controls: zoom into a single topic to see its internal concept graph
- Filter: show/hide concept nodes, show only topic-level graph

**Review Session**
- Flashcard UI: card with front text → reveal button → back text
- Rating buttons: Forgot / Hard / Good / Easy
- Progress indicator: "3 of 12"
- After completion: summary stats (reviewed, next review dates)
- "Re-study" link on each card → jumps to the lesson where the concept is taught

### 13.2 Real-Time Research Progress

When a research job is running, the dashboard shows a live progress view:

- Current topic being researched
- Current pass (survey / deep dive / exercises / validation)
- Modules planned and completed
- Prerequisites discovered so far
- Concepts identified so far
- Elapsed time

Implementation: polling the `/api/research/jobs/:id` endpoint every 5 seconds. No WebSocket needed — research takes minutes, not milliseconds.

### 13.3 Tech Stack (Frontend)

| Layer | Choice | Rationale |
|-------|--------|-----------|
| Framework | React + TypeScript | Standard, strong ecosystem. SPA is fine for single-user. |
| Routing | React Router v7 | Client-side routing. No SSR needed for single-user self-hosted. |
| Styling | Tailwind CSS | Rapid styling, utility-first. |
| Data Fetching | TanStack Query | Caching, background refetching, mutation handling. |
| State | Zustand (minimal) | Only for UI state (sidebar open/closed, active filters). Server state via TanStack Query. |
| Code Highlighting | Shiki | High-quality syntax highlighting, many language grammars. |
| Diagrams | Mermaid.js | Render mermaid diagrams from curriculum content. |
| Graph Visualization | D3.js (force-directed) | Concept map. Flexible, performant, well-documented. |
| Markdown | react-markdown + rehype plugins | Render `text` sections. |

---

## 14. Tech Stack Summary

| Layer | Choice | Rationale |
|-------|--------|-----------|
| Research Agents | Claude Code sessions (Opus) | Spawned per topic by the orchestrator. WebSearch + WebFetch built in. Zero marginal cost on Max plan. |
| Orchestrator / API Server | Go | Manages research queue, spawns Claude Code sessions, serves REST API. Strong concurrency. |
| Database | SQLite (via modernc.org/sqlite) | Single-file DB. FKs, JSON functions, FTS5. Zero infrastructure. |
| Frontend | React + TypeScript + Tailwind | Interactive course UI, wiki, concept map, review sessions. |
| Code Highlighting | Shiki | Syntax highlighting in lessons. |
| Diagrams | Mermaid.js | Architecture and flow diagrams in lessons. |
| Graph Viz | D3.js | Concept map force-directed graph. |
| Spaced Repetition | Custom SM-2 (Go backend) | Simple, proven algorithm. |
| Containerization | Docker Compose | Single `docker compose up` deployment. |

---

## 15. Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| `DATABASE_PATH` | `./data/apollo.db` | SQLite database file path |
| `SERVER_PORT` | `8080` | API server port |
| `CLAUDE_CODE_PATH` | `claude` | Path to the Claude Code CLI binary |
| `MAX_RESEARCH_DEPTH` | `3` | Maximum prerequisite recursion depth |
| `MAX_PARALLEL_AGENTS` | `3` | Maximum concurrent Claude Code research sessions |
| `TOPIC_SIZE_LIMIT` | `8` | Max modules per topic before splitting |
| `AUTO_EXPAND_PRIORITY` | `essential` | Which priority levels auto-recurse |
| `CURRICULUM_STALE_DAYS` | `180` | Days before a curriculum is flagged as potentially outdated |
| `MASTERY_THRESHOLD_DAYS` | `90` | Review interval at which a concept is marked "mastered" |
| `RESEARCH_WORK_DIR` | `./data/research` | Temporary directory for research session context/output files |

---

## 16. Risks & Mitigations

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Research agent produces inaccurate content | User learns wrong information | Medium | Self-review pass (Pass 4). User feedback mechanism (flag lesson as inaccurate). Curriculum is regenerable — if quality is poor, re-research. |
| Research agent can't produce valid structured output | Curriculum fails to parse, nothing renders | Low-Medium | Schema validation on ingest. Retry with explicit error feedback. The research skill prompt includes the full schema spec and examples. |
| WebSearch returns low-quality results for niche topics | Thin or shallow curricula | Medium | Agent uses multiple search strategies and follows up with WebFetch on specific URLs. Can search for official docs directly. If search quality becomes a bottleneck, Exa can be added as an MCP server — the research session gains a `exa_search` tool alongside WebSearch with zero architecture changes. |
| Knowledge pool grows large, concept deduplication becomes noisy | False merges or missed duplicates | Low (at personal scale) | Fuzzy matching with LLM confirmation. Conflict detection surfaces issues to the user. |
| Curriculum schema is too rigid for some topic types | Some topics don't fit the module/lesson/concept structure | Low | Schema is already flexible (TEXT (JSON) content sections, multiple exercise types). Can extend section types without breaking existing data. |
| Spaced repetition fatigue | Review queue grows overwhelming as pool expands | Medium | Mastery threshold removes concepts from active review. Daily review cap (configurable). Focus reviews on topics the user is actively studying. |
| Topic splitting produces awkward boundaries | Sub-topics feel arbitrary or have too much overlap | Medium | Research agent proposes splits for orchestrator review. User can override. Connection resolver handles cross-references between sub-topics. |
| Claude Code CLI interface changes | Orchestrator's session spawning breaks | Low | Pin Claude Code version. The CLI's `--print --output-format json` interface is stable. Alternatively, the Claude Code SDK (TypeScript) could be used if tighter integration is needed. |

---

## 17. Milestones

| Phase | Scope | Key Deliverable | Estimated Effort |
|-------|-------|-----------------|------------------|
| M1: Schema & Storage | SQLite schema + migrations, Go API server scaffold, CRUD endpoints for all entities, curriculum schema validation | API serves and validates curriculum data | 1 week |
| M2: Research Agent (Single Topic) | Research skill prompt, Claude Code session spawning, 4-pass research pipeline, structured JSON output, single topic end-to-end | Can generate a complete, validated curriculum for one topic | 1-2 weeks |
| M3: Orchestrator & Loop | Research queue, prerequisite extraction, depth control, parallel agent execution, knowledge pool checks, expansion queue | Full recursive research loop works end-to-end | 1 week |
| M4: Connection Resolver | Concept deduplication (exact + fuzzy), cross-reference injection, prerequisite validation, conflict detection | Concepts properly linked across curricula | 1 week |
| M5: Course Frontend | Course view — module sidebar, lesson rendering for all section types, syntax highlighting, mermaid diagrams, progress tracking, personal notes | Can learn from a generated curriculum in the browser | 1-2 weeks |
| M6: Retention System | SM-2 implementation, concept lifecycle, flashcard UI, review session, daily stats | Concepts retained via spaced repetition | 1 week |
| M7: Knowledge Wiki & Concept Map | Wiki view, concept index, concept detail pages, full-text search, D3.js concept map with filtering | Reference and exploration layer complete | 1-2 weeks |
| M8: Dashboard & Polish | Dashboard with topic cards, research progress, review queue widget, expansion triggers, curriculum refresh, Docker Compose deployment | Production-ready personal tool | 1 week |

**Critical path:** M1 → M2 → M3 → M4. If the schema, research agent, orchestrator loop, and connection resolver work, the product works. Everything after is rendering.

**Total estimated effort:** 8-11 weeks.

---

## 18. End-to-End Example: "Proxmox"

A complete walkthrough of what happens when the user submits "I want to learn everything about Proxmox."

### Step 1: User Submits Request

```
POST /api/research
{ "topic": "proxmox", "brief": "I just installed Proxmox on my home server. Teach me everything." }
```

Research job created. Status: `queued`.

### Step 2: Orchestrator Starts

Orchestrator normalizes the topic slug to `proxmox-ve`. Checks the knowledge pool — doesn't exist. Adds to research queue at depth 0.

### Step 3: Claude Code Session — Pass 1 (Survey)

Orchestrator spawns a Claude Code session with the research skill prompt and topic brief. Session uses WebSearch: "Proxmox VE documentation," "Proxmox tutorial beginner to advanced," "Proxmox administration guide."

Uses WebFetch to read official docs, community wiki, several tutorial series. Determines the topic fits within ~7 modules:
1. Introduction & Installation
2. Virtual Machines (KVM)
3. Containers (LXC)
4. Storage Management
5. Networking
6. Clustering & High Availability
7. Backup, Restore & Operations

No split needed (under the 8-module limit).

Identifies prerequisites:
- Essential: `linux-administration`, `networking-fundamentals`
- Helpful: `zfs`, `virtualization-concepts`
- Deep background: `operating-system-internals`

### Step 4: Claude Code Session — Pass 2 (Deep Dive)

For each of the 7 modules, the session does focused research. For "Networking":
- WebSearch: "Proxmox networking configuration," "Proxmox VLAN setup," "Proxmox bridge configuration," "Proxmox firewall guide"
- WebFetch on official networking docs, community guides, Stack Overflow threads
- Identifies concepts: `proxmox-vlan-aware-bridge`, `proxmox-sdn`, `proxmox-firewall`
- Notes existing concepts to reference: `vlan-tagging`, `linux-bridge`, `subnet-mask`

Produces 3-4 lessons per module with full content sections.

### Step 5: Claude Code Session — Pass 3 (Exercises & Assessment)

For the networking module, generates:
- Hands-on exercise: "Create an isolated management VLAN" (type: `hands_on`)
- Troubleshooting exercise: "VM can't reach the internet — diagnose" (type: `troubleshooting`)
- Exploration exercise: "Navigate the SDN panel and create a simple zone" (type: `exploration`)
- Review questions per lesson
- Module assessment with conceptual and practical questions

### Step 6: Claude Code Session — Pass 4 (Self-Review)

Session reviews its output:
- All 7 modules have learning objectives covered by lessons? Yes.
- Every lesson teaches or references at least one concept? Yes.
- Flashcard questions test understanding, not just terminology? Fixes 2 that were too surface-level.
- Exercises are actionable with clear success criteria? Yes.
- Prerequisites make sense? Yes.

Writes the validated `curriculum.json` to the output directory. Session exits.

### Step 7: Connection Resolver

- `proxmox-vlan-aware-bridge` is new — create it.
- Agent referenced `vlan-tagging` — pool doesn't have it yet (networking-fundamentals hasn't been researched). Create as `unresolved`. It will be resolved when networking-fundamentals is researched.
- Agent referenced `linux-bridge` — also `unresolved` for now.
- No conflicts detected.
- Cross-references created for all concept mentions.

Curriculum stored. Topic status: `published`.

### Step 8: Orchestrator Processes Prerequisites

Essential prerequisites:
- `linux-administration` — not in pool → add to queue at depth 1
- `networking-fundamentals` — not in pool → add to queue at depth 1

Helpful/deep — stored in `expansion_queue` as `available`.

### Step 9: Claude Code Sessions (Parallel) — Linux Admin & Networking

Orchestrator spawns two Claude Code sessions concurrently:
- Session A researches `linux-administration`. Produces 6 modules. Identifies essential prereqs: none (it's foundational). Identifies concepts: `file-permissions`, `systemd-unit-files`, `apt-package-manager`, etc.
- Session B researches `networking-fundamentals`. Produces 6 modules. Identifies essential prereqs: none. Identifies concepts: `vlan-tagging`, `linux-bridge`, `tcp-three-way-handshake`, `subnet-mask`, etc.

Both sessions use WebSearch + WebFetch independently. Each writes its own `curriculum.json`. Zero marginal cost — all included in the Max plan.

### Step 10: Connection Resolver (Second Pass)

- `vlan-tagging` is now defined by networking-fundamentals. The `unresolved` reference from the Proxmox curriculum is resolved → cross-reference created.
- `linux-bridge` now defined → resolved.
- `systemd-unit-files` defined by linux-admin. Proxmox Module 7 (operations) referenced systemd → cross-reference created.
- All concepts deduplicated. No conflicts.

### Step 11: Queue Empty — Research Complete

All essential prerequisites exist in the pool. No more items in the queue. The research job status is set to `published`.

The user's knowledge pool now contains:
- **Proxmox VE** (7 modules, ~42 concepts)
- **Linux Administration** (6 modules, ~35 concepts)
- **Networking Fundamentals** (6 modules, ~30 concepts)

Plus 3 topics available for expansion: ZFS, Virtualization Concepts, OS Internals.

### Step 12: User Opens the Course

The frontend shows the Proxmox course. Module sidebar on the left. The user starts with Module 1, Lesson 1. As they complete lessons, concepts enter the spaced repetition queue. When they reach the networking module, prerequisite callouts link to the networking fundamentals curriculum.

The concept map shows 3 topic nodes (Proxmox, Linux Admin, Networking) connected by prerequisite edges, with ~107 concept nodes floating between them, cross-referenced where they overlap.

---

## 19. Relationship to Other Projects

- **Separate from Mnemos.** Own service, own repo. Different purpose: Mnemos captures and classifies content you've already found (X bookmarks). Apollo generates structured knowledge from scratch. Future integration: Mnemos bookmarks tagged with a topic could be surfaced as supplementary reading in Apollo curricula.
- **Separate from Thoth.** Thoth is the ideas repo. Apollo is a product. Thoth's skill-based orchestration pattern (structured prompts guiding agents through a lifecycle) is a design inspiration.
- **Potential portfolio piece.** Demonstrates: Go backend, React frontend, AI agent orchestration, recursive system design, knowledge graph, spaced repetition. Strong candidate for Portfolio v2 featured projects.

---

## 20. Future Considerations (Not v1)

- **Mnemos integration:** Pull bookmarks into curricula as supplementary reading.
- **Interactive sandboxes:** Embedded terminal or VM for hands-on exercises (e.g., a Proxmox sandbox in the browser).
- **AI tutoring:** Chat with an AI about a specific lesson or concept, with curriculum context.
- **Community curricula:** Share and import curricula from other Apollo users.
- **Mobile app:** Review session on mobile for flashcard reviews on the go.
- **Browser extension:** "Learn about this" on any webpage → triggers research.
- **Voice input:** "Hey Apollo, I want to learn about Kubernetes networking" → queues research.

---

## 21. References & Inspiration

- **SM-2 Algorithm:** [SuperMemo's original paper](https://www.supermemo.com/en/archives1990-2015/english/ol/sm2)
- **Anki:** The gold standard for spaced repetition. Apollo's review session is inspired by Anki's card review flow, but with auto-generated cards.
- **Zettelkasten method:** The concept graph is inspired by Zettelkasten's principle of atomic, interconnected notes.
- **Khan Academy:** Course structure inspiration — modules with lessons, progress tracking, exercises inline.
- **Obsidian:** Knowledge wiki + graph view inspiration. Bidirectional links between concepts.
