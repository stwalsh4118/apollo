import type { Exercise } from "../api";
import ExerciseBlock from "./ExerciseBlock";

interface ExerciseListProps {
  exercises: Exercise[];
}

export default function ExerciseList({ exercises }: ExerciseListProps) {
  if (exercises.length === 0) {
    return null;
  }

  return (
    <section>
      <h2 className="mb-4 text-lg font-bold text-gray-900">Exercises</h2>
      <div className="space-y-4">
        {exercises.map((exercise, index) => (
          <ExerciseBlock key={index} exercise={exercise} index={index} />
        ))}
      </div>
    </section>
  );
}
