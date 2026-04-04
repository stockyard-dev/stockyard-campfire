package server
import ("encoding/json";"net/http";"github.com/stockyard-dev/stockyard-campfire/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux;limits Limits}
func New(db *store.DB,limits Limits)*Server{s:=&Server{db:db,mux:http.NewServeMux(),limits:limits}
s.mux.HandleFunc("GET /api/channels",s.listChannels)
s.mux.HandleFunc("POST /api/channels",s.createChannel)
s.mux.HandleFunc("DELETE /api/channels/{id}",s.deleteChannel)
s.mux.HandleFunc("GET /api/channels/{id}/messages",s.listMessages)
s.mux.HandleFunc("POST /api/channels/{id}/messages",s.postMessage)
s.mux.HandleFunc("DELETE /api/messages/{id}",s.deleteMessage)
s.mux.HandleFunc("GET /api/stats",s.stats)
s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /api/tier",func(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"tier":s.limits.Tier,"upgrade_url":"https://stockyard.dev/campfire/"})})
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard);s.mux.HandleFunc("GET /",s.root)
return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)listChannels(w http.ResponseWriter,r *http.Request){ch:=s.db.ListChannels();if ch==nil{ch=[]store.Channel{}};wj(w,200,map[string]any{"channels":ch})}
func(s *Server)createChannel(w http.ResponseWriter,r *http.Request){var c store.Channel;json.NewDecoder(r.Body).Decode(&c);if c.Name==""{we(w,400,"name required");return};s.db.CreateChannel(&c);wj(w,201,c)}
func(s *Server)deleteChannel(w http.ResponseWriter,r *http.Request){s.db.DeleteChannel(r.PathValue("id"));wj(w,200,map[string]string{"status":"deleted"})}
func(s *Server)listMessages(w http.ResponseWriter,r *http.Request){msgs:=s.db.ListMessages(r.PathValue("id"),100);if msgs==nil{msgs=[]store.Message{}};wj(w,200,map[string]any{"messages":msgs})}
func(s *Server)postMessage(w http.ResponseWriter,r *http.Request){var m store.Message;json.NewDecoder(r.Body).Decode(&m);m.ChannelID=r.PathValue("id");if m.Body==""{we(w,400,"body required");return};s.db.PostMessage(&m);wj(w,201,m)}
func(s *Server)deleteMessage(w http.ResponseWriter,r *http.Request){s.db.DeleteMessage(r.PathValue("id"));wj(w,200,map[string]string{"status":"deleted"})}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,s.db.Stats())}
func(s *Server)health(w http.ResponseWriter,r *http.Request){st:=s.db.Stats();wj(w,200,map[string]any{"service":"campfire","status":"ok","channels":st["channels"],"messages":st["messages"]})}
