<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import PageTitle from "$shared/components/PageTitle.svelte";

  import ColumnLayoutSwitcher from "$shared/components/ColumnLayoutSwitcher.svelte";
  import { Button } from "$shared/components/buttons";

  let {
    searchQuery = $bindable(""),
    selectedCount = 0,
    totalCount = 0,
    allSelected = false,
    onRefresh = () => {},
    onNewVolume = () => {},
    onToggleAll = () => {},
    onPrune = () => {},
    onDeleteSelected = () => {},
  }: {
    searchQuery?: string;
    selectedCount?: number;
    totalCount?: number;
    allSelected?: boolean;
    onRefresh?: () => void;
    onNewVolume?: () => void;
    onToggleAll?: () => void;
    onPrune?: () => void;
    onDeleteSelected?: () => void;
  } = $props();
</script>

<div class="space-y-6">
  <!-- Top Header -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <PageTitle title={t("volumes.title")} />

    <div class="flex items-center gap-2">
      <!-- 1, 2, 3 Colunas Switcher -->
      <ColumnLayoutSwitcher />

      <input
        type="text"
        placeholder={t("volumes.search_placeholder")}
        class="px-4 py-2 text-xs rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-[#0b1d2e] text-slate-855 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 transition-all w-60"
        bind:value={searchQuery}
      />
      <Button variant="neutral"
        size="sm"
        title={t("common.refresh")}
        onclick={onRefresh}
      >
        {t("common.refresh")}
      </Button>
      <Button variant="success"
        size="sm"
        onclick={onNewVolume}
      >
        {t("volumes.new_volume")}
      </Button>
    </div>
  </div>

  <!-- Action Toolbar -->
  <div
    class="flex flex-wrap items-center justify-between gap-3 bg-white dark:bg-[#0b1d2e] border border-slate-200/80 dark:border-slate-800/80 p-3.5 rounded-2xl shadow-sm"
  >
    <div class="flex items-center gap-2">
      <Button variant="success"
        size="sm"
        onclick={onToggleAll}
      >
        {allSelected ? t("common.deselect_all") : t("common.select_all")}
      </Button>

      {#if selectedCount > 0}
        <span
          class="text-xs font-semibold text-violet-600 dark:text-violet-400 animate-pulse px-2"
        >
          {selectedCount}
          {t("volumes.selected_count")}
        </span>
      {/if}

      <span
        class="text-xs text-slate-400 dark:text-slate-500 px-2 font-semibold"
      >
        {totalCount} {t("volumes.vol_count")}
      </span>
    </div>

    <div class="flex items-center gap-2">
      <Button variant="warning"
        size="sm"
        onclick={onPrune}
      >
        {t("volumes.prune_btn")}
      </Button>
      <Button variant="danger"
        size="sm"
        disabled={selectedCount === 0}
        onclick={onDeleteSelected}
      >
        {t("volumes.delete_selected")}
      </Button>
    </div>
  </div>
</div>
