<script lang="ts">
	import { page } from "$app/stores";
	import { onMount } from "svelte";
	import {
		store,
		setTheme,
		THEMES,
		pollGPU,
		fetchOSInfo,
		dismissToast,
		startScan,
		addToast,
		initStore,
		loadSchedules,
		startScheduler,
	} from "$lib/stores/app.svelte";
	import "../app.css";

	let { children } = $props();
	let sidebarOpen = $state(true);
	let gpuInterval: ReturnType<typeof setInterval>;

	onMount(async () => {
		const saved = localStorage.getItem("aeye-theme");
		if (saved && saved in THEMES) setTheme(saved as keyof typeof THEMES);
		else setTheme("obsidian");
		await fetchOSInfo();
		initStore(); // Load cached scan on boot
		gpuInterval = setInterval(pollGPU, 3000);
		pollGPU();
		loadSchedules();
		startScheduler();
		return () => clearInterval(gpuInterval);
	});

	// Nav structure: flat with section separators
	const nav = [
		{ href: "/", label: "Dashboard", icon: "⬡", section: null },
		{ href: "/scan", label: "Scan Center", icon: "◉", section: null },
		{
			href: "/services",
			label: "Live Services",
			icon: "⚙",
			section: "Intel",
		},
		{ href: "/agents", label: "Agents", icon: "◈", section: null },
		{ href: "/acp", label: "ACP Agents", icon: "▸", section: null },
		{ href: "/gateways", label: "Gateways", icon: "⟷", section: null },
		{ href: "/mcp", label: "MCP Servers", icon: "⬡", section: null },
		{ href: "/models", label: "Model Cache", icon: "◆", section: "Assets" },
		{ href: "/plugins", label: "Plugins", icon: "⬡", section: "Ops" },
		{ href: "/schedule", label: "Scheduler", icon: "◷", section: null },
		{
			href: "/alerts",
			label: "Alerts",
			icon: "◷",
			section: null,
			pro: true,
		},
		{
			href: "/export",
			label: "Export",
			icon: "◆",
			section: null,
			pro: true,
		},
		{
			href: "/remote",
			label: "Remote Scan",
			icon: "◈",
			section: null,
			pro: true,
		},
		{
			href: "/config",
			label: "Configuration",
			icon: "⚙",
			section: "System",
		},
		{ href: "/docs", label: "Documentation", icon: "≡", section: null },
		{ href: "/about", label: "About", icon: "◉", section: null },
	];

	function isActive(href: string) {
		if (href === "/") return $page.url.pathname === "/";
		return $page.url.pathname.startsWith(href);
	}

	let lastSection = "";
</script>

<div class="shell" data-theme={store.currentTheme}>
	<div class="grain"></div>

	<!-- Sidebar -->
	<aside class="sidebar" class:closed={!sidebarOpen}>
		<!-- Logo -->
		<div class="sidebar-logo">
			<a href="/" class="logo-mark">
				<svg
					class="eye-svg"
					viewBox="0 0 40 26"
					fill="none"
					xmlns="http://www.w3.org/2000/svg"
				>
					<ellipse
						cx="20"
						cy="13"
						rx="19"
						ry="11"
						stroke="currentColor"
						stroke-width="1.5"
						fill="none"
						opacity="0.7"
					/>
					<circle
						cx="20"
						cy="13"
						r="6"
						stroke="currentColor"
						stroke-width="1.5"
						fill="none"
					/>
					<circle
						cx="20"
						cy="13"
						r="2.5"
						fill="currentColor"
						opacity="0.9"
						class="pupil"
					/>
					<line
						x1="1"
						y1="13"
						x2="6"
						y2="13"
						stroke="currentColor"
						stroke-width="1"
						opacity="0.4"
					/>
					<line
						x1="34"
						y1="13"
						x2="39"
						y2="13"
						stroke="currentColor"
						stroke-width="1"
						opacity="0.4"
					/>
				</svg>
				{#if sidebarOpen}
					<div class="logo-text">
						<span class="logo-name">aEye</span>
						<span class="logo-sub">WATCHER</span>
					</div>
				{/if}
			</a>
			<button
				class="toggle-btn"
				onclick={() => (sidebarOpen = !sidebarOpen)}
			>
				{sidebarOpen ? "◀" : "▶"}
			</button>
		</div>

		<!-- OS badge -->
		{#if store.osInfo && sidebarOpen}
			<div class="os-badge">
				<span class="os-icon"
					>{store.osInfo.os === "linux"
						? "🐧"
						: store.osInfo.os === "darwin"
							? ""
							: "🪟"}</span
				>
				<div class="os-info">
					<span class="os-name"
						>{store.osInfo.distro ?? store.osInfo.os}</span
					>
					<span class="os-detail"
						>{store.osInfo.arch} · {store.osInfo.hostname}</span
					>
				</div>
				{#if store.osInfo.gpu_available}<span class="gpu-pill">GPU</span
					>{/if}
			</div>
		{/if}

		<!-- Nav -->
		<nav class="nav">
			{#each nav as item}
				{#if item.section && item.section !== lastSection}
					{#if sidebarOpen}
						<div class="nav-section-label">{item.section}</div>
					{:else}
						<div class="nav-divider"></div>
					{/if}
				{/if}
				{(lastSection = item.section ?? lastSection) && ""}
				<a
					href={item.href}
					class="nav-item"
					class:active={isActive(item.href)}
				>
					<span class="nav-icon">{item.icon}</span>
					{#if sidebarOpen}
						<span class="nav-label">{item.label}</span>
						{#if item.href === "/scan" && store.scanning}
							<span class="scanning-pip"></span>
						{/if}
						{#if "pro" in item && item.pro && !store.isPro}
							<span class="pro-pip">PRO</span>
						{/if}
					{/if}
				</a>
			{/each}
		</nav>

		<!-- Footer -->
		{#if sidebarOpen}
			<div class="sidebar-footer">
				<div class="theme-row">
					<span class="theme-label">THEME</span>
					<div class="theme-pills">
						{#each Object.entries(THEMES) as [key, val]}
							<button
								class="theme-pill"
								class:active={store.currentTheme === key}
								style="--pill:{val.accent}"
								onclick={() =>
									setTheme(key as keyof typeof THEMES)}
								title={val.label}
							></button>
						{/each}
					</div>
				</div>
				<a
					class="gh-link"
					href="https://github.com/NathanielRand/aEyeWatcher"
					target="_blank"
					rel="noopener"
				>
					⬡ GitHub
				</a>
			</div>
		{/if}
	</aside>

	<!-- Page content -->
	<div class="page-wrap">
		<!-- Topbar with GPU -->
		<div class="topbar">
			<div class="topbar-left">
				{#if store.scanning}
					<a href="/scan" class="scan-running-badge">
						<span class="scan-pip"></span> SCAN RUNNING — {store.scanPercent}%
						· {store.scanPhase.toUpperCase()}
					</a>
				{/if}
			</div>
			<div class="topbar-right">
				{#if store.gpuInfo && !store.gpuInfo.error}
					<div class="gpu-widget">
						<span class="gpu-label">VRAM</span>
						<div class="gpu-track">
							<div
								class="gpu-fill"
								style="width:{store.gpuInfo
									.use_pct}%;background:{store.gpuInfo
									.use_pct > 85
									? '#ef4444'
									: store.gpuInfo.use_pct > 60
										? '#f59e0b'
										: 'var(--accent)'}"
							></div>
						</div>
						<span class="gpu-val"
							>{store.gpuInfo.used_mb}/{store.gpuInfo
								.total_mb}M</span
						>
					</div>
				{/if}
				<a
					href="/scan"
					class="scan-cta-btn"
					class:scanning={store.scanning}
				>
					{#if store.scanning}<span class="spin">⟳</span>{:else}◉{/if}
					{store.scanning ? "SCANNING" : "NEW SCAN"}
				</a>
			</div>
		</div>

		<main class="main-content">
			{@render children()}
		</main>
	</div>

	<!-- Toast rack -->
	<div class="toast-rack">
		{#each store.toasts as toast (toast.id)}
			<div class="toast toast-{toast.type}" role="alert">
				<span class="toast-msg">{toast.msg}</span>
				{#if toast.link}
					<a href={toast.link} class="toast-link">View →</a>
				{/if}
				<button
					class="toast-close"
					onclick={() => dismissToast(toast.id)}>✕</button
				>
			</div>
		{/each}
	</div>
</div>

<style>
	@import url("https://fonts.googleapis.com/css2?family=Special+Elite&family=Courier+Prime:wght@400;700&family=Rajdhani:wght@400;500;600;700&display=swap");

	:global([data-theme="obsidian"]) {
		--accent: #c084fc;
		--accent-dim: rgba(192, 132, 252, 0.15);
		--bg: #08080f;
		--bg2: #0f0f1a;
		--bg3: #16162a;
		--border: rgba(192, 132, 252, 0.14);
		--text: #e8e6f8;
		--text-dim: #8a8aa8;
		--text-muted: #4a4a68;
	}
	:global([data-theme="clay"]) {
		--accent: #fb923c;
		--accent-dim: rgba(251, 146, 60, 0.13);
		--bg: #100a06;
		--bg2: #1a1008;
		--bg3: #241608;
		--border: rgba(251, 146, 60, 0.12);
		--text: #f2ebe2;
		--text-dim: #9a8070;
		--text-muted: #583a28;
	}
	:global([data-theme="forest"]) {
		--accent: #4ade80;
		--accent-dim: rgba(74, 222, 128, 0.12);
		--bg: #060e08;
		--bg2: #081a0e;
		--bg3: #0c2414;
		--border: rgba(74, 222, 128, 0.12);
		--text: #e2f0e5;
		--text-dim: #6a9870;
		--text-muted: #2a5035;
	}
	:global([data-theme="slate"]) {
		--accent: #7dd3fc;
		--accent-dim: rgba(125, 211, 252, 0.12);
		--bg: #060810;
		--bg2: #0a1020;
		--bg3: #0e1830;
		--border: rgba(125, 211, 252, 0.12);
		--text: #e2ecf8;
		--text-dim: #6a7a98;
		--text-muted: #283050;
	}
	:global([data-theme="amber"]) {
		--accent: #fbbf24;
		--accent-dim: rgba(251, 191, 36, 0.12);
		--bg: #0a0800;
		--bg2: #150f00;
		--bg3: #1f1600;
		--border: rgba(251, 191, 36, 0.12);
		--text: #f8f0d0;
		--text-dim: #b8982a;
		--text-muted: #6a5800;
	}

	:global(*, *::before, *::after) {
		box-sizing: border-box;
		margin: 0;
		padding: 0;
	}
	:global(body) {
		background: var(--bg);
		color: var(--text);
		font-family: "Rajdhani", monospace;
		font-size: 15.5px;
		overflow: hidden;
		height: 100vh;
	}
	:global(a) {
		color: inherit;
		text-decoration: none;
	}
	:global(h1, h2, h3) {
		font-family: "Special Elite", serif;
	}
	:global(code) {
		font-family: "Courier Prime", monospace;
		font-size: 0.85em;
		background: var(--bg3);
		border: 1px solid var(--border);
		padding: 1px 5px;
		border-radius: 3px;
		color: var(--accent);
	}
	:global(button) {
		font-family: "Rajdhani", monospace;
	}
	:global(::-webkit-scrollbar) {
		width: 4px;
		height: 4px;
	}
	:global(::-webkit-scrollbar-thumb) {
		background: var(--border);
		border-radius: 2px;
	}

	.grain {
		position: fixed;
		inset: 0;
		pointer-events: none;
		z-index: 9999;
		background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='0.04'/%3E%3C/svg%3E");
		background-size: 200px;
		opacity: 0.55;
		mix-blend-mode: overlay;
	}
	.shell {
		display: flex;
		height: 100vh;
		overflow: hidden;
	}

	/* Sidebar */
	.sidebar {
		width: 220px;
		min-width: 220px;
		background: var(--bg2);
		border-right: 1px solid var(--border);
		display: flex;
		flex-direction: column;
		transition:
			width 0.2s,
			min-width 0.2s;
		overflow: hidden;
		flex-shrink: 0;
		z-index: 50;
	}
	.sidebar.closed {
		width: 52px;
		min-width: 52px;
	}

	.sidebar-logo {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.85rem 0.75rem;
		border-bottom: 1px solid var(--border);
		flex-shrink: 0;
	}
	.logo-mark {
		display: flex;
		align-items: center;
		gap: 0.6rem;
		color: var(--accent);
	}
	.eye-svg {
		width: 32px;
		height: 21px;
		color: var(--accent);
		flex-shrink: 0;
	}

	/* Eye blink animation */
	.eye-svg {
		animation: eye-idle 8s ease-in-out infinite;
	}
	@keyframes eye-idle {
		0%,
		45%,
		55%,
		100% {
			transform: scaleY(1);
		}
		48%,
		52% {
			transform: scaleY(0.08);
		}
	}
	.pupil {
		animation: pupil-look 8s ease-in-out infinite;
	}
	@keyframes pupil-look {
		0%,
		40% {
			transform: translateX(0);
		}
		20% {
			transform: translateX(2px);
		}
		60%,
		100% {
			transform: translateX(0);
		}
		80% {
			transform: translateX(-2px);
		}
	}

	.logo-text {
		display: flex;
		flex-direction: column;
	}
	.logo-name {
		font-family: "Special Elite", serif;
		font-size: 1.1rem;
		color: var(--accent);
		letter-spacing: 0.08em;
		line-height: 1.1;
	}
	.logo-sub {
		font-size: 0.55rem;
		letter-spacing: 0.4em;
		color: var(--text-muted);
		text-transform: uppercase;
	}
	.toggle-btn {
		background: none;
		border: none;
		color: var(--text-muted);
		cursor: pointer;
		font-size: 0.7rem;
		padding: 3px 5px;
		transition: color 0.2s;
		flex-shrink: 0;
	}
	.toggle-btn:hover {
		color: var(--accent);
	}

	.os-badge {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.55rem 0.75rem;
		background: var(--accent-dim);
		border-bottom: 1px solid var(--border);
		flex-shrink: 0;
	}
	.os-icon {
		font-size: 0.95rem;
	}
	.os-info {
		flex: 1;
		min-width: 0;
	}
	.os-name {
		display: block;
		color: var(--accent);
		font-weight: 600;
		font-size: 0.78rem;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}
	.os-detail {
		font-size: 0.65rem;
		color: var(--text-dim);
	}
	.gpu-pill {
		background: var(--accent);
		color: var(--bg);
		font-size: 0.55rem;
		font-weight: 700;
		padding: 1px 5px;
		border-radius: 2px;
		letter-spacing: 0.1em;
		flex-shrink: 0;
	}

	.nav {
		flex: 1;
		overflow-y: auto;
		padding: 0.4rem;
		display: flex;
		flex-direction: column;
		gap: 1px;
	}
	.nav-section-label {
		font-size: 0.58rem;
		letter-spacing: 0.25em;
		color: var(--text-muted);
		padding: 0.6rem 0.6rem 0.2rem;
		text-transform: uppercase;
		font-family: "Rajdhani", monospace;
	}
	.nav-divider {
		height: 1px;
		background: var(--border);
		margin: 0.4rem 0.4rem;
	}
	.nav-item {
		display: flex;
		align-items: center;
		gap: 0.6rem;
		padding: 0.5rem 0.6rem;
		border-radius: 4px;
		color: var(--text-dim);
		font-size: 0.88rem;
		font-weight: 500;
		letter-spacing: 0.03em;
		transition: all 0.12s;
		position: relative;
	}
	.nav-item:hover {
		background: var(--accent-dim);
		color: var(--text);
	}
	.nav-item.active {
		background: var(--accent-dim);
		color: var(--accent);
		border-left: 2px solid var(--accent);
	}
	.nav-icon {
		width: 16px;
		text-align: center;
		font-size: 0.9rem;
		flex-shrink: 0;
	}
	.nav-label {
		flex: 1;
	}
	.scanning-pip {
		width: 7px;
		height: 7px;
		border-radius: 50%;
		background: var(--accent);
		animation: pip-pulse 1s ease-in-out infinite;
		flex-shrink: 0;
	}
	@keyframes pip-pulse {
		0%,
		100% {
			opacity: 1;
			transform: scale(1);
		}
		50% {
			opacity: 0.5;
			transform: scale(0.7);
		}
	}

	.sidebar-footer {
		padding: 0.7rem;
		border-top: 1px solid var(--border);
		display: flex;
		flex-direction: column;
		gap: 0.55rem;
		flex-shrink: 0;
	}
	.theme-row {
		display: flex;
		align-items: center;
		gap: 0.6rem;
	}
	.theme-label {
		font-size: 0.58rem;
		letter-spacing: 0.2em;
		color: var(--text-muted);
		white-space: nowrap;
	}
	.theme-pills {
		display: flex;
		gap: 5px;
	}
	.theme-pill {
		width: 15px;
		height: 15px;
		border-radius: 50%;
		border: 2px solid transparent;
		background: var(--pill);
		cursor: pointer;
		transition:
			transform 0.15s,
			border-color 0.15s;
		flex-shrink: 0;
	}
	.theme-pill:hover {
		transform: scale(1.25);
	}
	.theme-pill.active {
		border-color: var(--text);
		transform: scale(1.15);
	}
	.gh-link {
		display: flex;
		align-items: center;
		gap: 0.4rem;
		font-size: 0.75rem;
		color: var(--text-dim);
		padding: 0.3rem 0.4rem;
		border: 1px solid var(--border);
		border-radius: 4px;
		transition: all 0.15s;
	}
	.gh-link:hover {
		color: var(--accent);
		border-color: var(--accent);
	}

	/* Page wrap */
	.page-wrap {
		flex: 1;
		display: flex;
		flex-direction: column;
		overflow: hidden;
		min-width: 0;
	}

	.topbar {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.5rem 1.25rem;
		background: var(--bg2);
		border-bottom: 1px solid var(--border);
		flex-shrink: 0;
		min-height: 44px;
	}
	.topbar-left,
	.topbar-right {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}
	.scan-running-badge {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.72rem;
		color: var(--accent);
		letter-spacing: 0.1em;
		background: var(--accent-dim);
		border: 1px solid var(--accent);
		padding: 0.25rem 0.65rem;
		border-radius: 3px;
	}
	.scan-pip {
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: var(--accent);
		animation: pip-pulse 1s ease-in-out infinite;
	}
	.gpu-widget {
		display: flex;
		align-items: center;
		gap: 0.45rem;
		font-size: 0.68rem;
	}
	.gpu-label {
		color: var(--text-muted);
		letter-spacing: 0.08em;
	}
	.gpu-track {
		width: 72px;
		height: 4px;
		background: var(--bg3);
		border-radius: 2px;
		overflow: hidden;
	}
	.gpu-fill {
		height: 100%;
		border-radius: 2px;
		transition:
			width 0.5s,
			background 0.5s;
	}
	.gpu-val {
		color: var(--text-dim);
		font-family: "Courier Prime", monospace;
		font-size: 0.62rem;
	}
	.scan-cta-btn {
		display: flex;
		align-items: center;
		gap: 0.35rem;
		background: var(--accent-dim);
		border: 1px solid var(--accent);
		color: var(--accent);
		padding: 0.3rem 0.75rem;
		border-radius: 3px;
		font-size: 0.75rem;
		font-weight: 700;
		letter-spacing: 0.08em;
		transition: all 0.15s;
		cursor: pointer;
	}
	.scan-cta-btn:hover {
		background: var(--accent);
		color: var(--bg);
	}
	.scan-cta-btn.scanning {
		opacity: 0.7;
		cursor: default;
	}
	.spin {
		animation: spin 0.8s linear infinite;
		display: inline-block;
	}
	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.main-content {
		flex: 1;
		overflow-y: auto;
		overflow-x: hidden;
		scrollbar-width: thin;
		scrollbar-color: var(--border) transparent;
		min-height: 0;
	}

	/* Toasts */
	.toast-rack {
		position: fixed;
		bottom: 1.5rem;
		right: 1.5rem;
		z-index: 10000;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		pointer-events: none;
	}
	.toast {
		display: flex;
		align-items: center;
		gap: 0.65rem;
		padding: 0.6rem 0.9rem;
		border-radius: 6px;
		font-size: 0.82rem;
		pointer-events: all;
		min-width: 280px;
		max-width: 420px;
		border: 1px solid;
		animation: toast-in 0.2s ease;
		box-shadow: 0 4px 24px rgba(0, 0, 0, 0.5);
	}
	@keyframes toast-in {
		from {
			opacity: 0;
			transform: translateX(20px);
		}
		to {
			opacity: 1;
			transform: translateX(0);
		}
	}
	.toast-info {
		background: var(--bg2);
		border-color: var(--border);
		color: var(--text-dim);
	}
	.toast-success {
		background: color-mix(in srgb, #4ade80 10%, var(--bg2));
		border-color: #4ade80;
		color: #4ade80;
	}
	.toast-error {
		background: color-mix(in srgb, #ef4444 10%, var(--bg2));
		border-color: #ef4444;
		color: #ef4444;
	}
	.toast-warn {
		background: color-mix(in srgb, #fb923c 10%, var(--bg2));
		border-color: #fb923c;
		color: #fb923c;
	}
	.toast-msg {
		flex: 1;
		line-height: 1.4;
	}
	.toast-link {
		color: inherit;
		text-decoration: underline;
		font-size: 0.78rem;
		white-space: nowrap;
	}
	.toast-close {
		background: none;
		border: none;
		color: inherit;
		cursor: pointer;
		opacity: 0.6;
		font-size: 0.75rem;
		padding: 0 0 0 0.25rem;
		transition: opacity 0.15s;
	}
	.toast-close:hover {
		opacity: 1;
	}

	.pro-pip {
		font-size: 0.52rem;
		padding: 1px 4px;
		background: color-mix(in srgb, var(--accent) 15%, transparent);
		border: 1px solid var(--accent);
		color: var(--accent);
		border-radius: 2px;
		letter-spacing: 0.08em;
		flex-shrink: 0;
	}
	@media (max-width: 768px) {
		.sidebar {
			width: 52px;
			min-width: 52px;
		}
		.logo-text,
		.os-badge,
		.sidebar-footer,
		.nav-label,
		.nav-section-label {
			display: none;
		}
		.topbar {
			padding: 0.5rem 0.75rem;
		}
		.gpu-widget {
			display: none;
		}
	}
</style>
