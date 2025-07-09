package models

import (
	"database/sql"
	"errors"
	"net/mail"
	"time"
)

// User represents a user in the system
// Scans and JSON/DB mapping
type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest represents the payload for creating a user
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateUserRequest represents the payload for updating a user
type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

// Validate ensures a User has valid fields
func (u *User) Validate() error {
	if len(u.Name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("invalid email format")
	}
	return nil
}

// Validate ensures CreateUserRequest has valid fields
func (req *CreateUserRequest) Validate() error {
	if len(req.Name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return errors.New("invalid email format")
	}
	return nil
}

// ToUser converts CreateUserRequest into a User model
func (req *CreateUserRequest) ToUser() *User {
	now := time.Now()
	return &User{
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ScanRow scans a single sql.Row into the User struct
func (u *User) ScanRow(row *sql.Row) error {
	return row.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
}

// ScanUsers scans multiple sql.Rows into a slice of User
func ScanUsers(rows *sql.Rows) ([]User, error) {
	defer rows.Close()
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
