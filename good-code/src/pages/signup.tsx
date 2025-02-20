import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";

const SignUp = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");

  const navigate = useNavigate();

  const signupMutationFn = async (data: {
    email: string;
    password: string;
  }): Promise<{ success: boolean }> => {
    try {
      const response = await fetch("/api/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error("Invalid credentials");
      }

      return response.json();
    } catch (error) {
      throw error;
    }
  };

  const { mutate: signupUser, isPending } = useMutation<
    { success: boolean },
    Error,
    { email: string; password: string }
  >({
    mutationFn: signupMutationFn,
    onSuccess: (result) => {
      if (result.success) {
        navigate("/");
      }
    },
    onError: (error) => {
      console.error("Signup failed:", error);
      setError(
        error instanceof Error ? error.message : "An unknown error occurred"
      );
    },
  });

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError("");

    signupUser({ email, password });
  };

  return (
    <form onSubmit={handleSubmit}>
      <label>
        Email:
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
      </label>
      <label>
        Password:
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
      </label>
      <label>
        Confirm Password:
        <input
          type="password"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
        />
      </label>
      <button type="submit" disabled={isPending}>
        Sign Up
      </button>
      {error && <p>{error}</p>}
    </form>
  );
};

export default SignUp;