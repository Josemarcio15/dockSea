import type {
  PortMapping,
  EnvVar,
  VolumeMapping,
} from "./types";
import { parseProfileToState, parseJsonToState, buildSavePayload } from "./utils";

export function createConfigModalState(
  getSavedConfigs: () => any[],
  getImage: () => string,
  onSaveProfileCallback: (profile: any) => void,
) {
  let containerName = $state("");
  let projectName = $state("");
  let ports = $state<PortMapping[]>([]);
  let envs = $state<EnvVar[]>([]);
  let volumes = $state<VolumeMapping[]>([]);
  let network = $state("");
  let restartPolicy = $state("");
  let commands = $state<string[]>([]);
  let description = $state("");
  let profileName = $state("");

  let loadedProfileId = $state<string | null>(null);
  let loadedProfileName = $state("");
  let isModified = $state(false);
  let jsonEditor = $state("");
  let jsonMessage = $state("");
  let jsonIsValid = $state(false);
  let editorTab = $state<"fields" | "json">("fields");

  function getEditorConfig() {
    return {
      name: profileName,
      image: getImage(),
      containerName,
      projectName,
      ports: ports.length
        ? ports.map((p) => ({ port: p.external, "port-intern": p.internal }))
        : [{ port: "", "port-intern": "" }],
      envs: envs.length ? envs : [{ name: "", value: "" }],
      volumes: volumes.length ? volumes : [{ host: "", container: "" }],
      network,
      restartPolicy,
      commands: commands.length ? commands : [""],
      description,
    };
  }

  function syncJsonEditor() {
    jsonEditor = JSON.stringify(getEditorConfig(), null, 2);
    jsonMessage = "";
  }

  async function copyJson() {
    await navigator.clipboard.writeText(jsonEditor);
    jsonMessage = "JSON copiado.";
  }

  function applyJson() {
    const { valid, state } = parseJsonToState(jsonEditor);
    if (!valid) {
      jsonIsValid = false;
      jsonMessage = "JSON inválido. Verifique a estrutura antes de aplicar.";
      return;
    }
    containerName = state.containerName || "";
    projectName = state.projectName || "";
    network = state.network || "";
    restartPolicy = state.restartPolicy || "";
    description = state.description || "";
    profileName = state.profileName || profileName;
    ports = state.ports || [];
    envs = state.envs || [];
    volumes = state.volumes || [];
    commands = state.commands || [""];

    const match = getSavedConfigs().find((s) => s.name === profileName);
    loadedProfileId = match?.id || null;
    loadedProfileName = match?.name || "";
    isModified = true;
    jsonIsValid = true;
    jsonMessage = "JSON aplicado ao formulário.";
  }

  function handleJsonChange(value: string) {
    jsonEditor = value;
    isModified = true;
    try {
      JSON.parse(value);
      applyJson();
    } catch {}
  }

  function loadExample() {
    jsonEditor = JSON.stringify(
      {
        name: "meu-perfil",
        image: "nginx:latest",
        containerName: "meu-container",
        projectName: "meu-projeto",
        ports: [{ port: "8080", "port-intern": "80" }],
        envs: [{ name: "APP_ENV", value: "production" }],
        volumes: [{ host: "./data", container: "/app/data" }],
        network: "app-network",
        restartPolicy: "unless-stopped",
        commands: ["nginx -g 'daemon off;'"],
        description: "Exemplo de configuração Docker",
      },
      null,
      2,
    );
    applyJson();
    editorTab = "json";
  }

  let hasEmptyVolume = $derived(
    volumes.some(
      (v) =>
        (v.host.trim() && !v.container.trim()) ||
        (!v.host.trim() && v.container.trim()),
    ),
  );

  let nameAlreadyExists = $derived(
    profileName.trim() !== "" &&
      getSavedConfigs().some((c) => c.name === profileName && c.id !== loadedProfileId),
  );

  let isNameChanged = $derived(
    loadedProfileId !== null && profileName !== loadedProfileName,
  );

  let saveDisabled = $derived(
    profileName.trim() === "" || nameAlreadyExists || hasEmptyVolume,
  );

  function loadProfile(cfg: any) {
    const s = parseProfileToState(cfg);
    containerName = s.containerName || "";
    projectName = s.projectName || "";
    network = s.network || "";
    restartPolicy = s.restartPolicy || "no";
    description = s.description || "";
    profileName = s.profileName || "";
    loadedProfileId = cfg.id;
    loadedProfileName = cfg.name || "";
    commands = s.commands || [""];
    ports = s.ports || [];
    envs = s.envs || [];
    volumes = s.volumes || [];
    isModified = false;
    syncJsonEditor();
  }

  function createNewProfile() {
    loadedProfileId = null;
    loadedProfileName = "";
    profileName = "";
    containerName = "";
    projectName = "";
    ports = [];
    envs = [];
    volumes = [];
    network = "";
    restartPolicy = "";
    commands = [""];
    description = "";
    jsonEditor = "";
    jsonMessage = "Novo perfil pronto para edição.";
    isModified = true;
    editorTab = "fields";
  }

  function triggerSave() {
    const payload = buildSavePayload({
      profileName,
      containerName,
      projectName,
      image: getImage(),
      ports,
      envs,
      volumes,
      network,
      restartPolicy,
      commands,
      description,
      loadedProfileId,
      isNameChanged,
    });
    onSaveProfileCallback(payload);
    isModified = false;
  }

  function saveJsonProfile() {
    applyJson();
    if (!jsonIsValid) return;
    if (saveDisabled) {
      jsonMessage = "Informe um nome de perfil válido antes de salvar.";
      return;
    }
    triggerSave();
    jsonMessage = "Perfil salvo no SQLite.";
  }

  function resetState() {
    loadedProfileId = null;
    loadedProfileName = "";
    isModified = false;
    containerName = "";
    projectName = "";
    ports = [];
    envs = [];
    volumes = [];
    network = "";
    restartPolicy = "";
    commands = [];
    description = "";
    profileName = "";
    jsonEditor = "";
    jsonMessage = "";
    editorTab = "fields";
  }

  function initForImage(img: string) {
    if (containerName === "") {
      const shortName = img.split(":")[0].split("/").pop() || "app";
      containerName = `${shortName}-container`;
    }
    if (commands.length === 0) {
      commands = [""];
    }
  }

  return {
    get containerName() { return containerName; },
    set containerName(v) { containerName = v; },
    get projectName() { return projectName; },
    set projectName(v) { projectName = v; },
    get ports() { return ports; },
    set ports(v) { ports = v; },
    get envs() { return envs; },
    set envs(v) { envs = v; },
    get volumes() { return volumes; },
    set volumes(v) { volumes = v; },
    get network() { return network; },
    set network(v) { network = v; },
    get restartPolicy() { return restartPolicy; },
    set restartPolicy(v) { restartPolicy = v; },
    get commands() { return commands; },
    set commands(v) { commands = v; },
    get description() { return description; },
    set description(v) { description = v; },
    get profileName() { return profileName; },
    set profileName(v) { profileName = v; },
    get loadedProfileId() { return loadedProfileId; },
    get loadedProfileName() { return loadedProfileName; },
    get isModified() { return isModified; },
    set isModified(v) { isModified = v; },
    get jsonEditor() { return jsonEditor; },
    set jsonEditor(v) { jsonEditor = v; },
    get jsonMessage() { return jsonMessage; },
    get jsonIsValid() { return jsonIsValid; },
    get editorTab() { return editorTab; },
    set editorTab(v) { editorTab = v; },
    get hasEmptyVolume() { return hasEmptyVolume; },
    get nameAlreadyExists() { return nameAlreadyExists; },
    get isNameChanged() { return isNameChanged; },
    get saveDisabled() { return saveDisabled; },
    syncJsonEditor,
    copyJson,
    applyJson,
    handleJsonChange,
    loadExample,
    loadProfile,
    createNewProfile,
    triggerSave,
    saveJsonProfile,
    resetState,
    initForImage,
  };
}
