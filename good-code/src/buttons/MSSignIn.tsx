import React from "react";
import MicrosoftLogin from "react-microsoft-login";

const MSSignInButton: React.FC = () => {
  const clientId = import.meta.env.VITE_MS_CLIENT_ID!;

  const authHandler = (err: any, data: any) => {
    if (err) {
      console.error("Authentication Error:", err);
    } else {
      console.log("Authentication Data:", data);
    }
  };

  return (
    <MicrosoftLogin
      clientId={clientId}
      authCallback={authHandler}
      buttonTheme="dark"
    >
      <button
        style={{
          display: "flex",
          alignItems: "center",
          color: "white",
          padding: "10px 20px",
          borderRadius: "4px",
          border: "none",
          cursor: "pointer",
        }}
      >
        <img
          src="https://learn.microsoft.com/en-us/entra/identity-platform/media/howto-add-branding-in-apps/ms-symbollockup_mssymbol_19.svg"
          alt="Microsoft Logo"
          style={{
            width: "20px",
            height: "20px",
            marginRight: "8px",
          }}
        />
        Sign in with Microsoft
      </button>
    </MicrosoftLogin>
  );
};

export default MSSignInButton;
