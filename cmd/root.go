package cmd

import (
	"flag"
	"fmt"

	"github.com/gavsidhu/miflo/internal/cli"
)

var (
	showRootVersion bool
	showRootHelp    bool
)

var Version = "dev"

func init() {
	rootCmd.Flags.Usage = printRootHelp
	rootCmd.Flags.BoolVar(&showRootVersion, "v", false, "Show version information")
	rootCmd.Flags.BoolVar(&showRootHelp, "h", false, "Show help information")

}

var rootCmd = cli.Command{
	Name:        "miflo",
	Description: "miflow is a simple migration manager tool for SQLite and PostgreSQL",
	SubCommands: make(map[string]*cli.Command),
	Flags:       flag.NewFlagSet("miflo", flag.ExitOnError),
	Run: func(cmd *cli.Command, args []string) {

		if showRootVersion {
			fmt.Printf("v%s", Version)
			return
		}

		if showRootHelp || len(args) == 0 {
			cmd.Flags.Usage()
			return
		}
	},
}

func printRootHelp() {

	fmt.Println(rootCmd.Description)
	fmt.Println("\nUsage:\n  miflo [command]")
	fmt.Println("Available Commands:")
	for name, subCmd := range rootCmd.SubCommands {
		fmt.Printf("  %s\t\t%s\n", name, subCmd.Description)
	}
	fmt.Println("\nFlags:")
	rootCmd.Flags.PrintDefaults()
	fmt.Println("\nUse \"miflo [command] -h\" for more information about a command.")

}

func Execute() {
	rootCmd.Execute()
}
