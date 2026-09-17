<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { ButtonBlue, ButtonRed } from "$shared/components/buttons";
  import type { EnvVar } from "./types";

  let {
    envs = $bindable([]),
    onModified = () => {},
  }: {
    envs: EnvVar[];
    onModified: () => void;
  } = $props();

  function addEnv() {
    envs = [...envs, { name: "", value: "" }];
    onModified();
  }

  function removeEnv(idx: number) {
    envs = envs.filter((_, i) => i !== idx);
    onModified();
  }
</script>

<div class="flex flex-col gap-2">
  <div class="flex items-center justify-between">
    <span
      class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
      >{t("images.config_envs")}</span
    >
    <ButtonBlue size="xs" onclick={addEnv}>
      {t("images.config_add_env")}
    </ButtonBlue>
  </div>
  {#each envs as env, i}
    <div class="flex gap-2 items-center">
      <input
        type="text"
        class="flex-1 px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
        placeholder={t("images.config_placeholder_env_name")}
        bind:value={env.name}
        oninput={onModified}
      />
      <span class="text-slate-400">=</span>
      <input
        type="text"
        class="flex-1 px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
        placeholder={t("images.config_placeholder_env_value")}
        bind:value={env.value}
        oninput={onModified}
      />
      <ButtonRed size="xs" onclick={() => removeEnv(i)}>
        ✕
      </ButtonRed>
    </div>
  {/each}
</div>
