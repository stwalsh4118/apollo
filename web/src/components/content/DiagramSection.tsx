import type { DiagramSection as DiagramSectionType } from "../../api";
import { useMermaid } from "../../hooks/useMermaid";

interface DiagramSectionProps {
  section: DiagramSectionType;
}

function MermaidDiagram({ source }: { source: string }) {
  const { svg, error, isLoading } = useMermaid(source);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center rounded-lg bg-gray-50 p-8 text-sm text-gray-400">
        Loading diagram…
      </div>
    );
  }

  if (error) {
    return (
      <div className="rounded-lg border border-red-200 bg-red-50 p-4">
        <p className="mb-2 text-sm font-medium text-red-700">
          Diagram render error
        </p>
        <pre className="overflow-x-auto text-xs text-red-600">{error}</pre>
        <details className="mt-3">
          <summary className="cursor-pointer text-xs text-gray-500">
            Show raw source
          </summary>
          <pre className="mt-2 overflow-x-auto rounded bg-gray-100 p-3 text-xs text-gray-700">
            {source}
          </pre>
        </details>
      </div>
    );
  }

  if (!svg) {
    return null;
  }

  return (
    /* Safe: SVG is produced by Mermaid's render with securityLevel: "strict" */
    <div
      className="flex justify-center overflow-x-auto [&>svg]:max-w-full"
      dangerouslySetInnerHTML={{ __html: svg }}
    />
  );
}

function ImageDiagram({ source, title }: { source: string; title?: string }) {
  return (
    <img
      src={source}
      alt={title ?? "Diagram"}
      className="rounded-lg"
      loading="lazy"
    />
  );
}

export default function DiagramSection({ section }: DiagramSectionProps) {
  return (
    <figure>
      {section.title && (
        <figcaption className="mb-2 text-sm font-medium text-gray-700">
          {section.title}
        </figcaption>
      )}
      {section.format === "mermaid" ? (
        <MermaidDiagram source={section.source} />
      ) : (
        <ImageDiagram source={section.source} title={section.title} />
      )}
    </figure>
  );
}
