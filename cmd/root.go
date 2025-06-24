package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stxkxs/ok-cli/cmd/prep"
	"github.com/stxkxs/ok-cli/cmd/tidy"
	"github.com/stxkxs/ok-cli/logger"
	"log"
	"os"
)

var debug bool

var file string
var environment string

var rootCmd = &cobra.Command{
	Use:  "ok",
	Long: `aws cli utility`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(logging)

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logs")
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "command conf file")
	rootCmd.PersistentFlags().StringVarP(&environment, "environment", "e", "prototype", "target environment")

	rootCmd.AddCommand(tidy.Cmd)
	rootCmd.AddCommand(prep.Cmd)
	rootCmd.AddCommand(whoami)

	viper.AutomaticEnv()
	err := viper.BindPFlags(rootCmd.Flags())
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("failed to bind global flags to viper")
		return
	}
}

func logging() {
	if err := logger.Setup(debug); err != nil {
		log.Fatalf("failed to set up logger: %v", err)
	}
}
