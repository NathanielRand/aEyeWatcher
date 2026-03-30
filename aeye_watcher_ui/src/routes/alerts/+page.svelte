<script lang="ts">
  import ProGate from '$lib/components/ProGate.svelte';
  import { store, addToast } from '$lib/stores/app.svelte';

  type AlertRule = { id:string; name:string; type:'new_model'|'orphan'|'service_down'|'service_up'|'new_mcp'; enabled:boolean; channel:'toast'|'webhook'; webhookUrl?:string; lastTriggered?:string };

  let rules = $state<AlertRule[]>([
    { id:'r1', name:'New orphaned model detected', type:'orphan', enabled:true, channel:'toast' },
    { id:'r2', name:'Service goes offline', type:'service_down', enabled:true, channel:'toast' },
    { id:'r3', name:'New model installed', type:'new_model', enabled:false, channel:'toast' },
    { id:'r4', name:'New MCP server configured', type:'new_mcp', enabled:false, channel:'toast' },
  ]);

  let showForm = $state(false);
  let form = $state({ name:'', type:'orphan' as AlertRule['type'], channel:'toast' as AlertRule['channel'], webhookUrl:'' });

  const typeLabels: Record<string,string> = {
    new_model:'New model installed', orphan:'Orphaned model detected',
    service_down:'Service went offline', service_up:'Service came online', new_mcp:'New MCP server found'
  };

  function saveRule() {
    if (!form.name.trim()) return;
    rules = [...rules, { id: Math.random().toString(36).slice(2), ...form, enabled:true }];
    showForm = false;
    addToast('Alert rule created', 'success');
  }
  function toggle(id: string) { rules = rules.map(r => r.id===id ? {...r, enabled:!r.enabled} : r); }
  function remove(id: string) { rules = rules.filter(r => r.id!==id); }
</script>

<ProGate feature="Change Detection Alerts">
<div class="page">
  <div class="page-hdr">
    <span class="hdr-icon">◷</span>
    <div>
      <h1>Change Detection Alerts</h1>
      <p class="hdr-sub">Get notified when your AI ecosystem changes between scans</p>
    </div>
    <button class="btn-primary" onclick={() => showForm=!showForm}>+ New Rule</button>
  </div>

  <div class="info-banner">
    <span>◈</span>
    Alerts fire by comparing the current scan result against the previous one. Run at least two scans to enable change detection.
    {#if store.scanMeta?.completedAt}
      Last scan: <strong>{new Date(store.scanMeta.completedAt).toLocaleString()}</strong>
    {:else}
      <strong>No scan data yet.</strong>
    {/if}
  </div>

  {#if showForm}
    <div class="form-card">
      <div class="form-title">New Alert Rule</div>
      <div class="form-fields">
        <div class="ff"><label>Name</label><input class="fi" type="text" bind:value={form.name} placeholder="e.g. Orphan alert"/></div>
        <div class="ff">
          <label>Trigger</label>
          <select class="fi" bind:value={form.type}>
            {#each Object.entries(typeLabels) as [k,v]}<option value={k}>{v}</option>{/each}
          </select>
        </div>
        <div class="ff">
          <label>Channel</label>
          <select class="fi" bind:value={form.channel}>
            <option value="toast">In-app toast notification</option>
            <option value="webhook">Webhook (Slack, Discord, custom)</option>
          </select>
        </div>
        {#if form.channel === 'webhook'}
          <div class="ff"><label>Webhook URL</label><input class="fi" type="text" bind:value={form.webhookUrl} placeholder="https://hooks.slack.com/…"/></div>
        {/if}
      </div>
      <div class="form-btns">
        <button class="btn-cancel" onclick={() => showForm=false}>Cancel</button>
        <button class="btn-primary" onclick={saveRule}>Save Rule</button>
      </div>
    </div>
  {/if}

  <div class="rules-list">
    {#each rules as rule}
      <div class="rule-row" class:disabled={!rule.enabled}>
        <div class="rr-info">
          <span class="rr-name">{rule.name}</span>
          <span class="rr-meta">{typeLabels[rule.type]} · {rule.channel}</span>
          {#if rule.lastTriggered}<span class="rr-last">Last triggered: {new Date(rule.lastTriggered).toLocaleString()}</span>{/if}
        </div>
        <div class="rr-actions">
          <label class="toggle"><input type="checkbox" checked={rule.enabled} onchange={() => toggle(rule.id)}/><span class="tt"></span></label>
          <button class="rr-del" onclick={() => remove(rule.id)}>✕</button>
        </div>
      </div>
    {/each}
    {#if rules.length === 0}<div class="empty-rules">No alert rules configured.</div>{/if}
  </div>
</div>
</ProGate>

<style>
  .page{display:flex;flex-direction:column;gap:1rem;padding:1.25rem;overflow-y:auto;min-height:100%;}
  .page-hdr{display:flex;align-items:center;gap:.75rem;flex-wrap:wrap;}
  .hdr-icon{font-size:1.5rem;color:var(--accent);}
  h1{font-family:'Special Elite',serif;font-size:1.15rem;color:var(--text);}
  .hdr-sub{font-size:.82rem;color:var(--text-dim);margin-top:2px;}
  .btn-primary{margin-left:auto;background:var(--accent);border:none;color:var(--bg);padding:.38rem .85rem;border-radius:4px;cursor:pointer;font-family:'Rajdhani',monospace;font-size:.84rem;font-weight:700;white-space:nowrap;}
  .btn-primary:hover{opacity:.85;}
  .btn-cancel{background:none;border:1px solid var(--border);color:var(--text-dim);padding:.38rem .85rem;border-radius:3px;cursor:pointer;font-family:'Rajdhani',monospace;font-size:.84rem;}
  .btn-cancel:hover{border-color:var(--accent);color:var(--accent);}
  .info-banner{background:var(--bg2);border:1px solid var(--border);border-radius:6px;padding:.75rem 1rem;font-size:.82rem;color:var(--text-dim);display:flex;gap:.5rem;align-items:baseline;flex-wrap:wrap;}
  .info-banner strong{color:var(--text);}
  .form-card{background:var(--bg2);border:1px solid var(--accent);border-radius:8px;padding:1.25rem;display:flex;flex-direction:column;gap:.85rem;}
  .form-title{font-family:'Special Elite',serif;font-size:.95rem;color:var(--accent);}
  .form-fields{display:flex;flex-direction:column;gap:.65rem;}
  .ff{display:flex;flex-direction:column;gap:.3rem;}
  .ff label{font-size:.78rem;font-weight:600;color:var(--text);}
  .fi{background:var(--bg3);border:1px solid var(--border);color:var(--text);padding:.38rem .65rem;border-radius:4px;font-family:'Rajdhani',monospace;font-size:.84rem;outline:none;}
  .fi:focus{border-color:var(--accent);}
  .form-btns{display:flex;justify-content:flex-end;gap:.5rem;}
  .rules-list{background:var(--bg2);border:1px solid var(--border);border-radius:8px;overflow:hidden;}
  .rule-row{display:flex;align-items:center;gap:1rem;padding:.75rem 1rem;border-bottom:1px solid var(--border);transition:opacity .15s;}
  .rule-row:last-child{border-bottom:none;}
  .rule-row.disabled{opacity:.5;}
  .rr-info{flex:1;display:flex;flex-direction:column;gap:2px;}
  .rr-name{font-weight:600;font-size:.9rem;color:var(--text);}
  .rr-meta{font-size:.74rem;color:var(--text-dim);}
  .rr-last{font-size:.68rem;color:var(--text-muted);font-family:'Courier Prime',monospace;}
  .rr-actions{display:flex;align-items:center;gap:.5rem;}
  .toggle{cursor:pointer;}
  .toggle input{display:none;}
  .tt{display:block;width:32px;height:16px;border-radius:8px;background:var(--bg3);border:1px solid var(--border);position:relative;transition:background .2s;}
  .tt::after{content:'';position:absolute;top:2px;left:2px;width:10px;height:10px;border-radius:50%;background:var(--text-muted);transition:transform .2s,background .2s;}
  .toggle input:checked+.tt{background:var(--accent-dim);border-color:var(--accent);}
  .toggle input:checked+.tt::after{transform:translateX(16px);background:var(--accent);}
  .rr-del{background:none;border:none;color:var(--text-muted);cursor:pointer;font-size:.8rem;transition:color .15s;}
  .rr-del:hover{color:#ef4444;}
  .empty-rules{padding:1.5rem;text-align:center;font-size:.85rem;color:var(--text-muted);}
</style>
