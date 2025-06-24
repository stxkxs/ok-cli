package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/stxkxs/ok-cli/logger"
	"golang.org/x/time/rate"
)

type CloudWatch struct {
	Region    string `mapstructure:"region"`
	BatchSize int    `mapstructure:"batchSize"`
	Retry     int    `mapstructure:"retry"`
}

var cloudWatchRateLimit = rate.NewLimiter(rate.Limit(10), 10)

func DestroyLogGroups(c CloudWatch) error {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(c.Region), config.WithRetryMaxAttempts(c.Retry))
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error loading default aws configurations")
		return err
	}

	cwl := cloudwatchlogs.NewFromConfig(cfg)

	err = cloudWatchRateLimit.Wait(ctx)
	if err != nil {
		return err
	}

	var logGroupNames []string
	err = getAllLogGroupNames(ctx, cwl, &logGroupNames)
	if err != nil {
		return err
	}

	if len(logGroupNames) == 0 {
		logger.Logger.Warn().Str("region", c.Region).Msg("no log groups found in region")
		return nil
	}

	for i := 0; i < len(logGroupNames); i += c.BatchSize {
		end := i + c.BatchSize
		if end > len(logGroupNames) {
			end = len(logGroupNames)
		}

		batch := logGroupNames[i:end]
		logger.Logger.Debug().Strs("logGroups", batch).Msg("destroying log groups")

		for _, logGroupName := range batch {
			_, err := cwl.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
				LogGroupName: &logGroupName,
			})
			if err != nil {
				return err
			}
			logger.Logger.Info().Str("logGroup", logGroupName).Msg("destroyed log group")
		}
	}

	return nil
}

func getAllLogGroupNames(ctx context.Context, cwl *cloudwatchlogs.Client, logGroupNames *[]string) error {
	input := &cloudwatchlogs.DescribeLogGroupsInput{}

	for {
		resp, err := cwl.DescribeLogGroups(ctx, input)
		if err != nil {
			logger.Logger.Error().Err(err).Msg("error describing log groups")
			return err
		}

		for _, lg := range resp.LogGroups {
			*logGroupNames = append(*logGroupNames, *lg.LogGroupName)
		}

		if resp.NextToken == nil {
			break
		}
		input.NextToken = resp.NextToken
	}

	return nil
}
