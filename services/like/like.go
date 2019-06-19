package like

import (
	"log"

	"bishack.dev/services/dynamo"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
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

// get gets a particular item from the likes table given the username and id
// key
func (c *Client) get(username, id string) *Like {
	ks := "id = :id and username = :username"
	vals := map[string]interface{}{
		":id":       id,
		":username": username,
	}

	out, err := c.Query("", ks, "", vals, false, 1)
	if err != nil {
		log.Println(errors.Wrap(err, "Like query error:"))
		return nil
	}

	if len(out.Items) == 0 {
		return nil
	}

	var like *Like
	_ = dynamodbattribute.UnmarshalMap(out.Items[0], like)
	return like
}

// like add new item on the likes table with the given username and id
func (c *Client) like(username, id string) error {
	item, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
		"id":       id,
		"username": username,
	})

	input := &dynamodb.PutItemInput{
		Item: item,
	}

	_, err := c.Provider.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
