package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Standup struct{
	ID string `json:"id"`
	Author string `json:"author"`
	Yesterday string `json:"yesterday"`
	Today string `json:"today"`
	Blockers string `json:"blockers"`
	Mood string `json:"mood"`
	Date string `json:"date"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"campfire.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS standups(id TEXT PRIMARY KEY,author TEXT NOT NULL,yesterday TEXT DEFAULT '',today TEXT DEFAULT '',blockers TEXT DEFAULT '',mood TEXT DEFAULT '',date TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Standup)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO standups(id,author,yesterday,today,blockers,mood,date,created_at)VALUES(?,?,?,?,?,?,?,?)`,e.ID,e.Author,e.Yesterday,e.Today,e.Blockers,e.Mood,e.Date,e.CreatedAt);return err}
func(d *DB)Get(id string)*Standup{var e Standup;if d.db.QueryRow(`SELECT id,author,yesterday,today,blockers,mood,date,created_at FROM standups WHERE id=?`,id).Scan(&e.ID,&e.Author,&e.Yesterday,&e.Today,&e.Blockers,&e.Mood,&e.Date,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Standup{rows,_:=d.db.Query(`SELECT id,author,yesterday,today,blockers,mood,date,created_at FROM standups ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Standup;for rows.Next(){var e Standup;rows.Scan(&e.ID,&e.Author,&e.Yesterday,&e.Today,&e.Blockers,&e.Mood,&e.Date,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM standups WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM standups`).Scan(&n);return n}
