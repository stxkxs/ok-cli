package helm

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stxkxs/ok-cli/aws/ecr"
	"github.com/stxkxs/ok-cli/env"
	"github.com/stxkxs/ok-cli/logger"
)

var file string
var environment string
var public bool
var private bool

var Cmd = &cobra.Command{
	Use:   "helm",
	Short: "create, update, or destroy public or private helm charts",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		file, err = cmd.Flags().GetString("file")
		if err != nil {
			logger.Logger.Error().
				Err(err).
				Msg("error fetching file flag")
			return
		}

		environment, err = cmd.Flags().GetString("environment")
		if err != nil {
			logger.Logger.Error().
				Err(err).
				Msg("error fetching environment flag")
			return
		}

		public, err = cmd.Flags().GetBool("public")
		if err != nil {
			logger.Logger.Error().
				Err(err).
				Msg("error fetching public flag")
			return
		}

		private, err = cmd.Flags().GetBool("private")
		if err != nil {
			logger.Logger.Error().
				Err(err).
				Msg("error fetching private flag")
			return
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !public && !private {
			logger.Logger.Error().
				Msg("choose public or private docker operations")
			return fmt.Errorf("choose public or private docker operations")
		}

		if public && private {
			logger.Logger.Error().
				Msg("only public or private can be used at a time, not both")
			return fmt.Errorf("only public or private can be used at a time, not both")
		}

		return nil
	},
}

func LoadPrepConf() (ecr.Prep, error) {
	conf, err := env.Decode[ecr.Prep](file, fmt.Sprintf(".ok.prep.%s", environment))
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("failed to decode prep conf")
		return conf, err
	}

	logger.Logger.Debug().
		Interface("decoded", conf).
		Msg("decoded prep conf")

	return conf, nil
}

func init() {
	Cmd.AddCommand(create)
	Cmd.AddCommand(destroy)

	err := viper.BindPFlags(Cmd.Flags())
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("failed to bind prep helm flags to viper")
		return
	}
}
