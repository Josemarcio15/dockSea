import { defaultTheme, type DockSeaTheme, type ButtonRouteStyle, type TextStyle } from "./defaultTheme";
import * as ThemeWails from "$bindings/theme/themeservice.js";

export function isPredefinedThemeName(name: string): boolean {
  const clean = name.toLowerCase().replace(/\.json$/, "").trim();
  return clean === "default" || clean === "docksea_dark_classic" || clean === "docksea_default";
}

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

  // 3. Variáveis de Textos / Labels
  if (theme.texts) {
    for (const [textKey, style] of Object.entries(theme.texts)) {
      const safeKey = textKey.replace(/[^a-zA-Z0-9_-]/g, "-");
      const prefix = `--ds-txt-${safeKey}`;
      if (style.color) {
        root.style.setProperty(`${prefix}-color`, style.color);
      }
      if (style.bg) {
        root.style.setProperty(`${prefix}-bg`, style.bg);
      }
      if (style.size) {
        root.style.setProperty(`${prefix}-size`, style.size);
      }
    }
  }
}

class ThemeStore {
  // Tema ativo no aplicativo
  currentTheme = $state<DockSeaTheme>(defaultTheme);
  activeThemeName = $state<string>("default");

  // Estado do Editor
  isEditorOpen = $state(false);
  editorStep = $state<"select" | "copy" | "edit">("select");
  
  // Tema em edição (draft) antes de salvar ou cancelar
  editingTheme = $state<DockSeaTheme>(JSON.parse(JSON.stringify(defaultTheme)));
  editingThemeName = $state<string>("");
  isEditingPredefined = $state<boolean>(false);

  availableDiskThemes = $state<string[]>([]);
  themesDirectoryPath = $state<string>("");

  constructor() {
    this.currentTheme = this.normalizeTheme(defaultTheme);
    applyThemeToDom(this.currentTheme);
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

  normalizeTheme(rawTheme: any): DockSeaTheme {
    const normalized: DockSeaTheme = {
      name: rawTheme?.name || defaultTheme.name,
      author: rawTheme?.author || defaultTheme.author,
      version: rawTheme?.version || defaultTheme.version,
      global: {
        ...defaultTheme.global,
        ...(rawTheme?.global || {}),
      },
      routes: {},
    };

    const allRouteKeys = new Set([
      ...Object.keys(defaultTheme.routes || {}),
      ...Object.keys(rawTheme?.routes || {}),
    ]);

    for (const routeKey of allRouteKeys) {
      const defRoute = defaultTheme.routes[routeKey] || {};
      const rawRoute = rawTheme?.routes?.[routeKey] || {};
      normalized.routes[routeKey] = {};

      const allBtnKeys = new Set([
        ...Object.keys(defRoute),
        ...Object.keys(rawRoute),
      ]);

      for (const btnKey of allBtnKeys) {
        const defBtn = defRoute[btnKey] || {};
        const rawBtn = rawRoute[btnKey] || {};

        const bg = rawBtn.bg || rawBtn.color || defBtn.bg || "#2563eb";
        const hover = rawBtn.hover || defBtn.hover || bg;
        const text = rawBtn.text || defBtn.text || "#ffffff";
        const textHover = rawBtn.textHover || defBtn.textHover || text;
        const size = rawBtn.size || defBtn.size || "sm";

        normalized.routes[routeKey][btnKey] = {
          bg,
          hover,
          text,
          textHover,
          size,
        };
      }
    }

    normalized.texts = {
      ...(defaultTheme.texts || {}),
      ...(rawTheme?.texts || {}),
    };

    return normalized;
  }

  // Iniciar fluxo ao clicar no botão "🎨 Personalizar Tema"
  openEditor() {
    this.refreshDiskThemes();
    this.editorStep = "select";
    this.isEditorOpen = true;
  }

  closeEditor() {
    this.isEditorOpen = false;
    // Restaura o DOM para o tema ativo caso estivesse visualizando preview
    applyThemeToDom(this.currentTheme);
  }

  // Aplicar tema diretamente (para usar no app sem necessariamente editar)
  async applyTheme(name: string): Promise<boolean> {
    if (isPredefinedThemeName(name)) {
      this.currentTheme = this.normalizeTheme(defaultTheme);
      this.activeThemeName = "default";
      applyThemeToDom(this.currentTheme);
      return true;
    }

    try {
      const jsonContent = await ThemeWails.LoadTheme(name);
      if (jsonContent) {
        const parsed = JSON.parse(jsonContent);
        this.currentTheme = this.normalizeTheme(parsed);
        this.activeThemeName = name;
        applyThemeToDom(this.currentTheme);
        return true;
      }
      return false;
    } catch (e) {
      console.error("Falha ao carregar tema do disco:", e);
      return false;
    }
  }

  // Escolher um tema no editor
  async selectThemeForEdit(name: string) {
    if (isPredefinedThemeName(name)) {
      // Tema predefinido: obriga a criar uma cópia
      this.editingTheme = JSON.parse(JSON.stringify(defaultTheme));
      this.editingThemeName = "";
      this.isEditingPredefined = true;
      this.editorStep = "copy";
    } else {
      // Tema customizado existente: carrega para editar diretamente ou copiar
      try {
        const jsonContent = await ThemeWails.LoadTheme(name);
        if (jsonContent) {
          this.editingTheme = this.normalizeTheme(JSON.parse(jsonContent));
          this.editingThemeName = name;
          this.isEditingPredefined = false;
          this.editorStep = "edit";
          // Aplica preview no DOM
          applyThemeToDom(this.editingTheme);
        }
      } catch (e) {
        console.error("Falha ao carregar tema:", e);
      }
    }
  }

  // Criar cópia com novo nome e começar a editar
  startEditNewCopy(newThemeName: string): boolean {
    const clean = newThemeName.trim();
    if (!clean) return false;
    if (isPredefinedThemeName(clean)) {
      alert("Não é permitido usar o nome de um tema predefinido.");
      return false;
    }
    this.editingThemeName = clean;
    this.editingTheme.name = clean;
    this.editorStep = "edit";
    applyThemeToDom(this.editingTheme);
    return true;
  }

  // Edição de propriedades durante o draft
  setDraftRouteButtonProp(
    route: string,
    btnKey: string,
    prop: keyof ButtonRouteStyle,
    value: string
  ) {
    if (!this.editingTheme.routes) {
      this.editingTheme.routes = {};
    }
    if (!this.editingTheme.routes[route]) {
      this.editingTheme.routes[route] = {};
    }
    if (!this.editingTheme.routes[route][btnKey]) {
      this.editingTheme.routes[route][btnKey] = {
        bg: "#2563eb",
        hover: "#1d4ed8",
        text: "#ffffff",
        textHover: "#ffffff",
        size: "sm",
      };
    }
    (this.editingTheme.routes[route][btnKey] as any)[prop] = value;
    // Forçar atualização reativa profunda no Svelte 5
    this.editingTheme = { ...this.editingTheme };
    applyThemeToDom(this.editingTheme);
  }

  setDraftTextProp(
    textKey: string,
    prop: keyof TextStyle,
    value: string
  ) {
    if (!this.editingTheme.texts) {
      this.editingTheme.texts = {};
    }
    if (!this.editingTheme.texts[textKey]) {
      this.editingTheme.texts[textKey] = {};
    }
    (this.editingTheme.texts[textKey] as any)[prop] = value;
    this.editingTheme = { ...this.editingTheme };
    applyThemeToDom(this.editingTheme);
  }

  setDraftGlobalColor(key: keyof DockSeaTheme["global"], value: string) {
    this.editingTheme.global[key] = value;
    this.editingTheme = { ...this.editingTheme };
    applyThemeToDom(this.editingTheme);
  }

  // Ação de SALVAR: grava no arquivo .json e torna o tema ativo (opcionalmente fecha o editor)
  async saveDraftToDisk(closeAfterSave = true): Promise<boolean> {
    try {
      const cleanName = this.editingThemeName.trim();
      if (!cleanName) {
        alert("Nome do tema inválido.");
        return false;
      }
      if (isPredefinedThemeName(cleanName)) {
        alert("Não é permitido sobrescrever temas predefinidos.");
        return false;
      }

      const jsonStr = JSON.stringify(this.editingTheme, null, 2);
      await ThemeWails.SaveTheme(cleanName, jsonStr);
      await this.refreshDiskThemes();

      // Torna o tema ativo
      this.currentTheme = JSON.parse(JSON.stringify(this.editingTheme));
      this.activeThemeName = cleanName;
      applyThemeToDom(this.currentTheme);
      if (closeAfterSave) {
        this.closeEditor();
      }
      return true;
    } catch (e) {
      console.error("Falha ao salvar tema no disco:", e);
      return false;
    }
  }

  // Ação de CANCELAR: descarta alterações e restaura o tema ativo anterior
  cancelDraft() {
    applyThemeToDom(this.currentTheme);
    this.closeEditor();
  }

  // Ação de RESTAURAR DRAFT para o padrão enquanto edita
  resetDraftToDefault() {
    this.editingTheme = this.normalizeTheme(defaultTheme);
    this.editingTheme.name = this.editingThemeName;
    applyThemeToDom(this.editingTheme);
  }

  async deleteThemeFromDisk(name: string): Promise<boolean> {
    if (isPredefinedThemeName(name)) {
      alert("Não é permitido excluir o tema predefinido.");
      return false;
    }
    try {
      await ThemeWails.DeleteTheme(name);
      await this.refreshDiskThemes();
      if (this.activeThemeName === name) {
        await this.applyTheme("default");
      }
      return true;
    } catch (e) {
      console.error("Falha ao deletar tema do disco:", e);
      return false;
    }
  }

  exportEditingJson(): string {
    return JSON.stringify(this.editingTheme, null, 2);
  }

  importEditingJson(jsonString: string): boolean {
    try {
      const parsed = JSON.parse(jsonString);
      if (!parsed || typeof parsed !== "object") return false;
      this.editingTheme = this.normalizeTheme(parsed);
      this.editingTheme.name = this.editingThemeName;
      applyThemeToDom(this.editingTheme);
      return true;
    } catch (e) {
      console.error("JSON inválido:", e);
      return false;
    }
  }

  toggleEditor() {
    if (this.isEditorOpen) {
      this.closeEditor();
    } else {
      this.openEditor();
    }
  }
}

export const themeStore = new ThemeStore();
