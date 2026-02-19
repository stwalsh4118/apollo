import { Link } from "react-router";
import type { TopicSummary } from "../api";
import DifficultyBadge from "./DifficultyBadge";

interface TopicCardProps {
  topic: TopicSummary;
}

export default function TopicCard({ topic }: TopicCardProps) {
  return (
    <Link
      to={`/topics/${topic.id}`}
      className="block rounded-lg border border-gray-200 bg-white p-6 shadow-sm transition-shadow hover:shadow-md"
    >
      <div className="flex items-start justify-between gap-2">
        <h2 className="text-lg font-semibold text-gray-900">{topic.title}</h2>
        {topic.difficulty && <DifficultyBadge difficulty={topic.difficulty} />}
      </div>
      {topic.description && (
        <p className="mt-2 line-clamp-2 text-sm text-gray-600">
          {topic.description}
        </p>
      )}
      <p className="mt-4 text-xs text-gray-500">
        {topic.module_count} {topic.module_count === 1 ? "module" : "modules"}
      </p>
    </Link>
  );
}
