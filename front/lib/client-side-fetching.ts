"use client";

import axios from "axios";

const clientFetch = axios.create({});

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