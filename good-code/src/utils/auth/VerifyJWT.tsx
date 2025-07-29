import React, { useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "../../hooks";
import { LoadingSpinner } from "../../components/Common";

interface PrivateRouteProps {
  children: React.ReactNode;
}

export const PrivateRoute: React.FC<PrivateRouteProps> = ({ children }) => {
  const { data, isLoading, isError } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    if (isError || (data && !data.loggedIn)) {
      navigate("/login", { state: { from: location }, replace: true });
    }
  }, [data, isError, navigate, location]);

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <LoadingSpinner size="large" />
      </div>
    );
  }

  if (isError || !data?.loggedIn) {
    return null;
  }

  return <>{children}</>;
};

export default PrivateRoute;
