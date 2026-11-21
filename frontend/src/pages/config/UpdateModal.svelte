<script lang="ts">
  import FormModal from "$shared/components/FormModal.svelte";

  interface ModalButton {
    label: string;
    variant?: "primary" | "secondary" | "danger" | "success" | "warning";
    onclick: () => void | Promise<void>;
    disabled?: boolean;
  }
  import { t } from "$shared/stores/locale.svelte";
  import { updaterApi, type UpdateInfo, type InstallResult } from "./updater.api";

  let { show = $bindable(false) }: { show: boolean } = $props();

  type Phase = "checking" | "idle" | "installing" | "done" | "error";

  let phase = $state<Phase>("checking");
  let info = $state<UpdateInfo | null>(null);
  let result = $state<InstallResult | null>(null);
  let errorMsg = $state("");

  async function check() {
    phase = "checking";
    errorMsg = "";
    result = null;
    try {
      info = await updaterApi.check();
      phase = "idle";
    } catch (e) {
      errorMsg = String((e as Error)?.message ?? e);
      phase = "error";
    }
  }

  async function install() {
    phase = "installing";
    errorMsg = "";
    try {
      result = await updaterApi.install();
      phase = result.cancelled ? "idle" : "done";
    } catch (e) {
      errorMsg = String((e as Error)?.message ?? e);
      phase = "error";
    }
  }

  async function restart() {
    try {
      await updaterApi.restart();
    } catch (e) {
      errorMsg = String((e as Error)?.message ?? e);
      phase = "error";
    }
  }

  $effect(() => {
    if (show) void check();
  });

  function formatSize(bytes: number) {
    if (!bytes) return "";
    return `${(bytes / 1024 / 1024).toFixed(1)} MB`;
  }

  function formatDate(iso: string) {
    if (!iso) return "";
    try {
      return new Date(iso).toLocaleDateString();
    } catch {
      return iso;
    }
  }

  const buttons = $derived.by<ModalButton[]>(() => {
    const list: ModalButton[] = [];
    if (info?.releaseUrl) {
      list.push({
        label: t("updater.view_release"),
        variant: "secondary",
        onclick: () => updaterApi.openUrl(info!.releaseUrl),
      });
    }
    if (phase === "error") {
      list.push({ label: t("updater.retry"), variant: "secondary", onclick: check });
    }
    if (phase === "idle" && info?.hasUpdate && info.supported) {
      list.push({ label: t("updater.install_btn"), variant: "primary", onclick: install });
    }
    if (phase === "installing") {
      list.push({ label: t("updater.installing"), variant: "primary", disabled: true, onclick: () => {} });
    }
    if (phase === "done" && result?.success) {
      list.push({ label: t("updater.restart_btn"), variant: "success", onclick: restart });
    }
    return list;
  });
</script>

<FormModal bind:show title={t("updater.title")} cancelLabel={t("updater.close")} {buttons}>
  <div class="py-4 space-y-4">
    {#if phase === "checking"}
      <div class="py-8 flex flex-col items-center justify-center space-y-3">
        <div class="w-8 h-8 border-4 border-violet-500/30 border-t-violet-500 rounded-full animate-spin"></div>
        <p class="text-sm text-slate-400">{t("updater.checking")}</p>
      </div>
    {:else if phase === "installing"}
      <div class="py-8 flex flex-col items-center justify-center space-y-3 text-center">
        <div class="w-8 h-8 border-4 border-violet-500/30 border-t-violet-500 rounded-full animate-spin"></div>
        <p class="text-sm text-slate-300">{t("updater.installing")}</p>
        <p class="text-xs text-slate-500">{t("updater.installing_hint")}</p>
      </div>
    {:else}
      {#if info}
        <div class="grid grid-cols-2 gap-3">
          <div class="p-3 rounded-xl bg-slate-100 dark:bg-white/[0.03] border border-slate-200 dark:border-white/[0.06]">
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t("updater.current")}</span>
            <div class="font-mono text-sm font-semibold">{info.currentVersion}</div>
          </div>
          <div class="p-3 rounded-xl border {info.hasUpdate ? 'bg-violet-500/10 border-violet-500/30' : 'bg-slate-100 dark:bg-white/[0.03] border-slate-200 dark:border-white/[0.06]'}">
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t("updater.latest")}</span>
            <div class="font-mono text-sm font-semibold {info.hasUpdate ? 'text-violet-400' : ''}">
              {info.latestVersion || "—"}
            </div>
          </div>
        </div>

        {#if phase === "idle"}
          {#if !info.supported}
            <div class="p-3 rounded-xl bg-amber-500/10 border border-amber-500/20 text-amber-500 text-sm">
              {t("updater.unsupported")}
            </div>
          {:else if info.hasUpdate}
            <div class="p-3 rounded-xl bg-violet-500/10 border border-violet-500/20 text-sm">
              {t("updater.available", { version: info.latestVersion })}
              <span class="block text-xs text-slate-400 mt-1">
                {info.assetName} · {formatSize(info.assetSize)} · {formatDate(info.publishedAt)}
              </span>
            </div>
          {:else}
            <div class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-500 text-sm">
              {t("updater.up_to_date")}
            </div>
          {/if}
          {#if result?.cancelled}
            <p class="text-xs text-amber-500">{t("updater.cancelled")}</p>
          {/if}
        {/if}

        {#if info.hasUpdate && info.releaseNotes && phase !== "done"}
          <div class="space-y-1">
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t("updater.notes")}</span>
            <div class="max-h-56 overflow-y-auto p-3 rounded-xl bg-white dark:bg-[#0c1220] border border-slate-200/60 dark:border-slate-800/80 text-xs whitespace-pre-wrap leading-relaxed text-slate-600 dark:text-slate-300">
              {info.releaseNotes}
            </div>
          </div>
        {/if}
      {/if}

      {#if phase === "done" && result}
        {#if result.success}
          <div class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-500 text-sm">
            {t("updater.success")}
          </div>
        {:else}
          <div class="p-3 rounded-xl bg-amber-500/10 border border-amber-500/20 text-sm space-y-2">
            <p class="text-amber-500 font-semibold">{t("updater.fallback_title")}</p>
            <p class="text-xs text-slate-400 whitespace-pre-wrap font-mono">{result.message}</p>
            {#if result.fallbackPath}
              <p class="text-xs text-slate-300">{t("updater.fallback_saved")}</p>
              <div class="p-2 rounded-lg bg-black/40 font-mono text-[11px] select-all break-all">
                sudo apt install "{result.fallbackPath}"
              </div>
            {/if}
          </div>
        {/if}
      {/if}

      {#if phase === "error"}
        <div class="p-3 rounded-xl bg-red-500/10 border border-red-500/20 text-red-500 text-sm whitespace-pre-wrap">
          {errorMsg}
        </div>
      {/if}
    {/if}
  </div>
</FormModal>
