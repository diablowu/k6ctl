package cmd

import (
	"github.com/loadimpact/k6/lib/consts"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)




var RootCmd = &cobra.Command{
	Use:           "k6ctl",
	Short:         "k6 controller",
	Long:          "k6ctl",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		//setupLoggers(logFmt)
		//if noColor {
		//	stdout.Writer = colorable.NewNonColorable(os.Stdout)
		//	stderr.Writer = colorable.NewNonColorable(os.Stderr)
		//}
		//log.SetOutput(logrus.StandardLogger().Writer())
		logrus.Debugf("k6 version: v%s", consts.FullVersion())
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		code := -1
		var logger logrus.FieldLogger = logrus.StandardLogger()
		//if e, ok := err.(ExitCode); ok {
		//	code = e.Code
		//	if e.Hint != "" {
		//		logger = logger.WithField("hint", e.Hint)
		//	}
		//}
		logger.Error(err)
		os.Exit(code)
	}
}