package main

import (
	"database/sql"
	"log"

	"github.com/YahyaBekirCanevi/message_ai_go_htmx/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"
)

func main() {
	err := initializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
		return
	}
}

func initializeDatabase() error {
	// Define the database file name
	dbFileName := "database.sqlite"

	// Open the database connection.
	log.Printf("Attempting to open database: %s", dbFileName)
	db, err := sql.Open("sqlite", "file:"+dbFileName+"?cache=shared&mode=rwc")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
		log.Println("Database connection closed.")
	}()

	// Ping the database to ensure the connection is valid
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return err
	}
	log.Println("Successfully connected to SQLite database.")

	// Run migrations to create tables
	log.Println("Running database migrations...")
	if err = MigrateDB(db); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	log.Println("Database migrations completed successfully.")

	log.Println("Starting server...")
	startServer(db)
	return nil
}

func startServer(db *sql.DB) {
	r := gin.Default(func(e *gin.Engine) {
		e.Use(gin.Recovery())
		e.Use(gin.Logger())
		// Rate limiting middleware is not available in gin by default; consider using a third-party package or implement your own if needed.
	})

	r.GET("/", handlers.ListConversations(db))

	r.POST("/chat/new", handlers.NewChatForm())
	r.POST("/chat/start", handlers.StartChat(db))
	r.GET("/chat/:id", handlers.LoadChat(db))

	r.POST("/message/send", handlers.SendMessage(db))

	r.Run(":8080")
	log.Println("Server is running on port 8080")
}
