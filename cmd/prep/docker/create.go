package docker

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
	Short: "create or update public or private docker images",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info().
			Strs("args", args).
			Bool("public", public).
			Bool("private", private).
			Msg("ok prep docker create")

		if public {
			decoded, _ := LoadPrepConf()
			client := ecr.NewPublicEcrClient(decoded.Public.Region)
			client.MaybeCreateRepository(decoded.Account, client.ConvertDockerImagesToRepositories(decoded.Public.Images))
			for _, r := range decoded.Public.Images {
				client.CreateUpdateDockerImage(r.Alias, decoded.Public.Region, r)
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
			client.MaybeCreateRepository(decoded.Account, client.ConvertDockerImagesToRepositories(decoded.Private.Images))
			for _, r := range decoded.Private.Images {
				client.CreateUpdateDockerImage(decoded.Account, decoded.Private.Region, r)
			}
		}
	},
}

func init() {
	err := viper.BindPFlags(create.Flags())
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("failed to bind prep docker create flags to viper")
		return
	}
}
