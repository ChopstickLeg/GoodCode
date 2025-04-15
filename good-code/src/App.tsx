import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import "./App.css";
import Signup from "./pages/signup";
import Login from "./pages/login";
import Home from "./pages/home";
import PrivateRoute from "./auth/VerifyJWT";
import PrPage from "./pages/pr";

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
              <Home />
            </PrivateRoute>
          }
        />
        <Route
          path="/pr/:prID"
          element={
            <PrivateRoute>
              <PrPage />
            </PrivateRoute>
          }
        />
      </Routes>
    </Router>
  );
}

export default App;
