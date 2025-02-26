import React from "react";

import { useMsal } from "@azure/msal-react";

const SignInButton: React.FC = () => {
  const { instance } = useMsal();

  const handleLogin = () => {
    instance
      .loginPopup({
        scopes: ["User.Read"],
      })
      .catch((error) => {
        console.error(error);
      });
  };

  return <button onClick={handleLogin}>Sign In with Microsoft</button>;
};

export default SignInButton;
