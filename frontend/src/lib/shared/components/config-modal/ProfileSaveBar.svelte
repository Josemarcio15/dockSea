<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { ButtonGreen } from "$shared/components/buttons";

  let {
    profileName = $bindable(""),
    saveDisabled = false,
    loadedProfileId = null,
    isNameChanged = false,
    nameAlreadyExists = false,
    onTriggerSave = () => {},
  }: {
    profileName: string;
    saveDisabled: boolean;
    loadedProfileId: string | null;
    isNameChanged: boolean;
    nameAlreadyExists: boolean;
    onTriggerSave: () => void;
  } = $props();
</script>

<div
  class="flex flex-col gap-1.5 p-3.5 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/20"
>
  <label
    for="config-profile-name"
    class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >{t("images.config_save_profile")}</label
  >
  <div class="flex gap-2 items-center">
    <input
      id="config-profile-name"
      type="text"
      class="flex-1 px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
      placeholder={t("images.config_placeholder_profile")}
      bind:value={profileName}
    />
    {#if profileName.trim()}
      <ButtonGreen
        size="sm"
        disabled={saveDisabled}
        onclick={onTriggerSave}
      >
        {loadedProfileId && !isNameChanged
          ? t("images.config_btn_update")
          : t("images.config_btn_create_profile")}
      </ButtonGreen>
    {/if}
  </div>
  {#if nameAlreadyExists}
    <div class="text-[11px] text-red-500 font-semibold mt-1">
      Porta ou Perfil já cadastrado com este nome.
    </div>
  {/if}
</div>
