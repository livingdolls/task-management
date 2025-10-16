import { CirclePlus } from "lucide-react";

interface PageHeaderProps {
  userName?: string;
  onCreateTask: () => void;
}

export const PageHeader = ({ userName, onCreateTask }: PageHeaderProps) => {
  return (
    <div className="flex flex-row gap-4 mb-6 items-center">
      <h1 className="text-2xl font-bold text-gray-900">
        {userName ? `${userName}'s Tasks` : "Tasks"}
      </h1>

      <button
        onClick={onCreateTask}
        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded flex flex-row gap-2 items-center transition-colors focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
      >
        Add Task
        <CirclePlus size={18} />
      </button>
    </div>
  );
};
