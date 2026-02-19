import { useLessonDetail } from "../api";
import ContentRenderer from "./content/ContentRenderer";
import ExerciseList from "./ExerciseList";
import ReviewQuestions from "./ReviewQuestions";

interface LessonContentProps {
  lessonId: string;
}

export default function LessonContent({ lessonId }: LessonContentProps) {
  const { data: lesson, isLoading, isError, error } = useLessonDetail(lessonId);

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
    </article>
  );
}
