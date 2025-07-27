package main

import (
	"fmt"
	"os"

	"github.com/Grafiters/archive/task/coment"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "app"}
	for _, command := range coment.Commands {
		cobraCmd := &cobra.Command{
			Use:   command.Use,
			Short: command.Short,
			Run:   command.Run,
		}

		rootCmd.AddCommand(cobraCmd)
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Available tasks:")
			for _, c := range coment.Commands {
				fmt.Printf(" - %s: %s\n", c.Use, c.Short)
			}
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
