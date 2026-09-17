import type { PortMapping, EnvVar, VolumeMapping, ConfigFormState } from "./types";
import { emptyConfigTemplate, getDefaultContainerPath } from "./types";

export function parseProfileToState(cfg: any): Partial<ConfigFormState> {
  const commands = cfg.command ? cfg.command.split(" && ") : [""];

  // Parse ports "host:container host2:container2"
  const ports: PortMapping[] = [];
  if (cfg.ports) {
    for (const entry of cfg.ports.split(/\s+/)) {
      const parts = entry.split(":");
      if (parts.length === 2) {
        ports.push({ external: parts[0], internal: parts[1] });
      }
    }
  }

  // Parse ENVs
  const envs: EnvVar[] = [];
  if (cfg.env) {
    const rawEnv = cfg.env.trim();
    const hasNewline = rawEnv.includes("\n");
    const tokens = rawEnv.split(/\s+/);
    const isLegacySpaceSeparated =
      !hasNewline && tokens.length > 1 && tokens.every((t: string) => t.includes("="));

    if (isLegacySpaceSeparated) {
      for (const entry of tokens) {
        const idx = entry.indexOf("=");
        if (idx > 0) envs.push({ name: entry.substring(0, idx), value: entry.substring(idx + 1) });
      }
    } else {
      for (const line of rawEnv.split(/\r?\n/)) {
        const trimmed = line.trim();
        if (!trimmed) continue;
        const idx = trimmed.indexOf("=");
        if (idx > 0) envs.push({ name: trimmed.substring(0, idx).trim(), value: trimmed.substring(idx + 1) });
      }
    }
  }

  // Parse volumes "host:container host2:container2"
  const volumes: VolumeMapping[] = [];
  if (cfg.volumes) {
    for (const entry of cfg.volumes.split(/\s+/)) {
      const parts = entry.split(":");
      if (parts.length === 2) {
        volumes.push({ host: parts[0], container: parts[1] });
      } else if (entry) {
        volumes.push({ host: entry, container: "" });
      }
    }
  }

  return {
    containerName: cfg.containerName || "",
    projectName: cfg.projectName || "",
    network: cfg.network || "",
    restartPolicy: cfg.restartPolicy || "no",
    description: cfg.description || "",
    profileName: cfg.name || "",
    commands,
    ports,
    envs,
    volumes,
  };
}

export function parseJsonToState(jsonStr: string): { state: Partial<ConfigFormState>; valid: boolean } {
  try {
    const parsed = JSON.parse(jsonStr);
    const cfg = { ...emptyConfigTemplate, ...(parsed.config || parsed) };
    return {
      valid: true,
      state: {
        containerName: cfg.containerName || "",
        projectName: cfg.projectName || "",
        image: cfg.image || "",
        network: cfg.network || "",
        restartPolicy: cfg.restartPolicy || "",
        description: cfg.description || "",
        profileName: cfg.name || "",
        ports: (Array.isArray(cfg.ports) ? cfg.ports : []).map((p: any) => ({
          external: p.port ?? p.external ?? "",
          internal: p["port-intern"] ?? p.internal ?? "",
        })),
        envs: (Array.isArray(cfg.envs) ? cfg.envs : []).map((e: any) => ({
          name: e.name ?? "",
          value: e.value ?? "",
        })),
        volumes: (Array.isArray(cfg.volumes) ? cfg.volumes : []).map((v: any) => ({
          host: v.host ?? "",
          container: v.container ?? "",
        })),
        commands: Array.isArray(cfg.commands)
          ? cfg.commands.map((c: any) => String(c ?? ""))
          : [""],
      },
    };
  } catch {
    return { valid: false, state: {} };
  }
}

export function buildSavePayload(
  state: {
    profileName: string;
    containerName: string;
    projectName: string;
    image: string;
    ports: PortMapping[];
    envs: EnvVar[];
    volumes: VolumeMapping[];
    network: string;
    restartPolicy: string;
    commands: string[];
    description: string;
    loadedProfileId: string | null;
    isNameChanged: boolean;
  },
) {
  const portsStr = state.ports.map((p) => `${p.external}:${p.internal}`).join(" ");
  const envStr = state.envs
    .filter((e) => e.name.trim() !== "")
    .map((e) => `${e.name}=${e.value}`)
    .join("\n");
  const volStr = state.volumes
    .filter((v) => v.host.trim() && v.container.trim())
    .map((v) => `${v.host}:${v.container || getDefaultContainerPath(state.image)}`)
    .join(" ");
  const commandStr = state.commands.filter((c) => c.trim() !== "").join(" && ");

  return {
    id: state.isNameChanged ? null : state.loadedProfileId,
    name: state.profileName,
    containerName: state.containerName,
    projectName: state.projectName,
    image: state.image,
    ports: portsStr || null,
    env: envStr || null,
    volumes: volStr || null,
    network: state.network || null,
    restartPolicy: state.restartPolicy,
    command: commandStr || null,
    description: state.description || null,
  };
}
