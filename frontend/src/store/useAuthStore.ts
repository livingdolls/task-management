import { create } from "zustand";
import type { TApiResponse } from "../types/api";
import axiosClient from "../lib/axios";

interface User {
  id: number;
  username: string;
  name: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  loading: boolean;
  setToken: (token: string) => void;
  fetchUser: () => Promise<void>;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: localStorage.getItem("token"),
  loading: true,

  setToken: (token) => {
    localStorage.setItem("token", token);
    set({ token });
  },

  fetchUser: async () => {
    try {
      const res = await axiosClient.get<TApiResponse<User>>("/profile");
      set({ user: res.data.data, loading: false });
    } catch {
      localStorage.removeItem("token");
      set({ user: null, token: null, loading: false });
    }
  },

  logout: () => {
    localStorage.removeItem("token");
    set({ user: null, token: null });
  },
}));
