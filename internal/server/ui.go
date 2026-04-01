package server

import "net/http"

func (s *Server) handleUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Campfire — Stockyard</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;600&family=Libre+Baskerville:wght@400;700&display=swap" rel="stylesheet">
<style>*{margin:0;padding:0;box-sizing:border-box}body{background:#1a1410;color:#f0e6d3;font-family:'Libre Baskerville',serif;padding:2rem}
.hdr{font-family:'JetBrains Mono',monospace;font-size:.7rem;color:#a0845c;letter-spacing:3px;text-transform:uppercase;margin-bottom:2rem;border-bottom:2px solid #8b3d1a;padding-bottom:.8rem}
.cards{display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;margin-bottom:2rem;font-family:'JetBrains Mono',monospace}.card{background:#241e18;border:1px solid #2e261e;padding:1rem}.card-val{font-size:1.6rem;font-weight:700;display:block}.card-lbl{font-size:.55rem;letter-spacing:2px;text-transform:uppercase;color:#a0845c;margin-top:.2rem}
.section{margin-bottom:2rem}.section h2{font-family:'JetBrains Mono',monospace;font-size:.65rem;letter-spacing:3px;text-transform:uppercase;color:#e8753a;margin-bottom:.8rem;border-bottom:1px solid #2e261e;padding-bottom:.4rem}
.cat{background:#241e18;padding:.8rem 1rem;margin-bottom:.5rem;border:1px solid #2e261e}.cat-name{font-family:'JetBrains Mono',monospace;font-size:.85rem;color:#f0e6d3;font-weight:600}.cat-desc{font-size:.78rem;color:#7a7060;margin-top:.2rem}.cat-count{font-family:'JetBrains Mono',monospace;font-size:.6rem;color:#a0845c;margin-top:.3rem}
.empty{color:#7a7060;text-align:center;padding:2rem;font-style:italic}
</style></head><body>
<div class="hdr">Stockyard · Campfire</div>
<div class="cards"><div class="card"><span class="card-val" id="s-cats">—</span><span class="card-lbl">Categories</span></div><div class="card"><span class="card-val" id="s-threads">—</span><span class="card-lbl">Threads</span></div><div class="card"><span class="card-val" id="s-replies">—</span><span class="card-lbl">Replies</span></div></div>
<div class="section"><h2>Categories</h2><div id="cat-list"></div></div>
<script>
async function refresh(){
  try{const s=await(await fetch('/api/status')).json();document.getElementById('s-cats').textContent=s.categories||0;document.getElementById('s-threads').textContent=s.threads||0;document.getElementById('s-replies').textContent=s.replies||0;}catch(e){}
  try{const d=await(await fetch('/api/categories')).json();const cs=d.categories||[];
  document.getElementById('cat-list').innerHTML=cs.length?cs.map(c=>'<div class="cat"><div class="cat-name">'+esc(c.name)+'</div><div class="cat-desc">'+esc(c.description)+'</div><div class="cat-count">'+c.thread_count+' threads</div></div>').join(''):'<div class="empty">No categories yet</div>';}catch(e){}
}
function esc(s){const d=document.createElement('div');d.textContent=s||'';return d.innerHTML;}
refresh();setInterval(refresh,8000);
</script></body></html>`))
}
