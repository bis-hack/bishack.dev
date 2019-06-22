package like

import (
	"time"

	"bishack.dev/services/dynamo"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

// Client ...
type Client struct {
	*dynamo.Client
}

// New ...
func New(
	tableName,
	endpoint string,
	provider dynamo.Provider,
) *Client {
	return &Client{
		dynamo.New(tableName, endpoint, provider),
	}
}

// GetLikes ...
func (c *Client) GetLikes(id string) ([]*Like, error) {
	ks := "id = :id and created > :created"
	vals := map[string]interface{}{
		":id":      id,
		":created": 0,
	}

	out, err := c.Query("", ks, "", vals, false, 0)
	if err != nil {
		return nil, errors.Wrap(err, "GetLikes/Query error")
	}

	if len(out.Items) == 0 {
		return nil, errors.New("GetLikes/NotFound")
	}

	var likes []*Like

	_ = dynamodbattribute.UnmarshalListOfMaps(out.Items, &likes)
	return likes, nil
}

// ToggleLike ...
func (c *Client) ToggleLike(id, username string) error {
	// if not found, we add
	l, err := c.GetLike(id, username)
	if err != nil {
		err := c.addLike(id, username)
		if err != nil {
			return errors.Wrap(err, "Like")
		}
		return nil
	}

	// otherwise, we delete
	err = c.removeLike(id, l.Created)
	if err != nil {
		return errors.Wrap(err, "Like")
	}
	return nil
}

// GetLike ...
func (c *Client) GetLike(id, username string) (*Like, error) {
	ks := "id = :id and created > :created"
	fs := "username = :username"
	vals := map[string]interface{}{
		":id":       id,
		":created":  0,
		":username": username,
	}

	out, err := c.Query("", ks, fs, vals, false, 0)
	if err != nil {
		return nil, errors.Wrap(err, "GetLike/Query error")
	}

	if len(out.Items) == 0 {
		return nil, errors.New("Not found")
	}

	var like Like
	_ = dynamodbattribute.UnmarshalMap(out.Items[0], &like)
	return &like, nil
}

// addLike adds a new item on the likes table with the given username and id
func (c *Client) addLike(id, username string) error {
	item, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
		"id":       id,
		"username": username,
		"created":  time.Now().Unix(),
	})

	input := &dynamodb.PutItemInput{
		Item: item,
	}
	input.SetTableName(c.TableName)

	_, err := c.Provider.PutItem(input)
	if err != nil {
		return errors.Wrap(err, "addLike/PutItem error")
	}

	return nil
}

// removeLike removes like
func (c *Client) removeLike(id string, created int64) error {
	keys, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
		"id":      id,
		"created": created,
	})

	input := &dynamodb.DeleteItemInput{}
	input.SetTableName(c.TableName)
	input.SetKey(keys)

	_, err := c.Provider.DeleteItem(input)
	if err != nil {
		return errors.Wrap(err, "removeLike/DeleteItem error")
	}

	return nil
}
