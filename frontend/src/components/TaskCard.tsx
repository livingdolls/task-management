import type { TTask } from "../types/task";
import {
  TASK_CARD_COLORS,
  getStatusStyles,
  formatDate,
} from "../utils/taskUtils";

interface TaskCardProps {
  task: TTask;
  index: number;
  onEdit: (task: TTask) => void;
  onDelete: (id: number) => void;
  isDeleting?: boolean;
}

export const TaskCard = ({
  task,
  index,
  onEdit,
  onDelete,
  isDeleting,
}: TaskCardProps) => {
  const cardColor = TASK_CARD_COLORS[index % TASK_CARD_COLORS.length];
  const statusStyles = getStatusStyles(task.status);

  return (
    <div
      className={`${cardColor} border rounded-lg p-6 shadow-md hover:shadow-lg transition-shadow`}
    >
      <div className="mb-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-2">
          {task.title}
        </h3>
        <p className="text-gray-600 text-sm mb-3 line-clamp-3">
          {task.description}
        </p>

        <div className="flex items-center justify-between mb-4">
          <span
            className={`px-3 py-1 text-xs font-semibold rounded-full ${statusStyles}`}
          >
            {task.status}
          </span>

          <span className="text-sm text-gray-700">
            {formatDate(task.deadline)}
          </span>
        </div>
      </div>

      <div className="flex justify-end space-x-2">
        <button
          onClick={() => onEdit(task)}
          className="bg-indigo-500 hover:bg-indigo-600 text-white px-4 py-2 rounded text-sm font-medium transition-colors disabled:opacity-50"
          disabled={isDeleting}
        >
          Edit
        </button>
        <button
          onClick={() => onDelete(task.id)}
          className="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded text-sm font-medium transition-colors disabled:opacity-50"
          disabled={isDeleting}
        >
          {isDeleting ? "Deleting..." : "Delete"}
        </button>
      </div>
    </div>
  );
};
