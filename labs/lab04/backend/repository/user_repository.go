// repository/user_repository.go
package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"lab04-backend/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(req *models.CreateUserRequest) (*models.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	user := req.ToUser()
	res, err := r.db.Exec(
		"INSERT INTO users (name, email, created_at, updated_at) VALUES (?, ?, ?, ?)",
		user.Name, user.Email, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = int(id)
	return user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	row := r.db.QueryRow(
		"SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?",
		id,
	)
	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	row := r.db.QueryRow(
		"SELECT id, name, email, created_at, updated_at FROM users WHERE email = ?",
		email,
	)
	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	rows, err := r.db.Query(
		"SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Update(id int, req *models.UpdateUserRequest) (*models.User, error) {
	sets := []string{}
	args := []interface{}{}
	if req.Name != nil {
		sets = append(sets, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Email != nil {
		sets = append(sets, "email = ?")
		args = append(args, *req.Email)
	}
	// always update timestamp
	sets = append(sets, "updated_at = ?")
	now := time.Now()
	args = append(args, now)

	if len(sets) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}
	args = append(args, id)
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(sets, ", "))
	res, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, sql.ErrNoRows
	}
	return r.GetByID(id)
}

func (r *UserRepository) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepository) Count() (int, error) {
	row := r.db.QueryRow("SELECT COUNT(*) FROM users")
	var cnt int
	if err := row.Scan(&cnt); err != nil {
		return 0, err
	}
	return cnt, nil
}
