package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costandusagereportservice"
	"github.com/aws/aws-sdk-go-v2/service/costandusagereportservice/types"
	"github.com/stxkxs/ok-cli/logger"
)

type Report struct {
	Name                     string                     `mapstructure:"name"`
	Bucket                   string                     `mapstructure:"bucket"`
	Prefix                   string                     `mapstructure:"prefix"`
	RefreshClosedReports     bool                       `mapstructure:"refreshclosedreports"`
	Region                   types.AWSRegion            `mapstructure:"region"`
	Format                   types.ReportFormat         `mapstructure:"format"`
	Compression              types.CompressionFormat    `mapstructure:"compression"`
	ReportVersioning         types.ReportVersioning     `mapstructure:"reportVersioning"`
	TimeUnit                 types.TimeUnit             `mapstructure:"timeunit"`
	AdditionalArtifacts      []types.AdditionalArtifact `mapstructure:"additionalArtifacts"`
	AdditionalSchemaElements []types.SchemaElement      `mapstructure:"additionalSchemaElements"`
	Tags                     map[string]string          `mapstructure:"tags"`
}

type CostUsageReporting interface{}

type CostUsageReportingClient struct {
	Api *costandusagereportservice.Client
}

func NewCostUsageReportingClient() *CostUsageReportingClient {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error loading default foundation configurations")
		return nil
	}

	api := costandusagereportservice.NewFromConfig(cfg)

	return &CostUsageReportingClient{Api: api}
}

func (client *CostUsageReportingClient) CreateUpdateReport(r Report) {
	var tags []types.Tag
	for k, v := range r.Tags {
		tags = append(tags, types.Tag{Key: &k, Value: &v})
	}

	if client.exists(r.Name) {
		modify(r, client)
	} else {
		create(r, tags, client)
	}
}

func (client *CostUsageReportingClient) DestroyReport(r string) {
	_, err := client.Api.DeleteReportDefinition(context.Background(),
		&costandusagereportservice.DeleteReportDefinitionInput{ReportName: &r})

	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("error deleting cost report definition")
		return
	}

	logger.Logger.Info().
		Str("report.definition", r).
		Msg("deleted cost report definition")
}

func create(r Report, tags []types.Tag, client *CostUsageReportingClient) bool {
	request := costandusagereportservice.PutReportDefinitionInput{
		ReportDefinition: &types.ReportDefinition{
			AdditionalSchemaElements: r.AdditionalSchemaElements,
			Compression:              r.Compression,
			Format:                   r.Format,
			ReportName:               aws.String(r.Name),
			S3Bucket:                 aws.String(r.Bucket),
			S3Prefix:                 aws.String(r.Prefix),
			S3Region:                 r.Region,
			TimeUnit:                 r.TimeUnit,
			AdditionalArtifacts:      r.AdditionalArtifacts,
			RefreshClosedReports:     aws.Bool(r.RefreshClosedReports),
			ReportVersioning:         r.ReportVersioning,
		},
		Tags: tags,
	}

	response, err := client.Api.PutReportDefinition(context.Background(), &request)
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Interface("request", request).
			Msg("error putting report definition")
		return true
	}

	logger.Logger.Info().
		Interface("report.definition", response).
		Msg("put report definition")
	return false
}

func modify(r Report, client *CostUsageReportingClient) bool {
	request := costandusagereportservice.ModifyReportDefinitionInput{
		ReportDefinition: &types.ReportDefinition{
			AdditionalSchemaElements: r.AdditionalSchemaElements,
			Compression:              r.Compression,
			Format:                   r.Format,
			ReportName:               aws.String(r.Name),
			S3Bucket:                 aws.String(r.Bucket),
			S3Prefix:                 aws.String(r.Prefix),
			S3Region:                 r.Region,
			TimeUnit:                 r.TimeUnit,
			AdditionalArtifacts:      r.AdditionalArtifacts,
			RefreshClosedReports:     aws.Bool(r.RefreshClosedReports),
			ReportVersioning:         r.ReportVersioning,
		},
		ReportName: aws.String(r.Name),
	}

	response, err := client.Api.ModifyReportDefinition(context.Background(), &request)
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Interface("request", request).
			Msg("error modifying report definition")
		return true
	}

	logger.Logger.Info().
		Interface("report.definition", response).
		Msg("modified report definition")
	return false
}

func (client *CostUsageReportingClient) exists(s string) bool {
	response, err := client.Api.DescribeReportDefinitions(context.Background(), &costandusagereportservice.DescribeReportDefinitionsInput{})
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("error describing report definitions. continuing.")
	}

	for _, d := range response.ReportDefinitions {
		if *d.ReportName == s {
			logger.Logger.Debug().Msg("report definition already exists")
			return true
		}
	}

	return false
}
