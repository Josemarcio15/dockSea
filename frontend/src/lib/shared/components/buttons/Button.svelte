<script lang="ts" module>
  export type ButtonSize = "xs" | "sm" | "md" | "lg";
</script>

<script lang="ts">
  import type { Snippet } from "svelte";
  import { themeStore } from "$shared/theme/theme.svelte.ts";

  let {
    type = "button",
    size,
    color,
    themeKey, // Ex: "stacks.deploy_btn"
    disabled = false,
    loading = false,
    title = "",
    class: customClass = "",
    style: customStyle = "",
    onclick,
    children,
    icon,
  }: {
    type?: "button" | "submit" | "reset";
    size?: ButtonSize;
    color?: string; // HEX direto
    themeKey?: string; // "rota.acao" ex: "stacks.deploy_btn"
    disabled?: boolean;
    loading?: boolean;
    title?: string;
    class?: string;
    style?: string;
    onclick?: (e: MouseEvent) => unknown;
    children?: Snippet;
    icon?: Snippet;
  } = $props();

  const sizeClasses: Record<ButtonSize, string> = {
    xs: "px-2.5 py-1 text-[11px] rounded-lg gap-1.5",
    sm: "px-3.5 py-1.5 text-xs rounded-xl gap-1.5 font-semibold",
    md: "px-4 py-2 text-xs rounded-xl gap-2 font-bold",
    lg: "px-5 py-2.5 text-sm rounded-xl gap-2.5 font-bold",
  };

  // Buscar dados do tema por rota/chave de forma reativa direta
  let routeConfig = $derived.by(() => {
    if (!themeKey) return undefined;
    const parts = themeKey.split(".");
    if (parts.length === 2) {
      const [route, btnKey] = parts;
      return themeStore.currentTheme.routes?.[route]?.[btnKey];
    }
    return undefined;
  });

  // Determinar tamanho final (prop > tema > default "md")
  let activeSize = $derived<ButtonSize>(size || routeConfig?.size || "md");

  // Cor de fundo calculada
  let activeBg = $derived.by(() => {
    if (routeConfig?.color) return routeConfig.color;
    if (color) return color;
    return "#2563eb";
  });

  // Cor do texto
  let activeTextColor = $derived.by(() => {
    if (routeConfig?.text) return routeConfig.text;
    return "#ffffff";
  });

  // Estilos inline de cores dinâmicas
  let dynamicStyle = $derived.by(() => {
    const bg = routeConfig?.bg || (routeConfig as any)?.color || color || "#2563eb";
    const hover = routeConfig?.hover || bg;
    const text = routeConfig?.text || "#ffffff";
    const textHover = routeConfig?.textHover || text;

    let styles: string[] = [];
    styles.push(`--btn-bg: ${bg}`);
    styles.push(`--btn-hover: ${hover}`);
    styles.push(`--btn-text: ${text}`);
    styles.push(`--btn-text-hover: ${textHover}`);
    styles.push(`background-color: var(--btn-bg)`);
    styles.push(`color: var(--btn-text)`);
    styles.push(`border-color: transparent`);
    styles.push(`box-shadow: 0 4px 14px ${bg}33`);

    if (customStyle) {
      styles.push(customStyle);
    }
    return styles.join("; ");
  });
</script>

<button
  {type}
  {title}
  disabled={disabled || loading}
  {onclick}
  style={dynamicStyle}
  class="ds-btn inline-flex items-center justify-center transition-all duration-150 select-none cursor-pointer border font-sans active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed disabled:shadow-none disabled:active:scale-100 {sizeClasses[
    activeSize
  ]} {customClass}"
>
  {#if loading}
    <div
      class="w-3.5 h-3.5 border-2 border-current border-t-transparent rounded-full animate-spin shrink-0"
    ></div>
  {:else if icon}
    {@render icon()}
  {/if}

  {#if children}
    {@render children()}
  {/if}
</button>
