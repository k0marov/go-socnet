package sql_db

import (
	"database/sql"
	"fmt"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	"github.com/k0marov/go-socnet/core/general/core_values"
	"time"

	"github.com/k0marov/go-socnet/features/posts/domain/models"
	"github.com/k0marov/go-socnet/features/posts/domain/values"
)

type SqlDB struct {
	sql       *sql.DB
	TableName table_name.TableName
}

func NewSqlDB(sql *sql.DB) (*SqlDB, error) {
	err := initSQL(sql)
	if err != nil {
		return nil, err
	}
	return &SqlDB{sql: sql, TableName: table_name.NewTableName("Post")}, nil
}

func initSQL(sql *sql.DB) error {
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
		return fmt.Errorf("while creating Post table: %w", err)
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
		return fmt.Errorf("while creating PostImage table: %w", err)
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
		return []models.PostModel{}, fmt.Errorf("while getting posts from db: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		post := models.PostModel{}
		var createdAt int64
		err = rows.Scan(&post.Id, &post.AuthorId, &post.Text, &createdAt)
		if err != nil {
			return []models.PostModel{}, fmt.Errorf("while scanning a post: %w", err)
		}
		post.CreatedAt = time.Unix(createdAt, 0).UTC()
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
		return "", fmt.Errorf("while inserting a post: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("while getting the inserted post id: %w", err)
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
		return fmt.Errorf("while inserting a post image: %w", err)
	}
	return nil
}

func (db *SqlDB) getImages(post values.PostId) (images []models.PostImageModel, err error) {
	rows, err := db.sql.Query(`
		SELECT path, ind FROM PostImage WHERE post_id = ?
    `, post)
	if err != nil {
		return []models.PostImageModel{}, fmt.Errorf("while SELECTing post images: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		image := models.PostImageModel{}
		err := rows.Scan(&image.Path, &image.Index)
		if err != nil {
			return []models.PostImageModel{}, fmt.Errorf("while scanning an image: %w", err)
		}
		images = append(images, image)
	}
	return images, nil
}
