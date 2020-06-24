package cmd

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "running test js",
	Long:  `running`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("cmd %s",cmd.Name())
		spew.Dump(cmd.Flags())
		return nil
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
	runCmd.Flags().SortFlags = false
	runCmd.Flags().AddFlagSet(runFlagSet())
}


func runFlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	flags.StringP("addr","a",":9191","master listen address")
	return flags
}
