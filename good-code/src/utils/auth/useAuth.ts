import { useQuery } from "@tanstack/react-query";

const fetchAuthStatus = async () => {
  const res = await fetch("/api/v1/auth/verifyJWT", { credentials: "include" });
  if (!res.ok) throw new Error("Failed to fetch auth status");
  return res.json();
};

export const useAuth = () => {
  return useQuery({
    queryKey: ["authStatus"],
    queryFn: fetchAuthStatus,
    retry: false,
    staleTime: 5 * 60 * 1000,
  });
};
