import { useState } from "react";
import { useNavigate } from "react-router-dom";
import SignInButton from "../buttons/MSSignIn.tsx";
import GoogleSignInButton from "../buttons/GoogleSignIn.tsx";
import GitHubSignInButton from "../buttons/GithubSignIn.tsx";

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [rememberMe, setRememberMe] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    setIsLoading(true);
    console.log(JSON.stringify({ email, password }));
    try {
      const response = await fetch("/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });
      if (response.ok) {
        navigate("/");
      } else {
        setError("Invalid credentials");
      }
    } catch (error) {
      console.error(error);
      setError("An error occurred while logging in");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div>
      <div className="flex items-center justify-center h-full">
        <form onSubmit={handleSubmit} className="">
          <h2 className="text-2xl font-bold text-center mb-5">Login</h2>

          {error && <div className="text-red-500 text-sm">{error}</div>}
          <div className="flex flex-col space-y-4">
            <div>
              <input
                type="email"
                id="email"
                placeholder="Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full px-3 py-2 border rounded-md focus:outline-none focus:border-blue-500"
                required
              />
            </div>
            <div>
              <input
                type="password"
                id="password"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full px-3 py-2 border rounded-md focus:outline-none focus:border-blue-500"
                required
              />
            </div>
          </div>
          <div className="flex items-center justify-between px-1">
            <label className="text-gray-700 mt-1 mb-1">
              <input
                type="checkbox"
                checked={rememberMe}
                onChange={(e) => setRememberMe(e.target.checked)}
                className="mr-2"
              />
              Remember Me
            </label>
          </div>

          <button
            type="submit"
            disabled={isLoading}
            className="w-full bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 mb-3"
          >
            {isLoading ? (
              <div className="flex items-center justify-center">
                <svg
                  className="animate-spin h-5 w-5 mr-3"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                  ></circle>
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
                </svg>
                Loading...
              </div>
            ) : (
              "Login"
            )}
          </button>

          <div className="text-center text-gray-600">
            Don't have an account?{" "}
            <button
              type="button"
              onClick={() => navigate("/signup")}
              className="font-semibold text-blue-500 hover:text-blue-600 w-25 h-12 ml-2"
            >
              Sign up
            </button>
          </div>
          <div className="flex justify-center items-center text-gray-600 mt-3">
            <SignInButton />
          </div>
          <div className="flex justify-center items-center text-gray-600 mt-3">
            <GoogleSignInButton />
          </div>
          <div className="flex justify-center items-center text-gray-600 mt-3">
            <GitHubSignInButton />
          </div>
        </form>
      </div>
    </div>
  );
};

export default Login;
