package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "dev"

func init() {
	rootCmd.Version = Version
}

var rootCmd = &cobra.Command{
	Use:   "miflo",
	Short: "Miflo is a database migration manager for SQLite, PostgreSQL & Turso",
	Long:  "A simple database migration manger for SQLite, PostgreSQL & Turso. Miflo is a stand-alone tool that can be used with any project.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := cmd.Help()
			if err != nil {
				fmt.Println(err)
			}
			os.Exit(0)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
