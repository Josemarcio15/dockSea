export interface PortMapping {
  external: string;
  internal: string;
}

export interface EnvVar {
  name: string;
  value: string;
}

export interface VolumeMapping {
  host: string;
  container: string;
}

export interface ConfigFormState {
  image: string;
  containerName: string;
  projectName: string;
  ports: PortMapping[];
  envs: EnvVar[];
  volumes: VolumeMapping[];
  network: string;
  restartPolicy: string;
  commands: string[];
  description: string;
  profileName: string;
}

export const emptyConfigTemplate = {
  name: "",
  image: "",
  containerName: "",
  projectName: "",
  ports: [{ port: "", "port-intern": "" }],
  envs: [{ name: "", value: "" }],
  volumes: [{ host: "", container: "" }],
  network: "",
  restartPolicy: "",
  commands: [""],
  description: "",
};

export function getDefaultContainerPath(imgName: string): string {
  const img = imgName.split(":")[0].split("/").pop()?.toLowerCase() || "";
  if (img.includes("postgres")) return "/var/lib/postgresql/data";
  if (img.includes("mysql") || img.includes("mariadb")) return "/var/lib/mysql";
  if (img.includes("mongo")) return "/data/db";
  if (img.includes("redis")) return "/data";
  if (img.includes("influx")) return "/var/lib/influxdb";
  if (img.includes("rabbitmq")) return "/var/lib/rabbitmq";
  if (img.includes("elasticsearch")) return "/usr/share/elasticsearch/data";
  if (img.includes("nginx")) return "/usr/share/nginx/html";
  if (img.includes("httpd")) return "/usr/local/apache2/htdocs";
  return "/data";
}
