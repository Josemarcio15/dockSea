<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { Button, EditButtonIcon } from "$shared/components/buttons";
  let { profile, active, canDelete, onSelect, onEdit, onDelete } = $props();
</script>

<div
  class="relative rounded-[22px] bg-gradient-to-b from-slate-900/90 to-[#090d16]/95 dark:from-[#111726]/90 dark:to-[#070b13]/95 border border-white/10 dark:border-white/5 backdrop-blur-xl shadow-[0_12px_32px_rgba(0,0,0,0.45)] hover:border-violet-500/40 hover:shadow-[0_16px_36px_rgba(99,102,241,0.15)] transition-all duration-300 flex flex-col justify-between p-5 gap-4 group text-slate-100 overflow-hidden w-full"
>
  <!-- Glow highlight top -->
  <div class="absolute -top-12 -left-12 w-32 h-32 bg-violet-600/15 rounded-full blur-2xl pointer-events-none"></div>

  <div class="space-y-3">
    <!-- Header Section (Glassmorphism Pill) -->
    <div
      class="relative flex items-center justify-between p-3.5 rounded-2xl bg-white/[0.04] dark:bg-white/[0.03] border border-white/[0.08] backdrop-blur-md shadow-inner"
    >
      <div class="flex items-center gap-3.5 min-w-0 flex-1">
        <!-- User / Identity Icon -->
        <div
          class="w-11 h-11 rounded-xl bg-violet-600/10 border border-violet-500/20 flex items-center justify-center text-xl shrink-0"
        >
          <span>👤</span>
        </div>

        <div class="flex flex-col min-w-0 flex-1">
          <h3 class="font-extrabold text-white text-lg tracking-tight truncate drop-shadow-sm">
            {profile.name}
          </h3>
          <span class="text-xs text-slate-400 font-mono">
            {profile.id ? profile.id.substring(0, 10) : "Profile"}
          </span>
        </div>
      </div>

      {#if active}
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
    </div>

    <!-- Tags de Idioma e Tema -->
    <div class="flex items-center gap-2 text-[11px] font-medium text-slate-400 px-1">
      <span class="inline-flex items-center gap-1 bg-slate-800/60 px-2 py-0.5 rounded-md border border-slate-700/50">
        🌐 {profile.locale || "pt-BR"}
      </span>
      <span class="inline-flex items-center gap-1 bg-slate-800/60 px-2 py-0.5 rounded-md border border-slate-700/50 font-mono">
        🎨 {profile.theme || "default"}.json
      </span>
    </div>
  </div>

  <div
    class="flex items-center gap-2 pt-3 border-t border-white/[0.08]"
  >
    {#if !active}
      <Button
        size="sm"
        themeKey="profiles.select_btn"
        class="flex-1"
        onclick={onSelect}
      >
        {t("profiles.select_btn")}
      </Button>
    {/if}

    <EditButtonIcon
      size="sm"
      themeKey="profiles.edit_btn"
      title={t("common.edit")}
      onclick={onEdit}
    >
      {t("common.edit")}
    </EditButtonIcon>

    {#if canDelete}
      <Button
        size="sm"
        themeKey="profiles.delete_btn"
        title={t("common.delete")}
        onclick={onDelete}
      >
        {t("common.delete")}
      </Button>
    {/if}
  </div>
</div>
