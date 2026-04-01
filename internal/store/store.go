package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"os"
	"path/filepath"
	"time"
	_ "modernc.org/sqlite"
)

type DB struct{ conn *sql.DB }

func Open(dataDir string) (*DB, error) {
	os.MkdirAll(dataDir, 0755)
	conn, err := sql.Open("sqlite", filepath.Join(dataDir, "campfire.db"))
	if err != nil { return nil, err }
	conn.Exec("PRAGMA journal_mode=WAL"); conn.Exec("PRAGMA busy_timeout=5000"); conn.SetMaxOpenConns(4)
	db := &DB{conn: conn}; return db, db.migrate()
}
func (db *DB) Close() error { return db.conn.Close() }

func (db *DB) migrate() error {
	_, err := db.conn.Exec(`
CREATE TABLE IF NOT EXISTS categories (
    id TEXT PRIMARY KEY, name TEXT NOT NULL, slug TEXT NOT NULL UNIQUE,
    description TEXT DEFAULT '', sort_order INTEGER DEFAULT 0,
    created_at TEXT DEFAULT (datetime('now'))
);
CREATE TABLE IF NOT EXISTS threads (
    id TEXT PRIMARY KEY, category_id TEXT NOT NULL, title TEXT NOT NULL,
    author TEXT DEFAULT 'anonymous', content TEXT DEFAULT '',
    pinned INTEGER DEFAULT 0, locked INTEGER DEFAULT 0,
    reply_count INTEGER DEFAULT 0, last_reply_at TEXT DEFAULT '',
    created_at TEXT DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_threads_cat ON threads(category_id);

CREATE TABLE IF NOT EXISTS replies (
    id TEXT PRIMARY KEY, thread_id TEXT NOT NULL, author TEXT DEFAULT 'anonymous',
    content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_replies_thread ON replies(thread_id);
`)
	return err
}

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	ThreadCount int    `json:"thread_count"`
	CreatedAt   string `json:"created_at"`
}

func (db *DB) CreateCategory(name, slug, description string) (*Category, error) {
	id := "cat_" + gid(6); now := time.Now().UTC().Format(time.RFC3339)
	_, err := db.conn.Exec("INSERT INTO categories (id,name,slug,description,created_at) VALUES (?,?,?,?,?)", id, name, slug, description, now)
	if err != nil { return nil, err }
	return &Category{ID: id, Name: name, Slug: slug, Description: description, CreatedAt: now}, nil
}

func (db *DB) ListCategories() ([]Category, error) {
	rows, err := db.conn.Query(`SELECT c.id, c.name, c.slug, c.description,
		(SELECT COUNT(*) FROM threads WHERE category_id=c.id), c.created_at
		FROM categories c ORDER BY c.sort_order, c.name`)
	if err != nil { return nil, err }; defer rows.Close()
	var out []Category
	for rows.Next() { var c Category; rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.ThreadCount, &c.CreatedAt); out = append(out, c) }
	return out, rows.Err()
}

func (db *DB) DeleteCategory(id string) {
	db.conn.Exec("DELETE FROM replies WHERE thread_id IN (SELECT id FROM threads WHERE category_id=?)", id)
	db.conn.Exec("DELETE FROM threads WHERE category_id=?", id)
	db.conn.Exec("DELETE FROM categories WHERE id=?", id)
}

type Thread struct {
	ID          string `json:"id"`
	CategoryID  string `json:"category_id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Content     string `json:"content"`
	Pinned      bool   `json:"pinned"`
	Locked      bool   `json:"locked"`
	ReplyCount  int    `json:"reply_count"`
	LastReplyAt string `json:"last_reply_at,omitempty"`
	CreatedAt   string `json:"created_at"`
}

func (db *DB) CreateThread(categoryID, title, author, content string) (*Thread, error) {
	id := "thr_" + gid(8); now := time.Now().UTC().Format(time.RFC3339)
	if author == "" { author = "anonymous" }
	_, err := db.conn.Exec("INSERT INTO threads (id,category_id,title,author,content,last_reply_at,created_at) VALUES (?,?,?,?,?,?,?)",
		id, categoryID, title, author, content, now, now)
	if err != nil { return nil, err }
	return &Thread{ID: id, CategoryID: categoryID, Title: title, Author: author, Content: content, LastReplyAt: now, CreatedAt: now}, nil
}

func (db *DB) ListThreads(categoryID string) ([]Thread, error) {
	rows, err := db.conn.Query("SELECT id,category_id,title,author,content,pinned,locked,reply_count,last_reply_at,created_at FROM threads WHERE category_id=? ORDER BY pinned DESC, last_reply_at DESC", categoryID)
	if err != nil { return nil, err }; defer rows.Close()
	var out []Thread
	for rows.Next() {
		var t Thread; var p, l int
		rows.Scan(&t.ID, &t.CategoryID, &t.Title, &t.Author, &t.Content, &p, &l, &t.ReplyCount, &t.LastReplyAt, &t.CreatedAt)
		t.Pinned = p == 1; t.Locked = l == 1
		out = append(out, t)
	}
	return out, rows.Err()
}

func (db *DB) GetThread(id string) (*Thread, error) {
	var t Thread; var p, l int
	err := db.conn.QueryRow("SELECT id,category_id,title,author,content,pinned,locked,reply_count,last_reply_at,created_at FROM threads WHERE id=?", id).
		Scan(&t.ID, &t.CategoryID, &t.Title, &t.Author, &t.Content, &p, &l, &t.ReplyCount, &t.LastReplyAt, &t.CreatedAt)
	t.Pinned = p == 1; t.Locked = l == 1
	return &t, err
}

func (db *DB) DeleteThread(id string) {
	db.conn.Exec("DELETE FROM replies WHERE thread_id=?", id)
	db.conn.Exec("DELETE FROM threads WHERE id=?", id)
}

type Reply struct {
	ID        string `json:"id"`
	ThreadID  string `json:"thread_id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func (db *DB) CreateReply(threadID, author, content string) (*Reply, error) {
	id := "rep_" + gid(8); now := time.Now().UTC().Format(time.RFC3339)
	if author == "" { author = "anonymous" }
	_, err := db.conn.Exec("INSERT INTO replies (id,thread_id,author,content,created_at) VALUES (?,?,?,?,?)", id, threadID, author, content, now)
	if err != nil { return nil, err }
	db.conn.Exec("UPDATE threads SET reply_count=reply_count+1, last_reply_at=? WHERE id=?", now, threadID)
	return &Reply{ID: id, ThreadID: threadID, Author: author, Content: content, CreatedAt: now}, nil
}

func (db *DB) ListReplies(threadID string) ([]Reply, error) {
	rows, err := db.conn.Query("SELECT id,thread_id,author,content,created_at FROM replies WHERE thread_id=? ORDER BY created_at ASC", threadID)
	if err != nil { return nil, err }; defer rows.Close()
	var out []Reply
	for rows.Next() { var r Reply; rows.Scan(&r.ID, &r.ThreadID, &r.Author, &r.Content, &r.CreatedAt); out = append(out, r) }
	return out, rows.Err()
}

func (db *DB) Stats() map[string]any {
	var cats, threads, replies int
	db.conn.QueryRow("SELECT COUNT(*) FROM categories").Scan(&cats)
	db.conn.QueryRow("SELECT COUNT(*) FROM threads").Scan(&threads)
	db.conn.QueryRow("SELECT COUNT(*) FROM replies").Scan(&replies)
	return map[string]any{"categories": cats, "threads": threads, "replies": replies}
}

func gid(n int) string { b := make([]byte, n); rand.Read(b); return hex.EncodeToString(b) }
