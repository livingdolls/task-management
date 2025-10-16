import type { TTask } from "../types/task";
import { TaskCard } from "./TaskCard";

interface TaskListProps {
  tasks: TTask[];
  onEditTask: (task: TTask) => void;
  onDeleteTask: (id: number) => void;
  isDeleting?: boolean;
}

export const TaskList = ({
  tasks,
  onEditTask,
  onDeleteTask,
  isDeleting,
}: TaskListProps) => {
  if (tasks.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="text-gray-500 text-lg mb-2">No tasks found</div>
        <div className="text-gray-400 text-sm">
          Create your first task to get started!
        </div>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {tasks.map((task, index) => (
        <TaskCard
          key={task.id}
          task={task}
          index={index}
          onEdit={onEditTask}
          onDelete={onDeleteTask}
          isDeleting={isDeleting}
        />
      ))}
    </div>
  );
};
