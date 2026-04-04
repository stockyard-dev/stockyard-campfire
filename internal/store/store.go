package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Channel struct{ID string `json:"id"`;Name string `json:"name"`;Description string `json:"description"`;CreatedAt string `json:"created_at"`;MessageCount int `json:"message_count"`}
type Message struct{ID string `json:"id"`;ChannelID string `json:"channel_id"`;Author string `json:"author"`;Body string `json:"body"`;CreatedAt string `json:"created_at"`}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"campfire.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS channels(id TEXT PRIMARY KEY,name TEXT UNIQUE NOT NULL,description TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
db.Exec(`CREATE TABLE IF NOT EXISTS messages(id TEXT PRIMARY KEY,channel_id TEXT NOT NULL,author TEXT DEFAULT 'anonymous',body TEXT NOT NULL,created_at TEXT DEFAULT(datetime('now')))`)
db.Exec(`CREATE INDEX IF NOT EXISTS idx_msg_chan ON messages(channel_id,created_at)`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)CreateChannel(c *Channel)error{c.ID=genID();c.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO channels(id,name,description,created_at)VALUES(?,?,?,?)`,c.ID,c.Name,c.Description,c.CreatedAt);return err}
func(d *DB)ListChannels()[]Channel{rows,_:=d.db.Query(`SELECT id,name,description,created_at FROM channels ORDER BY name`);if rows==nil{return nil};defer rows.Close()
var o []Channel;for rows.Next(){var c Channel;rows.Scan(&c.ID,&c.Name,&c.Description,&c.CreatedAt);d.db.QueryRow(`SELECT COUNT(*) FROM messages WHERE channel_id=?`,c.ID).Scan(&c.MessageCount);o=append(o,c)};return o}
func(d *DB)DeleteChannel(id string)error{d.db.Exec(`DELETE FROM messages WHERE channel_id=?`,id);_,err:=d.db.Exec(`DELETE FROM channels WHERE id=?`,id);return err}
func(d *DB)PostMessage(m *Message)error{m.ID=genID();m.CreatedAt=now();if m.Author==""{m.Author="anonymous"};_,err:=d.db.Exec(`INSERT INTO messages(id,channel_id,author,body,created_at)VALUES(?,?,?,?,?)`,m.ID,m.ChannelID,m.Author,m.Body,m.CreatedAt);return err}
func(d *DB)ListMessages(channelID string,limit int)[]Message{if limit<=0{limit=100};rows,_:=d.db.Query(`SELECT id,channel_id,author,body,created_at FROM messages WHERE channel_id=? ORDER BY created_at DESC LIMIT ?`,channelID,limit);if rows==nil{return nil};defer rows.Close()
var o []Message;for rows.Next(){var m Message;rows.Scan(&m.ID,&m.ChannelID,&m.Author,&m.Body,&m.CreatedAt);o=append(o,m)};return o}
func(d *DB)DeleteMessage(id string)error{_,err:=d.db.Exec(`DELETE FROM messages WHERE id=?`,id);return err}
func(d *DB)Stats()map[string]any{var ch,msg int;d.db.QueryRow(`SELECT COUNT(*) FROM channels`).Scan(&ch);d.db.QueryRow(`SELECT COUNT(*) FROM messages`).Scan(&msg);return map[string]any{"channels":ch,"messages":msg}}
