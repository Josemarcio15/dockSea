<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { Button } from "$shared/components/buttons";

  let {
    server,
    usage,
    isActive,
    isLoading,
    onRefresh,
    onActivate,
    onViewContainers,
  }: {
    server: any;
    usage?: any;
    isActive: boolean;
    isLoading: boolean;
    onRefresh: () => void;
    onActivate: () => void;
    onViewContainers: () => void;
  } = $props();

  function formatBytes(bytes: number, decimals = 1) {
    if (!bytes || bytes === 0) return "0 B";
    const k = 1024;
    const sizes = ["B", "KB", "MB", "GB", "TB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return (
      parseFloat((bytes / Math.pow(k, i)).toFixed(Math.max(0, decimals))) +
      " " +
      sizes[i]
    );
  }
</script>

<div
  class="relative rounded-[22px] bg-gradient-to-b from-slate-900/90 to-[#090d16]/95 dark:from-[#111726]/90 dark:to-[#070b13]/95 border border-white/10 dark:border-white/5 backdrop-blur-xl shadow-[0_12px_32px_rgba(0,0,0,0.45)] hover:border-violet-500/40 hover:shadow-[0_16px_36px_rgba(99,102,241,0.15)] transition-all duration-300 flex flex-col p-5 gap-4 group text-slate-100 overflow-hidden w-full"
>
  <!-- Glow highlight top -->
  <div class="absolute -top-12 -left-12 w-32 h-32 bg-violet-600/15 rounded-full blur-2xl pointer-events-none"></div>

  <!-- Server Main Header Pill -->
  <div
    class="relative flex items-center justify-between p-3.5 rounded-2xl bg-white/[0.04] dark:bg-white/[0.03] border border-white/[0.08] backdrop-blur-md shadow-inner"
  >
    <div class="flex items-center gap-3.5 min-w-0 flex-1">
      <!-- Server Icon Container -->
      <div
        class="w-11 h-11 rounded-xl bg-violet-600/10 border border-violet-500/20 flex items-center justify-center text-xl shrink-0"
      >
        <span>🖥️</span>
      </div>

      <div class="flex flex-col min-w-0 flex-1">
        <h3
          class="font-extrabold text-base text-white tracking-tight truncate drop-shadow-sm"
        >
          {server.name}
        </h3>
        <p
          class="text-xs text-indigo-300 font-mono mt-0.5 truncate font-medium"
        >
          {server.connectionType === "ssh"
            ? `${server.username}@${server.host}:${server.port}`
            : "Local Docker Engine"}
        </p>
      </div>
    </div>

    <div class="flex items-center gap-2">
      {#if isActive}
        <span
          class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-semibold bg-emerald-500/15 border border-emerald-400/30 text-emerald-300 shadow-xs"
        >
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 shadow-[0_0_6px_#34d399] animate-pulse"></span>
          {t("profiles.active_badge")}
        </span>
      {:else}
        <span
          class="px-2.5 py-1 rounded-full text-[11px] font-medium bg-slate-500/15 text-slate-400 border border-slate-400/30"
        >
          {t("profiles.inactive_badge")}
        </span>
      {/if}

      {#if isActive}
        <Button
          size="xs"
          title="Atualizar estatísticas de hardware"
          loading={isLoading}
          onclick={onRefresh}
        >
          <span aria-hidden="true">↻</span>
        </Button>
      {/if}
    </div>
  </div>

  {#if isLoading}
    <div
      class="p-4 rounded-2xl bg-white/[0.03] border border-white/[0.06] space-y-2 animate-pulse"
    >
      <div class="h-2.5 bg-white/10 rounded-full w-3/4"></div>
      <div class="h-2 bg-white/10 rounded-full w-1/2"></div>
    </div>
  {:else if usage}
    <!-- Metrics Container -->
    <div
      class="p-4 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm divide-y divide-white/[0.06] space-y-3.5"
    >
      {#each [{ label: "RAM", used: usage.memUsed, total: usage.memTotal, percent: usage.memUsagePerc, color: usage.memUsagePerc > 85 ? "bg-rose-500" : usage.memUsagePerc > 65 ? "bg-amber-500" : "bg-emerald-500", badge: "bg-purple-500/15 text-purple-300 border border-purple-500/30" }, { label: "Disco (/)", used: usage.diskUsed, total: usage.diskTotal, percent: usage.diskUsagePerc, color: usage.diskUsagePerc > 85 ? "bg-rose-500" : usage.diskUsagePerc > 70 ? "bg-amber-500" : "bg-sky-500", badge: "bg-sky-500/15 text-sky-300 border border-sky-500/30" }] as metric, i}
        <div class="space-y-2 {i > 0 ? 'pt-3.5' : ''}">
          <div class="flex justify-between items-center text-xs font-semibold">
            <span
              class="px-2 py-0.5 rounded-md font-bold text-[10px] uppercase tracking-wider {metric.badge}"
              >{metric.label}</span
            >
            <span
              class="text-white font-mono font-bold"
              >{formatBytes(metric.used)} / {formatBytes(metric.total)}
              <span class="text-indigo-300"
                >({metric.percent.toFixed(1)}%)</span
              ></span
            >
          </div>
          <div
            class="h-2 w-full bg-black/40 rounded-full overflow-hidden p-0.5 border border-white/10"
          >
            <div
              class="h-full rounded-full transition-all duration-500 {metric.color} shadow-xs"
              style="width: {Math.min(100, Math.max(0, metric.percent))}%"
            ></div>
          </div>
        </div>
      {/each}

      {#if usage.swapTotal > 0}
        <div class="pt-3.5 space-y-2">
          <div class="flex justify-between items-center text-xs font-semibold">
            <span
              class="px-2 py-0.5 rounded-md font-bold text-[10px] uppercase tracking-wider bg-indigo-500/15 text-indigo-300 border border-indigo-500/30"
              >Swap</span
            >
            <span
              class="text-white font-mono font-bold"
              >{formatBytes(usage.swapUsed)} / {formatBytes(usage.swapTotal)}
              <span class="text-indigo-300"
                >({usage.swapUsagePerc.toFixed(1)}%)</span
              ></span
            >
          </div>
          <div
            class="h-2 w-full bg-black/40 rounded-full overflow-hidden p-0.5 border border-white/10"
          >
            <div
              class="h-full rounded-full transition-all duration-500 {usage.swapUsagePerc >
              70
                ? 'bg-rose-500'
                : 'bg-indigo-500'} shadow-xs"
              style="width: {Math.min(
                100,
                Math.max(0, usage.swapUsagePerc),
              )}%"
            ></div>
          </div>
        </div>
      {/if}

      {#if usage.uptime}
        <div
          class="pt-3 flex items-center justify-between text-xs text-slate-300 font-medium"
        >
          <span class="truncate max-w-35" title={usage.uptime}
            >{usage.uptime}</span
          >{#if usage.cpuCount}<span
              class="font-mono font-bold text-indigo-300"
              >{usage.cpuCount} CPU(s)</span
            >{/if}
        </div>
      {/if}
    </div>
  {/if}

  <div
    class="flex items-center gap-2.5 border-t border-white/[0.08] pt-3"
  >
    {#if isActive}
      <Button size="sm" themeKey="servers.view_containers_btn" class="w-full flex items-center justify-center gap-2" onclick={onViewContainers}>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="w-4 h-4"
        >
          <path
            d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"
          />
          <polyline points="3.27 6.96 12 12.01 20.73 6.96" />
          <line x1="12" y1="22.08" x2="12" y2="12" />
        </svg>
        <span>{t("devices.view_containers")}</span>
      </Button>
    {:else}
      <Button
        themeKey="servers.connect_btn"
        class="w-full flex items-center justify-center gap-2"
        onclick={onActivate}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2.2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="w-4 h-4"
        >
          <path d="M18.36 6.64a9 9 0 1 1-12.73 0" />
          <line x1="12" y1="2" x2="12" y2="12" />
        </svg>
        <span>{t("devices.activate")}</span>
      </Button>
    {/if}
  </div>
</div>
