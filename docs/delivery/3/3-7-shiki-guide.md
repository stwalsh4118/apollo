# Shiki Syntax Highlighting Guide

> **Date:** 2026-02-18
> **Package:** `shiki` v3.x (latest: 3.22.0)
> **Docs:** <https://shiki.style/>
> **Repository:** <https://github.com/shikijs/shiki>

---

## 1. Installation

```bash
pnpm add shiki
```

No additional peer dependencies required. Shiki ships with bundled themes and
language grammars that are lazy-loaded on demand.

---

## 2. Creating a Highlighter Instance (Async / Lazy)

Shiki highlighters are **async** because they load grammar and theme data at
creation time.

### Shorthand (easiest, lazy-loads automatically)

```ts
import { codeToHtml } from 'shiki'

// Each call lazy-loads the requested theme + language on first use
const html = await codeToHtml('const x = 1', {
  lang: 'typescript',
  theme: 'github-dark',
})
```

The shorthand manages a singleton highlighter internally. Suitable when you only
need HTML output and do not need fine-grained control.

### Explicit Highlighter Instance (recommended for apps)

```ts
import { createHighlighter } from 'shiki'

// Create once, reuse everywhere
const highlighter = await createHighlighter({
  themes: ['github-dark'],
  langs: ['typescript', 'python', 'json'],
})

// Synchronous after creation -- no await needed
const html = highlighter.codeToHtml('const x = 1', {
  lang: 'typescript',
  theme: 'github-dark',
})
```

**Important:** Cache the highlighter as a singleton. `createHighlighter` is
expensive; never call it inside render loops or per-request handlers.

Load additional themes/languages dynamically after creation:

```ts
await highlighter.loadTheme('github-light')
await highlighter.loadLanguage('css')
```

Call `highlighter.dispose()` when the highlighter is no longer needed to free
resources.

---

## 3. Getting HTML Output

`codeToHtml` returns a self-contained HTML string with inline styles (no
external CSS required):

```ts
const html = await codeToHtml('print("hello")', {
  lang: 'python',
  theme: 'github-dark',
})

// Returns something like:
// <pre class="shiki github-dark" style="background-color:#24292e;color:#e1e4e8" ...>
//   <code><span class="line"><span style="color:#79B8FF">print</span>(...)</span></code>
// </pre>
```

---

## 4. Available Themes

Shiki bundles a large set of VS Code-compatible themes. GitHub-related themes:

| Theme ID | Description |
|---|---|
| `github-dark` | GitHub dark theme |
| `github-dark-default` | GitHub dark (default variant) |
| `github-dark-dimmed` | GitHub dark dimmed |
| `github-dark-high-contrast` | GitHub dark high contrast |
| `github-light` | GitHub light theme |
| `github-light-default` | GitHub light (default variant) |
| `github-light-high-contrast` | GitHub light high contrast |

Other popular themes include `nord`, `dracula`, `catppuccin-mocha`,
`one-dark-pro`, `vitesse-dark`, `material-theme-ocean`, and many more.

Full list: <https://shiki.style/themes>

---

## 5. Fine-Grained Bundles (Smaller Bundle Size)

The default `shiki` import includes lazy references to **all** bundled
languages and themes. For web apps where bundle size matters, use the
fine-grained approach to include only what you need.

### Using `shiki/core` + Explicit Imports

```ts
import { createHighlighterCore } from 'shiki/core'
import { createJavaScriptRegexEngine } from 'shiki/engine/javascript'

const highlighter = await createHighlighterCore({
  themes: [
    import('@shikijs/themes/github-dark'),
  ],
  langs: [
    import('@shikijs/langs/typescript'),
    import('@shikijs/langs/python'),
    import('@shikijs/langs/json'),
  ],
  engine: createJavaScriptRegexEngine(),
})
```

Key points:

- **`shiki/core`** ships with zero bundled themes or languages.
- **`@shikijs/themes/<name>`** and **`@shikijs/langs/<name>`** are direct
  imports -- only what you import ends up in your bundle.
- **`shiki/engine/javascript`** uses native JS `RegExp` instead of the
  Oniguruma WASM binary (~4 MB), significantly reducing bundle size. It
  supports 97%+ of built-in languages.
- **`shiki/engine/oniguruma`** is the default engine with full compatibility
  but requires loading a WASM file.

### Engine Comparison

| Engine | Bundle Impact | Compatibility | Best For |
|---|---|---|---|
| `createJavaScriptRegexEngine()` | Small (no WASM) | ~97% of languages | Browser / client-side |
| `createOnigurumaEngine(import('shiki/wasm'))` | +4 MB WASM | 100% | Node.js / SSR / build-time |

---

## 6. React Integration: `codeToHtml` + `dangerouslySetInnerHTML`

### Client-Side Component Pattern

```tsx
import { useEffect, useState } from 'react'
import { codeToHtml } from 'shiki'

interface CodeBlockProps {
  code: string
  language: string
}

export function CodeBlock({ code, language }: CodeBlockProps) {
  const [html, setHtml] = useState<string>('')

  useEffect(() => {
    let cancelled = false

    codeToHtml(code, {
      lang: language,
      theme: 'github-dark',
    }).then((result) => {
      if (!cancelled) {
        setHtml(result)
      }
    })

    return () => {
      cancelled = true
    }
  }, [code, language])

  if (!html) {
    // Fallback while shiki loads
    return (
      <pre>
        <code>{code}</code>
      </pre>
    )
  }

  return <div dangerouslySetInnerHTML={{ __html: html }} />
}
```

### With a Cached Highlighter Singleton (Recommended)

For better performance, create the highlighter once and reuse it:

```tsx
import { useEffect, useState } from 'react'
import type { Highlighter } from 'shiki'
import { createHighlighter } from 'shiki'

// Singleton promise -- created once, resolved once
let highlighterPromise: Promise<Highlighter> | null = null

function getHighlighter(): Promise<Highlighter> {
  if (!highlighterPromise) {
    highlighterPromise = createHighlighter({
      themes: ['github-dark'],
      langs: ['typescript', 'python', 'json', 'bash'],
    })
  }
  return highlighterPromise
}

interface CodeBlockProps {
  code: string
  language: string
}

export function CodeBlock({ code, language }: CodeBlockProps) {
  const [html, setHtml] = useState<string>('')

  useEffect(() => {
    let cancelled = false

    getHighlighter().then((highlighter) => {
      if (!cancelled) {
        const result = highlighter.codeToHtml(code, {
          lang: language,
          theme: 'github-dark',
        })
        setHtml(result)
      }
    })

    return () => {
      cancelled = true
    }
  }, [code, language])

  if (!html) {
    return (
      <pre>
        <code>{code}</code>
      </pre>
    )
  }

  return <div dangerouslySetInnerHTML={{ __html: html }} />
}
```

---

## 7. SSR / Client-Side Considerations

| Concern | Detail |
|---|---|
| **Async loading** | `createHighlighter` and `codeToHtml` are async. In SSR frameworks (Next.js App Router), you can `await` them directly in server components. In client-only React, use `useEffect` + state as shown above. |
| **WASM in SSR** | The Oniguruma engine loads WASM, which works in Node.js but can cause issues in edge runtimes. Use the JS engine for edge/serverless. |
| **Bundle size** | The full `shiki` bundle lazily references all languages (~150+). Use fine-grained imports (section 5) for client-side bundles. |
| **Hydration** | If highlighting only on the client, the initial render shows the plain-text fallback. This avoids hydration mismatches. |
| **Thread blocking** | Highlighting is CPU-intensive. For large code blocks, consider offloading to a Web Worker. |
| **`dangerouslySetInnerHTML`** | Safe when the input to `codeToHtml` is trusted or sanitised. Shiki produces inline-styled `<span>` elements; it does not pass through raw user HTML. |
