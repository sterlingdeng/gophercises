package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"task/db"
)

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "a cli based task manager",
}

var (
	store db.DB
	err   error
)

func init() {
	store, err = db.NewDB("my.db")
	if err != nil {
		log.Fatalf("failed to initialize db. %v", err)
	}
	_ = fmt.Sprint(store)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
