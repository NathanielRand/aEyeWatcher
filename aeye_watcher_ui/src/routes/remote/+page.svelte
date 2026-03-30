<script lang="ts">
  import ProGate from '$lib/components/ProGate.svelte';
  import { store, addToast } from '$lib/stores/app.svelte';

  type RemoteHost = { id:string; label:string; host:string; user:string; keyPath:string; status:'idle'|'scanning'|'done'|'error'; lastScan?:string };

  let hosts = $state<RemoteHost[]>([]);
  let showForm = $state(false);
  let form = $state({ label:'', host:'', user:'', keyPath:'~/.ssh/id_rsa' });
  let selectedResult = $state<string|null>(null);

  function addHost() {
    if (!form.host.trim()) { addToast('Host is required', 'warn'); return; }
    hosts = [...hosts, { id: Math.random().toString(36).slice(2), ...form, status:'idle' }];
    showForm = false;
    form = { label:'', host:'', user:'', keyPath:'~/.ssh/id_rsa' };
    addToast('Remote host added', 'success');
  }

  async function triggerRemoteScan(id: string) {
    hosts = hosts.map(h => h.id===id ? {...h, status:'scanning'} : h);
    addToast('Remote scan triggered — requires aeye-core on remote host', 'info');
    // In production: POST /api/remote/scan with host config
    await new Promise(r => setTimeout(r, 2500));
    hosts = hosts.map(h => h.id===id ? {...h, status:'done', lastScan: new Date().toISOString()} : h);
    addToast('Remote scan complete (demo)', 'success');
  }
</script>

<ProGate feature="Remote Host Scanning">
<div class="page">
  <div class="page-hdr">
    <span class="hdr-icon">◈</span>
    <div>
      <h1>Remote Scanning</h1>
      <p class="hdr-sub">Scan AI ecosystems on remote systems via SSH</p>
    </div>
    <button class="btn-primary" onclick={() => showForm=!showForm}>+ Add Host</button>
  </div>

  <div class="info-banner">
    <span>◈</span>
    <div>
      Remote scanning requires <code>aeye-core</code> installed on each host. The local core SSH's in, runs a scan, and streams results back.
      Install on remote: <code>curl -sSL https://github.com/NathanielRand/aEyeWatcher/releases/latest/download/install.sh | sh</code>
    </div>
  </div>

  {#if showForm}
    <div class="form-card">
      <div class="form-title">Add Remote Host</div>
      <div class="form-fields">
        <div class="ff"><label>Label</label><input class="fi" type="text" bind:value={form.label} placeholder="My GPU Server"/></div>
        <div class="ff"><label>Host / IP</label><input class="fi" type="text" bind:value={form.host} placeholder="192.168.1.100 or server.example.com"/></div>
        <div class="ff"><label>SSH User</label><input class="fi" type="text" bind:value={form.user} placeholder="ubuntu"/></div>
        <div class="ff"><label>SSH Key Path</label><input class="fi" type="text" bind:value={form.keyPath}/></div>
      </div>
      <div class="form-btns">
        <button class="btn-cancel" onclick={() => showForm=false}>Cancel</button>
        <button class="btn-primary" onclick={addHost}>Add Host</button>
      </div>
    </div>
  {/if}

  {#if hosts.length === 0 && !showForm}
    <div class="empty">
      <span class="ei">◈</span>
      <h2>No remote hosts configured</h2>
      <p>Add a remote host to scan its AI ecosystem from this dashboard.</p>
      <button class="btn-primary" onclick={() => showForm=true}>+ Add First Host</button>
    </div>
  {:else if hosts.length > 0}
    <div class="host-list">
      {#each hosts as host}
        <div class="host-card">
          <div class="hc-header">
            <div class="hc-dot {host.status}"></div>
            <div class="hc-info">
              <span class="hc-label">{host.label || host.host}</span>
              <span class="hc-host">{host.user ? host.user + '@' : ''}{host.host}</span>
            </div>
            <div class="hc-status-badge {host.status}">
              {host.status === 'scanning' ? '⟳ SCANNING' : host.status === 'done' ? '✓ COMPLETE' : host.status === 'error' ? '✗ ERROR' : '○ IDLE'}
            </div>
          </div>
          {#if host.lastScan}
            <div class="hc-last">Last scan: {new Date(host.lastScan).toLocaleString()}</div>
          {/if}
          <div class="hc-actions">
            <button class="action-btn" onclick={() => triggerRemoteScan(host.id)} disabled={host.status==='scanning'}>
              {host.status === 'scanning' ? '⟳ Scanning…' : '▶ Scan Now'}
            </button>
            <button class="action-btn secondary" onclick={() => hosts = hosts.filter(h=>h.id!==host.id)}>Remove</button>
          </div>
        </div>
      {/each}
    </div>
  {/if}

  <div class="arch-card">
    <div class="arch-hdr">◆ HOW REMOTE SCANNING WORKS</div>
    <div class="arch-steps">
      <div class="arch-step"><span class="as-num">01</span><div><strong>Install aeye-core on remote</strong><p>The remote host needs the Go binary accessible. Use the install script or copy the binary manually.</p></div></div>
      <div class="arch-step"><span class="as-num">02</span><div><strong>Trigger from this dashboard</strong><p>Click "Scan Now" — the local core SSH's into the remote host using your configured key and runs the scan.</p></div></div>
      <div class="arch-step"><span class="as-num">03</span><div><strong>Results stream back</strong><p>The remote core outputs JSON via stdout over the SSH tunnel. Results appear in this dashboard just like a local scan.</p></div></div>
    </div>
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
  .info-banner{background:var(--bg2);border:1px solid var(--border);border-radius:6px;padding:.75rem 1rem;font-size:.82rem;color:var(--text-dim);display:flex;gap:.6rem;line-height:1.7;}
  .info-banner code{font-size:.75rem;}
  .form-card{background:var(--bg2);border:1px solid var(--accent);border-radius:8px;padding:1.25rem;display:flex;flex-direction:column;gap:.85rem;}
  .form-title{font-family:'Special Elite',serif;font-size:.95rem;color:var(--accent);}
  .form-fields{display:grid;grid-template-columns:1fr 1fr;gap:.65rem;}
  .ff{display:flex;flex-direction:column;gap:.3rem;}
  .ff label{font-size:.78rem;font-weight:600;color:var(--text);}
  .fi{background:var(--bg3);border:1px solid var(--border);color:var(--text);padding:.38rem .65rem;border-radius:4px;font-family:'Rajdhani',monospace;font-size:.84rem;outline:none;}
  .fi:focus{border-color:var(--accent);}
  .form-btns{display:flex;justify-content:flex-end;gap:.5rem;}
  .empty{display:flex;flex-direction:column;align-items:center;text-align:center;padding:3rem;gap:.85rem;}
  .ei{font-size:2.5rem;color:var(--accent);opacity:.3;}
  .empty h2{font-family:'Special Elite',serif;font-size:1.2rem;color:var(--text);}
  .empty p{font-size:.88rem;color:var(--text-dim);max-width:400px;line-height:1.7;}
  .host-list{display:flex;flex-direction:column;gap:.65rem;}
  .host-card{background:var(--bg2);border:1px solid var(--border);border-radius:8px;padding:1rem;}
  .hc-header{display:flex;align-items:center;gap:.75rem;margin-bottom:.5rem;}
  .hc-dot{width:10px;height:10px;border-radius:50%;flex-shrink:0;}
  .hc-dot.idle{background:var(--text-muted);}
  .hc-dot.scanning{background:var(--accent);animation:glow 1s ease-in-out infinite;}
  .hc-dot.done{background:#4ade80;}
  .hc-dot.error{background:#ef4444;}
  @keyframes glow{0%,100%{opacity:1}50%{opacity:.4}}
  .hc-info{flex:1;}
  .hc-label{display:block;font-weight:700;font-size:.9rem;color:var(--text);}
  .hc-host{font-size:.72rem;color:var(--text-dim);font-family:'Courier Prime',monospace;}
  .hc-status-badge{font-size:.62rem;font-weight:700;letter-spacing:.1em;padding:2px 7px;border-radius:3px;background:var(--bg3);border:1px solid var(--border);color:var(--text-dim);}
  .hc-status-badge.done{background:rgba(74,222,128,.12);border-color:#4ade80;color:#4ade80;}
  .hc-status-badge.scanning{background:var(--accent-dim);border-color:var(--accent);color:var(--accent);}
  .hc-status-badge.error{background:rgba(239,68,68,.12);border-color:#ef4444;color:#ef4444;}
  .hc-last{font-size:.7rem;color:var(--text-muted);font-family:'Courier Prime',monospace;margin-bottom:.6rem;}
  .hc-actions{display:flex;gap:.5rem;}
  .action-btn{background:var(--accent-dim);border:1px solid var(--accent);color:var(--accent);padding:.3rem .7rem;border-radius:4px;cursor:pointer;font-family:'Rajdhani',monospace;font-size:.8rem;font-weight:600;transition:all .15s;}
  .action-btn:hover:not(:disabled){background:var(--accent);color:var(--bg);}
  .action-btn:disabled{opacity:.5;cursor:not-allowed;}
  .action-btn.secondary{background:none;border-color:var(--border);color:var(--text-dim);}
  .action-btn.secondary:hover{border-color:var(--accent);color:var(--accent);}
  .arch-card{background:var(--bg2);border:1px solid var(--border);border-radius:8px;overflow:hidden;}
  .arch-hdr{padding:.6rem 1rem;background:var(--accent-dim);border-bottom:1px solid var(--border);font-size:.7rem;letter-spacing:.15em;color:var(--accent);font-weight:700;}
  .arch-steps{display:flex;flex-direction:column;gap:0;}
  .arch-step{display:flex;gap:1rem;padding:.85rem 1rem;border-bottom:1px solid var(--border);align-items:flex-start;}
  .arch-step:last-child{border-bottom:none;}
  .as-num{font-family:'Courier Prime',monospace;font-size:.75rem;color:var(--text-muted);flex-shrink:0;width:24px;}
  .arch-step strong{display:block;font-size:.88rem;color:var(--text);margin-bottom:.2rem;}
  .arch-step p{font-size:.8rem;color:var(--text-dim);line-height:1.6;}
  @media(max-width:640px){.form-fields{grid-template-columns:1fr;}}
</style>
