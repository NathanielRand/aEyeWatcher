// Mirrors Go model types from aeye_watcher_core/internal/models/types.go

export type ComponentStatus = 'Running' | 'Installed' | 'Not Found' | 'Degraded';

export interface Component {
  name: string;
  type: string;
  status: ComponentStatus;
  path?: string;
  port?: string;
  version?: string;
  description?: string;
  url?: string;
  last_seen?: string;
  metadata?: Record<string, string>;
}

export interface AIModel {
  name: string;
  tag: string;
  size_bytes: number;
  provider: string;
  file_path: string;
  extension?: string;
  is_orphaned: boolean;
  format?: string;
}

export interface AgentFile {
  type: string;
  name: string;
  path: string;
  size_bytes: number;
}

export interface ExpectedFile {
  label: string;
  pattern: string;
  found: boolean;
  path?: string;
}

export interface AgentCard {
  name: string;
  ecosystem: string;
  status: 'detected' | 'partial' | 'not_found';
  description: string;
  icon: string;
  files: AgentFile[];
  expected_files: ExpectedFile[];
}

export interface MCPServer {
  name: string;
  url?: string;
  config?: string;
  status: string;
  transport: string;
  tools?: string[];
}

export interface ACPAgent {
  name: string;
  endpoint: string;
  status: string;
  protocol: string;
}

export interface DiskInfo {
  mount: string;
  device: string;
  fs_type: string;
  total_gb: number;
  used_gb: number;
  free_gb: number;
  use_pct: number;
}

export interface GPUInfo {
  name?: string;
  used_mb: number;
  total_mb: number;
  use_pct: number;
  temp_c?: number;
  error?: string;
  backend: string;
}

export interface ScanResult {
  components: Component[];
  models: AIModel[];
  agent_cards: AgentCard[];
  mcp_servers: MCPServer[];
  acp_agents: ACPAgent[];
  model_dirs: string[];
  disks: DiskInfo[];
  gpu: GPUInfo;
  scanned_at: string;
}

export type ScanEventType = 'checkpoint' | 'progress' | 'found' | 'complete' | 'error';

export interface ScanEvent {
  type: ScanEventType;
  phase: string;
  message: string;
  percent: number;
  detail?: unknown;
}

export interface OSInfo {
  os: string;
  arch: string;
  distro?: string;
  version?: string;
  hostname?: string;
  is_container: boolean;
  gpu_available: boolean;
}

export type Theme = 'obsidian' | 'clay' | 'forest' | 'slate' | 'amber';
