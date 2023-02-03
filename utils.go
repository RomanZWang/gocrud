package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func getDDBClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(opts *config.LoadOptions) error {
		opts.Region = CURRENT_REGION
		return nil
	})
	if err != nil {
		panic(err)
	}

	svc := dynamodb.NewFromConfig(cfg)
	return svc
}

func getTableBasics(svc *dynamodb.Client) TableBasics {

	tb := TableBasics{
		DynamoDbClient: svc,
		TableName:      TEST_TABLE_NAME,
	}
	return tb
}
