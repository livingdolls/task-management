import { create } from "zustand";

interface SearchStore {
  updateModal: boolean;
  setUpdateModal: (value: boolean) => void;
}

export const useUiStore = create<SearchStore>((set) => ({
  updateModal: false,
  setUpdateModal: (value: boolean) => set({ updateModal: value }),
}));

interface ModalCreateTask {
  createModal: boolean;
  setCreateModal: (value: boolean) => void;
}

export const useCreateTaskStore = create<ModalCreateTask>((set) => ({
  createModal: false,
  setCreateModal: (value: boolean) => set({ createModal: value }),
}));
