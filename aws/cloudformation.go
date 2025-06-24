package aws

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/stxkxs/ok-cli/logger"
)

type CloudFormation struct {
	Region string   `mapstructure:"region"`
	Retry  int      `mapstructure:"retry"`
	Prefix []string `mapstructure:"prefix"`
}

var cloudFormationRateLimit = rate.NewLimiter(rate.Limit(10), 10)

func DestroyStacks(c CloudFormation) error {
	ctx := context.Background()

	err, api := NewClient(c, ctx)
	if err != nil {
		return err
	}

	err = cloudFormationRateLimit.Wait(ctx)
	if err != nil {
		return err
	}

	var stackNames []string
	err = getAllStackNames(ctx, api, &stackNames)
	if err != nil {
		return err
	}

	if len(stackNames) == 0 {
		logger.Logger.Warn().Str("region", c.Region).Msg("no stacks found in region")
		return nil
	}

	for _, stackName := range stackNames {
		if !toolkit(stackName) && matchesPrefix(stackName, c.Prefix) {
			err := deleteStackAndWait(ctx, api, stackName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func NewClient(c CloudFormation, ctx context.Context) (error, *cloudformation.Client) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(c.Region), config.WithRetryMaxAttempts(c.Retry))
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error loading default aws configurations")
		return err, nil
	}

	api := cloudformation.NewFromConfig(cfg)
	return err, api
}

func deleteStackAndWait(ctx context.Context, api *cloudformation.Client, stackName string) error {
	_, err := api.DeleteStack(ctx, &cloudformation.DeleteStackInput{
		StackName: &stackName,
	})
	if err != nil {
		logger.Logger.Error().Err(err).Str("stack", stackName).Msg("error deleting stack")
		return err
	}
	logger.Logger.Info().Str("stack", stackName).Msg("initiated deletion of stack")

	for {
		status, err := getStackStatus(ctx, api, stackName)
		if err != nil {
			return err
		}

		if status == "DELETE_COMPLETE" {
			logger.Logger.Info().Str("stack", stackName).Msg("deleted stack")
			return nil
		} else if status == "DELETE_FAILED" || status == "DELETE_IN_PROGRESS" {
			logger.Logger.Info().Str("stack", stackName).Str("status", status).Msg("waiting for stack deletion")
			time.Sleep(20 * time.Second) // Wait for 20 seconds before checking again
		} else {
			return fmt.Errorf("unexpected stack status: %s", status)
		}
	}
}

func getStackStatus(ctx context.Context, api *cloudformation.Client, stackName string) (string, error) {
	resp, err := api.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: &stackName,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Stacks) == 0 {
		return "", fmt.Errorf("stack not found: %s", stackName)
	}

	return string(resp.Stacks[0].StackStatus), nil
}

func getAllStackNames(ctx context.Context, api *cloudformation.Client, stackNames *[]string) error {
	input := &cloudformation.DescribeStacksInput{}

	for {
		re, err := api.DescribeStacks(ctx, input)
		if err != nil {
			logger.Logger.Error().Err(err).Msg("error listing stacks")
			return err
		}

		for _, stack := range re.Stacks {
			if stack.StackName != nil {
				*stackNames = append(*stackNames, *stack.StackName)
			}
		}

		if re.NextToken == nil {
			break
		}
		input.NextToken = re.NextToken
	}

	return nil
}

func matchesPrefix(stackName string, prefixes []string) bool {
	if len(prefixes) == 0 {
		return true
	}

	for _, prefix := range prefixes {
		if hasPrefix(stackName, prefix) {
			return true
		}
	}

	return false
}

func hasPrefix(stackName, prefix string) bool {
	return len(stackName) >= len(prefix) && stackName[:len(prefix)] == prefix
}

func toolkit(s string) bool {
	if strings.EqualFold(s, "cdktoolkit") {
		return true
	} else {
		return false
	}
}
