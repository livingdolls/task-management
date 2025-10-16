import type { TaskUdate } from "../components/UpdateTaskModal";
import axiosClient from "../lib/axios";
import type { TaskFilter } from "../pages/TaskPage";
import type { TApiResponse } from "../types/api";
import type { TTask } from "../types/task";

export const TaskRepository = async (
  arg: TaskFilter
): Promise<TApiResponse<TTask[]>> => {
  const params = new URLSearchParams();

  if (arg.status && arg.status !== "All") {
    params.append("status", arg.status);
  }

  if (arg.deadline) {
    params.append("deadline", arg.deadline);
  }

  const url = params.toString() ? `/tasks/?${params.toString()}` : "/tasks/";
  const res = await axiosClient.get(url);

  if (res.status !== 200 && res.status !== 201) {
    throw new Error(res.data?.errors || "Failed to fetch tasks");
  }

  return res.data;
};

export const UpdateTaskRepository = async (
  id: number,
  data: TaskUdate
): Promise<TApiResponse<TTask>> => {
  const res = await axiosClient.put(`/tasks/${id}`, data);

  if (res.status !== 200 && res.status !== 201) {
    throw new Error(res.data?.errors || "Failed to update task");
  }

  return res.data;
};

export const DeleteTaskRepository = async (id: number) => {
  const res = await axiosClient.delete(`/tasks/${id}`);

  if (res.status !== 200 && res.status !== 201 && res.status !== 204) {
    throw new Error(res.data?.errors || "Failed to delete task");
  }

  return res.data;
};

export const CreateTaskRepository = async (
  data: Omit<TaskUdate, "id">
): Promise<TApiResponse<TTask>> => {
  const res = await axiosClient.post("/tasks/", data);

  if (res.status !== 200 && res.status !== 201) {
    throw new Error(res.data?.errors || "Failed to create task");
  }

  return res.data;
};
