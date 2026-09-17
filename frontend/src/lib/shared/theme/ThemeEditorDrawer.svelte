<script lang="ts">
  import { themeStore } from "./theme.svelte.ts";
  import { defaultTheme, type DockSeaTheme } from "./defaultTheme";
  import { navigation } from "$navigation/navigation.svelte";

  let activeTab = $state<"route" | "global" | "disk" | "json">("route");
  let jsonInput = $state("");
  let copyFeedback = $state(false);
  let importError = $state("");
  let themeNameInput = $state("");
  let diskActionFeedback = $state("");

  let currentRoute = $derived(navigation.currentRoute || "stacks");

  const routeTitles: Record<string, string> = {
    servers: "Dispositivos / Servidores",
    containers: "Containers",
    images: "Imagens",
    volumes: "Volumes",
    networks: "Redes",
    stacks: "Stacks (Compose)",
    builder: "Builder",
    extras: "Extras",
    profiles: "Perfis",
    config: "Configurações",
  };

  const currentRouteButtons = $derived.by(() => {
    const routeObj = themeStore.currentTheme.routes?.[currentRoute] || defaultTheme.routes?.[currentRoute] || {};
    return Object.entries(routeObj).map(([key, style]) => ({
      key,
      label: key.replace(/_/g, " "),
      bg: style.bg || (style as any).color || "#2563eb",
      hover: style.hover || style.bg || "#1d4ed8",
      text: style.text || "#ffffff",
      size: style.size || "sm",
    }));
  });

  function handleBgChange(route: string, btnKey: string, color: string) {
    themeStore.setRouteButtonColor(route, btnKey, color);
  }

  function handleHoverChange(route: string, btnKey: string, color: string) {
    themeStore.setRouteButtonHover(route, btnKey, color);
  }

  function handleTextChange(route: string, btnKey: string, color: string) {
    themeStore.setRouteButtonText(route, btnKey, color);
  }

  function handleSizeChange(route: string, btnKey: string, size: any) {
    themeStore.setRouteButtonSize(route, btnKey, size);
  }

  function handleCopyJson() {
    navigator.clipboard.writeText(themeStore.exportThemeJson());
    copyFeedback = true;
    setTimeout(() => (copyFeedback = false), 2000);
  }

  function handleImport() {
    importError = "";
    if (!jsonInput.trim()) return;
    const ok = themeStore.importTheme(jsonInput);
    if (!ok) {
      importError = "JSON inválido! Verifique a formatação.";
    } else {
      jsonInput = "";
      alert("Tema importado com sucesso!");
    }
  }

  function handleFileImport(e: Event) {
    const input = e.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const reader = new FileReader();
      reader.onload = (event) => {
        const content = event.target?.result as string;
        const ok = themeStore.importTheme(content);
        if (ok) {
          alert("Tema carregado com sucesso!");
        } else {
          alert("Arquivo JSON inválido.");
        }
      };
      reader.readAsText(input.files[0]);
    }
  }

  async function handleSaveToDisk() {
    if (!themeNameInput.trim()) {
      alert("Por favor, digite um nome para o arquivo de tema.");
      return;
    }
    const ok = await themeStore.saveThemeToDisk(themeNameInput.trim());
    if (ok) {
      diskActionFeedback = `Tema "${themeNameInput.trim()}.json" salvo em ~/Documents/DockSea/themes/`;
      themeNameInput = "";
      setTimeout(() => (diskActionFeedback = ""), 3000);
    } else {
      alert("Erro ao salvar tema no disco.");
    }
  }

  async function handleLoadFromDisk(name: string) {
    const ok = await themeStore.loadThemeFromDisk(name);
    if (ok) {
      diskActionFeedback = `Tema "${name}" carregado!`;
      setTimeout(() => (diskActionFeedback = ""), 3000);
    } else {
      alert(`Erro ao carregar tema ${name}.`);
    }
  }

  async function handleDeleteFromDisk(name: string) {
    if (confirm(`Excluir tema "${name}.json" do disco?`)) {
      await themeStore.deleteThemeFromDisk(name);
    }
  }
</script>

{#if themeStore.isEditorOpen}
  <div
    class="fixed bottom-5 right-5 z-50 w-96 max-h-[85vh] bg-slate-900/95 border border-slate-700/80 backdrop-blur-xl rounded-2xl shadow-2xl flex flex-col overflow-hidden text-slate-200 font-sans animate-in fade-in zoom-in-95 duration-200"
  >
    <!-- Header -->
    <div class="px-5 py-4 border-b border-slate-800 flex items-center justify-between bg-slate-950/40">
      <div class="flex items-center gap-2.5">
        <span class="text-xl">🎨</span>
        <div>
          <h3 class="text-sm font-bold text-white tracking-wide">Editor de Tema</h3>
          <p class="text-[11px] text-slate-400">~/Documents/DockSea/themes</p>
        </div>
      </div>
      <button
        onclick={() => themeStore.toggleEditor()}
        class="w-7 h-7 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white flex items-center justify-center transition-colors"
        title="Fechar Editor"
      >
        ✕
      </button>
    </div>

    <!-- Navegação de Abas -->
    <div class="flex border-b border-slate-800 bg-slate-950/20 text-xs font-semibold px-2 pt-2 gap-1 overflow-x-auto">
      <button
        onclick={() => (activeTab = "route")}
        class="px-2.5 py-2 rounded-t-lg border-b-2 whitespace-nowrap transition-all {activeTab === 'route'
          ? 'border-sky-500 text-sky-400 bg-slate-800/40'
          : 'border-transparent text-slate-400 hover:text-slate-200'}"
      >
        📍 Tela Atual
      </button>
      <button
        onclick={() => (activeTab = "global")}
        class="px-2.5 py-2 rounded-t-lg border-b-2 whitespace-nowrap transition-all {activeTab === 'global'
          ? 'border-sky-500 text-sky-400 bg-slate-800/40'
          : 'border-transparent text-slate-400 hover:text-slate-200'}"
      >
        🌐 Fundo
      </button>
      <button
        onclick={() => (activeTab = "disk")}
        class="px-2.5 py-2 rounded-t-lg border-b-2 whitespace-nowrap transition-all {activeTab === 'disk'
          ? 'border-sky-500 text-sky-400 bg-slate-800/40'
          : 'border-transparent text-slate-400 hover:text-slate-200'}"
      >
        💾 Temas Salvos
      </button>
      <button
        onclick={() => (activeTab = "json")}
        class="px-2.5 py-2 rounded-t-lg border-b-2 whitespace-nowrap transition-all {activeTab === 'json'
          ? 'border-sky-500 text-sky-400 bg-slate-800/40'
          : 'border-transparent text-slate-400 hover:text-slate-200'}"
      >
        {`{ } JSON`}
      </button>
    </div>

    <!-- Conteúdo da Aba -->
    <div class="p-5 overflow-y-auto max-h-[50vh] text-xs space-y-4">
      {#if activeTab === "route"}
        <div>
          <div class="flex items-center justify-between mb-3">
            <span class="font-bold text-slate-200 uppercase tracking-wider text-[11px]">
              {routeTitles[currentRoute] || currentRoute}
            </span>
            <span class="text-[10px] bg-slate-800 px-2 py-0.5 rounded text-slate-400">
              {currentRouteButtons.length} botões
            </span>
          </div>

          {#if currentRouteButtons.length > 0}
            <div class="space-y-3">
              <p class="text-slate-400 text-[11px]">
                Ajuste a cor e tamanho de cada botão desta tela:
              </p>
              {#each currentRouteButtons as btn}
                <div class="flex flex-col bg-slate-800/50 p-3 rounded-xl border border-slate-700/50 gap-2.5">
                  <div class="flex items-center justify-between">
                    <span class="font-bold text-white capitalize text-xs truncate">{btn.label}</span>
                    <!-- Size Selector -->
                    <select
                      value={btn.size}
                      onchange={(e) => handleSizeChange(currentRoute, btn.key, e.currentTarget.value)}
                      class="bg-slate-900 border border-slate-700 rounded-lg px-2 py-0.5 text-[10px] text-slate-300 focus:outline-none"
                    >
                      <option value="xs">Tam: XS</option>
                      <option value="sm">Tam: SM</option>
                      <option value="md">Tam: MD</option>
                      <option value="lg">Tam: LG</option>
                    </select>
                  </div>

                  <div class="grid grid-cols-3 gap-2 pt-1 border-t border-slate-700/40 text-[10px]">
                    <!-- Fundo Normal -->
                    <label class="flex flex-col items-center gap-1 bg-slate-900/60 p-1.5 rounded-lg border border-slate-800 cursor-pointer">
                      <span class="text-slate-400">Fundo</span>
                      <input
                        type="color"
                        value={btn.bg}
                        oninput={(e) => handleBgChange(currentRoute, btn.key, e.currentTarget.value)}
                        class="w-6 h-6 rounded cursor-pointer bg-transparent border-0"
                      />
                    </label>

                    <!-- Hover Fundo -->
                    <label class="flex flex-col items-center gap-1 bg-slate-900/60 p-1.5 rounded-lg border border-slate-800 cursor-pointer">
                      <span class="text-slate-400">Hover</span>
                      <input
                        type="color"
                        value={btn.hover}
                        oninput={(e) => handleHoverChange(currentRoute, btn.key, e.currentTarget.value)}
                        class="w-6 h-6 rounded cursor-pointer bg-transparent border-0"
                      />
                    </label>

                    <!-- Cor do Texto -->
                    <label class="flex flex-col items-center gap-1 bg-slate-900/60 p-1.5 rounded-lg border border-slate-800 cursor-pointer">
                      <span class="text-slate-400">Texto</span>
                      <input
                        type="color"
                        value={btn.text}
                        oninput={(e) => handleTextChange(currentRoute, btn.key, e.currentTarget.value)}
                        class="w-6 h-6 rounded cursor-pointer bg-transparent border-0"
                      />
                    </label>
                  </div>
                </div>
              {/each}
            </div>
          {:else}
            <div class="py-6 text-center text-slate-500">
              Nenhum botão registrado para esta rota ainda.
            </div>
          {/if}
        </div>
      {:else if activeTab === "global"}
        <div class="space-y-3">
          <p class="text-slate-400 text-[11px]">Cores de superfície e layout:</p>
          
          <div class="flex items-center justify-between bg-slate-800/40 p-2.5 rounded-xl border border-slate-700/40">
            <span>Fundo Geral</span>
            <input
              type="color"
              value={themeStore.currentTheme.global.appBg}
              oninput={(e) => themeStore.setGlobalColor("appBg", e.currentTarget.value)}
              class="w-8 h-8 rounded-lg cursor-pointer bg-transparent border-0"
            />
          </div>

          <div class="flex items-center justify-between bg-slate-800/40 p-2.5 rounded-xl border border-slate-700/40">
            <span>Fundo dos Cards</span>
            <input
              type="color"
              value={themeStore.currentTheme.global.cardBg}
              oninput={(e) => themeStore.setGlobalColor("cardBg", e.currentTarget.value)}
              class="w-8 h-8 rounded-lg cursor-pointer bg-transparent border-0"
            />
          </div>

          <div class="flex items-center justify-between bg-slate-800/40 p-2.5 rounded-xl border border-slate-700/40">
            <span>Sidebar (Topo)</span>
            <input
              type="color"
              value={themeStore.currentTheme.global.sidebarFrom}
              oninput={(e) => themeStore.setGlobalColor("sidebarFrom", e.currentTarget.value)}
              class="w-8 h-8 rounded-lg cursor-pointer bg-transparent border-0"
            />
          </div>

          <div class="flex items-center justify-between bg-slate-800/40 p-2.5 rounded-xl border border-slate-700/40">
            <span>Sidebar (Fim)</span>
            <input
              type="color"
              value={themeStore.currentTheme.global.sidebarTo}
              oninput={(e) => themeStore.setGlobalColor("sidebarTo", e.currentTarget.value)}
              class="w-8 h-8 rounded-lg cursor-pointer bg-transparent border-0"
            />
          </div>
        </div>
      {:else if activeTab === "disk"}
        <div class="space-y-4">
          <div>
            <span class="text-[11px] text-slate-400 block mb-1.5 font-medium">
              Salvar tema atual no disco:
            </span>
            <div class="flex gap-2">
              <input
                type="text"
                bind:value={themeNameInput}
                placeholder="Ex: meu_tema_escuro"
                class="flex-1 bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-sky-500"
              />
              <button
                onclick={handleSaveToDisk}
                class="px-3 py-2 bg-emerald-600 hover:bg-emerald-500 active:bg-emerald-700 text-white font-semibold rounded-xl transition-all"
              >
                Salvar .json
              </button>
            </div>
            {#if diskActionFeedback}
              <p class="text-emerald-400 text-[11px] mt-1.5">{diskActionFeedback}</p>
            {/if}
          </div>

          <div class="pt-2 border-t border-slate-800">
            <div class="flex items-center justify-between mb-2">
              <span class="text-[11px] text-slate-400 font-medium">Temas em ~/Documents/DockSea/themes:</span>
              <button
                onclick={() => themeStore.refreshDiskThemes()}
                class="text-[10px] text-sky-400 hover:underline"
              >
                Atualizar
              </button>
            </div>

            {#if themeStore.availableDiskThemes.length === 0}
              <div class="py-4 text-center text-slate-500 text-[11px] bg-slate-950/40 rounded-xl border border-slate-800/60">
                Nenhum tema salvo ainda na pasta themes.
              </div>
            {:else}
              <div class="space-y-1.5">
                {#each themeStore.availableDiskThemes as themeName}
                  <div class="flex items-center justify-between bg-slate-800/60 px-3 py-2 rounded-xl border border-slate-700/50">
                    <span class="font-medium text-slate-200 truncate">{themeName}.json</span>
                    <div class="flex items-center gap-1.5">
                      <button
                        onclick={() => handleLoadFromDisk(themeName)}
                        class="px-2.5 py-1 bg-sky-600/80 hover:bg-sky-500 text-white rounded-lg text-[10px] font-semibold transition-all"
                      >
                        Carregar
                      </button>
                      <button
                        onclick={() => handleDeleteFromDisk(themeName)}
                        class="p-1 hover:bg-rose-500/20 text-rose-400 rounded-lg transition-colors"
                        title="Excluir"
                      >
                        ✕
                      </button>
                    </div>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        </div>
      {:else if activeTab === "json"}
        <div class="space-y-3">
          <p class="text-slate-400 text-[11px]">
            Exporte para compartilhar com a comunidade ou cole um tema `.json`:
          </p>

          <div class="flex gap-2">
            <button
              onclick={handleCopyJson}
              class="flex-1 py-2 px-3 bg-sky-600 hover:bg-sky-500 active:bg-sky-700 text-white font-semibold rounded-xl transition-all text-center"
            >
              {copyFeedback ? "✓ Copiado!" : "📋 Copiar JSON do Tema"}
            </button>
            <label class="py-2 px-3 bg-slate-800 hover:bg-slate-700 text-slate-200 font-semibold rounded-xl cursor-pointer transition-all">
              📁 Arquivo
              <input type="file" accept=".json" onchange={handleFileImport} class="hidden" />
            </label>
          </div>

          <div class="space-y-1.5 pt-2">
            <span class="text-[11px] text-slate-400">Importar colando JSON:</span>
            <textarea
              bind:value={jsonInput}
              rows="5"
              placeholder="Cole seu tema .json aqui..."
              class="w-full bg-slate-950 border border-slate-800 rounded-xl p-2.5 font-mono text-[11px] text-slate-300 focus:outline-none focus:border-sky-500"
            ></textarea>
            {#if importError}
              <p class="text-rose-400 text-[11px]">{importError}</p>
            {/if}
            <button
              onclick={handleImport}
              class="w-full py-2 bg-indigo-600 hover:bg-indigo-500 active:bg-indigo-700 text-white font-semibold rounded-xl transition-all"
            >
              Aplicar Tema Colado
            </button>
          </div>
        </div>
      {/if}
    </div>

    <!-- Footer Actions -->
    <div class="px-5 py-3 border-t border-slate-800 bg-slate-950/40 flex items-center justify-between text-xs">
      <button
        onclick={() => {
          if (confirm("Deseja restaurar todas as cores para o tema padrão?")) {
            themeStore.resetToDefault();
          }
        }}
        class="text-slate-400 hover:text-rose-400 transition-colors"
      >
        Restaurar Padrão
      </button>
      <button
        onclick={() => themeStore.toggleEditor()}
        class="px-4 py-1.5 bg-slate-800 hover:bg-slate-700 text-white rounded-lg font-medium transition-all"
      >
        Concluir
      </button>
    </div>
  </div>
{/if}
