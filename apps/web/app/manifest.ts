import type { MetadataRoute } from "next";

export default function manifest(): MetadataRoute.Manifest {
  return {
    name: "Moodora",
    short_name: "Moodora",
    description: "Cosmic lifestyle guidance for daily reflection.",
    start_url: "/",
    scope: "/",
    display: "standalone",
    orientation: "portrait",
    background_color: "#F8F4FA",
    theme_color: "#A78BFA",
    categories: ["lifestyle", "wellness"],
    shortcuts: [
      {
        name: "Daily Guide",
        short_name: "Daily",
        url: "/daily",
        icons: [{ src: "/icons/icon-192.svg", sizes: "192x192", type: "image/svg+xml" }]
      },
      {
        name: "Check-in",
        short_name: "Check-in",
        url: "/wallet",
        icons: [{ src: "/icons/icon-192.svg", sizes: "192x192", type: "image/svg+xml" }]
      },
      {
        name: "Tarot",
        short_name: "Tarot",
        url: "/tarot",
        icons: [{ src: "/icons/icon-192.svg", sizes: "192x192", type: "image/svg+xml" }]
      }
    ],
    icons: [
      { src: "/icons/icon-192.svg", sizes: "192x192", type: "image/svg+xml" },
      { src: "/icons/icon-512.svg", sizes: "512x512", type: "image/svg+xml" }
    ]
  };
}
