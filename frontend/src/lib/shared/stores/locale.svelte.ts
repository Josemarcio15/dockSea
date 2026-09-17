import * as LocaleWails from "$bindings/locale/localeservice";

// Reactive translations map holding all language dictionaries loaded from ~/Documents/DockSea/locales/
const translations = $state<Record<string, any>>({});
let availableLocales = $state<string[]>([]);
let localesDirectoryPath = $state<string>("");

// Reactive locale state — Svelte 5 $state rune via .svelte.ts extension
let storedLocale = "pt-BR";
if (typeof localStorage !== "undefined") {
  storedLocale = localStorage.getItem("locale") || "pt-BR";
}
let currentLocale = $state(storedLocale || "pt-BR");

/**
 * Load all locale files from ~/Documents/DockSea/locales/ via Wails
 */
export async function initLocalesFromDisk(): Promise<void> {
  try {
    if (typeof LocaleWails.LoadAllLocales === "function") {
      const diskLocales = await LocaleWails.LoadAllLocales();
      if (diskLocales && typeof diskLocales === "object") {
        for (const [localeKey, jsonStr] of Object.entries(diskLocales)) {
          if (jsonStr) {
            try {
              const parsed = JSON.parse(jsonStr);
              translations[localeKey] = parsed;
            } catch (err) {
              console.warn(`[Locale] Falha ao analisar JSON do locale '${localeKey}':`, err);
            }
          }
        }
      }
    }

    if (typeof LocaleWails.ListLocales === "function") {
      const list = await LocaleWails.ListLocales();
      if (list && list.length > 0) {
        availableLocales = Array.from(new Set([...list, "pt-BR", "en-US"]));
      }
    }

    if (typeof LocaleWails.GetLocalesDir === "function") {
      localesDirectoryPath = (await LocaleWails.GetLocalesDir()) || "";
    }
  } catch (e) {
    console.warn("[Locale] Erro ao carregar locales do disco via Wails:", e);
  }
}

// Auto-run on startup if running in browser/webview environment
if (typeof window !== "undefined") {
  void initLocalesFromDisk();
}

export function getAvailableLocales(): string[] {
  return availableLocales;
}

export function getLocalesDirectory(): string {
  return localesDirectoryPath;
}

export function setLocale(locale: string) {
  const chosen = translations[locale] ? locale : "en-US";
  currentLocale = chosen;
  if (typeof localStorage !== "undefined") {
    localStorage.setItem("locale", chosen);
  }
}

export function getLocale(): string {
  return currentLocale;
}

/**
 * Translate a dot-notation key using the current reactive locale.
 * Components that call t() in templates will re-render automatically
 * when setLocale() changes currentLocale (because it's a $state rune).
 */
export function t(
  key: string,
  params?: Record<string, string | number>,
): string {
  const parts = key.split(".");
  let value: any = translations[currentLocale] || translations["en-US"];

  for (const part of parts) {
    if (value && typeof value === "object" && part in value) {
      value = value[part];
    } else {
      // Fallback to English
      let fallbackValue: any = translations["en-US"];
      for (const fbPart of parts) {
        if (
          fallbackValue &&
          typeof fallbackValue === "object" &&
          fbPart in fallbackValue
        ) {
          fallbackValue = fallbackValue[fbPart];
        } else {
          fallbackValue = key;
          break;
        }
      }
      value = fallbackValue;
      break;
    }
  }

  if (typeof value !== "string") {
    return key;
  }

  if (params) {
    let result = value;
    for (const [k, v] of Object.entries(params)) {
      result = result.replace(new RegExp(`{${k}}`, "g"), String(v));
    }
    return result;
  }

  return value;
}

/**
 * Retorna a chave do locale usada para um botão em determinada rota
 */
export function getButtonTranslationKey(route: string, btnKey: string): string {
  // Mapeamentos específicos por rota
  const routeSpecificKeys: Record<string, Record<string, string>> = {
    servers: {
      manage_vps_btn: "devices.manage_vps",
      add_first_btn: "devices.add_first",
      connect_btn: "devices.activate",
    },
    images: {
      create_container_btn: "images.create_container_btn",
      pull_btn: "images.pull_btn",
      prune_btn: "images.delete_selected",
      delete_btn: "images.delete_btn",
      select_all_btn: "common.select_all",
    },
    containers: {
      start_btn: "containers.start",
      stop_btn: "containers.stop",
      restart_btn: "containers.restart",
      delete_btn: "containers.delete",
      select_all_btn: "common.select_all",
      refresh_btn: "common.refresh",
      view_labels_btn: "containers.card_view_labels",
      view_env_btn: "containers.card_view_env",
      view_logs_btn: "containers.card_view_logs",
    },
    volumes: {
      new_volume_btn: "volumes.new_volume",
      prune_btn: "volumes.prune_btn",
      delete_btn: "volumes.delete_selected",
      select_all_btn: "common.select_all",
      refresh_btn: "common.refresh",
    },
    networks: {
      new_network_btn: "networks.new_network",
      prune_btn: "networks.prune_btn",
      delete_btn: "networks.delete_selected",
      select_all_btn: "common.select_all",
      refresh_btn: "common.refresh",
    },
    stacks: {
      new_stack_btn: "stacks.new_stack",
      deploy_btn: "stacks.deploy_btn",
      stop_btn: "stacks.stop_btn",
      logs_btn: "stacks.logs_btn",
      edit_btn: "stacks.edit_btn",
      remove_remote_btn: "stacks.remove_remote_btn",
      delete_local_btn: "stacks.delete_local_btn",
      browse_folder_btn: "stacks.browse_folder_btn",
    },
    builder: {
      build_btn: "builder.build_btn",
      browse_folder_btn: "builder.select_folder",
    },
    config: {
      edit_btn: "config.edit",
      delete_btn: "config.remove",
      backup_btn: "config.db_backup_btn",
      restore_btn: "config.db_restore_btn",
      reset_btn: "config.db_reset_btn",
    },
    profiles: {
      new_profile_btn: "profiles.new_profile",
      select_btn: "profiles.select_btn",
    },
  };

  if (routeSpecificKeys[route]?.[btnKey]) {
    return routeSpecificKeys[route][btnKey];
  }

  return `${route}.${btnKey}`;
}

/**
 * Atualiza o valor de uma chave no dicionário em memória e salva no arquivo .json em disco
 */
export async function updateLocaleTranslationKey(
  key: string,
  newValue: string,
  targetLocale?: string
): Promise<boolean> {
  const loc = targetLocale || currentLocale;
  if (!translations[loc]) {
    translations[loc] = {};
  }

  const parts = key.split(".");
  let currentObj = translations[loc];

  for (let i = 0; i < parts.length - 1; i++) {
    const part = parts[i];
    if (!currentObj[part] || typeof currentObj[part] !== "object") {
      currentObj[part] = {};
    }
    currentObj = currentObj[part];
  }

  const lastPart = parts[parts.length - 1];
  currentObj[lastPart] = newValue;

  // Persistir no arquivo ~/Documents/DockSea/locales/{loc}.json
  try {
    if (typeof LocaleWails.SaveLocale === "function") {
      const jsonContent = JSON.stringify(translations[loc], null, 2);
      await LocaleWails.SaveLocale(loc, jsonContent);
      return true;
    }
  } catch (err) {
    console.warn(`[Locale] Falha ao salvar arquivo '${loc}.json' em disco:`, err);
  }
  return false;
}
