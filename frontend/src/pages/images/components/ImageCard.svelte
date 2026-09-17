<script lang="ts">
  import { t } from "$shared/stores/locale.svelte";
  import type { DockerImage } from "$lib/domains/images";
  import { Button } from "$shared/components/buttons";

  let {
    img,
    checked = false,
    on_toggle = () => {},
    on_build = () => {},
  }: {
    img: DockerImage;
    checked?: boolean;
    on_toggle?: () => void;
    on_build?: () => void;
  } = $props();

  let expanded = $state(false);

  const isInUse = $derived(
    img.containersUsing && img.containersUsing.length > 0,
  );
  const containerCount = $derived(
    img.containersUsing ? img.containersUsing.length : 0,
  );
  const repoLower = $derived(img.repo.toLowerCase());
  const containersHint = $derived(
    (img.containersUsing || []).join(" ").toLowerCase(),
  );

  const iconUrl = $derived.by(() => {
    const searchTarget =
      repoLower === "<none>" ? `${repoLower} ${containersHint}` : repoLower;

    if (searchTarget.includes("wordpress"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/wordpress/wordpress-plain.svg";
    if (searchTarget.includes("postgres") || searchTarget.includes("postgre"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/postgresql/postgresql-original.svg";
    if (searchTarget.includes("mariadb"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mariadb/mariadb-original.svg";
    if (searchTarget.includes("mysql"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mysql/mysql-original.svg";
    if (searchTarget.includes("mongo"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mongodb/mongodb-original.svg";
    if (searchTarget.includes("redis"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/redis/redis-original.svg";
    if (searchTarget.includes("debian"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/debian/debian-original.svg";
    if (searchTarget.includes("ubuntu"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/ubuntu/ubuntu-plain.svg";
    if (searchTarget.includes("alpine"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/alpinejs/alpinejs-original.svg";
    if (searchTarget.includes("nginx"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nginx/nginx-original.svg";
    if (searchTarget.includes("apache") || searchTarget.includes("httpd"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/apache/apache-original.svg";
    if (searchTarget.includes("node"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nodejs/nodejs-original.svg";
    if (searchTarget.includes("python"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/python/python-original.svg";
    if (searchTarget.includes("golang") || searchTarget.includes("go:"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg";
    if (searchTarget.includes("rust"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/rust/rust-plain.svg";
    if (searchTarget.includes("php"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/php/php-original.svg";
    if (searchTarget.includes("traefik"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/traefik/traefik-original.svg";
    if (searchTarget.includes("grafana"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/grafana/grafana-original.svg";
    if (searchTarget.includes("prometheus"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/prometheus/prometheus-original.svg";
    if (searchTarget.includes("rabbitmq"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/rabbitmq/rabbitmq-original.svg";
    if (searchTarget.includes("elasticsearch") || searchTarget.includes("elastic"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/elasticsearch/elasticsearch-original.svg";
    if (searchTarget.includes("docker"))
      return "https://cdn.jsdelivr.net/gh/devicons/devicon/icons/docker/docker-original.svg";
    return null;
  });

  const createdStr = $derived.by(() => {
    if (!img.created) return "";
    const date = new Date(img.created * 1000);
    const day = String(date.getDate()).padStart(2, "0");
    const month = String(date.getMonth() + 1).padStart(2, "0");
    const year = date.getFullYear();
    return `${day}/${month}/${year}`;
  });
</script>

<div
  class="relative rounded-[22px] bg-gradient-to-b from-slate-900/90 to-[#090d16]/95 dark:from-[#111726]/90 dark:to-[#070b13]/95 border border-white/10 dark:border-white/5 backdrop-blur-xl shadow-[0_12px_32px_rgba(0,0,0,0.45)] hover:border-violet-500/40 hover:shadow-[0_16px_36px_rgba(99,102,241,0.15)] transition-all duration-300 flex flex-col p-4 gap-3.5 group text-slate-100 overflow-hidden w-full"
>
  <!-- Glow highlight top -->
  <div class="absolute -top-12 -left-12 w-32 h-32 bg-violet-600/15 rounded-full blur-2xl pointer-events-none"></div>

  <!-- Header Section (Glassmorphism Pill) -->
  <div
    class="relative flex items-center justify-between p-3.5 rounded-2xl bg-white/[0.04] dark:bg-white/[0.03] border border-white/[0.08] backdrop-blur-md shadow-inner"
  >
    <div class="flex items-center gap-3.5 min-w-0 flex-1">
      <!-- Checkbox Customizado -->
      <button
        type="button"
        class="w-5 h-5 rounded-lg border flex items-center justify-center cursor-pointer transition-all duration-200 shrink-0 {checked
          ? 'bg-gradient-to-tr from-violet-600 to-indigo-500 border-violet-400 text-white shadow-[0_0_10px_rgba(124,58,237,0.5)]'
          : 'border-white/20 bg-white/5 hover:border-violet-400'}"
        onclick={on_toggle}
        aria-label="Selecionar imagem"
      >
        {#if checked}
          <span class="text-white text-[11px] font-extrabold leading-none">✓</span>
        {/if}
      </button>

      <!-- Logo / Ícone Container -->
      <div
        class="w-11 h-11 rounded-xl bg-violet-600/10 border border-violet-500/20 p-2 flex items-center justify-center shrink-0"
      >
        {#if iconUrl}
          <img
            src={iconUrl}
            alt={img.repo}
            class="w-full h-full object-contain"
          />
        {:else}
          <span class="text-sm font-black text-violet-300">IMG</span>
        {/if}
      </div>

      <!-- Título + Tag + Status Compacto -->
      <div class="flex flex-col min-w-0 flex-1">
        <h3
          class="font-extrabold text-base text-white tracking-tight truncate drop-shadow-sm"
          title={img.repo}
        >
          {#if !img.repo || img.repo === "<none>"}
            <span class="text-slate-400 italic font-normal">&lt;sem tag&gt;</span>
          {:else}
            {img.repo}
          {/if}
        </h3>

        <div class="flex items-center gap-2 mt-1.5 flex-wrap">
          <!-- Tag Badge -->
          <span
            class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-violet-500/15 border border-violet-400/30 text-violet-300 font-mono shadow-xs truncate max-w-[120px]"
            title={img.tag}
          >
            {img.tag || "latest"}
          </span>

          <!-- Status Indicator Badge -->
          <span
            class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-semibold {isInUse
              ? 'bg-emerald-500/15 border border-emerald-400/30 text-emerald-300'
              : 'bg-slate-500/15 border border-slate-400/30 text-slate-300'} shadow-xs"
          >
            <span
              class="w-1.5 h-1.5 rounded-full {isInUse
                ? 'bg-emerald-400 shadow-[0_0_6px_#34d399]'
                : 'bg-slate-400'} shrink-0"
            ></span>
            {isInUse ? `${containerCount} em uso` : 'Livre'}
          </span>
        </div>
      </div>
    </div>

    <!-- Botão Expandir / Recolher Detalhes -->
    <button
      type="button"
      class="w-9 h-9 rounded-xl bg-white/[0.06] hover:bg-white/[0.12] active:scale-95 border border-white/10 text-white/80 hover:text-white flex items-center justify-center cursor-pointer transition-all duration-200 ml-2 shrink-0 shadow-sm"
      onclick={() => (expanded = !expanded)}
      title="Mais detalhes"
      aria-label="Expandir detalhes"
    >
      <svg
        class="w-4 h-4 transition-transform duration-300 {expanded ? 'rotate-0' : 'rotate-180'}"
        fill="currentColor"
        viewBox="0 0 24 24"
      >
        <path d="M12 8.586l6.293 6.293 1.414-1.414L12 5.758 4.293 13.465l1.414 1.414z"/>
      </svg>
    </button>
  </div>

  <!-- Detalhes Extras (visíveis apenas quando expandido) -->
  {#if expanded}
    <div class="flex flex-col gap-3 animate-in fade-in slide-in-from-top-2 duration-300">
      <!-- ID & Tamanho Grid -->
      <div class="grid grid-cols-2 gap-3">
        <!-- ID Card -->
        <div
          class="flex items-center gap-3 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
        >
          <div
            class="w-10 h-10 rounded-xl bg-blue-500/15 border border-blue-500/30 flex items-center justify-center text-blue-300 shrink-0 font-bold text-xs"
          >
            ID
          </div>
          <div class="flex flex-col min-w-0">
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t("images.card_image_id")}</span>
            <span class="font-mono font-extrabold text-sm text-blue-300 truncate">{img.id.substring(0, 12)}</span>
          </div>
        </div>

        <!-- Tamanho Card -->
        <div
          class="flex items-center gap-3 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
        >
          <div
            class="w-10 h-10 rounded-xl bg-violet-500/15 border border-violet-500/30 flex items-center justify-center text-violet-300 shrink-0"
          >
            <svg class="w-5 h-5 fill-current" viewBox="0 0 24 24">
              <path d="M12 2C6.48 2 2 4.02 2 6.5s4.48 4.5 10 4.5 10-2.02 10-4.5S17.52 2 12 2zm0 6c-3.86 0-7-1.12-7-2.5S8.14 3 12 3s7 1.12 7 2.5S15.86 8 12 8zm0 4.5c-5.52 0-10-2.02-10-4.5v3c0 2.48 4.48 4.5 10 4.5s10-2.02 10-4.5v-3c0 2.48-4.48 4.5-10 4.5zm0 6c-5.52 0-10-2.02-10-4.5v3c0 2.48 4.48 4.5 10 4.5s10-2.02 10-4.5v-3c0 2.48-4.48 4.5-10 4.5z"/>
            </svg>
          </div>
          <div class="flex flex-col min-w-0">
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t("images.card_size")}</span>
            <span class="font-extrabold text-sm text-white truncate">{img.size}</span>
          </div>
        </div>
      </div>

      <!-- Data de Criação -->
      <div
        class="flex items-center justify-between p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
      >
        <div class="flex items-center gap-3">
          <div
            class="w-10 h-10 rounded-xl bg-purple-500/15 border border-purple-500/30 flex items-center justify-center text-purple-300 shrink-0"
          >
            <svg class="w-5 h-5 fill-current" viewBox="0 0 24 24">
              <path d="M19 4h-1V2h-2v2H8V2H6v2H5c-1.11 0-1.99.9-1.99 2L3 20a2 2 0 0 0 2 2h14c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 16H5V9h14v11zM7 11h5v5H7z"/>
            </svg>
          </div>
          <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{t("images.card_created")}</span>
        </div>
        <span class="font-extrabold text-sm text-white tracking-wide">{createdStr || "—"}</span>
      </div>

      <!-- Containers Usando -->
      {#if img.containersUsing && img.containersUsing.length > 0}
        <div
          class="flex flex-col gap-2 p-3.5 rounded-2xl bg-white/[0.03] border border-white/[0.06] shadow-sm"
        >
          <span class="text-[10px] font-bold text-purple-300 uppercase tracking-wider">
            {t("images.card_containers_using")}
          </span>
          <div class="flex flex-wrap gap-1.5">
            {#each img.containersUsing as containerName}
              <span
                class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-xl bg-purple-500/10 border border-purple-500/20 text-xs font-semibold text-purple-200"
              >
                {containerName}
              </span>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  {/if}

  <!-- Ação Sempre Visível (mesmo recolhido) -->
  <div class="pt-2 border-t border-white/[0.08] flex justify-center">
    <Button
      size="sm"
      themeKey="images.create_container_btn"
      class="w-full"
      onclick={on_build}
    >
      {t("images.create_container_btn")}
    </Button>
  </div>
</div>
