import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuthStore } from "../store/useAuthStore";
import { useMutation } from "@tanstack/react-query";
import { LoginRepository } from "../repository/auth_repository";
import { LoaderCircle } from "lucide-react";

export const LoginPage = () => {
  const navigate = useNavigate();
  const setToken = useAuthStore((s) => s.setToken);
  const [form, setForm] = useState({ username: "", password: "" });
  const [error, setError] = useState("");
  const token = useAuthStore((state) => state.token);

  useEffect(() => {
    if (token) {
      navigate("/tasks");
    }
  }, [token]);

  const loginMutation = useMutation({
    mutationFn: LoginRepository,
    onSuccess: (res) => {
      if (!res.success || !res.data) return setError("Login failed");
      setToken(res.data.token);
    },
    onError: (err: any) => {
      setError(err.response?.data?.message || "Login failed");
    },
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    loginMutation.mutate(form);
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h1 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Sign in to your account
          </h1>
        </div>
        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="rounded-md shadow-sm -space-y-px">
            <div>
              <input
                type="text"
                placeholder="Username"
                value={form.username}
                onChange={(e) => setForm({ ...form, username: e.target.value })}
                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                required
              />
            </div>
            <div>
              <input
                type="password"
                placeholder="Password"
                value={form.password}
                onChange={(e) => setForm({ ...form, password: e.target.value })}
                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                required
              />
            </div>
          </div>

          {error && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded relative">
              {error}
            </div>
          )}

          <div>
            <button
              type="submit"
              disabled={loginMutation.isPending}
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loginMutation.isPending ? "Signing in..." : "Sign in"}
            </button>
          </div>
        </form>

        <div>
          <p className="text-sm text-center text-gray-600">
            Don't have an account?{" "}
            <button
              disabled={loginMutation.isPending}
              onClick={() => navigate("/register")}
              className="font-medium text-indigo-600 hover:text-indigo-500 cursor-pointer"
            >
              <LoaderCircle
                size={4}
                className={`${
                  loginMutation.isPending ? "animate-pulse block" : "hidden"
                }`}
              />
              Register
            </button>
          </p>
        </div>
      </div>
    </div>
  );
};
