package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Campfire</title>
<style>:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5;height:100vh;display:flex;flex-direction:column}
.hdr{padding:.8rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;flex-shrink:0}.hdr h1{font-size:.9rem;letter-spacing:2px}
.layout{display:grid;grid-template-columns:200px 1fr;flex:1;overflow:hidden}
.sidebar{border-right:1px solid var(--bg3);padding:.8rem;overflow-y:auto}
.sidebar-title{font-size:.6rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}
.ch{padding:.4rem .6rem;cursor:pointer;font-size:.75rem;color:var(--cd);border-left:2px solid transparent;margin-bottom:.1rem}
.ch:hover{color:var(--cream);background:var(--bg2)}.ch.active{color:var(--rust);border-color:var(--rust);background:var(--bg2)}
.ch-count{font-size:.55rem;color:var(--cm);float:right}
.chat{display:flex;flex-direction:column;overflow:hidden}
.messages{flex:1;overflow-y:auto;padding:1rem;display:flex;flex-direction:column-reverse}
.msg{margin-bottom:.8rem}
.msg-author{font-size:.7rem;color:var(--gold);margin-bottom:.1rem}
.msg-body{font-size:.82rem;color:var(--cd);line-height:1.6}
.msg-time{font-size:.55rem;color:var(--cm)}
.compose{border-top:1px solid var(--bg3);padding:.8rem;display:flex;gap:.5rem}
.compose input{flex:1;padding:.5rem .8rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.78rem}
.compose input:focus{border-color:var(--leather);outline:none}
.btn{font-size:.6rem;padding:.3rem .8rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg)}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem;flex:1;display:flex;align-items:center;justify-content:center}
@media(max-width:600px){.layout{grid-template-columns:1fr}.sidebar{display:none}}
</style></head><body>
<div class="hdr"><h1>CAMPFIRE</h1><button class="btn" onclick="newChannel()">+ Channel</button></div>
<div class="layout">
<div class="sidebar"><div class="sidebar-title">Channels</div><div id="channels"></div></div>
<div class="chat"><div class="messages" id="messages"><div class="empty">Select a channel to start chatting</div></div>
<div class="compose"><input id="author" placeholder="name" style="width:80px;flex:none"><input id="msg" placeholder="Type a message..." onkeydown="if(event.key==='Enter')send()"><button class="btn btn-p" onclick="send()">Send</button></div></div>
</div>
<script>
const A='/api';let channels=[],curCh='',msgs=[];
async function load(){const r=await fetch(A+'/channels').then(r=>r.json());channels=r.channels||[];renderChannels();if(curCh)loadMessages();}
function renderChannels(){let h='';channels.forEach(c=>{h+='<div class="ch'+(curCh===c.id?' active':'')+'" onclick="selectCh(\''+c.id+'\')"><span>#'+esc(c.name)+'</span><span class="ch-count">'+c.message_count+'</span></div>';});document.getElementById('channels').innerHTML=h;}
function selectCh(id){curCh=id;renderChannels();loadMessages();}
async function loadMessages(){const r=await fetch(A+'/channels/'+curCh+'/messages').then(r=>r.json());msgs=r.messages||[];renderMessages();}
function renderMessages(){if(!msgs.length){document.getElementById('messages').innerHTML='<div class="empty">No messages yet. Say something.</div>';return;}
let h='';msgs.forEach(m=>{h+='<div class="msg"><div class="msg-author">'+esc(m.author)+' <span class="msg-time">'+ft(m.created_at)+'</span></div><div class="msg-body">'+esc(m.body)+'</div></div>';});
document.getElementById('messages').innerHTML=h;}
async function send(){const body=document.getElementById('msg').value;if(!body||!curCh)return;const author=document.getElementById('author').value||'anonymous';await fetch(A+'/channels/'+curCh+'/messages',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({author,body})});document.getElementById('msg').value='';loadMessages();load();}
function newChannel(){const name=prompt('Channel name:');if(!name)return;fetch(A+'/channels',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name})}).then(()=>load());}
function ft(t){if(!t)return'';const d=new Date(t);return d.toLocaleTimeString([],{hour:'2-digit',minute:'2-digit'});}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
load();setInterval(()=>{if(curCh)loadMessages()},5000);
</script><script>
(function(){
  fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
    if(!cfg||typeof cfg!=='object')return;
    if(cfg.dashboard_title){
      document.title=cfg.dashboard_title;
      var h1=document.querySelector('h1');
      if(h1){
        var inner=h1.innerHTML;
        var firstSpan=inner.match(/<span[^>]*>[^<]*<\/span>/);
        if(firstSpan){h1.innerHTML=firstSpan[0]+' '+cfg.dashboard_title}
        else{h1.textContent=cfg.dashboard_title}
      }
    }
  }).catch(function(){});
})();
</script>
</body></html>`
