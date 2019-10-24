package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

func init() {
	rootCmd.AddCommand(doCmd)
}

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "marks a task complete",
	Run:   do,
}

func do(cmd *cobra.Command, args []string) {
	if len(args) > 1 {
		fmt.Println("can only delete 1 item at a time for now")
		return
	}
	idx, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("please provide an integer input")
		return
	}
	err = store.Delete(idx)
	if err != nil {
		fmt.Println("something went wrong")
	}
}
