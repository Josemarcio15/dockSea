<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import PageTitle from "$shared/components/PageTitle.svelte";

  import StatusBanner from "$shared/components/StatusBanner.svelte";
  import { Button } from "$shared/components/buttons";
  import ServerCard from "./components/ServerCard.svelte";
  import { createServersStore } from "./store.svelte";
  import type { Route } from "$navigation/navigation.types";

  let { data, navigate }: { data: any; navigate?: (route: Route) => void } =
    $props();
  const serverStore = createServersStore(() => data);
  $effect(() => {
    void serverStore.load();
  });
</script>

<div class="space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <PageTitle title={t("sidebar.devices")} />
    <Button variant="primary" size="sm" onclick={() => navigate?.("config")}
      >{t("devices.manage_vps")}</Button>
  </div>

  <StatusBanner />

  {#if data.servers.length === 0}
    <div
      class="flex flex-col items-center justify-center py-20 px-6 text-center border border-dashed border-slate-200 dark:border-slate-800 rounded-3xl bg-white/80 dark:bg-[#0b1d2e]/80 backdrop-blur-sm shadow-sm max-w-2xl mx-auto"
    >
      <h2 class="text-xl font-bold text-slate-800 dark:text-white mb-2">
        {t("devices.empty_title")}
      </h2>
      <p class="text-sm text-slate-500 dark:text-slate-400 mb-6 max-w-md">
        {t("devices.empty_desc")}
      </p>
      <Button variant="primary" onclick={() => navigate?.("config")}
        >{t("devices.add_first")}</Button>
    </div>
  {:else}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each data.servers as server (server.id)}
        <ServerCard
          {server}
          usage={serverStore.usageCache[server.id]}
          isActive={data.activeVps?.id === server.id}
          isLoading={serverStore.loadingUsage[server.id]}
          onRefresh={() => serverStore.fetchUsage(server, true)}
          onActivate={() => serverStore.activate(server)}
          onViewContainers={() => navigate?.("containers")}
        />
      {/each}
    </div>
  {/if}
</div>
