<script lang="ts">
  import { page } from '$app/stores';
  import { store, formatBytes, requestDelete, openPath, readFile, addToast } from '$lib/stores/app.svelte';

  let searchQuery = $state('');
  let filterOrphaned = $state($page.url.searchParams.get('filter') === 'orphaned');
  let filterProvider = $state('all');
  let deleteConfirm = $state<string|null>(null);
  let viewFile = $state<{path:string;content:string;name:string}|null>(null);
  let loadingFile = $state(false);

  let providers = $derived([...new Set((store.scanResult?.models ?? []).map(m => m.provider))].sort());

  let filtered = $derived((() => {
    let models = store.scanResult?.models ?? [];
    if (filterOrphaned) models = models.filter(m => m.is_orphaned);
    if (filterProvider !== 'all') models = models.filter(m => m.provider === filterProvider);
    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      models = models.filter(m => m.name.toLowerCase().includes(q) || m.provider.toLowerCase().includes(q) || m.file_path.toLowerCase().includes(q));
    }
    return models;
  })());

  let totalBytes = $derived(filtered.reduce((s,m) => s + m.size_bytes, 0));

  async function tryReadFile(path: string, name: string) {
    const readable = /\.(md|json|yaml|yml|txt|toml|env|cfg|ini|log)$/i.test(path);
    if (!readable) { openPath(path); return; }
    loadingFile = true;
    const content = await readFile(path);
    loadingFile = false;
    if (content) viewFile = { path, content, name };
    else { addToast('Could not read file — try opening directly.', 'warn'); openPath(path); }
  }
</script>

<div class="models-page">
  <div class="page-header">
    <div class="ph-left">
      <span class="ph-icon">◆</span>
      <div>
        <h1>Model Cache</h1>
        <p class="ph-sub">
          {filtered.length} model{filtered.length !== 1 ? 's' : ''} · {formatBytes(totalBytes)}
          {#if filterOrphaned} · <span class="warn-text">orphaned filter active</span>{/if}
        </p>
      </div>
    </div>
    <div class="ph-right">
      <input class="search" type="text" placeholder="Search models…" bind:value={searchQuery}/>
      <label class="toggle-label">
        <input type="checkbox" bind:checked={filterOrphaned}/>
        <span>Orphaned only</span>
      </label>
    </div>
  </div>

  <!-- Provider filter tabs -->
  {#if providers.length > 0}
    <div class="provider-tabs">
      <button class="ptab" class:active={filterProvider==='all'} onclick={() => filterProvider='all'}>
        All <span class="ptab-cnt">{store.scanResult?.models?.length ?? 0}</span>
      </button>
      {#each providers as p}
        {@const cnt = (store.scanResult?.models ?? []).filter(m => m.provider === p).length}
        <button class="ptab" class:active={filterProvider===p} onclick={() => filterProvider=p}>
          {p} <span class="ptab-cnt">{cnt}</span>
        </button>
      {/each}
    </div>
  {/if}

  <!-- Model list -->
  <div class="model-list">
    {#each filtered as model}
      <div class="model-row" class:orphaned={model.is_orphaned}>
        <div class="mr-left">
          <div class="mr-name">{model.name}</div>
          <!-- Clickable file path -->
          <button class="mr-path clickable" onclick={() => openPath(model.file_path)} title="Open in explorer">
            {model.file_path}
          </button>
        </div>
        <div class="mr-right">
          {#if model.is_orphaned}<span class="badge orphan">ORPHANED</span>{/if}
          <span class="badge format">{model.format ?? model.extension}</span>
          <span class="badge provider">{model.provider}</span>
          <span class="mr-size">{formatBytes(model.size_bytes)}</span>
          <button class="icon-btn" onclick={() => deleteConfirm = model.file_path} title="Request cleanup">🗑</button>
        </div>
      </div>
    {/each}
    {#if filtered.length === 0}
      <div class="empty">{store.scanResult ? 'No models match your filter.' : 'Run a scan to discover model files.'}</div>
    {/if}
  </div>

  <!-- Known dirs -->
  {#if (store.scanResult?.model_dirs?.length ?? 0) > 0}
    <div class="dirs-card">
      <div class="dc-hdr">◉ KNOWN CACHE DIRECTORIES</div>
      {#each store.scanResult!.model_dirs as dir}
        <button class="dir-row clickable" onclick={() => openPath(dir)}>{dir}</button>
      {/each}
    </div>
  {/if}
</div>

<!-- Delete modal -->
{#if deleteConfirm}
  <div class="modal-bg" role="presentation" onclick={() => deleteConfirm=null}>
    <div class="modal" onclick={(e)=>e.stopPropagation()}>
      <div class="modal-title">⚠ REQUEST CLEANUP</div>
      <code class="modal-path">{deleteConfirm}</code>
      <p class="modal-note">Safe mode is active by default. This logs the request — no files are deleted unless the core has <code>AEYE_UNSAFE_DELETE=1</code>.</p>
      <div class="modal-btns">
        <button class="btn-cancel" onclick={() => deleteConfirm=null}>Cancel</button>
        <button class="btn-confirm" onclick={async () => { await requestDelete(deleteConfirm!); addToast('Cleanup requested','warn'); deleteConfirm=null; }}>Request</button>
      </div>
    </div>
  </div>
{/if}

<!-- File viewer -->
{#if viewFile}
  <div class="modal-bg" role="presentation" onclick={() => viewFile=null}>
    <div class="file-viewer" onclick={(e)=>e.stopPropagation()}>
      <div class="fv-header">
        <span class="fv-name">{viewFile.name}</span>
        <button class="fv-path-btn" onclick={() => openPath(viewFile!.path)}>Open in explorer ↗</button>
        <button class="fv-close" onclick={() => viewFile=null}>✕</button>
      </div>
      <pre class="fv-content">{viewFile.content}</pre>
    </div>
  </div>
{/if}

<style>
  .models-page{display:flex;flex-direction:column;gap:.85rem;padding:1.25rem;overflow-y:auto;min-height:100%;}
  .page-header{display:flex;align-items:flex-start;justify-content:space-between;gap:1rem;flex-wrap:wrap;}
  .ph-left{display:flex;align-items:center;gap:.75rem;}
  .ph-icon{font-size:1.4rem;color:var(--accent);}
  .page-header h1{font-family:'Special Elite',serif;font-size:1.1rem;}
  .ph-sub{font-size:.78rem;color:var(--text-dim);margin-top:2px;}
  .warn-text{color:#fb923c;}
  .ph-right{display:flex;align-items:center;gap:.75rem;flex-wrap:wrap;}
  .search{background:var(--bg2);border:1px solid var(--border);color:var(--text);padding:.35rem .7rem;border-radius:4px;font-family:'Courier Prime',monospace;font-size:.82rem;width:200px;outline:none;}
  .search:focus{border-color:var(--accent);}
  .toggle-label{display:flex;align-items:center;gap:.4rem;font-size:.78rem;color:var(--text-dim);cursor:pointer;}
  .toggle-label input{accent-color:var(--accent);}

  .provider-tabs{display:flex;gap:.3rem;flex-wrap:wrap;}
  .ptab{background:var(--bg2);border:1px solid var(--border);color:var(--text-dim);padding:.3rem .7rem;border-radius:4px;cursor:pointer;font-size:.78rem;font-family:'Rajdhani',monospace;transition:all .15s;display:flex;align-items:center;gap:.35rem;}
  .ptab:hover{border-color:var(--accent);color:var(--text);}
  .ptab.active{background:var(--accent-dim);border-color:var(--accent);color:var(--accent);}
  .ptab-cnt{font-size:.65rem;background:var(--bg3);padding:1px 4px;border-radius:2px;}

  .model-list{background:var(--bg2);border:1px solid var(--border);border-radius:6px;overflow:hidden;}
  .model-row{display:flex;align-items:center;gap:1rem;padding:.6rem 1rem;border-bottom:1px solid var(--border);transition:background .1s;}
  .model-row:last-child{border-bottom:none;}
  .model-row:hover{background:var(--accent-dim);}
  .model-row.orphaned{border-left:3px solid #fb923c;}
  .mr-left{flex:1;min-width:0;}
  .mr-name{font-weight:600;font-size:.88rem;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;}
  .mr-path{font-size:.68rem;color:var(--text-dim);font-family:'Courier Prime',monospace;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;display:block;max-width:100%;background:none;border:none;cursor:pointer;text-align:left;padding:0;}
  .clickable{cursor:pointer;transition:color .15s;}
  .clickable:hover{color:var(--accent)!important;}
  .mr-right{display:flex;align-items:center;gap:.4rem;flex-shrink:0;}
  .badge{font-size:.6rem;padding:1px 5px;border-radius:2px;letter-spacing:.07em;}
  .badge.orphan{background:rgba(251,146,60,.15);border:1px solid #fb923c;color:#fb923c;}
  .badge.format{background:var(--bg3);border:1px solid var(--border);color:var(--text-dim);}
  .badge.provider{background:var(--accent-dim);border:1px solid var(--accent);color:var(--accent);}
  .mr-size{font-family:'Courier Prime',monospace;font-size:.75rem;min-width:55px;text-align:right;}
  .icon-btn{background:none;border:none;cursor:pointer;font-size:.85rem;opacity:.4;transition:opacity .15s;padding:2px 3px;}
  .icon-btn:hover{opacity:1;}
  .empty{padding:2rem;text-align:center;color:var(--text-dim);font-size:.85rem;}

  .dirs-card{background:var(--bg2);border:1px solid var(--border);border-radius:6px;overflow:hidden;}
  .dc-hdr{padding:.6rem 1rem;background:var(--accent-dim);border-bottom:1px solid var(--border);font-size:.7rem;letter-spacing:.15em;color:var(--accent);font-weight:700;}
  .dir-row{display:block;width:100%;padding:.5rem 1rem;border-bottom:1px solid var(--border);font-size:.75rem;font-family:'Courier Prime',monospace;color:var(--text-dim);background:none;border-left:none;border-right:none;border-top:none;cursor:pointer;text-align:left;transition:color .15s;}
  .dir-row:last-child{border-bottom:none;}
  .dir-row:hover{color:var(--accent);}

  /* Modal */
  .modal-bg{position:fixed;inset:0;background:rgba(0,0,0,.75);display:flex;align-items:center;justify-content:center;z-index:1000;backdrop-filter:blur(4px);}
  .modal{background:var(--bg2);border:1px solid var(--border);border-radius:8px;padding:1.5rem;width:min(480px,90vw);}
  .modal-title{font-family:'Special Elite',serif;color:#fb923c;font-size:1rem;margin-bottom:.85rem;}
  .modal-path{display:block;font-size:.72rem;background:var(--bg3);padding:.5rem .65rem;border-radius:4px;border:1px solid var(--border);margin-bottom:.75rem;word-break:break-all;color:var(--text-dim);}
  .modal-note{font-size:.78rem;color:var(--text-dim);line-height:1.6;}
  .modal-btns{display:flex;justify-content:flex-end;gap:.65rem;margin-top:1.25rem;}
  .btn-cancel{background:none;border:1px solid var(--border);color:var(--text-dim);padding:.38rem .85rem;border-radius:3px;cursor:pointer;font-family:'Rajdhani',monospace;font-size:.82rem;}
  .btn-cancel:hover{border-color:var(--accent);color:var(--accent);}
  .btn-confirm{background:#fb923c;border:none;color:#000;padding:.38rem .85rem;border-radius:3px;cursor:pointer;font-family:'Rajdhani',monospace;font-size:.82rem;font-weight:700;}

  .file-viewer{background:var(--bg2);border:1px solid var(--border);border-radius:8px;width:min(760px,95vw);max-height:80vh;display:flex;flex-direction:column;overflow:hidden;}
  .fv-header{display:flex;align-items:center;gap:.65rem;padding:.65rem 1rem;border-bottom:1px solid var(--border);flex-shrink:0;}
  .fv-name{flex:1;font-size:.85rem;font-weight:600;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;}
  .fv-path-btn{background:none;border:1px solid var(--border);color:var(--text-dim);padding:.25rem .6rem;border-radius:3px;cursor:pointer;font-size:.72rem;font-family:'Rajdhani',monospace;transition:all .15s;white-space:nowrap;}
  .fv-path-btn:hover{color:var(--accent);border-color:var(--accent);}
  .fv-close{background:none;border:none;color:var(--text-dim);cursor:pointer;font-size:.85rem;padding:2px 5px;transition:color .15s;}
  .fv-close:hover{color:var(--text);}
  .fv-content{flex:1;overflow:auto;padding:1rem;font-family:'Courier Prime',monospace;font-size:.78rem;color:var(--text-dim);line-height:1.6;white-space:pre-wrap;word-break:break-word;}
</style>
