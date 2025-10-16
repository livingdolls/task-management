export type TTask = {
  id: number;
  title: string;
  description: string;
  status: "To Do" | "In Progress" | "Done";
  deadline?: Date | null;
  created_at: Date;
};
