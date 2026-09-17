<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import FormModal from "$shared/components/FormModal.svelte";
  import { themeStore } from "$shared/theme/theme.svelte.ts";

  let {
    show = $bindable(false),
    id,
    name = $bindable(""),
    locale = $bindable("pt-BR"),
    theme = $bindable("default"),
    onSave,
  }: {
    show?: boolean;
    id?: string;
    name?: string;
    locale?: string;
    theme?: string;
    onSave?: () => void;
  } = $props();

  $effect(() => {
    if (show) {
      themeStore.refreshDiskThemes();
    }
  });
</script>

<FormModal
  bind:show
  title={id ? t("profiles.edit_title") : t("profiles.create_title")}
  buttons={[
    {
      label: t("common.save"),
      variant: "primary",
      themeKey: "profiles.save_btn",
      onclick: onSave,
      disabled: !name.trim(),
    },
  ]}
>
  <div class="flex flex-col gap-4">
    <!-- Nome do Perfil -->
    <div class="flex flex-col gap-1.5">
      <label
        for="profile-name"
        class="text-xs font-bold text-slate-400 uppercase tracking-wider"
      >
        {t("profiles.field_name")}
      </label>
      <input
        id="profile-name"
        bind:value={name}
        placeholder={t("profiles.placeholder_name")}
        class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500"
      />
    </div>

    <!-- Idioma Pré-definido -->
    <div class="flex flex-col gap-1.5">
      <label
        for="profile-locale"
        class="text-xs font-bold text-slate-400 uppercase tracking-wider"
      >
        Idioma / Language
      </label>
      <select
        id="profile-locale"
        bind:value={locale}
        class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 cursor-pointer"
      >
        <option value="pt-BR">🇧🇷 Português (Brasil)</option>
        <option value="en-US">🇺🇸 English (US)</option>
      </select>
    </div>

    <!-- Tema Pré-definido -->
    <div class="flex flex-col gap-1.5">
      <label
        for="profile-theme"
        class="text-xs font-bold text-slate-400 uppercase tracking-wider"
      >
        Tema Padrão do Usuário
      </label>
      <select
        id="profile-theme"
        bind:value={theme}
        class="w-full px-3.5 py-2.5 text-sm border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-slate-900/40 text-slate-900 dark:text-slate-100 focus:outline-none focus:border-violet-500 cursor-pointer font-mono"
      >
        <option value="default">DockSea Dark Classic (default.json)</option>
        {#each themeStore.availableDiskThemes as diskTheme}
          {#if diskTheme !== "default"}
            <option value={diskTheme}>{diskTheme}.json</option>
          {/if}
        {/each}
      </select>
      <p class="text-[10px] text-slate-500 dark:text-slate-400">
        Temas disponíveis na pasta ~/Documents/DockSea/themes/
      </p>
    </div>
  </div>
</FormModal>
