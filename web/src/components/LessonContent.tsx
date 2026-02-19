import { useState } from "react";
import { useLessonDetail, useUpdateLessonProgress } from "../api";
import type { ConceptSummary, ProgressStatus } from "../api";
import ConceptChips from "./ConceptChips";
import ContentRenderer from "./content/ContentRenderer";
import ExerciseList from "./ExerciseList";
import ReviewQuestions from "./ReviewQuestions";

interface LessonContentProps {
  lessonId: string;
  topicId: string;
  lessonStatus?: ProgressStatus;
  lessonNotes?: string;
  concepts?: ConceptSummary[];
}

export default function LessonContent({
  lessonId,
  topicId,
  lessonStatus,
  lessonNotes,
  concepts,
}: LessonContentProps) {
  const { data: lesson, isLoading, isError, error } = useLessonDetail(lessonId);
  const completeMutation = useUpdateLessonProgress(topicId);
  const notesMutation = useUpdateLessonProgress(topicId);

  const [notes, setNotes] = useState(lessonNotes ?? "");
  const [notesSaved, setNotesSaved] = useState(false);

  const isCompleted = lessonStatus === "completed";
  const currentStatus = lessonStatus ?? "not_started";

  function handleMarkComplete() {
    completeMutation.mutate({
      lessonId,
      input: { status: "completed", notes },
    });
  }

  function handleSaveNotes() {
    setNotesSaved(false);
    notesMutation.mutate(
      { lessonId, input: { status: currentStatus, notes } },
      {
        onSuccess: () => {
          setNotesSaved(true);
          setTimeout(() => setNotesSaved(false), 2000);
        },
      },
    );
  }

  if (isLoading) {
    return (
      <div className="py-12 text-center text-gray-500">Loading lesson...</div>
    );
  }

  if (isError) {
    return (
      <div className="rounded-md bg-red-50 p-4 text-sm text-red-700">
        Failed to load lesson: {error.message}
      </div>
    );
  }

  if (!lesson) {
    return null;
  }

  return (
    <article>
      <h1 className="text-2xl font-bold text-gray-900">{lesson.title}</h1>
      {lesson.estimated_minutes && (
        <p className="mt-1 text-sm text-gray-500">
          {lesson.estimated_minutes} min read
        </p>
      )}

      {concepts && concepts.length > 0 && (
        <ConceptChips concepts={concepts} />
      )}

      {lesson.content && lesson.content.length > 0 && (
        <div className="mt-6">
          <ContentRenderer sections={lesson.content} />
        </div>
      )}

      {lesson.exercises && lesson.exercises.length > 0 && (
        <div className="mt-8">
          <ExerciseList exercises={lesson.exercises} />
        </div>
      )}

      {lesson.review_questions && lesson.review_questions.length > 0 && (
        <div className="mt-8">
          <ReviewQuestions questions={lesson.review_questions} />
        </div>
      )}

      {/* Personal notes */}
      <div className="mt-8">
        <h2 className="text-sm font-semibold text-gray-700">Personal Notes</h2>
        <textarea
          value={notes}
          onChange={(e) => setNotes(e.target.value)}
          placeholder="Add your notes for this lesson..."
          rows={4}
          className="mt-2 w-full resize-y rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 placeholder-gray-400 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
        />
        <div className="mt-2 flex items-center gap-3">
          <button
            onClick={handleSaveNotes}
            disabled={notesMutation.isPending}
            className="rounded-md bg-gray-800 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-gray-700 disabled:opacity-50"
          >
            {notesMutation.isPending ? "Saving..." : "Save Notes"}
          </button>
          {notesSaved && (
            <span className="text-sm text-green-600">Saved</span>
          )}
          {notesMutation.isError && (
            <span className="text-sm text-red-600">Failed to save notes</span>
          )}
        </div>
      </div>

      {/* Mark Complete button */}
      <div className="mt-10 border-t border-gray-200 pt-6">
        {isCompleted || completeMutation.isSuccess ? (
          <div className="flex items-center gap-2 text-sm font-medium text-green-600">
            <svg
              className="size-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M5 13l4 4L19 7"
              />
            </svg>
            Completed
          </div>
        ) : (
          <button
            onClick={handleMarkComplete}
            disabled={completeMutation.isPending}
            className="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
          >
            {completeMutation.isPending ? "Saving..." : "Mark Complete"}
          </button>
        )}
        {completeMutation.isError && (
          <p className="mt-2 text-sm text-red-600">
            Failed to update progress. Please try again.
          </p>
        )}
      </div>
    </article>
  );
}
