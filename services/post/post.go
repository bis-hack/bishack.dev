package post

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

// Create creates a new post
func (c *Client) Create(params map[string]interface{}) *Post {
	// dates
	now := time.Now().Unix()
	params["created"] = now
	params["updated"] = now

	// parse title to create slug for id
	title := params["title"].(string)
	// remove extra spaces inbetween
	slug := regexp.MustCompile("[^a-zA-Z0-9 ]").ReplaceAllString(title, "")
	// remove outer spaces
	slug = strings.Trim(regexp.MustCompile(`\s+`).ReplaceAllString(slug, " "), " ")
	slug = strings.Replace(slug, " ", "-", -1)
	slug = strings.ToLower(slug)
	// combine
	params["id"] = fmt.Sprintf("%s-%d", slug, now)

	item, _ := dynamodbattribute.MarshalMap(params)

	input := &dynamodb.PutItemInput{}
	input.SetTableName(c.TableName)
	input.SetItem(item)
	input.SetReturnValues("ALL_OLD")

	_, err := c.Provider.PutItem(input)
	if err != nil {
		log.Println("PutItem error:", err.Error())
		return nil
	}

	var post Post
	_ = dynamodbattribute.UnmarshalMap(item, &post)
	return &post
}

// Get ...
func (c *Client) Get(id string) *Post {
	qs := "id = :id and created > :created"
	vals := map[string]interface{}{
		":id":      id,
		":created": 0,
	}

	// we set index name to blank since we're not querying
	// global secondary index
	posts := c.query("", qs, vals, false)
	if posts == nil || len(posts) == 0 {
		return nil
	}

	return posts[0]
}

// GetAll ...
func (c *Client) GetAll() []*Post {
	qs := "publish = :publish and created > :created"
	vals := map[string]interface{}{
		":publish": 1,
		":created": 0,
	}

	return c.query("publish_index", qs, vals, false)
}

// Query
// just pass in the following:
//
//  in      - index name (optional)
//  qs      - query string
//  vals    - expression attribute values
//  forward - to ascending or not
func (c *Client) query(in, qs string, vals map[string]interface{}, forward bool) []*Post {
	// marshal values
	values, _ := dynamodbattribute.MarshalMap(vals)

	// prepare input
	input := &dynamodb.QueryInput{
		TableName:                 &c.TableName,
		KeyConditionExpression:    &qs,
		ExpressionAttributeValues: values,
		ScanIndexForward:          aws.Bool(forward),
	}
	// if index name exists
	if in != "" {
		input.SetIndexName(in)
	}

	out, err := c.Provider.Query(input)
	if err != nil {
		fmt.Println("Query error:", err.Error())
		return nil
	}

	if len(out.Items) == 0 {
		fmt.Println("Query error: Not Found")
		return nil
	}

	var posts []*Post
	_ = dynamodbattribute.UnmarshalListOfMaps(out.Items, &posts)

	return posts
}
