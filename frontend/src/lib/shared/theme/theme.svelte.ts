import { defaultTheme, type DockSeaTheme, type ButtonRouteStyle } from "./defaultTheme";
import * as ThemeWails from "$bindings/theme/themeservice.js";

const THEME_STORAGE_KEY = "docksea_custom_theme";

function applyThemeToDom(theme: DockSeaTheme) {
  if (typeof document === "undefined") return;
  const root = document.documentElement;

  // 1. Variáveis Globais
  if (theme.global) {
    root.style.setProperty("--ds-app-bg", theme.global.appBg);
    root.style.setProperty("--ds-sidebar-from", theme.global.sidebarFrom);
    if (theme.global.sidebarVia) {
      root.style.setProperty("--ds-sidebar-via", theme.global.sidebarVia);
    }
    root.style.setProperty("--ds-sidebar-to", theme.global.sidebarTo);
    root.style.setProperty("--ds-card-bg", theme.global.cardBg);
    root.style.setProperty("--ds-card-border", theme.global.cardBorder);
    root.style.setProperty("--ds-card-shadow", theme.global.cardShadow);
    root.style.setProperty("--ds-text-primary", theme.global.textPrimary);
    root.style.setProperty("--ds-text-muted", theme.global.textMuted);
  }

  // 2. Variáveis de Rotas/Telas
  if (theme.routes) {
    for (const [routeName, buttons] of Object.entries(theme.routes)) {
      for (const [btnKey, style] of Object.entries(buttons)) {
        const prefix = `--ds-${routeName}-${btnKey}`;
        const bg = style.bg || (style as any).color;
        if (bg) {
          root.style.setProperty(`${prefix}-bg`, bg);
        }
        if (style.hover) {
          root.style.setProperty(`${prefix}-hover`, style.hover);
        }
        if (style.text) {
          root.style.setProperty(`${prefix}-text`, style.text);
        }
        if (style.textHover) {
          root.style.setProperty(`${prefix}-text-hover`, style.textHover);
        }
      }
    }
  }
}

class ThemeStore {
  currentTheme = $state<DockSeaTheme>(defaultTheme);
  isEditorOpen = $state(false);
  availableDiskThemes = $state<string[]>([]);
  themesDirectoryPath = $state<string>("");

  constructor() {
    this.loadTheme();
    this.refreshDiskThemes();
  }

  async refreshDiskThemes() {
    try {
      if (typeof ThemeWails.ListThemes === "function") {
        this.availableDiskThemes = (await ThemeWails.ListThemes()) || [];
        this.themesDirectoryPath = await ThemeWails.GetThemesDir();
      }
    } catch (e) {
      console.warn("Wails ThemeService não disponível no momento:", e);
    }
  }

  async saveThemeToDisk(name: string): Promise<boolean> {
    try {
      const cleanName = name.trim() || "default";
      await ThemeWails.SaveTheme(cleanName, this.exportThemeJson());
      await this.refreshDiskThemes();
      return true;
    } catch (e) {
      console.error("Falha ao salvar tema no disco:", e);
      return false;
    }
  }

  async loadThemeFromDisk(name: string): Promise<boolean> {
    try {
      const jsonContent = await ThemeWails.LoadTheme(name);
      if (jsonContent) {
        return this.importTheme(jsonContent);
      }
      return false;
    } catch (e) {
      console.error("Falha ao carregar tema do disco:", e);
      return false;
    }
  }

  async deleteThemeFromDisk(name: string): Promise<boolean> {
    try {
      await ThemeWails.DeleteTheme(name);
      await this.refreshDiskThemes();
      return true;
    } catch (e) {
      console.error("Falha ao deletar tema do disco:", e);
      return false;
    }
  }

  loadTheme() {
    if (typeof localStorage === "undefined") return;
    try {
      const saved = localStorage.getItem(THEME_STORAGE_KEY);
      if (saved) {
        const parsed = JSON.parse(saved);
        const mergedRoutes: Record<string, any> = {};
        const allRouteKeys = new Set([
          ...Object.keys(defaultTheme.routes || {}),
          ...Object.keys(parsed.routes || {}),
        ]);

        for (const routeKey of allRouteKeys) {
          mergedRoutes[routeKey] = {
            ...(defaultTheme.routes[routeKey] || {}),
            ...(parsed.routes?.[routeKey] || {}),
          };
        }

        this.currentTheme = {
          ...defaultTheme,
          ...parsed,
          global: { ...defaultTheme.global, ...(parsed.global || {}) },
          routes: mergedRoutes as any,
        };
      } else {
        this.currentTheme = JSON.parse(JSON.stringify(defaultTheme));
      }
    } catch (e) {
      console.warn("Erro ao carregar tema:", e);
      this.currentTheme = { ...defaultTheme };
    }
    applyThemeToDom(this.currentTheme);
  }

  getRouteButtonStyle(route: string, btnKey: string): ButtonRouteStyle | undefined {
    return this.currentTheme.routes?.[route]?.[btnKey];
  }

  setRouteButtonProp(
    route: string,
    btnKey: string,
    prop: keyof ButtonRouteStyle,
    value: string
  ) {
    if (!this.currentTheme.routes[route]) {
      this.currentTheme.routes[route] = {};
    }
    if (!this.currentTheme.routes[route][btnKey]) {
      this.currentTheme.routes[route][btnKey] = {
        bg: "#2563eb",
        hover: "#1d4ed8",
        text: "#ffffff",
        textHover: "#ffffff",
        size: "sm",
      };
    }
    (this.currentTheme.routes[route][btnKey] as any)[prop] = value;
    applyThemeToDom(this.currentTheme);
    this.saveTheme();
  }

  setRouteButtonColor(route: string, btnKey: string, color: string) {
    this.setRouteButtonProp(route, btnKey, "bg", color);
  }

  setRouteButtonHover(route: string, btnKey: string, hoverColor: string) {
    this.setRouteButtonProp(route, btnKey, "hover", hoverColor);
  }

  setRouteButtonText(route: string, btnKey: string, textColor: string) {
    this.setRouteButtonProp(route, btnKey, "text", textColor);
  }

  setRouteButtonSize(route: string, btnKey: string, size: "xs" | "sm" | "md" | "lg") {
    this.setRouteButtonProp(route, btnKey, "size", size);
  }

  setGlobalColor(key: keyof DockSeaTheme["global"], value: string) {
    this.currentTheme.global[key] = value;
    applyThemeToDom(this.currentTheme);
    this.saveTheme();
  }

  saveTheme() {
    if (typeof localStorage === "undefined") return;
    try {
      localStorage.setItem(THEME_STORAGE_KEY, JSON.stringify(this.currentTheme));
    } catch (e) {
      console.error("Falha ao salvar tema:", e);
    }
  }

  importTheme(jsonString: string): boolean {
    try {
      const parsed = JSON.parse(jsonString);
      if (!parsed || typeof parsed !== "object") return false;
      this.currentTheme = {
        ...defaultTheme,
        ...parsed,
        global: { ...defaultTheme.global, ...(parsed.global || {}) },
        routes: { ...defaultTheme.routes, ...(parsed.routes || {}) },
      };
      applyThemeToDom(this.currentTheme);
      this.saveTheme();
      return true;
    } catch (e) {
      console.error("JSON de tema inválido:", e);
      return false;
    }
  }

  resetToDefault() {
    this.currentTheme = JSON.parse(JSON.stringify(defaultTheme));
    applyThemeToDom(this.currentTheme);
    if (typeof localStorage !== "undefined") {
      localStorage.removeItem(THEME_STORAGE_KEY);
    }
  }

  exportThemeJson(): string {
    return JSON.stringify(this.currentTheme, null, 2);
  }

  toggleEditor() {
    this.isEditorOpen = !this.isEditorOpen;
    if (this.isEditorOpen) {
      this.refreshDiskThemes();
    }
  }
}

export const themeStore = new ThemeStore();
