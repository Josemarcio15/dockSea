<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { Button, EditButtonIcon, TrashButtonIcon } from "$shared/components/buttons";

  let {
    server,
    isActive = false,
    onSelect,
    onTest,
    onEdit,
    onDelete,
  }: {
    server: any;
    isActive?: boolean;
    onSelect: (server: any) => void;
    onTest: (server: any) => void;
    onEdit: (server: any) => void;
    onDelete: (server: any) => void;
  } = $props();
</script>
<div
  class="relative rounded-[22px] bg-gradient-to-b from-slate-900/90 to-[#090d16]/95 dark:from-[#111726]/90 dark:to-[#070b13]/95 border border-white/10 dark:border-white/5 backdrop-blur-xl shadow-[0_12px_32px_rgba(0,0,0,0.45)] hover:border-violet-500/40 hover:shadow-[0_16px_36px_rgba(99,102,241,0.15)] transition-all duration-300 flex flex-col sm:flex-row sm:items-center justify-between p-4 gap-4 group text-slate-100 overflow-hidden w-full"
>
  <!-- Glow highlight top -->
  <div class="absolute -top-12 -left-12 w-32 h-32 bg-violet-600/15 rounded-full blur-2xl pointer-events-none"></div>

  <!-- Info -->
  <div class="flex items-center gap-3.5 min-w-0">
    <div
      class="w-11 h-11 rounded-xl bg-violet-600/10 border border-violet-500/20 text-violet-300 flex items-center justify-center shrink-0"
    >
      {#if server.connectionType === "ssh"}
        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.75">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 21a9.004 9.004 0 0 0 8.716-6.747M12 21a9.004 9.004 0 0 1-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 0 1 7.843 4.582M12 3a8.997 8.997 0 0 0-7.843 4.582m15.686 0A11.953 11.953 0 0 1 12 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0 1 21 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0 1 12 16.5c-3.162 0-6.133-.815-8.716-2.247m0 0A9.015 9.015 0 0 1 3 12c0-1.605.42-3.113 1.157-4.418" />
        </svg>
      {:else}
        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.75">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9 17.25v1.007a3 3 0 0 1-.879 2.122L7.5 21h9l-.621-.621A3 3 0 0 1 15 18.257V17.25m6-12V15a2.25 2.25 0 0 1-2.25 2.25H5.25A2.25 2.25 0 0 1 3 15V5.25m18 0A2.25 2.25 0 0 0 18.75 3H5.25A2.25 2.25 0 0 0 3 5.25m18 0H3" />
        </svg>
      {/if}
    </div>

    <div class="flex flex-col min-w-0">
      <div class="flex items-center gap-2">
        <span class="font-extrabold text-base text-white truncate drop-shadow-sm">
          {server.name}
        </span>
        {#if isActive}
          <span
            class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-emerald-500/15 border border-emerald-400/30 text-emerald-300 shadow-xs"
          >
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 shadow-[0_0_6px_#34d399] animate-pulse"></span>
            {t("profiles.active_badge")}
          </span>
        {/if}
      </div>

      <span class="text-xs text-indigo-300 font-mono truncate font-medium mt-0.5">
        {server.connectionType === "ssh"
          ? `${server.username}@${server.host}:${server.port}`
          : "Local Docker Engine"}
      </span>
    </div>
  </div>

  <!-- Actions Row -->
  <div class="flex items-center gap-2 self-end sm:self-center shrink-0 pt-2 sm:pt-0 border-t sm:border-t-0 border-white/[0.08]">
    {#if !isActive}
      <Button size="sm" themeKey="config.activate_btn" onclick={() => onSelect(server)}>
        {t("devices.activate")}
      </Button>
    {/if}

    <Button size="sm" themeKey="config.test_conn_btn" onclick={() => onTest(server)}>
      {t("config.test_conn_btn")}
    </Button>

    <EditButtonIcon
      size="sm"
      title={t("common.edit")}
      onclick={() => onEdit(server)}
    />

    <TrashButtonIcon
      size="sm"
      title={t("common.delete")}
      onclick={() => onDelete(server)}
    />
  </div>
</div>
