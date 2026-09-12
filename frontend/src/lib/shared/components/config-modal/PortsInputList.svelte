<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { Button } from "$shared/components/buttons";
  import type { PortMapping } from "./types";

  let {
    ports = $bindable([]),
    onModified = () => {},
  }: {
    ports: PortMapping[];
    onModified: () => void;
  } = $props();

  function addPort() {
    ports = [...ports, { external: "", internal: "" }];
    onModified();
  }

  function removePort(idx: number) {
    ports = ports.filter((_, i) => i !== idx);
    onModified();
  }
</script>

<div class="flex flex-col gap-2">
  <div class="flex items-center justify-between">
    <span
      class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
      >{t("images.config_ports")}</span
    >
    <Button variant="primary" size="xs" onclick={addPort}>
      {t("images.config_add_port")}
    </Button>
  </div>
  {#each ports as port, i}
    <div class="flex gap-2 items-center">
      <input
        type="text"
        class="flex-1 px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0e2536] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
        placeholder={t("images.config_placeholder_port_ext")}
        bind:value={port.external}
        oninput={onModified}
      />
      <span class="text-slate-400">:</span>
      <input
        type="text"
        class="flex-1 px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0e2536] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
        placeholder={t("images.config_placeholder_port_int")}
        bind:value={port.internal}
        oninput={onModified}
      />
      <Button variant="danger" size="xs" onclick={() => removePort(i)}>
        ✕
      </Button>
    </div>
  {/each}
</div>
