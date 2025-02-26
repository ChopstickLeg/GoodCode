const AuthenticateUser = () => {
  const getAuthCookie = () => {
    const cookieName = "auth";
    const decodedCookie = decodeURIComponent(document.cookie);
    const cookieArray = decodedCookie.split(";");

    for (let i = 0; i < cookieArray.length; i++) {
      let cookie = cookieArray[i].trimStart();
      if (cookie.indexOf(cookieName) === 0) {
        return cookie.substring(cookieName.length, cookie.length);
      }
    }
    return null;
  };
  const cookie = getAuthCookie();
  const isLoggedIn = cookie == null ? true : false;
  return isLoggedIn;
};
export default AuthenticateUser;
