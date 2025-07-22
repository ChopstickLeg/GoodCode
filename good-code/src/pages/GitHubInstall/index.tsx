import React, { useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { LoadingSpinner } from "../../components/Common";
import { GitHubAppSetup } from "../../types";
import { useGitHubAppInstall } from "../../hooks";

const GitHubInstall: React.FC = () => {
  const params = new URLSearchParams(window.location.search);
  const installationId = params.get("installation_id");
  const setupAction = params.get("setup_action");

  const installData = {
    installation_id: installationId ? BigInt(installationId) : 0n,
    setup_action: setupAction || "",
  } as GitHubAppSetup;

  console.log("GitHub Install Data:", installData);

  const navigate = useNavigate();
  const location = useLocation();

  const { mutate, data, isError, isPending } = useGitHubAppInstall(installData);

  useEffect(() => {
    if (installData.installation_id) {
      mutate(installData);
    }
  }, []);

  useEffect(() => {
    if (isError || (data && !data.success)) {
      navigate("/login", { state: { from: location }, replace: true });
    } else if (data && data.success) {
      navigate("/", { replace: true });
    }
  }, [data, isError, navigate, location]);

  if (isPending) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <LoadingSpinner size="large" />
      </div>
    );
  }

  return null;
};

export default GitHubInstall;
