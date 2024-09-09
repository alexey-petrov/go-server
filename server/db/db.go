// db/db.go
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/alexey-petrov/go-server/server/structs"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB // Global DB variable to be used across the application

// ConnectDB initializes the connection to the PostgreSQL database
func ConnectDB() *sql.DB {
	// Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Retrieve environment variables
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    fmt.Printf("Connecting to database '%s' as user '%s' with password '%s'.\n", dbName, dbUser, dbPassword)
	// Connection string (replace with your actual PostgreSQL credentials)
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)

	// Open the connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open a DB connection:", err)
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping DB:", err)
	}

	fmt.Println("Connected to the database successfully")
	DB = db
	return DB
}

// CloseDB closes the DB connection
func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Fatal("Failed to close DB connection:", err)
		}
	}
}

func CreateTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		firstName TEXT NOT NULL,
		lastName TEXT NOT NULL,
		password TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table created successfully!")
}

func InsertUser(user structs.User) structs.User {
	ConnectDB()
	query := `INSERT INTO users (email, firstName, lastName, password) VALUES ($1, $2, $3, $4) RETURNING id`
	// Insert the user into the database
	_, err := DB.Exec(query, user.Email, user.FirstName, user.LastName, user.Password)
	if err != nil {

		log.Fatal(err)
	}

	// Get the ID of the newly inserted user
	var id int
	err = DB.QueryRow("SELECT id FROM users ORDER BY id DESC LIMIT 1").Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the user data from the database
	userData := structs.User{}
	err = DB.QueryRow("SELECT id, email, firstName, lastName FROM users WHERE id = $1", id).Scan(&userData.ID, &userData.Email, &userData.FirstName, &userData.LastName)
	if err != nil {
		log.Fatal(err)
	}

	defer CloseDB()

	// Return the user data
	return userData

}

func GetUsers(db *sql.DB) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email string
		err = rows.Scan(&id, &name, &email)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("User: %d | Name: %s | Email: %s\n", id, name, email)
	}

	// Check for errors from iterating over rows.
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}