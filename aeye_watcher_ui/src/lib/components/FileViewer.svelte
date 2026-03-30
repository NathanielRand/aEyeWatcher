<script lang="ts">
  let { file, onclose, onopen }: {
    file: { path: string; content: string; name: string };
    onclose: () => void;
    onopen: () => void;
  } = $props();

  const ext = file.name.split('.').pop()?.toLowerCase() ?? '';
  const isJson = ext === 'json';
  let formatted = $state(file.content);
  if (isJson) {
    try { formatted = JSON.stringify(JSON.parse(file.content), null, 2); } catch { /* raw */ }
  }
</script>

<div class="modal-bg" role="presentation" onclick={onclose}>
  <div class="viewer" onclick={(e)=>e.stopPropagation()}>
    <div class="v-header">
      <div class="v-name-wrap">
        <span class="v-ext">{ext}</span>
        <span class="v-name">{file.name}</span>
      </div>
      <div class="v-actions">
        <button class="v-btn" onclick={onopen}>Open in system ↗</button>
        <button class="v-close" onclick={onclose}>✕</button>
      </div>
    </div>
    <div class="v-path">{file.path}</div>
    <pre class="v-content">{formatted}</pre>
  </div>
</div>

<style>
  .modal-bg{position:fixed;inset:0;background:rgba(0,0,0,.78);display:flex;align-items:center;justify-content:center;z-index:2000;backdrop-filter:blur(6px);}
  .viewer{background:var(--bg2);border:1px solid var(--accent);border-radius:8px;width:min(800px,95vw);max-height:85vh;display:flex;flex-direction:column;box-shadow:0 0 40px rgba(0,0,0,.6);}
  .v-header{display:flex;align-items:center;justify-content:space-between;padding:.7rem 1rem;border-bottom:1px solid var(--border);flex-shrink:0;}
  .v-name-wrap{display:flex;align-items:center;gap:.5rem;}
  .v-ext{font-size:.6rem;padding:2px 6px;background:var(--accent-dim);border:1px solid var(--accent);border-radius:2px;color:var(--accent);letter-spacing:.1em;text-transform:uppercase;}
  .v-name{font-weight:600;font-size:.9rem;color:var(--text);}
  .v-actions{display:flex;gap:.5rem;align-items:center;}
  .v-btn{background:none;border:1px solid var(--border);color:var(--text-dim);padding:.25rem .6rem;border-radius:3px;cursor:pointer;font-size:.75rem;font-family:'Rajdhani',monospace;transition:all .15s;}
  .v-btn:hover{border-color:var(--accent);color:var(--accent);}
  .v-close{background:none;border:none;color:var(--text-muted);cursor:pointer;font-size:.85rem;padding:2px 5px;transition:color .15s;}
  .v-close:hover{color:var(--text);}
  .v-path{padding:.35rem 1rem;font-size:.65rem;color:var(--text-muted);font-family:'Courier Prime',monospace;background:var(--bg3);border-bottom:1px solid var(--border);flex-shrink:0;}
  .v-content{flex:1;overflow:auto;padding:1rem;font-family:'Courier Prime',monospace;font-size:.8rem;color:var(--text-dim);line-height:1.65;white-space:pre-wrap;word-break:break-word;margin:0;}
</style>
