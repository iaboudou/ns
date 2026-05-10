/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://backend:4001/api/:path*",
      },
      {
        source: "/pics/:path*",
        destination: "http://backend:4001/pics/:path*",
      },
      {
        source: "/hassession",
        destination: "http://backend:4001/hassession",
      },
    ];
  },
};

export default nextConfig;
