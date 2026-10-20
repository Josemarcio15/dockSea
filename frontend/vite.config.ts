import { defineConfig, type Plugin } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import wails from "@wailsio/runtime/plugins/vite";
import path from "path";
import fs from "fs";

// Read app version from single source of truth: root VERSION file
function getAppVersion(): string {
  try {
    const versionPath = path.resolve(import.meta.dirname, "../VERSION");
    if (fs.existsSync(versionPath)) {
      return fs.readFileSync(versionPath, "utf-8").trim();
    }
    const configPath = path.resolve(import.meta.dirname, "../build/config.yml");
    const content = fs.readFileSync(configPath, "utf-8");
    const match = content.match(/^version:\s*["']?([^"'\r\n]+)["']?/m);
    return match ? match[1] : "0.0.0";
  } catch {
    return "0.0.0";
  }
}

function versionWatcherPlugin(): Plugin {
  const versionFile = path.resolve(import.meta.dirname, "../VERSION");
  return {
    name: "watch-version-file",
    buildStart() {
      this.addWatchFile(versionFile);
    },
    handleHotUpdate({ file, server }) {
      if (path.resolve(file) === versionFile) {
        server.restart();
      }
    },
  };
}

// https://vitejs.dev/config/
export default defineConfig(() => {
  const currentVersion = getAppVersion();

  return {
    define: {
      __APP_VERSION__: JSON.stringify(currentVersion),
    },
    server: {
      host: "127.0.0.1",
      port: Number(process.env.WAILS_VITE_PORT) || 9245,
      strictPort: true,
      watch: {
        // Assegura que arquivos fora do root do frontend (como ../VERSION) sejam monitorados
        ignored: ["!**/VERSION"],
      },
    },
    resolve: {
      alias: {
        $lib: path.resolve(import.meta.dirname, "./src/lib"),
        $shared: path.resolve(import.meta.dirname, "./src/lib/shared"),
        $navigation: path.resolve(import.meta.dirname, "./src/lib/shared/navigation"),
        $session: path.resolve(import.meta.dirname, "./src/lib/shared/session"),
        $bindings: path.resolve(import.meta.dirname, "./bindings/go-walis/internal"),
      },
    },
    plugins: [versionWatcherPlugin(), tailwindcss(), svelte(), wails("./bindings")],
  };
});
