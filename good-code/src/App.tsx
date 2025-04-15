import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import "./App.css";
import Signup from "./pages/Signup";
import Login from "./pages/Login";
import Home from "./pages/Home";
import PrivateRoute from "./utils/auth/VerifyJWT";
import PrPage from "./pages/PR Overview";

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
            // <PrivateRoute>
            <PrPage />
            // </PrivateRoute>
          }
        />
      </Routes>
    </Router>
  );
}

export default App;
