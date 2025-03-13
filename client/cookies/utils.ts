import Cookies from "js-cookie";

// Set a cookie
export const setCookie = (key:any, value:any, days = 7) => {
  Cookies.set(key, value, { expires: days, secure: true, sameSite: "Strict" });
};

// Get a cookie
export const getCookie = (key:any) => {
  return Cookies.get(key);
};

// Remove a cookie
export const removeCookie = (key:any) => {
  Cookies.remove(key);
};
