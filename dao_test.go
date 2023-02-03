package main

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

var databaseClient = getDDBClient()
var basics = getTableBasics(databaseClient)
var expectedId string
var expectedData Data
var expectedUpdatedId string
var expectedOldData Data
var expectedUpdatedData Data
var expectedDeletedId string
var expectedDeletedData Data

func TestMain(m *testing.M) {
	expectedId = "testId"
	expectedData = Data{
		Id:              expectedId,
		Number:          1,
		AssociatedTypes: []string{"t1", "t2"},
	}

	expectedUpdatedId = "testIdUpdate"
	expectedOldData = Data{
		Id:              expectedUpdatedId,
		Number:          2,
		AssociatedTypes: []string{},
	}

	expectedDeletedId = "testIdDeleted"
	expectedDeletedData = Data{
		Id:              expectedDeletedId,
		Number:          2,
		AssociatedTypes: []string{},
	}

	addItemError := basics.CreateTestData(expectedData)
	if addItemError != nil {
		log.Fatalf("failed to insert, %v", addItemError)
	}

	addItemError2 := basics.CreateTestData(expectedOldData)
	if addItemError2 != nil {
		log.Fatalf("failed to insert, %v", addItemError2)
	}

	addItemError3 := basics.CreateTestData(expectedDeletedData)
	if addItemError3 != nil {
		log.Fatalf("failed to insert, %v", addItemError3)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestCreateTestData(t *testing.T) {
	var responseData Data
	response, err := basics.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       expectedData.GetKey(),
		TableName: aws.String(basics.TableName),
	})
	if err != nil {
		log.Printf("Couldn't retrieve %v because %v", expectedId, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &responseData)
		if err != nil {
			log.Printf("Problem unmarshalling response because %v", err)
		}
	}

	assert.EqualValues(t, expectedData, responseData)
}

func TestReadTestData(t *testing.T) {
	var responseData Data
	response, err := basics.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       expectedData.GetKey(),
		TableName: aws.String(basics.TableName),
	})
	if err != nil {
		log.Printf("Couldn't retrieve %v because %v", expectedId, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &responseData)
		if err != nil {
			log.Printf("Problem unmarshalling response because %v", err)
		}
	}

	assert.EqualValues(t, expectedData, responseData)
}

func TestUpdateTestData(t *testing.T) {
	expectedUpdatedData = Data{
		Id:              expectedUpdatedId,
		Number:          23,
		AssociatedTypes: []string{"t3"},
	}
	var responseData Data

	_, err := basics.UpdateTestData(expectedUpdatedData)
	if err != nil {
		log.Printf("Couldn't update because %v", err)
	}

	response, err := basics.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       expectedUpdatedData.GetKey(),
		TableName: aws.String(basics.TableName),
	})
	if err != nil {
		log.Printf("Couldn't retrieve %v because %v", expectedUpdatedId, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &responseData)
		if err != nil {
			log.Printf("Problem unmarshalling response because %v", err)
		}
	}

	assert.EqualValues(t, expectedUpdatedData, responseData)
}

func TestDeleteTestData(t *testing.T) {
	basics.DeleteTestData(expectedDeletedId)
	response, err := basics.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       expectedDeletedData.GetKey(),
		TableName: aws.String(basics.TableName),
	})
	if err != nil {
		log.Printf("Couldn't retrieve %v because %v", expectedDeletedId, err)
	}
	assert.True(t, len(response.Item) == 0)
}
