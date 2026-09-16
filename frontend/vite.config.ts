import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import wails from "@wailsio/runtime/plugins/vite";
import path from "path";
import fs from "fs";

// Read app version from single source of truth: root VERSION file
function getAppVersion(): string {
  try {
    const versionPath = path.resolve(__dirname, "../VERSION");
    if (fs.existsSync(versionPath)) {
      return fs.readFileSync(versionPath, "utf-8").trim();
    }
    const configPath = path.resolve(__dirname, "../build/config.yml");
    const content = fs.readFileSync(configPath, "utf-8");
    const match = content.match(/^version:\s*["']?([^"'\r\n]+)["']?/m);
    return match ? match[1] : "0.0.0";
  } catch {
    return "0.0.0";
  }
}

const appVersion = getAppVersion();

// https://vitejs.dev/config/
export default defineConfig({
  define: {
    __APP_VERSION__: JSON.stringify(appVersion),
  },
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
  resolve: {
    alias: {
      $lib: path.resolve(__dirname, "./src/lib"),
      $shared: path.resolve(__dirname, "./src/lib/shared"),
      $navigation: path.resolve(__dirname, "./src/lib/shared/navigation"),
      $session: path.resolve(__dirname, "./src/lib/shared/session"),
      $bindings: path.resolve(__dirname, "./bindings/go-walis/internal"),
    },
  },
  plugins: [tailwindcss(), svelte(), wails("./bindings")],
});
