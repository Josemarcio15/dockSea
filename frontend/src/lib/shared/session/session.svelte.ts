import * as ProfileService from "$bindings/profiles/service.js";
import * as ServerService from "$bindings/servers/service.js";
import { setLocale } from "$shared/stores/locale.svelte";
import { themeStore } from "$shared/theme/theme.svelte";
import type { SessionState } from "./session.types";

export const session = $state<SessionState>({
  servers: [],
  activeVps: null,
  profiles: [],
  activeProfile: {
    name: "Perfil Padrão",
    locale: "pt-BR",
    theme: "default",
  },
});

export async function loadSession(): Promise<void> {
  const servers = await ServerService.ListServers();
  session.servers = servers || [];
  session.activeVps =
    (servers && servers.find((server: any) => server.isActive)) || null;

  const profiles = await ProfileService.ListProfiles();
  if (profiles && profiles.length > 0) {
    session.profiles = profiles;
    const active = profiles.find((profile: any) => profile.isActive) || profiles[0];
    session.activeProfile = active;

    // Sincronizar idioma do perfil
    if (active.locale) {
      setLocale(active.locale);
    }
    // Sincronizar tema do perfil
    if (active.theme) {
      void themeStore.applyTheme(active.theme);
    }
  } else {
    session.profiles = [];
  }
}
