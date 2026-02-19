import { useCallback, useState } from "react";
import type { CodeSection as CodeSectionType } from "../../api";
import { useShikiHighlighter } from "../../hooks/useShikiHighlighter";

const COPY_FEEDBACK_MS = 2000;

interface CodeSectionProps {
  section: CodeSectionType;
}

export default function CodeSection({ section }: CodeSectionProps) {
  const { html, isLoading } = useShikiHighlighter(
    section.code,
    section.language,
  );
  const [copied, setCopied] = useState(false);

  const handleCopy = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(section.code);
      setCopied(true);
      setTimeout(() => setCopied(false), COPY_FEEDBACK_MS);
    } catch {
      // Clipboard API not available — ignore silently
    }
  }, [section.code]);

  return (
    <div>
      {/* Header bar */}
      <div className="flex items-center justify-between rounded-t-lg bg-gray-800 px-4 py-2">
        <div className="flex items-center gap-2">
          {section.title && (
            <span className="text-sm font-medium text-gray-200">
              {section.title}
            </span>
          )}
          <span className="rounded bg-gray-700 px-2 py-0.5 text-xs text-gray-400">
            {section.language}
          </span>
        </div>
        <button
          onClick={handleCopy}
          className="rounded p-1 text-gray-400 transition-colors hover:bg-gray-700 hover:text-gray-200"
          aria-label={copied ? "Copied" : "Copy code"}
        >
          {copied ? (
            <svg
              className="size-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M5 13l4 4L19 7"
              />
            </svg>
          ) : (
            <svg
              className="size-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
              />
            </svg>
          )}
        </button>
      </div>

      {/* Code block — HTML is produced by Shiki's codeToHtml which escapes all user content */}
      {!isLoading && html ? (
        <div
          className="overflow-x-auto rounded-b-lg text-sm [&>pre]:p-4"
          dangerouslySetInnerHTML={{ __html: html }}
        />
      ) : (
        <pre className="overflow-x-auto rounded-b-lg bg-gray-900 p-4 text-sm text-gray-300">
          <code>{section.code}</code>
        </pre>
      )}

      {/* Explanation */}
      {section.explanation && (
        <p className="mt-2 text-sm text-gray-600">{section.explanation}</p>
      )}
    </div>
  );
}
