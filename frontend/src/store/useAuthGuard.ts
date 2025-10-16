import { useEffect } from "react";
import { useAuthStore } from "./useAuthStore";

export function useAuthGuard() {
  const { token, user, loading, fetchUser } = useAuthStore();

  useEffect(() => {
    if (token && !user) {
      fetchUser();
    } else if (!token) {
      if (loading) {
        setTimeout(() => {
          useAuthStore.setState({ loading: false });
        }, 100);
      }
    }
  }, [token, user]);

  return { token, user, loading };
}
