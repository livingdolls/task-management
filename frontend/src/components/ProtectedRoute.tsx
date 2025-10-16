import { Navigate } from "react-router-dom";
import type { JSX } from "react";
import { useAuthGuard } from "../store/useAuthGuard";

export default function ProtectedRoute({
  children,
}: {
  children: JSX.Element;
}) {
  const { token, user, loading } = useAuthGuard();

  if (loading) return <p>Loading...</p>;
  if (!token || !user) return <Navigate to="/login" replace />;

  return children;
}
