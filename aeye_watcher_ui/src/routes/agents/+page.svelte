<script lang="ts">
  import { store, openPath, readFile, addToast } from '$lib/stores/app.svelte';
  import FileViewer from '$lib/components/FileViewer.svelte';

  let viewFile = $state<{path:string;content:string;name:string}|null>(null);
  let expandedCard = $state<string|null>(null);

  async function tryOpen(path: string, name: string) {
    const readable = /\.(md|json|yaml|yml|txt|toml|soul|env)$/i.test(path);
    if (readable) {
      const content = await readFile(path);
      if (content) { viewFile = {path, content, name}; return; }
    }
    openPath(path);
  }

  function progressPct(card: typeof store.scanResult.agent_cards[0]) {
    const found = card.expected_files?.filter(f=>f.found).length ?? 0;
    const total = card.expected_files?.length ?? 1;
    return { found, total, pct: (found/total)*100 };
  }

  let detected = $derived((store.scanResult?.agent_cards ?? []).filter(a => a.status !== 'not_found'));
  let notFound  = $derived((store.scanResult?.agent_cards ?? []).filter(a => a.status === 'not_found'));
</script>

<div class="page">
  <div class="page-hdr">
    <span class="hdr-icon">◈</span>
    <div>
      <h1>Agent Ecosystems</h1>
      <p class="hdr-sub">{detected.length} detected · {(store.scanResult?.agent_cards ?? []).length} profiled</p>
    </div>
  </div>

  {#if !store.scanResult}
    <div class="empty"><span class="ei">◈</span><p>Run a scan to profile agent ecosystems.</p><a href="/scan" class="scan-link">◉ Scan Center</a></div>
  {:else}
    <!-- Detected agents -->
    {#if detected.length > 0}
      <div class="section-label">DETECTED</div>
      <div class="card-grid">
        {#each detected as card}
          {@const prog = progressPct(card)}
          <div class="agent-card" class:detected={card.status==='detected'} class:partial={card.status==='partial'}>
            <div class="ac-head">
              <span class="ac-icon">{card.icon}</span>
              <span class="ac-status {card.status}">{card.status.replace('_',' ').toUpperCase()}</span>
            </div>
            <div class="ac-name">{card.name}</div>
            <div class="ac-eco">{card.ecosystem}</div>
            <div class="ac-desc">{card.description}</div>
            <div class="progress-row">
              <div class="prog-bar"><div class="prog-fill" style="width:{prog.pct}%"></div></div>
              <span class="prog-label">{prog.found}/{prog.total}</span>
            </div>
            <!-- Memory files -->
            {#if card.files?.some(f => ['soul','persona','memory','agent-def','knowledge'].includes(f.type))}
              <div class="mem-section">
                <div class="mem-hdr">◈ MEMORY & SOUL FILES</div>
                {#each card.files.filter(f => ['soul','persona','memory','agent-def','knowledge'].includes(f.type)) as mf}
                  <button class="mem-file" onclick={() => tryOpen(mf.path, mf.name)}>
                    <span class="mf-badge {mf.type}">{mf.type}</span>
                    <span class="mf-name">{mf.name}</span>
                    <span class="mf-open">↗</span>
                  </button>
                {/each}
              </div>
            {/if}
            <!-- Checklist toggle -->
            <button class="toggle-checklist" onclick={() => expandedCard = expandedCard === card.ecosystem ? null : card.ecosystem}>
              {expandedCard === card.ecosystem ? '▲' : '▼'} File Checklist ({prog.found}/{prog.total})
            </button>
            {#if expandedCard === card.ecosystem}
              <div class="checklist">
                {#each card.expected_files ?? [] as ef}
                  <div class="cl-item" class:found={ef.found}>
                    <span class="cl-check">{ef.found ? '✓' : '○'}</span>
                    {#if ef.found && ef.path}
                      <button class="cl-label clickable" onclick={() => tryOpen(ef.path!, ef.label)}>{ef.label}</button>
                    {:else}
                      <span class="cl-label">{ef.label}</span>
                    {/if}
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/each}
      </div>
    {/if}

    <!-- Not found agents (collapsed) -->
    {#if notFound.length > 0}
      <div class="section-label muted">NOT DETECTED ({notFound.length})</div>
      <div class="not-found-grid">
        {#each notFound as card}
          <div class="nf-card">
            <span class="nf-icon">{card.icon}</span>
            <span class="nf-name">{card.name}</span>
            <span class="nf-eco">{card.ecosystem}</span>
          </div>
        {/each}
      </div>
    {/if}
  {/if}
</div>

{#if viewFile}
  <FileViewer file={viewFile} onclose={() => viewFile=null} onopen={() => openPath(viewFile!.path)} />
{/if}

<style>
  .page{display:flex;flex-direction:column;gap:1rem;padding:1.25rem;overflow-y:auto;min-height:100%;}
  .page-hdr{display:flex;align-items:center;gap:.75rem;flex-shrink:0;}
  .hdr-icon{font-size:1.5rem;color:var(--accent);}
  h1{font-family:'Special Elite',serif;font-size:1.15rem;color:var(--text);}
  .hdr-sub{font-size:.82rem;color:var(--text-dim);margin-top:2px;}
  .empty{display:flex;flex-direction:column;align-items:center;text-align:center;padding:3rem;gap:.85rem;}
  .ei{font-size:2.5rem;color:var(--accent);opacity:.3;}
  .empty p{font-size:.9rem;color:var(--text-dim);}
  .scan-link{background:var(--accent);color:var(--bg);padding:.4rem .9rem;border-radius:4px;font-size:.84rem;font-weight:700;}
  .section-label{font-size:.65rem;letter-spacing:.2em;color:var(--accent);font-family:'Rajdhani',monospace;font-weight:700;padding:.25rem 0;}
  .section-label.muted{color:var(--text-muted);}
  .card-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(300px,1fr));gap:.85rem;}
  .agent-card{background:var(--bg2);border:1px solid var(--border);border-radius:8px;padding:1rem;}
  .agent-card.detected{border-color:color-mix(in srgb,var(--accent) 40%,transparent);}
  .agent-card.partial{border-color:color-mix(in srgb,#f59e0b 30%,transparent);}
  .ac-head{display:flex;justify-content:space-between;align-items:center;margin-bottom:.5rem;}
  .ac-icon{font-size:1.4rem;}
  .ac-status{font-size:.62rem;font-weight:700;letter-spacing:.1em;padding:2px 7px;border-radius:3px;}
  .ac-status.detected{color:var(--accent);background:var(--accent-dim);border:1px solid var(--accent);}
  .ac-status.partial{color:#f59e0b;background:rgba(245,158,11,.12);border:1px solid #f59e0b;}
  .ac-status.not_found{color:var(--text-muted);background:var(--bg3);}
  .ac-name{font-weight:700;font-size:1rem;color:var(--text);margin-bottom:2px;}
  .ac-eco{font-size:.68rem;color:var(--text-dim);letter-spacing:.1em;margin-bottom:.4rem;}
  .ac-desc{font-size:.82rem;color:var(--text-dim);line-height:1.55;margin-bottom:.65rem;}
  .progress-row{display:flex;align-items:center;gap:.5rem;margin-bottom:.65rem;}
  .prog-bar{flex:1;height:4px;background:var(--bg3);border-radius:2px;overflow:hidden;}
  .prog-fill{height:100%;background:var(--accent);transition:width .5s;}
  .prog-label{font-size:.68rem;color:var(--text-dim);font-family:'Courier Prime',monospace;flex-shrink:0;}
  .mem-section{background:var(--bg3);border:1px solid var(--border);border-radius:5px;overflow:hidden;margin-bottom:.6rem;}
  .mem-hdr{font-size:.62rem;letter-spacing:.12em;color:var(--accent);padding:.35rem .65rem;border-bottom:1px solid var(--border);font-weight:700;}
  .mem-file{display:flex;align-items:center;gap:.45rem;padding:.35rem .65rem;border-bottom:1px solid var(--border);width:100%;background:none;border-left:none;border-right:none;border-top:none;cursor:pointer;transition:background .1s;}
  .mem-file:last-child{border-bottom:none;}
  .mem-file:hover{background:var(--accent-dim);}
  .mf-badge{font-size:.58rem;padding:1px 5px;border-radius:2px;letter-spacing:.08em;text-transform:uppercase;flex-shrink:0;border:1px solid;}
  .mf-badge.soul{color:#c084fc;border-color:rgba(192,132,252,.4);background:rgba(192,132,252,.1);}
  .mf-badge.persona,.mf-badge.memory{color:var(--accent);border-color:var(--accent);background:var(--accent-dim);}
  .mf-badge.agent-def,.mf-badge.knowledge{color:#7dd3fc;border-color:rgba(125,211,252,.3);background:rgba(125,211,252,.08);}
  .mf-name{flex:1;font-size:.78rem;color:var(--text-dim);font-family:'Courier Prime',monospace;text-align:left;}
  .mf-open{font-size:.72rem;color:var(--text-muted);}
  .toggle-checklist{width:100%;background:none;border:1px solid var(--border);color:var(--text-dim);padding:.3rem .6rem;border-radius:4px;cursor:pointer;font-size:.75rem;font-family:'Rajdhani',monospace;text-align:left;transition:all .15s;}
  .toggle-checklist:hover{border-color:var(--accent);color:var(--accent);}
  .checklist{margin-top:.4rem;display:flex;flex-direction:column;gap:2px;}
  .cl-item{display:flex;align-items:center;gap:.5rem;font-size:.78rem;color:var(--text-muted);opacity:.5;padding:2px 0;}
  .cl-item.found{opacity:1;color:var(--text-dim);}
  .cl-check{font-size:.72rem;width:14px;flex-shrink:0;color:var(--text-muted);}
  .cl-item.found .cl-check{color:var(--accent);}
  .cl-label{background:none;border:none;cursor:default;padding:0;font-size:.78rem;font-family:'Rajdhani',monospace;text-align:left;color:var(--text-dim);}
  .cl-label.clickable{cursor:pointer;transition:color .15s;}
  .cl-label.clickable:hover{color:var(--accent);}
  .not-found-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(160px,1fr));gap:.5rem;}
  .nf-card{background:var(--bg2);border:1px solid var(--border);border-radius:5px;padding:.65rem .85rem;display:flex;align-items:center;gap:.5rem;opacity:.45;}
  .nf-icon{font-size:1rem;}
  .nf-name{font-size:.82rem;color:var(--text-dim);font-weight:600;}
  .nf-eco{font-size:.65rem;color:var(--text-muted);margin-left:auto;}
  @media(max-width:640px){.card-grid{grid-template-columns:1fr;}}
</style>
