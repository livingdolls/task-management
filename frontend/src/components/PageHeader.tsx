import { CirclePlus } from "lucide-react";

interface PageHeaderProps {
  userName?: string;
  onCreateTask: () => void;
  onLogout: () => void;
}

export const PageHeader = ({
  userName,
  onCreateTask,
  onLogout,
}: PageHeaderProps) => {
  return (
    <div className="flex justify-between">
      <div className="flex flex-row gap-4 mb-6 items-center">
        <h1 className="text-2xl font-bold text-gray-900">
          {userName ? `${userName}'s Tasks` : "Tasks"}
        </h1>

        <button
          onClick={onCreateTask}
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-5 px-5 md:py-2 md:px-4  flex flex-row gap-2 items-center transition-colors focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 fixed bottom-10 right-4 md:static rounded-full md:rounded"
        >
          <span className="hidden md:block">Add Task</span>
          <CirclePlus size={18} />
        </button>
      </div>

      <div>
        <button
          onClick={onLogout}
          className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded flex flex-row gap-2 items-center transition-colors focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
        >
          Logout
        </button>
      </div>
    </div>
  );
};
