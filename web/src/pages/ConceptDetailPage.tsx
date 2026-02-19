import { useParams } from "react-router";

export default function ConceptDetailPage() {
  const { id } = useParams<{ id: string }>();

  return (
    <div className="mx-auto max-w-4xl px-6 py-12 lg:px-12">
      <h1 className="text-2xl font-bold text-gray-900">Concept Detail</h1>
      <p className="mt-2 text-gray-500">
        Concept <span className="font-mono text-gray-700">{id}</span> — full
        detail page will be available in the Knowledge Wiki (PBI 11).
      </p>
    </div>
  );
}
