import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";

const useAuth = () => {
  const [authenticated, setAuthenticated] = useState<boolean | null>(null);

  useEffect(() => {
    fetch("/auth/check-session", { credentials: "include" })
      .then((res) => res.json())
      .then((data) => setAuthenticated(data.loggedIn))
      .catch(() => setAuthenticated(false));
  }, []);

  return authenticated;
};

export const PrivateRoute = ({ children }: { children: JSX.Element }) => {
  const auth = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    if (auth === false) {
      navigate("/login", { state: { from: location }, replace: true });
    }
  }, [auth]);

  if (auth === null) return null; // Show nothing while checking auth

  return children;
};

export default PrivateRoute;