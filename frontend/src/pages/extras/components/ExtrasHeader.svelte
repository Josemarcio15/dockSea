<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import PageTitle from "$shared/components/PageTitle.svelte";

  import { Button } from "$shared/components/buttons";
  import type { ExtrasMainTab } from "../types";

  let {
    mainTab = $bindable(),
    refresh,
  }: { mainTab: ExtrasMainTab; refresh: () => void } = $props();
</script>

<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
  <PageTitle title={t("extras.title")} />
  <div class="flex items-center gap-3">
    <div
      class="inline-flex rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-[#0b1d2e] p-1 shadow-xs"
    >
      <button
        type="button"
        class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer {mainTab ===
        'nginx'
          ? 'bg-violet-600 text-white shadow-xs'
          : 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}"
        onclick={() => (mainTab = "nginx")}>{t("extras.tab_nginx")}</button
      >
      <button
        type="button"
        class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer {mainTab ===
        'ports'
          ? 'bg-violet-600 text-white shadow-xs'
          : 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}"
        onclick={() => (mainTab = "ports")}>{t("extras.tab_ports")}</button
      >
      <button type="button" class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all cursor-pointer {mainTab === 'deploy_temp' ? 'bg-violet-600 text-white shadow-xs' : 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}" onclick={() => (mainTab = "deploy_temp")}>{t("extras.tab_deploy_temp")}</button>
    </div>
    <Button variant="neutral" size="sm" title={t("common.refresh")} onclick={refresh}
      >{t("common.refresh")}</Button>
  </div>
</div>
