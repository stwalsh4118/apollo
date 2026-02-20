import { useState } from "react";
import Markdown, { type Components } from "react-markdown";
import rehypeRaw from "rehype-raw";
import type { Exercise } from "../api";
import MarkdownCode from "./content/MarkdownCode";

const mdComponents: Components = {
  pre({ children }) {
    return <>{children}</>;
  },
  code({ className, children }) {
    const match = className?.match(/language-(\w+)/);
    const code = String(children).replace(/\n$/, "");

    if (match) {
      return <MarkdownCode code={code} language={match[1]} />;
    }

    return (
      <code className="rounded bg-gray-100 px-1.5 py-0.5 text-sm font-mono text-gray-800">
        {children}
      </code>
    );
  },
};

interface ExerciseBlockProps {
  exercise: Exercise;
  index: number;
}

export default function ExerciseBlock({ exercise, index }: ExerciseBlockProps) {
  const [hintsRevealed, setHintsRevealed] = useState(0);
  const hints = exercise.hints ?? [];
  const hasMoreHints = hintsRevealed < hints.length;

  return (
    <div className="rounded-lg border border-gray-200 bg-white p-5">
      <div className="mb-3 flex items-center gap-2">
        <span className="text-sm font-semibold text-gray-900">
          Exercise {index + 1}: {exercise.title}
        </span>
        <span className="rounded bg-indigo-100 px-2 py-0.5 text-xs font-medium text-indigo-700">
          {exercise.type}
        </span>
      </div>

      <div className="prose prose-sm prose-gray max-w-none text-gray-700">
        <Markdown rehypePlugins={[rehypeRaw]} components={mdComponents}>
          {exercise.instructions}
        </Markdown>
      </div>

      {exercise.environment && (
        <p className="mt-2 text-xs text-gray-500">
          <span className="font-medium">Environment:</span>{" "}
          {exercise.environment}
        </p>
      )}

      {exercise.success_criteria && exercise.success_criteria.length > 0 && (
        <div className="mt-4">
          <p className="text-xs font-medium text-gray-500">Success Criteria</p>
          <ul className="mt-1 list-inside list-disc space-y-1 text-sm text-gray-600">
            {exercise.success_criteria.map((criterion, i) => (
              <li key={i}>{criterion}</li>
            ))}
          </ul>
        </div>
      )}

      {hints.length > 0 && (
        <div className="mt-4 border-t border-gray-100 pt-3">
          <div className="space-y-2">
            {hints.slice(0, hintsRevealed).map((hint, i) => (
              <div
                key={i}
                className="prose prose-sm max-w-none rounded bg-amber-50 px-3 py-2 text-sm text-amber-800"
              >
                <span className="font-medium">Hint {i + 1}:</span>
                <Markdown rehypePlugins={[rehypeRaw]} components={mdComponents}>
                  {hint}
                </Markdown>
              </div>
            ))}
          </div>
          {hasMoreHints && (
            <button
              type="button"
              onClick={() => setHintsRevealed((prev) => prev + 1)}
              className="mt-2 text-sm font-medium text-indigo-600 hover:text-indigo-800"
            >
              Show hint {hintsRevealed + 1} of {hints.length}
            </button>
          )}
        </div>
      )}
    </div>
  );
}
