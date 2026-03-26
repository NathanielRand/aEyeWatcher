<script lang="ts">
  let monoTab = $state<'overview' | 'pro' | 'cloud' | 'services'>('overview');
</script>

<div class="page">
  <header class="page-header">
    <div class="page-title">
      <span class="page-icon">◉</span>
      <div>
        <h1>About aEye Watcher</h1>
        <p class="page-sub">Open source project · MIT License · Built for the AI community</p>
      </div>
    </div>
    <div class="header-links">
      <a href="https://github.com/NathanielRand/aEyeWatcher" target="_blank" rel="noopener" class="link-btn">⬡ GitHub</a>
      <a href="https://github.com/NathanielRand/aEyeWatcher/issues" target="_blank" rel="noopener" class="link-btn">◈ Issues</a>
      <a href="https://github.com/NathanielRand/aEyeWatcher/releases" target="_blank" rel="noopener" class="link-btn accent">◆ Releases</a>
    </div>
  </header>

  <div class="content">

    <!-- Hero -->
    <div class="hero-card">
      <div class="hero-left">
        <div class="hero-logo">◉</div>
        <div>
          <div class="hero-name">aEye Watcher</div>
          <div class="hero-tagline">"Every model leaves a trace. Every agent leaves a trail. We watch them all."</div>
          <div class="hero-meta">
            <span class="badge">v2.0.0</span>
            <span class="badge">MIT License</span>
            <span class="badge">Go + SvelteKit</span>
            <span class="badge accent">Open Source</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Project stats row -->
    <div class="stats-row">
      <div class="stat-box"><span class="sv">60+</span><span class="sl">Components Tracked</span></div>
      <div class="stat-box"><span class="sv">12</span><span class="sl">Agent Ecosystems</span></div>
      <div class="stat-box"><span class="sv">13</span><span class="sl">Model Formats</span></div>
      <div class="stat-box"><span class="sv">3</span><span class="sl">Platforms Supported</span></div>
    </div>

    <!-- Two column: about + tech stack -->
    <div class="two-col">
      <div class="info-card">
        <div class="card-header-about">◈ Project</div>
        <div class="info-body">
          <p>aEye Watcher started as a personal tool to answer one question: <em>"What the hell is actually running on my machine?"</em> As the AI tooling landscape exploded in 2024-2025, keeping track of what's installed, what's running, and what's just quietly eating disk space became its own problem.</p>
          <p>The project is fully open-source and community-driven. The core is a single Go binary with zero runtime dependencies. The UI is a SvelteKit SPA that talks to it over HTTP + WebSockets.</p>
          <p>Contributions, bug reports, and ecosystem additions are always welcome.</p>
        </div>
      </div>
      <div class="info-card">
        <div class="card-header-about">⚙ Tech Stack</div>
        <div class="stack-list">
          <div class="stack-row"><span class="stack-lang go">Go 1.22+</span><span>Core scanner + API server</span></div>
          <div class="stack-row"><span class="stack-lang svelte">Svelte 5</span><span>Frontend (runes mode)</span></div>
          <div class="stack-row"><span class="stack-lang ts">TypeScript</span><span>Full type safety across UI</span></div>
          <div class="stack-row"><span class="stack-lang tw">Tailwind CSS</span><span>Utility styling</span></div>
          <div class="stack-row"><span class="stack-lang ws">WebSocket</span><span>Live scan progress streaming</span></div>
          <div class="stack-row"><span class="stack-lang mit">MIT</span><span>License — use it however you want</span></div>
        </div>
      </div>
    </div>

    <!-- Monetization section -->
    <div class="mono-section">
      <div class="section-title-row">
        <span class="section-icon">◆</span>
        <h2>Monetization Strategy</h2>
        <span class="section-sub">Open source + sustainable — the 2026 playbook</span>
      </div>

      <div class="mono-tabs">
        {#each [
          { id: 'overview', label: 'Overview' },
          { id: 'pro',      label: 'Pro / Teams' },
          { id: 'cloud',    label: 'Cloud / SaaS' },
          { id: 'services', label: 'Services' },
        ] as tab}
          <button class="mono-tab" class:active={monoTab === tab.id} onclick={() => monoTab = tab.id}>
            {tab.label}
          </button>
        {/each}
      </div>

      <div class="mono-content">
        {#if monoTab === 'overview'}
          <div class="mono-intro">
            <p>Donations and ad-based models don't scale for developer tools in 2026. The industry has converged on a clear playbook: <strong>open-core</strong> with a paid hosted layer. The open-source core stays 100% free — monetization comes from the surrounding infrastructure and team features.</p>
          </div>
          <div class="strategy-grid">
            <div class="strategy-card primary">
              <div class="sc-rank">01</div>
              <div class="sc-title">Open-Core (Pro Tier)</div>
              <div class="sc-desc">Free core forever. Paid tier adds team dashboards, scheduled scans, change alerts, and export. ~$9-19/mo per user. This is the Tailscale / Grafana model.</div>
              <div class="sc-fit">Best fit for aEye</div>
            </div>
            <div class="strategy-card">
              <div class="sc-rank">02</div>
              <div class="sc-title">Hosted Cloud (SaaS)</div>
              <div class="sc-desc">aEye Cloud — run a lightweight agent on your machines that reports to a hosted dashboard. No local setup. Usage-based pricing. ~$0.01/scan/host.</div>
              <div class="sc-fit">Medium-term target</div>
            </div>
            <div class="strategy-card">
              <div class="sc-rank">03</div>
              <div class="sc-title">Enterprise License</div>
              <div class="sc-desc">Self-hosted deployment with SSO, RBAC, audit logs, SLA support. Annual contract. Targets AI ops teams at mid-large companies managing fleets of GPU nodes.</div>
              <div class="sc-fit">Long-term</div>
            </div>
            <div class="strategy-card">
              <div class="sc-rank">04</div>
              <div class="sc-title">Marketplace / Plugins</div>
              <div class="sc-desc">Paid scanner plugins for proprietary AI platforms (internal corp tools, specialized hardware). Community marketplace with revenue share.</div>
              <div class="sc-fit">Ecosystem play</div>
            </div>
          </div>

        {:else if monoTab === 'pro'}
          <h3>Pro & Teams Tier</h3>
          <p>The open-core model is the most proven path for dev tools. The core stays MIT — you monetize the layer above it.</p>
          <div class="compare-table">
            <div class="ct-row header">
              <span>Feature</span><span>Free (OSS)</span><span>Pro ($15/mo)</span><span>Teams ($49/seat)</span>
            </div>
            <div class="ct-row"><span>Full scan engine</span><span class="yes">✓</span><span class="yes">✓</span><span class="yes">✓</span></div>
            <div class="ct-row"><span>All agent ecosystems</span><span class="yes">✓</span><span class="yes">✓</span><span class="yes">✓</span></div>
            <div class="ct-row"><span>Scheduled scans</span><span class="no">—</span><span class="yes">✓</span><span class="yes">✓</span></div>
            <div class="ct-row"><span>Change detection alerts</span><span class="no">—</span><span class="yes">✓</span><span class="yes">✓</span></div>
            <div class="ct-row"><span>Export (JSON / PDF)</span><span class="no">—</span><span class="yes">✓</span><span class="yes">✓</span></div>
            <div class="ct-row"><span>Remote host scanning</span><span class="no">—</span><span class="yes">✓</span><span class="yes">✓</span></div>
            <div class="ct-row"><span>Multi-host dashboard</span><span class="no">—</span><span class="no">—</span><span class="yes">✓</span></div>
            <div class="ct-row"><span>RBAC + SSO</span><span class="no">—</span><span class="no">—</span><span class="yes">✓</span></div>
            <div class="ct-row"><span>Slack / webhook alerts</span><span class="no">—</span><span class="no">—</span><span class="yes">✓</span></div>
            <div class="ct-row"><span>Priority support</span><span class="no">—</span><span class="no">—</span><span class="yes">✓</span></div>
          </div>

        {:else if monoTab === 'cloud'}
          <h3>aEye Cloud</h3>
          <p>A lightweight agent (~5MB) runs on each machine you want to monitor. It calls home to a hosted aEye dashboard. No port forwarding, no SSH config — just install and forget.</p>
          <div class="cloud-arch">
            <div class="ca-box">Your Machine(s)<br><span>aeye-agent (tiny daemon)</span></div>
            <div class="ca-arrow">→ HTTPS every N min</div>
            <div class="ca-box accent">aEye Cloud<br><span>hosted dashboard</span></div>
          </div>
          <div class="pricing-cards">
            <div class="pricing-card">
              <div class="pc-name">Hobbyist</div>
              <div class="pc-price">Free</div>
              <div class="pc-detail">1 host · 24h history · community support</div>
            </div>
            <div class="pricing-card featured">
              <div class="pc-name">Builder</div>
              <div class="pc-price">$12<span>/mo</span></div>
              <div class="pc-detail">5 hosts · 30d history · alerts · exports</div>
            </div>
            <div class="pricing-card">
              <div class="pc-name">Studio</div>
              <div class="pc-price">$49<span>/mo</span></div>
              <div class="pc-detail">25 hosts · 90d history · RBAC · API access</div>
            </div>
          </div>

        {:else if monoTab === 'services'}
          <h3>Professional Services</h3>
          <p>Beyond the product itself, there are natural adjacent revenue streams for an AI ecosystem monitoring tool:</p>
          <div class="service-list">
            <div class="service-item">
              <div class="si-title">AI Infra Audit</div>
              <div class="si-desc">One-time engagement: run aEye across a team's infrastructure, produce a report of what's running, what's orphaned, what's misconfigured, and what's a security risk. $2-5K per engagement.</div>
            </div>
            <div class="service-item">
              <div class="si-title">Custom Scanner Plugins</div>
              <div class="si-desc">Build bespoke detection for proprietary internal AI platforms (e.g. a company's internal LLM deployment, custom agent framework). Hourly or project-based.</div>
            </div>
            <div class="service-item">
              <div class="si-title">Integration Consulting</div>
              <div class="si-desc">Help teams wire aEye into their CI/CD pipelines, monitoring stacks (Datadog, Grafana), and incident response workflows.</div>
            </div>
            <div class="service-item">
              <div class="si-title">Sponsored Ecosystem Support</div>
              <div class="si-desc">AI companies pay to ensure their agent/model is a first-class citizen in the scanner — accurate detection, rich agent cards, docs. Similar to OpenCollective's sponsor model but product-integrated.</div>
            </div>
          </div>
        {/if}
      </div>
    </div>

  </div>
</div>

<style>
  .page{display:flex;flex-direction:column;height:100%;overflow:hidden;}
  .page-header{display:flex;align-items:center;justify-content:space-between;flex-wrap:wrap;gap:.75rem;padding:.85rem 1.5rem;background:var(--bg2);border-bottom:1px solid var(--border);flex-shrink:0;}
  .page-title{display:flex;align-items:center;gap:.85rem;}
  .page-icon{font-size:1.4rem;color:var(--accent);}
  .page-title h1{font-family:'Special Elite',serif;font-size:1.15rem;}
  .page-sub{font-size:.8rem;color:var(--text-dim);margin-top:2px;}
  .header-links{display:flex;gap:.5rem;}
  .link-btn{background:var(--bg3);border:1px solid var(--border);color:var(--text-dim);padding:.35rem .75rem;border-radius:3px;font-size:.8rem;transition:all .15s;}
  .link-btn:hover{border-color:var(--accent);color:var(--accent);}
  .link-btn.accent{background:var(--accent-dim);border-color:var(--accent);color:var(--accent);}

  .content{flex:1;overflow-y:auto;padding:1.5rem;scrollbar-width:thin;scrollbar-color:var(--border) transparent;display:flex;flex-direction:column;gap:1rem;}

  /* Hero */
  .hero-card{background:var(--bg2);border:1px solid var(--border);border-radius:8px;padding:1.5rem;position:relative;overflow:hidden;}
  .hero-card::before{content:'◉';position:absolute;right:1.5rem;top:50%;transform:translateY(-50%);font-size:6rem;color:var(--accent);opacity:.06;pointer-events:none;}
  .hero-left{display:flex;align-items:center;gap:1.25rem;}
  .hero-logo{font-size:2.5rem;color:var(--accent);animation:pulse-eye 3s ease-in-out infinite;}
  @keyframes pulse-eye{0%,100%{opacity:1}50%{opacity:.6}}
  .hero-name{font-family:'Special Elite',serif;font-size:1.5rem;color:var(--text);margin-bottom:.25rem;}
  .hero-tagline{font-size:.85rem;color:var(--text-dim);font-style:italic;margin-bottom:.65rem;}
  .hero-meta{display:flex;flex-wrap:wrap;gap:.4rem;}
  .badge{font-size:.65rem;padding:2px 7px;border-radius:2px;background:var(--bg3);border:1px solid var(--border);color:var(--text-dim);letter-spacing:.08em;}
  .badge.accent{background:var(--accent-dim);border-color:var(--accent);color:var(--accent);}

  /* Stats */
  .stats-row{display:grid;grid-template-columns:repeat(4,1fr);gap:.75rem;}
  .stat-box{background:var(--bg2);border:1px solid var(--border);border-radius:6px;padding:.85rem;text-align:center;}
  .sv{display:block;font-family:'Special Elite',serif;font-size:1.6rem;color:var(--accent);}
  .sl{font-size:.72rem;color:var(--text-dim);letter-spacing:.08em;text-transform:uppercase;}

  /* Two col */
  .two-col{display:grid;grid-template-columns:1fr 1fr;gap:1rem;}
  .info-card{background:var(--bg2);border:1px solid var(--border);border-radius:8px;overflow:hidden;}
  .card-header-about{padding:.65rem 1rem;background:var(--accent-dim);border-bottom:1px solid var(--border);font-size:.75rem;letter-spacing:.15em;color:var(--accent);font-weight:700;font-family:'Rajdhani',monospace;}
  .info-body{padding:.85rem 1rem;}
  .info-body p{font-size:.86rem;color:var(--text-dim);line-height:1.7;margin-bottom:.65rem;}
  .info-body p:last-child{margin-bottom:0;}
  .info-body em{color:var(--text);}
  .stack-list{padding:.25rem 0;}
  .stack-row{display:flex;align-items:center;gap:.75rem;padding:.5rem 1rem;border-bottom:1px solid var(--border);font-size:.85rem;}
  .stack-row:last-child{border-bottom:none;}
  .stack-row span:last-child{color:var(--text-dim);}
  .stack-lang{font-size:.65rem;padding:2px 7px;border-radius:2px;font-weight:700;letter-spacing:.08em;flex-shrink:0;}
  .stack-lang.go{background:rgba(0,173,216,.15);border:1px solid rgba(0,173,216,.3);color:#00add8;}
  .stack-lang.svelte{background:rgba(255,62,0,.12);border:1px solid rgba(255,62,0,.3);color:#ff3e00;}
  .stack-lang.ts{background:rgba(49,120,198,.12);border:1px solid rgba(49,120,198,.3);color:#3178c6;}
  .stack-lang.tw{background:rgba(56,189,248,.12);border:1px solid rgba(56,189,248,.3);color:#38bdf8;}
  .stack-lang.ws{background:var(--accent-dim);border:1px solid var(--accent);color:var(--accent);}
  .stack-lang.mit{background:rgba(74,222,128,.12);border:1px solid rgba(74,222,128,.3);color:#4ade80;}

  /* Monetization */
  .mono-section{background:var(--bg2);border:1px solid var(--border);border-radius:8px;overflow:hidden;}
  .section-title-row{display:flex;align-items:center;gap:.75rem;padding:.85rem 1.25rem;border-bottom:1px solid var(--border);background:var(--accent-dim);}
  .section-icon{font-size:1rem;color:var(--accent);}
  .section-title-row h2{font-family:'Special Elite',serif;font-size:1.1rem;color:var(--accent);flex:1;}
  .section-sub{font-size:.75rem;color:var(--text-dim);}
  .mono-tabs{display:flex;gap:0;border-bottom:1px solid var(--border);}
  .mono-tab{background:none;border:none;color:var(--text-dim);padding:.6rem 1rem;cursor:pointer;font-family:'Rajdhani',monospace;font-size:.85rem;font-weight:500;border-bottom:2px solid transparent;transition:all .15s;}
  .mono-tab:hover{color:var(--text);}
  .mono-tab.active{color:var(--accent);border-bottom-color:var(--accent);}
  .mono-content{padding:1.25rem;}
  .mono-intro p{font-size:.88rem;color:var(--text-dim);line-height:1.75;margin-bottom:.85rem;}
  .mono-content :global(h3){font-family:'Rajdhani',monospace;font-size:1rem;font-weight:700;letter-spacing:.05em;text-transform:uppercase;color:var(--text);margin-bottom:.85rem;}
  .mono-content :global(p){font-size:.88rem;color:var(--text-dim);line-height:1.7;margin-bottom:.75rem;}

  .strategy-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(220px,1fr));gap:.75rem;}
  .strategy-card{background:var(--bg3);border:1px solid var(--border);border-radius:6px;padding:1rem;}
  .strategy-card.primary{border-color:var(--accent);}
  .sc-rank{font-family:'Courier Prime',monospace;font-size:.65rem;color:var(--text-muted);margin-bottom:.4rem;letter-spacing:.15em;}
  .sc-title{font-weight:700;font-size:.9rem;color:var(--text);margin-bottom:.4rem;}
  .sc-desc{font-size:.8rem;color:var(--text-dim);line-height:1.6;margin-bottom:.5rem;}
  .sc-fit{font-size:.68rem;color:var(--accent);letter-spacing:.1em;font-weight:600;}

  .compare-table{background:var(--bg3);border:1px solid var(--border);border-radius:6px;overflow:hidden;}
  .ct-row{display:grid;grid-template-columns:1fr 80px 80px 90px;gap:.5rem;padding:.5rem .85rem;border-bottom:1px solid var(--border);font-size:.82rem;align-items:center;}
  .ct-row:last-child{border-bottom:none;}
  .ct-row.header{background:var(--accent-dim);font-size:.72rem;letter-spacing:.1em;font-weight:700;color:var(--text);}
  .ct-row span:first-child{color:var(--text-dim);}
  .yes{color:#4ade80;text-align:center;font-weight:700;}
  .no{color:var(--text-muted);text-align:center;}

  .cloud-arch{display:flex;align-items:center;gap:1rem;padding:1rem;background:var(--bg3);border:1px solid var(--border);border-radius:6px;margin-bottom:1rem;}
  .ca-box{background:var(--bg2);border:1px solid var(--border);border-radius:5px;padding:.6rem 1rem;text-align:center;font-weight:700;font-size:.85rem;}
  .ca-box.accent{border-color:var(--accent);color:var(--accent);}
  .ca-box span{display:block;font-size:.72rem;color:var(--text-dim);font-weight:400;font-family:'Courier Prime',monospace;margin-top:2px;}
  .ca-arrow{font-size:.8rem;color:var(--text-dim);}
  .pricing-cards{display:grid;grid-template-columns:repeat(3,1fr);gap:.75rem;}
  .pricing-card{background:var(--bg3);border:1px solid var(--border);border-radius:6px;padding:1rem;text-align:center;}
  .pricing-card.featured{border-color:var(--accent);}
  .pc-name{font-size:.75rem;letter-spacing:.15em;color:var(--text-dim);text-transform:uppercase;margin-bottom:.4rem;}
  .pc-price{font-family:'Special Elite',serif;font-size:1.8rem;color:var(--accent);line-height:1;}
  .pc-price span{font-size:.85rem;color:var(--text-dim);}
  .pc-detail{font-size:.75rem;color:var(--text-dim);margin-top:.4rem;line-height:1.5;}

  .service-list{display:flex;flex-direction:column;gap:.65rem;}
  .service-item{background:var(--bg3);border:1px solid var(--border);border-radius:6px;padding:.85rem;}
  .si-title{font-weight:700;font-size:.9rem;color:var(--text);margin-bottom:.35rem;}
  .si-desc{font-size:.83rem;color:var(--text-dim);line-height:1.65;}

  @media(max-width:768px){
    .stats-row{grid-template-columns:repeat(2,1fr);}
    .two-col{grid-template-columns:1fr;}
    .pricing-cards{grid-template-columns:1fr;}
    .ct-row{grid-template-columns:1fr 60px 60px;}
    .ct-row .ct-row:nth-child(4){display:none;}
  }
  @media(max-width:480px){
    .stats-row{grid-template-columns:repeat(2,1fr);}
    .strategy-grid{grid-template-columns:1fr;}
  }
</style>
