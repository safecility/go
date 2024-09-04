package gbigquery

import (
	"cloud.google.com/go/bigquery"
	"context"
	"github.com/rs/zerolog/log"
)

type BQTableConfig struct {
	Dataset string `json:"dataset"`
	Table   string `json:"table"`
	Schema  struct {
		Name     string `json:"name"`
		FilePath string `json:"filePath"`
		Revision string `json:"revision"`
	} `json:"schema"`
}

type BQTable struct {
	client *bigquery.Client
}

func NewBQTable(client *bigquery.Client) *BQTable {
	return &BQTable{client}
}

func (bqt BQTable) CheckOrCreateBigqueryTable(config *BQTableConfig, metaData *bigquery.TableMetadata) (*bigquery.TableMetadata, error) {
	ctx := context.Background()

	tableRef := bqt.client.Dataset(config.Dataset).Table(config.Table)

	tableMetadata, err := tableRef.Metadata(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get table metadata")
	}
	if tableMetadata == nil {
		err = tableRef.Create(ctx, metaData)
		if err != nil {
			return nil, err
		}
		log.Info().Msg("Created bigquery table")
		tableMetadata, err = tableRef.Metadata(ctx)
	}

	return tableMetadata, err
}
