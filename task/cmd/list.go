package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list",
	Long:  "list all tasks",
	Run:   list,
}

func list(cmd *cobra.Command, args []string) {
	tasks, err := store.GetList()
	if err != nil {
		fmt.Printf("failed to get tasks: %v\n", err)
		return
	}

	fmt.Println("Listing tasks below...")
	for i, task := range tasks {
		fmt.Printf("%d. %s\n", i + 1, task.Value)
	}
}
