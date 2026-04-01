package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/stockyard-dev/stockyard-campfire/internal/store"
)

type Server struct {
	db     *store.DB
	mux    *http.ServeMux
	port   int
	limits Limits
}

func New(db *store.DB, port int, limits Limits) *Server {
	s := &Server{db: db, mux: http.NewServeMux(), port: port, limits: limits}
	s.mux.HandleFunc("POST /api/categories", s.hCreateCat)
	s.mux.HandleFunc("GET /api/categories", s.hListCats)
	s.mux.HandleFunc("DELETE /api/categories/{id}", s.hDelCat)

	s.mux.HandleFunc("POST /api/categories/{id}/threads", s.hCreateThread)
	s.mux.HandleFunc("GET /api/categories/{id}/threads", s.hListThreads)
	s.mux.HandleFunc("GET /api/threads/{id}", s.hGetThread)
	s.mux.HandleFunc("DELETE /api/threads/{id}", s.hDelThread)

	s.mux.HandleFunc("POST /api/threads/{id}/replies", s.hCreateReply)
	s.mux.HandleFunc("GET /api/threads/{id}/replies", s.hListReplies)

	s.mux.HandleFunc("GET /api/status", func(w http.ResponseWriter, r *http.Request) { wj(w, 200, s.db.Stats()) })
	s.mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) { wj(w, 200, map[string]string{"status": "ok"}) })
	s.mux.HandleFunc("GET /ui", s.handleUI)
	s.mux.HandleFunc("GET /api/version", func(w http.ResponseWriter, r *http.Request) {
		wj(w, 200, map[string]any{"product": "stockyard-campfire", "version": "0.1.0"})
	})
	return s
}

func (s *Server) Start() error {
	log.Printf("[campfire] :%d", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
}

func (s *Server) hCreateCat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.Name == "" || req.Slug == "" {
		wj(w, 400, map[string]string{"error": "name and slug required"})
		return
	}
	c, err := s.db.CreateCategory(req.Name, req.Slug, req.Description)
	if err != nil {
		wj(w, 500, map[string]string{"error": err.Error()})
		return
	}
	wj(w, 201, map[string]any{"category": c})
}

func (s *Server) hListCats(w http.ResponseWriter, r *http.Request) {
	cs, _ := s.db.ListCategories()
	if cs == nil { cs = []store.Category{} }
	wj(w, 200, map[string]any{"categories": cs, "count": len(cs)})
}

func (s *Server) hDelCat(w http.ResponseWriter, r *http.Request) {
	s.db.DeleteCategory(r.PathValue("id"))
	wj(w, 200, map[string]string{"status": "deleted"})
}

func (s *Server) hCreateThread(w http.ResponseWriter, r *http.Request) {
	catID := r.PathValue("id")
	var req struct {
		Title   string `json:"title"`
		Author  string `json:"author"`
		Content string `json:"content"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.Title == "" {
		wj(w, 400, map[string]string{"error": "title required"})
		return
	}
	t, err := s.db.CreateThread(catID, req.Title, req.Author, req.Content)
	if err != nil {
		wj(w, 500, map[string]string{"error": err.Error()})
		return
	}
	wj(w, 201, map[string]any{"thread": t})
}

func (s *Server) hListThreads(w http.ResponseWriter, r *http.Request) {
	ts, _ := s.db.ListThreads(r.PathValue("id"))
	if ts == nil { ts = []store.Thread{} }
	wj(w, 200, map[string]any{"threads": ts, "count": len(ts)})
}

func (s *Server) hGetThread(w http.ResponseWriter, r *http.Request) {
	t, err := s.db.GetThread(r.PathValue("id"))
	if err != nil {
		wj(w, 404, map[string]string{"error": "thread not found"})
		return
	}
	replies, _ := s.db.ListReplies(t.ID)
	if replies == nil { replies = []store.Reply{} }
	wj(w, 200, map[string]any{"thread": t, "replies": replies})
}

func (s *Server) hDelThread(w http.ResponseWriter, r *http.Request) {
	s.db.DeleteThread(r.PathValue("id"))
	wj(w, 200, map[string]string{"status": "deleted"})
}

func (s *Server) hCreateReply(w http.ResponseWriter, r *http.Request) {
	threadID := r.PathValue("id")
	var req struct {
		Author  string `json:"author"`
		Content string `json:"content"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.Content == "" {
		wj(w, 400, map[string]string{"error": "content required"})
		return
	}
	rep, err := s.db.CreateReply(threadID, req.Author, req.Content)
	if err != nil {
		wj(w, 500, map[string]string{"error": err.Error()})
		return
	}
	wj(w, 201, map[string]any{"reply": rep})
}

func (s *Server) hListReplies(w http.ResponseWriter, r *http.Request) {
	rs, _ := s.db.ListReplies(r.PathValue("id"))
	if rs == nil { rs = []store.Reply{} }
	wj(w, 200, map[string]any{"replies": rs, "count": len(rs)})
}

func wj(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
