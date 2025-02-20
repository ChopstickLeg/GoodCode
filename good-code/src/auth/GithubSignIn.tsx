const GITHUB_CLIENT_ID = import.meta.env.VITE_GITHUB_CLIENT_ID;
const REDIRECT_URI = "http://localhost:3000/auth/github/callback";
import GitHubLogo from "../github-mark-white.svg";

const GitHubSignInButton: React.FC = () => {
  const handleLogin = () => {
    console.log("GitHub Login button clicked");
    window.location.href = `https://github.com/login/oauth/authorize?client_id=${GITHUB_CLIENT_ID}&redirect_uri=${REDIRECT_URI}`;
  };

  return (
    <button
      onClick={handleLogin}
      style={{
        display: "flex",
        alignItems: "center",
        backgroundColor: "#24292e",
        color: "#ffffff",
        padding: "10px 20px",
        border: "none",
        borderRadius: "5px",
        cursor: "pointer",
        fontSize: "16px",
        fontWeight: "bold",
      }}
    >
      <img
        src={GitHubLogo}
        alt="GitHub Logo"
        style={{ width: "24px", height: "24px", marginRight: "10px" }}
      />
      Sign in with GitHub
    </button>
  );
};

export default GitHubSignInButton;
