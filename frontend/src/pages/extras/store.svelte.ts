import { t } from "$shared/stores/locale.svelte";
import {
  notifyError,
  notifySuccess,
  notifyWarning,
} from "$shared/stores/notification.svelte";
import { nginxTemplate, type NginxAction } from "./service";
import * as api from "./api";

let activeVps = $state<any>(null);
let site = $state("");
let activeTab = $state<"available" | "enabled">("available");
let available = $state<string[]>([]);
let enabled = $state<string[]>([]);
let loading = $state(true);
let error = $state("");
let content = $state("");
let editorKey = $state(0);
let busy = $state<string | null>(null);
let loadingFile = $state(false);

export const extrasStore = {
  get activeVps() {
    return activeVps;
  },
  get site() {
    return site;
  },
  set site(value: string) {
    site = value;
  },
  get activeTab() {
    return activeTab;
  },
  set activeTab(value: "available" | "enabled") {
    activeTab = value;
  },
  get available() {
    return available;
  },
  get enabled() {
    return enabled;
  },
  get loading() {
    return loading;
  },
  get loadingFile() {
    return loadingFile;
  },
  get error() {
    return error;
  },
  get content() {
    return content;
  },
  set content(value: string) {
    content = value;
  },
  get editorKey() {
    return editorKey;
  },
  get busy() {
    return busy;
  },

  setVps(value: any) {
    activeVps = value;
  },

  async load(silent = false) {
    if (!activeVps) {
      loading = false;
      return;
    }
    if (!silent) loading = true;
    error = "";
    try {
      const result = await api.listNginxSites(activeVps);
      available = result?.available || [];
      enabled = result?.enabled || [];
      if (!site && available.length) await this.open(available[0]);
    } catch (cause: any) {
      error =
        cause?.message ||
        String(cause) ||
        "Não foi possível listar os arquivos do Nginx.";
      if (!silent) {
        available = [];
        enabled = [];
      }
    } finally {
      loading = false;
    }
  },

  async open(filename: string) {
    if (!activeVps || !filename) return;
    site = filename;
    content = ""; // Immediately clears previous content from editor memory
    editorKey += 1;
    loadingFile = true;
    try {
      const fetchedContent = (await api.readNginxSite(activeVps, filename, activeTab)) || "";
      content = fetchedContent;
      editorKey += 1;
    } catch (cause: any) {
      notifyError(
        cause?.message || String(cause) || "Não foi possível abrir o arquivo.",
      );
    } finally {
      loadingFile = false;
    }
  },

  newSite() {
    site = "";
    content = nginxTemplate;
    editorKey += 1;
  },

  async remove() {
    if (!site.trim()) return notifyWarning(t("extras.select_file_warn"));
    if (!activeVps) return notifyError(t("common.no_vps"));
    busy = "delete";
    try {
      const result = await api.deleteNginxSite(activeVps, site, activeTab);
      if (!result?.success)
        return notifyError(result?.message || t("extras.delete_error", { error: "Falha ao apagar" }));
      site = "";
      content = "";
      editorKey += 1;
      notifySuccess(result.message || t("extras.delete_success"));
      await this.load();
    } catch (cause: any) {
      notifyError(t("extras.delete_error", { error: cause?.message || String(cause) }));
    } finally {
      busy = null;
    }
  },

  async run(action: NginxAction) {
    if (!activeVps) return notifyError(t("common.no_vps"));
    if ((action === "enable" || action === "save") && !site.trim())
      return notifyWarning(t("extras.input_name_warn"));
    busy = action;
    try {
      const result =
        action === "save"
          ? await api.saveNginxSite(activeVps, site, content)
          : action === "enable"
            ? await api.enableNginxSite(activeVps, site)
            : action === "test"
              ? await api.testNginxConfig(activeVps)
              : await api.restartNginx(activeVps);
      if (!result?.success)
        return notifyError(
          result?.message ||
            result?.output ||
            t("extras.exec_error", { error: "Falha na ação" }),
        );
      notifySuccess(
        result.message || result.output || t("extras.exec_success"),
      );
      if (action === "save" || action === "enable") await this.load(true);
    } catch (cause: any) {
      notifyError(
        t("extras.exec_error", { error: cause?.message || String(cause) }),
      );
    } finally {
      busy = null;
    }
  },
};
