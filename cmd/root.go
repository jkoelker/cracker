//
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jkoelker/cracker/hunter2"
)

var target  string
var list    string
var workers int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cracker [target] [list]",
	Short: "Bro, its like the gibson and thingies",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return hunter2.Search(args[0], args[1], workers)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().IntVar(&workers, "workers", 20, "number of workers")
}
