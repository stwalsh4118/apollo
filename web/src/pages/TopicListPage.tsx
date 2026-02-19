import { useTopics } from "../api";
import TopicCard from "../components/TopicCard";

export default function TopicListPage() {
  const { data: topics, isLoading, isError, error } = useTopics();

  return (
    <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <h1 className="text-2xl font-bold text-gray-900">Topics</h1>
      <p className="mt-1 text-gray-600">Select a topic to begin learning.</p>

      {isLoading && (
        <div className="mt-8 text-center text-gray-500">Loading topics...</div>
      )}

      {isError && (
        <div className="mt-8 rounded-md bg-red-50 p-4 text-sm text-red-700">
          Failed to load topics: {error.message}
        </div>
      )}

      {topics && topics.length === 0 && (
        <div className="mt-8 text-center text-gray-500">
          No topics yet. Generate a curriculum to get started.
        </div>
      )}

      {topics && topics.length > 0 && (
        <div className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {topics.map((topic) => (
            <TopicCard key={topic.id} topic={topic} />
          ))}
        </div>
      )}
    </div>
  );
}
