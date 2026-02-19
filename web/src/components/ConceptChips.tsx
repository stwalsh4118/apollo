import { Link } from "react-router";
import type { ConceptSummary } from "../api";

interface ConceptChipsProps {
  concepts: ConceptSummary[];
}

export default function ConceptChips({ concepts }: ConceptChipsProps) {
  if (concepts.length === 0) {
    return null;
  }

  return (
    <div className="mt-3 flex flex-wrap gap-2">
      {concepts.map((concept) => (
        <Link
          key={concept.id}
          to={`/concepts/${encodeURIComponent(concept.id)}`}
          className="inline-flex items-center rounded-full bg-blue-50 px-3 py-1 text-xs font-medium text-blue-700 transition-colors hover:bg-blue-100"
        >
          {concept.name}
        </Link>
      ))}
    </div>
  );
}
