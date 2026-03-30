<script lang="ts">
	import {
		store,
		startScan,
		formatBytes,
		formatGB,
		formatDuration,
	} from "$lib/stores/app.svelte";

	const phaseLabels: Record<string, string> = {
		system: "SYSTEM",
		gpu: "GPU",
		disk: "STORAGE",
		components: "COMPONENTS",
		mcp: "MCP",
		acp: "ACP",
		models: "MODELS",
		agents: "AGENTS",
		complete: "COMPLETE",
	};
	const phases = [
		"system",
		"gpu",
		"disk",
		"components",
		"mcp",
		"acp",
		"models",
		"agents",
	];

	let lastProcessed = 0;
	let activityLog = $state<{ msg: string; type: string; ts: string }[]>([]);

	$effect(() => {
		const len = store.scanEvents.length;
		if (len > lastProcessed) {
			const newEvts = store.scanEvents.slice(lastProcessed, len);
			lastProcessed = len;
			for (const e of newEvts) {
				if (e.type === "found" || e.type === "checkpoint") {
					const ts = new Date().toLocaleTimeString();
					activityLog = [
						{ msg: e.message, type: e.type, ts },
						...activityLog,
					].slice(0, 100);
				}
			}
		}
	});

	function phaseStatus(ph: string) {
		const done = store.scanEvents.some(
			(e) =>
				e.phase === ph && (e.type === "found" || e.type === "complete"),
		);
		const active = store.scanPhase === ph;
		return done ? "done" : active ? "active" : "idle";
	}

	let runningCount = $derived(
		(store.scanResult?.components ?? []).filter(
			(c) => c.status === "Running",
		).length,
	);
	let orphanCount = $derived(
		(store.scanResult?.models ?? []).filter((m) => m.is_orphaned).length,
	);
	let totalGB = $derived(
		(store.scanResult?.models ?? []).reduce((s, m) => s + m.size_bytes, 0) /
			1e9,
	);
</script>

<div class="scan-page">
	<!-- Header / CTA -->
	<div
		class="scan-hero"
		class:has-result={!!store.scanResult}
		class:scanning={store.scanning}
	>
		<div class="hero-left">
			<div class="hero-eye" class:scanning={store.scanning}>◉</div>
			<div>
				<h1 class="hero-title">
					{#if store.scanning}Scan In Progress{:else if store.scanResult}Intelligence
						Report{:else}aEye Scan Center{/if}
				</h1>
				<p class="hero-sub">
					{#if store.scanning}
						Phase: <strong
							>{phaseLabels[store.scanPhase] ??
								store.scanPhase}</strong
						>
						· {store.scanPercent}% complete
					{:else if store.scanResult && store.scanMeta}
						Completed {new Date(
							store.scanMeta.completedAt ?? "",
						).toLocaleTimeString()} · Duration: {formatDuration(
							store.scanMeta.durationMs ?? 0,
						)}
					{:else}
						Full ecosystem scan — providers, models, agents, MCP
						servers, disk, GPU
					{/if}
				</p>
			</div>
		</div>
		{#if !store.scanning}
			<a href="/history" class="scan-history-link"> View scan history </a>
			<button
				class="start-scan-btn"
				onclick={() => {
					lastProcessed = 0;
					activityLog = [];
					startScan();
				}}
			>
				◉ {store.scanResult ? "RE-SCAN" : "START SCAN"}
			</button>
		{:else}
			<div class="scanning-indicator">
				<span class="spin-ring"></span>
				<span>SCANNING…</span>
			</div>
		{/if}
	</div>

	<!-- Progress (only while scanning) -->
	{#if store.scanning}
		<div class="progress-section">
			<div class="phase-track">
				{#each phases as ph}
					{@const st = phaseStatus(ph)}
					<div
						class="phase-node"
						class:done={st === "done"}
						class:active={st === "active"}
					>
						<div class="phase-dot">
							{st === "done" ? "✓" : st === "active" ? "◉" : "○"}
						</div>
						<div class="phase-name">{phaseLabels[ph]}</div>
					</div>
					{#if ph !== phases[phases.length - 1]}
						<div
							class="phase-connector"
							class:filled={phaseStatus(ph) === "done"}
						></div>
					{/if}
				{/each}
			</div>
			<div class="progress-bar-wrap">
				<div
					class="progress-bar-fill"
					style="width:{store.scanPercent}%"
				>
					<div class="progress-shimmer"></div>
				</div>
				<span class="progress-pct">{store.scanPercent}%</span>
			</div>
			<div class="live-msg">
				{store.scanEvents.length > 0
					? store.scanEvents[store.scanEvents.length - 1].message
					: "Initializing…"}
			</div>
		</div>
	{/if}

	{#if !store.scanning && !store.scanResult}
		<!-- Empty state CTA -->
		<div class="empty-hero">
			<div class="empty-eye">◉</div>
			<h2>No scan data yet</h2>
			<p>
				Run a scan to discover every AI component, model, agent, and
				service on this system.
			</p>
			<div class="what-we-scan">
				{#each ["Providers & Servers", "Local Models", "Agent Ecosystems", "MCP Configurations", "ACP Agents", "Protocol Gateways", "Vector Databases", "Disk Storage", "GPU VRAM"] as item}
					<span class="scan-tag">◆ {item}</span>
				{/each}
			</div>
			<button class="start-scan-btn large" onclick={() => startScan()}
				>◉ START SCAN</button
			>
		</div>
	{/if}

	{#if store.scanResult && store.scanMeta}
		<!-- Summary cards -->
		<div class="summary-grid">
			<a href="/models" class="sum-card accent">
				<div class="sum-icon">◆</div>
				<div class="sum-val">{store.scanMeta.totalModels}</div>
				<div class="sum-label">Models Found</div>
				{#if store.scanMeta.orphanedModels > 0}
					<div class="sum-sub warn">
						{store.scanMeta.orphanedModels} orphaned
					</div>
				{/if}
			</a>
			<a href="/services" class="sum-card">
				<div class="sum-icon">⚙</div>
				<div class="sum-val">{store.scanMeta.runningServices}</div>
				<div class="sum-label">Live Services</div>
				<div class="sum-sub">
					of {store.scanResult.components.length} checked
				</div>
			</a>
			<a href="/models" class="sum-card">
				<div class="sum-icon">◉</div>
				<div class="sum-val">
					{totalGB.toFixed(1)}<span class="sum-unit">GB</span>
				</div>
				<div class="sum-label">Model Storage</div>
			</a>
			<a href="/mcp" class="sum-card">
				<div class="sum-icon">⬡</div>
				<div class="sum-val">{store.scanMeta.mcpServers}</div>
				<div class="sum-label">MCP Servers</div>
			</a>
			<a href="/agents" class="sum-card">
				<div class="sum-icon">◈</div>
				<div class="sum-val">{store.scanMeta.agentCards}</div>
				<div class="sum-label">Agent Ecosystems</div>
				<div class="sum-sub">
					{store.scanResult.agent_cards.filter(
						(a) => a.status !== "not_found",
					).length} detected
				</div>
			</a>
			<a href="/schedule" class="sum-card">
				<div class="sum-icon">◷</div>
				<div class="sum-val">
					{formatDuration(store.scanMeta.durationMs ?? 0)}
				</div>
				<div class="sum-label">Scan Duration</div>
				<div class="sum-sub">{store.scanMeta.eventsCount} events</div>
			</a>
		</div>

		<!-- Live services - show running + installed, link to full services page -->
		{@const detectedSvcs = store.scanResult.components.filter(
			(c) => c.status !== "Not Found",
		)}
		{#if detectedSvcs.length > 0}
			<div class="result-card">
				<div class="rc-header">
					<span class="rc-icon">⚙</span> DETECTED SERVICES
					<span class="rc-count">{detectedSvcs.length}</span>
					<a href="/services" class="rc-link">View all →</a>
				</div>
				<div class="service-chips">
					{#each detectedSvcs.slice(0, 12) as comp}
						<a
							href="/services"
							class="service-chip"
							class:chip-running={comp.status === "Running"}
							class:chip-installed={comp.status === "Installed"}
						>
							<span
								class="chip-dot"
								class:dot-run={comp.status === "Running"}
								class:dot-inst={comp.status === "Installed"}
							></span>
							<span class="chip-name">{comp.name}</span>
							{#if comp.port}<span class="chip-port"
									>:{comp.port}</span
								>{/if}
						</a>
					{/each}
					{#if detectedSvcs.length > 12}
						<a href="/services" class="chip-more"
							>+{detectedSvcs.length - 12} more →</a
						>
					{/if}
				</div>
			</div>
		{/if}

		<!-- Disk summary -->
		{#if store.scanResult.disks.length > 0}
			<div class="result-card">
				<div class="rc-header">
					<span class="rc-icon">◉</span> STORAGE VOLUMES
					<span class="rc-count">{store.scanResult.disks.length}</span
					>
				</div>
				<div class="disk-rows">
					{#each store.scanResult.disks as disk}
						<div
							class="disk-row"
							style="background-color: var(--bg-lighter);"
						>
							<div class="disk-label-col">
								<span class="disk-mount">{disk.mount}</span>
								<span class="disk-device">{disk.device}</span>
							</div>
							<div class="disk-bar-col">
								<!-- Use the calculated disk usage percentage -->
								<div
									class="disk-bar"
									style="width: {disk.use_pct}%;"
								>
									<div class="disk-fill"></div>
								</div>
								<div class="disk-stats-row">
									<span class="disk-used"
										>{formatGB(disk.used_gb)} used</span
									>
									<span class="disk-pct"
										>{disk.use_pct.toFixed(1)}%</span
									>
									<!-- Use disk.use_pct -->
									<span class="disk-free"
										>{formatGB(disk.free_gb)} free</span
									>
									<span class="disk-total"
										>{formatGB(disk.total_gb)} total</span
									>
								</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}

		<!-- Scan event log -->
		{#if activityLog.length > 0}
			<div class="result-card">
				<div class="rc-header">
					<span class="rc-icon">≡</span> SCAN LOG
					<span class="rc-count">{activityLog.length}</span>
				</div>
				<div class="log-list">
					{#each activityLog.slice(0, 40) as entry}
						<div
							class="log-row"
							class:is-found={entry.type === "found"}
						>
							<span class="log-ts">{entry.ts}</span>
							<span class="log-icon"
								>{entry.type === "found" ? "◆" : "▸"}</span
							>
							<span class="log-msg">{entry.msg}</span>
						</div>
					{/each}
				</div>
			</div>
		{/if}
	{/if}
</div>

<style>
	.scan-page {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		padding: 1.25rem;
		overflow-y: auto;
		min-height: 100%;
	}

	.scan-hero {
		display: flex;
		align-items: center;
		justify-content: space-between;
		background: var(--bg2);
		border: 1px solid var(--border);
		border-radius: 8px;
		padding: 1.5rem;
		gap: 1rem;
		flex-wrap: wrap;
		position: relative;
		overflow: hidden;
	}
	.scan-hero.scanning {
		border-color: var(--accent);
		box-shadow: 0 0 30px var(--accent-dim);
	}
	.scan-hero::before {
		content: "◉";
		position: absolute;
		right: 2rem;
		top: 50%;
		transform: translateY(-50%);
		font-size: 8rem;
		color: var(--accent);
		opacity: 0.04;
		pointer-events: none;
	}
	.scan-history-link {
		color: var(--accent);
		text-decoration: none;
		font-size: 1.1rem;
		align-self: flex-center;
	}
	.scan-history-link:hover {
		text-decoration: underline;
		opacity: 0.8;
	}
	.hero-left {
		display: flex;
		align-items: center;
		gap: 1rem;
	}
	.hero-eye {
		font-size: 2rem;
		color: var(--accent);
		opacity: 0.8;
	}
	.hero-eye.scanning {
		animation: spin 2s linear infinite;
	}
	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
	.hero-title {
		font-family: "Special Elite", serif;
		font-size: 1.3rem;
		color: var(--text);
	}
	.hero-sub {
		font-size: 1rem;
		color: var(--text-dim);
		margin-top: 0.25rem;
	}
	.hero-sub strong {
		color: var(--accent);
	}

	.start-scan-btn {
		background: var(--accent);
		border: none;
		color: var(--bg);
		padding: 0.6rem 1.5rem;
		border-radius: 5px;
		cursor: pointer;
		font-family: "Special Elite", serif;
		font-size: 0.95rem;
		letter-spacing: 0.08em;
		transition: opacity 0.15s;
		white-space: nowrap;
		flex-shrink: 0;
	}
	.start-scan-btn:hover {
		opacity: 0.85;
	}
	.start-scan-btn.large {
		font-size: 1.05rem;
		padding: 0.75rem 2rem;
		margin-top: 1rem;
	}

	.scanning-indicator {
		display: flex;
		align-items: center;
		gap: 0.6rem;
		color: var(--accent);
		font-size: 0.8rem;
		letter-spacing: 0.15em;
	}
	.spin-ring {
		width: 18px;
		height: 18px;
		border: 2px solid var(--accent-dim);
		border-top-color: var(--accent);
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
		display: block;
	}

	/* Progress */
	.progress-section {
		background: var(--bg2);
		border: 1px solid var(--border);
		border-radius: 8px;
		padding: 1.25rem;
		display: flex;
		flex-direction: column;
		gap: 0.85rem;
	}
	.phase-track {
		display: flex;
		align-items: center;
		gap: 0;
	}
	.phase-node {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 3px;
		flex-shrink: 0;
	}
	.phase-dot {
		font-size: 0.75rem;
		color: var(--text-muted);
	}
	.phase-node.done .phase-dot {
		color: var(--accent);
	}
	.phase-node.active .phase-dot {
		color: var(--accent);
		animation: pip-pulse 1s ease-in-out infinite;
	}
	@keyframes pip-pulse {
		0%,
		100% {
			opacity: 1;
		}
		50% {
			opacity: 0.5;
		}
	}
	.phase-name {
		font-size: 0.52rem;
		color: var(--text-muted);
		letter-spacing: 0.08em;
		text-transform: uppercase;
		white-space: nowrap;
	}
	.phase-node.done .phase-name,
	.phase-node.active .phase-name {
		color: var(--accent);
	}
	.phase-connector {
		flex: 1;
		height: 1px;
		background: var(--border);
		margin: 0 2px;
		align-self: flex-start;
		margin-top: 10px;
	}
	.phase-connector.filled {
		background: var(--accent);
	}

	.progress-bar-wrap {
		background: var(--bg3);
		height: 8px;
		border-radius: 4px;
		overflow: hidden;
		position: relative;
	}
	.progress-bar-fill {
		height: 100%;
		background: var(--accent);
		border-radius: 4px;
		box-shadow: 0 0 10px var(--accent);
		position: relative;
		overflow: hidden;
		transition: width 0.4s ease;
	}
	.progress-shimmer {
		position: absolute;
		inset: 0;
		background: linear-gradient(
			90deg,
			transparent,
			rgba(255, 255, 255, 0.3),
			transparent
		);
		animation: shimmer 1.5s infinite;
	}
	@keyframes shimmer {
		0% {
			transform: translateX(-100%);
		}
		100% {
			transform: translateX(200%);
		}
	}
	.progress-pct {
		position: absolute;
		right: 0;
		top: -18px;
		font-size: 0.68rem;
		color: var(--accent);
		font-family: "Courier Prime", monospace;
	}
	.live-msg {
		font-size: 0.75rem;
		color: var(--text-dim);
		font-family: "Courier Prime", monospace;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	/* Empty */
	.empty-hero {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		padding: 3rem 1rem;
		gap: 1rem;
	}
	.empty-eye {
		font-size: 3rem;
		color: var(--accent);
		opacity: 0.4;
	}
	.empty-hero h2 {
		font-family: "Special Elite", serif;
		font-size: 1.4rem;
		color: var(--text);
	}
	.empty-hero p {
		font-size: 0.9rem;
		color: var(--text-dim);
		max-width: 480px;
		line-height: 1.7;
	}
	.what-we-scan {
		display: flex;
		flex-wrap: wrap;
		gap: 0.4rem;
		justify-content: center;
		margin-top: 0.5rem;
	}
	.scan-tag {
		font-size: 0.72rem;
		color: var(--text-dim);
		background: var(--bg2);
		border: 1px solid var(--border);
		padding: 0.2rem 0.5rem;
		border-radius: 3px;
	}

	/* Summary */
	.summary-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
		gap: 0.75rem;
	}
	.sum-card {
		background: var(--bg2);
		border: 1px solid var(--border);
		border-radius: 6px;
		padding: 1rem;
		position: relative;
		overflow: hidden;
		display: block;
		cursor: pointer;
		transition:
			border-color 0.15s,
			transform 0.1s;
		text-decoration: none;
	}
	.sum-card:hover {
		border-color: var(--accent);
		transform: translateY(-1px);
	}
	.sum-card::before {
		content: "";
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		height: 2px;
		background: var(--border);
	}
	.sum-card.accent::before {
		background: var(--accent);
	}
	.sum-icon {
		font-size: 1rem;
		color: var(--accent);
		margin-bottom: 0.35rem;
		opacity: 0.7;
	}
	.sum-val {
		font-family: "Special Elite", serif;
		font-size: 1.9rem;
		color: var(--text);
		line-height: 1;
	}
	.sum-unit {
		font-size: 1rem;
		color: var(--text-dim);
	}
	.sum-label {
		font-size: 0.68rem;
		color: var(--text-dim);
		letter-spacing: 0.1em;
		margin-top: 0.3rem;
	}
	.sum-sub {
		font-size: 0.65rem;
		color: var(--text-muted);
		margin-top: 0.2rem;
	}
	.sum-sub.warn {
		color: #fb923c;
	}

	/* Result cards */
	.result-card {
		background: var(--bg2);
		border: 1px solid var(--border);
		border-radius: 6px;
		overflow: visible;
	}
	.rc-header {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.6rem 1rem;
		background: var(--accent-dim);
		border-bottom: 1px solid var(--border);
		font-size: 0.72rem;
		letter-spacing: 0.15em;
		color: var(--accent);
		font-weight: 700;
	}
	.rc-icon {
		font-size: 0.9rem;
	}
	.rc-count {
		margin-left: auto;
		background: var(--bg3);
		border: 1px solid var(--border);
		padding: 1px 6px;
		border-radius: 2px;
		color: var(--text-dim);
		font-size: 0.65rem;
	}

	.service-chips {
		display: flex;
		flex-wrap: wrap;
		gap: 0.4rem;
		padding: 0.75rem 1rem;
	}
	.service-chip {
		display: flex;
		align-items: center;
		gap: 0.35rem;
		font-size: 0.78rem;
		background: var(--bg3);
		border: 1px solid var(--accent);
		padding: 0.25rem 0.55rem;
		border-radius: 3px;
		color: var(--accent);
	}
	.chip-dot {
		width: 7px;
		height: 7px;
		border-radius: 50%;
		background: var(--accent);
		box-shadow: 0 0 5px var(--accent);
		animation: pip-pulse 2s ease-in-out infinite;
	}
	.chip-port {
		color: var(--text-dim);
		font-family: "Courier Prime", monospace;
		font-size: 0.7rem;
	}

	.disk-rows {
		padding: 0.25rem 0;
	}
	.disk-row {
		display: flex;
		align-items: flex-start;
		gap: 1rem;
		padding: 0.65rem 1rem;
		border-bottom: 1px solid var(--border);
	}
	.disk-row:last-child {
		border-bottom: none;
	}
	.disk-label-col {
		min-width: 80px;
		flex-shrink: 0;
	}
	.disk-mount {
		display: block;
		font-family: "Courier Prime", monospace;
		font-size: 0.82rem;
		font-weight: 600;
		color: var(--text);
	}
	.disk-device {
		display: block;
		font-size: 0.66rem;
		color: var(--text-dim);
		font-family: "Courier Prime", monospace;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
		max-width: 80px;
	}
	.disk-bar-col {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 0.3rem;
	}
	.disk-bar {
		height: 6px;
		background: var(--accent);
		border-radius: 3px;
		overflow: hidden;
	}
	.disk-fill {
		height: 100%;
		border-radius: 3px;
		transition: width 0.5s;
	}
	.disk-stats-row {
		display: flex;
		justify-content: space-between;
		font-size: 0.7rem;
		font-family: "Courier Prime", monospace;
	}
	.disk-used {
		color: var(--text-dim);
	}
	.disk-pct {
		color: var(--text);
		font-weight: 700;
	}
	.disk-free {
		color: var(--accent);
	}
	.disk-total {
		color: var(--text-muted);
	}

	.log-list {
		max-height: 600px;
		overflow-y: auto;
		min-height: 80px;
	}
	.log-row {
		display: flex;
		gap: 0.6rem;
		padding: 0.35rem 1rem;
		border-bottom: 1px solid var(--border);
		font-size: 0.73rem;
		align-items: baseline;
	}
	.log-row:last-child {
		border-bottom: none;
	}
	.log-row.is-found {
		color: var(--accent);
	}
	.log-ts {
		color: var(--text-muted);
		font-family: "Courier Prime", monospace;
		font-size: 0.65rem;
		flex-shrink: 0;
	}
	.log-icon {
		flex-shrink: 0;
		width: 12px;
	}
	.log-msg {
		font-family: "Courier Prime", monospace;
		flex: 1;
	}

	@media (max-width: 600px) {
		.summary-grid {
			grid-template-columns: repeat(2, 1fr);
		}
		.phase-name {
			display: none;
		}
		.scan-hero {
			padding: 1rem;
		}
	}
</style>
