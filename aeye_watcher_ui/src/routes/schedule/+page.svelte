<script lang="ts">
  import { store, saveSchedule, deleteSchedule, computeNextRun, startScan, addToast, formatDuration, type ScheduleRule } from '$lib/stores/app.svelte';

  let showForm = $state(false);
  let editId = $state<string|null>(null);
  let form = $state({ label: '', cron: '60 min', enabled: true });

  const presets = [
    { label: 'Every 30 minutes', cron: '30 min' },
    { label: 'Every hour',       cron: '60 min' },
    { label: 'Every 6 hours',    cron: '6 hour' },
    { label: 'Every 12 hours',   cron: '12 hour' },
    { label: 'Daily',            cron: '24 hour' },
  ];

  function openNew() {
    editId = null;
    form = { label: '', cron: '60 min', enabled: true };
    showForm = true;
  }

  function openEdit(rule: ScheduleRule) {
    editId = rule.id;
    form = { label: rule.label, cron: rule.cron, enabled: rule.enabled };
    showForm = true;
  }

  function saveForm() {
    if (!form.label.trim()) { addToast('Please enter a label', 'warn'); return; }
    const id = editId ?? Math.random().toString(36).slice(2);
    const nextRun = computeNextRun(form.cron);
    saveSchedule({ id, label: form.label, cron: form.cron, enabled: form.enabled, nextRun });
    addToast(`Schedule "${form.label}" saved`, 'success');
    showForm = false;
  }

  function toggleEnabled(rule: ScheduleRule) {
    saveSchedule({ ...rule, enabled: !rule.enabled });
  }

  function triggerNow(rule: ScheduleRule) {
    addToast(`Manual trigger: "${rule.label}"`, 'info', '/scan');
    startScan();
  }

  function fmtTime(iso?: string) {
    if (!iso) return '—';
    return new Date(iso).toLocaleString();
  }
</script>

<div class="sched-page">
  <div class="page-header">
    <div class="ph-left">
      <span class="ph-icon">◷</span>
      <div>
        <h1>Scheduler</h1>
        <p class="ph-sub">Automated background scans with change detection</p>
      </div>
    </div>
    <button class="btn-primary" onclick={openNew}>+ New Schedule</button>
  </div>

  <!-- Info banner -->
  <div class="info-banner">
    <span class="ib-icon">◈</span>
    <span>Schedules run client-side while the UI is open. For persistent scheduling, use a system cron or task scheduler to call <code>GET /api/scan/trigger</code> on the core.</span>
  </div>

  {#if store.schedules.length === 0 && !showForm}
    <div class="empty-state">
      <div class="es-icon">◷</div>
      <h2>No schedules yet</h2>
      <p>Set up automated scans to track changes in your AI ecosystem over time.</p>
      <button class="btn-primary" onclick={openNew}>+ Create First Schedule</button>
    </div>
  {:else}
    <div class="schedule-list">
      {#each store.schedules as rule}
        <div class="rule-card" class:disabled={!rule.enabled}>
          <div class="rule-top">
            <div class="rule-info">
              <span class="rule-label">{rule.label}</span>
              <span class="rule-cron">{rule.cron}</span>
            </div>
            <div class="rule-actions">
              <button class="action-btn" onclick={() => triggerNow(rule)} title="Run now">▶</button>
              <button class="action-btn" onclick={() => openEdit(rule)} title="Edit">✎</button>
              <button class="action-btn danger" onclick={() => deleteSchedule(rule.id)} title="Delete">✕</button>
              <label class="toggle-switch">
                <input type="checkbox" checked={rule.enabled} onchange={() => toggleEnabled(rule)} />
                <span class="toggle-track"></span>
              </label>
            </div>
          </div>
          <div class="rule-meta">
            <span>Last run: <strong>{fmtTime(rule.lastRun)}</strong></span>
            <span>Next run: <strong class:accent={rule.enabled}>{rule.enabled ? fmtTime(rule.nextRun) : 'Disabled'}</strong></span>
          </div>
        </div>
      {/each}
    </div>
  {/if}

  <!-- Scan history (from scan meta) -->
  {#if store.scanMeta}
    <div class="history-card">
      <div class="hc-hdr">◉ LAST SCAN</div>
      <div class="hc-body">
        <div class="hc-row"><span>Status</span><span class:green={store.scanMeta.status==='complete'} class:red={store.scanMeta.status==='error'}>{store.scanMeta.status.toUpperCase()}</span></div>
        <div class="hc-row"><span>Started</span><span>{fmtTime(store.scanMeta.startedAt)}</span></div>
        <div class="hc-row"><span>Completed</span><span>{fmtTime(store.scanMeta.completedAt)}</span></div>
        <div class="hc-row"><span>Duration</span><span>{store.scanMeta.durationMs ? formatDuration(store.scanMeta.durationMs) : '—'}</span></div>
        <div class="hc-row"><span>Models found</span><span>{store.scanMeta.totalModels} ({store.scanMeta.orphanedModels} orphaned)</span></div>
        <div class="hc-row"><span>Live services</span><span>{store.scanMeta.runningServices}</span></div>
        <div class="hc-row"><span>Events</span><span>{store.scanMeta.eventsCount}</span></div>
      </div>
    </div>
  {/if}
</div>

<!-- Form modal -->
{#if showForm}
  <div class="modal-bg" role="presentation" onclick={() => showForm=false}>
    <div class="modal" onclick={(e)=>e.stopPropagation()}>
      <div class="modal-title">{editId ? 'Edit' : 'New'} Schedule</div>
      <div class="form-field">
        <label>Label</label>
        <input class="field-input" type="text" bind:value={form.label} placeholder="e.g. Hourly check"/>
      </div>
      <div class="form-field">
        <label>Interval</label>
        <div class="preset-row">
          {#each presets as p}
            <button class="preset-btn" class:active={form.cron===p.cron} onclick={() => form.cron=p.cron}>{p.label}</button>
          {/each}
        </div>
        <input class="field-input" type="text" bind:value={form.cron} placeholder="e.g. 30 min, 6 hour, 1 day"/>
        <span class="field-hint">Formats: <code>N min</code> · <code>N hour</code> · <code>N day</code></span>
      </div>
      <div class="form-field">
        <label class="toggle-label-inline">
          <input type="checkbox" bind:checked={form.enabled}/>
          <span>Enabled immediately</span>
        </label>
      </div>
      <div class="modal-btns">
        <button class="btn-cancel" onclick={() => showForm=false}>Cancel</button>
        <button class="btn-primary" onclick={saveForm}>Save Schedule</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .sched-page{display:flex;flex-direction:column;gap:1rem;padding:1.25rem;overflow-y:auto;min-height:100%;}
  .page-header{display:flex;align-items:flex-start;justify-content:space-between;gap:1rem;}
  .ph-left{display:flex;align-items:center;gap:.75rem;}
  .ph-icon{font-size:1.4rem;color:var(--accent);}
  h1{font-family:'Special Elite',serif;font-size:1.1rem;}
  .ph-sub{font-size:.78rem;color:var(--text-dim);margin-top:2px;}
  .btn-primary{background:var(--accent);border:none;color:var(--bg);padding:.4rem .9rem;border-radius:4px;cursor:pointer;font-family:'Rajdhani',monospace;font-size:.82rem;font-weight:700;letter-spacing:.05em;transition:opacity .15s;}
  .btn-primary:hover{opacity:.85;}
  .btn-cancel{background:none;border:1px solid var(--border);color:var(--text-dim);padding:.38rem .85rem;border-radius:3px;cursor:pointer;font-family:'Rajdhani',monospace;font-size:.82rem;}
  .btn-cancel:hover{border-color:var(--accent);color:var(--accent);}

  .info-banner{display:flex;align-items:flex-start;gap:.6rem;background:var(--bg2);border:1px solid var(--border);border-radius:6px;padding:.75rem 1rem;font-size:.8rem;color:var(--text-dim);line-height:1.6;}
  .ib-icon{color:var(--accent);flex-shrink:0;margin-top:1px;}

  .empty-state{display:flex;flex-direction:column;align-items:center;text-align:center;padding:3rem;gap:.85rem;}
  .es-icon{font-size:2.5rem;color:var(--accent);opacity:.35;}
  .empty-state h2{font-family:'Special Elite',serif;font-size:1.2rem;}
  .empty-state p{font-size:.86rem;color:var(--text-dim);max-width:400px;line-height:1.7;}

  .schedule-list{display:flex;flex-direction:column;gap:.65rem;}
  .rule-card{background:var(--bg2);border:1px solid var(--border);border-radius:6px;padding:1rem;transition:border-color .15s;}
  .rule-card:hover{border-color:var(--accent);}
  .rule-card.disabled{opacity:.55;}
  .rule-top{display:flex;align-items:center;justify-content:space-between;gap:.75rem;margin-bottom:.5rem;}
  .rule-info{display:flex;flex-direction:column;gap:2px;}
  .rule-label{font-weight:700;font-size:.92rem;}
  .rule-cron{font-family:'Courier Prime',monospace;font-size:.72rem;color:var(--accent);}
  .rule-actions{display:flex;align-items:center;gap:.4rem;}
  .action-btn{background:none;border:1px solid var(--border);color:var(--text-dim);padding:.2rem .45rem;border-radius:3px;cursor:pointer;font-size:.8rem;transition:all .15s;}
  .action-btn:hover{border-color:var(--accent);color:var(--accent);}
  .action-btn.danger:hover{border-color:#ef4444;color:#ef4444;}
  .rule-meta{display:flex;gap:1.5rem;font-size:.75rem;color:var(--text-dim);}
  .rule-meta strong{color:var(--text);}
  .accent{color:var(--accent)!important;}

  .toggle-switch{position:relative;display:flex;align-items:center;cursor:pointer;}
  .toggle-switch input{display:none;}
  .toggle-track{width:32px;height:16px;border-radius:8px;background:var(--bg3);border:1px solid var(--border);position:relative;transition:background .2s;}
  .toggle-track::after{content:'';position:absolute;top:2px;left:2px;width:10px;height:10px;border-radius:50%;background:var(--text-muted);transition:transform .2s,background .2s;}
  .toggle-switch input:checked + .toggle-track{background:var(--accent-dim);border-color:var(--accent);}
  .toggle-switch input:checked + .toggle-track::after{transform:translateX(16px);background:var(--accent);}

  .history-card{background:var(--bg2);border:1px solid var(--border);border-radius:6px;overflow:hidden;}
  .hc-hdr{padding:.6rem 1rem;background:var(--accent-dim);border-bottom:1px solid var(--border);font-size:.72rem;letter-spacing:.15em;color:var(--accent);font-weight:700;}
  .hc-body{padding:.25rem 0;}
  .hc-row{display:flex;justify-content:space-between;padding:.45rem 1rem;border-bottom:1px solid var(--border);font-size:.82rem;}
  .hc-row:last-child{border-bottom:none;}
  .hc-row span:first-child{color:var(--text-dim);}
  .hc-row span:last-child{color:var(--text);font-family:'Courier Prime',monospace;font-size:.78rem;}
  .green{color:#4ade80!important;}
  .red{color:#ef4444!important;}

  .modal-bg{position:fixed;inset:0;background:rgba(0,0,0,.75);display:flex;align-items:center;justify-content:center;z-index:1000;backdrop-filter:blur(4px);}
  .modal{background:var(--bg2);border:1px solid var(--border);border-radius:8px;padding:1.5rem;width:min(480px,90vw);display:flex;flex-direction:column;gap:1rem;}
  .modal-title{font-family:'Special Elite',serif;font-size:1rem;color:var(--accent);}
  .form-field{display:flex;flex-direction:column;gap:.4rem;}
  .form-field label{font-size:.78rem;font-weight:600;color:var(--text);}
  .field-input{background:var(--bg3);border:1px solid var(--border);color:var(--text);padding:.4rem .65rem;border-radius:4px;font-family:'Courier Prime',monospace;font-size:.82rem;outline:none;}
  .field-input:focus{border-color:var(--accent);}
  .field-hint{font-size:.72rem;color:var(--text-dim);}
  .preset-row{display:flex;flex-wrap:wrap;gap:.3rem;margin-bottom:.4rem;}
  .preset-btn{background:var(--bg3);border:1px solid var(--border);color:var(--text-dim);padding:.25rem .55rem;border-radius:3px;cursor:pointer;font-size:.72rem;font-family:'Rajdhani',monospace;transition:all .15s;}
  .preset-btn:hover{border-color:var(--accent);color:var(--accent);}
  .preset-btn.active{background:var(--accent-dim);border-color:var(--accent);color:var(--accent);}
  .toggle-label-inline{display:flex;align-items:center;gap:.5rem;cursor:pointer;font-size:.82rem;color:var(--text-dim);}
  .toggle-label-inline input{accent-color:var(--accent);}
  .modal-btns{display:flex;justify-content:flex-end;gap:.65rem;}
</style>
