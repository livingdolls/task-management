export type TTask = {
  id: number;
  title: string;
  description: string;
  status: TTaskStatus;
  deadline?: Date | null;
  created_at: Date;
};

export type TTaskStatus = "To Do" | "In Progress" | "Done";
