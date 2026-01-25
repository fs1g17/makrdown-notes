"use client";

import axios from "axios";

const clientFetch = axios.create();

clientFetch.interceptors.response.use(
  (response) => {
    if (response.status >= 400 && response.status <= 403) {
      window.location.pathname = "/sign-in";
      return Promise.reject("Unauthorized");
    }
    return response;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default clientFetch;