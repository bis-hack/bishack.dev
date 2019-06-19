package post

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"bishack.dev/services/dynamo"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Client ...
type Client struct {
	*dynamo.Client
}

// New creates new Client instance
func New(
	tableName,
	endpoint string,
	provider dynamo.DBProvider,
) *Client {
	return &Client{
		dynamo.New(tableName, endpoint, provider),
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

// CreatePost creates a new post
func (c *Client) CreatePost(params map[string]interface{}) *Post {
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

// GetPost ...
func (c *Client) GetPost(username, id string) *Post {
	ks := "id = :id and created > :created"
	fs := "username = :username"
	vals := map[string]interface{}{
		":id":       id,
		":created":  0,
		":username": username,
	}

	// we set index name to blank since we're not querying
	// global secondary index
	out, err := c.Query("", ks, fs, vals, false, 0)
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

	return posts[0]
}

// GetUser gets all the posts from user
func (c *Client) GetUserPosts(username string) []*Post {
	ks := "publish = :publish and created > :created"
	fs := "username = :username"
	vals := map[string]interface{}{
		":publish":  1,
		":created":  0,
		":username": username,
	}

	out, err := c.Query("publish_index", ks, fs, vals, false, 0)
	if err != nil || len(out.Items) == 0 {
		return nil
	}

	var posts []*Post
	_ = dynamodbattribute.UnmarshalListOfMaps(out.Items, &posts)
	return posts
}

// GetAll ...
func (c *Client) GetPosts() []*Post {
	ks := "publish = :publish and created > :created"
	vals := map[string]interface{}{
		":publish": 1,
		":created": 0,
	}

	out, err := c.Query("publish_index", ks, "", vals, false, 0)
	if err != nil || len(out.Items) == 0 {
		return nil
	}

	var posts []*Post
	_ = dynamodbattribute.UnmarshalListOfMaps(out.Items, &posts)
	return posts
}
