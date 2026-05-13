import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
    "./lib/**/*.{ts,tsx}"
  ],
  theme: {
    extend: {
      colors: {
        ink: "#211A2D",
        mist: "#F8F4FA",
        lilac: "#D9C7FF",
        peach: "#FFD8C7",
        mint: "#C9F2E3",
        sky: "#C8E5FF"
      },
      boxShadow: {
        soft: "0 20px 60px rgba(82, 61, 112, 0.12)"
      }
    }
  },
  plugins: []
};

export default config;
