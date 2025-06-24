package tidy

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stxkxs/ok-cli/aws"
	"github.com/stxkxs/ok-cli/env"
	"github.com/stxkxs/ok-cli/logger"
)

type Tidy struct {
	CodeBuild      aws.CodeBuild      `mapstructure:"codebuild"`
	CloudWatch     aws.CloudWatch     `mapstructure:"cloudwatch"`
	CloudFormation aws.CloudFormation `mapstructure:"cloudformation"`
}

var file string

var Cmd = &cobra.Command{
	Use:   "tidy",
	Short: "aws resource cleanup",
	Long:  `removes codebuild build history, cloudwatch log groups, and cloudformation stacks`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		file, err = cmd.Flags().GetString("file")
		if err != nil {
			logger.Logger.Error().
				Err(err).
				Msg("error fetching file flag")
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Logger.Debug().
			Strs("args", args).
			Msg("ok tidy")

		c := LoadTidyConf()

		err := aws.DestroyBuildHistory(c.CodeBuild)
		if err != nil {
			logger.Logger.Error().
				Err(err).
				Msg("error destroying codebuild history")
			return
		}

		err = aws.DestroyLogGroups(c.CloudWatch)
		if err != nil {
			logger.Logger.Error().
				Err(err).
				Msg("error destroying cloudwatch logs")
			return
		}

		err = aws.DestroyStacks(c.CloudFormation)
		if err != nil {
			logger.Logger.Error().
				Err(err).
				Msg("error destroying cloudformation stacks")
			return
		}
	},
}

func LoadTidyConf() *Tidy {
	c, err := env.Decode[Tidy](file, ".ok.tidy")
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("error decoding tidy conf")
		return nil
	}

	logger.Logger.Debug().
		Interface("decoded", c).
		Msg("decoded tidy conf")

	return &c
}

func init() {
	err := viper.BindPFlags(Cmd.Flags())
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("failed to bind tidy flags to viper")
		return
	}
}
