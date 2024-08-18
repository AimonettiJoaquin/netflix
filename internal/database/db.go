package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Connect(databaseUrl string) (*sql.DB, error) {
	db, err := sql.Open("mysql", databaseUrl)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Connected to the database")
	return db, nil
}

func CreateUsersTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Table Users created successfully")
	return nil
}

func CreateCommentsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id INT AUTO_INCREMENT PRIMARY KEY,
			id_user INT NOT NULL,
			id_movie INT NOT NULL,
			description VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Table Comments created successfully")
	return nil
}

func CreateMovieTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS movies (
			id INT PRIMARY KEY,
			counter INT
		)
	`)
	if err != nil {
		return err
	}
	log.Println("Table Movie created successfully")
	return nil
}
