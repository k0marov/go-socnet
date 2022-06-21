package sql_db

import (
	"database/sql"
	"fmt"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/comments/domain/values"
	"github.com/k0marov/socnet/features/comments/store/models"
	post_values "github.com/k0marov/socnet/features/posts/domain/values"
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
	return nil
}

func (db *SqlDB) IsLiked(comment values.CommentId, caller core_values.UserId) (bool, error) {
	panic("unimplemented")
}
func (db *SqlDB) Like(comment values.CommentId, liker core_values.UserId) error {
	panic("unimplemented")
}
func (db *SqlDB) Unlike(comment values.CommentId, unliker core_values.UserId) error {
	panic("unimplemented")
}
func (db *SqlDB) GetComments(post post_values.PostId) ([]models.CommentModel, error) {
	panic("unimplemented")
}

func (db *SqlDB) Create(newComment values.NewCommentValue) (models.CommentModel, error) {
	panic("unimplemented")
}
