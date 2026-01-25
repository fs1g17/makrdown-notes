import type { NextConfig } from "next";

const API_URL = process.env.BACKEND_BASE_URL;

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: `${API_URL}/:path*`,
      },
    ];
  },
};

export default nextConfig;
