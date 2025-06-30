import { reactRouter } from "@react-router/dev/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig({
  plugins: [tailwindcss(), reactRouter(), tsconfigPaths()],
  server: {
    host: true,
    port: 5173,
    // proxy: {
    //   '/api': {
    //     target: 'http://server:4000',
    //     changeOrigin: true,
    //     secure: false,
    //   }
    // }
    proxy: {
      '/api': {
        target: process.env.VITE_API_URL ||
          (process.env.NODE_ENV === 'development'
            ? `http://localhost:${process.env.SERVER_PORT || 4000}`  // Local development
            : `http://server:${process.env.SERVER_PORT || 4000}`),   // Docker
        changeOrigin: true,
        secure: false,
      }
    }
  }
});
