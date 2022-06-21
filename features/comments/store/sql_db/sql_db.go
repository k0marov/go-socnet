package sql_db

import (
	"database/sql"
	"fmt"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/comments/domain/values"
	"github.com/k0marov/socnet/features/comments/store/models"
	post_values "github.com/k0marov/socnet/features/posts/domain/values"
	"time"
)

type SqlDB struct {
	sql *sql.DB
}

func NewSqlDB(db *sql.DB) (*SqlDB, error) {
	err := initSQL(db)
	if err != nil {
		return nil, fmt.Errorf("while initializing sql for comments: %w", err)
	}
	return &SqlDB{db}, nil
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
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS CommentLike(
		    comment_id INT NOT NULL, 
		    liker_id INT NOT NULL, 
		    FOREIGN KEY(comment_id) REFERENCES Comment(id) ON DELETE CASCADE, 
		    FOREIGN KEY(liker_id) REFERENCES Profile(id) ON DELETE CASCADE
		)	
    `)
	if err != nil {
		return fmt.Errorf("while creating CommentLike table: %w", err)
	}
	return nil
}

func (db *SqlDB) IsLiked(comment values.CommentId, caller core_values.UserId) (bool, error) {
	row := db.sql.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM CommentLike WHERE comment_id = ? AND liker_id = ?)
	`, comment, caller)
	isLiked := 0
	err := row.Scan(&isLiked)
	if err != nil {
		return false, fmt.Errorf("while SELECTing is post liked: %w", err)
	}
	return isLiked == 1, nil
}
func (db *SqlDB) Like(comment values.CommentId, liker core_values.UserId) error {
	_, err := db.sql.Exec(`
		INSERT INTO CommentLike(comment_id, liker_id) VALUES(?, ?)
    `, comment, liker)
	if err != nil {
		return fmt.Errorf("while INSERTing a new PostLike: %w", err)
	}
	return nil
}
func (db *SqlDB) Unlike(comment values.CommentId, unliker core_values.UserId) error {
	_, err := db.sql.Exec(`
		DELETE FROM CommentLike WHERE comment_id = ? AND liker_id = ?
	`, comment, unliker)
	if err != nil {
		return fmt.Errorf("while DELETEing a PostLike: %w", err)
	}
	return nil
}
func (db *SqlDB) GetComments(post post_values.PostId) ([]models.CommentModel, error) {
	rows, err := db.sql.Query(`
		SELECT id, author_id, textContent, createdAt, 
		    (SELECT COUNT(*) FROM CommentLike WHERE comment_id = Comment.id) AS likes
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
		err := rows.Scan(&comment.Id, &comment.Author, &comment.Text, &createdAtUnix, &comment.Likes)
		if err != nil {
			return []models.CommentModel{}, fmt.Errorf("while scanning a comment: %w", err)
		}
		comment.CreatedAt = time.Unix(createdAtUnix, 0).UTC()
		comments = append(comments, comment)
	}
	return comments, nil
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
