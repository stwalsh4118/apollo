# Mermaid.js Diagram Rendering Guide

> **Date:** 2026-02-18
> **Package:** `mermaid` v11.x (installed: ^11.12.3)
> **Docs:** <https://mermaid.js.org/>
> **Repository:** <https://github.com/mermaid-js/mermaid>
> **API Reference:** <https://mermaid.js.org/config/usage.html>
> **Theme Reference:** <https://mermaid.js.org/config/theming.html>

---

## 1. Installation

Already installed in the project:

```bash
pnpm add mermaid
```

Version in `web/package.json`:

| Package | Version |
|---------|---------|
| `mermaid` | ^11.12.3 |

Mermaid has no peer dependencies. It is a self-contained library that bundles
its own rendering engine (based on D3.js and dagre internally).

---

## 2. Core API Overview

Mermaid exposes four key functions:

| Function | Purpose |
|---|---|
| `mermaid.initialize(config)` | Set global configuration (theme, security, etc.) |
| `mermaid.render(id, definition)` | Programmatically render a diagram to an SVG string |
| `mermaid.run(options)` | Auto-discover and render DOM elements with Mermaid markup |
| `mermaid.parse(definition)` | Validate Mermaid syntax without rendering |

For React integration, **`mermaid.render()`** is the recommended approach because
it returns SVG markup as a string, giving full control over DOM insertion and
avoiding conflicts with React's virtual DOM.

---

## 3. Initialisation — `mermaid.initialize()`

Call `initialize()` once before any render calls. In a React app this should
happen at the module level or inside a one-time `useEffect`.

```ts
import mermaid from 'mermaid'

mermaid.initialize({
  startOnLoad: false,   // CRITICAL for React — prevents auto-rendering on DOMContentLoaded
  theme: 'dark',        // Built-in dark theme
  securityLevel: 'loose',
  fontFamily: 'inherit',
})
```

### Key `initialize()` Options

| Option | Type | Default | Notes |
|---|---|---|---|
| `startOnLoad` | `boolean` | `true` | **Must be `false`** in React apps. Prevents Mermaid from auto-scanning the DOM. |
| `theme` | `string` | `'default'` | One of: `default`, `dark`, `forest`, `neutral`, `base` |
| `securityLevel` | `string` | `'strict'` | Use `'loose'` if diagrams need click handlers. `'strict'` is safer. |
| `fontFamily` | `string` | `'"trebuchet ms"...'` | Set to `'inherit'` to use the app's font stack. |
| `themeVariables` | `object` | `{}` | Fine-tune colors when using the `base` theme. |
| `flowchart` | `object` | `{}` | Flowchart-specific config (e.g., `{ useMaxWidth: true, htmlLabels: true }`). |
| `logLevel` | `number` | `5` | `1` = debug, `5` = fatal. Lower values produce more console output. |

---

## 4. Programmatic Rendering — `mermaid.render()`

This is the primary API for React integration.

```ts
const { svg, bindFunctions } = await mermaid.render(id, definition)
```

### Parameters

| Parameter | Type | Description |
|---|---|---|
| `id` | `string` | A **unique** identifier used as the SVG element's `id` attribute. Must be unique across the entire page. |
| `definition` | `string` | The Mermaid diagram definition (e.g., `'graph TD; A-->B'`). |

### Return Value

| Property | Type | Description |
|---|---|---|
| `svg` | `string` | The rendered SVG markup as a string. |
| `bindFunctions` | `((element: Element) => void) \| undefined` | Optional callback to attach interactive event listeners (e.g., click handlers) to the rendered SVG. |

### Usage

```ts
import mermaid from 'mermaid'

// Unique ID per diagram instance
const id = `mermaid-diagram-${crypto.randomUUID()}`
const definition = `graph TD
  A[Start] --> B{Decision}
  B -->|Yes| C[Do Something]
  B -->|No| D[Do Nothing]`

const { svg, bindFunctions } = await mermaid.render(id, definition)

// Insert into DOM
container.innerHTML = svg
if (bindFunctions) {
  bindFunctions(container)
}
```

### Unique ID Generation

Each call to `mermaid.render()` **must** receive a unique `id`. If two diagrams
share the same `id`, their SVG internals (clip paths, gradient defs, marker IDs)
will collide and one or both diagrams will render incorrectly.

Strategies for generating unique IDs in React:

| Strategy | Example | Notes |
|---|---|---|
| `React.useId()` | `const id = useId()` | Built-in to React 18+. Generates stable IDs across server/client. Contains colons, so prefix with `mermaid-`. |
| `crypto.randomUUID()` | `crypto.randomUUID()` | Globally unique. Fine for client-only rendering. |
| Counter | `useRef(counter++)` | Simple; works if only one component tree. |

**Recommended approach** for this project: Use `useId()` from React since the
project runs React 19, ensuring stable and unique IDs.

---

## 5. Auto-Discovery Rendering — `mermaid.run()`

`mermaid.run()` is the alternative API that scans DOM elements and renders them
in place. It is **less suitable for React** because it mutates the DOM directly,
conflicting with React's reconciliation.

```ts
// Render all elements matching a CSS selector
await mermaid.run({ querySelector: '.mermaid' })

// Render specific DOM nodes
await mermaid.run({ nodes: [document.getElementById('my-diagram')] })

// Suppress errors (don't throw on invalid syntax)
await mermaid.run({ suppressErrors: true })
```

**When to use `mermaid.run()`**: Only if you need Mermaid's built-in DOM
mutation behaviour (e.g., in a non-React context or a static HTML page).

**For React: Always prefer `mermaid.render()`.**

---

## 6. Syntax Validation — `mermaid.parse()`

Validates Mermaid syntax without producing SVG output. Useful for pre-flight
checks before rendering.

```ts
try {
  const result = await mermaid.parse(definition)
  // result: { diagramType: string } on success
  console.log('Valid diagram type:', result.diagramType)
} catch (error) {
  // error contains details about what's wrong with the syntax
  console.error('Invalid Mermaid syntax:', error.message)
}
```

With `suppressErrors`:

```ts
const result = await mermaid.parse(definition, { suppressErrors: true })
if (result === false) {
  // Invalid syntax — handle gracefully
}
```

---

## 7. Error Handling

### Types of Errors

| Error Source | When It Occurs | How to Catch |
|---|---|---|
| Invalid syntax | Malformed diagram definition | `try/catch` around `mermaid.render()` or `mermaid.parse()` |
| Unknown diagram type | Unrecognised diagram prefix | `UnknownDiagramError` thrown by Mermaid's `detectType()` |
| Render failure | Valid syntax but rendering fails (rare) | `try/catch` around `mermaid.render()` |

### Error Handling Pattern

```ts
import mermaid from 'mermaid'

async function renderDiagram(
  id: string,
  definition: string,
): Promise<{ svg: string } | { error: string }> {
  try {
    const { svg } = await mermaid.render(id, definition)
    return { svg }
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Unknown Mermaid error'
    return { error: message }
  }
}
```

### Global Error Handler (Optional)

Mermaid supports a global `parseError` callback. This is useful for logging but
not recommended as the primary error handling mechanism in React — prefer
`try/catch` in the component.

```ts
mermaid.parseError = (err: string, hash: unknown) => {
  console.error('[Mermaid Parse Error]', err)
}
```

---

## 8. React Integration — Recommended Pattern

### `useMermaid` Hook

This hook handles initialisation, rendering, unique IDs, and error fallback:

```tsx
import { useEffect, useId, useRef, useState } from 'react'
import mermaid from 'mermaid'

// Initialise once at module level
let initialised = false
function ensureInitialised(): void {
  if (!initialised) {
    mermaid.initialize({
      startOnLoad: false,
      theme: 'dark',
      securityLevel: 'strict',
      fontFamily: 'inherit',
    })
    initialised = true
  }
}

interface UseMermaidResult {
  /** Ref to attach to the container element */
  containerRef: React.RefObject<HTMLDivElement | null>
  /** Error message if rendering failed; null on success */
  error: string | null
  /** Whether the diagram is currently rendering */
  loading: boolean
}

export function useMermaid(definition: string): UseMermaidResult {
  const containerRef = useRef<HTMLDivElement | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const reactId = useId()

  // Sanitise the React-generated ID (contains colons) for use as an SVG id
  const diagramId = `mermaid-${reactId.replace(/:/g, '-')}`

  useEffect(() => {
    let cancelled = false

    async function render() {
      ensureInitialised()
      setLoading(true)
      setError(null)

      try {
        const { svg, bindFunctions } = await mermaid.render(
          diagramId,
          definition,
        )

        if (cancelled || !containerRef.current) return

        containerRef.current.innerHTML = svg
        if (bindFunctions) {
          bindFunctions(containerRef.current)
        }
        setError(null)
      } catch (err) {
        if (cancelled) return
        const message =
          err instanceof Error ? err.message : 'Failed to render diagram'
        setError(message)

        // Clear any partial render
        if (containerRef.current) {
          containerRef.current.innerHTML = ''
        }
      } finally {
        if (!cancelled) {
          setLoading(false)
        }
      }
    }

    render()

    return () => {
      cancelled = true
    }
  }, [definition, diagramId])

  return { containerRef, error, loading }
}
```

### `MermaidDiagram` Component

A component that uses the hook:

```tsx
import { useMermaid } from '../hooks/useMermaid'

interface MermaidDiagramProps {
  /** The Mermaid diagram definition string */
  definition: string
  /** Optional title displayed above the diagram */
  title?: string
}

export function MermaidDiagram({ definition, title }: MermaidDiagramProps) {
  const { containerRef, error, loading } = useMermaid(definition)

  return (
    <figure>
      {title && <figcaption>{title}</figcaption>}

      {/* Loading state */}
      {loading && !error && (
        <div aria-busy="true">Rendering diagram...</div>
      )}

      {/* SVG render target — Mermaid injects SVG here */}
      <div ref={containerRef} style={{ display: error ? 'none' : 'block' }} />

      {/* Error fallback — show raw source in a pre block */}
      {error && (
        <div role="alert">
          <p>Diagram render error: {error}</p>
          <pre><code>{definition}</code></pre>
        </div>
      )}
    </figure>
  )
}
```

### Key Design Decisions

| Decision | Rationale |
|---|---|
| `startOnLoad: false` | Prevents Mermaid from scanning the DOM on page load. React controls when rendering happens. |
| `useId()` for unique IDs | React 18+ hook that produces stable IDs. Avoids collisions when multiple diagrams exist on the same page. |
| `useEffect` (not `useLayoutEffect`) | `mermaid.render()` is async and does not need to block paint. Using `useEffect` avoids unnecessary render delays. |
| `containerRef.innerHTML = svg` | Direct DOM mutation inside a ref-managed container. React does not track children inside refs, so this is safe. |
| Cancellation flag | Prevents state updates on unmounted components when the render promise resolves after navigation. |
| Module-level initialisation guard | Ensures `mermaid.initialize()` is called exactly once regardless of how many diagram components mount. |

---

## 9. Dark Theme Configuration

### Built-in `dark` Theme

The simplest option. Provides a dark background with light text/lines:

```ts
mermaid.initialize({
  startOnLoad: false,
  theme: 'dark',
})
```

The built-in `dark` theme uses:
- Dark grey background for nodes
- Light text
- Muted edge/line colours

### Custom Dark Theme (Matching `github-dark`)

For closer alignment with the `github-dark` code block theme used by Shiki in
this project, use the `base` theme with custom `themeVariables`:

```ts
mermaid.initialize({
  startOnLoad: false,
  theme: 'base',
  themeVariables: {
    darkMode: true,
    background: '#0d1117',          // github-dark page background
    primaryColor: '#1f2937',        // node fill — dark grey
    primaryTextColor: '#e6edf3',    // node text — light grey
    primaryBorderColor: '#30363d',  // node border
    lineColor: '#8b949e',           // edge/arrow colour
    secondaryColor: '#161b22',      // alt node fill
    tertiaryColor: '#21262d',       // subgraph/cluster fill
    textColor: '#e6edf3',           // general text
    mainBkg: '#1f2937',             // flowchart node background
    nodeBorder: '#30363d',          // flowchart node border
    clusterBkg: '#161b22',          // subgraph background
    clusterBorder: '#30363d',       // subgraph border
    titleColor: '#e6edf3',          // diagram title colour
    edgeLabelBackground: '#0d1117', // label background on edges
    noteTextColor: '#e6edf3',
    noteBkgColor: '#161b22',
    noteBorderColor: '#30363d',
  },
})
```

### Available Built-in Themes

| Theme | Description |
|---|---|
| `default` | Standard light theme |
| `dark` | Dark background with light elements |
| `forest` | Green-tinted theme |
| `neutral` | Greyscale, good for printing |
| `base` | **Only customisable theme.** Use with `themeVariables`. |

**Important:** Only the `base` theme supports `themeVariables` customisation.
The other four themes ignore `themeVariables` entirely. Colour values must be
hex format (`#ff0000`); colour names like `red` are not accepted.

---

## 10. Rendering Multiple Diagrams on One Page

Each diagram needs a unique `id` passed to `mermaid.render()`. Mermaid embeds
this ID in internal SVG elements (clip paths, gradients, markers). If two
diagrams share the same ID, their SVG definitions collide.

The `useMermaid` hook (Section 8) handles this automatically via `useId()`.

When rendering multiple diagrams in a list:

```tsx
{sections.map((section) => (
  <MermaidDiagram
    key={section.id}
    definition={section.source}
    title={section.title}
  />
))}
```

Each `MermaidDiagram` instance calls `useId()` independently, producing unique
IDs like `mermaid--r1-`, `mermaid--r2-`, etc.

---

## 11. SSR / Client-Side Considerations

| Concern | Detail |
|---|---|
| **No SSR support** | Mermaid requires a browser DOM to compute layout dimensions (via D3.js). It cannot render SVG on the server. |
| **Client-only rendering** | In SSR frameworks (Next.js, Remix), the Mermaid component must be loaded client-side only. With Vite + React Router (this project), the app is a client-side SPA, so no special handling is needed. |
| **Dynamic import** | For Next.js: `const Mermaid = dynamic(() => import('./Mermaid'), { ssr: false })`. Not needed in this project's Vite SPA. |
| **Hydration** | If SSR were used, the initial server render should show a fallback (loading placeholder or raw source). The client `useEffect` then replaces it with the rendered SVG. |
| **`useEffect` timing** | `mermaid.render()` is async. It runs after the component mounts and the browser has painted. The loading state covers the gap. |
| **DOM dependency** | `mermaid.render()` internally creates a temporary SVG element in the document to compute layout. It requires `document` and `window` to be available. |
| **Web Workers** | Mermaid cannot run in a Web Worker because it depends on DOM APIs. |

---

## 12. Complete Minimal Example

A self-contained example combining initialisation, rendering, error handling,
and dark theme:

```tsx
import { useEffect, useId, useRef, useState } from 'react'
import mermaid from 'mermaid'

// One-time global initialisation
mermaid.initialize({
  startOnLoad: false,
  theme: 'dark',
  securityLevel: 'strict',
  fontFamily: 'inherit',
})

interface DiagramProps {
  source: string
  title?: string
}

export function Diagram({ source, title }: DiagramProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const [error, setError] = useState<string | null>(null)
  const reactId = useId()
  const diagramId = `mermaid-${reactId.replace(/:/g, '-')}`

  useEffect(() => {
    let cancelled = false

    mermaid
      .render(diagramId, source)
      .then(({ svg, bindFunctions }) => {
        if (cancelled || !containerRef.current) return
        containerRef.current.innerHTML = svg
        bindFunctions?.(containerRef.current)
        setError(null)
      })
      .catch((err: unknown) => {
        if (cancelled) return
        setError(err instanceof Error ? err.message : 'Render failed')
      })

    return () => {
      cancelled = true
    }
  }, [source, diagramId])

  if (error) {
    return (
      <div>
        {title && <p>{title}</p>}
        <p>Failed to render diagram: {error}</p>
        <pre><code>{source}</code></pre>
      </div>
    )
  }

  return (
    <figure>
      {title && <figcaption>{title}</figcaption>}
      <div ref={containerRef} />
    </figure>
  )
}
```

---

## 13. API Quick Reference

```ts
import mermaid from 'mermaid'

// Configure (call once)
mermaid.initialize({ startOnLoad: false, theme: 'dark' })

// Validate syntax (does not render)
const parseResult = await mermaid.parse('graph TD; A-->B')
// Returns { diagramType: 'flowchart' } or throws

// Render to SVG string
const { svg, bindFunctions } = await mermaid.render('unique-id', 'graph TD; A-->B')
// svg: string — the SVG markup
// bindFunctions: optional callback for interactive diagrams

// Auto-render DOM elements (not recommended for React)
await mermaid.run({ querySelector: '.mermaid' })
await mermaid.run({ nodes: [element], suppressErrors: true })
```
