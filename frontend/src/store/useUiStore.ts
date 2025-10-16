import { create } from "zustand";

interface SearchStore {
  updateModal: boolean;
  setUpdateModal: (value: boolean) => void;
}

export const useUiStore = create<SearchStore>((set) => ({
  updateModal: false,
  setUpdateModal: (value: boolean) => set({ updateModal: value }),
}));
