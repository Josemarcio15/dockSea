<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { Button } from "$shared/components/buttons";
  import type { VolumeMapping } from "./types";
  import { getDefaultContainerPath } from "./types";

  let {
    image = "",
    volumes = $bindable([]),
    existingVolumes = [],
    hasEmptyVolume = false,
    onModified = () => {},
  }: {
    image: string;
    volumes: VolumeMapping[];
    existingVolumes: string[];
    hasEmptyVolume: boolean;
    onModified: () => void;
  } = $props();

  function addVolume() {
    volumes = [...volumes, { host: "", container: "" }];
    onModified();
  }

  function removeVolume(idx: number) {
    volumes = volumes.filter((_, i) => i !== idx);
    onModified();
  }
</script>

<div class="flex flex-col gap-2">
  <div class="flex items-center justify-between">
    <span
      class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
      >{t("images.config_volumes")}</span
    >
    <Button size="xs" onclick={addVolume}>
      {t("images.config_add_volume")}
    </Button>
  </div>
  {#each volumes as vol, i}
    <div class="flex gap-2 items-center">
      <input
        type="text"
        class="flex-1 px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
        placeholder={t("images.config_placeholder_vol_host")}
        bind:value={vol.host}
        oninput={onModified}
        list="volumes-datalist"
      />
      <span class="text-slate-400">:</span>
      <input
        type="text"
        class="flex-1 px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
        placeholder={vol.host
          ? getDefaultContainerPath(image)
          : t("images.config_placeholder_vol_container")}
        bind:value={vol.container}
        oninput={onModified}
      />
      <Button size="xs" onclick={() => removeVolume(i)}>
        ✕
      </Button>
    </div>
  {/each}
  <datalist id="volumes-datalist">
    {#each existingVolumes as vol}
      <option value={vol}></option>
    {/each}
  </datalist>
  {#if hasEmptyVolume}
    <div class="text-xs text-red-500 font-semibold mt-1">
      {t("images.config_empty_volume_warning")}
    </div>
  {/if}
</div>
