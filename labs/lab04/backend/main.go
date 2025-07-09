package main

import (
	"fmt"
	"log"

	"lab04-backend/database"
	"lab04-backend/models"
	"lab04-backend/repository"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize database connection
	db, err := database.InitDB() // :contentReference[oaicite:0]{index=0}
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Run migrations (using goose-based approach)
	if err := database.RunMigrations(db); err != nil { // :contentReference[oaicite:1]{index=1}
		log.Fatal("Failed to run migrations:", err)
	}

	// Create repository instances
	userRepo := repository.NewUserRepository(db) // :contentReference[oaicite:2]{index=2}
	postRepo := repository.NewPostRepository(db) // :contentReference[oaicite:3]{index=3}

	fmt.Println("Database initialized successfully!")
	fmt.Printf("User repository: %T\n", userRepo)
	fmt.Printf("Post repository: %T\n", postRepo)

	// --- Demo data operations ---

	// 1) Create a user
	userReq := &models.CreateUserRequest{
		Name:  "Alice Example",
		Email: "alice@example.com",
	}
	createdUser, err := userRepo.Create(userReq)
	if err != nil {
		log.Fatal("Failed to create demo user:", err)
	}
	fmt.Printf("Created user: %+v\n", createdUser)

	// 2) Create a post for that user
	postReq := &models.CreatePostRequest{
		UserID:    createdUser.ID,
		Title:     "Hello, World!",
		Content:   "This is my very first post.",
		Published: true,
	}
	createdPost, err := postRepo.Create(postReq)
	if err != nil {
		log.Fatal("Failed to create demo post:", err)
	}
	fmt.Printf("Created post: %+v\n", createdPost)

	// 3) Fetch and print all users
	allUsers, err := userRepo.GetAll()
	if err != nil {
		log.Fatal("Failed to fetch users:", err)
	}
	fmt.Printf("All users: %+v\n", allUsers)

	// 4) Count and print total posts
	totalPosts, err := postRepo.Count()
	if err != nil {
		log.Fatal("Failed to count posts:", err)
	}
	fmt.Printf("Total posts in database: %d\n", totalPosts)
}
