package helm

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stxkxs/ok-cli/aws/ecr"
	"github.com/stxkxs/ok-cli/logger"
	"github.com/stxkxs/ok-cli/terminal"
	"strings"
	"time"
)

var create = &cobra.Command{
	Use:   "create",
	Short: "create or update public or private helm charts",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info().
			Strs("args", args).
			Bool("public", public).
			Bool("private", private).
			Msg("ok prep helm create")

		if public {
			decoded, _ := LoadPrepConf()
			client := ecr.NewPublicEcrClient(decoded.Public.Region)
			client.MaybeCreateRepository(decoded.Account, client.ConvertHelmChartsToRepositories(decoded.Public.Charts))
			for _, r := range decoded.Public.Charts {
				client.CreateUpdateHelmChart(r.Alias, decoded.Public.Region, r)
			}

			b := fmt.Sprintf("aws ecr-public put-registry-catalog-data --region %s --display-name \"%s\"", decoded.Public.Region, decoded.Name)
			c := strings.Join(strings.Fields(b), " ")

			logger.Logger.Info().
				Str("command", b).
				Msg("put public registry name")

			err := terminal.ExecuteCommand(c, time.Duration(5)*time.Minute)
			if err != nil {
				logger.Logger.Error().Msg("error naming public registry")
			}
		}

		if private {
			decoded, _ := LoadPrepConf()
			client := ecr.NewPrivateEcrClient(decoded.Private.Region)
			client.MaybeCreateRepository(decoded.Account, client.ConvertHelmChartsToRepositories(decoded.Private.Charts))
			for _, r := range decoded.Private.Charts {
				client.CreateUpdateHelmChart(decoded.Account, decoded.Private.Region, r)
			}
		}
	},
}

func init() {
	err := viper.BindPFlags(create.Flags())
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("failed to bind prep helm create flags to viper")
		return
	}
}
