<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import type { Container } from "$lib/domains/containers";
  import { statsState } from "$shared/stores/stats.svelte";
  import { ButtonPurple, ButtonCyan, ButtonPink } from "$shared/components/buttons";
  import FormModal from "$shared/components/FormModal.svelte";

  let {
    container,
    checked = false,
    on_toggle = () => {},
    on_open_logs = (name: string) => {},
  }: {
    container: Container;
    checked?: boolean;
    on_toggle?: () => void;
    on_open_logs?: (name: string) => void;
  } = $props();

  let expanded = $state(false);
  let showEnv = $state(false);
  let showLabels = $state(false);

  const labels = $derived(Object.entries(container.labels || {}));

  const myStats = $derived(
    statsState?.stats
      ? statsState.stats.find(
          (s) =>
            s &&
            s.ID &&
            (s.ID === container.id ||
              s.ID.startsWith(container.id.substring(0, 12))),
        )
      : undefined,
  );

  // Helpers
  const [accentBorder, statusColor] = $derived.by(() => {
    const s = container.status || "";
    if (s.includes("Up")) {
      return [
        "border-l-emerald-400",
        "text-emerald-600 dark:text-emerald-400",
      ];
    } else if (s.includes("Exited") || s.includes("Paused")) {
      return ["border-l-amber-400", "text-amber-600 dark:text-amber-400"];
    } else {
      return ["border-l-red-400", "text-red-600 dark:text-red-400"];
    }
  });

  const statusDotColor = $derived(
    (container.status || "").includes("Up")
      ? "bg-emerald-500"
      : (container.status || "").includes("Exited") ||
          (container.status || "").includes("Paused")
        ? "bg-amber-500"
        : "bg-red-500",
  );

  const pulseClass = $derived(
    (container.status || "").includes("Up") ? "animate-pulse" : "",
  );

  const statusLabel = $derived.by(() => {
    const rawStatus = (container.status || "").trim();
    if (rawStatus.includes("Up")) {
      const uptime = rawStatus.replace(/^Up\s*/i, "").trim();
      return uptime ? `Ativo (${uptime})` : "Ativo";
    }
    return rawStatus.includes("Exited") ? t("containers.status_stop") : rawStatus.includes("Paused") ? "Pausado" : "Indisponível";
  });

  const statusBg = $derived(
    (container.status || "").includes("Up")
      ? "bg-emerald-50 dark:bg-emerald-950/40 text-emerald-600 dark:text-emerald-400 border-emerald-200 dark:border-emerald-900/60"
      : (container.status || "").includes("Exited") ||
          (container.status || "").includes("Paused")
        ? "bg-amber-50 dark:bg-amber-950/40 text-amber-600 dark:text-amber-400 border-amber-200 dark:border-amber-900/60"
        : "bg-red-50 dark:bg-red-950/40 text-red-600 dark:text-red-400 border-red-200 dark:border-red-900/60",
  );

  const portItems = $derived.by(() => {
    if (!container.ports) return [];
    return container.ports
      .split(", ")
      .map((entry) => {
        const e = entry.trim();
        if (!e) return null;
        const formatted = e.replace("->", " → ");
        const isIpv6 = e.includes("::");
        return {
          formatted,
          tag: isIpv6 ? "IPv6" : "IPv4",
          class: isIpv6
            ? "bg-purple-50 dark:bg-purple-950/30 border-purple-200 dark:border-purple-900/40 text-purple-700 dark:text-purple-400"
            : "bg-blue-50 dark:bg-blue-950/30 border-blue-200 dark:border-blue-900/40 text-blue-700 dark:text-blue-400",
        };
      })
      .filter((x) => x !== null);
  });

  const createdStr = $derived.by(() => {
    if (!container.created) return "";
    const date = new Date(container.created * 1000);
    const day = String(date.getDate()).padStart(2, "0");
    const month = String(date.getMonth() + 1).padStart(2, "0");
    const year = date.getFullYear();
    return `${day}/${month}/${year}`;
  });

  const networkItems = $derived.by(() => {
    if (!container.networks) return [];
    return Object.entries(container.networks).map(([netName, netEndpoint]) => {
      return {
        name: netName,
        ip: netEndpoint?.ipAddress || "—",
        gateway: netEndpoint?.gateway || "",
      };
    });
  });

  const shortCid = $derived((container.id || "").substring(0, 12));

  const restartPolicyDisplay = $derived.by(() => {
    const policy =
      container.restartPolicy || (container as any).restart_policy;
    if (!policy) return "—";
    switch (policy.toLowerCase()) {
      case "always":
        return t("containers.card_restart_always");
      case "unless-stopped":
        return t("containers.card_restart_unless_stopped");
      case "on-failure":
        return t("containers.card_restart_on_failure");
      case "no":
        return t("containers.card_restart_no");
      default:
        return policy;
    }
  });

  const iconUrl = $derived.by(() => {
    const target = `${container.name} ${container.image || ""}`.toLowerCase();
    if (target.includes("postgres") || target.includes("postgre"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/postgresql/postgresql-original.svg";
    if (target.includes("mariadb"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mariadb/mariadb-original.svg";
    if (target.includes("mysql"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mysql/mysql-original.svg";
    if (target.includes("mongo"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mongodb/mongodb-original.svg";
    if (target.includes("redis"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/redis/redis-original.svg";
    if (target.includes("wordpress"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/wordpress/wordpress-plain.svg";
    if (target.includes("nginx"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nginx/nginx-original.svg";
    if (target.includes("node"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nodejs/nodejs-original.svg";
    if (target.includes("python"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/python/python-original.svg";
    if (target.includes("golang") || target.includes("go"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg";
    if (target.includes("debian"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/debian/debian-original.svg";
    if (target.includes("alpine"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/alpinejs/alpinejs-original.svg";
    if (target.includes("ubuntu"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/ubuntu/ubuntu-plain.svg";
    return null;
  });
</script>

<div
  class="relative rounded-[22px] bg-gradient-to-b from-slate-900/90 to-[#090d16]/95 dark:from-[#111726]/90 dark:to-[#070b13]/95 border border-white/10 dark:border-white/5 backdrop-blur-xl shadow-[0_12px_32px_rgba(0,0,0,0.45)] hover:border-violet-500/40 hover:shadow-[0_16px_36px_rgba(99,102,241,0.15)] transition-all duration-300 flex flex-col p-4 gap-3.5 group text-slate-100 overflow-hidden w-full"
>
  <!-- Glow highlight top -->
  <div class="absolute -top-12 -left-12 w-32 h-32 bg-violet-600/15 rounded-full blur-2xl pointer-events-none"></div>

  <!-- Header Section (Glassmorphism Pill) -->
  <div
    class="relative flex items-center justify-between p-3.5 rounded-2xl bg-white/[0.04] dark:bg-white/[0.03] border border-white/[0.08] backdrop-blur-md shadow-inner"
  >
    <div class="flex items-center gap-3.5 min-w-0 flex-1">
      <!-- Checkbox -->
      <button
        type="button"
        class="w-5 h-5 rounded-lg border flex items-center justify-center cursor-pointer transition-all duration-200 shrink-0 {checked
          ? 'bg-gradient-to-tr from-violet-600 to-indigo-500 border-violet-400 text-white shadow-[0_0_10px_rgba(124,58,237,0.5)]'
          : 'border-white/20 bg-white/5 hover:border-violet-400'}"
        onclick={on_toggle}
        aria-label="Selecionar container"
      >
        {#if checked}
          <span class="text-white text-[11px] font-extrabold leading-none">✓</span>
        {/if}
      </button>

      <!-- Technology Icon Container -->
      <div
        class="w-11 h-11 rounded-xl bg-violet-600/10 border border-violet-500/20 p-2 flex items-center justify-center shrink-0"
      >
        {#if iconUrl}
          <img
            src={iconUrl}
            alt={container.name}
            class="w-full h-full object-contain"
          />
        {:else}
          <span class="text-xl">📦</span>
        {/if}
      </div>

      <!-- Container Name & Subtitle Badges -->
      <div class="flex flex-col min-w-0 flex-1">
        <h3
          class="font-extrabold text-base text-white tracking-tight truncate drop-shadow-sm"
          title={container.name}
        >
          {container.name}
        </h3>

        <div class="flex items-center gap-2 mt-1.5 flex-wrap">
          <!-- Image Tag Badge -->
          <span
            class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-violet-500/15 border border-violet-400/30 text-violet-300 font-mono shadow-xs truncate max-w-[140px]"
            title={container.image}
          >
            {container.image}
          </span>

          <!-- Status Indicator Badge -->
          <span
            class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-semibold {statusBg} shadow-xs"
          >
            <span
              class="w-1.5 h-1.5 rounded-full {statusDotColor} {pulseClass} shrink-0"
            ></span>
            {statusLabel}
          </span>
        </div>
      </div>
    </div>

    <!-- Toggle Chevron Button -->
    <button
      type="button"
      class="w-9 h-9 rounded-xl bg-white/[0.06] hover:bg-white/[0.12] active:scale-95 border border-white/10 text-white/80 hover:text-white flex items-center justify-center cursor-pointer transition-all duration-200 ml-2 shrink-0 shadow-sm"
      onclick={() => (expanded = !expanded)}
      title="Mais detalhes"
      aria-label="Expandir detalhes"
    >
      <svg
        class="w-4 h-4 transition-transform duration-300 {expanded ? 'rotate-0' : 'rotate-180'}"
        fill="currentColor"
        viewBox="0 0 24 24"
      >
        <path d="M12 8.586l6.293 6.293 1.414-1.414L12 5.758 4.293 13.465l1.414 1.414z"/>
      </svg>
    </button>
  </div>

  <!-- Expanded Details Section -->
  {#if expanded}
    <div class="flex flex-col gap-3 animate-in fade-in slide-in-from-top-2 duration-300">
      <!-- 2-Column Grid: ID & Criado Em -->
      <div class="grid grid-cols-2 gap-3">
        <!-- ID Card -->
        <div
          class="flex items-center gap-3 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
        >
          <div
            class="w-10 h-10 rounded-xl bg-indigo-500/15 border border-indigo-500/30 flex items-center justify-center text-indigo-300 shrink-0 font-bold text-xs"
          >
            ID
          </div>
          <div class="flex flex-col min-w-0">
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">CONTAINER ID</span>
            <span class="font-mono font-extrabold text-sm text-indigo-300 truncate">{shortCid}</span>
          </div>
        </div>

        <!-- Criado Em Card -->
        <div
          class="flex items-center gap-3 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
        >
          <div
            class="w-10 h-10 rounded-xl bg-purple-500/15 border border-purple-500/30 flex items-center justify-center text-purple-300 shrink-0"
          >
            <svg class="w-5 h-5 fill-current" viewBox="0 0 24 24">
              <path d="M19 4h-1V2h-2v2H8V2H6v2H5c-1.11 0-1.99.9-1.99 2L3 20a2 2 0 0 0 2 2h14c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 16H5V9h14v11zM7 11h5v5H7z"/>
            </svg>
          </div>
          <div class="flex flex-col min-w-0">
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t("containers.card_created")}</span>
            <span class="font-extrabold text-sm text-white truncate">{createdStr || "—"}</span>
          </div>
        </div>
      </div>

      <!-- Live Stats se ativo -->
      {#if myStats}
        <div class="grid grid-cols-2 gap-3">
          <div
            class="flex items-center gap-3 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
          >
            <div
              class="w-10 h-10 rounded-xl bg-purple-500/15 border border-purple-500/30 flex items-center justify-center text-purple-300 shrink-0 font-bold text-xs"
            >
              CPU
            </div>
            <div class="flex flex-col min-w-0">
              <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">USO CPU</span>
              <span class="font-mono font-extrabold text-sm text-purple-300">{myStats.CPUPerc || "0%"}</span>
            </div>
          </div>

          <div
            class="flex items-center gap-3 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
          >
            <div
              class="w-10 h-10 rounded-xl bg-sky-500/15 border border-sky-500/30 flex items-center justify-center text-sky-300 shrink-0 font-bold text-xs"
            >
              RAM
            </div>
            <div class="flex flex-col min-w-0">
              <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">USO MEMÓRIA</span>
              <span class="font-mono font-extrabold text-sm text-sky-300">{myStats.MemUsage || "0B"}</span>
            </div>
          </div>
        </div>
      {/if}

      <!-- Política de Reinício -->
      <div
        class="flex items-center justify-between p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
      >
        <div class="flex items-center gap-3">
          <div
            class="w-10 h-10 rounded-xl bg-teal-500/15 border border-teal-500/30 flex items-center justify-center text-teal-300 shrink-0"
          >
            <svg class="w-5 h-5 fill-current" viewBox="0 0 24 24">
              <path d="M12 4V1L8 5l4 4V6c3.31 0 6 2.69 6 6 0 1.01-.25 1.97-.7 2.8l1.46 1.46A7.93 7.93 0 0 0 20 12c0-4.42-3.58-8-8-8zm0 14c-3.31 0-6-2.69-6-6 0-1.01.25-1.97.7-2.8L5.24 7.74A7.93 7.93 0 0 0 4 12c0 4.42 3.58 8 8 8v3l4-4-4-4v3z"/>
            </svg>
          </div>
          <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t("containers.card_restart")}</span>
        </div>
        <span class="font-extrabold text-sm text-white tracking-wide">{restartPolicyDisplay}</span>
      </div>

      <!-- Redes Conectadas -->
      {#if networkItems.length > 0}
        <div
          class="flex flex-col gap-2 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
        >
          <div class="flex items-center gap-3">
            <div
              class="w-9 h-9 rounded-xl bg-purple-500/15 border border-purple-500/30 flex items-center justify-center text-purple-300 shrink-0"
            >
              🌐
            </div>
            <span class="text-[10px] font-bold text-purple-300 uppercase tracking-wider">{t("containers.card_networks")}</span>
          </div>
          <div class="flex flex-wrap gap-1.5">
            {#each networkItems as net}
              <span
                class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-xl bg-purple-500/10 border border-purple-500/20 text-xs font-semibold text-purple-200"
              >
                {net.name} {net.ip !== "—" ? `(${net.ip})` : ""}
              </span>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Portas Mapeadas -->
      {#if portItems.length > 0}
        <div
          class="flex flex-col gap-2 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
        >
          <div class="flex items-center gap-3">
            <div
              class="w-9 h-9 rounded-xl bg-blue-500/15 border border-blue-500/30 flex items-center justify-center text-blue-300 shrink-0"
            >
              🔌
            </div>
            <span class="text-[10px] font-bold text-blue-300 uppercase tracking-wider">{t("containers.card_ports")}</span>
          </div>
          <div class="flex flex-wrap gap-1.5">
            {#each portItems as port}
              {#if port}
                <span
                  class="px-2.5 py-1 rounded-xl bg-blue-500/10 border border-blue-500/20 font-mono text-xs font-semibold text-blue-200"
                >
                  {port.formatted}
                </span>
              {/if}
            {/each}
          </div>
        </div>
      {/if}

      <!-- Container Details & Actions Footer -->
      <div class="pt-2 border-t border-white/[0.08] grid grid-cols-3 gap-2">
        <ButtonPink
          size="sm"
          class="w-full whitespace-nowrap"
          onclick={() => (showLabels = true)}
        >
          {#snippet icon()}
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="w-3.5 h-3.5"
            >
              <path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z" />
              <line x1="7" y1="7" x2="7.01" y2="7" />
            </svg>
          {/snippet}
          {t("containers.card_view_labels")}
        </ButtonPink>
        <ButtonCyan
          size="sm"
          class="w-full whitespace-nowrap"
          onclick={() => (showEnv = true)}
        >
          {#snippet icon()}
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="w-3.5 h-3.5"
            >
              <rect width="18" height="18" x="3" y="3" rx="2" />
              <path d="M7 8h10M7 12h10M7 16h6" />
            </svg>
          {/snippet}
          {t("containers.card_view_env")}
        </ButtonCyan>
        <ButtonPurple
          size="sm"
          class="w-full whitespace-nowrap"
          onclick={() => on_open_logs(container.name)}
        >
          {#snippet icon()}
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="w-3.5 h-3.5"
            >
              <polyline points="4 17 10 11 4 5" />
              <line x1="12" x2="20" y1="19" y2="19" />
            </svg>
          {/snippet}
          {t("containers.card_view_logs")}
        </ButtonPurple>
      </div>
    </div>
  {/if}

  <!-- Environment Variables Modal -->
  <FormModal bind:show={showEnv} cancelLabel={t("common.close")} title={`${t("containers.card_env_title")} — ${container.name}`}>
    {#if container.env?.length}
      <div class="grid grid-cols-[minmax(8rem,0.8fr)_minmax(0,2fr)] overflow-hidden rounded-2xl border border-cyan-500/20 bg-black/40">
        <div class="bg-cyan-950/40 px-3 py-2 text-[10px] font-bold uppercase tracking-wider text-cyan-300">{t("containers.card_env_key")}</div>
        <div class="bg-cyan-950/40 px-3 py-2 text-[10px] font-bold uppercase tracking-wider text-cyan-300">{t("containers.card_env_value")}</div>
        {#each container.env as env}
          {@const separator = env.indexOf("=")}
          {@const key = separator >= 0 ? env.slice(0, separator) : env}
          {@const value = separator >= 0 ? env.slice(separator + 1) : ""}
          <code class="block break-all border-t border-white/5 bg-white/[0.02] px-3 py-2 text-xs font-semibold text-orange-400">{key}</code>
          <code class="block break-all border-t border-white/5 bg-white/[0.02] px-3 py-2 text-xs text-slate-300">{value}</code>
        {/each}
      </div>
    {:else}
      <p class="text-sm text-slate-400">{t("containers.card_env_empty")}</p>
    {/if}
  </FormModal>

  <!-- Labels Modal -->
  <FormModal bind:show={showLabels} cancelLabel={t("common.close")} title={`${t("containers.card_labels_title")} — ${container.name}`}>
    {#if labels.length}
      <div class="grid grid-cols-[minmax(8rem,1fr)_minmax(0,2fr)] overflow-hidden rounded-2xl border border-pink-500/20 bg-black/40">
        <div class="bg-pink-950/40 px-3 py-2 text-[10px] font-bold uppercase tracking-wider text-pink-300">{t("containers.card_env_key")}</div>
        <div class="bg-pink-950/40 px-3 py-2 text-[10px] font-bold uppercase tracking-wider text-pink-300">{t("containers.card_env_value")}</div>
        {#each labels as [k, v]}
          <code class="block break-all select-all border-t border-white/5 bg-white/[0.02] px-3 py-2 text-xs font-semibold text-pink-400">{k}</code>
          <code class="block break-all select-all border-t border-white/5 bg-white/[0.02] px-3 py-2 text-xs text-slate-300">{v || "—"}</code>
        {/each}
      </div>
    {:else}
      <p class="text-sm text-slate-400">{t("containers.card_labels_empty")}</p>
    {/if}
  </FormModal>
</div>
