import { Call, Browser } from "@wailsio/runtime";

// Chamadas por nome evitam depender de bindings regenerados.
const SERVICE = "go-walis/internal/updater.UpdaterService";

export interface UpdateInfo {
  hasUpdate: boolean;
  currentVersion: string;
  latestVersion: string;
  releaseNotes: string;
  releaseUrl: string;
  publishedAt: string;
  assetName: string;
  assetSize: number;
  supported: boolean;
}

export interface InstallResult {
  success: boolean;
  cancelled: boolean;
  fallbackPath: string;
  message: string;
}

export const updaterApi = {
  getVersion: () => Call.ByName(`${SERVICE}.GetVersion`) as Promise<string>,
  check: () => Call.ByName(`${SERVICE}.CheckForUpdates`) as Promise<UpdateInfo>,
  install: () => Call.ByName(`${SERVICE}.InstallUpdate`) as Promise<InstallResult>,
  restart: () => Call.ByName(`${SERVICE}.Restart`) as Promise<void>,
  openUrl: (url: string) => Browser.OpenURL(url),
};
