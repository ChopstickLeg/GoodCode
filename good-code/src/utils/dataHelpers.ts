import { Repository, AIRoast } from "../types";

export const extractRoastsFromRepositories = (
  repositories: Repository[]
): AIRoast[] => {
  return repositories.reduce<AIRoast[]>((acc, repo) => {
    if (repo.ai_roasts && Array.isArray(repo.ai_roasts)) {
      return acc.concat(repo.ai_roasts);
    }
    return acc;
  }, []);
};

export const sortRoastsByDate = (roasts: AIRoast[]): AIRoast[] => {
  return roasts.sort((a, b) => {
    const dateA = new Date(a.created_at).getTime();
    const dateB = new Date(b.created_at).getTime();
    return dateB - dateA;
  });
};

export const getRecentRoasts = (
  ownedRepositories: Repository[],
  collaboratingRepositories: Repository[],
  limit: number = 5
): AIRoast[] => {
  const ownedRoasts = extractRoastsFromRepositories(ownedRepositories);
  const collaboratingRoasts = extractRoastsFromRepositories(
    collaboratingRepositories
  );

  const allRoasts = [...ownedRoasts, ...collaboratingRoasts];
  const sortedRoasts = sortRoastsByDate(allRoasts);

  return sortedRoasts.slice(0, limit);
};

export const formatDate = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
};

export const formatRelativeTime = (dateString: string): string => {
  const date = new Date(dateString);
  const now = new Date();
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);

  if (diffInSeconds < 60) {
    return "just now";
  } else if (diffInSeconds < 3600) {
    const minutes = Math.floor(diffInSeconds / 60);
    return `${minutes} minute${minutes > 1 ? "s" : ""} ago`;
  } else if (diffInSeconds < 86400) {
    const hours = Math.floor(diffInSeconds / 3600);
    return `${hours} hour${hours > 1 ? "s" : ""} ago`;
  } else {
    const days = Math.floor(diffInSeconds / 86400);
    return `${days} day${days > 1 ? "s" : ""} ago`;
  }
};

export const truncateText = (text: string, maxLength: number): string => {
  if (text.length <= maxLength) return text;
  return text.substring(0, maxLength).trim() + "...";
};

export const isValidRepository = (repo: unknown): repo is Repository => {
  return (
    typeof repo === "object" &&
    repo !== null &&
    typeof (repo as Repository).id === "number" &&
    typeof (repo as Repository).name === "string" &&
    typeof (repo as Repository).owner === "string" &&
    (Array.isArray((repo as Repository).ai_roasts) ||
      (repo as Repository).ai_roasts === null)
  );
};

export const filterValidRepositories = (
  repositories: Repository[]
): Repository[] => {
  return !repositories ? [] : repositories.filter(isValidRepository);
};
