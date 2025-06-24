package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stxkxs/ok-cli/env"
	"github.com/stxkxs/ok-cli/logger"
)

var account string
var version string
var organization string
var name string
var alias string

var whoami = &cobra.Command{
	Use:   "whoami",
	Short: "base32 10 character id",
	Long:  "generates 10-character base32 id from account, region, environment, organization, name, and alias inputs",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		environment, err = cmd.Flags().GetString("environment")
		if err != nil {
			logger.Logger.Error().
				Err(err).
				Msg("error fetching environment flag")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info().
			Strs("args", args).
			Str("account", account).
			Str("environment", environment).
			Str("organization", organization).
			Str("name", name).
			Str("alias", alias).
			Msg("ok whoami")

		logger.Logger.Info().
			Msg(env.Id(account + environment + version + organization + name + alias))
	},
}

func init() {
	whoami.PersistentFlags().StringVar(&account, "account", "", "target account")
	whoami.PersistentFlags().StringVar(&version, "version", "", "target version")
	whoami.PersistentFlags().StringVar(&organization, "organization", "", "target organization")
	whoami.PersistentFlags().StringVar(&name, "name", "", "target name")
	whoami.PersistentFlags().StringVar(&alias, "alias", "", "target alias")

	viper.AutomaticEnv()
	err := viper.BindPFlags(whoami.Flags())
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("failed to bind whoami flags to viper")
		return
	}
}
