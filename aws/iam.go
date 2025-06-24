package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	sts "github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/stxkxs/ok-cli/logger"
)

type IamClient struct {
	Client *iam.Client
}

type Iam interface {
	SimulatePolicy(role string, actions []string) ([]types.EvaluationResult, error)
}

func NewIamClient() *IamClient {
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("error getting iam client using credentials")
		return nil
	}

	return &IamClient{Client: iam.NewFromConfig(cfg)}
}

func NewIamClientWithCredentials(c *sts.Credentials) *IamClient {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithCredentialsProvider(
		credentials.NewStaticCredentialsProvider(*c.AccessKeyId, *c.SecretAccessKey, *c.SessionToken)))

	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("error getting iam client using credentials")
		return nil
	}

	return &IamClient{Client: iam.NewFromConfig(cfg)}
}

func (c *IamClient) SimulatePolicy(role string, actions, resources []string, ctx []types.ContextEntry) ([]types.EvaluationResult, error) {
	re, err := c.Client.SimulatePrincipalPolicy(context.Background(),
		&iam.SimulatePrincipalPolicyInput{
			ActionNames:     actions,
			PolicySourceArn: &role,
			ContextEntries:  ctx,
			ResourceArns:    resources,
		})

	if err != nil {
		logger.Logger.Error().
			Err(err).
			Str("role", role).
			Strs("actions", actions).
			Strs("resources", resources).
			Interface("context", ctx).
			Msg("error simluating policy")
		return nil, err
	}

	return re.EvaluationResults, nil
}
