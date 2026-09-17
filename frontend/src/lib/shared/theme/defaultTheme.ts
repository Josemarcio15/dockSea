import defaultThemeJson from "./theme_default.json";

export interface ButtonRouteStyle {
  bg?: string;        // Cor de fundo padrão (ex: "#10b981")
  hover?: string;     // Cor de fundo no hover (ex: "#059669")
  text?: string;      // Cor do texto (ex: "#ffffff")
  textHover?: string; // Cor do texto no hover (opcional, ex: "#ffffff")
  size?: "xs" | "sm" | "md" | "lg";
}

export interface DockSeaTheme {
  name: string;
  author: string;
  version: string;
  global: {
    appBg: string;
    sidebarFrom: string;
    sidebarVia?: string;
    sidebarTo: string;
    cardBg: string;
    cardBorder: string;
    cardShadow: string;
    textPrimary: string;
    textMuted: string;
  };
  routes: Record<string, Record<string, ButtonRouteStyle>>;
}

// A fonte padrão é o arquivo JSON físico
export const defaultTheme: DockSeaTheme = defaultThemeJson as DockSeaTheme;
