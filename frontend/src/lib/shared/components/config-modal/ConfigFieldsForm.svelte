<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import type { PortMapping, EnvVar } from "./types";
  import PortsInputList from "./PortsInputList.svelte";
  import EnvsInputList from "./EnvsInputList.svelte";
  import VolumesInputList from "./VolumesInputList.svelte";
  import CommandsInputList from "./CommandsInputList.svelte";
  import ProfileSaveBar from "./ProfileSaveBar.svelte";

  let {
    image = "",
    containerName = $bindable(""),
    projectName = $bindable(""),
    ports = $bindable([]),
    envs = $bindable([]),
    volumes = $bindable([]),
    network = $bindable(""),
    restartPolicy = $bindable(""),
    commands = $bindable([]),
    description = $bindable(""),
    profileName = $bindable(""),
    existingVolumes = [],
    existingNetworks = [],
    hasEmptyVolume = false,
    saveDisabled = false,
    nameAlreadyExists = false,
    loadedProfileId = null,
    isNameChanged = false,
    onModified = () => {},
    onTriggerSave = () => {},
  }: {
    image: string;
    containerName: string;
    projectName: string;
    ports: PortMapping[];
    envs: EnvVar[];
    volumes: any[];
    network: string;
    restartPolicy: string;
    commands: string[];
    description: string;
    profileName: string;
    existingVolumes: string[];
    existingNetworks: string[];
    hasEmptyVolume: boolean;
    saveDisabled: boolean;
    nameAlreadyExists: boolean;
    loadedProfileId: string | null;
    isNameChanged: boolean;
    onModified: () => void;
    onTriggerSave: () => void;
  } = $props();
</script>

<!-- Image name -->
<div class="flex flex-col gap-1">
  <label
    for="config-image"
    class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >Imagem</label
  >
  <input
    id="config-image"
    type="text"
    class="w-full px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-100 dark:bg-[#0c101b] text-slate-400 dark:text-slate-500 cursor-not-allowed font-mono"
    value={image}
    readonly
  />
</div>

<!-- Container Name -->
<div class="flex flex-col gap-1">
  <label
    for="config-container-name"
    class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >{t("images.config_container_name")}</label
  >
  <input
    id="config-container-name"
    type="text"
    class="w-full px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
    bind:value={containerName}
    oninput={onModified}
    placeholder={t("images.config_placeholder_container")}
  />
</div>

<!-- Project Name -->
<div class="flex flex-col gap-1">
  <label
    for="config-project-name"
    class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >{t("images.config_project_name")}</label
  >
  <input
    id="config-project-name"
    type="text"
    class="w-full px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
    bind:value={projectName}
    oninput={onModified}
    placeholder={t("images.config_placeholder_project")}
  />
</div>

<!-- Ports mapping -->
<PortsInputList bind:ports {onModified} />

<!-- Environment variables -->
<EnvsInputList bind:envs {onModified} />

<!-- Volumes mapping -->
<VolumesInputList
  {image}
  bind:volumes
  {existingVolumes}
  {hasEmptyVolume}
  {onModified}
/>

<!-- Network selection with datalist -->
<div class="flex flex-col gap-1">
  <label
    for="config-network"
    class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >{t("images.config_network")}</label
  >
  <input
    id="config-network"
    type="text"
    class="w-full px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 placeholder-slate-400 focus:border-violet-500 focus:outline-none transition-colors"
    placeholder={t("images.config_placeholder_network")}
    list="networks-datalist"
    bind:value={network}
    oninput={onModified}
  />
  <datalist id="networks-datalist">
    {#each existingNetworks as net}
      <option value={net}></option>
    {/each}
  </datalist>
</div>

<!-- Restart policy -->
<div class="flex flex-col gap-1">
  <label
    for="config-restart"
    class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >{t("containers.card_restart")}</label
  >
  <select
    id="config-restart"
    class="w-full px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors"
    bind:value={restartPolicy}
  >
    <option value="no">{t("containers.card_restart_no")}</option>
    <option value="always">{t("containers.card_restart_always")}</option>
    <option value="unless-stopped">{t("containers.card_restart_unless_stopped")}</option>
    <option value="on-failure">{t("containers.card_restart_on_failure")}</option>
  </select>
</div>

<!-- Command -->
<CommandsInputList bind:commands {onModified} />

<!-- Description -->
<div class="flex flex-col gap-1">
  <label
    for="config-description"
    class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider"
    >{t("images.config_desc")}</label
  >
  <textarea
    id="config-description"
    class="w-full px-3.5 py-2 text-xs border border-slate-200 dark:border-slate-800 rounded-xl bg-slate-50 dark:bg-[#0c101b] text-slate-800 dark:text-slate-200 focus:border-violet-500 focus:outline-none transition-colors resize-y min-h-15"
    placeholder={t("images.config_placeholder_desc")}
    bind:value={description}
    oninput={onModified}
  ></textarea>
</div>

<!-- Save Profile Card -->
<ProfileSaveBar
  bind:profileName
  {saveDisabled}
  {loadedProfileId}
  {isNameChanged}
  {nameAlreadyExists}
  {onTriggerSave}
/>
