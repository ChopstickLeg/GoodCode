const AuthenticateUser = () => {
  const getAuthCookie = () => {
    const cookieName = "auth";
    const decodedCookie = decodeURIComponent(document.cookie);
    const cookieArray = decodedCookie.split(";");
    console.log(cookieArray);
    for (let i = 0; i < cookieArray.length; i++) {
      let cookie = cookieArray[i].trimStart();
      if (cookie.indexOf(cookieName) === 0) {
        return cookie.substring(cookieName.length, cookie.length);
      }
    }
    return null;
  };
  const cookie = getAuthCookie();
  return cookie == null ? true : false;
};
export default AuthenticateUser;
