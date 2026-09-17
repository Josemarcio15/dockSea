<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import { Button } from "$shared/components/buttons";
  import type { StackItem } from "$lib/domains/stacks";

  let {
    stack,
    on_deploy = () => {},
    on_stop = () => {},
    on_logs = () => {},
    on_edit = () => {},
    on_delete_local = () => {},
    on_remove_remote = () => {},
  }: {
    stack: StackItem;
    on_deploy?: () => void | Promise<void>;
    on_stop?: () => void | Promise<void>;
    on_logs?: () => void | Promise<void>;
    on_edit?: () => void | Promise<void>;
    on_delete_local?: () => void | Promise<void>;
    on_remove_remote?: () => void | Promise<void>;
  } = $props();

  let formattedDate = $derived(
    stack.updatedAt ? new Date(stack.updatedAt).toLocaleString() : "",
  );

  let formattedDeployDate = $derived(
    stack.lastDeployedAt ? new Date(stack.lastDeployedAt).toLocaleString() : "",
  );
</script>

<div
  class="relative rounded-[22px] bg-gradient-to-b from-slate-900/90 to-[#090d16]/95 dark:from-[#111726]/90 dark:to-[#070b13]/95 border border-white/10 dark:border-white/5 backdrop-blur-xl shadow-[0_12px_32px_rgba(0,0,0,0.45)] hover:border-violet-500/40 hover:shadow-[0_16px_36px_rgba(99,102,241,0.15)] transition-all duration-300 flex flex-col xl:flex-row xl:items-center justify-between p-4 gap-4 group text-slate-100 overflow-hidden w-full"
>
  <!-- Glow highlight top -->
  <div class="absolute -top-12 -left-12 w-32 h-32 bg-violet-600/15 rounded-full blur-2xl pointer-events-none"></div>

  <!-- Stack Main Information -->
  <div class="flex items-center gap-3.5 min-w-0 flex-1">
    <!-- Stack Icon Container -->
    <div
      class="w-11 h-11 rounded-xl bg-violet-600/10 border border-violet-500/20 flex items-center justify-center text-xs font-black text-violet-300 shrink-0"
    >
      <span>{stack.sourceType === "folder" ? "DIR" : "YAML"}</span>
    </div>

    <div class="min-w-0 flex-1 space-y-1.5">
      <!-- Row 1: Title and Status Badges -->
      <div class="flex items-center gap-2.5 flex-wrap">
        <h3
          class="font-extrabold text-base text-white tracking-tight truncate drop-shadow-sm m-0"
          title={stack.name}
        >
          {stack.name}
        </h3>

        {#if stack.sourceType === "folder"}
          <span class="px-2.5 py-0.5 text-[11px] font-semibold rounded-full bg-amber-500/15 text-amber-300 border border-amber-500/30 shrink-0">
            {t("stacks.source_type_folder")}
          </span>
        {:else}
          <span class="px-2.5 py-0.5 text-[11px] font-semibold rounded-full bg-blue-500/15 text-blue-300 border border-blue-500/30 shrink-0">
            {t("stacks.source_type_editor")}
          </span>
        {/if}

        {#if stack.lastDeployedAt}
          <span class="px-2.5 py-0.5 text-[11px] font-semibold rounded-full bg-emerald-500/15 text-emerald-300 border border-emerald-500/30 flex items-center gap-1.5 shrink-0">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 shadow-[0_0_6px_#34d399] animate-pulse"></span>
            {t("stacks.last_deployed_label")} {formattedDeployDate}
          </span>
        {:else}
          <span class="px-2.5 py-0.5 text-[11px] font-medium rounded-full bg-slate-500/15 text-slate-400 border border-slate-400/30 shrink-0">
            {t("stacks.never_deployed")}
          </span>
        {/if}
      </div>

      <!-- Row 2: Stack Metadata -->
      <div
        class="flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-slate-400"
      >
        <span
          class="font-mono bg-white/[0.04] border border-white/[0.08] px-2 py-0.5 rounded-md text-xs text-indigo-300 font-medium"
        >
          {stack.projectName || stack.id}
        </span>

        {#if stack.folderPath}
          <span class="text-slate-600">•</span>
          <span class="font-mono text-xs text-slate-400 truncate max-w-xs" title={stack.folderPath}>
            {stack.folderPath}
          </span>
        {/if}

        {#if stack.createdAt}
          <span class="text-slate-600">•</span>
          <span>
            {t("stacks.created_at")}: {new Date(stack.createdAt).toLocaleDateString()}
          </span>
        {/if}

        {#if formattedDate}
          <span class="text-slate-600">•</span>
          <span class="text-slate-400">
            {t("stacks.updated_at")}: {formattedDate}
          </span>
        {/if}
      </div>
    </div>
  </div>

  <!-- Stack Actions -->
  <div class="flex flex-wrap items-center gap-2 shrink-0 pt-3 xl:pt-0 border-t xl:border-t-0 border-white/[0.08]">
    <!-- Deploy (Green) -->
    <Button size="sm" themeKey="stacks.deploy_btn" onclick={on_deploy}>
      {#snippet icon()}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="currentColor"
          class="w-3.5 h-3.5"
        >
          <path
            fill-rule="evenodd"
            d="M4.5 5.653c0-1.427 1.529-2.33 2.779-1.643l11.54 6.347c1.295.712 1.295 2.573 0 3.286L7.28 19.99c-1.25.687-2.779-.217-2.779-1.643V5.653Z"
            clip-rule="evenodd"
          />
        </svg>
      {/snippet}
      {t("stacks.deploy_btn")}
    </Button>

    <!-- Stop / Remove Container (Yellow / Amber) -->
    <Button size="sm" themeKey="stacks.stop_btn" onclick={on_stop}>
      {#snippet icon()}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="currentColor"
          class="w-3.5 h-3.5"
        >
          <rect width="12" height="12" x="6" y="6" rx="2" />
        </svg>
      {/snippet}
      {t("stacks.stop_btn")}
    </Button>

    <!-- Logs (Blue) -->
    <Button size="sm" themeKey="stacks.logs_btn" onclick={on_logs}>
      {#snippet icon()}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="w-3.5 h-3.5"
        >
          <polyline points="4 17 10 11 4 5" />
          <line x1="12" x2="20" y1="19" y2="19" />
        </svg>
      {/snippet}
      {t("stacks.logs_btn")}
    </Button>

    <!-- Edit (Purple) -->
    <Button size="sm" themeKey="stacks.edit_btn" onclick={on_edit}>
      {#snippet icon()}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="w-3.5 h-3.5"
        >
          <path
            d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"
          />
          <path d="m15 5 4 4" />
        </svg>
      {/snippet}
      {t("stacks.edit_btn")}
    </Button>

    <!-- Removes from VPS / Down (Orange) -->
    <Button size="sm" themeKey="stacks.remove_remote_btn" onclick={on_remove_remote}>
      {#snippet icon()}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="w-3.5 h-3.5"
        >
          <path d="M12 2v10" />
          <path d="m16 8-4 4-4-4" />
          <path d="M2 17h20" />
          <path d="M6 21h12" />
        </svg>
      {/snippet}
      {t("stacks.remove_remote_btn")}
    </Button>

    <!-- Delete Local Definition (Red / Rose) -->
    <Button size="sm" themeKey="stacks.delete_local_btn" onclick={on_delete_local}>
      {#snippet icon()}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="w-3.5 h-3.5"
        >
          <path d="M3 6h18" />
          <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
          <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
        </svg>
      {/snippet}
      {t("stacks.delete_local_btn")}
    </Button>
  </div>
</div>
