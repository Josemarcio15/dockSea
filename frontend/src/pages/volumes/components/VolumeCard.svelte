<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import type { DockerVolume } from "$lib/domains/volumes";
  import type { VpsServer } from "../api";
  import { getVolumeSize } from "$lib/domains/volumes";
  import { notifySuccess } from "$shared/stores/notification.svelte";

  let {
    vol,
    server,
    checked = false,
    on_toggle = () => {},
  }: {
    vol: DockerVolume;
    server?: VpsServer;
    checked?: boolean;
    on_toggle?: () => void;
  } = $props();

  let expanded = $state(false);
  let fetchedSize = $state<string | null>(null);
  let loadingSize = $state(false);
  let displaySize = $derived(fetchedSize || vol.size || "—");

  // Fetch size when expanded
  $effect(() => {
    if (expanded && (displaySize === "—" || !displaySize) && !loadingSize && server) {
      loadingSize = true;
      getVolumeSize(server, vol.name)
        .then((s) => {
          if (s) {
            fetchedSize = s;
          }
        })
        .catch((err) => {
          console.error("Erro ao buscar tamanho do volume:", err);
        })
        .finally(() => {
          loadingSize = false;
        });
    }
  });

  const tagText = $derived(
    vol.inUse ? t("volumes.in_use") : t("volumes.free"),
  );

  function formatDate(iso: string) {
    if (!iso) return "";
    const date = new Date(iso);
    if (isNaN(date.getTime())) return iso;
    const months = [
      "jan",
      "fev",
      "mar",
      "abr",
      "mai",
      "jun",
      "jul",
      "ago",
      "set",
      "out",
      "nov",
      "dez",
    ];
    const day = date.getDate();
    const month = months[date.getMonth()];
    const year = date.getFullYear();
    return `${day} ${month} ${year}`;
  }

  const formattedDate = $derived(formatDate(vol.createdAt));
  const labels = $derived(Object.entries(vol.labels || {}));
  const containers = $derived(
    ((vol.containers || []).filter(Boolean) as string[][]).filter(
      (c) => c && c.length >= 2,
    ),
  );

  function copyMountpoint() {
    if (vol.mountpoint) {
      navigator.clipboard.writeText(vol.mountpoint);
      notifySuccess("Ponto de montagem copiado!");
    }
  }
</script>

<div
  class="relative rounded-[22px] bg-gradient-to-b from-slate-900/90 to-[#090d16]/95 dark:from-[#111726]/90 dark:to-[#070b13]/95 border border-white/10 dark:border-white/5 backdrop-blur-xl shadow-[0_12px_32px_rgba(0,0,0,0.45)] hover:border-violet-500/40 hover:shadow-[0_16px_36px_rgba(99,102,241,0.15)] transition-all duration-300 flex flex-col p-4 gap-3.5 group text-slate-100 overflow-hidden"
>
  <!-- Glow highlight top -->
  <div class="absolute -top-12 -left-12 w-32 h-32 bg-violet-600/15 rounded-full blur-2xl pointer-events-none"></div>

  <!-- Header Section (Glassmorphism Pill) -->
  <div
    class="relative flex items-center justify-between p-3.5 rounded-2xl bg-white/[0.04] dark:bg-white/[0.03] border border-white/[0.08] backdrop-blur-md shadow-inner"
  >
    <div class="flex items-center gap-3.5 min-w-0 flex-1">
      <!-- Checkbox opcional estilizado -->
      <button
        type="button"
        class="w-5 h-5 rounded-lg border flex items-center justify-center cursor-pointer transition-all duration-200 shrink-0 {checked
          ? 'bg-gradient-to-tr from-violet-600 to-indigo-500 border-violet-400 text-white shadow-[0_0_10px_rgba(124,58,237,0.5)]'
          : 'border-white/20 bg-white/5 hover:border-violet-400'}"
        onclick={on_toggle}
        aria-label="Selecionar volume"
      >
        {#if checked}
          <span class="text-white text-[11px] font-extrabold leading-none">✓</span>
        {/if}
      </button>

      <!-- Floppy Disk Icon Container -->
      <div
        class="w-11 h-11 rounded-xl bg-violet-600/10 border border-violet-500/20 flex items-center justify-center text-xl shrink-0"
      >
        <span>💾</span>
      </div>

      <!-- Volume Name & Subtitle Badges -->
      <div class="flex flex-col min-w-0 flex-1">
        <h3
          class="font-extrabold text-base text-white tracking-tight truncate drop-shadow-sm"
          title={vol.name}
        >
          {vol.name}
        </h3>

        <div class="flex items-center gap-2 mt-1.5 flex-wrap">
          <!-- Driver / Local Badge -->
          <span
            class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-blue-500/15 border border-blue-400/30 text-blue-300 shadow-xs"
          >
            <svg class="w-3 h-3 fill-current opacity-80" viewBox="0 0 24 24">
              <path d="M4 6h16a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2zm0 2v8h16V8H4zm4 12h8v2H8v-2z"/>
            </svg>
            {vol.driver || "Local"}
          </span>

          <!-- Status Indicator Badge -->
          <span
            class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-semibold {vol.inUse
              ? 'bg-emerald-500/15 border border-emerald-400/30 text-emerald-300'
              : 'bg-teal-500/15 border border-teal-400/30 text-teal-300'}"
          >
            <span
              class="w-1.5 h-1.5 rounded-full {vol.inUse
                ? 'bg-emerald-400 shadow-[0_0_6px_#34d399]'
                : 'bg-teal-400 shadow-[0_0_6px_#2dd4bf]'} shrink-0"
            ></span>
            {tagText}
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
      <!-- 2-Column Grid: Escopo & Tamanho -->
      <div class="grid grid-cols-2 gap-3">
        <!-- Escopo Card -->
        <div
          class="flex items-center gap-3 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
        >
          <div
            class="w-10 h-10 rounded-xl bg-violet-500/15 border border-violet-500/30 flex items-center justify-center text-violet-300 shrink-0"
          >
            <svg class="w-5 h-5 fill-current" viewBox="0 0 24 24">
              <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5a2.5 2.5 0 0 1 0-5 2.5 2.5 0 0 1 0 5z"/>
            </svg>
          </div>
          <div class="flex flex-col min-w-0">
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">ESCOPO</span>
            <span class="font-extrabold text-sm text-white capitalize truncate">{vol.scope || "Local"}</span>
          </div>
        </div>

        <!-- Tamanho Card -->
        <div
          class="flex items-center gap-3 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
        >
          <div
            class="w-10 h-10 rounded-xl bg-blue-500/15 border border-blue-500/30 flex items-center justify-center text-blue-300 shrink-0"
          >
            <svg class="w-5 h-5 fill-current" viewBox="0 0 24 24">
              <path d="M12 2C6.48 2 2 4.02 2 6.5s4.48 4.5 10 4.5 10-2.02 10-4.5S17.52 2 12 2zm0 6c-3.86 0-7-1.12-7-2.5S8.14 3 12 3s7 1.12 7 2.5S15.86 8 12 8zm0 4.5c-5.52 0-10-2.02-10-4.5v3c0 2.48 4.48 4.5 10 4.5s10-2.02 10-4.5v-3c0 2.48-4.48 4.5-10 4.5zm0 6c-5.52 0-10-2.02-10-4.5v3c0 2.48 4.48 4.5 10 4.5s10-2.02 10-4.5v-3c0 2.48-4.48 4.5-10 4.5z"/>
            </svg>
          </div>
          <div class="flex flex-col min-w-0">
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t("volumes.label_size")}</span>
            <span class="font-extrabold text-sm text-white truncate">
              {#if loadingSize}
                <span class="inline-block animate-pulse text-violet-400">...</span>
              {:else}
                {displaySize}
              {/if}
            </span>
          </div>
        </div>
      </div>

      <!-- Criado Em -->
      <div
        class="flex items-center justify-between p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
      >
        <div class="flex items-center gap-3">
          <div
            class="w-10 h-10 rounded-xl bg-purple-500/15 border border-purple-500/30 flex items-center justify-center text-purple-300 shrink-0"
          >
            <svg class="w-5 h-5 fill-current" viewBox="0 0 24 24">
              <path d="M19 4h-1V2h-2v2H8V2H6v2H5c-1.11 0-1.99.9-1.99 2L3 20a2 2 0 0 0 2 2h14c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 16H5V9h14v11zM7 11h5v5H7z"/>
            </svg>
          </div>
          <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t("volumes.label_created")}</span>
        </div>
        <span class="font-extrabold text-sm text-white tracking-wide">{formattedDate || "—"}</span>
      </div>

      <!-- Ponto de Montagem -->
      <div
        class="flex flex-col gap-2 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
      >
        <div class="flex items-center gap-3">
          <div
            class="w-9 h-9 rounded-xl bg-indigo-500/15 border border-indigo-500/30 flex items-center justify-center text-indigo-300 shrink-0"
          >
            <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
              <path d="M3.9 12c0-1.71 1.39-3.1 3.1-3.1h4V7H7c-2.76 0-5 2.24-5 5s2.24 5 5 5h4v-1.9H7c-1.71 0-3.1-1.39-3.1-3.1zM8 13h8v-2H8v2zm9-6h-4v1.9h4c1.71 0 3.1 1.39 3.1 3.1s-1.39 3.1-3.1 3.1h-4V17h4c2.76 0 5-2.24 5-5s-2.24-5-5-5z"/>
            </svg>
          </div>
          <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t("volumes.label_mount")}</span>
        </div>

        <div
          class="flex items-center justify-between gap-2 p-3 rounded-xl bg-black/40 border border-white/[0.05] group/mount"
        >
          <span
            class="font-mono text-xs text-slate-300 break-all select-all leading-relaxed"
          >
            {vol.mountpoint}
          </span>
          <button
            type="button"
            class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 text-slate-400 hover:text-white transition-colors cursor-pointer shrink-0"
            onclick={copyMountpoint}
            title="Copiar caminho"
            aria-label="Copiar ponto de montagem"
          >
            <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
              <path d="M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z"/>
            </svg>
          </button>
        </div>
      </div>

      <!-- Containers Vinculados (se houver) -->
      {#if containers.length > 0}
        <div
          class="flex flex-col gap-2 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
        >
          <span class="text-[10px] font-bold text-violet-300 uppercase tracking-wider">
            {t("volumes.connected_containers")}
          </span>
          <div class="flex flex-wrap gap-1.5">
            {#each containers as c}
              <span
                class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-xl bg-violet-500/10 border border-violet-500/20 text-xs font-semibold text-violet-200"
              >
                {c[0]} <span class="text-[10px] opacity-60 uppercase font-mono">({c[1]})</span>
              </span>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Labels (se houver) -->
      {#if labels.length > 0}
        <div
          class="flex flex-col gap-2 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
        >
          <span class="text-[10px] font-bold text-pink-300 uppercase tracking-wider">
            {t("volumes.label_labels")}
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
    </div>
  {/if}
</div>
