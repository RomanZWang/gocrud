package main

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

type Data struct {
	Id string `dynamodbav:"id"`
	Number int
	AssociatedTypes []string
}

func (data Data) GetKey() map[string]types.AttributeValue {
	id, err := attributevalue.Marshal(data.Id)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"id": id}
}