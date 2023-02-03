package main

import (
	"context"
	"log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

// ListTables lists the DynamoDB table names for the current account.
func (basics TableBasics) ListTables() ([]string, error) {
	var tableNames []string
	tables, err := basics.DynamoDbClient.ListTables(
		context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Printf("Couldn't list tables. Here's why: %v\n", err)
	} else {
		tableNames = tables.TableNames
	}
	return tableNames, err
}

func (basics TableBasics) CreateTestData(data Data) error {
	item, err := attributevalue.MarshalMap(data)
	if err != nil {
		panic(err)
	}
	_, err = basics.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(basics.TableName), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}

func (basics TableBasics) ReadTestData(id string) (Data, error) {
	data := Data{
		Id: id,
	}
	response, err := basics.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: data.GetKey(),
		TableName: aws.String(basics.TableName),
	})
	if err != nil {
		log.Printf("Couldn't retrieve %v because %v", id, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &data)
		if err != nil {
			log.Printf("Problem unmarshalling response because %v", err)
		}
	}
	return data, err
}

func (basics TableBasics) UpdateTestData(data Data) (map[string]interface{}, error) {
	var err error
	var response *dynamodb.UpdateItemOutput
	var attributeMap map[string]interface{}
	update := expression.Set(expression.Name("Number"), expression.Value(data.Number))
	update.Set(expression.Name("AssociatedTypes"), expression.Value(data.AssociatedTypes))
	
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Printf("Couldn't build expression for update. Here's why: %v\n", err)
	} else {
		response, err = basics.DynamoDbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
			TableName:                 aws.String(basics.TableName),
			Key:                       data.GetKey(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		})
		if err != nil {
			log.Printf("Couldn't update movie %v. Here's why: %v\n", data.Id, err)
		} else {
			err = attributevalue.UnmarshalMap(response.Attributes, &attributeMap)
			if err != nil {
				log.Printf("Couldn't unmarshall update response. Here's why: %v\n", err)
			}
		}
	}
	return attributeMap, err
}

func (basics TableBasics) DeleteTestData(id string) error {
	dataToDelete := Data{
		Id: id,
	}

	_, err := basics.DynamoDbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(basics.TableName),
		Key: dataToDelete.GetKey(),
	})

	if err != nil {
		log.Printf("Couldn't delete %v from table because %v\n", dataToDelete.Id, err)
	}
	return err
}