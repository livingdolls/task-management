import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  CreateTaskRepository,
  DeleteTaskRepository,
  TaskRepository,
  UpdateTaskRepository,
} from "../repository/task_repository";
import { useAuthStore } from "../store/useAuthStore";
import { useState } from "react";
import { useCreateTaskStore, useUiStore } from "../store/useUiStore";
import { UpdateTaskModal, type TaskUdate } from "../components/UpdateTaskModal";
import { type TTask } from "../types/task";
import { CreateTaskModal } from "../components/CreateTaskModal";

export type TaskFilter = {
  status?: string;
  deadline?: string;
};

export const TaskPage = () => {
  const queryClient = useQueryClient();
  const [statusFilter, setStatusFilter] = useState<TaskFilter>({});
  const { setUpdateModal } = useUiStore();
  const { setCreateModal } = useCreateTaskStore();
  const [task, setTask] = useState<
    Pick<TTask, "title" | "id" | "description" | "status" | "deadline">
  >({
    id: 0,
    title: "",
    description: "",
    status: "To Do",
    deadline: new Date(),
  });

  const { data, isLoading, error } = useQuery({
    queryKey: ["tasks", statusFilter],
    queryFn: () => TaskRepository(statusFilter),
  });

  const { user } = useAuthStore();

  const updateTask = useMutation({
    mutationFn: (taskData: TaskUdate) =>
      UpdateTaskRepository(taskData.id, taskData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
      alert("Task updated successfully");
    },
    onError: (error) => {
      alert((error as Error).message || "Failed to update task");
    },
  });

  const createTask = useMutation({
    mutationFn: CreateTaskRepository,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
      alert("Task created successfully");
    },
    onError: (error) => {
      alert((error as Error).message || "Failed to create task");
    },
  });

  const handleUpdate = (updatedTask: TaskUdate) => {
    updateTask.mutate(updatedTask);
  };

  const handleSetTask = (task: TTask) => {
    const data: TaskUdate = {
      id: task.id,
      title: task.title,
      description: task.description,
      status: task.status,
      deadline: task.deadline ? new Date(task.deadline) : undefined,
    };

    setTask(data);
    setUpdateModal(true);
  };

  const deleteTask = useMutation({
    mutationFn: DeleteTaskRepository,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
      alert("Task deleted successfully");
    },
    onError: (error) => {
      alert((error as Error).message || "Failed to delete task");
    },
  });

  const handleCreateTask = ({
    title,
    description,
    status,
    deadline,
  }: {
    title: string;
    description: string;
    status: "To Do" | "In Progress" | "Done";
    deadline?: Date | null;
  }) => {
    createTask.mutate({ title, description, status, deadline });
    setCreateModal(false);
  };

  const handleDeleteTask = (id: number) => {
    deleteTask.mutate(id);
  };

  const handleOpenCreateModal = () => {
    setCreateModal(true);
  };

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {(error as Error).message}</div>;

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">{user?.name} Tasks</h1>

      <button
        onClick={handleOpenCreateModal}
        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
      >
        Add Task
      </button>
      <div className="mb-4 grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label
            htmlFor="status-filter"
            className="block text-sm font-medium text-gray-700 mb-2"
          >
            Filter by Status:
          </label>
          <select
            id="status-filter"
            value={statusFilter.status || ""}
            onChange={(e) =>
              setStatusFilter({
                ...statusFilter,
                status: e.target.value || undefined,
              })
            }
            className="border border-gray-300 rounded-md px-3 py-2 bg-white w-full"
          >
            <option value="">All</option>
            <option value="To Do">To Do</option>
            <option value="In Progress">In Progress</option>
            <option value="Done">Done</option>
          </select>
        </div>

        <div>
          <label
            htmlFor="deadline-filter"
            className="block text-sm font-medium text-gray-700 mb-2"
          >
            Deadline Before:
          </label>
          <input
            type="date"
            id="deadline-filter"
            value={statusFilter.deadline || ""}
            onChange={(e) =>
              setStatusFilter({
                ...statusFilter,
                deadline: e.target.value || undefined,
              })
            }
            className="border border-gray-300 rounded-md px-3 py-2 bg-white w-full"
          />
        </div>
      </div>

      <div className="overflow-x-auto">
        <table className="min-w-full bg-white border border-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Title
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Deadline
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {(data?.data || []).map((task: any) => (
              <tr key={task.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {task.title}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span
                    className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                      task.status === "Done"
                        ? "bg-green-100 text-green-800"
                        : task.status === "In Progress"
                        ? "bg-blue-100 text-blue-800"
                        : "bg-gray-100 text-gray-800"
                    }`}
                  >
                    {task.status}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {task.deadline
                    ? isNaN(Date.parse(task.deadline))
                      ? ""
                      : new Date(task.deadline).toLocaleDateString()
                    : ""}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <button
                    onClick={() => handleSetTask(task)}
                    className="text-indigo-600 hover:text-indigo-900 mr-2"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => handleDeleteTask(task.id)}
                    className="text-red-600 hover:text-red-900"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Modal */}
      <UpdateTaskModal task={task} onUpdate={handleUpdate} />

      {/* Modal Create */}
      <CreateTaskModal onSubmit={handleCreateTask} />
    </div>
  );
};
