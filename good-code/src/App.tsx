import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import "./App.css";
import Signup from "./pages/Signup";
import Login from "./pages/Login";
import Home from "./pages/Home";
import Repositories from "./pages/Repositories";
import RepositoryDetails from "./pages/RepositoryDetails";
import PrivateRoute from "./utils/auth/VerifyJWT";
import Layout from "./components/Layout";
import GitHubInstall from "./pages/GitHubInstall";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Signup />} />
        <Route
          path="/"
          element={
            <PrivateRoute>
              <Layout>
                <Home />
              </Layout>
            </PrivateRoute>
          }
        />
        <Route
          path="/repositories"
          element={
            <PrivateRoute>
              <Layout>
                <Repositories />
              </Layout>
            </PrivateRoute>
          }
        />
        <Route
          path="/repositories/:repoId"
          element={
            <PrivateRoute>
              <Layout>
                <RepositoryDetails />
              </Layout>
            </PrivateRoute>
          }
        />
        <Route
          path="/github/install"
          element={
            <Layout>
              <GitHubInstall>
                <Home />
              </GitHubInstall>
            </Layout>
          }
        />
      </Routes>
    </Router>
  );
}

export default App;
