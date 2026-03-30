<script lang="ts">
  import { store, activatePro } from '$lib/stores/app.svelte';
  let { feature, children }: { feature: string; children?: any } = $props();
  let showUpgrade = $state(false);
  let licenseKey = $state('');
</script>

{#if store.isPro}
  {@render children?.()}
{:else}
  <div class="pro-gate">
    <div class="pg-card">
      <div class="pg-icon">◆</div>
      <h2>Pro Feature</h2>
      <p><strong>{feature}</strong> is part of the aEye Watcher Pro tier.</p>
      <div class="pg-features">
        <div class="pgf">✓ Change Detection Alerts</div>
        <div class="pgf">✓ Scan Export (JSON / PDF)</div>
        <div class="pgf">✓ Remote Host Scanning</div>
        <div class="pgf">✓ Advanced Scheduler</div>
        <div class="pgf">✓ Priority Support</div>
      </div>
      <button class="pg-btn" onclick={() => showUpgrade = !showUpgrade}>
        {showUpgrade ? '✕ Cancel' : '◉ Activate Pro'}
      </button>
      {#if showUpgrade}
        <div class="pg-form">
          <input class="pg-input" type="text" placeholder="License key (try: AEYE-PRO-DEMO)" bind:value={licenseKey} />
          <button class="pg-activate" onclick={() => activatePro(licenseKey)}>Activate</button>
        </div>
        <p class="pg-hint">Get a license at <a href="https://github.com/NathanielRand/aEyeWatcher" target="_blank">github.com/NathanielRand/aEyeWatcher</a></p>
      {/if}
    </div>
  </div>
{/if}

<style>
  .pro-gate{display:flex;align-items:center;justify-content:center;min-height:400px;padding:2rem;}
  .pg-card{background:var(--bg2);border:1px solid var(--accent);border-radius:12px;padding:2.5rem;max-width:440px;width:100%;text-align:center;box-shadow:0 0 40px var(--accent-dim);}
  .pg-icon{font-size:2rem;color:var(--accent);margin-bottom:.75rem;}
  .pg-card h2{font-family:'Special Elite',serif;font-size:1.4rem;color:var(--text);margin-bottom:.5rem;}
  .pg-card p{font-size:.88rem;color:var(--text-dim);margin-bottom:1.25rem;line-height:1.6;}
  .pg-card p strong{color:var(--accent);}
  .pg-features{display:flex;flex-direction:column;gap:.35rem;margin-bottom:1.5rem;text-align:left;}
  .pgf{font-size:.84rem;color:var(--text-dim);padding:.2rem 0;}
  .pgf:first-child{color:var(--accent);}
  .pg-btn{background:var(--accent);border:none;color:var(--bg);padding:.55rem 1.5rem;border-radius:5px;cursor:pointer;font-family:'Special Elite',serif;font-size:.95rem;letter-spacing:.05em;transition:opacity .15s;margin-bottom:1rem;}
  .pg-btn:hover{opacity:.85;}
  .pg-form{display:flex;gap:.5rem;margin-bottom:.5rem;}
  .pg-input{flex:1;background:var(--bg3);border:1px solid var(--border);color:var(--text);padding:.38rem .65rem;border-radius:4px;font-family:'Courier Prime',monospace;font-size:.82rem;outline:none;}
  .pg-input:focus{border-color:var(--accent);}
  .pg-activate{background:var(--accent-dim);border:1px solid var(--accent);color:var(--accent);padding:.38rem .75rem;border-radius:4px;cursor:pointer;font-family:'Rajdhani',monospace;font-size:.82rem;font-weight:700;white-space:nowrap;}
  .pg-hint{font-size:.72rem;color:var(--text-muted);}
  .pg-hint a{color:var(--accent);}
</style>
