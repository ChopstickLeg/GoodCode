export interface UserLogin {
  id: bigint;
  email: string;
  name: string;
  github_id: bigint;
  enabled: boolean;
  created_at: string;
  updated_at: string;
  owned_repositories: Repository[];
  collaborating_repositories: Repository[];
}

export interface Repository {
  id: bigint;
  name: string;
  owner: string;
  owner_id: bigint;
  installation_id: bigint;
  created_at: string;
  updated_at: string;
  ai_roasts: AIRoast[];
}

export interface AIRoast {
  id: bigint;
  repo_id: bigint;
  pull_request_number: number;
  content: string;
  is_open: boolean;
  created_at: string;
  updated_at: string;
  repository: Repository;
}

export interface UserRepositoryCollaborator {
  id: bigint;
  repository_id: bigint;
  github_user_id: bigint;
  github_login: string;
  role: string;
  is_good_code_user: boolean;
  user_login_id?: bigint;
  github_avatar_url?: string;
}

export interface RepositoryDetails {
  repo: Repository;
  collaborators: UserRepositoryCollaborator[];
  roasts: AIRoast[];
}

export interface AuthStatus {
  loggedIn: boolean;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface SignupRequest {
  email: string;
  name: string;
  password: string;
}

export interface LoginResponse {
  success: boolean;
  message?: string;
}

export interface SignupResponse {
  success: boolean;
  message?: string;
}

export interface GitHubAppSetup {
  installation_id: string;
  setup_action: string;
}

export interface GitHubAppInstallResponse {
  success: boolean;
  message?: string;
}
