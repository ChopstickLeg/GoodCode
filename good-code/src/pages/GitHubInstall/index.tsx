import React, { useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { LoadingSpinner } from "../../components/Common";
import { GitHubAppSetup } from "../../types";
import { useGitHubAppInstall } from "../../hooks";

interface GitHubInstallProps {
  children?: React.ReactNode;
}

const GitHubInstall: React.FC<GitHubInstallProps> = ({ children }) => {
  const params = new URLSearchParams(window.location.search);
  const installationId = params.get("installation_id");
  const setupAction = params.get("setup_action");

  const installData = {
    installation_id: installationId || "",
    setup_action: setupAction || "",
  } as GitHubAppSetup;

  const { data, isError, isPending } = useGitHubAppInstall(installData);
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    if (isError || (data && !data.success)) {
      navigate("/login", { state: { from: location }, replace: true });
    }
  }, [data, isError, navigate, location]);

  if (isPending) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <LoadingSpinner size="large" />
      </div>
    );
  }

  if (isError || !data?.success) {
    return null;
  }

  return <>{children}</>;
};

export default GitHubInstall;
