import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  CreateTaskRepository,
  DeleteTaskRepository,
  TaskRepository,
  UpdateTaskRepository,
} from "../repository/task_repository";
import type { TaskUdate } from "../components/UpdateTaskModal";
import type { TTaskStatus } from "../types/task";

export type TaskFilter = {
  status?: string;
  deadline?: string;
};

export type CreateTaskData = {
  title: string;
  description: string;
  status: TTaskStatus;
  deadline?: Date | null;
};

export const useTasks = (filters: TaskFilter) => {
  const queryClient = useQueryClient();

  const {
    data: tasksData,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["tasks", filters],
    queryFn: () => TaskRepository(filters),
  });

  const updateTaskMutation = useMutation({
    mutationFn: (taskData: TaskUdate) =>
      UpdateTaskRepository(taskData.id, taskData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
    },
    onError: (error) => {
      throw new Error((error as Error).message || "Failed to update task");
    },
  });

  const createTaskMutation = useMutation({
    mutationFn: CreateTaskRepository,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
    },
    onError: (error) => {
      throw new Error((error as Error).message || "Failed to create task");
    },
  });

  const deleteTaskMutation = useMutation({
    mutationFn: DeleteTaskRepository,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
    },
    onError: (error) => {
      throw new Error((error as Error).message || "Failed to delete task");
    },
  });

  return {
    tasks: tasksData?.data || [],
    isLoading,
    error,
    updateTask: updateTaskMutation.mutate,
    createTask: createTaskMutation.mutate,
    deleteTask: deleteTaskMutation.mutate,
    isUpdating: updateTaskMutation.isPending,
    isCreating: createTaskMutation.isPending,
    isDeleting: deleteTaskMutation.isPending,
  };
};
