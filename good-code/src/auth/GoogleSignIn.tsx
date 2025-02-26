import React from "react";
import GoogleButton from "react-google-button";

const GoogleSignInButton: React.FC = () => {
  const handleLogin = () => {
    console.log("Google button clicked");
  };

  return (
    <GoogleButton
      onClick={handleLogin}
      label="Sign in with Google"
      disabled={false}
      type="dark"
    />
  );
};

export default GoogleSignInButton;
