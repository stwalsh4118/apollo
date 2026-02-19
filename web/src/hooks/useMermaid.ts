import { useEffect, useState } from "react";

const MERMAID_THEME = "dark";

let mermaidPromise: Promise<typeof import("mermaid")> | null = null;
let mermaidInitialised = false;
let idCounter = 0;

function getMermaid(): Promise<typeof import("mermaid")> {
  if (!mermaidPromise) {
    mermaidPromise = import("mermaid");
  }
  return mermaidPromise;
}

function nextId(): string {
  idCounter += 1;
  return `mermaid-diagram-${idCounter}`;
}

interface MermaidResult {
  svg: string | null;
  error: string | null;
  isLoading: boolean;
}

export function useMermaid(source: string): MermaidResult {
  const [result, setResult] = useState<{
    svg: string | null;
    error: string | null;
    source: string;
  } | null>(null);

  useEffect(() => {
    let cancelled = false;
    const renderId = nextId();

    getMermaid()
      .then(async (mod) => {
        if (cancelled) return;
        const mermaid = mod.default;

        if (!mermaidInitialised) {
          mermaid.initialize({
            startOnLoad: false,
            theme: MERMAID_THEME,
            securityLevel: "strict",
          });
          mermaidInitialised = true;
        }

        const { svg } = await mermaid.render(renderId, source);
        if (!cancelled) {
          setResult({ svg, error: null, source });
        }
      })
      .catch((err: unknown) => {
        if (cancelled) return;
        const message =
          err instanceof Error ? err.message : "Failed to render diagram";
        setResult({ svg: null, error: message, source });
      });

    return () => {
      cancelled = true;
    };
  }, [source]);

  const isStale = result === null || result.source !== source;

  return {
    svg: isStale ? null : result.svg,
    error: isStale ? null : result.error,
    isLoading: isStale,
  };
}
