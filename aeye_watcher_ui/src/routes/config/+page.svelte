<script lang="ts">
  import { store } from '$lib/stores/app.svelte';

  type Config = {
    coreUrl: string;
    scanOnLoad: boolean;
    gpuPollInterval: number;
    deepScanEnabled: boolean;
    deepScanMaxDepth: number;
    safeMode: boolean;
    minModelSizeMB: number;
    ignoreDirs: string;
    extraModelDirs: string;
    extraPorts: string;
    remoteEnabled: boolean;
    remoteHosts: string;
    remoteSshUser: string;
    remoteSshKeyPath: string;
  };

  const DEFAULTS: Config = {
    coreUrl: 'http://localhost:8080',
    scanOnLoad: true,
    gpuPollInterval: 3,
    deepScanEnabled: true,
    deepScanMaxDepth: 12,
    safeMode: true,
    minModelSizeMB: 10,
    ignoreDirs: 'node_modules,.git,.npm,Windows,System32',
    extraModelDirs: '',
    extraPorts: '',
    remoteEnabled: false,
    remoteHosts: '',
    remoteSshUser: '',
    remoteSshKeyPath: '~/.ssh/id_rsa',
  };

  function load(): Config {
    try {
      const saved = localStorage.getItem('aeye-config');
      if (saved) return { ...DEFAULTS, ...JSON.parse(saved) };
    } catch {}
    return { ...DEFAULTS };
  }

  let cfg = $state<Config>(load());
  let saved = $state(false);
  let saveTimer: ReturnType<typeof setTimeout>;

  function save() {
    localStorage.setItem('aeye-config', JSON.stringify(cfg));
    saved = true;
    clearTimeout(saveTimer);
    saveTimer = setTimeout(() => saved = false, 2000);
  }

  function reset() {
    cfg = { ...DEFAULTS };
    save();
  }
</script>

<div class="page">
  <header class="page-header">
    <div class="page-title">
      <span class="page-icon">⚙</span>
      <div>
        <h1>Configuration</h1>
        <p class="page-sub">Scanner settings, remote hosts, and behavior controls</p>
      </div>
    </div>
    <div class="header-actions">
      {#if saved}<span class="saved-badge">✓ SAVED</span>{/if}
      <button class="btn-secondary" onclick={reset}>Reset Defaults</button>
      <button class="btn-primary" onclick={save}>Save Settings</button>
    </div>
  </header>

  <div class="content">
    <div class="config-grid">

      <!-- Core Connection -->
      <div class="config-card">
        <div class="card-header-cfg"><span class="cfg-icon">⬡</span> Core Connection</div>
        <div class="field-group">
          <label class="field-label">Core API URL</label>
          <input class="field-input" type="text" bind:value={cfg.coreUrl} placeholder="http://localhost:8080" />
          <span class="field-hint">Address of the aeye_watcher_core Go server</span>
        </div>
        <div class="field-group">
          <label class="field-label toggle-label">
            <input type="checkbox" bind:checked={cfg.scanOnLoad} />
            <span>Auto-scan on page load</span>
          </label>
        </div>
        <div class="field-group">
          <label class="field-label">GPU Poll Interval (seconds)</label>
          <input class="field-input narrow" type="number" min="1" max="60" bind:value={cfg.gpuPollInterval} />
        </div>
      </div>

      <!-- Scan Behavior -->
      <div class="config-card">
        <div class="card-header-cfg"><span class="cfg-icon">◆</span> Scan Behavior</div>
        <div class="field-group">
          <label class="field-label toggle-label">
            <input type="checkbox" bind:checked={cfg.deepScanEnabled} />
            <span>Enable deep filesystem scan</span>
          </label>
          <span class="field-hint">Walks home directory and /opt looking for model files outside known provider dirs</span>
        </div>
        <div class="field-group">
          <label class="field-label">Max scan depth</label>
          <input class="field-input narrow" type="number" min="3" max="20" bind:value={cfg.deepScanMaxDepth} />
        </div>
        <div class="field-group">
          <label class="field-label">Minimum model file size (MB)</label>
          <input class="field-input narrow" type="number" min="1" bind:value={cfg.minModelSizeMB} />
          <span class="field-hint">Files smaller than this are ignored (avoids test stubs)</span>
        </div>
        <div class="field-group">
          <label class="field-label">Additional model directories (one per line)</label>
          <textarea class="field-textarea" rows="3" bind:value={cfg.extraModelDirs} placeholder="/data/models&#10;D:\AI\models"></textarea>
        </div>
        <div class="field-group">
          <label class="field-label">Directories to skip (comma-separated)</label>
          <input class="field-input" type="text" bind:value={cfg.ignoreDirs} />
        </div>
      </div>

      <!-- Safety & Cleanup -->
      <div class="config-card">
        <div class="card-header-cfg"><span class="cfg-icon">⚠</span> Safety & Cleanup</div>
        <div class="field-group">
          <label class="field-label toggle-label">
            <input type="checkbox" bind:checked={cfg.safeMode} />
            <span>Safe mode (no actual file deletion)</span>
          </label>
          <span class="field-hint warn">
            When enabled, cleanup requests are only logged in the Go terminal. Disable with caution — also requires <code>AEYE_UNSAFE_DELETE=1</code> env var on the core server.
          </span>
        </div>
        {#if !cfg.safeMode}
          <div class="warn-box">
            ⚠ Safe mode is OFF. File deletions will be permanent. The Go core must also have <code>AEYE_UNSAFE_DELETE=1</code> set.
          </div>
        {/if}
      </div>

      <!-- Port Probing -->
      <div class="config-card">
        <div class="card-header-cfg"><span class="cfg-icon">◉</span> Port Probing</div>
        <div class="field-group">
          <label class="field-label">Extra ports to probe (comma-separated)</label>
          <input class="field-input" type="text" bind:value={cfg.extraPorts} placeholder="9876,12345" />
          <span class="field-hint">Ports beyond the built-in list to check for running AI services</span>
        </div>
      </div>

      <!-- Remote Scanning -->
      <div class="config-card full-width">
        <div class="card-header-cfg">
          <span class="cfg-icon">◈</span> Remote System Scanning
          <span class="beta-badge">EXPERIMENTAL</span>
        </div>
        <div class="field-group">
          <label class="field-label toggle-label">
            <input type="checkbox" bind:checked={cfg.remoteEnabled} />
            <span>Enable remote host scanning</span>
          </label>
          <span class="field-hint">
            Connects to remote systems via SSH and runs the aeye_watcher_core binary there, streaming results back.
            The remote host must have the core binary installed and accessible via SSH.
          </span>
        </div>
        {#if cfg.remoteEnabled}
          <div class="remote-fields">
            <div class="field-group">
              <label class="field-label">Remote Hosts (one per line)</label>
              <textarea class="field-textarea" rows="4" bind:value={cfg.remoteHosts}
                placeholder="192.168.1.100&#10;myserver.example.com&#10;user@10.0.0.50"></textarea>
              <span class="field-hint">Format: hostname or user@hostname. SSH key auth required.</span>
            </div>
            <div class="field-row">
              <div class="field-group">
                <label class="field-label">SSH User (default)</label>
                <input class="field-input" type="text" bind:value={cfg.remoteSshUser} placeholder="ubuntu" />
              </div>
              <div class="field-group">
                <label class="field-label">SSH Key Path</label>
                <input class="field-input" type="text" bind:value={cfg.remoteSshKeyPath} placeholder="~/.ssh/id_rsa" />
              </div>
            </div>
            <div class="info-box">
              <strong>How remote scanning works:</strong> The core server SSH's into each host, invokes
              <code>aeye-core --mode=scan --json</code>, and streams the result back as a WebSocket event.
              Install the binary on each remote host: <code>curl -sSL https://github.com/NathanielRand/aEyeWatcher/releases/latest/download/install.sh | sh</code>
            </div>
          </div>
        {/if}
      </div>

    </div>
  </div>
</div>

<style>
  .page{display:flex;flex-direction:column;height:100%;overflow:hidden;}
  .page-header{
    display:flex;align-items:center;justify-content:space-between;
    padding:.85rem 1.5rem;background:var(--bg2);border-bottom:1px solid var(--border);
    flex-shrink:0;flex-wrap:wrap;gap:.75rem;
  }
  .page-title{display:flex;align-items:center;gap:.85rem;}
  .page-icon{font-size:1.4rem;color:var(--accent);}
  .page-title h1{font-family:'Special Elite',serif;font-size:1.15rem;color:var(--text);}
  .page-sub{font-size:.8rem;color:var(--text-dim);margin-top:2px;}
  .header-actions{display:flex;align-items:center;gap:.75rem;}
  .saved-badge{font-size:.72rem;color:#4ade80;letter-spacing:.1em;font-weight:700;}
  .btn-primary{background:var(--accent);border:none;color:var(--bg);padding:.4rem 1rem;border-radius:3px;cursor:pointer;font-family:'Rajdhani',monospace;font-size:.85rem;font-weight:700;letter-spacing:.05em;transition:opacity .15s;}
  .btn-primary:hover{opacity:.85;}
  .btn-secondary{background:none;border:1px solid var(--border);color:var(--text-dim);padding:.4rem 1rem;border-radius:3px;cursor:pointer;font-family:'Rajdhani',monospace;font-size:.85rem;transition:all .15s;}
  .btn-secondary:hover{border-color:var(--accent);color:var(--accent);}

  .content{flex:1;overflow-y:auto;padding:1.5rem;scrollbar-width:thin;scrollbar-color:var(--border) transparent;}

  .config-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(360px,1fr));gap:1rem;}
  .full-width{grid-column:1/-1;}

  .config-card{background:var(--bg2);border:1px solid var(--border);border-radius:8px;overflow:hidden;}
  .card-header-cfg{
    display:flex;align-items:center;gap:.5rem;
    padding:.7rem 1rem;background:var(--accent-dim);border-bottom:1px solid var(--border);
    font-size:.75rem;letter-spacing:.15em;color:var(--accent);font-weight:700;
    font-family:'Rajdhani',monospace;
  }
  .cfg-icon{font-size:.9rem;}
  .beta-badge{margin-left:auto;font-size:.6rem;background:rgba(251,146,60,.2);border:1px solid #fb923c;color:#fb923c;padding:1px 6px;border-radius:2px;letter-spacing:.1em;}

  .field-group{padding:.85rem 1rem;border-bottom:1px solid var(--border);}
  .field-group:last-child{border-bottom:none;}
  .field-label{display:block;font-size:.82rem;font-weight:600;color:var(--text);margin-bottom:.4rem;letter-spacing:.03em;}
  .toggle-label{display:flex;align-items:center;gap:.5rem;cursor:pointer;}
  .toggle-label input[type="checkbox"]{accent-color:var(--accent);width:14px;height:14px;}
  .field-input{
    width:100%;background:var(--bg3);border:1px solid var(--border);color:var(--text);
    padding:.4rem .65rem;border-radius:4px;font-family:'Courier Prime',monospace;font-size:.82rem;
    outline:none;transition:border-color .15s;
  }
  .field-input:focus{border-color:var(--accent);}
  .field-input.narrow{width:100px;}
  .field-textarea{
    width:100%;background:var(--bg3);border:1px solid var(--border);color:var(--text);
    padding:.4rem .65rem;border-radius:4px;font-family:'Courier Prime',monospace;font-size:.8rem;
    outline:none;resize:vertical;transition:border-color .15s;
  }
  .field-textarea:focus{border-color:var(--accent);}
  .field-hint{display:block;font-size:.75rem;color:var(--text-dim);margin-top:.35rem;line-height:1.5;}
  .field-hint.warn{color:#fb923c;}
  .field-hint code{background:var(--bg3);padding:1px 4px;border-radius:2px;font-size:.75rem;}

  .warn-box{margin:.5rem 1rem .85rem;padding:.65rem .85rem;background:rgba(251,146,60,.1);border:1px solid #fb923c;border-radius:4px;font-size:.8rem;color:#fb923c;line-height:1.5;}
  .warn-box code{background:rgba(251,146,60,.15);padding:1px 4px;border-radius:2px;}

  .remote-fields{padding:.25rem 0;}
  .field-row{display:grid;grid-template-columns:1fr 1fr;gap:0;}

  .info-box{margin:.5rem 1rem .85rem;padding:.65rem .85rem;background:var(--bg3);border:1px solid var(--border);border-radius:4px;font-size:.8rem;color:var(--text-dim);line-height:1.7;}
  .info-box strong{color:var(--text);}
  .info-box code{background:var(--bg2);padding:2px 5px;border-radius:2px;font-family:'Courier Prime',monospace;font-size:.75rem;}

  @media(max-width:640px){
    .config-grid{grid-template-columns:1fr;}
    .field-row{grid-template-columns:1fr;}
    .page-header{flex-direction:column;align-items:flex-start;}
  }
</style>
