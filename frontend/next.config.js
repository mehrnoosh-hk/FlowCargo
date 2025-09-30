/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Enable standalone output for better deployment
  output: "standalone",
  // Configure image domains if needed
  images: {
    remotePatterns: [],
  },
  // Environment variables available to the browser
  env: {
    NEXT_PUBLIC_API_URL:
      process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080",
  },
};

module.exports = nextConfig;
