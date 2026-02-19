import { useEffect, useState } from "react";
import type { HighlighterCore } from "shiki";

const THEME = "github-dark";

const BUNDLED_LANGUAGES = [
  "bash",
  "json",
  "yaml",
  "go",
  "javascript",
  "typescript",
] as const;

let highlighterPromise: Promise<HighlighterCore> | null = null;

function getHighlighter(): Promise<HighlighterCore> {
  if (!highlighterPromise) {
    highlighterPromise = import("shiki")
      .then((shiki) =>
        shiki.createHighlighter({
          themes: [THEME],
          langs: [...BUNDLED_LANGUAGES],
        }),
      )
      .catch((err) => {
        highlighterPromise = null;
        throw err;
      });
  }
  return highlighterPromise;
}

interface CachedResult {
  html: string | null;
  code: string;
  language: string;
}

interface HighlightResult {
  html: string | null;
  isLoading: boolean;
}

export function useShikiHighlighter(
  code: string,
  language: string,
): HighlightResult {
  const [result, setResult] = useState<CachedResult | null>(null);

  useEffect(() => {
    let cancelled = false;

    getHighlighter()
      .then((highlighter) => {
        if (cancelled) return;
        const loadedLangs = highlighter.getLoadedLanguages();
        const lang = loadedLangs.includes(language) ? language : "text";
        const highlighted = highlighter.codeToHtml(code, {
          lang,
          theme: THEME,
        });
        setResult({ html: highlighted, code, language });
      })
      .catch(() => {
        if (cancelled) return;
        setResult({ html: null, code, language });
      });

    return () => {
      cancelled = true;
    };
  }, [code, language]);

  const isStale =
    result === null || result.code !== code || result.language !== language;

  return { html: isStale ? null : result.html, isLoading: isStale };
}
