package sql_db

import (
	"database/sql"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/posts/domain/entities"
	"github.com/k0marov/socnet/features/posts/domain/values"
	"github.com/k0marov/socnet/features/posts/store/models"
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
	return nil
}

//DBPostsGetter  func(core_values.UserId) ([]entities.Post, error)
//DBLiker        func(values.PostId, core_values.UserId) error
//DBUnliker      func(values.PostId, core_values.UserId) error
//DBLikeChecker  func(values.PostId, core_values.UserId) (bool, error)
//DBAuthorGetter func(values.PostId) (core_values.UserId, error)
//
//DBPostCreator     func(newPost models.PostToCreate) (values.PostId, error)
//DBPostImagesAdder func(values.PostId, []core_values.StaticFilePath) error
//DBPostDeleter     func(values.PostId) error

func (db *SqlDB) GetPosts(author core_values.UserId) ([]entities.Post, error) {
	panic("unimplemented")
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
	panic("unimplemented")
}
func (db *SqlDB) CreatePost(newPost models.PostToCreate) (values.PostId, error) {
	panic("unimplemented")
}
func (db *SqlDB) AddPostImages(post values.PostId, imagePaths []core_values.StaticFilePath) error {
	panic("unimplemented")
}
func (db *SqlDB) PostDeleter(post values.PostId) error {
	panic("unimplemented")
}
