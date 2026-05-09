/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://localhost:4001/api/:path*",
      },
      {
        source: "/pics/:path*",
        destination: "http://localhost:4001/pics/:path*",
      },
      {
        source: "/hassession",
        destination: "http://localhost:4001/hassession",
      },
    ];
  },
};

export default nextConfig;
