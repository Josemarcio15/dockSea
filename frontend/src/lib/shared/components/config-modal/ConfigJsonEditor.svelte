<script lang="ts">
  import CodeEditor from "$shared/components/CodeEditor.svelte";
  import { Button } from "$shared/components/buttons";

  let {
    jsonEditor = $bindable(""),
    jsonMessage = "",
    jsonIsValid = false,
    saveDisabled = false,
    onLoadExample,
    onCopyJson,
    onSaveProfile,
    onJsonChange,
  }: {
    jsonEditor: string;
    jsonMessage: string;
    jsonIsValid: boolean;
    saveDisabled: boolean;
    onLoadExample: () => void;
    onCopyJson: () => void;
    onSaveProfile: () => void;
    onJsonChange: (val: string) => void;
  } = $props();
</script>

<div
  class="space-y-2 rounded-xl border border-violet-200 dark:border-violet-900/50 bg-violet-50/40 dark:bg-violet-950/10 p-3"
>
  <div class="flex items-center justify-between gap-2">
    <div>
      <span
        class="text-[10px] font-bold uppercase tracking-wider text-violet-700 dark:text-violet-300"
        >JSON do perfil</span
      >
      <p class="text-[10px] text-slate-500">
        Selecione um perfil acima, depois edite ou cole as variáveis aqui.
      </p>
    </div>
    <div class="flex gap-2 shrink-0">
      <Button variant="primary" size="xs" onclick={onLoadExample}>
        Exemplo
      </Button>
      <Button variant="primary" size="xs" onclick={onCopyJson}>
        Copiar
      </Button>
      <Button
        variant="success"
        size="xs"
        disabled={!jsonIsValid || saveDisabled}
        onclick={onSaveProfile}
      >
        Salvar
      </Button>
    </div>
  </div>
  <CodeEditor
    value={jsonEditor}
    mode="json"
    onchange={onJsonChange}
  />
  {#if jsonMessage}
    <p
      class="text-[10px] font-semibold {jsonMessage.includes('inválido')
        ? 'text-red-500'
        : 'text-emerald-600'}"
    >
      {jsonMessage}
    </p>
  {/if}
</div>
