import type { TTaskStatus } from "../types/task";

export const TASK_STATUS_OPTIONS: { value: TTaskStatus | ""; label: string }[] =
  [
    { value: "", label: "All" },
    { value: "To Do", label: "To Do" },
    { value: "In Progress", label: "In Progress" },
    { value: "Done", label: "Done" },
  ];

export const TASK_CARD_COLORS = [
  "bg-red-100 border-red-200",
  "bg-blue-100 border-blue-200",
  "bg-green-100 border-green-200",
  "bg-yellow-100 border-yellow-200",
  "bg-purple-100 border-purple-200",
  "bg-pink-100 border-pink-200",
  "bg-indigo-100 border-indigo-200",
  "bg-orange-100 border-orange-200",
];

export const getStatusStyles = (status: TTaskStatus): string => {
  switch (status) {
    case "Done":
      return "bg-green-500 text-white";
    case "In Progress":
      return "bg-blue-500 text-white";
    case "To Do":
    default:
      return "bg-gray-500 text-white";
  }
};

export const formatDate = (date: Date | string | null | undefined): string => {
  if (!date) return "No deadline";

  const parsedDate = typeof date === "string" ? new Date(date) : date;
  return isNaN(parsedDate.getTime())
    ? "Invalid date"
    : parsedDate.toLocaleDateString();
};
