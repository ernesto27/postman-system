package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang.org/x/exp/rand"
)

type DynamoDB struct {
	client    *dynamodb.Client
	tableName string
}

type Collection struct {
	UserID       string `json:"user_id"`
	CollectionID int    `json:"collection_id"`
	Name         string `json:"name"`
}

func NewDynamoDB(key string, secret string, region string, tableName string) *DynamoDB {
	creds := credentials.NewStaticCredentialsProvider(key, secret, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		log.Fatal(err)
	}

	ddbClient := dynamodb.NewFromConfig(cfg)
	return &DynamoDB{
		client:    ddbClient,
		tableName: tableName,
	}
}

func (d *DynamoDB) GetCollections() ([]Collection, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(d.tableName),
	}

	result, err := d.client.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var collections []Collection
	for _, item := range result.Items {
		userIDAttr, ok := item["user_id"].(*types.AttributeValueMemberS)
		if !ok {
			return nil, fmt.Errorf("invalid type for user_id")
		}

		collectionIDAttr, ok := item["collection_id"].(*types.AttributeValueMemberN)
		if !ok {
			return nil, fmt.Errorf("invalid type for collection_id")
		}
		collectionIDValue, err := strconv.Atoi(collectionIDAttr.Value)
		if err != nil {
			return nil, err
		}

		collections = append(collections, Collection{
			UserID:       userIDAttr.Value,
			CollectionID: collectionIDValue,
			Name:         item["name"].(*types.AttributeValueMemberS).Value,
		})
	}

	return collections, nil
}

func (d *DynamoDB) CreateCollection(name string, userID string) error {
	collectionID := rand.Intn(1000000) + 1
	input := &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item: map[string]types.AttributeValue{
			"user_id":       &types.AttributeValueMemberS{Value: userID},
			"collection_id": &types.AttributeValueMemberN{Value: strconv.Itoa(collectionID)},
			"name":          &types.AttributeValueMemberS{Value: name},
		},
	}

	_, err := d.client.PutItem(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}
