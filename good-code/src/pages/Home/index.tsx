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

  const { owned_repositories, collaborating_repositories } = data;

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
                  d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
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
                  d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"
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
                  d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
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
