<script lang="ts">
	import { onMount } from "svelte";
	import { store, formatDuration, formatGB } from "$lib/stores/app.svelte";

	let history = $state<
		Array<{ file: string; timestamp: string; size: number }>
	>([]);
	let selectedScan = $state<any>(null);
	let previousScan = $state<any>(null);
	let loading = $state(false);

	onMount(() => {
		loadHistory();
	});

	async function loadHistory() {
		const r = await fetch(`${store.config.coreUrl}/api/scan/history`);
		if (r.ok) history = await r.json();
	}

	async function loadScan(file: string) {
		loading = true;
		selectedScan = null;
		previousScan = null;

		// Fetch current selection
		const r = await fetch(
			`${store.config.coreUrl}/api/scan/history/${file}`,
		);
		if (r.ok) selectedScan = await r.json();

		// Find and fetch the scan immediately before this one for comparison
		const idx = history.findIndex((h) => h.file === file);
		if (idx < history.length - 1) {
			const prevR = await fetch(
				`${store.config.coreUrl}/api/scan/history/${history[idx + 1].file}`,
			);
			if (prevR.ok) previousScan = await prevR.json();
		}

		loading = false;
	}

	function compareScans(oldScan: any, newScan: any) {
		// Return diff: new models, removed models, storage changes, etc.
		const oldModels = new Set(
			oldScan.models?.map((m: any) => m.file_path) || [],
		);
		const newModels = new Set(
			newScan.models?.map((m: any) => m.file_path) || [],
		);
		return {
			newModels: [...newModels].filter((p) => !oldModels.has(p)).length,
			removedModels: [...oldModels].filter((p) => !newModels.has(p))
				.length,
			storageDelta:
				(newScan.models?.reduce(
					(s: number, m: any) => s + m.size_bytes,
					0,
				) || 0) -
				(oldScan.models?.reduce(
					(s: number, m: any) => s + m.size_bytes,
					0,
				) || 0),
		};
	}
</script>

<div class="history-page">
	<!-- Go back -->
	<a href="/scan" class="mb-2 text-lg">← Back to scan</a>

	<!-- Header -->
	<div class="page-hdr">
		<span class="hdr-icon">◷</span>
		<div>
			<h1>Scan History</h1>
			<p class="hdr-sub">
				{history.length} scans archived · Trend analysis & change detection
			</p>
		</div>
	</div>

	<div class="history-layout">
		<!-- History list -->
		<div class="history-list">
			{#each history as item}
				<div class="hl-item" onclick={() => loadScan(item.file)}>
					<span class="hl-date"
						>{new Date(item.timestamp).toLocaleDateString()}</span
					>
					<span class="hl-time"
						>{new Date(item.timestamp).toLocaleTimeString()}</span
					>
					<span class="hl-size"
						>{(item.size / 1024).toFixed(0)} KB</span
					>
				</div>
			{/each}
		</div>

		<!-- Scan detail / comparison -->
		<div class="scan-detail">
			{#if loading}
				<div class="loading">Loading scan...</div>
			{:else if selectedScan}
				<div class="sd-header">
					<h2>
						Scan from {new Date(
							selectedScan.scanned_at,
						).toLocaleString()}
					</h2>
					<div class="sd-stats">
						<span>{selectedScan.models?.length || 0} models</span>
						<span
							>{selectedScan.components?.filter(
								(c: any) => c.status === "Running",
							).length || 0} running</span
						>
						<span
							>{formatGB(
								selectedScan.models?.reduce(
									(s: number, m: any) => s + m.size_bytes,
									0,
								) / 1e9,
							)}</span
						>
					</div>
				</div>

				<!-- Compare with previous -->
				{#if history.length > 1}
					{@const prevIdx = history.findIndex(
						(h) => h.file !== selectedScan.file,
					)}
					<div class="compare-section">
						<h3>Changes from previous scan</h3>
						{#if previousScan}
							{@const diff = compareScans(
								previousScan,
								selectedScan,
							)}
							<div class="diff-grid">
								<div class="diff-card new">
									+{diff.newModels} new models
								</div>
								<div class="diff-card removed">
									-{diff.removedModels} removed
								</div>
								<div class="diff-card storage">
									{diff.storageDelta > 0 ? "+" : ""}{(
										diff.storageDelta / 1e9
									).toFixed(2)} GB
								</div>
							</div>
						{:else}
							<p class="no-compare">
								No previous scan found for comparison.
							</p>
						{/if}
					</div>
				{/if}
			{:else}
				<div class="empty">Select a scan to view details</div>
			{/if}
		</div>
	</div>
</div>

<style>
	.history-page {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		padding: 1.25rem;
	}
	.page-hdr {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}
	.hdr-icon {
		font-size: 1.5rem;
		color: var(--accent);
	}
	h1 {
		font-family: "Special Elite", serif;
		font-size: 1.6rem;
		color: var(--text);
	}
	.hdr-sub {
		font-size: 1rem;
		color: var(--text-dim);
	}
	.history-layout {
		display: grid;
		grid-template-columns: 280px 1fr;
		gap: 1rem;
		min-height: 500px;
	}
	.history-list {
		background: var(--bg2);
		border: 1px solid var(--border);
		border-radius: 8px;
		overflow: auto;
		max-height: 600px;
	}
	.hl-item {
		display: grid;
		grid-template-columns: 1fr 1fr auto;
		gap: 0.5rem;
		padding: 0.6rem 1rem;
		border-bottom: 1px solid var(--border);
		cursor: pointer;
		font-size: 0.8rem;
	}
	.hl-item:hover {
		background: var(--accent-dim);
	}
	.hl-date {
		color: var(--text);
	}
	.hl-time {
		color: var(--text-dim);
	}
	.hl-size {
		font-family: "Courier Prime", monospace;
		color: var(--text-muted);
	}
	.scan-detail {
		background: var(--bg2);
		border: 1px solid var(--border);
		border-radius: 8px;
		padding: 1.25rem;
	}
	.sd-header h2 {
		font-family: "Special Elite", serif;
		font-size: 1rem;
		color: var(--text);
		margin-bottom: 0.75rem;
	}
	.sd-stats {
		display: flex;
		gap: 1rem;
		font-size: 0.78rem;
		color: var(--text-dim);
	}
	.sd-stats span {
		background: var(--bg3);
		padding: 2px 8px;
		border-radius: 3px;
	}
	.compare-section {
		margin-top: 1.5rem;
	}
	.compare-section h3 {
		font-size: 0.85rem;
		color: var(--accent);
		margin-bottom: 0.75rem;
		letter-spacing: 0.1em;
	}
	.diff-grid {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 0.75rem;
	}
	.diff-card {
		padding: 0.75rem;
		border-radius: 6px;
		text-align: center;
		font-weight: 700;
		font-size: 0.9rem;
	}
	.diff-card.new {
		background: rgba(74, 222, 128, 0.15);
		border: 1px solid #4ade80;
		color: #4ade80;
	}
	.diff-card.removed {
		background: rgba(239, 68, 68, 0.15);
		border: 1px solid #ef4444;
		color: #ef4444;
	}
	.diff-card.storage {
		background: var(--accent-dim);
		border: 1px solid var(--accent);
		color: var(--accent);
	}
	.loading,
	.empty {
		padding: 3rem;
		text-align: center;
		color: var(--text-dim);
	}
	@media (max-width: 768px) {
		.history-layout {
			grid-template-columns: 1fr;
		}
	}
</style>
