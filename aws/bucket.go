package aws

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/stxkxs/ok-cli/logger"
	"os"
)

type Bucket struct {
	Name      string            `mapstructure:"name"`
	Region    string            `mapstructure:"region"`
	Ownership string            `mapstructure:"ownership"`
	Policy    string            `mapstructure:"policy"`
	Tags      map[string]string `mapstructure:"tags"`
}

type S3 interface {
	CreateBucket(b Bucket) bool
}

type BucketClient struct {
	Api *s3.Client
}

func NewBucketClient() *BucketClient {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error loading default foundation configurations")
		return nil
	}

	api := s3.NewFromConfig(cfg)

	return &BucketClient{Api: api}
}

func (client *BucketClient) CreateUpdateBucket(b Bucket) {
	found, err := client.Exists(b.Name)
	if err != nil {
		logger.Logger.
			Err(err).
			Msg("error checking if bucket exists. trying anyways.")
	}

	if found {
		logger.Logger.Warn().Msg("you own this bucket! updating it.")
	} else {
		createRequest := s3.CreateBucketInput{
			Bucket: aws.String(b.Name),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(b.Region),
			},
		}

		_, err = client.Api.CreateBucket(context.Background(), &createRequest)
		if err != nil {
			logger.Logger.Error().
				Err(err).
				Interface("request", createRequest).
				Msg("error creating bucket")
			os.Exit(1)
			return
		}
	}

	var tags []types.Tag
	for k, v := range b.Tags {
		tags = append(tags, types.Tag{Key: &k, Value: &v})
	}

	taggingRequest := s3.PutBucketTaggingInput{
		Bucket:  aws.String(b.Name),
		Tagging: &types.Tagging{TagSet: tags},
	}

	_, err = client.Api.PutBucketTagging(context.Background(), &taggingRequest)
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Interface("request", taggingRequest).
			Msg("error tagging bucket")
		os.Exit(1)
	}

	policyRequest := s3.PutBucketPolicyInput{
		Bucket: aws.String(b.Name),
		Policy: &b.Policy,
	}

	_, err = client.Api.PutBucketPolicy(context.Background(), &policyRequest)
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Interface("request", policyRequest).
			Msg("error enabling access to bucket")
		os.Exit(1)
	}

	logger.Logger.Info().
		Msg("created, tagged, and permitted bucket")
}

func (client *BucketClient) Exists(n string) (bool, error) {
	_, err := client.Api.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(n),
	})

	e := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			var notFound *types.NotFound
			switch {
			case errors.As(apiError, &notFound):
				logger.Logger.Debug().
					Str("bucket", n).
					Msg("bucket is available")
				e = false
				err = nil
			default:
				logger.Logger.Warn().
					Str("bucket", n).
					Msg("you dont have access to the bucket you are looking for")
			}
		}
	} else {
		logger.Logger.Debug().
			Str("bucket", n).
			Msg("you own this bucket already")
	}

	return e, err
}

func (client *BucketClient) Destroy(n string) {
	response, err := client.Api.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{Bucket: &n})
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("error getting buckets objects")
		return
	}

	for _, object := range response.Contents {
		_, err = client.Api.DeleteObject(context.Background(),
			&s3.DeleteObjectInput{Bucket: aws.String(n), Key: object.Key})

		if err != nil {
			logger.Logger.Error().
				Err(err).
				Msg("error deleting s3 object")
			return
		}
	}

	_, err = client.Api.DeleteBucket(context.Background(),
		&s3.DeleteBucketInput{Bucket: &n})

	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("error deleting s3 bucket")
		return
	}

	logger.Logger.Info().
		Str("bucket", n).
		Msg("deleted s3 bucket")
}
