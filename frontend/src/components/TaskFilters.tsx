import type { TaskFilter } from "../hooks/useTasks";
import { TASK_STATUS_OPTIONS } from "../utils/taskUtils";

interface TaskFiltersProps {
  filters: TaskFilter;
  onFiltersChange: (filters: TaskFilter) => void;
}

export const TaskFilters = ({ filters, onFiltersChange }: TaskFiltersProps) => {
  const handleStatusChange = (status: string) => {
    onFiltersChange({
      ...filters,
      status: status || undefined,
    });
  };

  const handleDeadlineChange = (deadline: string) => {
    onFiltersChange({
      ...filters,
      deadline: deadline || undefined,
    });
  };

  return (
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
          value={filters.status || ""}
          onChange={(e) => handleStatusChange(e.target.value)}
          className="border border-gray-300 rounded-md px-3 py-2 bg-white w-full focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
        >
          {TASK_STATUS_OPTIONS.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
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
          value={filters.deadline || ""}
          onChange={(e) => handleDeadlineChange(e.target.value)}
          className="border border-gray-300 rounded-md px-3 py-2 bg-white w-full focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
        />
      </div>
    </div>
  );
};
