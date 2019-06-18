package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DBProvider ...
type DBProvider interface {
	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	UpdateItem(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
	DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
	Query(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	DescribeTable(*dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error)
}

// Client ...
type Client struct {
	TableName string
	Provider  DBProvider
}

// New creates new Client instance
func New(
	tableName,
	endpoint string,
	provider DBProvider,
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

// Query ...
// just pass in the following:
//
//  in      - index name (optional)
//  fs      - filter expression string
//  ks      - key expression string
//  vals    - expression attribute values
//  forward - to ascending or not
//  limit   - total number of returned rows
func (c *Client) Query(
	in,
	ks,
	fs string,
	vals map[string]interface{},
	forward bool,
	limit int64,
) (*dynamodb.QueryOutput, error) {
	// marshal values
	values, _ := dynamodbattribute.MarshalMap(vals)

	// prepare input
	input := &dynamodb.QueryInput{
		TableName:                 &c.TableName,
		KeyConditionExpression:    &ks,
		ExpressionAttributeValues: values,
		ScanIndexForward:          aws.Bool(forward),
	}
	if limit != 0 {
		input.SetLimit(limit)
	}
	// if limit exists
	// if index name exists
	if in != "" {
		input.SetIndexName(in)
	}
	// if fs exists
	if fs != "" {
		input.SetFilterExpression(fs)
	}

	return c.Provider.Query(input)
}
