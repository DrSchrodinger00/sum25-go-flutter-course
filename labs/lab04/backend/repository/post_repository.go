// repository/post_repository.go
package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"fmt"

	"lab04-backend/models"

	"github.com/georgysavva/scany/sqlscan"
)

// PostRepository handles database operations for posts using scany
type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(req *models.CreatePostRequest) (*models.Post, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	query := `
		INSERT INTO posts (user_id, title, content, published, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, title, content, published, created_at, updated_at`
	now := time.Now()
	post := &models.Post{}
	if err := sqlscan.Get(
		context.Background(),
		r.db,
		post,
		query,
		req.UserID,
		req.Title,
		req.Content,
		req.Published,
		now,
		now,
	); err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostRepository) GetByID(id int) (*models.Post, error) {
	post := &models.Post{}
	if err := sqlscan.Get(
		context.Background(),
		r.db,
		post,
		"SELECT id, user_id, title, content, published, created_at, updated_at FROM posts WHERE id = $1",
		id,
	); err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostRepository) GetByUserID(userID int) ([]models.Post, error) {
	var posts []models.Post
	if err := sqlscan.Select(
		context.Background(),
		r.db,
		&posts,
		"SELECT id, user_id, title, content, published, created_at, updated_at FROM posts WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) GetPublished() ([]models.Post, error) {
	var posts []models.Post
	if err := sqlscan.Select(
		context.Background(),
		r.db,
		&posts,
		"SELECT id, user_id, title, content, published, created_at, updated_at FROM posts WHERE published = TRUE ORDER BY created_at DESC",
	); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) GetAll() ([]models.Post, error) {
	var posts []models.Post
	if err := sqlscan.Select(
		context.Background(),
		r.db,
		&posts,
		"SELECT id, user_id, title, content, published, created_at, updated_at FROM posts ORDER BY created_at DESC",
	); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) Update(id int, req *models.UpdatePostRequest) (*models.Post, error) {
	// Build dynamic SET clause
	sets := []string{}
	args := []interface{}{}
	if req.Title != nil {
		sets = append(sets, "title = ?")
		args = append(args, *req.Title)
	}
	if req.Content != nil {
		sets = append(sets, "content = ?")
		args = append(args, *req.Content)
	}
	if req.Published != nil {
		sets = append(sets, "published = ?")
		args = append(args, *req.Published)
	}
	now := time.Now()
	sets = append(sets, "updated_at = ?")
	args = append(args, now)
	if len(sets) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}
	args = append(args, id)
	query := fmt.Sprintf(
		"UPDATE posts SET %s WHERE id = ? RETURNING id, user_id, title, content, published, created_at, updated_at",
		strings.Join(sets, ", "),
	)
	post := &models.Post{}
	if err := sqlscan.Get(context.Background(), r.db, post, query, args...); err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PostRepository) Count() (int, error) {
	var cnt int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&cnt); err != nil {
		return 0, err
	}
	return cnt, nil
}

func (r *PostRepository) CountByUserID(userID int) (int, error) {
	var cnt int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = $1", userID).Scan(&cnt); err != nil {
		return 0, err
	}
	return cnt, nil
}
