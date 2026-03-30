<script lang="ts">
  let activeSection = $state('overview');

  const sections = [
    { id: 'overview',      label: 'Overview' },
    { id: 'install',       label: 'Installation' },
    { id: 'quickstart',    label: 'Quick Start' },
    { id: 'architecture',  label: 'Architecture' },
    { id: 'scanner',       label: 'What It Scans' },
    { id: 'config',        label: 'Configuration' },
    { id: 'contributing',  label: 'Contributing' },
    { id: 'faq',           label: 'FAQ' },
  ];
</script>

<div class="page">
  <header class="page-header">
    <div class="page-title">
      <span class="page-icon">◈</span>
      <div>
        <h1>Documentation</h1>
        <p class="page-sub">Installation, architecture, and contribution guides</p>
      </div>
    </div>
    <a class="gh-btn" href="https://github.com/NathanielRand/aEyeWatcher" target="_blank" rel="noopener">
      ⬡ View on GitHub
    </a>
  </header>

  <div class="doc-layout">
    <nav class="doc-nav">
      {#each sections as sec}
        <button
          class="doc-nav-item"
          class:active={activeSection === sec.id}
          onclick={() => activeSection = sec.id}
        >{sec.label}</button>
      {/each}
    </nav>

    <article class="doc-content">

      {#if activeSection === 'overview'}
        <h2>What is aEye Watcher?</h2>
        <p>aEye Watcher is a full-stack, open-source AI ecosystem monitor — a living command center that scans your local (and optionally remote) systems for every AI-related component: providers, models, agents, MCP servers, ACP agents, vector databases, gateways, skills, soul files, and more.</p>
        <p>Think of it as a detective that never sleeps, logging every AI footprint on your machine so you know exactly what's installed, what's running, where your model weights live, and what's quietly orphaned.</p>

        <div class="feature-grid">
          <div class="feature-card"><span class="fc-icon">⚙</span><strong>30+ Components</strong><span>Providers, agents, gateways, vector DBs, GUIs, TTS engines</span></div>
          <div class="feature-card"><span class="fc-icon">◆</span><strong>Deep Model Hunt</strong><span>GGUF, Safetensors, ONNX, PyTorch — across your whole disk</span></div>
          <div class="feature-card"><span class="fc-icon">◈</span><strong>Agent Cards</strong><span>Character sheets per ecosystem with soul/memory file checklists</span></div>
          <div class="feature-card"><span class="fc-icon">⬡</span><strong>MCP + ACP</strong><span>Discovers server configs and probes live ACP agent endpoints</span></div>
          <div class="feature-card"><span class="fc-icon">◉</span><strong>Disk & GPU</strong><span>Real-time VRAM usage and storage volume mapping</span></div>
          <div class="feature-card"><span class="fc-icon">▸</span><strong>Cross-Platform</strong><span>Linux, macOS, Windows — detects OS on start</span></div>
        </div>

      {:else if activeSection === 'install'}
        <h2>Installation</h2>
        <h3>Requirements</h3>
        <ul>
          <li><strong>Go 1.22+</strong> — for the core backend</li>
          <li><strong>Node 18+</strong> — for the UI (optional if embedding the built UI)</li>
          <li><code>nvidia-smi</code> — optional, for NVIDIA GPU monitoring</li>
        </ul>

        <h3>Option A — Pre-built Binaries (Recommended)</h3>
        <div class="code-block">
          <div class="code-label">Linux / macOS</div>
          <pre>curl -sSL https://github.com/NathanielRand/aEyeWatcher/releases/latest/download/install.sh | sh</pre>
        </div>
        <div class="code-block">
          <div class="code-label">Windows (PowerShell)</div>
          <pre>irm https://github.com/NathanielRand/aEyeWatcher/releases/latest/download/install.ps1 | iex</pre>
        </div>

        <h3>Option B — Build from Source</h3>
        <div class="code-block">
          <div class="code-label">Clone & build</div>
          <pre>git clone https://github.com/NathanielRand/aEyeWatcher
cd aEyeWatcher

# Build the Go core
cd aeye_watcher_core
go mod tidy
go build -o aeye-core main.go

# Build & serve the UI
cd ../aeye_watcher_ui
npm install
npm run build</pre>
        </div>

        <h3>Option C — Docker</h3>
        <div class="code-block">
          <div class="code-label">Docker</div>
          <pre>docker run -p 8080:8080 -v $HOME:/scan-home:ro ghcr.io/nathanielrand/aeye-watcher:latest</pre>
        </div>

      {:else if activeSection === 'quickstart'}
        <h2>Quick Start</h2>
        <h3>1. Start the Core</h3>
        <div class="code-block">
          <pre>cd aeye_watcher_core
go run main.go
# or: ./aeye-core
# → Core API: http://localhost:8080</pre>
        </div>
        <h3>2. Start the UI</h3>
        <div class="code-block">
          <pre>cd aeye_watcher_ui
npm run dev
# → http://localhost:5173</pre>
        </div>
        <h3>3. Run a Scan</h3>
        <p>The dashboard auto-scans on load. Hit the <strong>◉ SCAN</strong> button in the topbar to re-scan at any time. Watch the overlay for live phase-by-phase progress.</p>

        <h3>Custom Port</h3>
        <div class="code-block">
          <pre>./aeye-core -port 9090</pre>
        </div>
        <p>Then update <strong>Configuration → Core API URL</strong> in the UI to match.</p>

        <h3>Enable Real Deletion</h3>
        <div class="code-block">
          <div class="code-label">Safe mode is ON by default. To enable:</div>
          <pre>AEYE_UNSAFE_DELETE=1 ./aeye-core</pre>
        </div>
        <p>Then toggle safe mode off in <strong>Configuration → Safety &amp; Cleanup</strong>.</p>

      {:else if activeSection === 'architecture'}
        <h2>Architecture</h2>
        <p>The two apps are fully decoupled. The core exposes REST + WebSocket; any client can consume it.</p>
        <div class="arch-diagram">
          <div class="arch-box client">SvelteKit UI<br><span>localhost:5173</span></div>
          <div class="arch-arrow">⟷ REST + WS</div>
          <div class="arch-box core">Go Core<br><span>localhost:8080</span></div>
          <div class="arch-arrow">↓ scans</div>
          <div class="arch-box system">Local System<br><span>filesystem, ports, processes</span></div>
        </div>

        <h3>Core Internal Packages</h3>
        <div class="pkg-table">
          <div class="pkg-row"><code>internal/models</code><span>Shared types mirrored in TypeScript</span></div>
          <div class="pkg-row"><code>internal/os_detect</code><span>OS, arch, hostname, GPU availability</span></div>
          <div class="pkg-row"><code>internal/disk</code><span>Cross-platform volume scanner</span></div>
          <div class="pkg-row"><code>internal/scanner/components</code><span>60+ CLI + port checks</span></div>
          <div class="pkg-row"><code>internal/scanner/models</code><span>Ollama / LM Studio / HF / deep walk</span></div>
          <div class="pkg-row"><code>internal/scanner/mcp_acp</code><span>MCP config parsing + ACP probing</span></div>
          <div class="pkg-row"><code>internal/scanner/agents</code><span>Ecosystem cards + memory file hunting</span></div>
          <div class="pkg-row"><code>internal/scanner/gpu</code><span>NVIDIA / AMD / Apple Silicon stats</span></div>
          <div class="pkg-row"><code>internal/server</code><span>HTTP + WebSocket scan engine</span></div>
        </div>

        <h3>API Endpoints</h3>
        <div class="pkg-table">
          <div class="pkg-row"><code>WS /ws/scan</code><span>Live scan — streams ScanEvent JSON</span></div>
          <div class="pkg-row"><code>GET /api/gpu</code><span>Current GPU stats (pollable every N sec)</span></div>
          <div class="pkg-row"><code>GET /api/os</code><span>OS / system info</span></div>
          <div class="pkg-row"><code>GET /api/scan/cached</code><span>Last scan result as JSON</span></div>
          <div class="pkg-row"><code>POST /api/delete</code><span>Request model file cleanup</span></div>
          <div class="pkg-row"><code>GET /health</code><span>Core health check</span></div>
        </div>

      {:else if activeSection === 'scanner'}
        <h2>What aEye Watcher Scans</h2>

        <h3>AI Providers</h3>
        <p>Checks both CLI presence and live port: Ollama, LM Studio, LocalAI, llama.cpp server, KoboldCPP, vLLM, Oobabooga, Jan, GPT4All, Llamafile.</p>

        <h3>Agents</h3>
        <p>Claude Code, OpenClaw (+ gateway), Hermes (+ gateway), OpenInterpreter, Aider, AutoGPT, CrewAI, LangChain, Goose, SWE-agent, Devin, Amp.</p>

        <h3>Gateways</h3>
        <p>LiteLLM, Portkey, HelixML, OpenClaw Gateway (auto-started on boot), Hermes Gateway — all probed on their known ports plus alternates.</p>

        <h3>Vector Databases</h3>
        <p>ChromaDB, Qdrant, Weaviate, Milvus, LanceDB, pgvector.</p>

        <h3>Model Files</h3>
        <p>Extensions hunted: <code>.gguf .ggml .safetensors .onnx .pt .pth .bin .pkl .mlmodel .tflite .engine .llamafile</code>. Provider-aware: Ollama manifests (layer-accurate sizes), LM Studio hub (<code>~/.lmstudio/hub/models/{company}/{model}</code>), HuggingFace hub cache. Files outside known provider directories are tagged <strong>orphaned</strong>.</p>

        <h3>Agent Ecosystem Files</h3>
        <p>For each supported agent, aEye checks for: soul files (<code>SOUL.md</code>), persona files (<code>USER.md, PERSONA.md, IDENTITY.md</code>), memory files (<code>MEMORY.md, SUMMARY.md, CONTEXT.md, HISTORY.md</code>), agent definitions (<code>AGENTS.md, RULES.md, INSTRUCTIONS.md</code>), skills directories, tool configs, and MCP server configs.</p>

        <h3>MCP Servers</h3>
        <p>Parses <code>claude_desktop_config.json</code>, Cursor settings, Windsurf config, Zed settings, Continue config, and any <code>mcp*.json</code> in your home directory. Also probes common MCP ports.</p>

        <h3>ACP Agents</h3>
        <p>Probes <code>/.well-known/agent.json</code> on standard agent ports per the ACP spec.</p>

        <h3>Disk Volumes</h3>
        <p>All real mounted volumes with accurate used/free/total — <code>df -k</code> on Unix, PowerShell <code>Get-PSDrive</code> (with <code>wmic</code> fallback) on Windows.</p>

      {:else if activeSection === 'config'}
        <h2>Configuration Reference</h2>
        <p>Settings are persisted to <code>localStorage</code> and applied to the UI. The Go core reads its own environment variables.</p>

        <div class="pkg-table">
          <div class="pkg-row header"><strong>UI Setting</strong><strong>Description</strong></div>
          <div class="pkg-row"><code>Core API URL</code><span>Where the Go core is running</span></div>
          <div class="pkg-row"><code>Auto-scan on load</code><span>Triggers a scan when the UI opens</span></div>
          <div class="pkg-row"><code>GPU Poll Interval</code><span>Seconds between VRAM polls</span></div>
          <div class="pkg-row"><code>Deep Scan</code><span>Walk filesystem for orphaned model files</span></div>
          <div class="pkg-row"><code>Max Depth</code><span>Directory recursion limit for deep scan</span></div>
          <div class="pkg-row"><code>Min Model Size</code><span>Ignore model files below this (MB)</span></div>
          <div class="pkg-row"><code>Extra Model Dirs</code><span>Additional paths to include as non-orphaned</span></div>
          <div class="pkg-row"><code>Ignore Dirs</code><span>Directory names to skip entirely</span></div>
          <div class="pkg-row"><code>Safe Mode</code><span>Prevents actual file deletion</span></div>
          <div class="pkg-row"><code>Extra Ports</code><span>Additional ports to probe for AI services</span></div>
        </div>

        <h3>Environment Variables (Core)</h3>
        <div class="pkg-table">
          <div class="pkg-row"><code>AEYE_UNSAFE_DELETE=1</code><span>Enable real file deletion (requires safe mode off in UI too)</span></div>
          <div class="pkg-row"><code>PORT=8080</code><span>Override default port (same as -port flag)</span></div>
        </div>

      {:else if activeSection === 'contributing'}
        <h2>Contributing</h2>
        <p>aEye Watcher is fully open-source under the MIT license. Contributions of all kinds are welcome.</p>

        <h3>Getting Started</h3>
        <div class="code-block">
          <pre>git clone https://github.com/NathanielRand/aEyeWatcher
cd aEyeWatcher
# Core
cd aeye_watcher_core && go mod tidy && go run main.go
# UI (new terminal)
cd aeye_watcher_ui && npm install && npm run dev</pre>
        </div>

        <h3>Ways to Contribute</h3>
        <ul>
          <li><strong>Add a new component</strong> → append to <code>componentDefs</code> in <code>internal/scanner/components.go</code></li>
          <li><strong>Add a new agent ecosystem</strong> → append to <code>getAgentDefs()</code> in <code>internal/scanner/agents.go</code></li>
          <li><strong>Add a new model extension</strong> → append to <code>modelExtensions</code> in <code>internal/scanner/models.go</code></li>
          <li><strong>Add a new MCP config location</strong> → append to <code>getMCPConfigPaths()</code> in <code>internal/scanner/mcp_acp.go</code></li>
          <li><strong>UI improvements</strong> → everything is in <code>aeye_watcher_ui/src/routes/</code></li>
          <li><strong>Bug reports & feature requests</strong> → open a GitHub Issue</li>
        </ul>

        <h3>Code Style</h3>
        <ul>
          <li>Go: standard <code>gofmt</code> formatting, no external deps beyond gorilla/websocket</li>
          <li>Svelte: Svelte 5 runes mode, TypeScript strict, Tailwind utilities or scoped CSS variables</li>
          <li>Keep the decoupled architecture — no coupling between the Go core and the UI beyond the JSON API</li>
        </ul>

        <h3>Roadmap</h3>
        <ul>
          <li>Remote host scanning via SSH tunnel</li>
          <li>Scheduled / background scanning with change detection alerts</li>
          <li>Plugin system for custom scanners</li>
          <li>Export scan reports as JSON / PDF</li>
          <li>Desktop app (Tauri wrapper)</li>
          <li>Team dashboard for shared AI infra visibility</li>
        </ul>

      {:else if activeSection === 'faq'}
        <h2>FAQ</h2>

        <div class="faq-item">
          <div class="faq-q">Does aEye Watcher upload anything to the internet?</div>
          <div class="faq-a">No. Everything runs locally. The Go core talks only to localhost ports and the local filesystem. No telemetry, no analytics, no external calls.</div>
        </div>
        <div class="faq-item">
          <div class="faq-q">Why are my disk used/free values wrong?</div>
          <div class="faq-a">On Windows, the core uses PowerShell <code>Get-PSDrive</code> (falling back to <code>wmic</code>). On Unix, it uses <code>df -k</code>. If values still look off, check which drives are being reported — pseudo-filesystems are filtered out but some edge cases exist.</div>
        </div>
        <div class="faq-item">
          <div class="faq-q">Why aren't my LM Studio models showing up?</div>
          <div class="faq-a">On Windows, LM Studio stores models in <code>~\.lmstudio\hub\models\{company}\{model}\</code>. The scanner checks this path directly. If your models are elsewhere, add the path in <strong>Configuration → Extra Model Directories</strong>.</div>
        </div>
        <div class="faq-item">
          <div class="faq-q">The scan seems frozen — how do I tell if it's still running?</div>
          <div class="faq-a">The scan overlay shows a phase step strip, a live last-message line below the progress bar, and a shimmer animation on the bar itself. If all three are static for more than 30 seconds, the WebSocket may have dropped — refresh and re-scan.</div>
        </div>
        <div class="faq-item">
          <div class="faq-q">Can I run the UI without the Go core?</div>
          <div class="faq-a">The UI will load but all data will be empty. It degrades gracefully — no crashes, just empty states everywhere until the core is reachable.</div>
        </div>
        <div class="faq-item">
          <div class="faq-q">Is there a desktop app?</div>
          <div class="faq-a">Not yet, but a Tauri wrapper is on the roadmap. Since the UI uses <code>adapter-static</code> compatible patterns, the migration is straightforward.</div>
        </div>
      {/if}

    </article>
  </div>
</div>

<style>
  .page{display:flex;flex-direction:column;min-height:100%;overflow:hidden;}
  .page-header{display:flex;align-items:center;justify-content:space-between;padding:.85rem 1.5rem;background:var(--bg2);border-bottom:1px solid var(--border);flex-shrink:0;}
  .page-title{display:flex;align-items:center;gap:.85rem;}
  .page-icon{font-size:1.4rem;color:var(--accent);}
  .page-title h1{font-family:'Special Elite',serif;font-size:1.15rem;}
  .page-sub{font-size:.8rem;color:var(--text-dim);margin-top:2px;}
  .gh-btn{background:var(--accent-dim);border:1px solid var(--accent);color:var(--accent);padding:.4rem .9rem;border-radius:3px;font-size:.82rem;letter-spacing:.05em;transition:all .15s;}
  .gh-btn:hover{background:var(--accent);color:var(--bg);}

  .doc-layout{display:flex;flex:1;overflow:hidden;}

  .doc-nav{
    width:180px;min-width:180px;
    background:var(--bg2);border-right:1px solid var(--border);
    padding:.6rem;display:flex;flex-direction:column;gap:2px;
    overflow-y:auto;
  }
  .doc-nav-item{
    background:none;border:none;color:var(--text-dim);text-align:left;
    padding:.5rem .65rem;border-radius:4px;cursor:pointer;
    font-family:'Rajdhani',monospace;font-size:.88rem;font-weight:500;
    transition:all .15s;
  }
  .doc-nav-item:hover{background:var(--accent-dim);color:var(--text);}
  .doc-nav-item.active{background:var(--accent-dim);color:var(--accent);border-left:2px solid var(--accent);}

  .doc-content{overflow-y:auto!important;
    flex:1;overflow-y:auto;min-height:0;padding:2rem;
    scrollbar-width:thin;scrollbar-color:var(--border) transparent;
    max-width:820px;
  }

  .doc-content :global(h2){font-family:'Special Elite',serif;font-size:1.4rem;color:var(--accent);margin-bottom:1rem;padding-bottom:.5rem;border-bottom:1px solid var(--border);}
  .doc-content :global(h3){font-family:'Rajdhani',monospace;font-size:1rem;font-weight:700;color:var(--text);margin:1.5rem 0 .6rem;letter-spacing:.05em;text-transform:uppercase;}
  .doc-content :global(p){font-size:.9rem;color:var(--text-dim);line-height:1.75;margin-bottom:.85rem;}
  .doc-content :global(ul){padding-left:1.5rem;margin-bottom:.85rem;}
  .doc-content :global(li){font-size:.88rem;color:var(--text-dim);line-height:1.7;margin-bottom:.3rem;}
  .doc-content :global(li strong){color:var(--text);}
  .doc-content :global(code){background:var(--bg3);border:1px solid var(--border);padding:1px 5px;border-radius:3px;font-family:'Courier Prime',monospace;font-size:.82rem;color:var(--accent);}

  .feature-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:.75rem;margin:1rem 0;}
  .feature-card{background:var(--bg2);border:1px solid var(--border);border-radius:6px;padding:.85rem;display:flex;flex-direction:column;gap:.3rem;}
  .fc-icon{font-size:1.1rem;color:var(--accent);}
  .feature-card strong{font-size:.88rem;color:var(--text);}
  .feature-card span{font-size:.78rem;color:var(--text-dim);line-height:1.5;}

  .code-block{background:var(--bg3);border:1px solid var(--border);border-radius:6px;overflow:hidden;margin:.75rem 0;}
  .code-label{font-size:.68rem;letter-spacing:.12em;color:var(--text-dim);padding:.35rem .75rem;background:var(--bg2);border-bottom:1px solid var(--border);}
  .code-block pre{padding:.85rem .9rem;font-family:'Courier Prime',monospace;font-size:.82rem;color:var(--text-dim);line-height:1.6;overflow-x:auto;white-space:pre;}

  .arch-diagram{display:flex;flex-direction:column;align-items:center;gap:.5rem;padding:1.5rem;background:var(--bg2);border:1px solid var(--border);border-radius:8px;margin:1rem 0;}
  .arch-box{background:var(--bg3);border:1px solid var(--accent);border-radius:6px;padding:.65rem 1.25rem;text-align:center;font-weight:700;font-size:.88rem;}
  .arch-box.client{border-color:var(--accent);}
  .arch-box.core{border-color:var(--accent);}
  .arch-box.system{border-color:var(--border);}
  .arch-box span{display:block;font-size:.72rem;color:var(--text-dim);font-weight:400;font-family:'Courier Prime',monospace;margin-top:2px;}
  .arch-arrow{font-size:.85rem;color:var(--text-dim);}

  .pkg-table{background:var(--bg2);border:1px solid var(--border);border-radius:6px;overflow:hidden;margin:.75rem 0;}
  .pkg-row{display:grid;grid-template-columns:240px 1fr;gap:1rem;padding:.55rem .85rem;border-bottom:1px solid var(--border);font-size:.85rem;align-items:baseline;}
  .pkg-row:last-child{border-bottom:none;}
  .pkg-row.header{background:var(--accent-dim);font-size:.75rem;letter-spacing:.1em;}
  .pkg-row code{color:var(--accent);background:transparent;border:none;font-size:.82rem;}
  .pkg-row span{color:var(--text-dim);}

  .faq-item{background:var(--bg2);border:1px solid var(--border);border-radius:6px;margin-bottom:.75rem;overflow:hidden;}
  .faq-q{padding:.75rem 1rem;font-weight:700;font-size:.9rem;color:var(--text);background:var(--accent-dim);border-bottom:1px solid var(--border);}
  .faq-a{padding:.75rem 1rem;font-size:.86rem;color:var(--text-dim);line-height:1.65;}

  @media(max-width:640px){
    .doc-nav{width:120px;min-width:120px;}
    .pkg-row{grid-template-columns:1fr;gap:.25rem;}
    .doc-content{padding:1rem;}
  }
</style>
