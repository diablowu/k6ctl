package main

import (
	"github.com/sirupsen/logrus"
	"k6ctl/cmd"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	cmd.Execute()
}
