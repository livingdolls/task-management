import React, { useState, useEffect } from "react";
import { useUiStore } from "../store/useUiStore";

export type TaskUdate = {
  id: number;
  title: string;
  description: string;
  status: "To Do" | "In Progress" | "Done";
  deadline?: Date | null | undefined;
};

interface UpdateTaskModalProps {
  task?: TaskUdate;
  onUpdate: (task: TaskUdate) => void;
}

export const UpdateTaskModal: React.FC<UpdateTaskModalProps> = ({
  task,
  onUpdate,
}) => {
  const { updateModal, setUpdateModal } = useUiStore();
  const [formData, setFormData] = useState<Partial<TaskUdate>>({
    id: task?.id,
    title: "",
    description: "",
    status: "To Do",
    deadline: new Date(),
  });

  useEffect(() => {
    if (task) {
      setFormData(task);
    }
  }, [task]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (formData.title && task) {
      onUpdate({ ...task, ...formData } as TaskUdate);
      setUpdateModal(false);
    }
  };

  const handleClose = () => {
    setUpdateModal(false);
  };

  if (!updateModal) return null;

  return (
    <div className="fixed inset-0 bg-black/30 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-full max-w-md">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold">Update Task</h2>
          <button
            onClick={handleClose}
            className="text-gray-400 hover:text-gray-600"
          >
            ✕
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">Title</label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) =>
                setFormData({ ...formData, title: e.target.value })
              }
              className="w-full border rounded px-3 py-2"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">
              Description
            </label>
            <textarea
              value={formData.description}
              onChange={(e) =>
                setFormData({ ...formData, description: e.target.value })
              }
              className="w-full border rounded px-3 py-2 h-24"
              rows={3}
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Status</label>
            <select
              value={formData.status}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  status: e.target.value as TaskUdate["status"],
                })
              }
              className="w-full border rounded px-3 py-2"
            >
              <option value="To Do">To Do</option>
              <option value="In Progress">In Progress</option>
              <option value="Done">Done</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Due Date</label>
            <input
              type="date"
              value={
                formData.deadline instanceof Date
                  ? formData.deadline.toISOString().split("T")[0]
                  : formData.deadline || ""
              }
              onChange={(e) =>
                setFormData({ ...formData, deadline: new Date(e.target.value) })
              }
              className="w-full border rounded px-3 py-2"
            />
          </div>

          <div className="flex gap-2 pt-4">
            <button
              type="button"
              onClick={handleClose}
              className="flex-1 px-4 py-2 border rounded hover:bg-gray-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="flex-1 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Update Task
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
