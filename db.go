package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

type Post struct {
	ID        int
	Title     string
	Slug      string
	Content   string
	Summary   string
	Published bool
	CreatedAt time.Time
}

func initDB() {
	dbPath := globalConfig.Database.DBPath
	if dbPath == "" {
		dbPath = "data/artdotech.db"
	}

	// Asegurar que el directorio existe
	dir := filepath.Dir(dbPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		slug TEXT UNIQUE NOT NULL,
		content TEXT NOT NULL,
		summary TEXT,
		published BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Error creando tabla posts: %v", err)
	}

	log.Println("Base de datos SQLite inicializada correctamente.")
}

// CRUD Operations

func getAllPosts(onlyPublished bool) ([]Post, error) {
	query := "SELECT id, title, slug, content, summary, published, created_at FROM posts ORDER BY created_at DESC"
	if onlyPublished {
		query = "SELECT id, title, slug, content, summary, published, created_at FROM posts WHERE published = 1 ORDER BY created_at DESC"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.Summary, &p.Published, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func getPostBySlug(slug string) (*Post, error) {
	var p Post
	err := db.QueryRow("SELECT id, title, slug, content, summary, published, created_at FROM posts WHERE slug = ?", slug).
		Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.Summary, &p.Published, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func createPost(p Post) error {
	_, err := db.Exec("INSERT INTO posts (title, slug, content, summary, published) VALUES (?, ?, ?, ?, ?)",
		p.Title, p.Slug, p.Content, p.Summary, p.Published)
	return err
}

func updatePost(p Post) error {
	_, err := db.Exec("UPDATE posts SET title = ?, slug = ?, content = ?, summary = ?, published = ? WHERE id = ?",
		p.Title, p.Slug, p.Content, p.Summary, p.Published, p.ID)
	return err
}

func deletePost(id int) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", id)
	return err
}
