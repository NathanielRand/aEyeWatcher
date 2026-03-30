<script lang="ts">
  import { store, openPath, readFile } from '$lib/stores/app.svelte';
  import FileViewer from '$lib/components/FileViewer.svelte';

  let viewFile = $state<{path:string;content:string;name:string}|null>(null);

  async function tryOpen(path: string, name: string) {
    const readable = /\.(md|json|yaml|yml|txt|toml|soul)$/i.test(path);
    if (readable) {
      const c = await readFile(path);
      if (c) { viewFile = {path,content:c,name}; return; }
    }
    openPath(path);
  }

  let items = $derived(
    'acp' === 'gateways'
      ? (store.scanResult?.components ?? []).filter(c => c.type === 'Gateway' && c.status !== 'Not Found')
      : 'acp' === 'acp'
        ? (store.scanResult?.acp_agents ?? [])
        : (store.scanResult?.mcp_servers ?? [])
  );
  let running = $derived(
    'acp' === 'gateways'
      ? items.filter((i: any) => i.status === 'Running').length
      : 'acp' === 'acp'
        ? items.filter((i: any) => i.status === 'running').length
        : items.filter((i: any) => i.status === 'running').length
  );
</script>

<div class="page">
  <div class="page-hdr">
    <span class="hdr-icon">▸</span>
    <div>
      <h1>ACP Agents</h1>
      <p class="hdr-sub">{items.length} found · {running} active</p>
    </div>
  </div>

  {#if !store.scanResult}
    <div class="empty"><span class="ei">▸</span><p>Run a scan to detect ACP Agents.</p><a href="/scan" class="scan-link">◉ Scan Center</a></div>
  {:else if items.length === 0}
    <div class="empty-inner">No ACP agents detected on standard ports. ACP agents expose <code>/.well-known/agent.json</code>.</div>
  {:else}
    <div class="item-list">
      {#each items as item}
        <div class="item-row">
          <div class="dot {'status' in item && (item.status === 'Running' || item.status === 'running') ? 'dot-run' : 'status' in item && (item.status === 'Installed' || item.status === 'configured') ? 'dot-inst' : 'dot-off'}"></div>
          <div class="item-info">
            <span class="item-name">{'name' in item ? item.name : '?'}</span>
            {'description' in item && item.description ? `<span class="item-desc">${item.description}</span>` : ''}
            {#if 'url' in item && item.url}
              <button class="item-path clickable" onclick={() => openPath(item.url!)}>{'url' in item ? item.url : ''}</button>
            {/if}
            {#if 'config' in item && item.config}
              <button class="item-path clickable" onclick={() => tryOpen(item.config!, ('name' in item ? item.name : '') + ' config')}>Config: {'config' in item ? item.config : ''}</button>
            {/if}
            {#if 'endpoint' in item && item.endpoint}
              <button class="item-path clickable" onclick={() => openPath(item.endpoint)}>{'endpoint' in item ? item.endpoint : ''}</button>
            {/if}
            {#if 'path' in item && item.path}
              <button class="item-path clickable" onclick={() => openPath(item.path!)}>{'path' in item ? item.path : ''}</button>
            {/if}
            {#if 'version' in item && item.version}
              <span class="item-version">{'version' in item ? item.version : ''}</span>
            {/if}
          </div>
          <div class="item-right">
            {#if 'port' in item && item.port}
              <span class="badge port">:{'port' in item ? item.port : ''}</span>
            {/if}
            {#if 'transport' in item && item.transport}
              <span class="badge transport">{'transport' in item ? item.transport : ''}</span>
            {/if}
            {#if 'protocol' in item && item.protocol}
              <span class="badge transport">{'protocol' in item ? item.protocol : ''}</span>
            {/if}
            <span class="badge status-badge" class:running={'status' in item && (item.status === 'Running' || item.status === 'running')} class:installed={'status' in item && (item.status === 'Installed' || item.status === 'configured')}>
              {'status' in item ? item.status : 'unknown'}
            </span>
          </div>
        </div>
      {/each}
    </div>
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
  .scan-link{background:var(--accent);color:var(--bg);padding:.4rem .9rem;border-radius:4px;font-size:.84rem;font-weight:700;letter-spacing:.06em;}
  .empty-inner{background:var(--bg2);border:1px dashed var(--border);border-radius:8px;padding:2rem;font-size:.88rem;color:var(--text-dim);line-height:1.7;text-align:center;}
  .item-list{background:var(--bg2);border:1px solid var(--border);border-radius:8px;overflow:hidden;}
  .item-row{display:flex;align-items:flex-start;gap:.75rem;padding:.75rem 1rem;border-bottom:1px solid var(--border);transition:background .1s;}
  .item-row:last-child{border-bottom:none;}
  .item-row:hover{background:var(--accent-dim);}
  .dot{width:9px;height:9px;border-radius:50%;flex-shrink:0;margin-top:5px;}
  .dot.dot-run{background:var(--accent);box-shadow:0 0 6px var(--accent);animation:glow 2s ease-in-out infinite;}
  .dot.dot-inst{background:#60a5fa;}
  .dot.dot-off{background:var(--text-muted);}
  @keyframes glow{0%,100%{opacity:1}50%{opacity:.5}}
  .item-info{flex:1;display:flex;flex-direction:column;gap:3px;min-width:0;}
  .item-name{font-weight:600;font-size:.92rem;color:var(--text);}
  .item-desc{font-size:.76rem;color:var(--text-dim);}
  .item-path{font-size:.7rem;color:var(--text-muted);font-family:'Courier Prime',monospace;background:none;border:none;cursor:pointer;padding:0;text-align:left;transition:color .15s;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;max-width:100%;}
  .item-path:hover{color:var(--accent);}
  .item-version{font-size:.7rem;color:var(--text-dim);font-family:'Courier Prime',monospace;}
  .clickable{cursor:pointer;transition:color .15s;}
  .clickable:hover{color:var(--accent)!important;}
  .item-right{display:flex;flex-direction:column;align-items:flex-end;gap:4px;flex-shrink:0;}
  .badge{font-size:.62rem;padding:2px 6px;border-radius:2px;border:1px solid var(--border);color:var(--text-dim);letter-spacing:.06em;}
  .badge.port,.badge.transport{background:var(--bg3);font-family:'Courier Prime',monospace;}
  .badge.status-badge{background:var(--bg3);}
  .badge.status-badge.running{background:color-mix(in srgb,var(--accent) 12%,transparent);border-color:var(--accent);color:var(--accent);}
  .badge.status-badge.installed{background:rgba(96,165,250,.1);border-color:#60a5fa;color:#60a5fa;}
</style>
