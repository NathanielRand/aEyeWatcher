<script lang="ts">
  import { store, formatGB, formatBytes } from '$lib/stores/app.svelte';

  let r = $derived(store.scanResult);
  let running = $derived((r?.components ?? []).filter(c => c.status === 'Running'));
  let orphaned = $derived((r?.models ?? []).filter(m => m.is_orphaned));
  let totalGB = $derived((r?.models ?? []).reduce((s,m) => s + m.size_bytes,0) / 1e9);
  let detected = $derived((r?.agent_cards ?? []).filter(a => a.status !== 'not_found'));
</script>

<div class="dashboard">
  {#if !r && !store.scanning}
    <!-- No scan yet -->
    <div class="no-scan-banner">
      <div class="nsb-eye">◉</div>
      <div class="nsb-text">
        <strong>No intelligence data yet.</strong>
        Run a scan to discover your AI ecosystem — models, agents, providers, MCP servers, and more.
      </div>
      <a href="/scan" class="nsb-btn">◉ Go to Scan Center</a>
    </div>
  {:else if store.scanning}
    <div class="scanning-banner">
      <span class="spin">⟳</span>
      <span>Scan running in background — <a href="/scan">view progress →</a></span>
    </div>
  {/if}

  <!-- Stat grid -->
  <div class="stat-grid">
    <a href="/scan" class="stat-card" class:dim={!r}>
      <div class="sc-icon">⚙</div>
      <div class="sc-val">{running.length || '—'}</div>
      <div class="sc-label">Live Services</div>
    </a>
    <a href="/models" class="stat-card" class:dim={!r}>
      <div class="sc-icon">◆</div>
      <div class="sc-val">{r?.models?.length ?? '—'}</div>
      <div class="sc-label">Models Found</div>
    </a>
    <a href="/models?filter=orphaned" class="stat-card" class:warn={orphaned.length > 0} class:dim={!r}>
      <div class="sc-icon">⚠</div>
      <div class="sc-val">{orphaned.length || '—'}</div>
      <div class="sc-label">Orphaned</div>
    </a>
    <a href="/mcp" class="stat-card" class:dim={!r}>
      <div class="sc-icon">⬡</div>
      <div class="sc-val">{r?.mcp_servers?.length ?? '—'}</div>
      <div class="sc-label">MCP Servers</div>
    </a>
    <a href="/agents" class="stat-card" class:dim={!r}>
      <div class="sc-icon">◈</div>
      <div class="sc-val">{detected.length || '—'}</div>
      <div class="sc-label">Agents</div>
    </a>
    <a href="/models" class="stat-card" class:dim={!r}>
      <div class="sc-icon">◉</div>
      <div class="sc-val">{r ? totalGB.toFixed(1) : '—'}<span class="sc-unit">{r ? 'GB' : ''}</span></div>
      <div class="sc-label">Model Storage</div>
    </a>
  </div>

  {#if running.length > 0}
    <div class="section-card">
      <div class="card-hdr"><span>⚙ LIVE SERVICES</span><span class="cnt">{running.length}</span></div>
      <div class="chip-wrap">
        {#each running as c}
          <div class="chip running">
            <span class="cdot"></span>{c.name}
            {#if c.port}<span class="cport">:{c.port}</span>{/if}
          </div>
        {/each}
      </div>
    </div>
  {/if}

  {#if r?.disks && r.disks.length > 0}
    <div class="section-card">
      <div class="card-hdr"><span>◉ STORAGE</span><span class="cnt">{r.disks.length}</span></div>
      {#each r.disks as disk}
        <div class="disk-row">
          <div class="dr-left">
            <span class="dr-mount">{disk.mount}</span>
            <span class="dr-dev">{disk.device}</span>
          </div>
          <div class="dr-bar"><div class="dr-fill" style="width:{disk.use_pct}%;background:{disk.use_pct>90?'#ef4444':disk.use_pct>75?'#f59e0b':'var(--accent)'}"></div></div>
          <div class="dr-right">
            <span class="dr-pct">{disk.use_pct.toFixed(0)}%</span>
            <span class="dr-free">{formatGB(disk.free_gb)} free</span>
          </div>
        </div>
      {/each}
    </div>
  {/if}

  {#if orphaned.length > 0}
    <div class="section-card warn-card">
      <div class="card-hdr"><span>⚠ ORPHANED MODELS</span><span class="cnt">{orphaned.length}</span></div>
      {#each orphaned.slice(0,5) as m}
        <div class="orphan-row">
          <span class="om-name">{m.name}</span>
          <span class="om-size">{formatBytes(m.size_bytes)}</span>
        </div>
      {/each}
      {#if orphaned.length > 5}
        <a href="/models?filter=orphaned" class="see-all">See all {orphaned.length} orphaned →</a>
      {/if}
    </div>
  {/if}
</div>

<style>
  .dashboard{display:flex;flex-direction:column;gap:1rem;padding:1.25rem;}
  .no-scan-banner{display:flex;align-items:center;gap:1rem;background:var(--bg2);border:1px dashed var(--border);border-radius:8px;padding:1.25rem;flex-wrap:wrap;}
  .nsb-eye{font-size:1.8rem;color:var(--accent);opacity:.5;flex-shrink:0;}
  .nsb-text{flex:1;font-size:.88rem;color:var(--text-dim);line-height:1.6;}
  .nsb-text strong{color:var(--text);}
  .nsb-btn{background:var(--accent);border:none;color:var(--bg);padding:.45rem 1rem;border-radius:4px;font-size:.82rem;font-weight:700;letter-spacing:.06em;white-space:nowrap;flex-shrink:0;font-family:'Rajdhani',monospace;}
  .nsb-btn:hover{opacity:.85;}
  .scanning-banner{display:flex;align-items:center;gap:.65rem;background:var(--accent-dim);border:1px solid var(--accent);border-radius:6px;padding:.65rem 1rem;font-size:.82rem;color:var(--accent);}
  .scanning-banner a{color:var(--accent);text-decoration:underline;}
  .spin{animation:spin .8s linear infinite;display:inline-block;}
  @keyframes spin{to{transform:rotate(360deg)}}

  .stat-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(150px,1fr));gap:.75rem;}
  .stat-card{background:var(--bg2);border:1px solid var(--border);border-radius:6px;padding:.85rem;text-align:center;position:relative;overflow:hidden;cursor:pointer;transition:border-color .15s,transform .1s;display:block;}
  .stat-card::before{content:'';position:absolute;top:0;left:0;right:0;height:2px;background:var(--accent);opacity:.5;}
  .stat-card.warn::before{background:#fb923c;}
  .stat-card.dim{opacity:.5;}
  .stat-card:hover:not(.dim){border-color:var(--accent);transform:translateY(-1px);}
  .sc-icon{font-size:.95rem;color:var(--accent);margin-bottom:.35rem;opacity:.7;}
  .sc-val{font-family:'Special Elite',serif;font-size:1.75rem;color:var(--text);line-height:1;}
  .sc-unit{font-size:1rem;color:var(--text-dim);}
  .sc-label{font-size:.68rem;color:var(--text-dim);letter-spacing:.1em;margin-top:.3rem;}

  .section-card{background:var(--bg2);border:1px solid var(--border);border-radius:6px;overflow:hidden;}
  .warn-card{border-color:rgba(251,146,60,.3);}
  .card-hdr{display:flex;justify-content:space-between;align-items:center;padding:.6rem 1rem;background:var(--accent-dim);border-bottom:1px solid var(--border);font-size:.72rem;letter-spacing:.14em;color:var(--accent);font-weight:700;}
  .cnt{background:var(--bg3);border:1px solid var(--border);padding:1px 5px;border-radius:2px;color:var(--text-dim);font-size:.65rem;}
  .chip-wrap{display:flex;flex-wrap:wrap;gap:.4rem;padding:.75rem 1rem;}
  .chip{display:flex;align-items:center;gap:.35rem;font-size:.78rem;background:var(--bg3);border:1px solid var(--border);padding:.25rem .55rem;border-radius:3px;}
  .chip.running{border-color:var(--accent);}
  .cdot{width:7px;height:7px;border-radius:50%;background:var(--accent);box-shadow:0 0 4px var(--accent);}
  .cport{color:var(--text-dim);font-family:'Courier Prime',monospace;font-size:.7rem;}
  .disk-row{display:flex;align-items:center;gap:.75rem;padding:.55rem 1rem;border-bottom:1px solid var(--border);}
  .disk-row:last-child{border-bottom:none;}
  .dr-left{min-width:100px;}
  .dr-mount{font-size:.82rem;font-weight:600;display:block;}
  .dr-dev{font-size:.65rem;color:var(--text-dim);font-family:'Courier Prime',monospace;}
  .dr-bar{flex:1;height:4px;background:var(--bg3);border-radius:2px;overflow:hidden;}
  .dr-fill{height:100%;border-radius:2px;}
  .dr-right{min-width:140px;text-align:right;}
  .dr-pct{font-size:.78rem;font-family:'Courier Prime',monospace;}
  .dr-free{font-size:.68rem;color:var(--text-dim);margin-left:.4rem;}
  .orphan-row{display:flex;justify-content:space-between;align-items:center;padding:.45rem 1rem;border-bottom:1px solid var(--border);font-size:.82rem;}
  .orphan-row:last-of-type{border-bottom:none;}
  .om-name{color:var(--text-dim);white-space:nowrap;overflow:hidden;text-overflow:ellipsis;max-width:70%;}
  .om-size{color:#fb923c;font-family:'Courier Prime',monospace;font-size:.72rem;flex-shrink:0;}
  .see-all{display:block;padding:.5rem 1rem;font-size:.75rem;color:var(--accent);border-top:1px solid var(--border);text-align:center;}
  @media(max-width:600px){.stat-grid{grid-template-columns:repeat(2,1fr);}}
</style>
