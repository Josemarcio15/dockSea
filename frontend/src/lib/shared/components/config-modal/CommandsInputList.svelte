<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { Button } from "$shared/components/buttons";

  let {
    commands = $bindable([]),
    onModified = () => {},
  }: {
    commands: string[];
    onModified: () => void;
  } = $props();

  function addCommand() {
    commands = [...commands, ""];
    onModified();
  }

  function removeCommand(idx: number) {
    commands = commands.filter((_, i) => i !== idx);
    onModified();
  }
</script>

<div class="flex flex-col gap-2">
  <div class="flex items-center justify-between">
    <span
      class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
      >{t("images.config_command")}</span
    >
    <Button size="xs" onclick={addCommand}>
      {t("images.config_add_command")}
    </Button>
  </div>

  {#each commands as cmd, i}
    <div class="flex gap-2 items-center">
      <input
        type="text"
        class="flex-1 px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors font-mono"
        placeholder={t("images.config_placeholder_command")}
        bind:value={commands[i]}
        oninput={onModified}
      />
      <Button size="xs" onclick={() => removeCommand(i)}>
        ✕
      </Button>
    </div>
  {/each}
  <span class="text-[10px] text-slate-400 dark:text-slate-500 mt-0.5 leading-relaxed">
    {t("images.config_command_hint")}
  </span>
</div>
