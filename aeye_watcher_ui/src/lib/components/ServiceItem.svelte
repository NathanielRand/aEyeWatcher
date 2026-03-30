<script lang="ts">
  import type { Component } from '$lib/types';
  let { comp, onclick }: { comp: Component; onclick?: () => void } = $props();
  
  const statusColor = (s: string) => s==='Running' ? 'var(--accent)' : s==='Installed' ? '#60a5fa' : 'var(--text-muted)';
  const dotClass = (s: string) => s==='Running' ? 'dot-run' : s==='Installed' ? 'dot-inst' : 'dot-off';
</script>

<div class="item" class:clickable={!!onclick} onclick={onclick}>
  <div class="dot {dotClass(comp.status)}"></div>
  <div class="info">
    <span class="name">{comp.name}</span>
    {#if comp.description}<span class="desc">{comp.description}</span>{/if}
    {#if comp.version}<span class="version">{comp.version}</span>{/if}
    {#if comp.url || comp.path}<span class="path">{comp.url ?? comp.path}</span>{/if}
  </div>
  <div class="right">
    {#if comp.port}<span class="badge port">:{comp.port}</span>{/if}
    <span class="badge status" style="color:{statusColor(comp.status)};border-color:{statusColor(comp.status)}">{comp.status}</span>
  </div>
</div>

<style>
  .item{display:flex;align-items:flex-start;gap:.75rem;padding:.7rem 1rem;border-bottom:1px solid var(--border);transition:background .1s;}
  .item:last-child{border-bottom:none;}
  .item.clickable{cursor:pointer;}
  .item:hover{background:var(--accent-dim);}
  .dot{width:9px;height:9px;border-radius:50%;flex-shrink:0;margin-top:4px;}
  .dot.dot-run{background:var(--accent);box-shadow:0 0 6px var(--accent);animation:glow 2s ease-in-out infinite;}
  .dot.dot-inst{background:#60a5fa;}
  .dot.dot-off{background:var(--text-muted);}
  @keyframes glow{0%,100%{opacity:1}50%{opacity:.55}}
  .info{flex:1;display:flex;flex-direction:column;gap:2px;min-width:0;}
  .name{font-weight:600;font-size:.92rem;color:var(--text);}
  .desc{font-size:.75rem;color:var(--text-dim);}
  .version{font-size:.7rem;color:var(--text-dim);font-family:'Courier Prime',monospace;}
  .path{font-size:.68rem;color:var(--text-muted);font-family:'Courier Prime',monospace;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;}
  .right{display:flex;flex-direction:column;align-items:flex-end;gap:3px;flex-shrink:0;}
  .badge{font-size:.62rem;padding:1px 6px;border-radius:2px;border:1px solid var(--border);color:var(--text-dim);letter-spacing:.06em;}
  .badge.port{font-family:'Courier Prime',monospace;}
  .badge.status{background:transparent;}
</style>
