import { BASE_URL } from "./src/config.mjs";

/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: `${BASE_URL}/api/:path*`,
      },
      {
        source: "/pics/:path*",
        destination: `${BASE_URL}/pics/:path*`,
      },
      {
        source: "/hassession",
        destination: `${BASE_URL}/hassession`,
      },
    ];
  },
  devIndicators: false,
};

export default nextConfig;
