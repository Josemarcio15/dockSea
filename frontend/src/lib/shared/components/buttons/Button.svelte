<script lang="ts">
  import type { Snippet } from "svelte";

  export type ButtonSize = "xs" | "sm" | "md" | "lg";
  export type ButtonVariant =
    | "primary"
    | "success"
    | "danger"
    | "warning"
    | "neutral"
    | "ghost";

  let {
    type = "button",
    size = "md",
    variant = "ghost",
    disabled = false,
    loading = false,
    title = "",
    class: customClass = "",
    onclick,
    children,
    icon,
  }: {
    type?: "button" | "submit" | "reset";
    size?: ButtonSize;
    variant?: ButtonVariant;
    disabled?: boolean;
    loading?: boolean;
    title?: string;
    class?: string;
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

  const variantClasses: Record<ButtonVariant, string> = {
    primary:
      "bg-violet-600 hover:bg-violet-700 active:bg-violet-800 text-white border-transparent shadow-md shadow-violet-500/20",
    success:
      "bg-emerald-600 hover:bg-emerald-700 active:bg-emerald-800 text-white border-transparent shadow-md shadow-emerald-500/20",
    danger:
      "bg-red-500 hover:bg-red-600 active:bg-red-700 text-white border-transparent shadow-md shadow-red-500/20",
    warning:
      "bg-amber-500 hover:bg-amber-600 active:bg-amber-700 text-white border-transparent shadow-md shadow-amber-500/20",
    neutral:
      "bg-slate-200 hover:bg-slate-300 active:bg-slate-400 text-slate-700 border-slate-300 shadow-sm dark:bg-slate-800 dark:hover:bg-slate-700 dark:active:bg-slate-600 dark:text-slate-200 dark:border-slate-700",
    ghost:
      "bg-transparent hover:bg-slate-100 dark:hover:bg-white/10 text-slate-600 dark:text-slate-300 border-transparent",
  };
</script>

<button
  {type}
  {title}
  disabled={disabled || loading}
  {onclick}
  class="inline-flex items-center justify-center transition-all duration-150 select-none cursor-pointer border font-sans disabled:opacity-50 disabled:cursor-not-allowed disabled:shadow-none {sizeClasses[
    size
  ]} {variantClasses[variant]} {customClass}"
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
