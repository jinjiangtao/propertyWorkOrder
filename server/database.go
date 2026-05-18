package main

import (
	"database/sql"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

type RepairRequest struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Username     string    `json:"username"`
	RepairType   string    `json:"repair_type"`
	Description  string    `json:"description"`
	ImageURL     string    `json:"image_url"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	WorkerID     int64     `json:"worker_id"`
	WorkerName   string    `json:"worker_name"`
	RepairResult string    `json:"repair_result"`
	RepairImgs   string    `json:"repair_imgs"`
}

type Worker struct {
	ID        int64     `json:"id"`
	WorkNo    string    `json:"work_no"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Password  string    `json:"-"`
	SkillType string    `json:"skill_type"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_sync=OFF")
	if err != nil {
		return err
	}

	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)
	DB.SetConnMaxLifetime(0)

	if _, err := DB.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return err
	}

	if _, err := DB.Exec("PRAGMA synchronous=NORMAL;"); err != nil {
		return err
	}

	if _, err := DB.Exec("PRAGMA cache_size=10000;"); err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	if err = createTables(); err != nil {
		return err
	}

	if err = createAdminUser(); err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

func createTables() error {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		is_admin BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	repairRequestsTable := `
	CREATE TABLE IF NOT EXISTS repair_requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		repair_type TEXT NOT NULL,
		description TEXT NOT NULL,
		image_url TEXT,
		status TEXT DEFAULT '未处理',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		worker_id INTEGER DEFAULT 0,
		worker_name TEXT,
		repair_result TEXT,
		repair_imgs TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	workersTable := `
	CREATE TABLE IF NOT EXISTS workers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		work_no TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		phone TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		skill_type TEXT NOT NULL,
		status INTEGER DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := DB.Exec(usersTable); err != nil {
		return err
	}

	if _, err := DB.Exec(repairRequestsTable); err != nil {
		return err
	}

	if _, err := DB.Exec(workersTable); err != nil {
		return err
	}

	addColumnIfNotExists("repair_requests", "worker_id", "INTEGER DEFAULT 0")
	addColumnIfNotExists("repair_requests", "worker_name", "TEXT")
	addColumnIfNotExists("repair_requests", "repair_result", "TEXT")
	addColumnIfNotExists("repair_requests", "repair_imgs", "TEXT")

	return nil
}

func addColumnIfNotExists(tableName, columnName, columnType string) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM pragma_table_info(?) WHERE name = ?`, tableName, columnName).Scan(&count)
	if err != nil || count == 0 {
		DB.Exec(`ALTER TABLE ` + tableName + ` ADD COLUMN ` + columnName + ` ` + columnType)
	}
}

func createAdminUser() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? AND is_admin = 1", "admin").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		_, err = DB.Exec(
			"INSERT INTO users (username, password, is_admin) VALUES (?, ?, 1)",
			"admin",
			string(hashedPassword),
		)
		if err != nil {
			return err
		}
		log.Println("Admin user created")
	}

	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
