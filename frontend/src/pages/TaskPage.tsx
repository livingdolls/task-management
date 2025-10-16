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
import { type TTask, type TTaskStatus } from "../types/task";
import { CreateTaskModal } from "../components/CreateTaskModal";
import { CirclePlus } from "lucide-react";

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
    status: TTaskStatus;
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
      <div className="flex flex-row gap-4 mb-4 items-center">
        <h1 className="text-2xl font-bold ">{user?.name} Tasks</h1>

        <button
          onClick={handleOpenCreateModal}
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded flex flex-row gap-4 items-center"
        >
          Add Task
          <CirclePlus />
        </button>
      </div>
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

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {(data?.data || []).map((task: any, index: number) => {
          const colors = [
            "bg-red-100 border-red-200",
            "bg-blue-100 border-blue-200",
            "bg-green-100 border-green-200",
            "bg-yellow-100 border-yellow-200",
            "bg-purple-100 border-purple-200",
            "bg-pink-100 border-pink-200",
            "bg-indigo-100 border-indigo-200",
            "bg-orange-100 border-orange-200",
          ];
          const cardColor = colors[index % colors.length];

          return (
            <div
              key={task.id}
              className={`${cardColor} border rounded-lg p-6 shadow-md hover:shadow-lg transition-shadow`}
            >
              <div className="mb-4">
                <h3 className="text-lg font-semibold text-gray-900 mb-2">
                  {task.title}
                </h3>
                <p className="text-gray-600 text-sm mb-3">{task.description}</p>

                <div className="flex items-center justify-between mb-4">
                  <span
                    className={`px-3 py-1 text-xs font-semibold rounded-full ${
                      task.status === "Done"
                        ? "bg-green-500 text-white"
                        : task.status === "In Progress"
                        ? "bg-blue-500 text-white"
                        : "bg-gray-500 text-white"
                    }`}
                  >
                    {task.status}
                  </span>

                  <span className="text-sm text-gray-700">
                    {task.deadline
                      ? isNaN(Date.parse(task.deadline))
                        ? ""
                        : new Date(task.deadline).toLocaleDateString()
                      : "No deadline"}
                  </span>
                </div>
              </div>

              <div className="flex justify-end space-x-2">
                <button
                  onClick={() => handleSetTask(task)}
                  className="bg-indigo-500 hover:bg-indigo-600 text-white px-4 py-2 rounded text-sm font-medium transition-colors"
                >
                  Edit
                </button>
                <button
                  onClick={() => handleDeleteTask(task.id)}
                  className="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded text-sm font-medium transition-colors"
                >
                  Delete
                </button>
              </div>
            </div>
          );
        })}
      </div>

      {/* Modal */}
      <UpdateTaskModal task={task} onUpdate={handleUpdate} />

      {/* Modal Create */}
      <CreateTaskModal onSubmit={handleCreateTask} />
    </div>
  );
};
