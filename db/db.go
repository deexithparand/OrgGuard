package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// DB is the global database connection
var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	connStr := "host=localhost port=5432 user=admin password=admin@1234 dbname=orgguard sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	fmt.Println("Database connected successfully!")
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
		fmt.Println("Database connection closed!")
	}
}

// RunMigrations runs basic migrations
func RunMigrations() {
	query := `
	CREATE TABLE IF NOT EXISTS logs (
		id SERIAL PRIMARY KEY,
		action VARCHAR(50),
		org_name VARCHAR(255),
		sender VARCHAR(255),
		target_user VARCHAR(255), -- Renamed from 'user' to 'target_user'
		membership_state VARCHAR(50),
		membership_url TEXT,
		opa_validation_result TEXT, -- Stores 'passed' or 'failed'
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Renamed from 'time'
	);
	`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Database migration completed!")
}

// Save event logs to database
func SaveLog(action, orgName, sender, targetUser, membershipState, membershipURL, opaResult string) error {
	query := `
		INSERT INTO logs (action, org_name, sender, target_user, membership_state, membership_url, opa_validation_result)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := DB.Exec(query, action, orgName, sender, targetUser, membershipState, membershipURL, opaResult)
	return err
}
