package coment

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/Grafiters/archive/configs"
	database "github.com/Grafiters/archive/db"
	"github.com/Grafiters/archive/db/seed"

	"fmt"

	"github.com/spf13/cobra"
)

type CommandDefinition struct {
	Use   string
	Short string
	Run   func(cmd *cobra.Command, args []string)
}

var Commands = []CommandDefinition{
	{
		Use:   "db:create",
		Short: "Run database create",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running Create database...")
			database.Create()
			fmt.Println("Done creating database")
		},
	},
	{
		Use:   "db:drop",
		Short: "Run database drop",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running database migrations...")
			database.Drop()
			fmt.Println("Done Droping database")
		},
	},
	{
		Use:   "db:migrate",
		Short: "Run database migrations",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running database migrations...")
			err := database.Migrate("up", "")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Done database migrations...")
		},
	},
	{
		Use:   "db:rollback:all",
		Short: "rollback all table",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Are you sure you want to rollback all tables? (y/n): ")
			confirmation, _ := reader.ReadString('\n')
			confirmation = strings.ToLower(strings.TrimSpace(confirmation))

			if confirmation == "y" {
				err := database.Migrate("down_all", "")
				if err != nil {
					log.Fatalf("Error rolling back all migrations: %v", err)
				}
				fmt.Println("All migrations rolled back successfully!")
			} else {
				fmt.Println("Rollback operation canceled.")
			}
		},
	},
	{
		Use:   "db:rollback",
		Short: "Run single table",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("Are you sure you want to rollback all tables? (y/n): ")
			if len(args) < 1 {
				fmt.Printf("Rolling back migrations for table: %s\n", "")
				err := database.Migrate("down", "")
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Done database migrations...")
			} else {
				tableName := args[0]
				reader := bufio.NewReader(os.Stdin)
				confirmation, _ := reader.ReadString('\n')
				confirmation = strings.ToLower(strings.TrimSpace(confirmation))

				if confirmation == "y" {
					fmt.Printf("Rolling back migrations for table: %s\n", tableName)
					err := database.Migrate("down", tableName)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Done database migrations...")
				} else {
					fmt.Println("Rollback operation canceled.")
				}
			}
		},
	},
	{
		Use:   "db:seed",
		Short: "Seed the database with initial data",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Seeding the database...")
			if err := configs.Initialize(); err != nil {
				log.Fatal(err)
				return
			}
			seed.SeederData()
		},
	},
	{
		Use:   "generate:migration",
		Short: "Run generate migration table",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) >= 1 {
				fmt.Println("Running generator migrations...")
				err := database.GenerateMigration(args[0])
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Done database migrations...")
			} else {
				fmt.Println("!No generator running")
			}

		},
	},
}
