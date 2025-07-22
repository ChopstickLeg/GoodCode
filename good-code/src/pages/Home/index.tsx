import React from "react";
import { Helmet } from "react-helmet";
import { useDashboardData } from "../../hooks";
import {
  getRecentRoasts,
  filterValidRepositories,
} from "../../utils/dataHelpers";
import { LoadingSpinner, ErrorMessage } from "../../components/Common";
import { RepositoryList } from "../../components/Repository";
import { RoastList } from "../../components/Roast";
import GitHubAppInstall from "../../components/GitHubAppInstall";

const Home: React.FC = () => {
  const { data, error, isLoading, refetch } = useDashboardData();

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
        <p className="text-gray-500">No data available</p>
      </div>
    );
  }

  const { owned_repositories = [], collaborating_repositories = [] } = data;

  if (data.github_id == BigInt(0)) {
    return (
      <div className="max-w-md mx-auto mt-8">
        <GitHubAppInstall />
      </div>
    );
  }

  const validOwnedRepos = filterValidRepositories(owned_repositories);
  const validCollaboratingRepos = filterValidRepositories(
    collaborating_repositories
  );

  const recentRoasts = getRecentRoasts(
    validOwnedRepos,
    validCollaboratingRepos,
    5
  );

  return (
    <div className="min-h-screen">
      <Helmet>
        <title>Dashboard - Good Code</title>
        <meta
          name="description"
          content="Dashboard showing your repositories and recent AI roasts"
        />
      </Helmet>

      <div className="container mx-auto px-4 py-8">
        <header className="mb-12 text-center">
          <h1 className="text-5xl font-bold mb-4 bg-gradient-to-r from-blue-600 via-amber-600 to-blue-600 dark:from-blue-400 dark:via-amber-400 dark:to-blue-400 bg-clip-text text-transparent">
            Dashboard
          </h1>
          <p className="text-xl text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
            Welcome back! Here's what's happening with your repositories and
            recent AI insights.
          </p>
        </header>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
          <div className="bg-white dark:bg-gray-800 rounded-2xl p-8 text-center shadow-lg border border-gray-200 dark:border-gray-700 card-hover">
            <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-blue-600 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg
                className="w-8 h-8 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M3 8.2C3 7.07989 3 6.51984 3.21799 6.09202C3.40973 5.71569 3.71569 5.40973 4.09202 5.21799C4.51984 5 5.0799 5 6.2 5H9.67452C10.1637 5 10.4083 5 10.6385 5.05526C10.8425 5.10425 11.0376 5.18506 11.2166 5.29472C11.4184 5.4184 11.5914 5.59135 11.9373 5.93726L12.0627 6.06274C12.4086 6.40865 12.5816 6.5816 12.7834 6.70528C12.9624 6.81494 13.1575 6.89575 13.3615 6.94474C13.5917 7 13.8363 7 14.3255 7H17.8C18.9201 7 19.4802 7 19.908 7.21799C20.2843 7.40973 20.5903 7.71569 20.782 8.09202C21 8.51984 21 9.0799 21 10.2V15.8C21 16.9201 21 17.4802 20.782 17.908C20.5903 18.2843 20.2843 18.5903 19.908 18.782C19.4802 19 18.9201 19 17.8 19H6.2C5.07989 19 4.51984 19 4.09202 18.782C3.71569 18.5903 3.40973 18.2843 3.21799 17.908C3 17.4802 3 16.9201 3 15.8V8.2Z"
                />
              </svg>
            </div>
            <h3 className="text-3xl font-bold text-blue-600 dark:text-blue-400 mb-2">
              {validOwnedRepos.length}
            </h3>
            <p className="text-gray-600 dark:text-gray-400 font-medium">
              Owned Repositories
            </p>
          </div>

          <div className="bg-white dark:bg-gray-800 rounded-2xl p-8 text-center shadow-lg border border-gray-200 card-hover">
            <div className="w-16 h-16 bg-gradient-to-br from-green-500 to-green-600 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg
                className="w-8 h-8 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M18 19V18C18 15.7909 16.2091 14 14 14H10C7.79086 14 6 15.7909 6 18V19M23 19V18C23 15.7909 21.2091 14 19 14H18.5M1 19V18C1 15.7909 2.79086 14 5 14H5.5M17 11C18.6569 11 20 9.65685 20 8C20 6.34315 18.6569 5 17 5M7 11C5.34315 11 4 9.65685 4 8C4 6.34315 5.34315 5 7 5M15 8C15 9.65685 13.6569 11 12 11C10.3431 11 9 9.65685 9 8C9 6.34315 10.3431 5 12 5C13.6569 5 15 6.34315 15 8Z"
                />
              </svg>
            </div>
            <h3 className="text-3xl font-bold text-green-600 dark:text-green-400 mb-2">
              {validCollaboratingRepos.length}
            </h3>
            <p className="text-gray-600 dark:text-gray-400 font-medium">
              Collaborating Repositories
            </p>
          </div>

          <div className="bg-white dark:bg-gray-800 rounded-2xl p-8 text-center shadow-lg border border-gray-200 dark:border-gray-700 card-hover">
            <div className="w-16 h-16 bg-gradient-to-br from-amber-500 to-amber-600 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg
                className="w-8 h-8 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M8.5 14.5A2.5 2.5 0 0011 12c0-1.38-.5-2-1-3-1.072-2.143-.224-4.054 2-6 .5 2.5 2 4.9 4 6.5 2 1.6 3 3.5 3 5.5a7 7 0 11-14 0c0-1.153.433-2.294 1-3a2.5 2.5 0 002.5 2.5z"
                />
              </svg>
            </div>
            <h3 className="text-3xl font-bold text-amber-600 dark:text-amber-400 mb-2">
              {recentRoasts.length}
            </h3>
            <p className="text-gray-600 dark:text-gray-400 font-medium">
              Recent Roasts
            </p>
          </div>
        </div>

        <div className="grid grid-cols-1 xl:grid-cols-2 gap-8 mb-12">
          <div className="bg-white dark:bg-gray-800 rounded-2xl p-8 shadow-lg border border-gray-200 dark:border-gray-700">
            <RepositoryList
              repositories={validOwnedRepos}
              title="Your Repositories"
              limit={5}
              showViewAll={true}
              viewAllLink="/repositories"
            />
          </div>

          <div className="bg-white dark:bg-gray-800 rounded-2xl p-8 shadow-lg border border-gray-200 dark:border-gray-700">
            <RoastList roasts={recentRoasts} title="Recent Roasts" limit={5} />
          </div>
        </div>

        {validCollaboratingRepos.length > 0 && (
          <div className="bg-white dark:bg-gray-800 rounded-2xl p-8 shadow-lg border border-gray-200 dark:border-gray-700">
            <RepositoryList
              repositories={validCollaboratingRepos}
              title="Collaborating Repositories"
              limit={5}
              showViewAll={true}
              viewAllLink="/repositories"
            />
          </div>
        )}
      </div>
    </div>
  );
};

export default Home;
