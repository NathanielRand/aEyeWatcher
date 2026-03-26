import type { ScanResult, ScanEvent, OSInfo, GPUInfo, Theme } from '$lib/types';

const CORE_URL = 'http://localhost:8080';
const CORE_WS = 'ws://localhost:8080';

const _state = $state({
  scanResult: null as ScanResult | null,
  scanEvents: [] as ScanEvent[],
  scanPercent: 0,
  scanning: false,
  scanPhase: '',
  osInfo: null as OSInfo | null,
  gpuInfo: null as GPUInfo | null,
  currentTheme: 'obsidian' as Theme,
  activeSection: 'dashboard',
  sidebarOpen: true,
});

export const store = _state;

export function setTheme(t: Theme) {
  _state.currentTheme = t;
  if (typeof document !== 'undefined') {
    document.documentElement.setAttribute('data-theme', t);
    localStorage.setItem('aeye-theme', t);
  }
}

export function setSection(s: string) {
  _state.activeSection = s;
}

export function toggleSidebar() {
  _state.sidebarOpen = !_state.sidebarOpen;
}

export async function fetchOSInfo() {
  try {
    const res = await fetch(`${CORE_URL}/api/os`);
    if (res.ok) _state.osInfo = await res.json();
  } catch {
    _state.osInfo = null;
  }
}

export async function pollGPU() {
  try {
    const res = await fetch(`${CORE_URL}/api/gpu`);
    if (res.ok) _state.gpuInfo = await res.json();
  } catch { /* silent */ }
}

export function startScan(): Promise<void> {
  return new Promise((resolve) => {
    if (_state.scanning) return;
    _state.scanning = true;
    _state.scanEvents = [];
    _state.scanPercent = 0;
    _state.scanResult = null;

    const ws = new WebSocket(`${CORE_WS}/ws/scan`);

    ws.onmessage = (e) => {
      try {
        const evt: ScanEvent = JSON.parse(e.data);
        _state.scanEvents = [..._state.scanEvents, evt];
        _state.scanPercent = evt.percent;
        _state.scanPhase = evt.phase;

        if (evt.type === 'complete' && evt.detail) {
          _state.scanResult = evt.detail as ScanResult;
          _state.scanning = false;
          resolve();
        }
        if (evt.type === 'error') {
          _state.scanning = false;
          resolve();
        }
      } catch { /* ignore parse errors */ }
    };

    ws.onerror = () => { _state.scanning = false; resolve(); };
    ws.onclose = () => { if (_state.scanning) { _state.scanning = false; resolve(); } };
  });
}

export async function requestDelete(filePath: string): Promise<boolean> {
  const res = await fetch(`${CORE_URL}/api/delete`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ file_path: filePath, confirm: true }),
  });
  return res.ok;
}

export const THEMES: Record<Theme, { label: string; accent: string; bg: string }> = {
  obsidian: { label: 'Obsidian', accent: '#c084fc', bg: '#0a0a0f' },
  clay:     { label: 'Clay',     accent: '#fb923c', bg: '#0f0a08' },
  forest:   { label: 'Forest',   accent: '#4ade80', bg: '#080f09' },
  slate:    { label: 'Slate',    accent: '#7dd3fc', bg: '#080c12' },
  amber:    { label: 'Amber',    accent: '#fbbf24', bg: '#0f0c05' },
};

export function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`;
}

export function formatGB(gb: number): string {
  if (gb >= 1024) return `${(gb / 1024).toFixed(1)} TB`;
  return `${gb.toFixed(1)} GB`;
}
