import React from "react";
import { useParams } from "react-router-dom";
import { Helmet } from "react-helmet";
import { useRepositoryDetails } from "../../hooks";
import { LoadingSpinner, ErrorMessage } from "../../components/Common";
import { RoastList } from "../../components/Roast";
import { UserRepositoryCollaborator } from "../../types";

const RepositoryDetails: React.FC = () => {
  const { repoId } = useParams<{ repoId: string }>();
  const { data, error, isLoading, refetch } = useRepositoryDetails(repoId);

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
        <p className="text-gray-500">Repository not found</p>
      </div>
    );
  }

  const { repo, collaborators, roasts } = data;
  const recentRoasts = roasts.filter(
    (r) =>
      new Date(r.created_at) > new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)
  );

  return (
    <div className="min-h-screen">
      <Helmet>
        <title>{repo.name} - Good Code</title>
        <meta
          name="description"
          content={`Repository details for ${repo.name}`}
        />
        <link rel="icon" type="image/svg+xml" href="./GoodCode Logo.svg" />
      </Helmet>

      <div className="container mx-auto px-4 py-8">
        <header className="mb-12">
          <div className="bg-white dark:bg-gray-800 rounded-2xl p-8 shadow-lg border border-gray-200 dark:border-gray-700">
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <h1 className="text-4xl font-bold mb-3 bg-gradient-to-r from-blue-600 to-amber-600 dark:from-blue-400 dark:to-amber-400 bg-clip-text text-transparent">
                  {repo.name}
                </h1>
                <p className="text-gray-600 dark:text-gray-400 mb-4 text-lg">
                  Owned by{" "}
                  <span className="font-semibold text-blue-600 dark:text-blue-400">
                    {repo.owner}
                  </span>
                </p>
                <div className="flex items-center space-x-6 text-sm text-gray-500 dark:text-gray-400">
                  <div className="flex items-center space-x-1">
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
                        d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                      />
                    </svg>
                    <span>
                      Created: {new Date(repo.created_at).toLocaleDateString()}
                    </span>
                  </div>
                  <div className="flex items-center space-x-1">
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
                        d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                    <span>
                      Updated: {new Date(repo.updated_at).toLocaleDateString()}
                    </span>
                  </div>
                </div>
              </div>
              <div className="flex items-center space-x-2">
                <div className="w-4 h-4 bg-green-500 rounded-full"></div>
                <span className="text-sm font-medium text-green-600">
                  Active
                </span>
              </div>
            </div>
          </div>
        </header>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
          <div className="bg-white dark:bg-gray-800 rounded-2xl p-8 text-center shadow-lg border border-gray-200 card-hover">
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
                  d="M18 19V18C18 15.7909 16.2091 14 14 14H10C7.79086 14 6 15.7909 6 18V19M23 19V18C23 15.7909 21.2091 14 19 14H18.5M1 19V18C1 15.7909 2.79086 14 5 14H5.5M17 11C18.6569 11 20 9.65685 20 8C20 6.34315 18.6569 5 17 5M7 11C5.34315 11 4 9.65685 4 8C4 6.34315 5.34315 5 7 5M15 8C15 9.65685 13.6569 11 12 11C10.3431 11 9 9.65685 9 8C9 6.34315 10.3431 5 12 5C13.6569 5 15 6.34315 15 8Z"
                />
              </svg>
            </div>
            <h3 className="text-3xl font-bold text-blue-600 mb-2">
              {collaborators.length}
            </h3>
            <p className="text-gray-600 dark:text-gray-400 font-medium">
              Collaborators
            </p>
          </div>

          <div className="bg-white dark:bg-gray-800 rounded-2xl p-8 text-center shadow-lg border border-gray-200 card-hover">
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
            <h3 className="text-3xl font-bold text-amber-600 mb-2">
              {roasts.length}
            </h3>
            <p className="text-gray-600 dark:text-gray-400 font-medium">
              Total Roasts
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
                  d="M12 7V12H15M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12Z"
                />
              </svg>
            </div>
            <h3 className="text-3xl font-bold text-green-600 mb-2">
              {recentRoasts.length}
            </h3>
            <p className="text-gray-600 dark:text-gray-400 font-medium">
              Recent Roasts
            </p>
          </div>
        </div>

        <div className="grid grid-cols-1 xl:grid-cols-2 gap-8">
          <div className="bg-white dark:bg-gray-800 rounded-2xl p-8 shadow-lg border border-gray-200">
            <h2 className="text-2xl font-bold mb-6 text-gray-900 dark:text-gray-100">
              Collaborators
            </h2>
            {collaborators.length === 0 ? (
              <div className="text-center py-12">
                <div className="w-24 h-24 mx-auto mb-4 bg-gray-100 dark:bg-gray-700 rounded-full flex items-center justify-center">
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
                      d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"
                    />
                  </svg>
                </div>
                <p className="text-gray-600 text-lg">No collaborators found.</p>
              </div>
            ) : (
              <div className="space-y-4">
                {collaborators.map(
                  (collaborator: UserRepositoryCollaborator) => (
                    <div
                      key={collaborator.id}
                      className="flex items-center p-4 bg-gray-50 dark:bg-gray-800 rounded-xl border border-gray-200 hover:shadow-md transition-all duration-200"
                    >
                      <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-blue-600 rounded-full flex items-center justify-center mr-4">
                        {collaborator.github_avatar_url ? (
                          <img
                            src={collaborator.github_avatar_url}
                            alt={collaborator.github_login}
                            className="w-10 h-10 rounded-full"
                          />
                        ) : (
                          <span className="text-white text-xl font-bold">
                            {collaborator.github_login?.charAt(0).toUpperCase()}
                          </span>
                        )}
                      </div>
                      <div className="flex-1">
                        <p className="font-semibold text-gray-900 dark:text-gray-100 text-lg">
                          {collaborator.github_login}
                        </p>
                        <p className="text-sm text-gray-600 dark:text-gray-400 capitalize">
                          {collaborator.role}
                        </p>
                      </div>
                      <div className="flex items-center space-x-2">
                        <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800 capitalize">
                          {collaborator.role}
                        </span>
                      </div>
                    </div>
                  )
                )}
              </div>
            )}
          </div>

          <div className="bg-white dark:bg-gray-800 rounded-2xl p-8 shadow-lg border border-gray-200">
            <RoastList roasts={roasts} title="AI Roasts" />
          </div>
        </div>
      </div>
    </div>
  );
};

export default RepositoryDetails;
