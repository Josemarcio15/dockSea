<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import {
    notifySuccess,
    notifyError,
  } from "$shared/stores/notification.svelte";
  import {
    ButtonGreen,
    ButtonYellow,
    ButtonPurple,
  } from "$shared/components/buttons";
  import { folderNameFromPath } from "../service";
  import type { BuilderFolder, BuilderStore } from "../types";

  let { store }: { store: BuilderStore } = $props();
  let editablePath = $state("");
  let lastSyncedPath = $state("");

  $effect(() => {
    if (store.currentPath && store.currentPath !== lastSyncedPath) {
      editablePath = store.currentPath;
      lastSyncedPath = store.currentPath;
    }
  });

  async function navigateToPath() {
    const path = editablePath.trim();
    if (path && path !== store.currentPath) await store.browse(path);
  }
  async function save() {
    try {
      await store.saveCurrentPath();
      notifySuccess("Pasta salva com sucesso!");
    } catch (error: any) {
      notifyError(`Erro ao salvar pasta: ${error?.message || error}`);
    }
  }
  async function remove(path: string, event: Event) {
    event.stopPropagation();
    try {
      await store.removeSavedPath(path);
      notifySuccess("Pasta removida dos favoritos!");
    } catch (error: any) {
      notifyError(`Erro ao remover pasta: ${error?.message || error}`);
    }
  }
</script>

<div class="space-y-4">
  <div
    class="relative rounded-[22px] bg-gradient-to-b from-slate-900/90 to-[#090d16]/95 dark:from-[#111726]/90 dark:to-[#070b13]/95 border border-white/10 dark:border-white/5 backdrop-blur-xl shadow-[0_12px_32px_rgba(0,0,0,0.45)] p-5 space-y-4 text-slate-100 overflow-hidden"
  >
    <!-- Glow highlight top -->
    <div class="absolute -top-12 -left-12 w-32 h-32 bg-violet-600/15 rounded-full blur-2xl pointer-events-none"></div>

    <!-- Header Section (Glassmorphism Pill) -->
    <div
      class="relative flex items-center justify-between p-3.5 rounded-2xl bg-white/[0.04] dark:bg-white/[0.03] border border-white/[0.08] backdrop-blur-md shadow-inner"
    >
      <div class="flex items-center gap-3.5 min-w-0 flex-1">
        <!-- 3D Folder Icon -->
        <div
          class="w-11 h-11 rounded-xl bg-violet-600/10 border border-violet-500/20 flex items-center justify-center text-xl shrink-0"
        >
          <span>📁</span>
        </div>

        <div class="flex flex-col min-w-0 flex-1">
          <h3
            class="font-extrabold text-base text-white tracking-tight truncate drop-shadow-sm m-0"
          >
            {t("builder.nav_title")}
          </h3>
          <span class="text-xs text-slate-400 font-mono truncate">
            {store.currentPath || "Navegador de diretórios"}
          </span>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <ButtonPurple size="xs" onclick={() => store.browse()}
          >{t("builder.nav_home")}</ButtonPurple
        >
        {#if store.parentPath}
          <ButtonYellow
            size="xs"
            onclick={() => store.browse(store.parentPath ?? "")}
            >{t("builder.nav_up")}</ButtonYellow
          >
        {/if}
      </div>
    </div>

    <!-- Saved shortcuts -->
    {#if store.savedPaths.length > 0}
      <div class="relative flex items-center gap-2 flex-wrap pt-2 border-t border-white/[0.08]">
        {#each store.savedPaths as path}
          <div class="relative inline-block group">
            <ButtonGreen size="xs" onclick={() => store.browse(path)} title={path}>
              <span class="mr-1">📌</span>{folderNameFromPath(path)}
            </ButtonGreen>
            <span
              class="absolute -top-1.5 -right-1.5 w-3.5 h-3.5 rounded-full bg-rose-500 text-white text-[8px] flex items-center justify-center opacity-0 group-hover:opacity-100 cursor-pointer z-10 shadow-xs"
              role="button"
              tabindex="0"
              onclick={(event) => remove(path, event)}
              onkeydown={(event) => {
                if (event.key === "Enter" || event.key === " ")
                  remove(path, event);
              }}
              title="Remover atalho">✕</span
            >
          </div>
        {/each}
      </div>
    {/if}

    <!-- Path input -->
    {#if store.currentPath}
      <input
        type="text"
        aria-label="Caminho da pasta do projeto"
        class="relative w-full px-3.5 py-2.5 rounded-xl bg-black/40 border border-white/10 text-xs font-mono text-slate-300 focus:outline-none focus:border-violet-500 transition-colors"
        bind:value={editablePath}
        title={store.currentPath}
        onkeydown={(event) => { if (event.key === "Enter") void navigateToPath(); }}
        onblur={() => void navigateToPath()}
      />
    {/if}

    <!-- Folders list -->
    <div
      class="relative border border-white/10 rounded-2xl max-h-64 overflow-y-auto bg-black/30 backdrop-blur-md divide-y divide-white/5"
    >
      {#if store.loading}
        <div class="flex items-center justify-center py-8">
          <span
            class="animate-spin inline-block w-5 h-5 border-2 border-violet-500 border-t-transparent rounded-full"
          ></span>
        </div>
      {:else if store.folders.length === 0 && !store.hasDockerfile}
        <div
          class="text-center py-8 text-xs text-slate-400 italic"
        >
          Nenhuma subpasta encontrada.
        </div>
      {:else}
        {#each store.folders as folder: BuilderFolder}
          <button
            type="button"
            class="w-full flex items-center gap-2.5 px-4 py-2.5 text-xs text-left hover:bg-white/[0.05] cursor-pointer transition-colors"
            onclick={() => store.browse(folder.path)}
          >
            <span class="text-base shrink-0">📂</span>
            <span
              class="font-mono font-semibold text-slate-200 truncate"
            >
              {folder.name}
            </span>
          </button>
        {/each}
      {/if}
    </div>

    <!-- Detection badges -->
    {#if store.hasDockerfile}
      <div class="relative flex flex-wrap items-center gap-2 pt-1">
        <span
          class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-semibold bg-emerald-500/15 text-emerald-300 border border-emerald-400/30 shadow-xs"
        >
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 shadow-[0_0_6px_#34d399] animate-pulse"></span>
          ✓ {t("builder.detect_dockerfile")}
        </span>
        {#if store.hasDockerignore}
          <span
            class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-semibold bg-blue-500/15 text-blue-300 border border-blue-400/30 shadow-xs"
          >
            ✓ {t("builder.detect_dockerignore")}
          </span>
        {:else}
          <span
            class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-medium bg-amber-500/15 text-amber-300 border border-amber-400/30"
          >
            ⚠️ {t("builder.warn_no_dockerignore")}
          </span>
        {/if}
      </div>
    {/if}

    {#if store.hasDockerignore && store.ignoredFiles.length > 0}
      <details class="relative rounded-2xl border border-blue-500/20 bg-blue-950/20 p-3">
        <summary class="cursor-pointer text-xs font-bold text-blue-300">Arquivos ignorados ({store.ignoredFiles.length})</summary>
        <ul class="mt-2 max-h-32 overflow-y-auto space-y-1 text-[11px] font-mono text-slate-300">{#each store.ignoredFiles as file}<li class="break-all">{file}</li>{/each}</ul>
      </details>
    {/if}
  </div>
</div>
