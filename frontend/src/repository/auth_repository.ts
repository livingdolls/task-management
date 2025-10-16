import axiosClient from "../lib/axios";
import type { TApiResponse } from "../types/api";
import type { TLogin, TLoginResponse, TRegisterResponse } from "../types/auth";

export const LoginRepository = async (
  arg: TLogin
): Promise<TApiResponse<TLoginResponse>> => {
  const rest = await axiosClient.post("/auth/login", arg);

  if (rest.status !== 200 && rest.status !== 201) {
    throw new Error(rest.data?.errors || "Login failed");
  }

  return rest.data;
};

export const RegisterRepository = async (
  arg: TLogin
): Promise<TApiResponse<TRegisterResponse>> => {
  const rest = await axiosClient.post("/auth/register", arg);

  if (rest.status !== 200 && rest.status !== 201) {
    throw new Error(rest.data?.errors || "Registration failed");
  }

  return rest.data;
};
