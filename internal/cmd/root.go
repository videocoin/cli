package cmd

import (
	"github.com/VideoCoin/cli/internal/config"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "VideoCoin Network Testnet Client",
	Long:  "VideoCoin Network Testnet Client",
}

var c config.Config

func Execute(b, v string) {
	Build = b
	Version = v
	rootCmd.AddCommand(cmdVersion)

	err := envconfig.Process("cli", &c)
	if err != nil {
		logrus.Fatal(err)
	}

	err = c.InitLogger(b, v)
	if err != nil {
		logrus.Fatal(err)
	}

	cmdStart.Flags().StringP("account", "a", "", "account file path")
	err = cmdStart.MarkFlagRequired("account")
	if err != nil {
		logrus.WithError(err).Panic()
	}

	cmdStart.Flags().StringP("password", "p", "", "private key password")

	rootCmd.AddCommand(cmdStart)

	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Panic()
	}
}
