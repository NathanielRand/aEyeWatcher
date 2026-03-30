import type { ScanResult, ScanEvent, OSInfo, GPUInfo, Theme } from "$lib/types";

// ─── Config (loaded from localStorage, drives actual scanner behavior) ────────
export type AppConfig = {
	coreUrl: string;
	coreWs: string;
	scanOnLoad: boolean;
	gpuPollInterval: number;
	deepScanEnabled: boolean;
	deepScanMaxDepth: number;
	safeMode: boolean;
	minModelSizeMB: number;
	ignoreDirs: string;
	extraModelDirs: string;
	extraPorts: string;
	remoteEnabled: boolean;
	remoteHosts: string;
	remoteSshUser: string;
	remoteSshKeyPath: string;
};

const CONFIG_DEFAULTS: AppConfig = {
	coreUrl: "http://localhost:8080",
	coreWs: "ws://localhost:8080",
	scanOnLoad: false,
	gpuPollInterval: 3,
	deepScanEnabled: true,
	deepScanMaxDepth: 12,
	safeMode: true,
	minModelSizeMB: 10,
	ignoreDirs: "node_modules,.git,.npm,Windows,System32",
	extraModelDirs: "",
	extraPorts: "",
	remoteEnabled: false,
	remoteHosts: "",
	remoteSshUser: "",
	remoteSshKeyPath: "~/.ssh/id_rsa",
};

function loadConfig(): AppConfig {
	try {
		const s = localStorage.getItem("aeye-config");
		if (s) return { ...CONFIG_DEFAULTS, ...JSON.parse(s) };
	} catch {
		/* ignore */
	}
	return { ...CONFIG_DEFAULTS };
}

export function saveConfig(cfg: Partial<AppConfig>) {
	const current = _state.config;
	const next = { ...current, ...cfg };
	_state.config = next;
	// Derive WS URL from coreUrl
	_state.config.coreWs = next.coreUrl.replace(/^http/, "ws");
	localStorage.setItem("aeye-config", JSON.stringify(_state.config));
}

// ─── Types ────────────────────────────────────────────────────────────────────
export type Toast = {
	id: string;
	msg: string;
	type: "info" | "success" | "error" | "warn";
	link?: string;
};
export type ScheduleRule = {
	id: string;
	label: string;
	cron: string;
	enabled: boolean;
	lastRun?: string;
	nextRun?: string;
};
export type ScanMeta = {
	startedAt: string;
	completedAt?: string;
	durationMs?: number;
	totalModels: number;
	orphanedModels: number;
	runningServices: number;
	mcpServers: number;
	agentCards: number;
	diskVolumes: number;
	eventsCount: number;
	status: "running" | "complete" | "error";
};
export type ProFeature =
	| "alerts"
	| "export"
	| "remote"
	| "scheduler_advanced"
	| "plugins";
export type Plugin = {
	id: string;
	name: string;
	version: string;
	enabled: boolean;
	description: string;
	author: string;
	repoUrl?: string;
};

// ─── Central reactive state ───────────────────────────────────────────────────
const _state = $state({
	config: loadConfig() as AppConfig,
	scanResult: null as ScanResult | null,
	scanEvents: [] as ScanEvent[],
	scanPercent: 0,
	scanning: false,
	scanPhase: "",
	scanMeta: null as ScanMeta | null,
	osInfo: null as OSInfo | null,
	gpuInfo: null as GPUInfo | null,
	currentTheme: "obsidian" as Theme,
	toasts: [] as Toast[],
	schedules: [] as ScheduleRule[],
	scheduleTimer: null as ReturnType<typeof setInterval> | null,
	// Pro tier: in production this would come from auth/license check
	isPro: false,
	plugins: [] as Plugin[],
});

export const store = _state;

// ─── Pro tier ─────────────────────────────────────────────────────────────────
export function isProFeature(f: ProFeature): boolean {
	return _state.isPro;
}
export function activatePro(licenseKey: string) {
	// In production: validate against license server
	if (licenseKey === "AEYE-PRO-DEMO") {
		_state.isPro = true;
		localStorage.setItem("aeye-pro", "true");
		addToast("Pro tier activated!", "success");
	} else {
		addToast("Invalid license key", "error");
	}
}
export function loadProStatus() {
	if (typeof localStorage !== "undefined") {
		_state.isPro = localStorage.getItem("aeye-pro") === "true";
	}
}

// ─── Theme ────────────────────────────────────────────────────────────────────
export function setTheme(t: Theme) {
	_state.currentTheme = t;
	if (typeof document !== "undefined") {
		document.documentElement.setAttribute("data-theme", t);
		localStorage.setItem("aeye-theme", t);
	}
}

// ─── Toast ────────────────────────────────────────────────────────────────────
export function addToast(
	msg: string,
	type: Toast["type"] = "info",
	link?: string,
	duration = 5000,
) {
	const id = Math.random().toString(36).slice(2);
	_state.toasts = [..._state.toasts, { id, msg, type, link }];
	setTimeout(() => dismissToast(id), duration);
}
export function dismissToast(id: string) {
	_state.toasts = _state.toasts.filter((t) => t.id !== id);
}

// ─── OS / GPU ─────────────────────────────────────────────────────────────────
export async function fetchOSInfo() {
	try {
		const r = await fetch(`${_state.config.coreUrl}/api/os`);
		if (r.ok) _state.osInfo = await r.json();
	} catch {
		_state.osInfo = null;
	}
}

/**
 * Fetch the last successful scan result from the backend's persistent cache.
 */
export async function initStore() {
	try {
		const r = await fetch(`${_state.config.coreUrl}/api/scan/cached`);
		if (r.ok) {
			_state.scanResult = await r.json();
		}
	} catch {
		/* ignore */
	}
}

export async function pollGPU() {
	try {
		const r = await fetch(`${_state.config.coreUrl}/api/gpu`);
		if (r.ok) _state.gpuInfo = await r.json();
	} catch {
		/* silent */
	}
}

// ─── Scan (config-driven) ────────────────────────────────────────────────────
export function startScan(): Promise<void> {
	return new Promise((resolve) => {
		if (_state.scanning) return;
		_state.scanning = true;
		_state.scanEvents = [];
		_state.scanPercent = 0;
		_state.scanResult = null;
		_state.scanMeta = {
			startedAt: new Date().toISOString(),
			status: "running",
			totalModels: 0,
			orphanedModels: 0,
			runningServices: 0,
			mcpServers: 0,
			agentCards: 0,
			diskVolumes: 0,
			eventsCount: 0,
		};

		// Pass config params as query string so the Go core can respect them
		const params = new URLSearchParams({
			deepScan: String(_state.config.deepScanEnabled),
			maxDepth: String(_state.config.deepScanMaxDepth),
			minModelMB: String(_state.config.minModelSizeMB),
			ignoreDirs: _state.config.ignoreDirs,
			extraDirs: _state.config.extraModelDirs,
		});
		const wsUrl = `${_state.config.coreWs}/ws/scan?${params}`;
		const ws = new WebSocket(wsUrl);
		const startMs = Date.now();

		ws.onmessage = (e) => {
			try {
				const evt: ScanEvent = JSON.parse(e.data);
				_state.scanEvents = [..._state.scanEvents, evt];
				_state.scanPercent = evt.percent;
				_state.scanPhase = evt.phase;
				if (_state.scanMeta)
					_state.scanMeta.eventsCount = _state.scanEvents.length;

				if (evt.type === "complete" && evt.detail) {
					const result = evt.detail as ScanResult;
					_state.scanResult = result;
					_state.scanning = false;
					_state.scanMeta = {
						startedAt: _state.scanMeta!.startedAt,
						completedAt: new Date().toISOString(),
						durationMs: Date.now() - startMs,
						status: "complete",
						totalModels: result.models?.length ?? 0,
						orphanedModels:
							result.models?.filter((m) => m.is_orphaned)
								.length ?? 0,
						runningServices:
							result.components?.filter(
								(c) => c.status === "Running",
							).length ?? 0,
						mcpServers: result.mcp_servers?.length ?? 0,
						agentCards: result.agent_cards?.length ?? 0,
						diskVolumes: result.disks?.length ?? 0,
						eventsCount: _state.scanEvents.length,
					};
					addToast(
						"Scan complete — intelligence report ready.",
						"success",
						"/scan",
					);
					resolve();
				}
				if (evt.type === "error") {
					_state.scanning = false;
					if (_state.scanMeta) _state.scanMeta.status = "error";
					addToast("Scan error. Check core logs.", "error", "/scan");
					resolve();
				}
			} catch {
				/* ignore parse errors */
			}
		};
		ws.onerror = () => {
			_state.scanning = false;
			if (_state.scanMeta) _state.scanMeta.status = "error";
			addToast("WebSocket error — is the core running?", "error");
			resolve();
		};
		ws.onclose = () => {
			if (_state.scanning) {
				_state.scanning = false;
				resolve();
			}
		};
	});
}

// ─── File ops ─────────────────────────────────────────────────────────────────
export async function openPath(p: string) {
	try {
		await fetch(`${_state.config.coreUrl}/api/open`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ path: p }),
		});
	} catch {
		addToast("Could not open path — is the core running?", "warn");
	}
}
export async function readFile(p: string): Promise<string | null> {
	try {
		const r = await fetch(`${_state.config.coreUrl}/api/read`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ path: p }),
		});
		if (r.ok) return await r.text();
	} catch {
		/* silent */
	}
	return null;
}
export async function requestDelete(p: string): Promise<boolean> {
	const r = await fetch(`${_state.config.coreUrl}/api/delete`, {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ file_path: p, confirm: true }),
	});
	return r.ok;
}

// ─── Schedules ────────────────────────────────────────────────────────────────
export function loadSchedules() {
	try {
		const s = localStorage.getItem("aeye-schedules");
		if (s) _state.schedules = JSON.parse(s);
	} catch {
		_state.schedules = [];
	}
}
export function saveSchedule(rule: ScheduleRule) {
	const idx = _state.schedules.findIndex((s) => s.id === rule.id);
	_state.schedules =
		idx >= 0
			? _state.schedules.map((s, i) => (i === idx ? rule : s))
			: [..._state.schedules, rule];
	localStorage.setItem("aeye-schedules", JSON.stringify(_state.schedules));
}
export function deleteSchedule(id: string) {
	_state.schedules = _state.schedules.filter((s) => s.id !== id);
	localStorage.setItem("aeye-schedules", JSON.stringify(_state.schedules));
}
export function startScheduler() {
	if (_state.scheduleTimer) return;
	_state.scheduleTimer = setInterval(() => {
		const now = new Date();
		for (const rule of _state.schedules) {
			if (!rule.enabled) continue;
			const next = rule.nextRun ? new Date(rule.nextRun) : null;
			if (next && now >= next && !_state.scanning) {
				addToast(
					`Scheduled scan "${rule.label}" starting…`,
					"info",
					"/scan",
				);
				startScan().then(() =>
					saveSchedule({
						...rule,
						lastRun: new Date().toISOString(),
						nextRun: computeNextRun(rule.cron),
					}),
				);
			}
		}
	}, 60_000);
}
export function computeNextRun(cron: string): string {
	const m = cron.match(/(\d+)\s*(min|hour|day)/i);
	if (m) {
		const n = parseInt(m[1]),
			unit = m[2].toLowerCase();
		const ms = unit.startsWith("min")
			? n * 60_000
			: unit.startsWith("hour")
				? n * 3_600_000
				: n * 86_400_000;
		return new Date(Date.now() + ms).toISOString();
	}
	return new Date(Date.now() + 3_600_000).toISOString();
}

// ─── Plugins ─────────────────────────────────────────────────────────────────
export function loadPlugins() {
	try {
		const s = localStorage.getItem("aeye-plugins");
		if (s) _state.plugins = JSON.parse(s);
	} catch {
		_state.plugins = [];
	}
}
export function togglePlugin(id: string) {
	_state.plugins = _state.plugins.map((p) =>
		p.id === id ? { ...p, enabled: !p.enabled } : p,
	);
	localStorage.setItem("aeye-plugins", JSON.stringify(_state.plugins));
}
export function installPlugin(plugin: Plugin) {
	if (_state.plugins.some((p) => p.id === plugin.id)) {
		addToast("Plugin already installed", "warn");
		return;
	}
	_state.plugins = [..._state.plugins, { ...plugin, enabled: true }];
	localStorage.setItem("aeye-plugins", JSON.stringify(_state.plugins));
	addToast(`Plugin "${plugin.name}" installed`, "success");
}
export function uninstallPlugin(id: string) {
	_state.plugins = _state.plugins.filter((p) => p.id !== id);
	localStorage.setItem("aeye-plugins", JSON.stringify(_state.plugins));
	addToast("Plugin removed", "info");
}

// ─── Helpers ─────────────────────────────────────────────────────────────────
export const THEMES: Record<Theme, { label: string; accent: string }> = {
	obsidian: { label: "Obsidian", accent: "#c084fc" },
	clay: { label: "Clay", accent: "#fb923c" },
	forest: { label: "Forest", accent: "#4ade80" },
	slate: { label: "Slate", accent: "#7dd3fc" },
	amber: { label: "Amber", accent: "#fbbf24" },
};
export function formatBytes(b: number): string {
	if (!b) return "0 B";
	const k = 1024,
		s = ["B", "KB", "MB", "GB", "TB"],
		i = Math.floor(Math.log(b) / Math.log(k));
	return `${(b / Math.pow(k, i)).toFixed(1)} ${s[i]}`;
}
export function formatGB(gb: number): string {
	return gb >= 1024 ? `${(gb / 1024).toFixed(1)} TB` : `${gb.toFixed(1)} GB`;
}
export function formatDuration(ms: number): string {
	if (ms < 1000) return `${ms}ms`;
	if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`;
	return `${Math.floor(ms / 60000)}m ${Math.floor((ms % 60000) / 1000)}s`;
}
