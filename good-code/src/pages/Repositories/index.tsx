import React from "react";
import { Helmet } from "react-helmet";
import { useRepositories } from "../../hooks";
import { filterValidRepositories } from "../../utils/dataHelpers";
import { LoadingSpinner, ErrorMessage } from "../../components/Common";
import { RepositoryList } from "../../components/Repository";

const Repositories: React.FC = () => {
  const { data, error, isLoading, refetch } = useRepositories();

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

  return (
    <div className="min-h-screen">
      <Helmet>
        <title>Repositories - Good Code</title>
        <meta
          name="description"
          content="Manage your repositories and collaborations"
        />
      </Helmet>

      <div className="container mx-auto px-4 py-8">
        <header className="mb-12 text-center">
          <h1 className="text-5xl font-bold mb-4 bg-gradient-to-r from-blue-600 to-amber-600 dark:from-blue-400 dark:to-amber-400 bg-clip-text text-transparent">
            Repositories
          </h1>
          <p className="text-xl text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
            Manage your repositories and collaborations all in one place
          </p>
        </header>

        <div className="grid grid-cols-1 xl:grid-cols-2 gap-8 mb-12">
          <div className="bg-white dark:bg-gray-800 rounded-2xl p-8 shadow-lg border border-gray-200 dark:border-gray-700">
            <RepositoryList
              repositories={validOwnedRepos}
              title="Owned Repositories"
            />
          </div>

          <div className="bg-white dark:bg-gray-800 rounded-2xl p-8 shadow-lg border border-gray-200 dark:border-gray-700">
            <RepositoryList
              repositories={validCollaboratingRepos}
              title="Collaborating Repositories"
            />
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-2xl p-8 shadow-lg border border-gray-200 dark:border-gray-700">
          <h2 className="text-3xl font-bold mb-8 text-center text-gray-900 dark:text-gray-100">
            Repository Summary
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div className="text-center p-6 bg-gradient-to-br from-blue-50 to-blue-100 dark:from-blue-900/30 dark:to-blue-800/30 rounded-xl">
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
              <h3 className="text-4xl font-bold text-blue-600 dark:text-blue-400 mb-2">
                {validOwnedRepos.length}
              </h3>
              <p className="text-gray-600 dark:text-gray-400 font-medium text-lg">
                Owned Repositories
              </p>
            </div>
            <div className="text-center p-6 bg-gradient-to-br from-green-50 to-green-100 dark:from-green-900/30 dark:to-green-800/30 rounded-xl">
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
              <h3 className="text-4xl font-bold text-green-600 dark:text-green-400 mb-2">
                {validCollaboratingRepos.length}
              </h3>
              <p className="text-gray-600 dark:text-gray-400 font-medium text-lg">
                Collaborating Repositories
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Repositories;
