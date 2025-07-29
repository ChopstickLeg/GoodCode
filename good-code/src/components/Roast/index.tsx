import React from "react";
import { AIRoast } from "../../types";
import { formatRelativeTime, truncateText } from "../../utils/dataHelpers";
import { useRepositoryDetails } from "../../hooks";
import { LoadingSpinner, ErrorMessage } from "../../components/Common";

interface RoastCardProps {
  roast: AIRoast;
  className?: string;
}

export const RoastCard: React.FC<RoastCardProps> = ({
  roast,
  className = "",
}) => {
  const { data, error, isLoading, refetch } = useRepositoryDetails(
    roast.repo_id.toString()
  );

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <LoadingSpinner size="large" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-md mx-auto mt-8">
        <ErrorMessage message={error.message} onRetry={() => refetch()} />
      </div>
    );
  }

  if (!data) {
    return (
      <div className="text-center py-8">
        <p className="text-gray-500 dark:text-gray-400">No data available</p>
      </div>
    );
  }

  const repository = data.repo;

  const linkToPR = `https://github.com/${repository.owner}/${repository.name}/pull/${roast.pull_request_number}`;

  return (
    <a href={linkToPR} target="_blank" rel="noopener noreferrer">
      <div
        className={`p-6 bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 hover:shadow-lg transition-all duration-300 ${className}`}
      >
        <div className="flex items-start justify-between mb-3">
          <div className="flex-1">
            <div className="flex items-center space-x-2 mb-3">
              <div className="w-3 h-3 bg-amber-500 rounded-full"></div>
              <p className="text-sm text-gray-600 dark:text-gray-400 font-medium">
                Pull Request #{roast.pull_request_number} in{" "}
                <span className="font-bold text-blue-600 dark:text-blue-400">
                  {repository.name}
                </span>
              </p>
            </div>
            <p className="text-gray-900 dark:text-gray-100 text-sm leading-relaxed bg-gray-50 dark:bg-gray-700 p-3 rounded-lg border border-gray-100 dark:border-gray-600">
              {truncateText(roast.content, 150)}
            </p>
          </div>
        </div>
        <div className="flex items-center justify-between mt-4">
          <span className="text-xs text-gray-500 dark:text-gray-400 flex items-center space-x-1">
            <svg
              className="w-3 h-3"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <span>{formatRelativeTime(roast.created_at)}</span>
          </span>
          <div className="flex items-center space-x-2">
            <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-900/50 text-blue-800 dark:text-blue-300">
              PR #{roast.pull_request_number}
            </span>
          </div>
        </div>
      </div>
    </a>
  );
};

interface RoastListProps {
  roasts: AIRoast[];
  title: string;
  limit?: number;
  className?: string;
}

export const RoastList: React.FC<RoastListProps> = ({
  roasts,
  title,
  limit,
  className = "",
}) => {
  const displayedRoasts = limit ? roasts.slice(0, limit) : roasts;

  if (roasts == null || roasts.length === 0) {
    return (
      <div className={className}>
        <h2 className="text-2xl font-bold mb-6 text-gray-900 dark:text-gray-100">
          {title}
        </h2>
        <div className="text-center py-12">
          <div className="w-24 h-24 mx-auto mb-4 bg-gray-100 dark:bg-gray-700 rounded-full flex items-center justify-center">
            <svg
              className="w-12 h-12 text-gray-400 dark:text-gray-500"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
              />
            </svg>
          </div>
          <p className="text-gray-600 dark:text-gray-400 text-lg">
            No roasts found.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className={className}>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          {title}
        </h2>
        {limit && roasts.length > limit && (
          <span className="text-blue-600 text-sm font-medium">
            Showing {limit} of {roasts.length} roasts
          </span>
        )}
      </div>
      <div className="space-y-4">
        {displayedRoasts.map((roast) => (
          <RoastCard key={roast.id} roast={roast} />
        ))}
      </div>
    </div>
  );
};
