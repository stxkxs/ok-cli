package ecr

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrprivatetypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic"
	ecrpublictypes "github.com/aws/aws-sdk-go-v2/service/ecrpublic/types"
	"github.com/stxkxs/ok-cli/logger"
	"os"
	"strings"
)

func (client *PrivateClient) MaybeCreateRepository(id string, repos []PrivateRepository) {
	if len(repos) == 0 {
		return
	}

	re, err := client.Api.DescribeRepositories(context.Background(), &ecr.DescribeRepositoriesInput{RegistryId: &id})
	if err != nil {
		logger.Logger.Err(err).Msg("error retrieving private ecr repositories")
		return
	}

MaybeCreate:
	for _, repository := range repos {
		for _, r := range re.Repositories {
			if *r.RepositoryName == repository.Name {
				break MaybeCreate
			}
		}

		tags := make([]ecrprivatetypes.Tag, 0, len(repository.Tags))
		for k, v := range repository.Tags {
			tags = append(tags, ecrprivatetypes.Tag{Key: &k, Value: &v})
		}

		request := &ecr.CreateRepositoryInput{
			RegistryId:                 &id,
			RepositoryName:             &repository.Name,
			ImageScanningConfiguration: &ecrprivatetypes.ImageScanningConfiguration{ScanOnPush: repository.ScanOnPush},
			ImageTagMutability:         ecrprivatetypes.ImageTagMutability(strings.ToUpper(repository.Mutability)),
			Tags:                       tags,
		}

		response, err := client.Api.CreateRepository(context.Background(), request)
		if err != nil {
			logger.Logger.Err(err).
				Interface("request", request).
				Msg("error creating private ecr repository")
			os.Exit(1)
		}

		logger.Logger.Info().
			Interface("request", request).
			Interface("response", response).
			Msg("created private ecr repository")
	}
}

func (client *PrivateClient) ConvertDockerImagesToRepositories(images []PrivateDockerImage) []PrivateRepository {
	repos := make([]PrivateRepository, len(images))
	for i, image := range images {
		repos[i] = PrivateRepository{
			Name:       image.Name,
			Version:    image.Version,
			ScanOnPush: image.ScanOnPush,
			Mutability: image.Mutability,
			Tags:       image.Tags,
		}
	}
	return repos
}

func (client *PrivateClient) ConvertHelmChartsToRepositories(charts []PrivateHelmChart) []PrivateRepository {
	repos := make([]PrivateRepository, len(charts))
	for i, chart := range charts {
		repos[i] = PrivateRepository{
			Name:       chart.Name,
			Version:    chart.Version,
			ScanOnPush: chart.ScanOnPush,
			Mutability: chart.Mutability,
			Tags:       chart.Tags,
		}
	}
	return repos
}

func (client *PublicClient) MaybeCreateRepository(id string, repos []PublicRepository) {
	re, err := client.Api.DescribeRepositories(context.Background(), &ecrpublic.DescribeRepositoriesInput{RegistryId: &id})
	if err != nil {
		logger.Logger.Err(err).Msg("error retrieving public ecr repositories")
		return
	}

MaybeCreate:
	for _, repository := range repos {
		for _, r := range re.Repositories {
			if *r.RepositoryName == repository.Name {
				request := &ecrpublic.PutRepositoryCatalogDataInput{
					RepositoryName: &repository.Name,
					CatalogData: &ecrpublictypes.RepositoryCatalogDataInput{
						Description:      &repository.Description,
						AboutText:        &repository.About,
						UsageText:        &repository.Usage,
						Architectures:    repository.Architectures,
						OperatingSystems: repository.OperatingSystems,
					},
				}
				_, err = client.Api.PutRepositoryCatalogData(context.Background(), request)

				if err != nil {
					logger.Logger.Err(err).
						Interface("request", request).
						Msg("error updating public ecr repository")
					os.Exit(1)
				}

				break MaybeCreate
			}
		}

		tags := make([]ecrpublictypes.Tag, 0, len(repository.Tags))
		for k, v := range repository.Tags {
			key, value := k, v
			tags = append(tags, ecrpublictypes.Tag{Key: &key, Value: &value})
		}

		request := &ecrpublic.CreateRepositoryInput{
			RepositoryName: &repository.Name,
			CatalogData: &ecrpublictypes.RepositoryCatalogDataInput{
				Description:      &repository.Description,
				AboutText:        &repository.About,
				UsageText:        &repository.Usage,
				Architectures:    repository.Architectures,
				OperatingSystems: repository.OperatingSystems,
			},
			Tags: tags,
		}

		response, err := client.Api.CreateRepository(context.Background(), request)
		if err != nil {
			logger.Logger.Err(err).
				Interface("request", request).
				Msg("error creating public ecr repository")
			os.Exit(1)
		}

		logger.Logger.Info().
			Interface("request", request).
			Interface("response", response).
			Msg("created public ecr repository")
	}
}

func (client *PublicClient) ConvertDockerImagesToRepositories(images []PublicDockerImage) []PublicRepository {
	repos := make([]PublicRepository, len(images))
	for i, image := range images {
		repos[i] = PublicRepository{
			Name:             image.Name,
			Version:          image.Version,
			ScanOnPush:       image.ScanOnPush,
			Mutability:       image.Mutability,
			Alias:            image.Alias,
			Description:      image.Description,
			About:            image.About,
			Usage:            image.Usage,
			Architectures:    image.Architectures,
			OperatingSystems: image.OperatingSystems,
			Tags:             image.Tags,
		}
	}
	return repos
}

func (client *PublicClient) ConvertHelmChartsToRepositories(charts []PublicHelmChart) []PublicRepository {
	repos := make([]PublicRepository, len(charts))
	for i, chart := range charts {
		repos[i] = PublicRepository{
			Name:             chart.Name,
			Version:          chart.Version,
			ScanOnPush:       chart.ScanOnPush,
			Mutability:       chart.Mutability,
			Alias:            chart.Alias,
			Description:      chart.Description,
			About:            chart.About,
			Usage:            chart.Usage,
			Architectures:    chart.Architectures,
			OperatingSystems: chart.OperatingSystems,
			Tags:             chart.Tags,
		}
	}
	return repos
}

func (client *PrivateClient) DestroyRepository(id, repository string) {
	_, err := client.Api.DeleteRepository(context.Background(), &ecr.DeleteRepositoryInput{
		RegistryId:     &id,
		RepositoryName: &repository,
		Force:          true,
	})

	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("error deleting private ecr repository")
		return
	}

	logger.Logger.Info().
		Interface("id", id).
		Interface("name", repository).
		Msg("destroyed private ecr repository")
}

func (client *PublicClient) DestroyRepository(id, repository string) {
	_, err := client.Api.DeleteRepository(context.Background(), &ecrpublic.DeleteRepositoryInput{
		RegistryId:     &id,
		RepositoryName: &repository,
		Force:          true,
	})

	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("error deleting public ecr repository")
		return
	}

	logger.Logger.Info().
		Interface("id", id).
		Interface("name", repository).
		Msg("destroyed public ecr repository")
}
