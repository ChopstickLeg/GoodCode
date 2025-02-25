import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Helmet } from "react-helmet";
// import SignInButton from "../auth/MSSignIn.tsx";
// import GoogleSignInButton from "../auth/GoogleSignIn.tsx";
// import GitHubSignInButton from "../auth/GithubSignIn.tsx";
import { useMutation } from "@tanstack/react-query";

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const navigate = useNavigate();

  const loginMutationFn = async (data: { email: string; password: string }) => {
    try {
      const response = await fetch("/api/account/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error(await response.text());
      }

      const result = await response.json();
      return { success: result.success };
    } catch (error) {
      throw error;
    }
  };

  const { mutate: loginUser, isPending } = useMutation<
    { success: boolean },
    Error,
    { email: string; password: string }
  >({
    mutationFn: loginMutationFn,
    onSuccess: (result) => {
      if (result.success) {
        console.log("Login successful");
        navigate("/home");
      }
    },

    onError: (error) => {
      console.error("Login failed:", error);
      setError(
        error instanceof Error ? error.message : "An unknown error occurred"
      );
    },
  });

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    console.log("Submit clicked");
    setError("");

    loginUser({ email, password });
  };

  return (
    <div>
      <Helmet>
        <title>Login</title>
      </Helmet>
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
          <button
            type="submit"
            disabled={isPending}
            className="w-full bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 mb-3"
          >
            {isPending ? (
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
          {/* <div className="flex justify-center items-center text-gray-600 mt-3">
            <SignInButton />
          </div>
          <div className="flex justify-center items-center text-gray-600 mt-3">
            <GoogleSignInButton />
          </div>
          <div className="flex justify-center items-center text-gray-600 mt-3">
            <GitHubSignInButton />
          </div>
          TODO: Add login with Microsoft, Google, and GitHub if possible*/}
        </form>
      </div>
    </div>
  );
};

export default Login;
