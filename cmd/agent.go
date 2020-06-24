package cmd

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Start agent node",
	Long:  `Start agent node`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("cmd %s",cmd.Name())
		spew.Dump(cmd.Flags())
		return nil
	},
}

func agentFlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	flags.Int64P("vus", "u", 1, "number of virtual users")
	flags.Int64P("max", "m", 0, "max available virtual users")
	flags.DurationP("duration", "d", 0, "test duration limit")
	flags.Int64P("iterations", "i", 0, "script total iteration limit (among all VUs)")
	flags.StringSliceP("stage", "s", nil, "add a `stage`, as `[duration]:[target]`")
	flags.BoolP("paused", "p", false, "start the test in a paused state")
	return flags
}

func init() {
	RootCmd.AddCommand(agentCmd)
	agentCmd.Flags().SortFlags = false
	agentCmd.Flags().AddFlagSet(agentFlagSet())

}
