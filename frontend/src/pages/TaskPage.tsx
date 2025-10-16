import { useState, useCallback } from "react";
import { useAuthStore } from "../store/useAuthStore";
import { useCreateTaskStore, useUiStore } from "../store/useUiStore";
import { UpdateTaskModal, type TaskUdate } from "../components/UpdateTaskModal";
import { CreateTaskModal } from "../components/CreateTaskModal";
import { TaskFilters } from "../components/TaskFilters";
import { TaskList } from "../components/TaskList";
import { PageHeader } from "../components/PageHeader";
import { ErrorDisplay } from "../components/ErrorDisplay";
import { LoadingSpinner } from "../components/LoadingSpinner";
import {
  useTasks,
  type TaskFilter,
  type CreateTaskData,
} from "../hooks/useTasks";
import { notifications } from "../utils/notifications";
import type { TTask } from "../types/task";

export const TaskPage = () => {
  const [filters, setFilters] = useState<TaskFilter>({});
  const [selectedTask, setSelectedTask] = useState<TaskUdate>({
    id: 0,
    title: "",
    description: "",
    status: "To Do",
    deadline: new Date(),
  });

  const { user } = useAuthStore();
  const { setUpdateModal } = useUiStore();
  const { setCreateModal } = useCreateTaskStore();

  const {
    tasks,
    isLoading,
    error,
    updateTask,
    createTask,
    deleteTask,
    isDeleting,
  } = useTasks(filters);

  const handleFiltersChange = useCallback((newFilters: TaskFilter) => {
    setFilters(newFilters);
  }, []);

  const handleUpdateTask = useCallback(
    (updatedTask: TaskUdate) => {
      try {
        updateTask(updatedTask);
        notifications.taskUpdated();
      } catch (error) {
        notifications.taskUpdateFailed((error as Error).message);
      }
    },
    [updateTask]
  );

  const handleEditTask = useCallback(
    (task: TTask) => {
      const taskData: TaskUdate = {
        id: task.id,
        title: task.title,
        description: task.description,
        status: task.status,
        deadline: task.deadline ? new Date(task.deadline) : undefined,
      };

      setSelectedTask(taskData);
      setUpdateModal(true);
    },
    [setUpdateModal]
  );

  const handleCreateTask = useCallback(
    (taskData: CreateTaskData) => {
      try {
        createTask(taskData);
        setCreateModal(false);
        notifications.taskCreated();
      } catch (error) {
        notifications.taskCreateFailed((error as Error).message);
      }
    },
    [createTask, setCreateModal]
  );

  const handleDeleteTask = useCallback(
    (id: number) => {
      if (window.confirm("Are you sure you want to delete this task?")) {
        try {
          deleteTask(id);
          notifications.taskDeleted();
        } catch (error) {
          notifications.taskDeleteFailed((error as Error).message);
        }
      }
    },
    [deleteTask]
  );

  const handleOpenCreateModal = useCallback(() => {
    setCreateModal(true);
  }, [setCreateModal]);

  if (isLoading) {
    return (
      <div className="container mx-auto p-4">
        <LoadingSpinner />
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto p-4">
        <ErrorDisplay
          message={(error as Error).message}
          onRetry={() => window.location.reload()}
        />
      </div>
    );
  }

  return (
    <div className="container mx-auto p-4">
      <PageHeader userName={user?.name} onCreateTask={handleOpenCreateModal} />

      <TaskFilters filters={filters} onFiltersChange={handleFiltersChange} />

      <TaskList
        tasks={tasks}
        onEditTask={handleEditTask}
        onDeleteTask={handleDeleteTask}
        isDeleting={isDeleting}
      />

      <UpdateTaskModal task={selectedTask} onUpdate={handleUpdateTask} />
      <CreateTaskModal onSubmit={handleCreateTask} />
    </div>
  );
};
