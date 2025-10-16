import { Navigate } from "react-router-dom";
import { useAuthStore } from "../store/useAuthStore";

export const HomeRedirect = () => {
  const token = useAuthStore((s) => s.token);
  return <Navigate to={token ? "/tasks" : "/login"} replace />;
};
