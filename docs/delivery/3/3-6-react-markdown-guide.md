# react-markdown Usage Guide

**Date:** 2026-02-18
**Package:** `react-markdown` v10.1.0
**Docs:** <https://github.com/remarkjs/react-markdown>

---

## 1. Installation

The project already has these packages installed:

```bash
pnpm add react-markdown rehype-raw
```

Versions in `web/package.json`:

| Package | Version |
|---------|---------|
| `react-markdown` | `^10.1.0` |
| `rehype-raw` | `^7.0.0` |
| `@tailwindcss/typography` | `^0.5.19` |

> `react-markdown` v10 is **ESM-only** and requires **React 18+** (project uses React 19).

---

## 2. Basic Usage

```tsx
import Markdown from 'react-markdown';

function CourseDescription({ content }: { content: string }) {
  return <Markdown>{content}</Markdown>;
}
```

The `<Markdown>` component accepts a markdown string as `children` and renders it as React elements. It does **not** use `dangerouslySetInnerHTML` -- it builds a virtual DOM from a syntax tree, so React only patches what changed.

### Key Props

| Prop | Type | Description |
|------|------|-------------|
| `children` | `string` | Markdown source to render |
| `components` | `Components` | Map HTML tag names to custom React components |
| `remarkPlugins` | `Plugin[]` | Remark plugins (transform markdown AST) |
| `rehypePlugins` | `Plugin[]` | Rehype plugins (transform HTML AST) |
| `allowedElements` | `string[]` | Whitelist of allowed HTML tag names |
| `disallowedElements` | `string[]` | Blacklist of HTML tag names |
| `skipHtml` | `boolean` | Ignore raw HTML in markdown (default: `false`) |
| `urlTransform` | `UrlTransform` | Custom URL transformation function |

---

## 3. Using rehype-raw for HTML-in-Markdown

By default, `react-markdown` escapes raw HTML in markdown for security. The `rehype-raw` plugin re-parses and renders it.

```tsx
import Markdown from 'react-markdown';
import rehypeRaw from 'rehype-raw';

const content = `
# Heading

<div class="callout">
  This contains **bold markdown** inside raw HTML.
</div>
`;

function RichContent() {
  return (
    <Markdown rehypePlugins={[rehypeRaw]}>
      {content}
    </Markdown>
  );
}
```

**Important:** Block-level HTML that contains markdown must have blank lines around it to be parsed correctly (CommonMark rule). Only use `rehype-raw` with **trusted content** -- it enables arbitrary HTML rendering.

---

## 4. Custom Component Overrides

The `components` prop maps HTML element names to React components. Every component receives `node` (the AST node), `children`, and the standard HTML attributes as props.

### Overridable elements

Markdown maps to these HTML elements: `a`, `blockquote`, `br`, `code`, `em`, `h1`-`h6`, `hr`, `img`, `li`, `ol`, `p`, `pre`, `strong`, `ul`.

### Example: Styling headings, links, and code blocks

```tsx
import Markdown, { type Components } from 'react-markdown';
import rehypeRaw from 'rehype-raw';

const components: Components = {
  // Headings
  h1: ({ children }) => (
    <h1 className="text-3xl font-bold mt-8 mb-4">{children}</h1>
  ),
  h2: ({ children }) => (
    <h2 className="text-2xl font-semibold mt-6 mb-3">{children}</h2>
  ),
  h3: ({ children }) => (
    <h3 className="text-xl font-medium mt-4 mb-2">{children}</h3>
  ),

  // Links -- open external in new tab
  a: ({ href, children }) => (
    <a
      href={href}
      target="_blank"
      rel="noopener noreferrer"
      className="text-blue-600 underline hover:text-blue-800"
    >
      {children}
    </a>
  ),

  // Code blocks -- distinguish inline vs fenced
  code: ({ className, children, ...rest }) => {
    const match = /language-(\w+)/.exec(className || '');
    const isBlock = Boolean(match);

    if (isBlock) {
      return (
        <code
          className={`block bg-gray-900 text-gray-100 rounded-lg p-4 overflow-x-auto text-sm ${className}`}
          {...rest}
        >
          {children}
        </code>
      );
    }

    return (
      <code
        className="bg-gray-100 text-pink-600 rounded px-1.5 py-0.5 text-sm"
        {...rest}
      >
        {children}
      </code>
    );
  },

  // Block quotes
  blockquote: ({ children }) => (
    <blockquote className="border-l-4 border-blue-500 pl-4 italic text-gray-600 my-4">
      {children}
    </blockquote>
  ),
};

function StyledMarkdown({ content }: { content: string }) {
  return (
    <Markdown
      rehypePlugins={[rehypeRaw]}
      components={components}
    >
      {content}
    </Markdown>
  );
}
```

---

## 5. Tailwind CSS Typography (`prose`) Integration

The project uses Tailwind CSS v4 with `@tailwindcss/typography` already configured:

```css
/* src/index.css */
@import "tailwindcss";
@plugin "@tailwindcss/typography";
```

### Using `prose` classes (recommended approach)

Wrapping `<Markdown>` in a `prose` container gives you well-styled typography with zero custom components:

```tsx
import Markdown from 'react-markdown';
import rehypeRaw from 'rehype-raw';

function CourseContent({ markdown }: { markdown: string }) {
  return (
    <div className="prose prose-lg max-w-none">
      <Markdown rehypePlugins={[rehypeRaw]}>
        {markdown}
      </Markdown>
    </div>
  );
}
```

### Useful `prose` modifiers

| Class | Effect |
|-------|--------|
| `prose` | Base typographic styles |
| `prose-sm` / `prose-lg` / `prose-xl` | Size variants |
| `prose-gray` / `prose-slate` / `prose-zinc` | Gray scale themes |
| `prose-invert` | Dark mode (light text on dark background) |
| `max-w-none` | Remove the default max-width constraint |

### Fine-tuning with element modifiers

Tailwind typography supports per-element overrides via modifier classes:

```tsx
<div className="prose prose-headings:text-blue-900 prose-a:text-blue-600 prose-a:no-underline hover:prose-a:underline prose-code:text-pink-600 prose-blockquote:border-blue-500">
  <Markdown rehypePlugins={[rehypeRaw]}>
    {markdown}
  </Markdown>
</div>
```

### Dark mode

```tsx
<div className="prose dark:prose-invert">
  <Markdown>{markdown}</Markdown>
</div>
```

### Combining `prose` with custom components

You can use `prose` for baseline styles and override specific elements with `components` when you need behavior changes (e.g., opening links in new tabs):

```tsx
import Markdown, { type Components } from 'react-markdown';
import rehypeRaw from 'rehype-raw';

const components: Components = {
  // Only override elements where you need custom behavior
  a: ({ href, children }) => (
    <a href={href} target="_blank" rel="noopener noreferrer">
      {children}
    </a>
  ),
};

function CourseContent({ markdown }: { markdown: string }) {
  return (
    <div className="prose prose-lg max-w-none dark:prose-invert">
      <Markdown rehypePlugins={[rehypeRaw]} components={components}>
        {markdown}
      </Markdown>
    </div>
  );
}
```

This lets Tailwind typography handle all visual styling while `components` handles behavioral overrides.

---

## TypeScript Exports

The package exports these types:

| Export | Description |
|--------|-------------|
| `Markdown` | Main synchronous component (default export) |
| `Components` | Type for the component override map |
| `Options` | Main configuration/props type |
| `ExtraProps` | Additional props passed to custom components |
| `UrlTransform` | Type for URL transformation function |
| `defaultUrlTransform` | Built-in URL sanitization function |

---

## Sources

- [react-markdown on GitHub](https://github.com/remarkjs/react-markdown)
- [react-markdown on npm](https://www.npmjs.com/package/react-markdown)
- [rehype-raw on npm](https://www.npmjs.com/package/rehype-raw)
- [@tailwindcss/typography on GitHub](https://github.com/tailwindlabs/tailwindcss-typography)
