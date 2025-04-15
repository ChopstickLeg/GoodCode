import { useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "./useAuth";

export const PrivateRoute = ({ children }: { children: JSX.Element }) => {
  const { data, isLoading, isError } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  console.log("Auth Data:", data);
  console.log("isLoading:", isLoading);
  console.log("isError:", isError);

  useEffect(() => {
    if (isError || (data && !data.loggedIn)) {
      navigate("/login", { state: { from: location }, replace: true });
    }
  }, [data, isError, navigate, location]);

  if (isLoading) return <div>Loading...</div>;

  return children;
};

export default PrivateRoute;
