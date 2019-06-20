package like

import (
	"bishack.dev/services/dynamo"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

type Client struct {
	*dynamo.Client
}

func New(
	tableName,
	endpoint string,
	provider dynamo.DBProvider,
) *Client {
	return &Client{
		dynamo.New(tableName, endpoint, provider),
	}
}

// GetLikes ...
func (c *Client) GetLikes(id string) ([]*Like, error) {
	ks := "id = :id"
	vals := map[string]interface{}{
		":id": id,
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
	if _, err := c.getLike(id, username); err != nil {
		err := c.addLike(id, username)
		if err != nil {
			return errors.Wrap(err, "Like")
		}
		return nil
	}

	// otherwise, we delete
	err := c.removeLike(id, username)
	if err != nil {
		return errors.Wrap(err, "Like")
	}
	return nil
}

func (c *Client) getLike(id, username string) (*Like, error) {
	ks := "id = :id"
	fs := "username = :username"
	vals := map[string]interface{}{
		":id":       id,
		":username": username,
	}

	out, err := c.Query("", ks, fs, vals, false, 1)
	if err != nil {
		return nil, errors.Wrap(err, "getLike/Query error")
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
	})

	input := &dynamodb.PutItemInput{
		Item: item,
	}

	_, err := c.Provider.PutItem(input)
	if err != nil {
		return errors.Wrap(err, "addLike/PutItem error")
	}

	return nil
}

// removeLike removes like
func (c *Client) removeLike(id, username string) error {
	keys, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
		"id":       id,
		"username": username,
	})

	input := &dynamodb.DeleteItemInput{}
	input.SetKey(keys)

	_, err := c.Provider.DeleteItem(input)
	if err != nil {
		return errors.Wrap(err, "removeLike/DeleteItem error")
	}

	return nil
}
