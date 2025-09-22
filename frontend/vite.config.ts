import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react-swc";
import electron from "vite-plugin-electron";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");

  return {
    plugins: [
      react(),
      electron([
        {
          entry: "electron/main.ts",
          vite: {
            define: {
              "process.env.NODE_ENV": JSON.stringify(env.NODE_ENV),
            },
          },
        },
        {
          entry: "electron/preload.ts",
          onstart(options) {
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
  };
});
