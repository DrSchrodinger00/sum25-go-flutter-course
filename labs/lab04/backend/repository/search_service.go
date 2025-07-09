// repository/search_service.go
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"lab04-backend/models"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/sqlscan"
)

type SearchService struct {
	db   *sql.DB
	psql squirrel.StatementBuilderType
}

type SearchFilters struct {
	Query        string
	UserID       *int
	Published    *bool
	MinWordCount *int
	Limit        int
	Offset       int
	OrderBy      string
	OrderDir     string
}

func NewSearchService(db *sql.DB) *SearchService {
	return &SearchService{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (s *SearchService) SearchPosts(ctx context.Context, filters SearchFilters) ([]models.Post, error) {
	base := s.psql.
		Select("id", "user_id", "title", "content", "published", "created_at", "updated_at").
		From("posts")

	// dynamic WHERE
	query := s.BuildDynamicQuery(base, filters)

	// order, limit, offset
	if filters.OrderBy != "" && (filters.OrderDir == "ASC" || filters.OrderDir == "DESC") {
		query = query.OrderBy(fmt.Sprintf("%s %s", filters.OrderBy, filters.OrderDir))
	}
	if filters.Limit > 0 {
		query = query.Limit(uint64(filters.Limit))
	}
	if filters.Offset > 0 {
		query = query.Offset(uint64(filters.Offset))
	}

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var posts []models.Post
	if err := sqlscan.Select(ctx, s.db, &posts, sqlStr, args...); err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *SearchService) SearchUsers(ctx context.Context, nameQuery string, limit int) ([]models.User, error) {
	builder := s.psql.
		Select("id", "name", "email", "created_at", "updated_at").
		From("users").
		Where(squirrel.Like{"name": "%" + nameQuery + "%"}).
		OrderBy("name")

	if limit > 0 {
		builder = builder.Limit(uint64(limit))
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var users []models.User
	if err := sqlscan.Select(ctx, s.db, &users, sqlStr, args...); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *SearchService) BuildDynamicQuery(base squirrel.SelectBuilder, f SearchFilters) squirrel.SelectBuilder {
	q := base
	if f.Query != "" {
		term := "%" + f.Query + "%"
		q = q.Where(squirrel.Or{
			squirrel.ILike{"title": term},
			squirrel.ILike{"content": term},
		})
	}
	if f.UserID != nil {
		q = q.Where(squirrel.Eq{"user_id": *f.UserID})
	}
	if f.Published != nil {
		q = q.Where(squirrel.Eq{"published": *f.Published})
	}
	if f.MinWordCount != nil {
		q = q.Where("length(content) - length(replace(content, ' ', '')) + 1 >= ?", *f.MinWordCount)
	}
	return q
}
