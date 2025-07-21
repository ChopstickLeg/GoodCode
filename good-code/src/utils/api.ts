import {
  RepositoryDetails,
  AuthStatus,
  LoginRequest,
  SignupRequest,
  LoginResponse,
  SignupResponse,
  Repository,
  UserRepositoryCollaborator,
  AIRoast,
  UserLogin,
  GitHubAppInstallResponse,
  GitHubAppSetup,
} from "../types";

export class APIError extends Error {
  constructor(
    message: string,
    public status?: number,
    public response?: Response
  ) {
    super(message);
    this.name = "APIError";
  }
}

const apiFetch = async <T>(
  url: string,
  options: RequestInit = {}
): Promise<T> => {
  try {
    const response = await fetch(url, {
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
      ...options,
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new APIError(
        errorText || `HTTP error! status: ${response.status}`,
        response.status,
        response
      );
    }

    return (await response.json()) as T;
  } catch (error) {
    if (error instanceof APIError) {
      throw error;
    }
    throw new APIError(
      error instanceof Error ? error.message : "An unknown error occurred"
    );
  }
};

export const fetchRepositoryDetails = async (
  repoId: string
): Promise<RepositoryDetails> => {
  const [repoRes, collaboratorsRes, roastsRes] = await Promise.all([
    apiFetch<Repository>(`/api/repositories/${repoId}`),
    apiFetch<UserRepositoryCollaborator[]>(
      `/api/repositories/${repoId}/collaborators`
    ),
    apiFetch<AIRoast[]>(`/api/repositories/${repoId}/roasts`),
  ]);

  return {
    repo: repoRes,
    collaborators: collaboratorsRes,
    roasts: roastsRes,
  };
};

export const fetchAuthStatus = async (): Promise<AuthStatus> => {
  return apiFetch<AuthStatus>("/api/auth/verifyJWT");
};

export const loginUser = async (
  credentials: LoginRequest
): Promise<LoginResponse> => {
  return apiFetch<LoginResponse>("/api/account/login", {
    method: "POST",
    body: JSON.stringify(credentials),
  });
};

export const signupUser = async (
  userData: SignupRequest
): Promise<SignupResponse> => {
  return apiFetch<SignupResponse>("/api/account/signup", {
    method: "POST",
    body: JSON.stringify(userData),
  });
};

export const fetchRepositories = async (): Promise<UserLogin> => {
  const result = await apiFetch<UserLogin>("/api/repositories");
  console.table(result);
  return result;
};

export const postGitHubAppInstall = async (
  data: GitHubAppSetup
): Promise<GitHubAppInstallResponse> => {
  const result = await apiFetch<GitHubAppInstallResponse>(
    "/api/github/install",
    {
      method: "POST",
      body: JSON.stringify(data),
    }
  );
  return result;
};
