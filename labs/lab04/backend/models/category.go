package models

import (
	"errors"
	"log"
	"regexp"
	"time"

	"gorm.io/gorm"
)

// Category represents a blog post category using GORM model conventions
// This model demonstrates GORM ORM patterns and relationships
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null;uniqueIndex"`
	Description string         `json:"description" gorm:"size:500"`
	Color       string         `json:"color" gorm:"size:7"` // Hex color code
	Active      bool           `json:"active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// GORM Associations (demonstrates ORM relationships)
	Posts []Post `json:"posts,omitempty" gorm:"many2many:post_categories;"`
}

// CreateCategoryRequest represents the payload for creating a category
type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// UpdateCategoryRequest represents the payload for updating a category
type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
	Active      *bool   `json:"active,omitempty"`
}

// BeforeCreate hook: validate and set defaults
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if len(c.Name) < 2 {
		return errors.New("category name must be at least 2 characters")
	}
	if c.Color == "" {
		c.Color = "#007bff"
	} else {
		matched, _ := regexp.MatchString(`^#([A-Fa-f0-9]{6})$`, c.Color)
		if !matched {
			return errors.New("color must be a valid hex code, e.g. #RRGGBB")
		}
	}
	return nil
}

// AfterCreate hook: log creation
func (c *Category) AfterCreate(tx *gorm.DB) error {
	log.Printf("Category created: %s (ID: %d)", c.Name, c.ID)
	return nil
}

// BeforeUpdate hook: validate updates
func (c *Category) BeforeUpdate(tx *gorm.DB) error {
	if len(c.Name) > 0 && len(c.Name) < 2 {
		return errors.New("category name must be at least 2 characters")
	}
	if c.Color != "" {
		matched, _ := regexp.MatchString(`^#([A-Fa-f0-9]{6})$`, c.Color)
		if !matched {
			return errors.New("color must be a valid hex code, e.g. #RRGGBB")
		}
	}
	return nil
}

// Validate ensures CreateCategoryRequest has valid fields
func (req *CreateCategoryRequest) Validate() error {
	if len(req.Name) < 2 || len(req.Name) > 100 {
		return errors.New("name must be between 2 and 100 characters")
	}
	if len(req.Description) > 500 {
		return errors.New("description must be at most 500 characters")
	}
	if req.Color != "" {
		matched, _ := regexp.MatchString(`^#([A-Fa-f0-9]{6})$`, req.Color)
		if !matched {
			return errors.New("color must be a valid hex code, e.g. #RRGGBB")
		}
	}
	return nil
}

// ToCategory converts CreateCategoryRequest to Category model
func (req *CreateCategoryRequest) ToCategory() *Category {
	c := &Category{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		Active:      true,
	}
	if c.Color == "" {
		c.Color = "#007bff"
	}
	return c
}

// Validate ensures UpdateCategoryRequest has valid fields
func (req *UpdateCategoryRequest) Validate() error {
	if req.Name != nil {
		if len(*req.Name) < 2 || len(*req.Name) > 100 {
			return errors.New("name must be between 2 and 100 characters")
		}
	}
	if req.Description != nil {
		if len(*req.Description) > 500 {
			return errors.New("description must be at most 500 characters")
		}
	}
	if req.Color != nil {
		matched, _ := regexp.MatchString(`^#([A-Fa-f0-9]{6})$`, *req.Color)
		if !matched {
			return errors.New("color must be a valid hex code, e.g. #RRGGBB")
		}
	}
	return nil
}

// GORM scope: only active categories
func ActiveCategories(db *gorm.DB) *gorm.DB {
	return db.Where("active = ?", true)
}

// GORM scope: categories that have at least one post
func CategoriesWithPosts(db *gorm.DB) *gorm.DB {
	return db.Joins("JOIN post_categories ON post_categories.category_id = categories.id").
		Joins("JOIN posts ON posts.id = post_categories.post_id").
		Group("categories.id")
}

// Check if category is active
func (c *Category) IsActive() bool {
	return c.Active
}

// Get count of posts in this category
func (c *Category) PostCount(db *gorm.DB) int64 {
	return db.Model(c).Association("Posts").Count()
}
