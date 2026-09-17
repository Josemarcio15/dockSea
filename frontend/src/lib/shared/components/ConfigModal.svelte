<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { Button } from "$shared/components/buttons";
  import * as VolumeService from "$bindings/volumes/volumeservice.js";
  import * as NetworkService from "$bindings/networks/networkservice.js";
  import type { VpsServer } from "$bindings/core/db/models.js";

  import ConfigSavedProfiles from "./config-modal/ConfigSavedProfiles.svelte";
  import ConfigJsonEditor from "./config-modal/ConfigJsonEditor.svelte";
  import ConfigFieldsForm from "./config-modal/ConfigFieldsForm.svelte";
  import ConfigPersistenceDialog from "./config-modal/ConfigPersistenceDialog.svelte";
  import { getDefaultContainerPath } from "./config-modal/types";
  import { createConfigModalState } from "./config-modal/formState.svelte";

  let {
    show = $bindable(false),
    image = "",
    savedConfigs = [],
    serverId = "",
    activeVps,
    onsubmit = (config: any) => {},
    onsaveprofile = (profile: any) => {},
    ondeleteprofile = (profileId: string) => {},
  }: {
    show: boolean;
    image: string;
    savedConfigs: any[];
    serverId: string;
    activeVps?: VpsServer;
    onsubmit: (config: any) => void;
    onsaveprofile: (profile: any) => void;
    ondeleteprofile: (profileId: string) => void;
  } = $props();

  let existingVolumes = $state<string[]>([]);
  let existingNetworks = $state<string[]>([]);
  let showConfirmNoVolume = $state(false);

  const form = createConfigModalState(
    () => savedConfigs,
    () => image,
    (p) => onsaveprofile(p),
  );

  $effect(() => {
    if (!show || !activeVps) return;
    existingVolumes = [];
    existingNetworks = [];
    VolumeService.ListVolumes(activeVps)
      .then((vols) => {
        existingVolumes = (vols || []).map((v) => v.name);
      })
      .catch(() => {});
    NetworkService.ListNetworks(activeVps)
      .then((nets) => {
        existingNetworks = (nets || []).map((n) => n.name);
      })
      .catch(() => {});
  });

  $effect(() => {
    if (show) form.initForImage(image);
    if (!show) form.resetState();
  });

  function validateAndSubmit() {
    if (form.volumes.length === 0) {
      showConfirmNoVolume = true;
    } else {
      submitForm();
    }
  }

  function submitForm() {
    showConfirmNoVolume = false;
    show = false;

    const volList = form.volumes
      .filter((v) => v.host.trim() && v.container.trim())
      .map((v) => `${v.host}:${v.container || getDefaultContainerPath(image)}`);
    const commandStr = form.commands.filter((c) => c.trim() !== "").join(" && ");

    onsubmit({
      containerName: form.containerName,
      image,
      ports: form.ports,
      envs: form.envs,
      volumes: volList,
      network: form.network,
      restartPolicy: form.restartPolicy,
      projectName: form.projectName,
      command: commandStr || undefined,
    });
  }
</script>

{#if show}
  <div
    class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
  >
    <div
      class="bg-white dark:bg-[#111827] border border-slate-200 dark:border-slate-800 rounded-2xl shadow-2xl w-[90vw] h-[90vh] max-w-none max-h-none flex flex-col text-slate-800 dark:text-slate-100 overflow-hidden"
    >
      <!-- Header -->
      <div
        class="flex justify-between items-center px-6 py-4 border-b border-slate-200 dark:border-slate-850 bg-slate-50 dark:bg-slate-900/50"
      >
        <h2 class="text-base font-bold text-slate-850 dark:text-slate-100">
          {t("images.config_title")}
        </h2>
        <button
          type="button"
          class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 text-lg bg-transparent border-none cursor-pointer transition-colors"
          onclick={() => (show = false)}
        >
          ✕
        </button>
      </div>

      <!-- Body -->
      <div
        class="p-6 overflow-y-auto flex-1 grid grid-cols-1 lg:grid-cols-[280px_minmax(0,1fr)] gap-6 min-h-0 pr-3"
      >
        <ConfigSavedProfiles
          {savedConfigs}
          onSelectProfile={form.loadProfile}
          onCreateNew={form.createNewProfile}
          onDeleteProfile={ondeleteprofile}
        />

        <div class="lg:col-start-2 flex flex-col min-w-0 gap-4">
          <div
            class="flex rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/40 p-1"
          >
            <button
              type="button"
              class="flex-1 rounded-lg px-3 py-2 text-xs font-bold {form.editorTab === 'fields'
                ? 'bg-white dark:bg-slate-800 text-violet-600 shadow-sm'
                : 'text-slate-500'}"
              onclick={() => (form.editorTab = "fields")}>Campos</button
            >
            <button
              type="button"
              class="flex-1 rounded-lg px-3 py-2 text-xs font-bold {form.editorTab === 'json'
                ? 'bg-white dark:bg-slate-800 text-violet-600 shadow-sm'
                : 'text-slate-500'}"
              onclick={() => {
                form.editorTab = "json";
                form.syncJsonEditor();
              }}>JSON</button
            >
          </div>

          {#if form.editorTab === "json"}
            <ConfigJsonEditor
              bind:jsonEditor={form.jsonEditor}
              jsonMessage={form.jsonMessage}
              jsonIsValid={form.jsonIsValid}
              saveDisabled={form.saveDisabled}
              onLoadExample={form.loadExample}
              onCopyJson={form.copyJson}
              onSaveProfile={form.saveJsonProfile}
              onJsonChange={form.handleJsonChange}
            />
          {/if}

          {#if form.editorTab === "fields"}
            <ConfigFieldsForm
              {image}
              bind:containerName={form.containerName}
              bind:projectName={form.projectName}
              bind:ports={form.ports}
              bind:envs={form.envs}
              bind:volumes={form.volumes}
              bind:network={form.network}
              bind:restartPolicy={form.restartPolicy}
              bind:commands={form.commands}
              bind:description={form.description}
              bind:profileName={form.profileName}
              {existingVolumes}
              {existingNetworks}
              hasEmptyVolume={form.hasEmptyVolume}
              saveDisabled={form.saveDisabled}
              nameAlreadyExists={form.nameAlreadyExists}
              loadedProfileId={form.loadedProfileId}
              isNameChanged={form.isNameChanged}
              onModified={() => (form.isModified = true)}
              onTriggerSave={form.triggerSave}
            />
          {/if}
        </div>
      </div>

      <!-- Footer -->
      <div
        class="flex justify-end gap-3 px-6 py-4 border-t border-slate-200 dark:border-slate-850 bg-slate-50 dark:bg-slate-900/50"
      >
        <Button onclick={() => (show = false)}>
          {t("common.cancel")}
        </Button>
        <Button
          disabled={form.hasEmptyVolume}
          onclick={validateAndSubmit}
        >
          {t("images.config_create_container")}
        </Button>
      </div>
    </div>
  </div>
{/if}

<ConfigPersistenceDialog
  bind:show={showConfirmNoVolume}
  onConfirm={submitForm}
/>
