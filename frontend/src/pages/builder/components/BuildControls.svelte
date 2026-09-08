<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { Button } from "$shared/components/buttons";
  import type { BuilderStore } from "../types";
  let { store, goToImages }: { store: BuilderStore; goToImages: () => void } =
    $props();
</script>

{#if store.hasDockerfile}
  <div class="flex flex-wrap items-center gap-2 mt-4">
    <div
      class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-slate-50 dark:bg-slate-800/80 border border-slate-200 dark:border-slate-700/80 focus-within:border-violet-500 transition-colors w-full"
    >
      <span
        class="text-xs font-semibold text-slate-500 select-none whitespace-nowrap"
        >🏷️ {t("builder.tag_label")}</span
      ><input
        type="text"
        class="bg-transparent text-xs font-mono font-bold text-violet-600 focus:outline-none flex-1 placeholder:text-slate-400"
        placeholder={store.defaultTag}
        bind:value={store.customTag}
        disabled={store.status === "building"}
      /><span class="text-xs font-mono text-slate-400 select-none">:latest</span
      >{#if store.customTag.trim()}<Button variant="warning"
          size="xs"
          onclick={() => (store.customTag = "")}>Restaurar</Button>{/if}
    </div>
    {#if store.status === "building"}<span
        class="text-xs text-amber-500 font-semibold"
        >{t("builder.building")}</span
      >{/if}
    <div class="flex items-center gap-2 w-full">
      <Button variant="success"
        size="md"
        class="flex-1"
        disabled={!store.canBuild}
        loading={store.status === "building"}
        onclick={() => store.build()}>{t("builder.build_btn")}</Button>{#if store.currentPath}<Button variant="primary"
          size="md"
          disabled={store.savedPaths.includes(store.currentPath)}
          onclick={() => store.saveCurrentPath()}
          >{t("builder.save_path")}</Button>{/if}
    </div>
    {#if store.status === "success"}<div
        class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 w-full"
      >
        <p class="text-xs font-bold">
          {t("builder.image_ready", { name: store.builtImage })}
        </p>
      </div>
      <Button variant="primary" size="md" class="w-full" onclick={goToImages}
        >{t("builder.transfer_btn")}</Button>{/if}
  </div>
{:else if store.currentPath && !store.loading}<span
    class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-bold bg-slate-100 text-slate-500 border border-slate-200 mt-4"
    >{t("builder.no_dockerfile")}</span
  >{/if}
