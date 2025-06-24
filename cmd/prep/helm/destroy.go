package helm

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stxkxs/ok-cli/aws/ecr"
	"github.com/stxkxs/ok-cli/logger"
)

var destroy = &cobra.Command{
	Use:   "destroy",
	Short: "destroy public or private helm charts",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info().
			Strs("args", args).
			Bool("public", public).
			Bool("private", private).
			Msg("ok prep helm destroy")

		if public {
			decoded, _ := LoadPrepConf()
			client := ecr.NewPublicEcrClient(decoded.Public.Region)
			for _, i := range decoded.Public.Charts {
				client.DestroyRepository(decoded.Account, i.Name)
			}
		}

		if private {
			decoded, _ := LoadPrepConf()
			client := ecr.NewPrivateEcrClient(decoded.Private.Region)
			for _, i := range decoded.Private.Charts {
				client.DestroyRepository(decoded.Account, i.Name)
			}
		}
	},
}

func init() {
	viper.AutomaticEnv()
	err := viper.BindPFlags(Cmd.Flags())
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("failed to bind prep helm destroy flags to viper")
		return
	}
}
