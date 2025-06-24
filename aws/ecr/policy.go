package ecr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/stxkxs/ok-cli/logger"
	"os"
)

func (client *PrivateClient) GetRepositoryPolicy(account string, repository string) *ecr.GetRepositoryPolicyOutput {
	response, err := client.Api.GetRepositoryPolicy(context.Background(),
		&ecr.GetRepositoryPolicyInput{
			RegistryId:     &account,
			RepositoryName: &repository,
		})

	if err != nil {
		var rnf *types.RepositoryPolicyNotFoundException
		if errors.As(err, &rnf) {
			logger.Logger.Warn().Msg("repository policy does not exist. going to move ahead.")
			return nil
		}

		logger.Logger.Error().
			Err(err).
			Msg("error getting ecr repository policy")
	}

	return response
}

func (client *PrivateClient) AddAccountToPolicy(account, repository, policy string) {
	response := client.GetRepositoryPolicy(account, repository)

	if response != nil && response.PolicyText != nil {
		client.addAccountToPolicy(account, repository, *response.PolicyText, policy)
	} else {
		client.CreateRepositoryPolicy(account, repository, policy)
	}
}

func (client *PrivateClient) CreateRepositoryPolicy(account, repository, policy string) {
	re, err := client.Api.SetRepositoryPolicy(context.Background(),
		&ecr.SetRepositoryPolicyInput{
			RegistryId:     &account,
			RepositoryName: &repository,
			PolicyText:     aws.String(fmt.Sprintf(repositoryPolicy, policy)),
			Force:          false,
		})

	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("error creating ecr repository policy")
		os.Exit(1)
	}

	logger.Logger.Info().
		Interface("response", re).
		Msg("created ecr repository policy")
}

func (client *PrivateClient) RemoveAccountFromPolicy(account, repository, remove string) {
	response := client.GetRepositoryPolicy(account, repository)

	if response == nil {
		logger.Logger.Warn().Msg("ecr repository policy is already empty")
		return
	}

	t := fmt.Sprintf("arn:aws:iam::%s:root", remove)
	updated := client.removeAccountFromPolicy(t, *response.PolicyText)

	if updated == nil {
		deleted, err := client.Api.DeleteRepositoryPolicy(context.Background(),
			&ecr.DeleteRepositoryPolicyInput{
				RegistryId:     &account,
				RepositoryName: &repository,
			})

		if err != nil {
			logger.Logger.Error().
				Err(err).
				Msg("error getting ecr repository policy")
			os.Exit(1)
		}

		logger.Logger.Info().
			Interface("response", deleted).
			Msg("ecr repository policy deleted because it was empty")
	} else {
		re, err := client.Api.SetRepositoryPolicy(context.Background(),
			&ecr.SetRepositoryPolicyInput{
				RegistryId:     &account,
				RepositoryName: &repository,
				PolicyText:     aws.String(string(updated)),
				Force:          false,
			})

		if err != nil {
			logger.Logger.Error().
				Err(err).
				Msg("error removing account from ecr repository policy")
			os.Exit(1)
		}

		logger.Logger.Info().
			Interface("response", re).
			Msg("ecr repository policy set")
	}
}

func (client *PrivateClient) addAccountToPolicy(account, repository, existing, extension string) {
	var current map[string]interface{}
	err := json.Unmarshal([]byte(existing), &current)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error unmarshaling ecr repository policy json")
	}

	var e map[string]interface{}
	err = json.Unmarshal([]byte(extension), &e)
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("error unmarshaling ecr repository policy extension json")
		os.Exit(1)
	}

	statements := current["Statement"].([]interface{})
	statements = append(statements, e)
	current["Statement"] = statements

	u, err := json.Marshal(current)
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("error marshaling ecr repository policy")
		os.Exit(1)
	}

	re, err := client.Api.SetRepositoryPolicy(context.Background(),
		&ecr.SetRepositoryPolicyInput{
			RegistryId:     &account,
			RepositoryName: &repository,
			PolicyText:     aws.String(string(u)),
			Force:          false,
		})

	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("error adding account to ecr repository policy")
		os.Exit(1)
	}

	logger.Logger.Info().
		Interface("response", re).
		Msg("ecr repository policy set")
}

func (client *PrivateClient) removeAccountFromPolicy(target, p string) []byte {
	var policy map[string]interface{}
	if err := json.Unmarshal([]byte(p), &policy); err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("failed to read ecr repository policy")
		return nil
	}

	statements := policy["Statement"].([]interface{})
	var latest []interface{}
	for _, stmt := range statements {
		statement := stmt.(map[string]interface{})
		principal := statement["Principal"].(map[string]interface{})["AWS"]
		if principal != target {
			latest = append(latest, statement)
		}
	}
	policy["Statement"] = latest

	updated, err := json.Marshal(policy)
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("failed to marshal updated ecr repository policy")
		os.Exit(1)
	}

	if latest == nil {
		return nil
	} else {
		return updated
	}
}
