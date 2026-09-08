<script lang="ts">
  import { onMount } from "svelte";
  import PageTitle from "$shared/components/PageTitle.svelte";

  import { t } from "$shared/stores/locale.svelte";
  import StatusBanner from "$shared/components/StatusBanner.svelte";
  import VpsSelectWarning from "$shared/components/VpsSelectWarning.svelte";
  import { Events } from "@wailsio/runtime";
  import { builderStore as store } from "./store.svelte";
  import FolderBrowser from "./components/FolderBrowser.svelte";
  import BuildControls from "./components/BuildControls.svelte";
  import TaskProgressModal from "$shared/components/TaskProgressModal.svelte";

  let { data } = $props();

  function goToImages() {
    window.location.href = `/images?highlight=${encodeURIComponent(store.builtImage)}`;
  }

  onMount(() => {
    const unsubscribeProgress = Events.On("builder:progress", (event) => {
      if (event.data?.line) store.appendLog(event.data.line);
    });
    const unsubscribeComplete = Events.On("builder:complete", (event) => {
      store.completeBuild(
        event.data as { success?: boolean; image?: string; message?: string },
      );
    });

    store.loadSavedPaths();
    store.browse();

    return () => {
      unsubscribeProgress();
      unsubscribeComplete();
    };
  });
</script>

{#if !data.activeVps}
  <VpsSelectWarning />
{:else}
  <div class="space-y-6 max-w-4xl mx-auto">
    <div
      class="flex flex-col sm:flex-row sm:items-center justify-between gap-4"
    >
      <PageTitle title={t("sidebar.builder")} />
    </div>

    <StatusBanner />

    <div class="space-y-4">
      <FolderBrowser {store} />
      <BuildControls {store} {goToImages} />
    </div>
  </div>

  <TaskProgressModal
    bind:show={store.showProgressModal}
    title={`Build de Imagem: ${store.effectiveTag}`}
    eventPrefix="builder"
    oncomplete={() => {
      store.completeBuild({ success: true, image: store.effectiveTag });
    }}
  />
{/if}
