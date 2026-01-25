"use client";

import axios from "axios";

const clientFetch = axios.create({
  baseURL: process.env.NEXT_PUBLIC_BACKEND_BASE_URL,
});

clientFetch.interceptors.response.use(
  (response) => {
    if (response.status === 401) {
      window.location.pathname = "/app/sign-in";
      return Promise.reject("Unauthorized");
    }
    return response;
  },
);

export default clientFetch;