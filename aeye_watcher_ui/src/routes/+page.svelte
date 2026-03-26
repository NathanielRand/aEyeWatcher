<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import type { AgentCard } from '$lib/types';
  import {
    store,
    fetchOSInfo, pollGPU, startScan, requestDelete,
    setTheme, setSection, toggleSidebar,
    THEMES, formatBytes, formatGB
  } from '$lib/stores/app.svelte';

  // ---- Local state ----
  let gpuInterval: ReturnType<typeof setInterval>;
  let deleteConfirm = $state<string | null>(null);
  let deleteLoading = $state(false);
  let filterOrphaned = $state(false);
  let searchQuery = $state('');
  let selectedAgent = $state<AgentCard | null>(null);
  let activityLog = $state<{ time: string; msg: string; type: string }[]>([]);

  // Plain (non-reactive) cursor tracking which scanEvents we've consumed.
  // Not $state intentionally — writing it must NOT trigger $effect re-runs.
  let lastProcessedEventIdx = 0;

  // Phase labels for richer scan overlay header
  const phaseLabels: Record<string, string> = {
    system:     'SYSTEM PROFILING',
    gpu:        'GPU INTERROGATION',
    disk:       'STORAGE MAPPING',
    components: 'COMPONENT SWEEP',
    mcp:        'MCP CONFIG HUNT',
    acp:        'ACP PROBE',
    models:     'MODEL RECON',
    agents:     'AGENT PROFILING',
    complete:   'COMPILING REPORT',
  };

  // ---- Derived ----
  let filteredModels = $derived((() => {
    let models = store.scanResult?.models ?? [];
    if (filterOrphaned) models = models.filter(m => m.is_orphaned);
    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      models = models.filter(m =>
        m.name.toLowerCase().includes(q) || m.provider.toLowerCase().includes(q)
      );
    }
    return models;
  })());

  let runningComponents = $derived(
    (store.scanResult?.components ?? []).filter(c => c.status === 'Running')
  );

  let orphanedCount = $derived(
    (store.scanResult?.models ?? []).filter(m => m.is_orphaned).length
  );

  let totalModelGB = $derived(
    (store.scanResult?.models ?? []).reduce((s, m) => s + m.size_bytes, 0) / (1024 ** 3)
  );

  // ---- Lifecycle ----
  onMount(async () => {
    const saved = localStorage.getItem('aeye-theme');
    if (saved && saved in THEMES) setTheme(saved as keyof typeof THEMES);
    else setTheme('obsidian');

    await fetchOSInfo();
    gpuInterval = setInterval(pollGPU, 3000);
    pollGPU();
    handleScan();
  });

  onDestroy(() => clearInterval(gpuInterval));

  async function handleScan() {
    if (store.scanning) return;
    lastProcessedEventIdx = 0;
    addActivity('Initiating full system scan...', 'info');
    await startScan();
    addActivity('Scan complete. Intelligence report compiled.', 'success');
  }

  function addActivity(msg: string, type = 'info') {
    const time = new Date().toLocaleTimeString();
    activityLog = [{ time, msg, type }, ...activityLog].slice(0, 50);
  }

  // Reads only store.scanEvents.length (a number) as the reactive trigger.
  // Writes to lastProcessedEventIdx (plain let) and activityLog ($state).
  // activityLog is NOT read inside this effect, so no feedback loop.
  $effect(() => {
    const len = store.scanEvents.length;
    if (len > lastProcessedEventIdx) {
      const newEvents = store.scanEvents.slice(lastProcessedEventIdx, len);
      lastProcessedEventIdx = len;
      for (const evt of newEvents) {
        if (evt.type === 'found' || evt.type === 'checkpoint') {
          addActivity(evt.message, evt.type === 'found' ? 'found' : 'info');
        }
      }
    }
  });

  async function handleDelete(filePath: string) {
    deleteConfirm = filePath;
  }

  async function confirmDelete() {
    if (!deleteConfirm) return;
    deleteLoading = true;
    await requestDelete(deleteConfirm);
    addActivity(`Cleanup requested: ${deleteConfirm}`, 'warn');
    deleteConfirm = null;
    deleteLoading = false;
  }

  function componentsByType(types: string[]) {
    return (store.scanResult?.components ?? []).filter(c => types.includes(c.type));
  }

  function statusColor(status: string) {
    if (status === 'Running') return 'var(--accent)';
    if (status === 'Installed') return '#60a5fa';
    return '#374151';
  }

  function statusDot(status: string) {
    if (status === 'Running') return 'dot-running';
    if (status === 'Installed') return 'dot-installed';
    return 'dot-none';
  }

  function agentStatusColor(status: string) {
    if (status === 'detected') return 'var(--accent)';
    if (status === 'partial') return '#f59e0b';
    return '#374151';
  }

  const navItems = [
    { id: 'dashboard', label: 'Command Center', icon: '⬡' },
    { id: 'providers', label: 'Providers',       icon: '⚙' },
    { id: 'agents',   label: 'Agents',           icon: '◈' },
    { id: 'mcp',      label: 'MCP / ACP',        icon: '⬡' },
    { id: 'models',   label: 'Model Cache',      icon: '◆' },
    { id: 'storage',  label: 'Storage',          icon: '◉' },
    { id: 'activity', label: 'Activity Feed',    icon: '≡' },
  ];
</script>

<div class="app" data-theme={store.currentTheme}>
  <div class="grain"></div>

  <!-- ── Sidebar ── -->
  <aside class="sidebar" class:closed={!store.sidebarOpen}>
    <div class="sidebar-logo">
      <div class="logo-mark">
        <span class="logo-eye">◉</span>
        {#if store.sidebarOpen}
          <div>
            <div class="logo-name">aEye</div>
            <div class="logo-sub">WATCHER</div>
          </div>
        {/if}
      </div>
      <button class="toggle-btn" onclick={toggleSidebar}>
        {store.sidebarOpen ? '◀' : '▶'}
      </button>
    </div>

    {#if store.osInfo && store.sidebarOpen}
      <div class="os-badge">
        <span class="os-icon">
          {store.osInfo.os === 'linux' ? '🐧' : store.osInfo.os === 'darwin' ? '' : '🪟'}
        </span>
        <div class="os-info">
          <span class="os-name">{store.osInfo.distro ?? store.osInfo.os}</span>
          <span class="os-arch">{store.osInfo.arch} · {store.osInfo.hostname}</span>
        </div>
        {#if store.osInfo.gpu_available}
          <span class="gpu-badge">GPU</span>
        {/if}
      </div>
    {/if}

    <nav class="nav">
      {#each navItems as item}
        <button
          class="nav-item"
          class:active={store.activeSection === item.id}
          onclick={() => setSection(item.id)}
        >
          <span class="nav-icon">{item.icon}</span>
          {#if store.sidebarOpen}
            <span class="nav-label">{item.label}</span>
            {#if item.id === 'models' && orphanedCount > 0}
              <span class="nav-badge orphan">{orphanedCount}</span>
            {/if}
            {#if item.id === 'providers' && runningComponents.length > 0}
              <span class="nav-badge running">{runningComponents.length}</span>
            {/if}
          {/if}
        </button>
      {/each}
    </nav>

    {#if store.sidebarOpen}
      <div class="theme-switcher">
        <div class="theme-label">THEME</div>
        <div class="theme-pills">
          {#each Object.entries(THEMES) as [key, val]}
            <button
              class="theme-pill"
              class:active={store.currentTheme === key}
              style="--pill-color: {val.accent}"
              onclick={() => setTheme(key as keyof typeof THEMES)}
              title={val.label}
            ></button>
          {/each}
        </div>
      </div>
    {/if}
  </aside>

  <!-- ── Main ── -->
  <main class="main">
    <header class="topbar">
      <div class="topbar-left">
        <span class="section-title">
          {navItems.find(n => n.id === store.activeSection)?.label ?? 'Dashboard'}
        </span>
        {#if store.scanning}
          <span class="scanning-badge">● SCANNING</span>
        {:else if store.scanResult}
          <span class="scanned-badge">✓ SCANNED {new Date(store.scanResult.scanned_at).toLocaleTimeString()}</span>
        {/if}
      </div>
      <div class="topbar-right">
        {#if store.gpuInfo && !store.gpuInfo.error}
          <div class="gpu-widget">
            <span class="gpu-label">VRAM</span>
            <div class="gpu-bar-wrap">
              <div class="gpu-bar-fill" style="width: {store.gpuInfo.use_pct}%; background: {store.gpuInfo.use_pct > 85 ? '#ef4444' : store.gpuInfo.use_pct > 60 ? '#f59e0b' : 'var(--accent)'}"></div>
            </div>
            <span class="gpu-val">{store.gpuInfo.used_mb}/{store.gpuInfo.total_mb}MB</span>
          </div>
        {/if}
        <button class="scan-btn" onclick={handleScan} disabled={store.scanning}>
          {#if store.scanning}
            <span class="spin">⟳</span> SCANNING...
          {:else}
            ◉ SCAN
          {/if}
        </button>
      </div>
    </header>

    <!-- Scan Progress Overlay -->
    {#if store.scanning}
      <div class="scan-overlay">
        <div class="scan-panel">
          <div class="scan-header">
            <div class="scan-title-row">
              <span class="scan-pulse">◉</span>
              <span class="scan-title">ACTIVE SCAN IN PROGRESS</span>
            </div>
            <span class="scan-phase-badge">{phaseLabels[store.scanPhase] ?? store.scanPhase.toUpperCase()}</span>
          </div>

          <div class="scan-phase-steps">
            {#each ['system','gpu','disk','components','mcp','acp','models','agents'] as ph}
              {@const isDone = store.scanEvents.some(e => e.phase === ph && (e.type === 'found' || e.type === 'complete'))}
              {@const isActive = store.scanPhase === ph}
              <div class="phase-step" class:done={isDone} class:active={isActive}>
                <span class="phase-step-dot">{isDone ? '✓' : isActive ? '◉' : '○'}</span>
                <span class="phase-step-label">{ph}</span>
              </div>
            {/each}
          </div>

          <div class="progress-track">
            <div class="progress-fill" style="width: {store.scanPercent}%">
              <div class="progress-shimmer"></div>
            </div>
          </div>
          <div class="progress-meta">
            <span class="progress-pct-label">{store.scanPercent}%</span>
            <span class="progress-phase-msg">
              {store.scanEvents.length > 0 ? store.scanEvents[store.scanEvents.length - 1].message : 'Initializing...'}
            </span>
          </div>

          <div class="scan-feed">
            <div class="scan-feed-header">LIVE FEED</div>
            {#each [...store.scanEvents].reverse().slice(0, 10) as evt}
              <div class="feed-line" class:checkpoint={evt.type === 'checkpoint'} class:found={evt.type === 'found'} class:progress={evt.type === 'progress'}>
                <span class="feed-ts">{String(store.scanEvents.indexOf(evt)).padStart(3,'0')}</span>
                <span class="feed-icon">{evt.type === 'checkpoint' ? '▸' : evt.type === 'found' ? '◆' : '·'}</span>
                <span class="feed-phase">[{evt.phase}]</span>
                <span class="feed-msg">{evt.message}</span>
              </div>
            {/each}
            {#if store.scanEvents.length === 0}
              <div class="feed-line"><span class="feed-icon">·</span><span class="feed-msg">Waiting for scan data...</span></div>
            {/if}
          </div>
        </div>
      </div>
    {/if}

    <!-- ── Content ── -->
    <div class="content">

      <!-- DASHBOARD -->
      {#if store.activeSection === 'dashboard'}
        <div class="dashboard">
          <div class="stat-grid">
            <div class="stat-card">
              <div class="stat-icon">⚙</div>
              <div class="stat-val">{runningComponents.length}</div>
              <div class="stat-label">Live Services</div>
            </div>
            <div class="stat-card">
              <div class="stat-icon">◆</div>
              <div class="stat-val">{store.scanResult?.models?.length ?? '—'}</div>
              <div class="stat-label">Models Found</div>
            </div>
            <div class="stat-card" class:warn={orphanedCount > 0}>
              <div class="stat-icon">⚠</div>
              <div class="stat-val">{orphanedCount}</div>
              <div class="stat-label">Orphaned</div>
            </div>
            <div class="stat-card">
              <div class="stat-icon">⬡</div>
              <div class="stat-val">{store.scanResult?.mcp_servers?.length ?? '—'}</div>
              <div class="stat-label">MCP Servers</div>
            </div>
            <div class="stat-card">
              <div class="stat-icon">◈</div>
              <div class="stat-val">{store.scanResult?.agent_cards?.filter(a => a.status !== 'not_found').length ?? '—'}</div>
              <div class="stat-label">Agents Detected</div>
            </div>
            <div class="stat-card">
              <div class="stat-icon">◉</div>
              <div class="stat-val">{totalModelGB.toFixed(1)} GB</div>
              <div class="stat-label">Model Storage</div>
            </div>
          </div>

          {#if runningComponents.length > 0}
            <div class="section-card">
              <div class="card-header">
                <span>⚙ LIVE PROVIDERS</span>
                <span class="card-count">{runningComponents.length}</span>
              </div>
              <div class="provider-chips">
                {#each runningComponents as comp}
                  <div class="provider-chip running">
                    <span class="dot-running"></span>
                    {comp.name}
                    {#if comp.port}<span class="chip-port">:{comp.port}</span>{/if}
                  </div>
                {/each}
              </div>
            </div>
          {/if}

          {#if store.scanResult?.disks && store.scanResult.disks.length > 0}
            <div class="section-card">
              <div class="card-header">
                <span>◉ STORAGE VOLUMES</span>
                <span class="card-count">{store.scanResult.disks.length}</span>
              </div>
              {#each store.scanResult.disks as disk}
                <div class="disk-row">
                  <div class="disk-left">
                    <span class="disk-mount">{disk.mount}</span>
                    <span class="disk-device">{disk.device}</span>
                  </div>
                  <div class="disk-bar-wrap">
                    <div class="disk-bar-fill" style="width: {disk.use_pct}%; background: {disk.use_pct > 90 ? '#ef4444' : disk.use_pct > 75 ? '#f59e0b' : 'var(--accent)'}"></div>
                  </div>
                  <div class="disk-right">
                    <span class="disk-pct">{disk.use_pct.toFixed(0)}%</span>
                    <span class="disk-free">{formatGB(disk.free_gb)} free</span>
                  </div>
                </div>
              {/each}
            </div>
          {/if}

          <div class="section-card">
            <div class="card-header"><span>≡ RECENT ACTIVITY</span></div>
            <div class="activity-list">
              {#each activityLog.slice(0, 10) as entry}
                <div class="activity-row" class:act-found={entry.type === 'found'} class:act-warn={entry.type === 'warn'} class:act-success={entry.type === 'success'}>
                  <span class="act-time">{entry.time}</span>
                  <span class="act-msg">{entry.msg}</span>
                </div>
              {/each}
              {#if activityLog.length === 0}
                <div class="empty-state">No activity yet. Run a scan.</div>
              {/if}
            </div>
          </div>
        </div>

      <!-- PROVIDERS -->
      {:else if store.activeSection === 'providers'}
        <div class="providers-view">
          {#each [
            { label: 'Core Providers',    icon: '⚙', types: ['Provider'] },
            { label: 'Routing & Gateways',icon: '▸', types: ['Gateway'] },
            { label: 'Vector Databases',  icon: '◆', types: ['Vector DB'] },
            { label: 'GUI Clients',       icon: '◉', types: ['GUI Client'] },
            { label: 'Voice & TTS',       icon: '♪', types: ['TTS / Audio'] },
            { label: 'Observability',     icon: '◈', types: ['Observability'] },
          ] as group}
            {@const comps = componentsByType(group.types)}
            {@const active = comps.filter(c => c.status !== 'Not Found')}
            {#if active.length > 0}
              <div class="section-card">
                <div class="card-header">
                  <span>{group.icon} {group.label}</span>
                  <span class="card-count">{active.length} / {comps.length}</span>
                </div>
                <div class="comp-list">
                  {#each active as comp}
                    <div class="comp-row">
                      <div class="comp-dot {statusDot(comp.status)}"></div>
                      <div class="comp-info">
                        <span class="comp-name">{comp.name}</span>
                        {#if comp.description}<span class="comp-desc">{comp.description}</span>{/if}
                        {#if comp.version}<span class="comp-version">{comp.version}</span>{/if}
                        {#if comp.path || comp.url}<span class="comp-path">{comp.url ?? comp.path}</span>{/if}
                      </div>
                      <div class="comp-status" style="color: {statusColor(comp.status)}">
                        {comp.status}
                        {#if comp.port}<span class="comp-port">:{comp.port}</span>{/if}
                      </div>
                    </div>
                  {/each}
                </div>
              </div>
            {/if}
          {/each}
          {#if !store.scanResult}
            <div class="empty-state">Run a scan to detect providers.</div>
          {/if}
        </div>

      <!-- AGENTS -->
      {:else if store.activeSection === 'agents'}
        <div class="agents-view">
          {#if selectedAgent}
            <div class="agent-detail">
              <button class="back-btn" onclick={() => selectedAgent = null}>← Back</button>
              <div class="agent-card-header">
                <span class="agent-big-icon">{selectedAgent.icon}</span>
                <div>
                  <div class="agent-big-name">{selectedAgent.name}</div>
                  <div class="agent-ecosystem">{selectedAgent.ecosystem}</div>
                  <div class="agent-desc">{selectedAgent.description}</div>
                </div>
                <div class="agent-status-badge" style="color: {agentStatusColor(selectedAgent.status)}">
                  {selectedAgent.status.toUpperCase()}
                </div>
              </div>

              <div class="section-card">
                <div class="card-header"><span>◈ EXPECTED ECOSYSTEM FILES</span></div>
                <div class="expected-grid">
                  {#each selectedAgent.expected_files ?? [] as ef}
                    <div class="expected-item" class:found={ef.found}>
                      <span class="expected-check">{ef.found ? '✓' : '○'}</span>
                      <div class="expected-info">
                        <span class="expected-label">{ef.label}</span>
                        {#if ef.path}<span class="expected-path">{ef.path}</span>{/if}
                      </div>
                    </div>
                  {/each}
                </div>
              </div>

              {#if selectedAgent.files?.length > 0}
                <div class="section-card">
                  <div class="card-header"><span>◆ DISCOVERED FILES</span></div>
                  {#each selectedAgent.files as f}
                    <div class="agent-file-row">
                      <span class="file-type-badge">{f.type}</span>
                      <span class="file-name">{f.name}</span>
                      <span class="file-size">{formatBytes(f.size_bytes)}</span>
                    </div>
                  {/each}
                </div>
              {/if}
            </div>
          {:else}
            <div class="agent-grid">
              {#each store.scanResult?.agent_cards ?? [] as card}
                {@const foundCount = card.expected_files?.filter(f => f.found).length ?? 0}
                {@const totalCount = card.expected_files?.length ?? 1}
                <button
                  class="agent-card"
                  class:detected={card.status === 'detected'}
                  class:partial={card.status === 'partial'}
                  onclick={() => selectedAgent = card}
                >
                  <div class="ac-top">
                    <span class="ac-icon">{card.icon}</span>
                    <div class="ac-status-dot" style="background: {agentStatusColor(card.status)}"></div>
                  </div>
                  <div class="ac-name">{card.name}</div>
                  <div class="ac-eco">{card.ecosystem}</div>
                  <div class="ac-progress">
                    <div class="ac-bar">
                      <div class="ac-bar-fill" style="width: {(foundCount / totalCount) * 100}%"></div>
                    </div>
                    <span class="ac-found">{foundCount}/{totalCount}</span>
                  </div>
                  <div class="ac-status-label" style="color: {agentStatusColor(card.status)}">
                    {card.status.replace('_', ' ').toUpperCase()}
                  </div>
                </button>
              {/each}
              {#if !store.scanResult}
                <div class="empty-state">Run a scan to profile agent ecosystems.</div>
              {/if}
            </div>
          {/if}
        </div>

      <!-- MCP / ACP -->
      {:else if store.activeSection === 'mcp'}
        <div class="mcp-view">
          {#if (store.scanResult?.mcp_servers?.length ?? 0) > 0}
            <div class="section-card">
              <div class="card-header">
                <span>⬡ MCP SERVERS</span>
                <span class="card-count">{store.scanResult?.mcp_servers?.length}</span>
              </div>
              <div class="comp-list">
                {#each store.scanResult?.mcp_servers ?? [] as srv}
                  <div class="comp-row">
                    <div class="comp-dot {srv.status === 'running' ? 'dot-running' : srv.status === 'configured' ? 'dot-installed' : 'dot-none'}"></div>
                    <div class="comp-info">
                      <span class="comp-name">{srv.name}</span>
                      {#if srv.url}<span class="comp-path">{srv.url}</span>{/if}
                      {#if srv.config}<span class="comp-desc">Config: {srv.config}</span>{/if}
                    </div>
                    <div class="comp-right">
                      <span class="transport-badge">{srv.transport}</span>
                      <span class="comp-status" style="color: {srv.status === 'running' ? 'var(--accent)' : '#60a5fa'}">{srv.status}</span>
                    </div>
                  </div>
                {/each}
              </div>
            </div>
          {:else}
            <div class="section-card">
              <div class="empty-state">No MCP server configurations detected.<br><span class="hint">Configure MCP servers in your Claude Desktop, Cursor, or Continue config files.</span></div>
            </div>
          {/if}

          {#if (store.scanResult?.acp_agents?.length ?? 0) > 0}
            <div class="section-card">
              <div class="card-header">
                <span>◈ ACP AGENTS</span>
                <span class="card-count">{store.scanResult?.acp_agents?.length}</span>
              </div>
              <div class="comp-list">
                {#each store.scanResult?.acp_agents ?? [] as agent}
                  <div class="comp-row">
                    <div class="comp-dot dot-running"></div>
                    <div class="comp-info">
                      <span class="comp-name">{agent.name}</span>
                      <span class="comp-path">{agent.endpoint}</span>
                    </div>
                    <div class="comp-right">
                      <span class="transport-badge">{agent.protocol}</span>
                      <span class="comp-status" style="color: var(--accent)">{agent.status}</span>
                    </div>
                  </div>
                {/each}
              </div>
            </div>
          {:else if store.scanResult}
            <div class="section-card">
              <div class="empty-state">No ACP agents detected.<br><span class="hint">ACP agents are discovered via /.well-known/agent.json</span></div>
            </div>
          {/if}
        </div>

      <!-- MODELS -->
      {:else if store.activeSection === 'models'}
        <div class="models-view">
          <div class="models-toolbar">
            <input class="search-input" type="text" placeholder="Search models..." bind:value={searchQuery} />
            <label class="filter-toggle">
              <input type="checkbox" bind:checked={filterOrphaned} />
              <span>Orphaned only ({orphanedCount})</span>
            </label>
            <span class="model-total">{filteredModels.length} model(s) · {totalModelGB.toFixed(1)} GB total</span>
          </div>

          <div class="model-list">
            {#each filteredModels as model}
              <div class="model-row" class:orphaned={model.is_orphaned}>
                <div class="model-left">
                  <div class="model-name">{model.name}</div>
                  <div class="model-path">{model.file_path}</div>
                </div>
                <div class="model-right">
                  {#if model.is_orphaned}<span class="orphan-badge">ORPHANED</span>{/if}
                  <span class="format-badge">{model.format ?? model.extension}</span>
                  <span class="provider-badge">{model.provider}</span>
                  <span class="model-size">{formatBytes(model.size_bytes)}</span>
                  <button class="delete-btn" onclick={() => handleDelete(model.file_path)}>🗑</button>
                </div>
              </div>
            {/each}
            {#if filteredModels.length === 0}
              <div class="empty-state">{store.scanResult ? 'No models match your filter.' : 'Run a scan to discover model files.'}</div>
            {/if}
          </div>

          {#if (store.scanResult?.model_dirs?.length ?? 0) > 0}
            <div class="section-card" style="margin-top: 1.5rem">
              <div class="card-header"><span>◉ KNOWN CACHE DIRECTORIES</span></div>
              {#each store.scanResult?.model_dirs ?? [] as dir}
                <div class="dir-row"><code>{dir}</code></div>
              {/each}
            </div>
          {/if}
        </div>

      <!-- STORAGE -->
      {:else if store.activeSection === 'storage'}
        <div class="storage-view">
          {#each store.scanResult?.disks ?? [] as disk}
            <div class="section-card disk-card">
              <div class="disk-header">
                <div>
                  <div class="disk-mount-big">{disk.mount}</div>
                  <div class="disk-meta">{disk.device} · {disk.fs_type}</div>
                </div>
                <div class="disk-pct-big" style="color: {disk.use_pct > 90 ? '#ef4444' : disk.use_pct > 75 ? '#f59e0b' : 'var(--accent)'}">
                  {disk.use_pct.toFixed(1)}%
                </div>
              </div>
              <div class="disk-bar-big-wrap">
                <div class="disk-bar-big-fill" style="width: {disk.use_pct}%; background: {disk.use_pct > 90 ? '#ef4444' : disk.use_pct > 75 ? '#f59e0b' : 'var(--accent)'}"></div>
              </div>
              <div class="disk-stats">
                <div class="disk-stat">
                  <span class="ds-val">{formatGB(disk.used_gb)}</span>
                  <span class="ds-label">Used</span>
                </div>
                <div class="disk-stat">
                  <span class="ds-val">{formatGB(disk.free_gb)}</span>
                  <span class="ds-label">Free</span>
                </div>
                <div class="disk-stat">
                  <span class="ds-val">{formatGB(disk.total_gb)}</span>
                  <span class="ds-label">Total</span>
                </div>
              </div>
            </div>
          {/each}
          {#if !store.scanResult}
            <div class="empty-state">Run a scan to map storage volumes.</div>
          {/if}
        </div>

      <!-- ACTIVITY -->
      {:else if store.activeSection === 'activity'}
        <div class="activity-view">
          <div class="section-card">
            <div class="card-header">
              <span>≡ FULL ACTIVITY LOG</span>
              <button class="clear-btn" onclick={() => activityLog = []}>Clear</button>
            </div>
            <div class="activity-log-full">
              {#each activityLog as entry}
                <div class="activity-row" class:act-found={entry.type === 'found'} class:act-warn={entry.type === 'warn'} class:act-success={entry.type === 'success'}>
                  <span class="act-time">{entry.time}</span>
                  <span class="act-bullet">◆</span>
                  <span class="act-msg">{entry.msg}</span>
                </div>
              {/each}
              {#if activityLog.length === 0}
                <div class="empty-state">No activity yet.</div>
              {/if}
            </div>
          </div>
        </div>
      {/if}

    </div><!-- /content -->
  </main>

  <!-- Delete confirm modal -->
  {#if deleteConfirm}
    <div class="modal-overlay" role="presentation" onclick={() => deleteConfirm = null}>
      <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
      <div class="modal" onclick={(e) => e.stopPropagation()}>
        <div class="modal-title">⚠ CONFIRM CLEANUP</div>
        <div class="modal-body">
          <code>{deleteConfirm}</code>
          <p>This logs a cleanup request. SAFE MODE is active by default — no files are deleted until you enable it in the server config.</p>
        </div>
        <div class="modal-actions">
          <button class="modal-cancel" onclick={() => deleteConfirm = null}>Cancel</button>
          <button class="modal-confirm" onclick={confirmDelete} disabled={deleteLoading}>
            {deleteLoading ? 'Requesting...' : 'Request Cleanup'}
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  @import url('https://fonts.googleapis.com/css2?family=Special+Elite&family=Courier+Prime:wght@400;700&family=Rajdhani:wght@400;500;600;700&display=swap');

  :global([data-theme="obsidian"]) {
    --accent: #c084fc; --accent-dim: rgba(192,132,252,0.15);
    --bg: #08080f; --bg2: #0f0f1a; --bg3: #16162a;
    --border: rgba(192,132,252,0.12);
    --text: #e2e0f0; --text-dim: #6b6a85; --text-muted: #38374d;
  }
  :global([data-theme="clay"]) {
    --accent: #fb923c; --accent-dim: rgba(251,146,60,0.12);
    --bg: #0f0905; --bg2: #1a1008; --bg3: #241608;
    --border: rgba(251,146,60,0.1);
    --text: #f0e8e0; --text-dim: #85705a; --text-muted: #3d2e1e;
  }
  :global([data-theme="forest"]) {
    --accent: #4ade80; --accent-dim: rgba(74,222,128,0.1);
    --bg: #05100a; --bg2: #081a0e; --bg3: #0c2414;
    --border: rgba(74,222,128,0.1);
    --text: #e0f0e4; --text-dim: #5a8560; --text-muted: #1e3d24;
  }
  :global([data-theme="slate"]) {
    --accent: #7dd3fc; --accent-dim: rgba(125,211,252,0.1);
    --bg: #06090f; --bg2: #0a1020; --bg3: #0e1830;
    --border: rgba(125,211,252,0.1);
    --text: #e0ecf8; --text-dim: #5a6a85; --text-muted: #1e2840;
  }
  :global([data-theme="amber"]) {
    --accent: #fbbf24; --accent-dim: rgba(251,191,36,0.1);
    --bg: #0a0800; --bg2: #150f00; --bg3: #1f1600;
    --border: rgba(251,191,36,0.1);
    --text: #f8f0d0; --text-dim: #856a2a; --text-muted: #3d2e00;
  }

  :global(*, *::before, *::after) { box-sizing: border-box; margin: 0; padding: 0; }
  :global(body) {
    background: var(--bg); color: var(--text);
    font-family: 'Rajdhani', 'Courier Prime', monospace;
    overflow: hidden; height: 100vh;
  }

  .grain {
    position: fixed; inset: 0; pointer-events: none; z-index: 9999;
    background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)' opacity='0.04'/%3E%3C/svg%3E");
    background-size: 200px 200px; opacity: 0.6; mix-blend-mode: overlay;
  }

  .app { display: flex; height: 100vh; overflow: hidden; position: relative; }

  /* Sidebar */
  .sidebar {
    width: 220px; min-width: 220px;
    background: var(--bg2); border-right: 1px solid var(--border);
    display: flex; flex-direction: column;
    transition: width 0.2s, min-width 0.2s;
    overflow: hidden; z-index: 100;
  }
  .sidebar.closed { width: 52px; min-width: 52px; }
  .sidebar-logo {
    display: flex; align-items: center; justify-content: space-between;
    padding: 1rem 0.75rem; border-bottom: 1px solid var(--border);
  }
  .logo-mark { display: flex; align-items: center; gap: 0.6rem; }
  .logo-eye {
    font-size: 1.4rem; color: var(--accent);
    animation: pulse-eye 3s ease-in-out infinite;
  }
  @keyframes pulse-eye { 0%,100%{opacity:1;transform:scale(1)} 50%{opacity:.7;transform:scale(.95)} }
  .logo-name { font-family:'Special Elite',serif; font-size:1.1rem; color:var(--accent); letter-spacing:.1em; line-height:1; }
  .logo-sub { font-size:.55rem; letter-spacing:.35em; color:var(--text-dim); }
  .toggle-btn { background:none; border:none; color:var(--text-dim); cursor:pointer; font-size:.7rem; padding:2px 4px; transition:color .2s; }
  .toggle-btn:hover { color:var(--accent); }

  .os-badge { display:flex; align-items:center; gap:.5rem; padding:.6rem .75rem; background:var(--accent-dim); border-bottom:1px solid var(--border); font-size:.7rem; }
  .os-icon { font-size:1rem; }
  .os-info { display:flex; flex-direction:column; flex:1; }
  .os-name { color:var(--accent); font-weight:600; font-size:.7rem; }
  .os-arch { color:var(--text-dim); font-size:.6rem; }
  .gpu-badge { background:var(--accent); color:var(--bg); font-size:.55rem; font-weight:700; padding:1px 4px; border-radius:2px; letter-spacing:.1em; }

  .nav { flex:1; display:flex; flex-direction:column; gap:2px; padding:.5rem; }
  .nav-item {
    display:flex; align-items:center; gap:.6rem; padding:.5rem .6rem;
    background:none; border:none; border-radius:4px; color:var(--text-dim);
    cursor:pointer; font-family:'Rajdhani',monospace; font-size:.8rem; font-weight:500;
    letter-spacing:.05em; text-align:left; transition:all .15s; width:100%;
  }
  .nav-item:hover { background:var(--accent-dim); color:var(--text); }
  .nav-item.active { background:var(--accent-dim); color:var(--accent); border-left:2px solid var(--accent); }
  .nav-icon { font-size:.9rem; width:16px; text-align:center; flex-shrink:0; }
  .nav-label { flex:1; }
  .nav-badge { font-size:.6rem; padding:1px 5px; border-radius:2px; font-weight:700; }
  .nav-badge.running { background:rgba(74,222,128,.2); color:#4ade80; }
  .nav-badge.orphan { background:rgba(251,146,60,.2); color:#fb923c; }

  .theme-switcher { padding:.75rem; border-top:1px solid var(--border); }
  .theme-label { font-size:.55rem; letter-spacing:.2em; color:var(--text-muted); margin-bottom:.4rem; }
  .theme-pills { display:flex; gap:6px; }
  .theme-pill { width:16px; height:16px; border-radius:50%; border:2px solid transparent; background:var(--pill-color); cursor:pointer; transition:transform .15s, border-color .15s; }
  .theme-pill:hover { transform:scale(1.2); }
  .theme-pill.active { border-color:var(--text); transform:scale(1.15); }

  /* Main */
  .main { flex:1; display:flex; flex-direction:column; overflow:hidden; min-width:0; }

  .topbar { display:flex; align-items:center; justify-content:space-between; padding:.6rem 1.25rem; background:var(--bg2); border-bottom:1px solid var(--border); flex-shrink:0; }
  .topbar-left,.topbar-right { display:flex; align-items:center; gap:.75rem; }
  .section-title { font-family:'Special Elite',serif; font-size:.9rem; color:var(--text); letter-spacing:.05em; }
  .scanning-badge { font-size:.65rem; color:var(--accent); letter-spacing:.15em; animation:blink 1s steps(1) infinite; }
  @keyframes blink { 50%{opacity:0} }
  .scanned-badge { font-size:.65rem; color:var(--text-dim); letter-spacing:.1em; }

  .gpu-widget { display:flex; align-items:center; gap:.5rem; font-size:.65rem; }
  .gpu-label { color:var(--text-dim); letter-spacing:.1em; }
  .gpu-bar-wrap { width:80px; height:4px; background:var(--bg3); border-radius:2px; overflow:hidden; }
  .gpu-bar-fill { height:100%; transition:width .5s,background .5s; border-radius:2px; }
  .gpu-val { color:var(--text-dim); font-family:'Courier Prime',monospace; font-size:.6rem; }

  .scan-btn {
    background:var(--accent-dim); border:1px solid var(--accent); color:var(--accent);
    padding:.35rem .85rem; border-radius:3px; cursor:pointer;
    font-family:'Rajdhani',monospace; font-size:.75rem; font-weight:700; letter-spacing:.1em;
    transition:all .15s; display:flex; align-items:center; gap:.4rem;
  }
  .scan-btn:hover:not(:disabled) { background:var(--accent); color:var(--bg); }
  .scan-btn:disabled { opacity:.6; cursor:not-allowed; }
  .spin { display:inline-block; animation:spin .8s linear infinite; }
  @keyframes spin { to{transform:rotate(360deg)} }

  /* Scan overlay */
  .scan-overlay {
    position:absolute; inset:0; background:rgba(0,0,0,.7);
    display:flex; align-items:center; justify-content:center;
    z-index:200; backdrop-filter:blur(4px);
  }
  .scan-panel {
    background:var(--bg2); border:1px solid var(--accent); border-radius:8px;
    padding:1.5rem; width:min(600px,90vw);
    box-shadow:0 0 40px var(--accent-dim),0 0 80px rgba(0,0,0,.8);
  }
  /* --- SCAN OVERLAY (improved) --- */
  .scan-header { display:flex; justify-content:space-between; align-items:center; margin-bottom:1rem; }
  .scan-title-row { display:flex; align-items:center; gap:.5rem; }
  .scan-pulse { color:var(--accent); font-size:1rem; animation:glow-pulse 1s ease-in-out infinite; }
  .scan-title { color:var(--accent); font-family:'Special Elite',serif; font-size:.9rem; letter-spacing:.1em; }
  .scan-phase-badge { font-size:.6rem; letter-spacing:.2em; color:var(--bg); background:var(--accent); padding:2px 8px; border-radius:2px; font-weight:700; }

  .scan-phase-steps { display:flex; gap:0; margin-bottom:1rem; border:1px solid var(--border); border-radius:4px; overflow:hidden; }
  .phase-step { flex:1; display:flex; flex-direction:column; align-items:center; padding:5px 2px; gap:2px; border-right:1px solid var(--border); background:var(--bg3); transition:background .3s; }
  .phase-step:last-child { border-right:none; }
  .phase-step.done { background:color-mix(in srgb,var(--accent) 12%,var(--bg3)); }
  .phase-step.active { background:color-mix(in srgb,var(--accent) 20%,var(--bg3)); }
  .phase-step-dot { font-size:.65rem; color:var(--text-muted); line-height:1; }
  .phase-step.done .phase-step-dot { color:var(--accent); }
  .phase-step.active .phase-step-dot { color:var(--accent); animation:glow-pulse 1s ease-in-out infinite; }
  .phase-step-label { font-size:.45rem; letter-spacing:.1em; color:var(--text-muted); text-transform:uppercase; }
  .phase-step.done .phase-step-label, .phase-step.active .phase-step-label { color:var(--accent); }

  .progress-track { background:var(--bg3); height:10px; border-radius:5px; overflow:hidden; margin-bottom:.4rem; position:relative; }
  .progress-fill { height:100%; background:var(--accent); transition:width .4s ease; border-radius:5px; box-shadow:0 0 12px var(--accent); position:relative; overflow:hidden; }
  .progress-shimmer { position:absolute; inset:0; background:linear-gradient(90deg,transparent 0%,rgba(255,255,255,.25) 50%,transparent 100%); animation:shimmer 1.5s ease-in-out infinite; }
  @keyframes shimmer { 0%{transform:translateX(-100%)} 100%{transform:translateX(200%)} }
  .progress-meta { display:flex; justify-content:space-between; align-items:baseline; margin-bottom:.75rem; }
  .progress-pct-label { font-family:'Courier Prime',monospace; font-size:.75rem; color:var(--accent); font-weight:700; }
  .progress-phase-msg { font-size:.65rem; color:var(--text-dim); font-family:'Courier Prime',monospace; max-width:75%; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }

  .scan-feed { border:1px solid var(--border); border-radius:4px; overflow:hidden; }
  .scan-feed-header { padding:3px 8px; background:var(--bg3); font-size:.55rem; letter-spacing:.2em; color:var(--text-muted); border-bottom:1px solid var(--border); }
  .feed-line { display:flex; gap:.4rem; font-size:.65rem; color:var(--text-muted); padding:3px 8px; border-bottom:1px solid color-mix(in srgb,var(--border) 50%,transparent); align-items:baseline; animation:fadeSlide .15s ease; }
  .feed-line:last-child { border-bottom:none; }
  @keyframes fadeSlide { from{opacity:0;transform:translateY(-3px)} to{opacity:1;transform:translateY(0)} }
  .feed-line.checkpoint { color:var(--text-dim); background:color-mix(in srgb,var(--accent) 5%,transparent); }
  .feed-line.found { color:var(--accent); }
  .feed-line.progress { color:var(--text-muted); }
  .feed-ts { color:var(--text-muted); font-family:'Courier Prime',monospace; font-size:.55rem; flex-shrink:0; opacity:.5; }
  .feed-icon { flex-shrink:0; width:10px; }
  .feed-phase { color:var(--text-muted); font-family:'Courier Prime',monospace; font-size:.6rem; flex-shrink:0; opacity:.6; }
  .feed-msg { font-family:'Courier Prime',monospace; flex:1; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }

  /* Content */
  .content { flex:1; overflow-y:auto; padding:1.25rem; scrollbar-width:thin; scrollbar-color:var(--border) transparent; }
  .content::-webkit-scrollbar { width:4px; }
  .content::-webkit-scrollbar-thumb { background:var(--border); border-radius:2px; }

  /* Cards */
  .section-card { background:var(--bg2); border:1px solid var(--border); border-radius:6px; margin-bottom:1rem; overflow:hidden; }
  .card-header { display:flex; justify-content:space-between; align-items:center; padding:.65rem 1rem; border-bottom:1px solid var(--border); font-size:.7rem; letter-spacing:.15em; color:var(--accent); font-family:'Rajdhani',monospace; font-weight:700; background:var(--accent-dim); }
  .card-count { background:var(--bg3); border:1px solid var(--border); padding:1px 6px; border-radius:2px; color:var(--text-dim); font-size:.65rem; }
  .empty-state { padding:2rem; text-align:center; color:var(--text-muted); font-size:.8rem; font-style:italic; line-height:1.8; }
  .hint { font-size:.7rem; color:var(--text-dim); }

  /* Dots */
  :global(.dot-running) { width:8px; height:8px; border-radius:50%; background:var(--accent); box-shadow:0 0 6px var(--accent); flex-shrink:0; animation:glow-pulse 2s ease-in-out infinite; }
  @keyframes glow-pulse { 0%,100%{opacity:1} 50%{opacity:.6} }
  :global(.dot-installed) { width:8px; height:8px; border-radius:50%; background:#60a5fa; flex-shrink:0; }
  :global(.dot-none) { width:8px; height:8px; border-radius:50%; background:var(--text-muted); flex-shrink:0; }

  /* Dashboard */
  .dashboard { display:flex; flex-direction:column; gap:0; }
  .stat-grid { display:grid; grid-template-columns:repeat(auto-fill,minmax(160px,1fr)); gap:.75rem; margin-bottom:1rem; }
  .stat-card { background:var(--bg2); border:1px solid var(--border); border-radius:6px; padding:1rem; text-align:center; position:relative; overflow:hidden; transition:border-color .2s; }
  .stat-card::before { content:''; position:absolute; top:0; left:0; right:0; height:2px; background:var(--accent); opacity:.5; }
  .stat-card.warn::before { background:#fb923c; }
  .stat-card:hover { border-color:var(--accent); }
  .stat-icon { font-size:1rem; color:var(--accent); margin-bottom:.4rem; opacity:.7; }
  .stat-val { font-family:'Special Elite',serif; font-size:1.8rem; color:var(--text); line-height:1; }
  .stat-label { font-size:.65rem; color:var(--text-dim); letter-spacing:.1em; margin-top:.3rem; }

  .provider-chips { display:flex; flex-wrap:wrap; gap:.5rem; padding:.75rem 1rem; }
  .provider-chip { display:flex; align-items:center; gap:.4rem; padding:.3rem .6rem; border-radius:3px; font-size:.75rem; background:var(--bg3); border:1px solid var(--border); }
  .provider-chip.running { border-color:var(--accent); }
  .chip-port { color:var(--text-dim); font-family:'Courier Prime',monospace; font-size:.65rem; }

  .disk-row { display:flex; align-items:center; gap:.75rem; padding:.6rem 1rem; border-bottom:1px solid var(--border); }
  .disk-row:last-child { border-bottom:none; }
  .disk-left { min-width:120px; }
  .disk-mount { font-size:.8rem; font-weight:600; display:block; }
  .disk-device { font-size:.65rem; color:var(--text-dim); font-family:'Courier Prime',monospace; }
  .disk-bar-wrap { flex:1; height:4px; background:var(--bg3); border-radius:2px; overflow:hidden; }
  .disk-bar-fill { height:100%; border-radius:2px; transition:width .5s; }
  .disk-right { min-width:120px; text-align:right; }
  .disk-pct { font-size:.75rem; color:var(--text); font-family:'Courier Prime',monospace; }
  .disk-free { font-size:.65rem; color:var(--text-dim); margin-left:.5rem; }

  .activity-list { padding:.5rem 0; }
  .activity-row { display:flex; gap:.75rem; padding:.4rem 1rem; font-size:.7rem; border-bottom:1px solid var(--border); align-items:baseline; }
  .activity-row:last-child { border-bottom:none; }
  .activity-row.act-found { color:var(--accent); }
  .activity-row.act-warn { color:#fb923c; }
  .activity-row.act-success { color:#4ade80; }
  .act-time { color:var(--text-muted); font-family:'Courier Prime',monospace; font-size:.65rem; flex-shrink:0; }
  .act-bullet { color:var(--text-muted); flex-shrink:0; }
  .act-msg { flex:1; }

  /* Providers */
  .providers-view,.mcp-view,.storage-view,.activity-view { display:flex; flex-direction:column; gap:0; }
  .comp-list { padding:.25rem 0; }
  .comp-row { display:flex; align-items:flex-start; gap:.75rem; padding:.6rem 1rem; border-bottom:1px solid var(--border); transition:background .1s; }
  .comp-row:last-child { border-bottom:none; }
  .comp-row:hover { background:var(--accent-dim); }
  .comp-info { flex:1; display:flex; flex-direction:column; gap:2px; }
  .comp-name { font-weight:600; font-size:.85rem; }
  .comp-desc { font-size:.65rem; color:var(--text-dim); }
  .comp-version { font-size:.65rem; color:var(--text-dim); font-family:'Courier Prime',monospace; }
  .comp-path { font-size:.6rem; color:var(--text-muted); font-family:'Courier Prime',monospace; }
  .comp-status { font-size:.75rem; font-weight:600; letter-spacing:.05em; text-align:right; flex-shrink:0; }
  .comp-port { font-family:'Courier Prime',monospace; font-size:.65rem; color:var(--text-dim); }
  .comp-right { display:flex; flex-direction:column; align-items:flex-end; gap:4px; }
  .transport-badge { font-size:.6rem; padding:1px 5px; border-radius:2px; background:var(--bg3); border:1px solid var(--border); color:var(--text-dim); letter-spacing:.1em; }

  /* Agents */
  .agents-view { }
  .agent-grid { display:grid; grid-template-columns:repeat(auto-fill,minmax(200px,1fr)); gap:.75rem; }
  .agent-card { background:var(--bg2); border:1px solid var(--border); border-radius:6px; padding:1rem; cursor:pointer; transition:all .2s; text-align:left; color:var(--text); font-family:'Rajdhani',monospace; }
  .agent-card:hover { border-color:var(--accent); transform:translateY(-2px); box-shadow:0 4px 20px var(--accent-dim); }
  .agent-card.detected { border-color:color-mix(in srgb,var(--accent) 30%,transparent); }
  .agent-card.partial { border-color:color-mix(in srgb,#f59e0b 30%,transparent); }
  .ac-top { display:flex; justify-content:space-between; align-items:center; margin-bottom:.5rem; }
  .ac-icon { font-size:1.5rem; }
  .ac-status-dot { width:8px; height:8px; border-radius:50%; }
  .ac-name { font-weight:700; font-size:.9rem; }
  .ac-eco { font-size:.65rem; color:var(--text-dim); letter-spacing:.1em; margin-bottom:.5rem; }
  .ac-progress { display:flex; align-items:center; gap:.5rem; margin-bottom:.4rem; }
  .ac-bar { flex:1; height:3px; background:var(--bg3); border-radius:2px; overflow:hidden; }
  .ac-bar-fill { height:100%; background:var(--accent); border-radius:2px; transition:width .5s; }
  .ac-found { font-size:.65rem; color:var(--text-dim); font-family:'Courier Prime',monospace; }
  .ac-status-label { font-size:.6rem; font-weight:700; letter-spacing:.15em; }

  .agent-detail { display:flex; flex-direction:column; gap:1rem; }
  .back-btn { background:none; border:1px solid var(--border); color:var(--text-dim); padding:.35rem .75rem; border-radius:3px; cursor:pointer; font-family:'Rajdhani',monospace; font-size:.75rem; align-self:flex-start; transition:all .15s; }
  .back-btn:hover { border-color:var(--accent); color:var(--accent); }
  .agent-card-header { display:flex; align-items:flex-start; gap:1rem; background:var(--bg2); border:1px solid var(--border); border-radius:6px; padding:1.25rem; }
  .agent-big-icon { font-size:2.5rem; }
  .agent-big-name { font-family:'Special Elite',serif; font-size:1.3rem; color:var(--text); }
  .agent-ecosystem { font-size:.65rem; color:var(--text-dim); letter-spacing:.2em; text-transform:uppercase; margin:2px 0; }
  .agent-desc { font-size:.75rem; color:var(--text-dim); }
  .agent-status-badge { font-weight:700; font-size:.7rem; letter-spacing:.15em; margin-left:auto; }
  .expected-grid { display:flex; flex-direction:column; gap:2px; padding:.5rem 0; }
  .expected-item { display:flex; align-items:flex-start; gap:.75rem; padding:.4rem 1rem; border-bottom:1px solid var(--border); opacity:.4; transition:opacity .15s; }
  .expected-item.found { opacity:1; }
  .expected-item:last-child { border-bottom:none; }
  .expected-check { color:var(--text-dim); font-size:.8rem; }
  .expected-item.found .expected-check { color:var(--accent); }
  .expected-info { display:flex; flex-direction:column; }
  .expected-label { font-size:.8rem; font-weight:600; }
  .expected-path { font-size:.6rem; color:var(--text-dim); font-family:'Courier Prime',monospace; }
  .agent-file-row { display:flex; align-items:center; gap:.75rem; padding:.4rem 1rem; border-bottom:1px solid var(--border); font-size:.75rem; }
  .agent-file-row:last-child { border-bottom:none; }
  .file-type-badge { font-size:.6rem; padding:1px 5px; border-radius:2px; background:var(--bg3); border:1px solid var(--border); color:var(--accent); letter-spacing:.1em; text-transform:uppercase; flex-shrink:0; }
  .file-name { flex:1; font-family:'Courier Prime',monospace; }
  .file-size { color:var(--text-dim); font-family:'Courier Prime',monospace; font-size:.65rem; }

  /* Models */
  .models-view { }
  .models-toolbar { display:flex; align-items:center; gap:1rem; margin-bottom:.75rem; flex-wrap:wrap; }
  .search-input { background:var(--bg2); border:1px solid var(--border); color:var(--text); padding:.35rem .75rem; border-radius:3px; font-family:'Rajdhani',monospace; font-size:.8rem; width:220px; outline:none; transition:border-color .15s; }
  .search-input:focus { border-color:var(--accent); }
  .filter-toggle { display:flex; align-items:center; gap:.4rem; font-size:.75rem; color:var(--text-dim); cursor:pointer; }
  .filter-toggle input { accent-color:var(--accent); }
  .model-total { margin-left:auto; font-size:.7rem; color:var(--text-dim); font-family:'Courier Prime',monospace; }
  .model-list { display:flex; flex-direction:column; background:var(--bg2); border:1px solid var(--border); border-radius:6px; overflow:hidden; }
  .model-row { display:flex; align-items:center; gap:1rem; padding:.6rem 1rem; border-bottom:1px solid var(--border); transition:background .1s; }
  .model-row:last-child { border-bottom:none; }
  .model-row:hover { background:var(--accent-dim); }
  .model-row.orphaned { border-left:3px solid #fb923c; }
  .model-left { flex:1; min-width:0; }
  .model-name { font-weight:600; font-size:.85rem; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }
  .model-path { font-size:.6rem; color:var(--text-muted); font-family:'Courier Prime',monospace; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }
  .model-right { display:flex; align-items:center; gap:.5rem; flex-shrink:0; }
  .orphan-badge { font-size:.6rem; padding:1px 5px; background:rgba(251,146,60,.15); border:1px solid #fb923c; color:#fb923c; border-radius:2px; letter-spacing:.1em; }
  .format-badge { font-size:.6rem; padding:1px 5px; background:var(--bg3); border:1px solid var(--border); color:var(--text-dim); border-radius:2px; }
  .provider-badge { font-size:.6rem; padding:1px 5px; background:var(--accent-dim); border:1px solid var(--accent); color:var(--accent); border-radius:2px; }
  .model-size { font-family:'Courier Prime',monospace; font-size:.75rem; min-width:60px; text-align:right; }
  .delete-btn { background:none; border:none; cursor:pointer; font-size:.85rem; opacity:.4; transition:opacity .15s; padding:2px 4px; }
  .delete-btn:hover { opacity:1; }
  .dir-row { padding:.5rem 1rem; border-bottom:1px solid var(--border); }
  .dir-row:last-child { border-bottom:none; }
  .dir-row code { font-size:.7rem; color:var(--text-dim); font-family:'Courier Prime',monospace; }

  /* Storage */
  .disk-card { }
  .disk-header { display:flex; justify-content:space-between; align-items:center; padding:1rem 1rem .5rem; }
  .disk-mount-big { font-family:'Special Elite',serif; font-size:1.1rem; }
  .disk-meta { font-size:.65rem; color:var(--text-dim); font-family:'Courier Prime',monospace; margin-top:2px; }
  .disk-pct-big { font-family:'Special Elite',serif; font-size:1.8rem; }
  .disk-bar-big-wrap { margin:0 1rem .75rem; height:8px; background:var(--bg3); border-radius:4px; overflow:hidden; }
  .disk-bar-big-fill { height:100%; transition:width .5s,background .5s; border-radius:4px; }
  .disk-stats { display:flex; padding:.75rem 1rem; border-top:1px solid var(--border); }
  .disk-stat { flex:1; text-align:center; }
  .ds-val { display:block; font-family:'Courier Prime',monospace; font-size:1rem; font-weight:700; color:var(--text); }
  .ds-label { font-size:.6rem; color:var(--text-dim); letter-spacing:.1em; text-transform:uppercase; }

  .activity-log-full { max-height:60vh; overflow-y:auto; }
  .clear-btn { background:none; border:1px solid var(--border); color:var(--text-dim); padding:2px 8px; border-radius:2px; cursor:pointer; font-size:.65rem; font-family:'Rajdhani',monospace; transition:all .15s; }
  .clear-btn:hover { border-color:var(--accent); color:var(--accent); }

  /* Modal */
  .modal-overlay { position:fixed; inset:0; background:rgba(0,0,0,.8); display:flex; align-items:center; justify-content:center; z-index:999; backdrop-filter:blur(4px); }
  .modal { background:var(--bg2); border:1px solid var(--border); border-radius:6px; padding:1.5rem; width:min(480px,90vw); box-shadow:0 0 40px rgba(0,0,0,.8); }
  .modal-title { font-family:'Special Elite',serif; color:#fb923c; font-size:1rem; margin-bottom:1rem; }
  .modal-body code { display:block; font-size:.7rem; color:var(--text-dim); background:var(--bg3); padding:.5rem; border-radius:3px; border:1px solid var(--border); margin-bottom:.75rem; word-break:break-all; font-family:'Courier Prime',monospace; }
  .modal-body p { font-size:.75rem; color:var(--text-dim); line-height:1.6; }
  .modal-actions { display:flex; justify-content:flex-end; gap:.75rem; margin-top:1.25rem; }
  .modal-cancel { background:none; border:1px solid var(--border); color:var(--text-dim); padding:.4rem 1rem; border-radius:3px; cursor:pointer; font-family:'Rajdhani',monospace; font-size:.8rem; transition:all .15s; }
  .modal-cancel:hover { border-color:var(--accent); color:var(--accent); }
  .modal-confirm { background:#fb923c; border:none; color:var(--bg); padding:.4rem 1rem; border-radius:3px; cursor:pointer; font-family:'Rajdhani',monospace; font-size:.8rem; font-weight:700; transition:all .15s; }
  .modal-confirm:hover:not(:disabled) { background:#f97316; }
  .modal-confirm:disabled { opacity:.6; cursor:not-allowed; }
</style>
