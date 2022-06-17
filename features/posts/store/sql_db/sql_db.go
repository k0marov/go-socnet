package sql_db

import (
	"database/sql"
	"fmt"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/posts/domain/values"
	"github.com/k0marov/socnet/features/posts/store/models"
	"time"
)

type SqlDB struct {
	sql *sql.DB
}

func NewSqlDB(sql *sql.DB) (*SqlDB, error) {
	err := initSQL(sql)
	if err != nil {
		return nil, err
	}
	return &SqlDB{sql: sql}, nil
}

func initSQL(sql *sql.DB) error {
	_, err := sql.Exec(`
		CREATE TABLE IF NOT EXISTS Post(
		    id INTEGER PRIMARY KEY, 
			author_id INT NOT NULL, 
			textContent TEXT NOT NULL, 
			createdAt INT NOT NULL, 
			FOREIGN KEY(author_id) REFERENCES Profile(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("while creating Post table: %w", err)
	}
	return nil
}

func (db *SqlDB) GetPosts(author core_values.UserId) (posts []models.PostModel, err error) {
	rows, err := db.sql.Query(`
		SELECT id, author_id, textContent, createdAt FROM Post 
		WHERE author_id = ?
	`, author)
	defer rows.Close()
	for rows.Next() {
		post := models.PostModel{}
		var createdAt int64
		err = rows.Scan(&post.Id, &post.Author, &post.Text, &createdAt)
		post.CreatedAt = time.Unix(createdAt, 0)
		post.Images = []core_values.FileURL{}
		posts = append(posts, post)
	}
	return posts, nil
}
func (db *SqlDB) LikePost(post values.PostId, fromUser core_values.UserId) error {
	panic("unimplemented")
}
func (db *SqlDB) UnlikePost(post values.PostId, fromUser core_values.UserId) error {
	panic("unimplemented")
}
func (db *SqlDB) IsLiked(post values.PostId, byUser core_values.UserId) error {
	panic("unimplemented")
}
func (db *SqlDB) GetAuthor(post values.PostId) (core_values.UserId, error) {
	row := db.sql.QueryRow(`
		SELECT author_id FROM Post 
		WHERE id = ?
    `, post)
	var authorId int
	row.Scan(&authorId)
	return fmt.Sprintf("%d", authorId), nil
}
func (db *SqlDB) CreatePost(newPost models.PostToCreate) (values.PostId, error) {
	res, _ := db.sql.Exec(`
		INSERT INTO Post(author_id, textContent, createdAt) VALUES (?, ?, ?)
	`, newPost.Author, newPost.Text, newPost.CreatedAt.Unix())
	id, _ := res.LastInsertId()
	return fmt.Sprintf("%d", id), nil
}
func (db *SqlDB) AddPostImages(post values.PostId, imagePaths []core_values.StaticFilePath) error {
	panic("unimplemented")
}
func (db *SqlDB) DeletePost(post values.PostId) error {
	panic("unimplemented")
}
