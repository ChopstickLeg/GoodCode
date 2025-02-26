import { useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "./useAuth";

export const PrivateRoute = ({ children }: { children: JSX.Element }) => {
  const { data, isLoading, isError } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    if (isError || (data && !data.loggedIn)) {
      navigate("/login", { state: { from: location }, replace: true });
    }
  }, [data, isError]);

  if (isLoading) return (<div>Loading...</div>); // Show loading state while checking auth

  return children;
};

export default PrivateRoute;
