import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import electron from "vite-plugin-electron";

export default defineConfig({
  plugins: [
    react(),
    electron([
      {
        // Main-Process entry file of the Electron App.
        entry: "electron/main.ts",
      },
      {
        // Preload-Scripts entry file (optional)
        entry: "electron/preload.ts",
        onstart(options) {
          // Notify the Renderer-Process to reload the page when the Preload-Scripts build is complete,
          // instead of restarting the entire Electron App.
          options.reload();
        },
      },
    ]),
  ],
  build: {
    rollupOptions: {
      external: ["electron"],
    },
  },
});
