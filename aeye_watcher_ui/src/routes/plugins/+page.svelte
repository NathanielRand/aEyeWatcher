<script lang="ts">
  import { store, installPlugin, uninstallPlugin, togglePlugin, loadPlugins, isProFeature, type Plugin } from '$lib/stores/app.svelte';
  import { onMount } from 'svelte';
  onMount(() => loadPlugins());

  // Built-in registry of available plugins
  const registry: Plugin[] = [
    { id:'aeye-ollama-deep', name:'Ollama Deep Inspector', version:'1.0.0', enabled:false, description:'Extended Ollama metrics: per-model load times, request counts, memory breakdowns.', author:'aEye Team', repoUrl:'https://github.com/NathanielRand/aEyeWatcher' },
    { id:'aeye-lmstudio-bridge', name:'LM Studio Bridge', version:'1.0.0', enabled:false, description:'Direct integration with LM Studio API for real-time model load/unload status.', author:'aEye Team', repoUrl:'https://github.com/NathanielRand/aEyeWatcher' },
    { id:'aeye-hf-hub', name:'HuggingFace Hub Scanner', version:'1.0.0', enabled:false, description:'Scans HF Hub cache with model card metadata, license info, and update checks.', author:'aEye Team' },
    { id:'aeye-docker-ai', name:'Docker AI Container Scanner', version:'0.9.0', enabled:false, description:'Detects AI workloads running in Docker containers (vLLM, TGI, LocalAI images).', author:'Community' },
    { id:'aeye-grafana', name:'Grafana Exporter', version:'0.8.0', enabled:false, description:'Exports aEye scan results as Prometheus metrics for Grafana dashboards.', author:'Community' },
    { id:'aeye-slack', name:'Slack Notifier', version:'1.0.0', enabled:false, description:'Send scan summaries and change alerts to a Slack channel via webhook.', author:'aEye Team' },
    { id:'aeye-discord', name:'Discord Notifier', version:'1.0.0', enabled:false, description:'Post scan events and orphan alerts to a Discord channel.', author:'Community' },
    { id:'aeye-custom-scanner', name:'Custom Scanner SDK', version:'0.5.0', enabled:false, description:'Template and SDK for writing your own component detectors in Go.', author:'aEye Team', repoUrl:'https://github.com/NathanielRand/aEyeWatcher' },
  ];

  let tab = $state<'installed'|'browse'>('browse');
  let search = $state('');
  let filtered = $derived(registry.filter(p =>
    !search || p.name.toLowerCase().includes(search.toLowerCase()) || p.description.toLowerCase().includes(search.toLowerCase())
  ));
  let installedPlugins = $derived(store.plugins);
</script>

<div class="page">
  <div class="page-hdr">
    <span class="hdr-icon">⬡</span>
    <div>
      <h1>Plugins</h1>
      <p class="hdr-sub">Extend aEye Watcher with custom scanners and integrations</p>
    </div>
  </div>

  <div class="tabs">
    <button class="tab" class:active={tab==='browse'} onclick={() => tab='browse'}>Browse Registry ({registry.length})</button>
    <button class="tab" class:active={tab==='installed'} onclick={() => tab='installed'}>Installed ({installedPlugins.length})</button>
  </div>

  {#if tab === 'browse'}
    <div class="search-row">
      <input class="search" type="text" placeholder="Search plugins…" bind:value={search} />
    </div>
    <div class="plugin-grid">
      {#each filtered as plugin}
        {@const isInstalled = installedPlugins.some(p => p.id === plugin.id)}
        <div class="plugin-card" class:installed={isInstalled}>
          <div class="pc-header">
            <div class="pc-name">{plugin.name}</div>
            <span class="pc-version">v{plugin.version}</span>
          </div>
          <div class="pc-author">by {plugin.author}</div>
          <div class="pc-desc">{plugin.description}</div>
          <div class="pc-footer">
            {#if plugin.repoUrl}
              <a class="pc-repo" href={plugin.repoUrl} target="_blank" rel="noopener">GitHub ↗</a>
            {/if}
            {#if isInstalled}
              <button class="pc-btn uninstall" onclick={() => uninstallPlugin(plugin.id)}>Remove</button>
            {:else}
              <button class="pc-btn install" onclick={() => installPlugin(plugin)}>Install</button>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  {:else}
    {#if installedPlugins.length === 0}
      <div class="empty"><span class="ei">⬡</span><p>No plugins installed. Browse the registry to add extensions.</p></div>
    {:else}
      <div class="plugin-list">
        {#each installedPlugins as p}
          <div class="pl-row">
            <div class="pl-info">
              <span class="pl-name">{p.name}</span>
              <span class="pl-desc">{p.description}</span>
            </div>
            <div class="pl-actions">
              <span class="pl-version">v{p.version}</span>
              <label class="toggle">
                <input type="checkbox" checked={p.enabled} onchange={() => togglePlugin(p.id)} />
                <span class="toggle-track"></span>
              </label>
              <button class="pl-remove" onclick={() => uninstallPlugin(p.id)}>✕</button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  {/if}

  <div class="sdk-banner">
    <span class="sdk-icon">◈</span>
    <div>
      <strong>Build a plugin</strong>
      <p>The aEye Plugin SDK lets you write custom scanners in Go. Any CLI tool, port, or file pattern can be detected and surfaced in the dashboard.</p>
    </div>
    <a href="https://github.com/NathanielRand/aEyeWatcher" target="_blank" class="sdk-btn">View SDK →</a>
  </div>
</div>

<style>
  .page{display:flex;flex-direction:column;gap:1rem;padding:1.25rem;overflow-y:auto;min-height:100%;}
  .page-hdr{display:flex;align-items:center;gap:.75rem;}
  .hdr-icon{font-size:1.5rem;color:var(--accent);}
  h1{font-family:'Special Elite',serif;font-size:1.15rem;color:var(--text);}
  .hdr-sub{font-size:.82rem;color:var(--text-dim);margin-top:2px;}
  .tabs{display:flex;gap:0;border-bottom:1px solid var(--border);}
  .tab{background:none;border:none;border-bottom:2px solid transparent;color:var(--text-dim);padding:.5rem 1rem;cursor:pointer;font-family:'Rajdhani',monospace;font-size:.88rem;font-weight:500;transition:all .15s;}
  .tab:hover{color:var(--text);}
  .tab.active{color:var(--accent);border-bottom-color:var(--accent);}
  .search-row{display:flex;}
  .search{background:var(--bg2);border:1px solid var(--border);color:var(--text);padding:.38rem .75rem;border-radius:4px;font-family:'Courier Prime',monospace;font-size:.84rem;width:260px;outline:none;}
  .search:focus{border-color:var(--accent);}
  .plugin-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:.85rem;}
  .plugin-card{background:var(--bg2);border:1px solid var(--border);border-radius:8px;padding:1rem;display:flex;flex-direction:column;gap:.4rem;transition:border-color .15s;}
  .plugin-card:hover{border-color:var(--accent);}
  .plugin-card.installed{border-color:color-mix(in srgb,var(--accent) 35%,transparent);}
  .pc-header{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem;}
  .pc-name{font-weight:700;font-size:.95rem;color:var(--text);}
  .pc-version{font-size:.65rem;color:var(--text-muted);font-family:'Courier Prime',monospace;flex-shrink:0;}
  .pc-author{font-size:.72rem;color:var(--text-dim);}
  .pc-desc{font-size:.82rem;color:var(--text-dim);line-height:1.6;flex:1;}
  .pc-footer{display:flex;justify-content:space-between;align-items:center;margin-top:.25rem;}
  .pc-repo{font-size:.72rem;color:var(--text-dim);transition:color .15s;}
  .pc-repo:hover{color:var(--accent);}
  .pc-btn{border:none;padding:.3rem .75rem;border-radius:3px;cursor:pointer;font-family:'Rajdhani',monospace;font-size:.8rem;font-weight:700;transition:all .15s;}
  .pc-btn.install{background:var(--accent);color:var(--bg);}
  .pc-btn.install:hover{opacity:.85;}
  .pc-btn.uninstall{background:rgba(239,68,68,.15);border:1px solid #ef4444;color:#ef4444;}
  .pc-btn.uninstall:hover{background:#ef4444;color:#fff;}
  .plugin-list{background:var(--bg2);border:1px solid var(--border);border-radius:8px;overflow:hidden;}
  .pl-row{display:flex;align-items:center;gap:1rem;padding:.7rem 1rem;border-bottom:1px solid var(--border);}
  .pl-row:last-child{border-bottom:none;}
  .pl-info{flex:1;display:flex;flex-direction:column;gap:2px;}
  .pl-name{font-weight:600;font-size:.9rem;color:var(--text);}
  .pl-desc{font-size:.76rem;color:var(--text-dim);}
  .pl-actions{display:flex;align-items:center;gap:.6rem;}
  .pl-version{font-size:.68rem;color:var(--text-muted);font-family:'Courier Prime',monospace;}
  .toggle{cursor:pointer;display:flex;align-items:center;}
  .toggle input{display:none;}
  .toggle-track{width:32px;height:16px;border-radius:8px;background:var(--bg3);border:1px solid var(--border);position:relative;transition:background .2s;}
  .toggle-track::after{content:'';position:absolute;top:2px;left:2px;width:10px;height:10px;border-radius:50%;background:var(--text-muted);transition:transform .2s,background .2s;}
  .toggle input:checked+.toggle-track{background:var(--accent-dim);border-color:var(--accent);}
  .toggle input:checked+.toggle-track::after{transform:translateX(16px);background:var(--accent);}
  .pl-remove{background:none;border:none;color:var(--text-muted);cursor:pointer;font-size:.8rem;transition:color .15s;}
  .pl-remove:hover{color:#ef4444;}
  .empty{display:flex;flex-direction:column;align-items:center;text-align:center;padding:3rem;gap:.85rem;}
  .ei{font-size:2.5rem;color:var(--accent);opacity:.3;}
  .empty p{font-size:.9rem;color:var(--text-dim);}
  .sdk-banner{display:flex;align-items:center;gap:1rem;background:var(--bg2);border:1px dashed var(--border);border-radius:8px;padding:1.1rem;flex-wrap:wrap;}
  .sdk-icon{font-size:1.4rem;color:var(--accent);flex-shrink:0;}
  .sdk-banner div{flex:1;}
  .sdk-banner strong{display:block;font-size:.9rem;color:var(--text);margin-bottom:.2rem;}
  .sdk-banner p{font-size:.8rem;color:var(--text-dim);line-height:1.6;}
  .sdk-btn{background:var(--accent-dim);border:1px solid var(--accent);color:var(--accent);padding:.38rem .85rem;border-radius:4px;font-size:.82rem;font-weight:700;font-family:'Rajdhani',monospace;white-space:nowrap;flex-shrink:0;}
  .sdk-btn:hover{background:var(--accent);color:var(--bg);}
  @media(max-width:640px){.plugin-grid{grid-template-columns:1fr;}}
</style>
