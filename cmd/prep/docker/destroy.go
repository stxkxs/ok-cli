package docker

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stxkxs/ok-cli/aws/ecr"
	"github.com/stxkxs/ok-cli/logger"
)

var destroy = &cobra.Command{
	Use:   "destroy",
	Short: "destroy public or private docker images",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info().
			Strs("args", args).
			Bool("public", public).
			Bool("private", private).
			Msg("ok prep docker destroy")

		if public {
			decoded, _ := LoadPrepConf()
			client := ecr.NewPublicEcrClient(decoded.Public.Region)
			for _, i := range decoded.Public.Images {
				client.DestroyRepository(decoded.Account, i.Name)
			}
		}

		if private {
			decoded, _ := LoadPrepConf()
			client := ecr.NewPrivateEcrClient(decoded.Private.Region)
			for _, i := range decoded.Private.Images {
				client.DestroyRepository(decoded.Account, i.Name)
			}
		}
	},
}

func init() {
	viper.AutomaticEnv()
	err := viper.BindPFlags(create.Flags())
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("failed to bind prep docker destroy flags to viper")
		return
	}
}
