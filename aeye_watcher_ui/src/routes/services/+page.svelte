<script lang="ts">
  import { store, pollGPU, formatBytes } from '$lib/stores/app.svelte';
  import { onMount, onDestroy } from 'svelte';

  let interval: ReturnType<typeof setInterval>;
  let selected = $state<string | null>(null);

  onMount(() => { interval = setInterval(pollGPU, 2000); });
  onDestroy(() => clearInterval(interval));

  // All components that exist (not just running)
  let allServices = $derived(
    (store.scanResult?.components ?? []).filter(c =>
      ['Provider','Agent','Gateway','Vector DB','Observability'].includes(c.type)
    )
  );
  let running  = $derived(allServices.filter(c => c.status === 'Running'));
  let installed = $derived(allServices.filter(c => c.status === 'Installed'));
  let offline  = $derived(allServices.filter(c => c.status === 'Not Found'));

  let selectedComp = $derived(allServices.find(c => c.name === selected) ?? null);

  const typeColor: Record<string,string> = {
    Provider: '#c084fc', Agent: '#4ade80', Gateway: '#7dd3fc',
    'Vector DB': '#fb923c', Observability: '#fbbf24',
  };
  const typeIcon: Record<string,string> = {
    Provider: '⚙', Agent: '◈', Gateway: '⟷', 'Vector DB': '◆', Observability: '◷',
  };
</script>

<div class="services-page">
  <div class="page-hdr">
    <span class="hdr-icon">⚙</span>
    <div>
      <h1>Live Services</h1>
      <p class="hdr-sub">
        <span class="pill green">{running.length} running</span>
        <span class="pill blue">{installed.length} installed</span>
        <span class="pill dim">{offline.length} not found</span>
      </p>
    </div>
    {#if store.gpuInfo && !store.gpuInfo.error}
      <div class="gpu-card">
        <div class="gc-label">VRAM</div>
        <div class="gc-bar"><div class="gc-fill" style="width:{store.gpuInfo.use_pct}%;background:{store.gpuInfo.use_pct>85?'#ef4444':store.gpuInfo.use_pct>60?'#f59e0b':'var(--accent)'}"></div></div>
        <div class="gc-val">{store.gpuInfo.used_mb} / {store.gpuInfo.total_mb} MB</div>
        {#if store.gpuInfo.name}<div class="gc-name">{store.gpuInfo.name}</div>{/if}
      </div>
    {/if}
  </div>

  {#if !store.scanResult}
    <div class="empty"><span class="ei">⚙</span><p>Run a scan to discover services.</p><a href="/scan" class="scan-link">◉ Scan Center</a></div>
  {:else}

  <div class="services-layout">
    <!-- Service list -->
    <div class="svc-list">
      {#if running.length > 0}
        <div class="group-hdr running-hdr">● RUNNING <span class="cnt">{running.length}</span></div>
        {#each running as svc}
          <button class="svc-row running" class:selected={selected===svc.name} onclick={() => selected = selected===svc.name ? null : svc.name}>
            <div class="svc-dot running-dot"></div>
            <div class="svc-info">
              <span class="svc-name">{svc.name}</span>
              <span class="svc-type" style="color:{typeColor[svc.type] ?? 'var(--text-dim)'}">{typeIcon[svc.type] ?? '·'} {svc.type}</span>
            </div>
            <div class="svc-right">
              {#if svc.port}<span class="svc-port">:{svc.port}</span>{/if}
              <span class="svc-status-badge running-badge">LIVE</span>
            </div>
          </button>
        {/each}
      {/if}

      {#if installed.length > 0}
        <div class="group-hdr installed-hdr">◉ INSTALLED <span class="cnt">{installed.length}</span></div>
        {#each installed as svc}
          <button class="svc-row" class:selected={selected===svc.name} onclick={() => selected = selected===svc.name ? null : svc.name}>
            <div class="svc-dot installed-dot"></div>
            <div class="svc-info">
              <span class="svc-name">{svc.name}</span>
              <span class="svc-type" style="color:{typeColor[svc.type] ?? 'var(--text-dim)'}">{typeIcon[svc.type] ?? '·'} {svc.type}</span>
            </div>
            <div class="svc-right">
              <span class="svc-status-badge installed-badge">OFFLINE</span>
            </div>
          </button>
        {/each}
      {/if}

      {#if offline.length > 0}
        <div class="group-hdr offline-hdr">○ NOT DETECTED <span class="cnt">{offline.length}</span></div>
        {#each offline as svc}
          <button class="svc-row dim" class:selected={selected===svc.name} onclick={() => selected = selected===svc.name ? null : svc.name}>
            <div class="svc-dot offline-dot"></div>
            <div class="svc-info">
              <span class="svc-name">{svc.name}</span>
              <span class="svc-type">{typeIcon[svc.type] ?? '·'} {svc.type}</span>
            </div>
            <div class="svc-right">
              <span class="svc-status-badge offline-badge">—</span>
            </div>
          </button>
        {/each}
      {/if}
    </div>

    <!-- Detail panel -->
    <div class="detail-panel">
      {#if selectedComp}
        <div class="detail-header">
          <span class="detail-icon">{typeIcon[selectedComp.type] ?? '·'}</span>
          <div>
            <div class="detail-name">{selectedComp.name}</div>
            <div class="detail-type" style="color:{typeColor[selectedComp.type] ?? 'var(--text-dim)'}">{selectedComp.type}</div>
          </div>
          <div class="detail-status-badge" class:run={selectedComp.status==='Running'} class:inst={selectedComp.status==='Installed'}>
            {selectedComp.status}
          </div>
        </div>

        <div class="detail-meta">
          {#if selectedComp.description}
            <div class="dm-row"><span class="dm-label">Description</span><span class="dm-val">{selectedComp.description}</span></div>
          {/if}
          {#if selectedComp.port}
            <div class="dm-row"><span class="dm-label">Port</span><span class="dm-val mono">{selectedComp.port}</span></div>
          {/if}
          {#if selectedComp.url}
            <div class="dm-row"><span class="dm-label">Endpoint</span><a class="dm-link" href={selectedComp.url} target="_blank">{selectedComp.url}</a></div>
          {/if}
          {#if selectedComp.path}
            <div class="dm-row"><span class="dm-label">Binary</span><span class="dm-val mono">{selectedComp.path}</span></div>
          {/if}
          {#if selectedComp.version}
            <div class="dm-row"><span class="dm-label">Version</span><span class="dm-val mono">{selectedComp.version}</span></div>
          {/if}
        </div>

        {#if selectedComp.status === 'Running'}
          <!-- Live telemetry section -->
          <div class="telem-section">
            <div class="telem-hdr">◉ LIVE TELEMETRY</div>
            <div class="telem-grid">
              <div class="telem-card">
                <div class="tc-icon">●</div>
                <div class="tc-val">LIVE</div>
                <div class="tc-label">Status</div>
              </div>
              {#if store.gpuInfo && !store.gpuInfo.error}
                <div class="telem-card">
                  <div class="tc-icon">⚡</div>
                  <div class="tc-val">{store.gpuInfo.use_pct.toFixed(0)}<span class="tc-unit">%</span></div>
                  <div class="tc-label">VRAM Used</div>
                </div>
                <div class="telem-card">
                  <div class="tc-icon">◆</div>
                  <div class="tc-val">{store.gpuInfo.used_mb}<span class="tc-unit">MB</span></div>
                  <div class="tc-label">VRAM MB</div>
                </div>
                {#if store.gpuInfo.temp_c}
                  <div class="telem-card" class:hot={store.gpuInfo.temp_c > 80}>
                    <div class="tc-icon">🌡</div>
                    <div class="tc-val">{store.gpuInfo.temp_c}<span class="tc-unit">°C</span></div>
                    <div class="tc-label">GPU Temp</div>
                  </div>
                {/if}
              {/if}
            </div>

            {#if store.gpuInfo && !store.gpuInfo.error}
              <div class="vram-bar-wrap">
                <div class="vram-label">
                  <span>VRAM Usage</span>
                  <span>{store.gpuInfo.used_mb} / {store.gpuInfo.total_mb} MB</span>
                </div>
                <div class="vram-track">
                  <div class="vram-fill" style="width:{store.gpuInfo.use_pct}%;background:{store.gpuInfo.use_pct>85?'#ef4444':store.gpuInfo.use_pct>60?'#f59e0b':'var(--accent)'}">
                    <div class="vram-shimmer"></div>
                  </div>
                </div>
              </div>
            {/if}

            <div class="live-note">
              <span class="live-dot"></span>
              Polling every 2s · GPU stats are system-wide (not per-service)
            </div>
          </div>

          <!-- Quick actions -->
          <div class="actions-section">
            <div class="actions-hdr">◈ QUICK ACTIONS</div>
            <div class="action-btns">
              {#if selectedComp.url}
                <a href={selectedComp.url} target="_blank" class="action-btn">Open API ↗</a>
              {/if}
              <button class="action-btn secondary" onclick={() => selected = null}>Close</button>
            </div>
          </div>
        {:else if selectedComp.status === 'Installed'}
          <div class="offline-hint">
            <div class="oh-icon">◉</div>
            <div>
              <strong>{selectedComp.name} is installed but not running.</strong>
              <p>Start the service and re-scan to see live telemetry.</p>
              {#if selectedComp.path}
                <code class="oh-cmd">{selectedComp.path}</code>
              {/if}
            </div>
          </div>
        {:else}
          <div class="offline-hint">
            <div class="oh-icon" style="opacity:.3">○</div>
            <div>
              <strong>{selectedComp.name} was not detected on this system.</strong>
              <p>{selectedComp.description ?? 'Install and run, then re-scan.'}</p>
            </div>
          </div>
        {/if}
      {:else}
        <div class="detail-empty">
          <span class="de-icon">⚙</span>
          <p>Select a service to view details and live telemetry</p>
        </div>
      {/if}
    </div>
  </div>
  {/if}
</div>

<style>
  .services-page{display:flex;flex-direction:column;gap:1rem;padding:1.25rem;overflow-y:auto;min-height:100%;}
  .page-hdr{display:flex;align-items:center;gap:.85rem;flex-wrap:wrap;flex-shrink:0;}
  .hdr-icon{font-size:1.5rem;color:var(--accent);}
  h1{font-family:'Special Elite',serif;font-size:1.15rem;color:var(--text);}
  .hdr-sub{display:flex;gap:.4rem;margin-top:.3rem;flex-wrap:wrap;}
  .pill{font-size:.68rem;padding:2px 7px;border-radius:3px;font-weight:700;letter-spacing:.05em;}
  .pill.green{background:rgba(74,222,128,.12);border:1px solid #4ade80;color:#4ade80;}
  .pill.blue{background:rgba(96,165,250,.1);border:1px solid #60a5fa;color:#60a5fa;}
  .pill.dim{background:var(--bg3);border:1px solid var(--border);color:var(--text-dim);}
  .gpu-card{margin-left:auto;background:var(--bg2);border:1px solid var(--border);border-radius:6px;padding:.6rem .85rem;min-width:180px;}
  .gc-label{font-size:.6rem;letter-spacing:.15em;color:var(--text-muted);margin-bottom:.3rem;}
  .gc-bar{height:5px;background:var(--bg3);border-radius:3px;overflow:hidden;margin-bottom:.3rem;}
  .gc-fill{height:100%;border-radius:3px;transition:width .5s,background .5s;}
  .gc-val{font-family:'Courier Prime',monospace;font-size:.75rem;color:var(--text);}
  .gc-name{font-size:.65rem;color:var(--text-dim);margin-top:2px;}
  .empty{display:flex;flex-direction:column;align-items:center;text-align:center;padding:3rem;gap:.85rem;}
  .ei{font-size:2.5rem;color:var(--accent);opacity:.3;}
  .empty p{font-size:.9rem;color:var(--text-dim);}
  .scan-link{background:var(--accent);color:var(--bg);padding:.4rem .9rem;border-radius:4px;font-size:.84rem;font-weight:700;}
  .services-layout{display:grid;grid-template-columns:300px 1fr;gap:1rem;min-height:500px;}
  .svc-list{background:var(--bg2);border:1px solid var(--border);border-radius:8px;overflow:hidden;display:flex;flex-direction:column;}
  .group-hdr{padding:.45rem .85rem;font-size:.62rem;letter-spacing:.18em;font-weight:700;border-bottom:1px solid var(--border);}
  .running-hdr{background:color-mix(in srgb,var(--accent) 10%,var(--bg2));color:var(--accent);}
  .installed-hdr{background:color-mix(in srgb,#60a5fa 8%,var(--bg2));color:#60a5fa;}
  .offline-hdr{color:var(--text-muted);background:var(--bg2);}
  .cnt{float:right;background:var(--bg3);padding:0 5px;border-radius:2px;font-size:.6rem;}
  .svc-row{display:flex;align-items:center;gap:.65rem;padding:.6rem .85rem;border-bottom:1px solid var(--border);cursor:pointer;background:none;border-left:none;border-right:none;border-top:none;transition:background .1s;text-align:left;width:100%;}
  .svc-row:last-child{border-bottom:none;}
  .svc-row:hover{background:var(--accent-dim);}
  .svc-row.selected{background:var(--accent-dim);border-left:2px solid var(--accent);}
  .svc-row.dim{opacity:.55;}
  .svc-dot{width:8px;height:8px;border-radius:50%;flex-shrink:0;}
  .running-dot{background:var(--accent);box-shadow:0 0 6px var(--accent);animation:glow 2s ease-in-out infinite;}
  @keyframes glow{0%,100%{opacity:1}50%{opacity:.5}}
  .installed-dot{background:#60a5fa;}
  .offline-dot{background:var(--text-muted);}
  .svc-info{flex:1;display:flex;flex-direction:column;gap:1px;min-width:0;}
  .svc-name{font-size:.88rem;font-weight:600;color:var(--text);white-space:nowrap;overflow:hidden;text-overflow:ellipsis;}
  .svc-type{font-size:.68rem;}
  .svc-right{display:flex;flex-direction:column;align-items:flex-end;gap:2px;flex-shrink:0;}
  .svc-port{font-family:'Courier Prime',monospace;font-size:.65rem;color:var(--text-muted);}
  .svc-status-badge{font-size:.6rem;padding:1px 5px;border-radius:2px;font-weight:700;letter-spacing:.08em;}
  .running-badge{background:color-mix(in srgb,var(--accent) 15%,transparent);border:1px solid var(--accent);color:var(--accent);}
  .installed-badge{background:rgba(96,165,250,.1);border:1px solid #60a5fa;color:#60a5fa;}
  .offline-badge{color:var(--text-muted);border:1px solid var(--border);}
  .detail-panel{background:var(--bg2);border:1px solid var(--border);border-radius:8px;overflow:hidden;display:flex;flex-direction:column;}
  .detail-header{display:flex;align-items:center;gap:.85rem;padding:1.25rem;border-bottom:1px solid var(--border);flex-shrink:0;}
  .detail-icon{font-size:1.8rem;color:var(--accent);}
  .detail-name{font-family:'Special Elite',serif;font-size:1.2rem;color:var(--text);}
  .detail-type{font-size:.75rem;margin-top:2px;}
  .detail-status-badge{margin-left:auto;font-size:.72rem;font-weight:700;letter-spacing:.1em;padding:3px 9px;border-radius:3px;}
  .detail-status-badge.run{background:color-mix(in srgb,var(--accent) 15%,transparent);border:1px solid var(--accent);color:var(--accent);}
  .detail-status-badge.inst{background:rgba(96,165,250,.1);border:1px solid #60a5fa;color:#60a5fa;}
  .detail-meta{padding:.5rem 0;}
  .dm-row{display:flex;justify-content:space-between;align-items:baseline;padding:.45rem 1.25rem;border-bottom:1px solid var(--border);font-size:.84rem;gap:1rem;}
  .dm-label{color:var(--text-dim);flex-shrink:0;}
  .dm-val{color:var(--text);text-align:right;}
  .dm-val.mono{font-family:'Courier Prime',monospace;font-size:.78rem;}
  .dm-link{color:var(--accent);font-family:'Courier Prime',monospace;font-size:.78rem;}
  .telem-section{border-top:1px solid var(--border);padding:1rem 1.25rem;display:flex;flex-direction:column;gap:.85rem;}
  .telem-hdr{font-size:.68rem;letter-spacing:.18em;color:var(--accent);font-weight:700;}
  .telem-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(90px,1fr));gap:.5rem;}
  .telem-card{background:var(--bg3);border:1px solid var(--border);border-radius:6px;padding:.65rem;text-align:center;}
  .telem-card.hot{border-color:#ef4444;}
  .tc-icon{font-size:.85rem;color:var(--accent);margin-bottom:.25rem;}
  .tc-val{font-family:'Special Elite',serif;font-size:1.3rem;color:var(--text);}
  .tc-unit{font-size:.7rem;color:var(--text-dim);}
  .tc-label{font-size:.62rem;color:var(--text-muted);letter-spacing:.08em;margin-top:.2rem;}
  .vram-bar-wrap{display:flex;flex-direction:column;gap:.3rem;}
  .vram-label{display:flex;justify-content:space-between;font-size:.73rem;color:var(--text-dim);}
  .vram-track{height:8px;background:var(--bg3);border-radius:4px;overflow:hidden;position:relative;}
  .vram-fill{height:100%;border-radius:4px;transition:width .5s,background .5s;position:relative;overflow:hidden;}
  .vram-shimmer{position:absolute;inset:0;background:linear-gradient(90deg,transparent,rgba(255,255,255,.25),transparent);animation:shimmer 2s infinite;}
  @keyframes shimmer{0%{transform:translateX(-100%)}100%{transform:translateX(200%)}}
  .live-note{display:flex;align-items:center;gap:.4rem;font-size:.7rem;color:var(--text-muted);}
  .live-dot{width:6px;height:6px;border-radius:50%;background:var(--accent);animation:glow 2s ease-in-out infinite;flex-shrink:0;}
  .actions-section{border-top:1px solid var(--border);padding:1rem 1.25rem;}
  .actions-hdr{font-size:.68rem;letter-spacing:.18em;color:var(--text-dim);font-weight:700;margin-bottom:.65rem;}
  .action-btns{display:flex;gap:.5rem;flex-wrap:wrap;}
  .action-btn{background:var(--accent-dim);border:1px solid var(--accent);color:var(--accent);padding:.38rem .85rem;border-radius:4px;cursor:pointer;font-size:.82rem;font-family:'Rajdhani',monospace;font-weight:600;transition:all .15s;}
  .action-btn:hover{background:var(--accent);color:var(--bg);}
  .action-btn.secondary{background:none;border-color:var(--border);color:var(--text-dim);}
  .action-btn.secondary:hover{border-color:var(--accent);color:var(--accent);}
  .offline-hint{display:flex;gap:.85rem;padding:1.25rem;align-items:flex-start;}
  .oh-icon{font-size:1.5rem;color:var(--text-dim);flex-shrink:0;}
  .offline-hint strong{display:block;font-size:.9rem;color:var(--text);margin-bottom:.3rem;}
  .offline-hint p{font-size:.82rem;color:var(--text-dim);line-height:1.6;}
  .oh-cmd{display:block;margin-top:.5rem;font-family:'Courier Prime',monospace;font-size:.75rem;background:var(--bg3);padding:.3rem .5rem;border-radius:3px;border:1px solid var(--border);}
  .detail-empty{display:flex;flex-direction:column;align-items:center;justify-content:center;gap:.75rem;flex:1;text-align:center;padding:2rem;color:var(--text-muted);}
  .de-icon{font-size:2rem;opacity:.3;}
  .detail-empty p{font-size:.85rem;max-width:200px;line-height:1.6;}
  @media(max-width:768px){.services-layout{grid-template-columns:1fr;}.gpu-card{margin-left:0;width:100%;}}
</style>
