package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/rs/zerolog/log"
	"github.com/stxkxs/ok-cli/logger"
)

type StsClient struct {
	Client *sts.Client
}

func NewStsClient() *StsClient {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("error loading default aws configurations")
		return nil
	}

	client := sts.NewFromConfig(cfg)
	return &StsClient{Client: client}
}

func NewStsClientFromRole(role, externalId, session string) *StsClient {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("error loading default aws configurations")
		return nil
	}

	client := sts.NewFromConfig(cfg)
	self := &StsClient{Client: client}

	subscriberRoleAssumed, err := self.AssumeRole(role, externalId, session)
	if err != nil {
		logger.Logger.Error().Err(err).
			Str("role", role).
			Msg("error assuming subscriber role")
		return nil
	}

	if subscriberRoleAssumed.Credentials == nil {
		logger.Logger.Error().Msg("no credentials returned from subscriber role")
		return nil
	}

	p := credentials.NewStaticCredentialsProvider(
		*subscriberRoleAssumed.Credentials.AccessKeyId,
		*subscriberRoleAssumed.Credentials.SecretAccessKey,
		*subscriberRoleAssumed.Credentials.SessionToken)

	cfg, err = config.LoadDefaultConfig(context.Background(), config.WithCredentialsProvider(p))
	if err != nil {
		log.Error().Err(err).Msg("error loading aws configurations with custom credentials")
		return nil
	}

	client = sts.NewFromConfig(cfg)
	return &StsClient{Client: client}
}

func (client *StsClient) AssumeRole(role, id, session string) (*sts.AssumeRoleOutput, error) {
	input := &sts.AssumeRoleInput{
		ExternalId:      aws.String(id),
		RoleArn:         aws.String(role),
		RoleSessionName: aws.String(session),
	}

	assumed, err := client.Client.AssumeRole(context.Background(), input)
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Interface("request", input).
			Msg("error assuming role")
		return nil, err
	}

	return assumed, nil
}
