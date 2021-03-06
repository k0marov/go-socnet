package sql_db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
	"github.com/k0marov/go-socnet/features/posts/domain/models"
	"github.com/k0marov/go-socnet/features/posts/domain/values"
)

type SqlDB struct {
	sql       *sqlx.DB
	TableName table_name.TableName
}

func NewSqlDB(sql *sqlx.DB) (*SqlDB, error) {
	err := initSQL(sql)
	if err != nil {
		return nil, err
	}
	return &SqlDB{sql: sql, TableName: table_name.NewTableName("Post")}, nil
}

func initSQL(sql *sqlx.DB) error {
	_, err := sql.Exec(`
		CREATE TABLE IF NOT EXISTS Post(
		    id INTEGER PRIMARY KEY, 
			owner_id INT NOT NULL, 
			textContent TEXT NOT NULL, 
			createdAt INT NOT NULL, 
			FOREIGN KEY(owner_id) REFERENCES Profile(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return core_err.Rethrow("creating Post table", err)
	}
	_, err = sql.Exec(`
		CREATE TABLE IF NOT EXISTS PostImage(
		    post_id INT NOT NULL, 
		    path VARCHAR(255), 
		    ind INT,
		    FOREIGN KEY(post_id) REFERENCES Post(id) ON DELETE CASCADE
		)	
	`)
	if err != nil {
		return core_err.Rethrow("creating PostImage table", err)
	}
	return nil
}

func (db *SqlDB) GetPosts(author core_values.UserId) (posts []models.PostModel, err error) {
	rows, err := db.sql.Query(`
		SELECT id, owner_id, textContent, createdAt
		FROM Post 
		WHERE owner_id = ?
		ORDER BY createdAt DESC
	`, author)
	if err != nil {
		return []models.PostModel{}, core_err.Rethrow("getting posts from db", err)
	}
	defer rows.Close()
	for rows.Next() {
		post := models.PostModel{}
		err = rows.Scan(&post.Id, &post.AuthorId, &post.Text, &post.CreatedAt)
		if err != nil {
			return []models.PostModel{}, core_err.Rethrow("scanning a post", err)
		}
		post.Images, err = db.getImages(post.Id)
		if err != nil {
			return []models.PostModel{}, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (db *SqlDB) CreatePost(newPost models.PostToCreate) (values.PostId, error) {
	res, err := db.sql.Exec(`
		INSERT INTO Post(owner_id, textContent, createdAt) VALUES (?, ?, ?)
	`, newPost.Author, newPost.Text, newPost.CreatedAt.Unix())
	if err != nil {
		return "", core_err.Rethrow("inserting a post", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", core_err.Rethrow("getting the inserted post id", err)
	}
	return fmt.Sprintf("%d", id), nil
}

func (db *SqlDB) AddPostImages(post values.PostId, images []models.PostImageModel) error {
	for _, image := range images {
		err := db.addImage(post, image)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *SqlDB) addImage(post values.PostId, image models.PostImageModel) error {
	_, err := db.sql.Exec(`
		INSERT INTO PostImage(post_id, path, ind) VALUES (?, ?, ?)
   `, post, image.Path, image.Index)
	if err != nil {
		return core_err.Rethrow("inserting a post image", err)
	}
	return nil
}

func (db *SqlDB) getImages(post values.PostId) (images []models.PostImageModel, err error) {
	err = db.sql.Select(&images, `
		SELECT path, ind FROM PostImage WHERE post_id = ?
    `, post)
	if err != nil {
		return []models.PostImageModel{}, core_err.Rethrow("SELECTing post images", err)
	}
	return images, nil
}
