package sql_db

import (
	"database/sql"
	"fmt"
	"github.com/k0marov/go-socnet/core/abstract/likeable/table_name"
	"github.com/k0marov/go-socnet/core/general/core_errors"
	"github.com/k0marov/go-socnet/core/general/core_values"
	"time"

	"github.com/k0marov/go-socnet/features/comments/domain/models"
	"github.com/k0marov/go-socnet/features/comments/domain/values"
	post_values "github.com/k0marov/go-socnet/features/posts/domain/values"
)

type SqlDB struct {
	sql       *sql.DB
	TableName table_name.TableName
}

func NewSqlDB(db *sql.DB) (*SqlDB, error) {
	err := initSQL(db)
	if err != nil {
		return nil, fmt.Errorf("while initializing sql for comments: %w", err)
	}
	return &SqlDB{db, table_name.NewTableName("Comment")}, nil
}

func initSQL(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Comment(
		    id INTEGER PRIMARY KEY, 
		   	post_id INT NOT NULL, 
		   	author_id INT NOT NULL, 
		   	textContent TEXT NOT NULL, 
		   	createdAt INT NOT NULL, 
		   	FOREIGN KEY(post_id) REFERENCES Post(id) ON DELETE CASCADE, 
		   	FOREIGN KEY(author_id) REFERENCES Profile(id) ON DELETE CASCADE
		)
   `)
	if err != nil {
		return fmt.Errorf("while creating Comment table: %w", err)
	}
	return nil
}

func (db *SqlDB) GetComments(post post_values.PostId) ([]models.CommentModel, error) {
	rows, err := db.sql.Query(`
		SELECT id, author_id, textContent, createdAt
		FROM Comment 
		WHERE post_id = ?
		ORDER BY createdAt DESC
    `, post)
	if err != nil {
		return []models.CommentModel{}, fmt.Errorf("while SELECTing post comments: %w", err)
	}
	var comments []models.CommentModel
	for rows.Next() {
		comment := models.CommentModel{}
		var createdAtUnix int64
		err := rows.Scan(&comment.Id, &comment.AuthorId, &comment.Text, &createdAtUnix)
		if err != nil {
			return []models.CommentModel{}, fmt.Errorf("while scanning a comment: %w", err)
		}
		comment.CreatedAt = time.Unix(createdAtUnix, 0).UTC()
		comments = append(comments, comment)
	}
	return comments, nil
}

func (db *SqlDB) GetAuthor(comment values.CommentId) (core_values.UserId, error) {
	row := db.sql.QueryRow(`
		SELECT author_id FROM Comment WHERE id = ?
    `, comment)
	var authorId core_values.UserId
	err := row.Scan(&authorId)
	if err == sql.ErrNoRows {
		return "", core_errors.ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("while scanning the SELECTed author_id: %w", err)
	}
	return authorId, nil
}

func (db *SqlDB) Create(newComment values.NewCommentValue, createdAt time.Time) (values.CommentId, error) {
	res, err := db.sql.Exec(`
		INSERT INTO Comment(post_id, author_id, textContent, createdAt) 
		VALUES (?, ?, ?, ?)
    `, newComment.Post, newComment.Author, newComment.Text, createdAt.Unix())
	if err != nil {
		return "", fmt.Errorf("while INSERTing a new comment: %w", err)
	}
	newId, err := res.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("while getting the ID of newly inserted comment: %w", err)
	}
	return fmt.Sprintf("%d", newId), nil
}
