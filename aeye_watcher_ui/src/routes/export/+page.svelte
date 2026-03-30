<script lang="ts">
  import ProGate from '$lib/components/ProGate.svelte';
  import { store, formatBytes, formatDuration } from '$lib/stores/app.svelte';

  let exporting = $state(false);
  let exportDone = $state(false);

  function exportJSON() {
    if (!store.scanResult) return;
    exporting = true;
    const blob = new Blob([JSON.stringify({
      meta: store.scanMeta,
      result: store.scanResult,
      exportedAt: new Date().toISOString(),
      version: '2.0.0',
    }, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `aeye-scan-${new Date().toISOString().split('T')[0]}.json`;
    a.click();
    URL.revokeObjectURL(url);
    exporting = false;
    exportDone = true;
    setTimeout(() => exportDone = false, 3000);
  }

  function exportCSV() {
    if (!store.scanResult) return;
    const rows = [
      ['Name','Provider','Format','Size (bytes)','Orphaned','Path'],
      ...(store.scanResult.models ?? []).map(m => [
        m.name, m.provider, m.format ?? m.extension ?? '', String(m.size_bytes),
        String(m.is_orphaned), m.file_path
      ])
    ];
    const csv = rows.map(r => r.map(c => `"${c.replace(/"/g,'""')}"`).join(',')).join('\n');
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url; a.download = `aeye-models-${new Date().toISOString().split('T')[0]}.csv`;
    a.click(); URL.revokeObjectURL(url);
  }

  let r = $derived(store.scanResult);
  let totalBytes = $derived((r?.models ?? []).reduce((s,m) => s + m.size_bytes, 0));
</script>

<ProGate feature="Scan Export">
<div class="page">
  <div class="page-hdr">
    <span class="hdr-icon">◆</span>
    <div>
      <h1>Export</h1>
      <p class="hdr-sub">Download scan results as JSON, CSV, or generate reports</p>
    </div>
  </div>

  {#if !r}
    <div class="empty"><span class="ei">◆</span><p>Run a scan first to enable exports.</p><a href="/scan" class="scan-link">◉ Scan Center</a></div>
  {:else}
    <div class="export-grid">
      <div class="export-card">
        <div class="ec-icon">{ }</div>
        <div class="ec-name">Full Scan Report (JSON)</div>
        <div class="ec-desc">Complete scan result including all components, models, agents, MCP servers, and scan metadata. Ideal for programmatic processing.</div>
        <div class="ec-meta">
          ~{formatBytes(JSON.stringify(r).length)} · {(r.models?.length ?? 0)} models · {(r.components?.length ?? 0)} components
        </div>
        <button class="ec-btn" onclick={exportJSON} disabled={exporting}>
          {exportDone ? '✓ Downloaded!' : exporting ? '…' : '↓ Export JSON'}
        </button>
      </div>

      <div class="export-card">
        <div class="ec-icon">⋮⋮</div>
        <div class="ec-name">Model Cache (CSV)</div>
        <div class="ec-desc">Spreadsheet of all discovered model files: name, provider, format, size, orphan status, and full path.</div>
        <div class="ec-meta">{r.models?.length ?? 0} models · {formatBytes(totalBytes)} total</div>
        <button class="ec-btn" onclick={exportCSV}>↓ Export CSV</button>
      </div>

      <div class="export-card coming">
        <div class="ec-icon">◈</div>
        <div class="ec-name">PDF Report <span class="soon-badge">SOON</span></div>
        <div class="ec-desc">Formatted PDF report with executive summary, charts, and recommendations. Great for sharing with your team.</div>
        <div class="ec-meta">Coming in v2.1</div>
        <button class="ec-btn" disabled>↓ Export PDF</button>
      </div>

      <div class="export-card coming">
        <div class="ec-icon">◷</div>
        <div class="ec-name">Prometheus Metrics <span class="soon-badge">SOON</span></div>
        <div class="ec-desc">Export metrics in Prometheus format for ingestion into Grafana or other observability platforms.</div>
        <div class="ec-meta">Coming in v2.1</div>
        <button class="ec-btn" disabled>↓ Export Metrics</button>
      </div>
    </div>

    <!-- Scan summary preview -->
    <div class="preview-card">
      <div class="pv-hdr">◉ CURRENT SCAN SNAPSHOT</div>
      <div class="pv-grid">
        <div class="pv-stat"><span class="pvs-val">{r.models?.length ?? 0}</span><span class="pvs-label">Models</span></div>
        <div class="pv-stat"><span class="pvs-val">{(r.models?.filter(m=>m.is_orphaned).length ?? 0)}</span><span class="pvs-label">Orphaned</span></div>
        <div class="pv-stat"><span class="pvs-val">{r.components?.filter(c=>c.status==='Running').length ?? 0}</span><span class="pvs-label">Running</span></div>
        <div class="pv-stat"><span class="pvs-val">{r.mcp_servers?.length ?? 0}</span><span class="pvs-label">MCP</span></div>
        <div class="pv-stat"><span class="pvs-val">{store.scanMeta ? formatDuration(store.scanMeta.durationMs ?? 0) : '—'}</span><span class="pvs-label">Duration</span></div>
        <div class="pv-stat"><span class="pvs-val">{r.scanned_at ? new Date(r.scanned_at).toLocaleDateString() : '—'}</span><span class="pvs-label">Scan Date</span></div>
      </div>
    </div>
  {/if}
</div>
</ProGate>

<style>
  .page{display:flex;flex-direction:column;gap:1rem;padding:1.25rem;overflow-y:auto;min-height:100%;}
  .page-hdr{display:flex;align-items:center;gap:.75rem;}
  .hdr-icon{font-size:1.5rem;color:var(--accent);}
  h1{font-family:'Special Elite',serif;font-size:1.15rem;color:var(--text);}
  .hdr-sub{font-size:.82rem;color:var(--text-dim);margin-top:2px;}
  .empty{display:flex;flex-direction:column;align-items:center;text-align:center;padding:3rem;gap:.85rem;}
  .ei{font-size:2.5rem;color:var(--accent);opacity:.3;}
  .empty p{font-size:.9rem;color:var(--text-dim);}
  .scan-link{background:var(--accent);color:var(--bg);padding:.4rem .9rem;border-radius:4px;font-size:.84rem;font-weight:700;}
  .export-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(260px,1fr));gap:.85rem;}
  .export-card{background:var(--bg2);border:1px solid var(--border);border-radius:8px;padding:1.1rem;display:flex;flex-direction:column;gap:.5rem;transition:border-color .15s;}
  .export-card:hover:not(.coming){border-color:var(--accent);}
  .export-card.coming{opacity:.6;}
  .ec-icon{font-size:1.2rem;color:var(--accent);}
  .ec-name{font-weight:700;font-size:.95rem;color:var(--text);display:flex;align-items:center;gap:.4rem;}
  .ec-desc{font-size:.8rem;color:var(--text-dim);line-height:1.6;flex:1;}
  .ec-meta{font-size:.7rem;color:var(--text-muted);font-family:'Courier Prime',monospace;}
  .ec-btn{background:var(--accent);border:none;color:var(--bg);padding:.38rem .85rem;border-radius:4px;cursor:pointer;font-family:'Rajdhani',monospace;font-size:.84rem;font-weight:700;transition:opacity .15s;margin-top:.25rem;}
  .ec-btn:hover:not(:disabled){opacity:.85;}
  .ec-btn:disabled{opacity:.4;cursor:not-allowed;}
  .soon-badge{font-size:.58rem;padding:1px 5px;background:rgba(251,146,60,.15);border:1px solid #fb923c;color:#fb923c;border-radius:2px;letter-spacing:.08em;}
  .preview-card{background:var(--bg2);border:1px solid var(--border);border-radius:8px;overflow:hidden;}
  .pv-hdr{padding:.6rem 1rem;background:var(--accent-dim);border-bottom:1px solid var(--border);font-size:.7rem;letter-spacing:.15em;color:var(--accent);font-weight:700;}
  .pv-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(100px,1fr));gap:0;}
  .pv-stat{padding:.85rem 1rem;border-right:1px solid var(--border);text-align:center;}
  .pv-stat:last-child{border-right:none;}
  .pvs-val{display:block;font-family:'Special Elite',serif;font-size:1.5rem;color:var(--text);}
  .pvs-label{font-size:.65rem;color:var(--text-dim);letter-spacing:.1em;}
</style>
