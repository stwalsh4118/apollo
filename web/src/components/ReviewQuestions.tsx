import type { ReviewQuestion } from "../api";

interface ReviewQuestionsProps {
  questions: ReviewQuestion[];
}

export default function ReviewQuestions({ questions }: ReviewQuestionsProps) {
  if (questions.length === 0) {
    return null;
  }

  return (
    <details className="rounded-lg border border-gray-200 bg-white">
      <summary className="cursor-pointer px-5 py-4 text-lg font-bold text-gray-900">
        Review Questions ({questions.length})
      </summary>
      <div className="space-y-3 px-5 pb-5">
        {questions.map((q, index) => (
          <details key={index} className="rounded border border-gray-100 bg-gray-50">
            <summary className="cursor-pointer px-4 py-3 text-sm font-medium text-gray-800">
              {index + 1}. {q.question}
            </summary>
            <div className="border-t border-gray-100 px-4 py-3 text-sm text-gray-700">
              {q.answer}
            </div>
          </details>
        ))}
      </div>
    </details>
  );
}
