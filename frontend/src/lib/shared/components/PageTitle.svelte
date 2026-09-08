<script lang="ts">
  import type { Snippet } from "svelte";

  let {
    title,
    icon,
  }: {
    title: string;
    icon?: Snippet;
  } = $props();
</script>

<div
  class="inline-flex items-center gap-2.5 px-5 py-2.5 rounded-2xl bg-linear-to-br from-violet-100 to-fuchsia-100 dark:from-violet-950/40 dark:to-fuchsia-950/40 border border-violet-200/50 dark:border-violet-800/50 self-start shadow-sm"
>
  {#if icon}
    <span class="text-violet-600 dark:text-fuchsia-300 shrink-0">
      {@render icon()}
    </span>
  {/if}
  <h1
    class="text-2xl font-bold bg-linear-to-r from-violet-600 to-fuchsia-500 bg-clip-text text-transparent m-0 flex items-center gap-2"
  >
    {title}
  </h1>
</div>
