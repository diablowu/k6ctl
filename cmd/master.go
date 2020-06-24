package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k6ctl/internal/admin"
)






var masterCmd = &cobra.Command{
	Use:   "master",
	Short: "Start master node",
	Long:  `Start master node`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logrus.Debug("cmd: {}",cmd.Name())

		admin.Staring(cmd.Flag("addr").Value.String())
		return nil
	},
}

func init() {
	RootCmd.AddCommand(masterCmd)
	masterCmd.Flags().SortFlags = false
	masterCmd.Flags().AddFlagSet(masterFlagSet())
}


func masterFlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	flags.StringP("addr","a",":9191","master listen address")
	return flags
}


