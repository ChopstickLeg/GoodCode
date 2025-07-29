import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import {
  fetchRepositoryDetails,
  fetchAuthStatus,
  loginUser,
  signupUser,
  fetchRepositories,
  postGitHubAppInstall,
} from "../utils/api";
import {
  UserLogin,
  RepositoryDetails,
  AuthStatus,
  LoginRequest,
  SignupRequest,
  LoginResponse,
  SignupResponse,
  GitHubAppSetup,
  GitHubAppInstallResponse,
} from "../types";

export const useDashboardData = () => {
  return useQuery<UserLogin, Error>({
    queryKey: ["dashboardData"],
    queryFn: fetchRepositories,
    staleTime: 5 * 60 * 1000,
    retry: 2,
  });
};

export const useRepositoryDetails = (repoId: string | undefined) => {
  return useQuery<RepositoryDetails, Error>({
    queryKey: ["repositoryDetails", repoId],
    queryFn: () => {
      if (!repoId) throw new Error("Repository ID is required");
      return fetchRepositoryDetails(repoId);
    },
    enabled: !!repoId,
    staleTime: 5 * 60 * 1000,
    retry: 2,
  });
};

export const useRepositories = () => {
  return useQuery<UserLogin, Error>({
    queryKey: ["repositories"],
    queryFn: fetchRepositories,
    staleTime: 5 * 60 * 1000,
    retry: 2,
  });
};

export const useAuth = () => {
  return useQuery<AuthStatus, Error>({
    queryKey: ["authStatus"],
    queryFn: fetchAuthStatus,
    retry: false,
    staleTime: 5 * 60 * 1000,
  });
};

export const useLogin = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  return useMutation<LoginResponse, Error, LoginRequest>({
    mutationFn: loginUser,
    onSuccess: (response) => {
      if (response.success) {
        queryClient.invalidateQueries({ queryKey: ["authStatus"] });
        navigate("/");
      }
    },
  });
};

export const useSignup = () => {
  const navigate = useNavigate();

  return useMutation<SignupResponse, Error, SignupRequest>({
    mutationFn: signupUser,
    onSuccess: (response) => {
      if (response.success) {
        navigate("/login");
      }
    },
  });
};

export const useLogout = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  return useMutation<void, Error, void>({
    mutationFn: async () => {
      await fetch("/api/auth/logout", {
        method: "POST",
        credentials: "include",
      });
    },
    onSuccess: () => {
      queryClient.clear();
      navigate("/login");
    },
  });
};

export const useGitHubAppInstall = (githubAppSetup: GitHubAppSetup) => {
  return useMutation<GitHubAppInstallResponse, Error, GitHubAppSetup>({
    mutationFn: () => {
      if (!githubAppSetup.installation_id) {
        throw new Error(
          "Installation ID is required for GitHub App installation"
        );
      }
      return postGitHubAppInstall(githubAppSetup);
    },
    onSuccess: (response) => {
      if (response.success) {
      }
    },
    onError: (error) => {
      console.error("Error installing GitHub App:", error);
    },
  });
};
