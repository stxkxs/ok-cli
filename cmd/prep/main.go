package prep

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stxkxs/ok-cli/cmd/prep/docker"
	"github.com/stxkxs/ok-cli/cmd/prep/helm"
	"github.com/stxkxs/ok-cli/logger"
)

var public bool
var private bool

var Cmd = &cobra.Command{
	Use:   "prep",
	Short: "prepare dependencies",
	Long:  `prepare ecr repositories, docker images, and helm charts`,
}

func init() {
	Cmd.AddCommand(docker.Cmd)
	Cmd.AddCommand(helm.Cmd)

	Cmd.PersistentFlags().BoolVar(&public, "public", false, "manages public docker images when true")
	Cmd.PersistentFlags().BoolVar(&private, "private", false, "manages private docker images when true")

	viper.AutomaticEnv()
	err := viper.BindPFlags(Cmd.Flags())
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("failed to bind prep flags to viper")
		return
	}
}
