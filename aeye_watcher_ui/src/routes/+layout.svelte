<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import { store, setTheme, toggleSidebar, THEMES, pollGPU, fetchOSInfo } from '$lib/stores/app.svelte';
  import '../app.css';

  let { children } = $props();
  let gpuInterval: ReturnType<typeof setInterval>;

  onMount(async () => {
    const saved = localStorage.getItem('aeye-theme');
    if (saved && saved in THEMES) setTheme(saved as keyof typeof THEMES);
    else setTheme('obsidian');
    await fetchOSInfo();
    gpuInterval = setInterval(pollGPU, 3000);
    pollGPU();
    return () => clearInterval(gpuInterval);
  });

  const navItems = [
    { href: '/',       label: 'Command Center', icon: '⬡' },
    { href: '/config', label: 'Configuration',  icon: '⚙' },
    { href: '/docs',   label: 'Documentation',  icon: '◈' },
    { href: '/about',  label: 'About / OSS',    icon: '◉' },
  ];

  function isActive(href: string) {
    if (href === '/') return $page.url.pathname === '/';
    return $page.url.pathname.startsWith(href);
  }
</script>

<div class="shell" data-theme={store.currentTheme}>
  <div class="grain"></div>

  <aside class="sidebar" class:closed={!store.sidebarOpen}>
    <div class="sidebar-logo">
      <a href="/" class="logo-mark">
        <span class="logo-eye">◉</span>
        {#if store.sidebarOpen}
          <div>
            <div class="logo-name">aEye</div>
            <div class="logo-sub">WATCHER</div>
          </div>
        {/if}
      </a>
      <button class="toggle-btn" onclick={toggleSidebar}>
        {store.sidebarOpen ? '◀' : '▶'}
      </button>
    </div>

    {#if store.osInfo && store.sidebarOpen}
      <div class="os-badge">
        <span class="os-icon">
          {store.osInfo.os === 'linux' ? '🐧' : store.osInfo.os === 'darwin' ? '' : '🪟'}
        </span>
        <div class="os-info">
          <span class="os-name">{store.osInfo.distro ?? store.osInfo.os}</span>
          <span class="os-arch">{store.osInfo.arch} · {store.osInfo.hostname}</span>
        </div>
        {#if store.osInfo.gpu_available}<span class="gpu-badge">GPU</span>{/if}
      </div>
    {/if}

    <nav class="nav">
      {#each navItems as item}
        <a href={item.href} class="nav-item" class:active={isActive(item.href)}>
          <span class="nav-icon">{item.icon}</span>
          {#if store.sidebarOpen}
            <span class="nav-label">{item.label}</span>
          {/if}
        </a>
      {/each}
    </nav>

    {#if store.sidebarOpen}
      <div class="sidebar-footer">
        <div class="theme-label">THEME</div>
        <div class="theme-pills">
          {#each Object.entries(THEMES) as [key, val]}
            <button
              class="theme-pill"
              class:active={store.currentTheme === key}
              style="--pill-color: {val.accent}"
              onclick={() => setTheme(key as keyof typeof THEMES)}
              title={val.label}
            ></button>
          {/each}
        </div>
        <a class="gh-link" href="https://github.com/NathanielRand/aEyeWatcher" target="_blank" rel="noopener">
          <span>⬡</span> GitHub
        </a>
      </div>
    {/if}
  </aside>

  <div class="page-wrap">
    {@render children()}
  </div>
</div>

<style>
  @import url('https://fonts.googleapis.com/css2?family=Special+Elite&family=Courier+Prime:wght@400;700&family=Rajdhani:wght@400;500;600;700&display=swap');

  :global([data-theme="obsidian"]){--accent:#c084fc;--accent-dim:rgba(192,132,252,0.15);--bg:#08080f;--bg2:#0f0f1a;--bg3:#16162a;--border:rgba(192,132,252,0.14);--text:#e8e6f8;--text-dim:#8a8aa8;--text-muted:#45445c;}
  :global([data-theme="clay"]){--accent:#fb923c;--accent-dim:rgba(251,146,60,0.13);--bg:#100a06;--bg2:#1a1008;--bg3:#241608;--border:rgba(251,146,60,0.12);--text:#f2ebe2;--text-dim:#9a8070;--text-muted:#483220;}
  :global([data-theme="forest"]){--accent:#4ade80;--accent-dim:rgba(74,222,128,0.12);--bg:#060e08;--bg2:#081a0e;--bg3:#0c2414;--border:rgba(74,222,128,0.12);--text:#e2f0e5;--text-dim:#6a9870;--text-muted:#22402a;}
  :global([data-theme="slate"]){--accent:#7dd3fc;--accent-dim:rgba(125,211,252,0.12);--bg:#060810;--bg2:#0a1020;--bg3:#0e1830;--border:rgba(125,211,252,0.12);--text:#e2ecf8;--text-dim:#6a7a98;--text-muted:#202840;}
  :global([data-theme="amber"]){--accent:#fbbf24;--accent-dim:rgba(251,191,36,0.12);--bg:#0a0800;--bg2:#150f00;--bg3:#1f1600;--border:rgba(251,191,36,0.12);--text:#f8f0d0;--text-dim:#987a30;--text-muted:#3d2e00;}

  :global(*,*::before,*::after){box-sizing:border-box;margin:0;padding:0;}
  :global(body){background:var(--bg);color:var(--text);font-family:'Rajdhani','Courier Prime',monospace;font-size:15px;overflow:hidden;height:100vh;}
  :global(a){color:inherit;text-decoration:none;}
  :global(h1,h2,h3){font-family:'Special Elite',serif;color:var(--text);}
  :global(code,pre){font-family:'Courier Prime',monospace;}

  .grain{position:fixed;inset:0;pointer-events:none;z-index:9999;background-image:url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='0.04'/%3E%3C/svg%3E");background-size:200px 200px;opacity:.55;mix-blend-mode:overlay;}
  .shell{display:flex;height:100vh;overflow:hidden;position:relative;}
  .sidebar{width:230px;min-width:230px;background:var(--bg2);border-right:1px solid var(--border);display:flex;flex-direction:column;transition:width .2s,min-width .2s;overflow:hidden;z-index:100;flex-shrink:0;}
  .sidebar.closed{width:56px;min-width:56px;}
  .sidebar-logo{display:flex;align-items:center;justify-content:space-between;padding:1rem .85rem;border-bottom:1px solid var(--border);}
  .logo-mark{display:flex;align-items:center;gap:.65rem;}
  .logo-eye{font-size:1.5rem;color:var(--accent);animation:pulse-eye 3s ease-in-out infinite;}
  @keyframes pulse-eye{0%,100%{opacity:1;transform:scale(1)}50%{opacity:.7;transform:scale(.94)}}
  .logo-name{font-family:'Special Elite',serif;font-size:1.15rem;color:var(--accent);letter-spacing:.1em;line-height:1;}
  .logo-sub{font-size:.6rem;letter-spacing:.38em;color:var(--text-dim);}
  .toggle-btn{background:none;border:none;color:var(--text-dim);cursor:pointer;font-size:.75rem;padding:2px 5px;transition:color .2s;}
  .toggle-btn:hover{color:var(--accent);}
  .os-badge{display:flex;align-items:center;gap:.5rem;padding:.65rem .85rem;background:var(--accent-dim);border-bottom:1px solid var(--border);}
  .os-icon{font-size:1rem;}
  .os-info{display:flex;flex-direction:column;flex:1;}
  .os-name{color:var(--accent);font-weight:600;font-size:.82rem;}
  .os-arch{color:var(--text-dim);font-size:.72rem;}
  .gpu-badge{background:var(--accent);color:var(--bg);font-size:.62rem;font-weight:700;padding:1px 5px;border-radius:2px;letter-spacing:.1em;}
  .nav{flex:1;display:flex;flex-direction:column;gap:3px;padding:.6rem;}
  .nav-item{display:flex;align-items:center;gap:.65rem;padding:.6rem .65rem;background:none;border:none;border-radius:5px;color:var(--text-dim);cursor:pointer;font-family:'Rajdhani',monospace;font-size:.92rem;font-weight:500;letter-spacing:.05em;transition:all .15s;width:100%;}
  .nav-item:hover{background:var(--accent-dim);color:var(--text);}
  .nav-item.active{background:var(--accent-dim);color:var(--accent);border-left:2px solid var(--accent);}
  .nav-icon{font-size:1rem;width:18px;text-align:center;flex-shrink:0;}
  .nav-label{flex:1;}
  .sidebar-footer{padding:.85rem;border-top:1px solid var(--border);display:flex;flex-direction:column;gap:.65rem;}
  .theme-label{font-size:.62rem;letter-spacing:.2em;color:var(--text-muted);}
  .theme-pills{display:flex;gap:7px;}
  .theme-pill{width:17px;height:17px;border-radius:50%;border:2px solid transparent;background:var(--pill-color);cursor:pointer;transition:transform .15s,border-color .15s;}
  .theme-pill:hover{transform:scale(1.2);}
  .theme-pill.active{border-color:var(--text);transform:scale(1.15);}
  .gh-link{display:flex;align-items:center;gap:.4rem;font-size:.78rem;color:var(--text-dim);letter-spacing:.05em;padding:.35rem .5rem;border-radius:4px;border:1px solid var(--border);transition:all .15s;}
  .gh-link:hover{color:var(--accent);border-color:var(--accent);}
  .page-wrap{flex:1;display:flex;flex-direction:column;overflow:hidden;min-width:0;}

  @media(max-width:768px){
    .sidebar{width:56px;min-width:56px;}
    .sidebar.closed{width:0;min-width:0;border:none;}
    .logo-name,.logo-sub,.nav-label,.os-badge,.sidebar-footer{display:none;}
  }
</style>
