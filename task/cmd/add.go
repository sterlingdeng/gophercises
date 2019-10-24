package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a task",
	Run:   add,
}

func add(cmd *cobra.Command, args []string) {
	err := store.AddTasks(args...)
	if err != nil {
		fmt.Printf("failed to add tasks: %v\n", err)
	}
}
