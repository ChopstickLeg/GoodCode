import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
  useLocation,
} from "react-router-dom";
import "./App.css";
import Signup from "./pages/signup";
import Login from "./pages/login";
import Home from "./pages/home";
import AuthenticateUser from "./auth/VerifyJWT";

const PrivateRoute = ({ children }: { children: JSX.Element }) => {
  const location = useLocation();
  if (!AuthenticateUser()) {
    console.log(true)
    return <Navigate to="/login" state={{ from: location }} replace />;
  }
  console.log(false)
  return children;
};

function App() {
  return (
    <Router>
      <Routes>
        <Route
          path="/"
          element={
            <PrivateRoute>
              <Home />
            </PrivateRoute>
          }
        />
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Signup />} />
      </Routes>
    </Router>
  );
}

export default App;
