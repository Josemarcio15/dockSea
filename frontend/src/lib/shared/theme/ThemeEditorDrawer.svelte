<script lang="ts">
  import { themeStore, isPredefinedThemeName } from "./theme.svelte.ts";
  import { defaultTheme } from "./defaultTheme";
  import { navigation } from "$navigation/navigation.svelte";
  import {
    t,
    getButtonTranslationKey,
    updateLocaleTranslationKey,
    getFlattenedTranslations,
    localeVersionState,
  } from "$shared/stores/locale.svelte";
  import { Button, EditButtonIcon, TrashButtonIcon } from "$shared/components/buttons";

  let activeTab = $state<"route" | "text" | "global" | "json">("route");
  let newCopyName = $state("");
  let copyFeedback = $state(false);
  let importFeedback = $state("");

  // Text Tab state
  let textSearchQuery = $state("");
  let textCategory = $state("current");
  let editingTextKey = $state<string | null>(null);
  let editingTextValue = $state("");

  // Route selector in Route tab (defaults to currentRoute)
  let selectedRouteOverride = $state<string>("");

  function focusElement(node: HTMLElement) {
    node.focus();
  }

  let currentRoute = $derived(navigation.currentRoute || "stacks");
  let activeEditorRoute = $derived(selectedRouteOverride || currentRoute);

  const routeTitles = $derived<Record<string, string>>({
    servers: t("sidebar.devices"),
    containers: t("sidebar.containers"),
    images: t("sidebar.images"),
    volumes: t("sidebar.volumes"),
    networks: t("sidebar.networks"),
    stacks: t("sidebar.stacks"),
    builder: t("sidebar.builder"),
    extras: t("extras.title"),
    profiles: t("sidebar.profiles"),
    config: t("sidebar.configs"),
  });

  // Helper uniforme para buscar o alias traduzido do botão baseado na rota e chave
  function getButtonAlias(route: string, key: string): string {
    const translationKey = getButtonTranslationKey(route, key);
    const translated = t(translationKey);
    if (translated && translated !== translationKey) {
      return translated;
    }
    return key.replace(/_/g, " ");
  }

  const currentRouteButtons = $derived.by(() => {
    const routeObj =
      themeStore.editingTheme.routes?.[activeEditorRoute] ||
      defaultTheme.routes?.[activeEditorRoute] ||
      {};
    return Object.entries(routeObj).map(([key, style]) => ({
      key,
      label: key,
      alias: getButtonAlias(activeEditorRoute, key),
      bg: style.bg || (style as any).color || "#2563eb",
      hover: style.hover || style.bg || "#1d4ed8",
      text: style.text || "#ffffff",
      textHover: style.textHover || style.text || "#ffffff",
      size: style.size || "sm",
    }));
  });

  let editingAliasKey = $state<string | null>(null);
  let editingAliasValue = $state("");

  function startEditAlias(btn: { key: string; alias: string }) {
    editingAliasKey = btn.key;
    editingAliasValue = btn.alias;
  }

  async function saveAlias(route: string, btnKey: string) {
    if (editingAliasValue.trim()) {
      const localeKey = getButtonTranslationKey(route, btnKey);
      await updateLocaleTranslationKey(localeKey, editingAliasValue.trim());
    }
    editingAliasKey = null;
    editingAliasValue = "";
  }

  function cancelEditAlias() {
    editingAliasKey = null;
    editingAliasValue = "";
  }

  function handleBgChange(route: string, btnKey: string, color: string) {
    themeStore.setDraftRouteButtonProp(route, btnKey, "bg", color);
  }

  function handleHoverChange(route: string, btnKey: string, color: string) {
    themeStore.setDraftRouteButtonProp(route, btnKey, "hover", color);
  }

  function handleTextChange(route: string, btnKey: string, color: string) {
    themeStore.setDraftRouteButtonProp(route, btnKey, "text", color);
  }

  function handleTextHoverChange(route: string, btnKey: string, color: string) {
    themeStore.setDraftRouteButtonProp(route, btnKey, "textHover", color);
  }

  function handleSizeChange(route: string, btnKey: string, size: any) {
    themeStore.setDraftRouteButtonProp(route, btnKey, "size", size);
  }

  // --- Funções da Aba Text ---
  const allFlattenedTexts = $derived(getFlattenedTranslations());

  const filteredTexts = $derived.by(() => {
    const query = textSearchQuery.toLowerCase().trim();
    return allFlattenedTexts.filter((item) => {
      // Filtro por Categoria / Rota
      if (textCategory === "current") {
        const route = activeEditorRoute;
        const prefix = route === "servers" ? "devices" : route;
        if (!item.key.startsWith(prefix + ".") && !item.key.startsWith(route + ".")) {
          return false;
        }
      } else if (textCategory !== "all") {
        if (!item.key.startsWith(textCategory + ".") && item.key !== textCategory) {
          return false;
        }
      }

      // Filtro por Texto ou Chave
      if (query) {
        return (
          item.key.toLowerCase().includes(query) ||
          item.value.toLowerCase().includes(query)
        );
      }
      return true;
    });
  });

  function startEditText(key: string, val: string) {
    editingTextKey = key;
    editingTextValue = val;
  }

  async function saveText(key: string) {
    if (editingTextValue.trim()) {
      await updateLocaleTranslationKey(key, editingTextValue.trim());
    }
    editingTextKey = null;
    editingTextValue = "";
  }

  function cancelEditText() {
    editingTextKey = null;
    editingTextValue = "";
  }

  function handleTextColorChange(key: string, color: string) {
    themeStore.setDraftTextProp(key, "color", color);
  }

  function handleTextBgChange(key: string, bg: string) {
    themeStore.setDraftTextProp(key, "bg", bg);
  }

  function handleTextSizeChange(key: string, size: string) {
    themeStore.setDraftTextProp(key, "size", size);
  }

  // --- Exportar / Importar ---
  function handleExportFile() {
    const jsonStr = themeStore.exportEditingJson();
    const blob = new Blob([jsonStr], { type: "application/json" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${themeStore.editingThemeName || "theme"}.json`;
    a.click();
    URL.revokeObjectURL(url);
  }

  function handleCopyJson() {
    navigator.clipboard.writeText(themeStore.exportEditingJson());
    copyFeedback = true;
    setTimeout(() => (copyFeedback = false), 2000);
  }

  function handleFileImport(e: Event) {
    const input = e.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const reader = new FileReader();
      reader.onload = (event) => {
        const content = event.target?.result as string;
        const ok = themeStore.importEditingJson(content);
        if (ok) {
          importFeedback = t("theme_editor.import_success");
          setTimeout(() => (importFeedback = ""), 3000);
        } else {
          alert(t("theme_editor.import_invalid"));
        }
      };
      reader.readAsText(input.files[0]);
    }
  }

  function handleConfirmCopy() {
    if (!newCopyName.trim()) {
      alert(t("theme_editor.alert_enter_name"));
      return;
    }
    const ok = themeStore.startEditNewCopy(newCopyName.trim());
    if (ok) {
      newCopyName = "";
    }
  }
</script>

{#if themeStore.isEditorOpen}
  <div
    class="fixed bottom-5 right-5 z-50 w-[420px] max-w-[95vw] h-[60vh] min-h-[460px] max-h-[85vh] bg-slate-900/95 border border-slate-700/80 backdrop-blur-xl rounded-2xl shadow-2xl flex flex-col overflow-hidden text-slate-200 font-sans animate-in fade-in zoom-in-95 duration-200"
  >
    <!-- Header -->
    <div class="px-5 py-3.5 border-b border-slate-800 flex items-center justify-between bg-slate-950/50 shrink-0">
      <div class="flex items-center gap-2.5">
        <div>
          <h3 class="text-sm font-bold text-white tracking-wide">{t("theme_editor.title")}</h3>
          <p class="text-[10px] text-slate-400">
            {#if themeStore.editorStep === "select"}
              {t("theme_editor.step_select")}
            {:else if themeStore.editorStep === "copy"}
              {t("theme_editor.step_copy")}
            {:else}
              {t("theme_editor.step_editing", { name: themeStore.editingThemeName })}
            {/if}
          </p>
        </div>
      </div>
      <button
        onclick={() => themeStore.cancelDraft()}
        class="w-7 h-7 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white flex items-center justify-center transition-colors text-xs"
        title={t("common.close")}
      >
        ✕
      </button>
    </div>

    <!-- ETAPA 1: SELEÇÃO DE TEMA -->
    {#if themeStore.editorStep === "select"}
      <div class="p-5 flex-1 overflow-y-auto space-y-4 text-xs">
        <p class="text-slate-400 text-[11px]">
          {t("theme_editor.select_prompt")}
        </p>

        <!-- Tema Padrão do Sistema (Protegido) -->
        <div class="bg-slate-800/60 p-3.5 rounded-xl border border-slate-700/60 space-y-2">
          <div class="flex items-center justify-between">
            <div>
              <span class="font-bold text-white text-xs">DockSea Dark Classic</span>
              <span class="ml-2 text-[9px] bg-slate-700 px-1.5 py-0.5 rounded text-slate-300 uppercase tracking-wider">
                {t("theme_editor.default_badge")}
              </span>
            </div>
            {#if themeStore.activeThemeName === "default"}
              <span class="text-[10px] text-emerald-400 font-semibold">{t("theme_editor.active_badge")}</span>
            {/if}
          </div>
          <div class="flex gap-2 pt-1">
            {#if themeStore.activeThemeName !== "default"}
              <Button
                size="sm"
                class="flex-1 bg-slate-700 hover:bg-slate-600 text-white"
                onclick={() => themeStore.applyTheme("default")}
              >
                {t("theme_editor.btn_activate")}
              </Button>
            {/if}
            <Button
              size="sm"
              class="flex-1 bg-sky-600 hover:bg-sky-500 text-white"
              onclick={() => themeStore.selectThemeForEdit("default")}
            >
              {t("theme_editor.btn_create_copy_edit")}
            </Button>
          </div>
        </div>

        <!-- Temas Customizados na pasta ~/Documents/DockSea/themes -->
        <div class="space-y-2">
          <span class="text-[11px] font-semibold text-slate-400 block">{t("theme_editor.custom_themes_title")}</span>
          {#if themeStore.availableDiskThemes.length === 0}
            <div class="py-4 text-center text-slate-500 text-[11px] bg-slate-950/40 rounded-xl border border-slate-800/60">
              {t("theme_editor.no_custom_themes")}
            </div>
          {:else}
            {#each themeStore.availableDiskThemes as tName}
              <div class="bg-slate-800/60 p-3 rounded-xl border border-slate-700/50 flex items-center justify-between">
                <div>
                  <span class="font-bold text-white text-xs">{tName}.json</span>
                  {#if themeStore.activeThemeName === tName}
                    <span class="ml-2 text-[10px] text-emerald-400 font-semibold">{t("theme_editor.active_badge")}</span>
                  {/if}
                </div>
                <div class="flex items-center gap-1.5">
                  {#if themeStore.activeThemeName !== tName}
                    <Button
                      size="xs"
                      class="bg-slate-700 hover:bg-slate-600 text-white"
                      onclick={() => themeStore.applyTheme(tName)}
                    >
                      {t("theme_editor.btn_activate")}
                    </Button>
                  {/if}
                  <EditButtonIcon
                    size="xs"
                    class="bg-sky-600 hover:bg-sky-500 text-white"
                    onclick={() => themeStore.selectThemeForEdit(tName)}
                  >
                    {t("common.edit")}
                  </EditButtonIcon>
                  <TrashButtonIcon
                    size="xs"
                    class="bg-rose-600/80 hover:bg-rose-600 text-white"
                    onclick={() => {
                      if (confirm(t("theme_editor.delete_theme_confirm", { name: tName }))) {
                        themeStore.deleteThemeFromDisk(tName);
                      }
                    }}
                  />
                </div>
              </div>
            {/each}
          {/if}
        </div>
      </div>

    <!-- ETAPA 2: OBRIGATÓRIO CRIAR CÓPIA PARA EDITAR TEMA PREDEFINIDO -->
    {:else if themeStore.editorStep === "copy"}
      <div class="p-5 flex-1 overflow-y-auto space-y-4 text-xs">
        <div class="bg-amber-500/10 border border-amber-500/30 p-3 rounded-xl text-amber-300 text-[11px] leading-relaxed">
          <strong>{t("theme_editor.protected_theme_warn")}</strong>
        </div>

        <div class="space-y-2">
          <label for="new-copy-theme-name" class="block text-slate-300 font-semibold text-xs">{t("theme_editor.new_theme_name_label")}</label>
          <input
            id="new-copy-theme-name"
            type="text"
            bind:value={newCopyName}
            placeholder={t("theme_editor.new_theme_name_placeholder")}
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3.5 py-2 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-sky-500 font-mono"
          />
          <p class="text-[10px] text-slate-500">
            {t("theme_editor.new_theme_save_path_hint", { name: newCopyName.trim() || 'meu_tema' })}
          </p>
        </div>
      </div>

      <!-- Footer da tela de cópia -->
      <div class="px-5 py-3 border-t border-slate-800 bg-slate-950/60 flex items-center justify-between text-xs shrink-0">
        <Button
          size="sm"
          class="bg-transparent hover:bg-slate-800 text-slate-400 hover:text-white"
          onclick={() => (themeStore.editorStep = "select")}
        >
          {t("theme_editor.btn_back")}
        </Button>
        <Button
          size="sm"
          class="bg-sky-600 hover:bg-sky-500 active:bg-sky-700 text-white font-semibold"
          onclick={handleConfirmCopy}
        >
          {t("theme_editor.btn_start_editing")}
        </Button>
      </div>

    <!-- ETAPA 3: EDIÇÃO COM ABAS: TELA ATUAL | TEXTOS | FUNDO | JSON -->
    {:else if themeStore.editorStep === "edit"}
      <!-- Navegação de Abas -->
      <div class="flex border-b border-slate-800 bg-slate-950/30 text-xs font-semibold px-2 pt-1.5 gap-1 shrink-0 overflow-x-auto">
        <button
          onclick={() => (activeTab = "route")}
          class="px-3 py-1.5 rounded-t-lg border-b-2 whitespace-nowrap transition-all {activeTab === 'route'
            ? 'border-sky-500 text-sky-400 bg-slate-800/40 font-bold'
            : 'border-transparent text-slate-400 hover:text-slate-200'}"
        >
          {t("theme_editor.tab_current_screen")}
        </button>
        <button
          onclick={() => (activeTab = "text")}
          class="px-3 py-1.5 rounded-t-lg border-b-2 whitespace-nowrap transition-all {activeTab === 'text'
            ? 'border-sky-500 text-sky-400 bg-slate-800/40 font-bold'
            : 'border-transparent text-slate-400 hover:text-slate-200'}"
        >
          {t("theme_editor.tab_texts")}
        </button>
        <button
          onclick={() => (activeTab = "global")}
          class="px-3 py-1.5 rounded-t-lg border-b-2 whitespace-nowrap transition-all {activeTab === 'global'
            ? 'border-sky-500 text-sky-400 bg-slate-800/40 font-bold'
            : 'border-transparent text-slate-400 hover:text-slate-200'}"
        >
          {t("theme_editor.tab_background")}
        </button>
        <button
          onclick={() => (activeTab = "json")}
          class="px-3 py-1.5 rounded-t-lg border-b-2 whitespace-nowrap transition-all {activeTab === 'json'
            ? 'border-sky-500 text-sky-400 bg-slate-800/40 font-bold'
            : 'border-transparent text-slate-400 hover:text-slate-200'}"
        >
          {t("theme_editor.tab_json")}
        </button>
      </div>

      <!-- Conteúdo da Aba -->
      <div class="p-4 flex-1 overflow-y-auto text-xs space-y-3">
        <!-- ABA 1: BOTÕES DA ROTA / TELA ATUAL -->
        {#if activeTab === "route"}
          <div>
            <div class="flex items-center justify-between mb-2 gap-2">
              <select
                value={activeEditorRoute}
                onchange={(e) => (selectedRouteOverride = e.currentTarget.value)}
                class="bg-slate-900 border border-slate-700 font-bold text-sky-400 rounded-lg px-2 py-1 text-xs focus:outline-none cursor-pointer"
              >
                <option value="stacks">{routeTitles.stacks || "Stacks"}</option>
                <option value="containers">{routeTitles.containers || "Containers"}</option>
                <option value="images">{routeTitles.images || "Images"}</option>
                <option value="volumes">{routeTitles.volumes || "Volumes"}</option>
                <option value="networks">{routeTitles.networks || "Networks"}</option>
                <option value="servers">{routeTitles.servers || "Servers"}</option>
                <option value="builder">{routeTitles.builder || "Builder"}</option>
                <option value="config">{routeTitles.config || "Config"}</option>
                <option value="profiles">{routeTitles.profiles || "Profiles"}</option>
                <option value="extras">{routeTitles.extras || "Extras"}</option>
              </select>
              <span class="text-[10px] bg-slate-800 px-2 py-0.5 rounded text-slate-400 shrink-0">
                {t("theme_editor.buttons_count", { count: currentRouteButtons.length })}
              </span>
            </div>

            {#if currentRouteButtons.length > 0}
              <div class="space-y-2.5">
                {#each currentRouteButtons as btn}
                  <div class="flex flex-col bg-slate-800/50 p-2.5 rounded-xl border border-slate-700/50 gap-2">
                    <!-- Top: Key do Botão & Seletor de Tamanho -->
                    <div class="flex items-center justify-between gap-2">
                      <span class="font-mono font-bold text-sky-400 text-xs truncate">{btn.label}</span>
                      <!-- Size Selector -->
                      <select
                        value={btn.size}
                        onchange={(e) => handleSizeChange(activeEditorRoute, btn.key, e.currentTarget.value)}
                        class="bg-slate-900 border border-slate-700 rounded-lg px-1.5 py-0.5 text-[10px] text-slate-300 focus:outline-none shrink-0 cursor-pointer"
                      >
                        <option value="xs">{t("theme_editor.size_xs")}</option>
                        <option value="sm">{t("theme_editor.size_sm")}</option>
                        <option value="md">{t("theme_editor.size_md")}</option>
                        <option value="lg">{t("theme_editor.size_lg")}</option>
                      </select>
                    </div>

                    <!-- Linha divisória -->
                    <div class="border-t border-slate-700/40"></div>

                    <!-- Bottom: Alias / Tradução & Botão de Edição -->
                    <div class="flex items-center justify-between gap-2 min-h-[24px]">
                      {#if editingAliasKey === btn.key}
                        <form
                          onsubmit={(e) => {
                            e.preventDefault();
                            saveAlias(activeEditorRoute, btn.key);
                          }}
                          class="flex items-center gap-1 flex-1 min-w-0"
                        >
                          <input
                            type="text"
                            bind:value={editingAliasValue}
                            class="w-full bg-slate-950 border border-sky-500 rounded px-1.5 py-0.5 text-xs text-white focus:outline-none font-sans"
                            placeholder={t("theme_editor.new_alias_placeholder")}
                            use:focusElement
                          />
                          <button
                            type="submit"
                            class="text-emerald-400 hover:text-emerald-300 px-1 text-xs cursor-pointer font-bold"
                            title={t("theme_editor.btn_save_alias")}
                          >
                            ✓
                          </button>
                          <button
                            type="button"
                            onclick={cancelEditAlias}
                            class="text-slate-500 hover:text-slate-300 px-1 text-xs cursor-pointer"
                            title={t("common.cancel")}
                          >
                            ✕
                          </button>
                        </form>
                      {:else}
                        <span class="text-slate-300 font-medium text-xs truncate">{btn.alias}</span>
                        <EditButtonIcon
                          size="xs"
                          title={t("theme_editor.edit_alias_title")}
                          onclick={() => startEditAlias(btn)}
                          class="!p-1 bg-transparent hover:bg-slate-700 text-slate-400 hover:text-sky-400 shadow-none shrink-0"
                        />
                      {/if}
                    </div>

                    <div class="grid grid-cols-4 gap-1.5 pt-1 border-t border-slate-700/40 text-[9px]">
                      <!-- Fundo Normal -->
                      <label class="flex flex-col items-center gap-1 bg-slate-900/60 p-1 rounded-lg border border-slate-800 cursor-pointer">
                        <span class="text-slate-400">{t("theme_editor.color_bg")}</span>
                        <input
                          type="color"
                          value={btn.bg}
                          oninput={(e) => handleBgChange(activeEditorRoute, btn.key, e.currentTarget.value)}
                          class="w-5 h-5 rounded cursor-pointer bg-transparent border-0"
                        />
                      </label>

                      <!-- Hover Fundo -->
                      <label class="flex flex-col items-center gap-1 bg-slate-900/60 p-1 rounded-lg border border-slate-800 cursor-pointer">
                        <span class="text-slate-400">{t("theme_editor.color_hover")}</span>
                        <input
                          type="color"
                          value={btn.hover}
                          oninput={(e) => handleHoverChange(activeEditorRoute, btn.key, e.currentTarget.value)}
                          class="w-5 h-5 rounded cursor-pointer bg-transparent border-0"
                        />
                      </label>

                      <!-- Cor do Texto -->
                      <label class="flex flex-col items-center gap-1 bg-slate-900/60 p-1 rounded-lg border border-slate-800 cursor-pointer">
                        <span class="text-slate-400">{t("theme_editor.color_text")}</span>
                        <input
                          type="color"
                          value={btn.text}
                          oninput={(e) => handleTextChange(activeEditorRoute, btn.key, e.currentTarget.value)}
                          class="w-5 h-5 rounded cursor-pointer bg-transparent border-0"
                        />
                      </label>

                      <!-- Hover Texto -->
                      <label class="flex flex-col items-center gap-1 bg-slate-900/60 p-1 rounded-lg border border-slate-800 cursor-pointer">
                        <span class="text-slate-400">{t("theme_editor.color_text_hover")}</span>
                        <input
                          type="color"
                          value={btn.textHover}
                          oninput={(e) => handleTextHoverChange(activeEditorRoute, btn.key, e.currentTarget.value)}
                          class="w-5 h-5 rounded cursor-pointer bg-transparent border-0"
                        />
                      </label>
                    </div>
                  </div>
                {/each}
              </div>
            {:else}
              <div class="py-6 text-center text-slate-500">
                {t("theme_editor.no_buttons_route")}
              </div>
            {/if}
          </div>

        <!-- ABA 2: TEXTOS / LABELS DO JSON DE TRADUÇÃO -->
        {:else if activeTab === "text"}
          <div class="space-y-3">
            <!-- Barra de Filtro e Busca -->
            <div class="flex flex-col gap-2 bg-slate-950/40 p-2.5 rounded-xl border border-slate-800/80">
              <div class="flex items-center gap-2">
                <input
                  type="text"
                  bind:value={textSearchQuery}
                  placeholder={t("theme_editor.search_texts_placeholder")}
                  class="flex-1 bg-slate-900 border border-slate-700 rounded-lg px-2.5 py-1.5 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-sky-500"
                />
                {#if textSearchQuery}
                  <button
                    onclick={() => (textSearchQuery = "")}
                    class="text-slate-400 hover:text-white text-xs px-1"
                  >
                    ✕
                  </button>
                {/if}
              </div>

              <div class="flex items-center justify-between gap-2">
                <select
                  bind:value={textCategory}
                  class="flex-1 bg-slate-900 border border-slate-700 rounded-lg px-2 py-1 text-[11px] text-slate-300 focus:outline-none cursor-pointer"
                >
                  <option value="current">{t("theme_editor.category_current", { route: activeEditorRoute })}</option>
                  <option value="all">{t("theme_editor.category_all")}</option>
                  <option value="common">common</option>
                  <option value="sidebar">sidebar</option>
                  <option value="containers">containers</option>
                  <option value="images">images</option>
                  <option value="volumes">volumes</option>
                  <option value="networks">networks</option>
                  <option value="stacks">stacks</option>
                  <option value="builder">builder</option>
                  <option value="devices">devices / servers</option>
                  <option value="config">config</option>
                  <option value="profiles">profiles</option>
                  <option value="extras">extras</option>
                  <option value="theme_editor">theme_editor</option>
                </select>
                <span class="text-[10px] bg-slate-800 px-2 py-0.5 rounded text-slate-400 shrink-0">
                  {t("theme_editor.texts_count", { count: filteredTexts.length })}
                </span>
              </div>
            </div>

            <!-- Lista de Textos -->
            {#if filteredTexts.length > 0}
              <div class="space-y-2.5 max-h-[50vh] overflow-y-auto pr-1">
                {#each filteredTexts as item (item.key)}
                  {@const textStyle = themeStore.editingTheme.texts?.[item.key] || {}}
                  <div class="flex flex-col bg-slate-800/50 p-2.5 rounded-xl border border-slate-700/50 gap-2">
                    <!-- Top: Key do Texto & Seletor de Tamanho de Fonte -->
                    <div class="flex items-center justify-between gap-2">
                      <span class="font-mono font-bold text-sky-400 text-[11px] truncate select-all" title={item.key}>
                        {item.key}
                      </span>
                      <select
                        value={textStyle.size || ""}
                        onchange={(e) => handleTextSizeChange(item.key, e.currentTarget.value)}
                        class="bg-slate-900 border border-slate-700 rounded-lg px-1.5 py-0.5 text-[10px] text-slate-300 focus:outline-none shrink-0 cursor-pointer"
                      >
                        <option value="">{t("theme_editor.size_sm")}</option>
                        <option value="xs">{t("theme_editor.size_xs")}</option>
                        <option value="sm">{t("theme_editor.size_sm")}</option>
                        <option value="md">{t("theme_editor.size_md")}</option>
                        <option value="lg">{t("theme_editor.size_lg")}</option>
                        <option value="xl">{t("theme_editor.size_xl")}</option>
                        <option value="2xl">{t("theme_editor.size_2xl")}</option>
                      </select>
                    </div>

                    <!-- Linha divisória -->
                    <div class="border-t border-slate-700/40"></div>

                    <!-- Middle: Tradução Atual & Edição de Tradução -->
                    <div class="flex items-center justify-between gap-2 min-h-[24px]">
                      {#if editingTextKey === item.key}
                        <form
                          onsubmit={(e) => {
                            e.preventDefault();
                            saveText(item.key);
                          }}
                          class="flex items-center gap-1 flex-1 min-w-0"
                        >
                          <input
                            type="text"
                            bind:value={editingTextValue}
                            class="w-full bg-slate-950 border border-sky-500 rounded px-1.5 py-0.5 text-xs text-white focus:outline-none font-sans"
                            placeholder={t("theme_editor.new_text_placeholder")}
                            use:focusElement
                          />
                          <button
                            type="submit"
                            class="text-emerald-400 hover:text-emerald-300 px-1 text-xs cursor-pointer font-bold"
                            title={t("theme_editor.btn_save_text")}
                          >
                            ✓
                          </button>
                          <button
                            type="button"
                            onclick={cancelEditText}
                            class="text-slate-500 hover:text-slate-300 px-1 text-xs cursor-pointer"
                            title={t("common.cancel")}
                          >
                            ✕
                          </button>
                        </form>
                      {:else}
                        <span class="text-slate-200 font-medium text-xs break-words line-clamp-2" title={item.value}>
                          {item.value}
                        </span>
                        <EditButtonIcon
                          size="xs"
                          title={t("theme_editor.edit_text_title")}
                          onclick={() => startEditText(item.key, item.value)}
                          class="!p-1 bg-transparent hover:bg-slate-700 text-slate-400 hover:text-sky-400 shadow-none shrink-0"
                        />
                      {/if}
                    </div>

                    <!-- Bottom: Configuração de Cores (Texto & Fundo) -->
                    <div class="grid grid-cols-2 gap-2 pt-1 border-t border-slate-700/40 text-[10px]">
                      <!-- Cor do Texto -->
                      <label class="flex items-center justify-between bg-slate-900/60 px-2 py-1 rounded-lg border border-slate-800 cursor-pointer">
                        <span class="text-slate-400">{t("theme_editor.color_text")}</span>
                        <input
                          type="color"
                          value={textStyle.color || "#f8fafc"}
                          oninput={(e) => handleTextColorChange(item.key, e.currentTarget.value)}
                          class="w-5 h-5 rounded cursor-pointer bg-transparent border-0"
                        />
                      </label>

                      <!-- Cor do Fundo -->
                      <label class="flex items-center justify-between bg-slate-900/60 px-2 py-1 rounded-lg border border-slate-800 cursor-pointer">
                        <span class="text-slate-400">{t("theme_editor.color_bg")}</span>
                        <input
                          type="color"
                          value={textStyle.bg || "#111827"}
                          oninput={(e) => handleTextBgChange(item.key, e.currentTarget.value)}
                          class="w-5 h-5 rounded cursor-pointer bg-transparent border-0"
                        />
                      </label>
                    </div>
                  </div>
                {/each}
              </div>
            {:else}
              <div class="py-8 text-center text-slate-500">
                {t("theme_editor.no_texts_found")}
              </div>
            {/if}
          </div>

        <!-- ABA 3: CORES GLOBAIS / FUNDO -->
        {:else if activeTab === "global"}
          <div class="space-y-2.5">
            <p class="text-slate-400 text-[11px]">{t("theme_editor.global_colors_desc")}</p>
            
            <div class="flex items-center justify-between bg-slate-800/40 p-2 rounded-xl border border-slate-700/40">
              <span>{t("theme_editor.global_app_bg")}</span>
              <input
                type="color"
                value={themeStore.editingTheme.global.appBg}
                oninput={(e) => themeStore.setDraftGlobalColor("appBg", e.currentTarget.value)}
                class="w-7 h-7 rounded-lg cursor-pointer bg-transparent border-0"
              />
            </div>

            <div class="flex items-center justify-between bg-slate-800/40 p-2 rounded-xl border border-slate-700/40">
              <span>{t("theme_editor.global_card_bg")}</span>
              <input
                type="color"
                value={themeStore.editingTheme.global.cardBg}
                oninput={(e) => themeStore.setDraftGlobalColor("cardBg", e.currentTarget.value)}
                class="w-7 h-7 rounded-lg cursor-pointer bg-transparent border-0"
              />
            </div>

            <div class="flex items-center justify-between bg-slate-800/40 p-2 rounded-xl border border-slate-700/40">
              <span>{t("theme_editor.global_sidebar_from")}</span>
              <input
                type="color"
                value={themeStore.editingTheme.global.sidebarFrom}
                oninput={(e) => themeStore.setDraftGlobalColor("sidebarFrom", e.currentTarget.value)}
                class="w-7 h-7 rounded-lg cursor-pointer bg-transparent border-0"
              />
            </div>

            <div class="flex items-center justify-between bg-slate-800/40 p-2 rounded-xl border border-slate-700/40">
              <span>{t("theme_editor.global_sidebar_to")}</span>
              <input
                type="color"
                value={themeStore.editingTheme.global.sidebarTo}
                oninput={(e) => themeStore.setDraftGlobalColor("sidebarTo", e.currentTarget.value)}
                class="w-7 h-7 rounded-lg cursor-pointer bg-transparent border-0"
              />
            </div>
          </div>

        <!-- ABA 4: IMPORTAR / EXPORTAR JSON -->
        {:else if activeTab === "json"}
          <div class="space-y-4 py-1">
            <p class="text-slate-400 text-[11px] leading-relaxed">
              {t("theme_editor.json_tab_desc")}
            </p>

            <!-- Exportar Tema -->
            <div class="bg-slate-800/50 p-3 rounded-xl border border-slate-700/50 space-y-2">
              <span class="text-[11px] font-bold text-white block">{t("theme_editor.export_theme_title")}</span>
              <div class="flex gap-2">
                <Button
                  size="sm"
                  class="flex-1 bg-sky-600 hover:bg-sky-500 active:bg-sky-700 text-white font-semibold"
                  onclick={handleExportFile}
                >
                  {t("theme_editor.btn_download_json")}
                </Button>
                <Button
                  size="sm"
                  class="bg-slate-700 hover:bg-slate-600 text-slate-200"
                  onclick={handleCopyJson}
                  title={t("theme_editor.copy_json_title")}
                >
                  {copyFeedback ? t("theme_editor.btn_copied") : t("theme_editor.btn_copy_json")}
                </Button>
              </div>
            </div>

            <!-- Importar Tema de Arquivo -->
            <div class="bg-slate-800/50 p-3 rounded-xl border border-slate-700/50 space-y-2">
              <span class="text-[11px] font-bold text-white block">{t("theme_editor.import_theme_title")}</span>
              <label class="w-full py-2.5 px-3 bg-indigo-600 hover:bg-indigo-500 active:bg-indigo-700 text-white font-semibold rounded-xl cursor-pointer text-xs flex items-center justify-center gap-2 transition-all shadow-sm">
                <span>{t("theme_editor.btn_select_json")}</span>
                <input type="file" accept=".json" onchange={handleFileImport} class="hidden" />
              </label>
              {#if importFeedback}
                <p class="text-emerald-400 text-[11px] font-medium text-center">{importFeedback}</p>
              {/if}
            </div>
          </div>
        {/if}
      </div>

      <!-- Footer da Edição: RESTAURAR PADRÃO | CANCELAR | APLICAR | SALVAR -->
      <div class="px-4 py-2.5 border-t border-slate-800 bg-slate-950/60 flex items-center justify-between text-xs shrink-0">
        <Button
          size="sm"
          class="bg-transparent hover:bg-rose-500/10 text-slate-400 hover:text-rose-400 text-[11px] shadow-none"
          onclick={() => {
            if (confirm(t("theme_editor.reset_default_confirm"))) {
              themeStore.resetDraftToDefault();
            }
          }}
        >
          {t("theme_editor.btn_reset_default")}
        </Button>
        <div class="flex items-center gap-1.5">
          <Button
            size="sm"
            class="bg-slate-800 hover:bg-slate-700 text-slate-300"
            onclick={() => themeStore.cancelDraft()}
          >
            {t("theme_editor.btn_cancel")}
          </Button>
          <Button
            size="sm"
            class="bg-sky-600 hover:bg-sky-500 active:bg-sky-700 text-white font-semibold"
            onclick={() => themeStore.saveDraftToDisk(false)}
          >
            {t("theme_editor.btn_apply")}
          </Button>
          <Button
            size="sm"
            class="bg-emerald-600 hover:bg-emerald-500 active:bg-emerald-700 text-white font-bold shadow-lg shadow-emerald-900/30"
            onclick={() => themeStore.saveDraftToDisk(true)}
          >
            {t("theme_editor.btn_save")}
          </Button>
        </div>
      </div>
    {/if}
  </div>
{/if}
