package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	"github.com/stxkxs/ok-cli/logger"
	"golang.org/x/time/rate"
	"strings"
)

type CodeBuild struct {
	Region    string   `mapstructure:"region"`
	BatchSize int      `mapstructure:"batchSize"`
	Prefix    []string `mapstructure:"prefix"`
	Retry     int      `mapstructure:"retry"`
}

var codeBuildRateLimit = rate.NewLimiter(rate.Limit(10), 10)

func DestroyBuildHistory(c CodeBuild) error {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(c.Region), config.WithRetryMaxAttempts(c.Retry))
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error loading default aws configurations")
		return err
	}

	cb := codebuild.NewFromConfig(cfg)
	input := &codebuild.ListBuildsInput{}

	err = codeBuildRateLimit.Wait(ctx)
	if err != nil {
		return err
	}

	builds, err := getAllBuilds(ctx, cb, input, c.Prefix)
	if err != nil {
		return err
	}

	err, done := maybeDeleteBuilds(c, builds, cb, ctx)
	if done {
		return err
	}

	return nil
}

func maybeDeleteBuilds(c CodeBuild, builds []string, cb *codebuild.Client, ctx context.Context) (error, bool) {
	if len(builds) == 0 {
		logger.Logger.Warn().Str("region", c.Region).Msg("no builds found in region")
		return nil, true
	}

	for i := 0; i < len(builds); i += c.BatchSize {
		end := i + c.BatchSize
		if end > len(builds) {
			end = len(builds)
		}

		batch := builds[i:end]
		logger.Logger.Debug().Strs("ids", batch).Msg("deleting builds")

		deleted, err := cb.BatchDeleteBuilds(ctx, &codebuild.BatchDeleteBuildsInput{Ids: batch})
		if err != nil {
			return err, true
		}

		logger.Logger.Info().
			Strs("ids", batch).
			Interface("deleted", deleted).
			Msg("deleted builds")
	}

	return nil, false
}

func getAllBuilds(ctx context.Context, cb *codebuild.Client, input *codebuild.ListBuildsInput, prefixes []string) ([]string, error) {
	var builds []string

	for {
		found, err := cb.ListBuilds(ctx, input)
		if err != nil {
			logger.Logger.Error().Err(err).Msg("error listing codebuild builds")
			return nil, err
		}

		for _, id := range found.Ids {
			for _, prefix := range prefixes {
				if strings.HasPrefix(id, prefix) {
					builds = append(builds, id)
				}
			}
		}

		if found.NextToken == nil {
			break
		}

		input.NextToken = found.NextToken
	}

	return builds, nil
}
