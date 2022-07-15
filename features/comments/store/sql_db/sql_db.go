package sql_db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"time"

	"github.com/k0marov/go-socnet/features/comments/domain/models"
	"github.com/k0marov/go-socnet/features/comments/domain/values"
	post_values "github.com/k0marov/go-socnet/features/posts/domain/values"
)

type SqlDB struct {
	sql       *sqlx.DB
	TableName table_name.TableName
}

func NewSqlDB(db *sqlx.DB) (*SqlDB, error) {
	err := initSQL(db)
	if err != nil {
		return nil, core_err.Rethrow("initializing sql for comments", err)
	}
	return &SqlDB{db, table_name.NewTableName("Comment")}, nil
}

func initSQL(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Comment(
		    id INTEGER PRIMARY KEY, 
		   	post_id INT NOT NULL, 
		   	owner_id INT NOT NULL, 
		   	textContent TEXT NOT NULL, 
		   	createdAt INT NOT NULL, 
		   	FOREIGN KEY(post_id) REFERENCES Post(id) ON DELETE CASCADE, 
		   	FOREIGN KEY(owner_id) REFERENCES Profile(id) ON DELETE CASCADE
		)
   `)
	if err != nil {
		return core_err.Rethrow("creating Comment table", err)
	}
	return nil
}

func (db *SqlDB) GetComments(post post_values.PostId) ([]models.CommentModel, error) {
	rows, err := db.sql.Query(`
		SELECT id, owner_id, textContent, createdAt
		FROM Comment 
		WHERE post_id = ?
		ORDER BY createdAt DESC
    `, post)
	if err != nil {
		return []models.CommentModel{}, core_err.Rethrow("SELECTing post comments", err)
	}
	var comments []models.CommentModel
	for rows.Next() {
		comment := models.CommentModel{}
		var createdAtUnix int64
		err := rows.Scan(&comment.Id, &comment.AuthorId, &comment.Text, &createdAtUnix)
		if err != nil {
			return []models.CommentModel{}, core_err.Rethrow("scanning a comment", err)
		}
		comment.CreatedAt = time.Unix(createdAtUnix, 0).UTC()
		comments = append(comments, comment)
	}
	return comments, nil
}

func (db *SqlDB) Create(newComment values.NewCommentValue, createdAt time.Time) (values.CommentId, error) {
	res, err := db.sql.Exec(`
		INSERT INTO Comment(post_id, owner_id, textContent, createdAt) 
		VALUES (?, ?, ?, ?)
    `, newComment.Post, newComment.Author, newComment.Text, createdAt.Unix())
	if err != nil {
		return "", core_err.Rethrow("INSERTing a new comment", err)
	}
	newId, err := res.LastInsertId()
	if err != nil {
		return "", core_err.Rethrow("getting the ID of newly inserted comment", err)
	}
	return fmt.Sprintf("%d", newId), nil
}
