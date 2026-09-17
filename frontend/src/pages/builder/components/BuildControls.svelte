<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { Button } from "$shared/components/buttons";
  import type { BuilderStore } from "../types";
  let { store, goToImages }: { store: BuilderStore; goToImages: () => void } =
    $props();
</script>

{#if store.hasDockerfile}
  <div
    class="relative rounded-[22px] bg-gradient-to-b from-slate-900/90 to-[#090d16]/95 dark:from-[#111726]/90 dark:to-[#070b13]/95 border border-white/10 dark:border-white/5 backdrop-blur-xl shadow-[0_12px_32px_rgba(0,0,0,0.45)] p-5 flex flex-col gap-3.5 text-slate-100 overflow-hidden mt-4"
  >
    <!-- Glow highlight top -->
    <div class="absolute -top-12 -left-12 w-28 h-28 bg-violet-600/15 rounded-full blur-2xl pointer-events-none"></div>

    <div
      class="relative flex items-center gap-2 p-3 rounded-2xl bg-white/[0.04] border border-white/[0.08] focus-within:border-violet-400/50 transition-colors w-full"
    >
      <span
        class="text-xs font-semibold text-slate-400 select-none whitespace-nowrap"
      >
        🏷️ {t("builder.tag_label")}
      </span>
      <input
        type="text"
        class="bg-transparent text-xs font-mono font-bold text-violet-300 focus:outline-none flex-1 placeholder:text-slate-500"
        placeholder={store.defaultTag}
        bind:value={store.customTag}
        disabled={store.status === "building"}
      />
      <span class="text-xs font-mono text-slate-500 select-none">:latest</span>
      {#if store.customTag.trim()}
        <Button
          size="xs"
          onclick={() => (store.customTag = "")}
        >
          Restaurar
        </Button>
      {/if}
    </div>

    {#if store.status === "building"}
      <span class="text-xs text-amber-400 font-semibold flex items-center gap-2">
        <span class="animate-spin text-sm">↻</span>
        {t("builder.building")}
      </span>
    {/if}

    <div class="relative flex items-center gap-2 w-full pt-1">
      <Button
        size="md"
        class="flex-1"
        disabled={!store.canBuild}
        loading={store.status === "building"}
        onclick={() => store.build()}
      >
        {t("builder.build_btn")}
      </Button>
      {#if store.currentPath}
        <Button
          size="md"
          disabled={store.savedPaths.includes(store.currentPath)}
          onclick={() => store.saveCurrentPath()}
        >
          {t("builder.save_path")}
        </Button>
      {/if}
    </div>

    {#if store.status === "success"}
      <div
        class="relative p-3.5 rounded-2xl bg-emerald-500/15 border border-emerald-500/30 text-emerald-300 w-full shadow-xs"
      >
        <p class="text-xs font-extrabold m-0">
          {t("builder.image_ready", { name: store.builtImage })}
        </p>
      </div>
      <Button size="md" class="w-full" onclick={goToImages}>
        {t("builder.transfer_btn")}
      </Button>
    {/if}
  </div>
{:else if store.currentPath && !store.loading}
  <span
    class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-[11px] font-semibold bg-white/5 border border-white/10 text-slate-400 mt-4"
  >
    {t("builder.no_dockerfile")}
  </span>
{/if}
