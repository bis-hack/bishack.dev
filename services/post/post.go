package post

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoDBProvider ...
type DynamoDBProvider interface {
	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	UpdateItem(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
	DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
	Query(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	DescribeTable(*dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error)
}

// Client ...
type Client struct {
	TableName string
	Provider  DynamoDBProvider
}

// New creates new Client instance
func New(
	tableName,
	endpoint string,
	provider DynamoDBProvider,
) *Client {
	if provider != nil {
		return &Client{
			TableName: tableName,
			Provider:  provider,
		}
	}

	conf := &aws.Config{
		Region: aws.String("us-east-1"),
	}

	// if endpoint is present
	if endpoint != "" {
		conf.Endpoint = aws.String(endpoint)
	}

	client := session.Must(session.NewSession(conf))

	provider = dynamodb.New(client)

	return &Client{
		TableName: tableName,
		Provider:  provider,
	}
}

// GetCount gets the total number of item in the table
func (c *Client) GetCount() int64 {
	// input
	in := &dynamodb.DescribeTableInput{}
	in.SetTableName(c.TableName)

	resp, err := c.Provider.DescribeTable(in)
	if err != nil {
		// we log non-actionable errors
		fmt.Println("DescribeTable Error:", c.TableName, err.Error())
		return 0
	}

	return *resp.Table.ItemCount
}
