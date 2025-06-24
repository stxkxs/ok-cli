package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bcmdataexports"
	"github.com/aws/aws-sdk-go-v2/service/bcmdataexports/types"
	"github.com/stxkxs/ok-cli/logger"
)

type ExportConfiguration struct {
	Export       ExportDetails       `json:"export"`
	ResourceTags []types.ResourceTag `json:"resourceTags"`
}

type ExportDetails struct {
	Name                      string                          `json:"name"`
	Description               string                          `json:"description"`
	ExportArn                 string                          `json:"exportArn"`
	DataQuery                 types.DataQuery                 `json:"dataQuery"`
	DestinationConfigurations types.DestinationConfigurations `json:"destinationConfigurations"`
	RefreshCadence            types.RefreshCadence            `json:"refreshCadence"`
}

type BillingCostManagement interface{}

type BillingCostManagementClient struct {
	Api *bcmdataexports.Client
}

func NewBillingCostManagementClient() *BillingCostManagementClient {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error loading default foundation configurations")
		return nil
	}

	api := bcmdataexports.NewFromConfig(cfg)

	return &BillingCostManagementClient{Api: api}
}

func (client *BillingCostManagementClient) CreateUpdateExport(conf ExportConfiguration) {
	_, err := client.Api.CreateExport(context.Background(), &bcmdataexports.CreateExportInput{
		Export: &types.Export{
			DataQuery:                 &conf.Export.DataQuery,
			DestinationConfigurations: &conf.Export.DestinationConfigurations,
			Name:                      &conf.Export.Name,
			RefreshCadence:            &conf.Export.RefreshCadence,
			Description:               &conf.Export.Description,
			ExportArn:                 &conf.Export.ExportArn,
		},
		ResourceTags: conf.ResourceTags,
	})

	if err != nil {
		logger.Logger.Error().
			Err(err).
			Msg("error creating billing and cost management export")
		return
	}
}

func (client *BillingCostManagementClient) DestroyReport(r string) {

}
