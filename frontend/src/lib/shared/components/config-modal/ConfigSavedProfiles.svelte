<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { ButtonBlue, ButtonPurple, ButtonRed } from "$shared/components/buttons";

  let {
    savedConfigs = [],
    onSelectProfile,
    onCreateNew,
    onDeleteProfile,
  }: {
    savedConfigs: any[];
    onSelectProfile: (cfg: any) => void;
    onCreateNew: () => void;
    onDeleteProfile: (id: string) => void;
  } = $props();
</script>

{#if savedConfigs.length > 0}
  <div class="flex flex-col gap-2">
    <div class="flex items-center justify-between gap-2">
      <span
        class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
      >
        {t("images.config_saved_profiles")}
      </span>
      <ButtonBlue size="xs" onclick={onCreateNew}>
        + Novo perfil
      </ButtonBlue>
    </div>
    <div
      class="flex flex-col gap-2 max-h-[calc(92vh-180px)] overflow-y-auto pr-1"
    >
      {#each savedConfigs as cfg (cfg.id)}
        <div
          class="flex items-center gap-2 p-3 rounded-xl bg-slate-50 dark:bg-slate-900 border border-slate-150 dark:border-slate-800 shadow-sm text-xs"
        >
          <div
            class="flex-1 text-slate-600 dark:text-slate-350 flex flex-col min-w-0 gap-0.5"
          >
            <div>
              <span class="font-bold text-slate-700 dark:text-slate-300"
                >{t("images.config_profile_name")}</span
              >
              <span class="font-semibold text-violet-600 dark:text-violet-400"
                >{cfg.name}</span
              >
            </div>
            {#if cfg.description}
              <div class="truncate text-slate-400 dark:text-slate-500">
                <span class="font-bold text-slate-700 dark:text-slate-300"
                  >{t("images.config_profile_desc")}</span
                >
                <span>{cfg.description}</span>
              </div>
            {/if}
          </div>
          <ButtonPurple
            size="xs"
            onclick={() => onSelectProfile(cfg)}
          >
            {t("images.config_load")}
          </ButtonPurple>
          <ButtonRed
            size="xs"
            onclick={() => onDeleteProfile(cfg.id)}
          >
            ✕
          </ButtonRed>
        </div>
      {/each}
    </div>
  </div>
{:else}
  <div
    class="flex flex-col gap-3 rounded-xl border border-slate-200 dark:border-slate-800 p-3"
  >
    <span class="text-[10px] font-bold uppercase tracking-wider text-slate-400"
      >Perfis salvos</span
    >
    <ButtonBlue size="xs" onclick={onCreateNew}>
      + Novo perfil
    </ButtonBlue>
  </div>
{/if}
