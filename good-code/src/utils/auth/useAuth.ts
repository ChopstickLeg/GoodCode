import { useQuery } from "@tanstack/react-query";
import { fetchAuthStatus } from "../api";

export const useAuth = () => {
  return useQuery({
    queryKey: ["authStatus"],
    queryFn: fetchAuthStatus,
    retry: false,
    staleTime: 5 * 60 * 1000,
  });
};
