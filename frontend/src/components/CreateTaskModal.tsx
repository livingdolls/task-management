import React, { useState } from "react";
import { useCreateTaskStore } from "../store/useUiStore";

interface CreateTaskModalProps {
  onSubmit: (task: {
    title: string;
    description: string;
    status: "To Do" | "In Progress" | "Done";
    deadline?: Date | null;
  }) => void;
}

export const CreateTaskModal: React.FC<CreateTaskModalProps> = ({
  onSubmit,
}) => {
  const { createModal, setCreateModal } = useCreateTaskStore();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [status, setStatus] = useState<"To Do" | "In Progress" | "Done">(
    "To Do"
  );
  const [deadline, setDeadline] = useState<Date | null>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (title.trim()) {
      onSubmit({ title, description, status, deadline });
      setTitle("");
      setDescription("");
      setStatus("To Do");
      setDeadline(null);
    }
  };

  console.log("Create Modal State:", createModal);

  if (!createModal) return null;

  return (
    <div
      className="fixed inset-0 bg-black/30  flex items-center justify-center z-50"
      onClick={() => setCreateModal(false)}
    >
      <div
        className="bg-white rounded-lg shadow-lg w-full max-w-md mx-4"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between p-6">
          <h2 className="text-xl font-semibold text-gray-800">
            Create New Task
          </h2>
          <button
            className="text-gray-400 hover:text-gray-600 text-2xl font-bold"
            onClick={() => setCreateModal(false)}
          >
            ×
          </button>
        </div>
        <form onSubmit={handleSubmit} className="p-6">
          <div className="mb-4">
            <label
              htmlFor="title"
              className="block text-sm font-medium text-gray-700 mb-2"
            >
              Title *
            </label>
            <input
              id="title"
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Enter task title"
              required
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
          <div className="mb-4">
            <label
              htmlFor="description"
              className="block text-sm font-medium text-gray-700 mb-2"
            >
              Description
            </label>
            <textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Enter task description"
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            />
          </div>
          <div className="mb-6">
            <label
              htmlFor="status"
              className="block text-sm font-medium text-gray-700 mb-2"
            >
              Priority
            </label>
            <select
              id="status"
              value={status}
              onChange={(e) =>
                setStatus(e.target.value as "To Do" | "In Progress" | "Done")
              }
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="To Do">To Do</option>
              <option value="In Progress">In Progress</option>
              <option value="Done">Done</option>
            </select>
          </div>

          <div className="mb-6">
            <label className="block text-sm font-medium mb-1">Due Date</label>
            <input
              type="date"
              value={deadline ? deadline.toISOString().split("T")[0] : ""}
              onChange={(e) =>
                setDeadline(e.target.value ? new Date(e.target.value) : null)
              }
              className="w-full border rounded px-3 py-2"
            />
          </div>

          <div className="flex gap-3 justify-end">
            <button
              type="button"
              onClick={() => setCreateModal(false)}
              className="px-4 py-2 text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-gray-300"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              Create Task
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
