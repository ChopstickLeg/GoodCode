import React from "react";
import { Link } from "react-router-dom";
import { Repository } from "../../types";

interface RepositoryCardProps {
  repository: Repository;
  className?: string;
}

export const RepositoryCard: React.FC<RepositoryCardProps> = ({
  repository,
  className = "",
}) => {
  return (
    <Link
      to={`/repositories/${repository.id}`}
      className={`block p-6 bg-white rounded-xl border border-gray-200 hover:border-blue-300 transition-all duration-300 card-hover group ${className}`}
    >
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <h3 className="text-xl font-bold text-gray-900 group-hover:text-blue-600 transition-colors">
            {repository.name}
          </h3>
          <p className="text-gray-600 mt-1 font-medium">{repository.owner}</p>
        </div>
        <div className="flex items-center space-x-2">
          <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
            {repository.ai_roasts.length} roast
            {repository.ai_roasts.length !== 1 ? "s" : ""}
          </span>
        </div>
      </div>
      <div className="mt-4 flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-1">
            <div className="w-3 h-3 bg-green-500 rounded-full"></div>
            <span className="text-sm text-gray-600">Active</span>
          </div>
        </div>
        <span className="text-sm text-gray-500">
          Updated {new Date(repository.updated_at).toLocaleDateString()}
        </span>
      </div>
    </Link>
  );
};

interface RepositoryListProps {
  repositories: Repository[];
  title: string;
  limit?: number;
  showViewAll?: boolean;
  viewAllLink?: string;
  className?: string;
}

export const RepositoryList: React.FC<RepositoryListProps> = ({
  repositories,
  title,
  limit,
  showViewAll = false,
  viewAllLink = "/repositories",
  className = "",
}) => {
  const displayedRepos = limit ? repositories.slice(0, limit) : repositories;

  if (repositories.length === 0) {
    return (
      <div className={className}>
        <h2 className="text-2xl font-bold mb-6 text-gray-900">{title}</h2>
        <div className="text-center py-12">
          <div className="w-24 h-24 mx-auto mb-4 bg-gray-100 rounded-full flex items-center justify-center">
            <svg
              className="w-12 h-12 text-gray-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
              />
            </svg>
          </div>
          <p className="text-gray-600 text-lg">No repositories found.</p>
        </div>
      </div>
    );
  }

  return (
    <div className={className}>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-900">{title}</h2>
        {showViewAll && repositories.length > (limit || 0) && (
          <Link
            to={viewAllLink}
            className="text-blue-600 hover:text-blue-700 font-medium transition-colors duration-200 flex items-center space-x-1"
          >
            <span>View all ({repositories.length})</span>
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 5l7 7-7 7"
              />
            </svg>
          </Link>
        )}
      </div>
      <div className="space-y-4">
        {displayedRepos.map((repo) => (
          <RepositoryCard key={repo.id} repository={repo} />
        ))}
      </div>
    </div>
  );
};
