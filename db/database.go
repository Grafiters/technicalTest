package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Grafiters/archive/db/config"
)

func Create() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func Drop() {
	db, err := config.ConfigDBBase()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	error := config.DeleteDatabase(db)
	if error != nil {
		log.Fatal(error)
	}

	fmt.Println("DONE Delete Database")
}

func Migrate(direction string, table string) error {
	db, err := config.ConfigDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = createSchemaMigrations(db)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	path := "db/migration"
	if table != "" {
		path = fmt.Sprintf("db/migration/%s", table)
	}

	lastVersion, err := getLastAppliedVersion(db)
	if err != nil {
		log.Fatal(err)
	}

	if direction == "up" {
		migrationFiles, err := getMigrationFiles(path, ".up.sql")
		if err != nil {
			log.Fatal(err)
		}

		// Apply pending migrations
		for _, file := range migrationFiles {
			version := extractVersion(file)
			if version > lastVersion {
				if err := applyMigration(db, file); err != nil {
					log.Fatalf("Error applying migration %s: %v", file, err)
				}
				fmt.Printf("Applied migration: %s\n", file)
			}
		}

		fmt.Println("All pending migrations applied successfully!")
	} else if direction == "down_all" {
		migrationFiles, err := getMigrationFiles(path, ".down.sql")
		if err != nil {
			log.Fatal(err)
		}

		reverseArrayString(migrationFiles)

		for _, file := range migrationFiles {
			if err := RemoveMigration(db, file); err != nil {
				log.Fatalf("Error rolling back migration %s: %v", file, err)
			}

			if _, err := db.Exec("DELETE FROM schema_migrations WHERE version = $1", lastVersion); err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Rolled back migration: %d\n", lastVersion)
		}
	} else {
		var rollbackProcess bool = true
		var count int
		migrationFiles, err := getMigrationFiles(path, ".down.sql")
		if err != nil {
			log.Fatal(err)
		}

		reverseArrayString(migrationFiles)

		for _, file := range migrationFiles {
			lastVersion := extractVersion(file)
			if _, err := db.Exec("SELECT * FROM schema_migrations WHERE version = $1", lastVersion); err != nil {
				fmt.Println(err)
			}

			row := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", lastVersion)
			if err := row.Scan(&count); err != nil {
				log.Fatal(err)
			}

			if count >= 1 && rollbackProcess {
				if err := RemoveMigration(db, file); err != nil {
					log.Fatalf("Error rolling back migration %s: %v", file, err)
				}

				if _, err := db.Exec("DELETE FROM schema_migrations WHERE version = $1", lastVersion); err != nil {
					log.Fatal(err)
				}

				rollbackProcess = false

				fmt.Printf("Rolled back migration: %d\n", lastVersion)
			}
		}
	}

	return nil
}

func GenerateMigration(name string) error {
	path := "db/migration"
	migrateCmd := exec.Command("migrate", "create", "-ext", "sql", "-dir", path, fmt.Sprintf("%s", name+"_table"))

	migrateCmd.Stdout = os.Stdout
	migrateCmd.Stderr = os.Stderr
	if err := migrateCmd.Run(); err != nil {
		log.Println(err)
		log.Fatalf("Failed to run migrations: %v", err)
		return err
	}

	return nil
}

func createSchemaMigrations(db *sql.DB) error {
	var exists bool

	query := `SELECT EXISTS (
        SELECT FROM information_schema.tables 
        WHERE  table_schema = 'public'
        AND    table_name   = 'schema_migrations'
    );`

	db.QueryRow(query).Scan(&exists)
	if !exists {
		// Create the schema_migrations table
		createTableQuery := `
        CREATE TABLE schema_migrations (
            version bigint PRIMARY KEY,
            dirty boolean NULL
        );`

		_, err := db.Exec(createTableQuery)
		if err != nil {
			log.Fatalf("Unable to create schema_migrations table: %v", err)
			return err
		}
		fmt.Println("schema_migrations table created successfully.")
	}

	return nil
}

func getLastAppliedVersion(db *sql.DB) (int, error) {
	var version int
	row := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations")
	if err := row.Scan(&version); err != nil {
		return 0, err
	}

	return version, nil
}

func getMigrationFiles(dir, ext string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ext) {
			migrationFiles = append(migrationFiles, filepath.Join(dir, file.Name()))
		}
	}

	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

func extractVersion(filename string) int {
	base := filepath.Base(filename)
	parts := strings.Split(base, "_")
	var version int
	fmt.Sscanf(parts[0], "%d", &version)
	return version
}

func applyMigration(db *sql.DB, filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(content)); err != nil {
		return err
	}

	version := extractVersion(filename)
	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		return err
	}

	return tx.Commit()
}

func RemoveMigration(db *sql.DB, filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	tx.Exec(string(content))

	version := extractVersion(filename)
	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", version); err != nil {
		return err
	}

	return tx.Commit()
}

func reverseArrayString(data []string) []string {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}

	return data
}
