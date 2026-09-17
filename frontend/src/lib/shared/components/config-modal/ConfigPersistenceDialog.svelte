<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { Button } from "$shared/components/buttons";

  let {
    show = $bindable(false),
    onConfirm = () => {},
  }: {
    show: boolean;
    onConfirm: () => void;
  } = $props();
</script>

{#if show}
  <div
    class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-60"
  >
    <div
      class="bg-white dark:bg-[#1b2330] border border-slate-200 dark:border-[#3f3f46] rounded-2xl p-6 shadow-2xl max-w-sm w-full mx-4 text-slate-800 dark:text-slate-100"
    >
      <h3 class="text-sm font-bold text-slate-800 dark:text-slate-100 mb-2">
        {t("images.config_confirm_persistence")}
      </h3>
      <p class="text-xs text-red-600 dark:text-red-400 font-semibold mb-6">
        {t("images.config_persistence_warning")}
      </p>
      <div class="flex gap-3 justify-end">
        <Button
          size="sm"
          onclick={() => (show = false)}
        >
          {t("images.config_no")}
        </Button>
        <Button size="sm" onclick={onConfirm}>
          {t("images.config_yes")}
        </Button>
      </div>
    </div>
  </div>
{/if}
