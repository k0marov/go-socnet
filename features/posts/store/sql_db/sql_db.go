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
	_, err = sql.Exec(`
		CREATE TABLE IF NOT EXISTS PostLike(
		    post_id INT NOT NULL,
			profile_id INT NOT NULL,
		   	FOREIGN KEY(post_id) REFERENCES Post(id) ON DELETE CASCADE, 
		   	FOREIGN KEY(profile_id) REFERENCES Profile(id) ON DELETE CASCADE 
		)	
    `)
	return nil
}

func (db *SqlDB) GetPosts(author core_values.UserId) (posts []models.PostModel, err error) {
	rows, err := db.sql.Query(`
		SELECT id, author_id, textContent, createdAt FROM Post 
		WHERE author_id = ?
		ORDER BY createdAt DESC
	`, author)
	if err != nil {
		return []models.PostModel{}, fmt.Errorf("while getting posts from db: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		post := models.PostModel{}
		var createdAt int64
		err = rows.Scan(&post.Id, &post.Author, &post.Text, &createdAt)
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
func (db *SqlDB) LikePost(post values.PostId, fromUser core_values.UserId) error {
	_, err := db.sql.Exec(`
		INSERT INTO PostLike(post_id, profile_id) VALUES(?, ?)
    `, post, fromUser)
	if err != nil {
		return fmt.Errorf("while INSERTing a new PostLike: %w", err)
	}
	return nil
}
func (db *SqlDB) UnlikePost(post values.PostId, fromUser core_values.UserId) error {
	_, err := db.sql.Exec(`
		DELETE FROM PostLike WHERE post_id = ? AND profile_id = ?
    `, post, fromUser)
	if err != nil {
		return fmt.Errorf("while DELETEing a PostLike: %w", err)
	}
	return nil
}
func (db *SqlDB) IsLiked(post values.PostId, byUser core_values.UserId) (bool, error) {
	row := db.sql.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM PostLike WHERE post_id = ? AND profile_id = ?)
    `, post, byUser)
	isLiked := 0
	err := row.Scan(&isLiked)
	if err != nil {
		return false, fmt.Errorf("while SELECTing is post liked: %w", err)
	}
	return isLiked == 1, nil
}
func (db *SqlDB) GetAuthor(post values.PostId) (core_values.UserId, error) {
	row := db.sql.QueryRow(`
		SELECT author_id FROM Post 
		WHERE id = ?
    `, post)
	var authorId int
	err := row.Scan(&authorId)
	if err != nil {
		return "", fmt.Errorf("while SELECTing a post author: %w", err)
	}
	return fmt.Sprintf("%d", authorId), nil
}
func (db *SqlDB) CreatePost(newPost models.PostToCreate) (values.PostId, error) {
	res, err := db.sql.Exec(`
		INSERT INTO Post(author_id, textContent, createdAt) VALUES (?, ?, ?)
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
func (db *SqlDB) AddPostImages(post values.PostId, images []values.PostImage) error {
	for _, image := range images {
		err := db.addImage(post, image)
		if err != nil {
			return err
		}
	}
	return nil
}
func (db *SqlDB) DeletePost(post values.PostId) error {
	_, err := db.sql.Exec(`
		DELETE FROM Post WHERE id = ?
    `, post)
	if err != nil {
		return fmt.Errorf("while DELETEing a post by id: %w", err)
	}
	return nil
}

func (db *SqlDB) addImage(post values.PostId, image values.PostImage) error {
	_, err := db.sql.Exec(`
		INSERT INTO PostImage(post_id, path, ind) VALUES (?, ?, ?)
   `, post, image.Path, image.Index)
	if err != nil {
		return fmt.Errorf("while inserting a post image: %w", err)
	}
	return nil
}
func (db *SqlDB) getImages(post values.PostId) (images []values.PostImage, err error) {
	rows, err := db.sql.Query(`
		SELECT path, ind FROM PostImage WHERE post_id = ?
    `, post)
	if err != nil {
		return []values.PostImage{}, fmt.Errorf("while SELECTing post images: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		image := values.PostImage{}
		err := rows.Scan(&image.Path, &image.Index)
		if err != nil {
			return []values.PostImage{}, fmt.Errorf("while scanning an image: %w", err)
		}
		images = append(images, image)
	}
	return images, nil
}
