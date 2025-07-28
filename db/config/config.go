package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	db            *sql.DB
	DnsConnection string
)

func ConfigDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbName, host, port)

	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		errorText := err
		fmt.Println(errorText)
		return nil, err
	}

	db := conn
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	DnsConnection = fmt.Sprintf("%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	return db, nil
}

func InitDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	fmt.Println("====== Checking connection DB =======")
	connectionString := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable", user, password, host, port)

	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	fmt.Println("====== Succes connection DB =======")

	db := conn
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("====== Create INIT DATABASE =======")
	createDatabase(db, dbName)
	fmt.Println("====== SUCCESS INIT DATABASE =======")
	return db, nil
}

func ConfigDBBase() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")

	connectionString := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable", user, password, host, port)

	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		errorText := err
		fmt.Println(errorText)
		return nil, err
	}

	db := conn
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDatabase(db *sql.DB, dbName string) {
	createDBQuery := fmt.Sprintf("CREATE DATABASE %s", dbName)

	_, err := db.Exec(createDBQuery)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Database '%s' created successfully\n", dbName)
}

func DatabaseExists() (bool, error) {
	dbName := os.Getenv("DB_NAME")
	rows, err := db.Query(`SELECT 1 FROM pg_database WHERE datname = $1`, dbName)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func DeleteDatabase(db *sql.DB) error {
	dbName := os.Getenv("DB_NAME")
	dropDBQuery := fmt.Sprintf("DROP DATABASE %s", dbName)
	_, err := db.Exec(dropDBQuery)
	fmt.Println(err)
	return err
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}
func GetDB() *sql.DB {
	return db
}
