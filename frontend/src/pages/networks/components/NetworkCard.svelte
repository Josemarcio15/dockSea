<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { ButtonGreen, ButtonRed } from "$shared/components/buttons";
  import type { DockerNetwork } from "$lib/domains/networks";
  import { isDefaultNetwork as checkIsDefault } from "$lib/domains/networks";

  let {
    network,
    checked = false,
    on_toggle = () => {},
    on_disconnect = (containerName: string) => {},
    on_delete = () => {},
    on_connect = () => {},
  }: {
    network: DockerNetwork;
    checked?: boolean;
    on_toggle?: () => void;
    on_disconnect?: (containerName: string) => void | Promise<void>;
    on_delete?: () => void | Promise<void>;
    on_connect?: () => void;
  } = $props();

  let expanded = $state(false);

  const isDefaultNetwork = $derived(checkIsDefault(network.name));

  // Curated accent border and badge based on driver
  const [accentClass, driverBadgeClass] = $derived.by(() => {
    switch (network.driver) {
      case "bridge":
        return [
          "border-l-blue-500",
          "bg-blue-50 dark:bg-blue-950/30 border-blue-200 dark:border-blue-900/40 text-blue-700 dark:text-blue-400",
        ];
      case "host":
        return [
          "border-l-amber-500",
          "bg-amber-50 dark:bg-amber-950/30 border-amber-200 dark:border-amber-900/40 text-amber-700 dark:text-amber-400",
        ];
      case "overlay":
        return [
          "border-l-purple-500",
          "bg-purple-50 dark:bg-purple-950/30 border-purple-200 dark:border-purple-900/40 text-purple-700 dark:text-purple-400",
        ];
      case "macvlan":
        return [
          "border-l-rose-500",
          "bg-rose-50 dark:bg-rose-950/30 border-rose-200 dark:border-rose-900/40 text-rose-700 dark:text-rose-400",
        ];
      default:
        return [
          "border-l-slate-400",
          "bg-slate-50 dark:bg-slate-800 border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-400",
        ];
    }
  });

  const containersCount = $derived(network.containers?.length || 0);
  const labels = $derived(Object.entries(network.labels || {}));
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
        class="w-5 h-5 rounded-lg border flex items-center justify-center transition-all duration-200 shrink-0 {checked
          ? 'bg-gradient-to-tr from-violet-600 to-indigo-500 border-violet-400 text-white shadow-[0_0_10px_rgba(124,58,237,0.5)]'
          : isDefaultNetwork
            ? 'border-white/10 bg-white/5 cursor-not-allowed opacity-30'
            : 'border-white/20 bg-white/5 hover:border-violet-400 cursor-pointer'}"
        disabled={isDefaultNetwork}
        onclick={on_toggle}
        aria-label="Selecionar rede"
      >
        {#if checked}
          <span class="text-white text-[11px] font-extrabold leading-none">✓</span>
        {/if}
      </button>

      <!-- Network Icon Container -->
      <div
        class="w-11 h-11 rounded-xl bg-violet-600/10 border border-violet-500/20 flex items-center justify-center text-xl shrink-0"
      >
        <span>🌐</span>
      </div>

      <!-- Nome + Driver + Status -->
      <div class="flex flex-col min-w-0 flex-1">
        <h3
          class="font-extrabold text-base text-white tracking-tight truncate drop-shadow-sm"
          title={network.name}
        >
          {network.name}
        </h3>

        <div class="flex items-center gap-2 mt-1.5 flex-wrap">
          <!-- Driver Badge -->
          <span
            class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-blue-500/15 border border-blue-400/30 text-blue-300 font-mono shadow-xs"
          >
            {network.driver}
          </span>

          <!-- Status Indicator Badge -->
          <span
            class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-semibold {containersCount > 0
              ? 'bg-emerald-500/15 border border-emerald-400/30 text-emerald-300'
              : 'bg-slate-500/15 border border-slate-400/30 text-slate-300'} shadow-xs"
          >
            <span
              class="w-1.5 h-1.5 rounded-full {containersCount > 0
                ? 'bg-emerald-400 shadow-[0_0_6px_#34d399]'
                : 'bg-slate-400'} shrink-0"
            ></span>
            {containersCount > 0 ? `${containersCount} container(s)` : 'Sem containers'}
          </span>
        </div>
      </div>
    </div>

    <!-- Botão Expandir -->
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

  <!-- Expanded Details -->
  {#if expanded}
    <div class="flex flex-col gap-3 animate-in fade-in slide-in-from-top-2 duration-300">
      <!-- ID & Subnet Grid -->
      <div class="grid grid-cols-2 gap-3">
        <!-- ID Card -->
        <div
          class="flex items-center gap-3 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
        >
          <div
            class="w-10 h-10 rounded-xl bg-blue-500/15 border border-blue-500/30 flex items-center justify-center text-blue-300 shrink-0 font-bold text-xs"
          >
            ID
          </div>
          <div class="flex flex-col min-w-0">
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">NETWORK ID</span>
            <span class="font-mono font-extrabold text-sm text-blue-300 truncate">{network.id.substring(0, 12)}</span>
          </div>
        </div>

        <!-- Subnet Card -->
        <div
          class="flex items-center gap-3 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
        >
          <div
            class="w-10 h-10 rounded-xl bg-purple-500/15 border border-purple-500/30 flex items-center justify-center text-purple-300 shrink-0"
          >
            <svg class="w-5 h-5 fill-current" viewBox="0 0 24 24">
              <path d="M4 6h16a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2zm0 2v8h16V8H4zm4 12h8v2H8v-2z"/>
            </svg>
          </div>
          <div class="flex flex-col min-w-0">
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">SUBNET</span>
            <span class="font-mono font-extrabold text-sm text-white truncate">{network.subnet || "—"}</span>
          </div>
        </div>
      </div>

      <!-- Escopo e Gateway -->
      <div
        class="flex items-center justify-between p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
      >
        <div class="flex items-center gap-3">
          <div
            class="w-10 h-10 rounded-xl bg-violet-500/15 border border-violet-500/30 flex items-center justify-center text-violet-300 shrink-0"
          >
            <svg class="w-5 h-5 fill-current" viewBox="0 0 24 24">
              <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5a2.5 2.5 0 0 1 0-5 2.5 2.5 0 0 1 0 5z"/>
            </svg>
          </div>
          <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t("networks.card_scope")}</span>
        </div>
        <span class="font-extrabold text-sm text-white tracking-wide capitalize">{network.scope || "local"}</span>
      </div>

      {#if network.gateway}
        <div
          class="flex items-center justify-between p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
        >
          <div class="flex items-center gap-3">
            <div
              class="w-10 h-10 rounded-xl bg-teal-500/15 border border-teal-500/30 flex items-center justify-center text-teal-300 shrink-0"
            >
              <svg class="w-5 h-5 fill-current" viewBox="0 0 24 24">
                <path d="M3.9 12c0-1.71 1.39-3.1 3.1-3.1h4V7H7c-2.76 0-5 2.24-5 5s2.24 5 5 5h4v-1.9H7c-1.71 0-3.1-1.39-3.1-3.1zM8 13h8v-2H8v2zm9-6h-4v1.9h4c1.71 0 3.1 1.39 3.1 3.1s-1.39 3.1-3.1 3.1h-4V17h4c2.76 0 5-2.24 5-5s-2.24-5-5-5z"/>
              </svg>
            </div>
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t("networks.card_gateway")}</span>
          </div>
          <span class="font-mono font-extrabold text-sm text-white tracking-wide">{network.gateway}</span>
        </div>
      {/if}

      <!-- Connected Containers -->
      <div
        class="flex flex-col gap-2 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
      >
        <span class="text-[10px] font-bold text-purple-300 uppercase tracking-wider">
          {t("networks.card_containers_connected")}
        </span>
        {#if network.containers && network.containers.length > 0}
          <div class="flex flex-col gap-1.5">
            {#each network.containers as container}
              <div class="flex items-center justify-between text-xs p-2 rounded-xl bg-white/[0.02] border border-white/5">
                <div class="flex items-center gap-1.5 truncate">
                  <span
                    class="font-semibold text-slate-200 text-xs truncate"
                  >
                    {container.name}
                  </span>
                  {#if container.ip}
                    <span
                      class="font-mono text-[10px] text-slate-400"
                    >
                      ({container.ip})
                    </span>
                  {/if}
                </div>
                <ButtonRed
                  size="xs"
                  onclick={() => on_disconnect(container.name)}
                >
                  {t("networks.card_disconnect_btn")}
                </ButtonRed>
              </div>
            {/each}
          </div>
        {:else}
          <p class="text-xs italic text-slate-400 my-0.5">
            {t("networks.card_no_containers")}
          </p>
        {/if}
      </div>

      <!-- Labels -->
      {#if labels.length > 0}
        <div
          class="flex flex-col gap-2 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
        >
          <span class="text-[10px] font-bold text-pink-300 uppercase tracking-wider">
            {t("networks.card_labels")}
          </span>
          <div class="flex flex-col gap-1.5">
            {#each labels as [k, v]}
              <span
                class="font-mono text-[11px] px-2.5 py-1.5 rounded-xl bg-pink-500/10 border border-pink-500/20 text-slate-200 break-all select-all"
              >
                <span class="font-bold text-pink-400">{k}:</span> {v}
              </span>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Actions -->
      <div class="flex gap-2 pt-1 border-t border-white/[0.08]">
        <ButtonGreen
          size="sm"
          class="flex-1"
          onclick={on_connect}
        >
          {t("networks.connect_title")}
        </ButtonGreen>
        {#if !isDefaultNetwork}
          <ButtonRed
            size="sm"
            onclick={on_delete}
          >
            {t("common.delete")}
          </ButtonRed>
        {/if}
      </div>
    </div>
  {/if}
</div>
