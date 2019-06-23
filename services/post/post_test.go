package post

import (
	"errors"
	"testing"

	test "bishack.dev/testing"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	t.Run("without default provider", func(t *testing.T) {
		c := New("beep", "boop", nil)
		assert.NotNil(t, c.Provider.DescribeTable)
	})

	t.Run("with default provider", func(t *testing.T) {
		provider := new(test.DynamoProviderMock)
		c := New("beep", "boop", provider)
		assert.NotNil(t, c.Provider.DescribeTable)
	})
}

func TestUpdatePost(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		provider := new(test.DynamoProviderMock)
		c := New("beep", "boop", provider)

		provider.On("UpdateItem", mock.MatchedBy(func(input *dynamodb.UpdateItemInput) bool {
			return true
		})).Return(nil, errors.New(""))

		err := c.UpdatePost("test", "test", "test", int64(42))
		assert.NotNil(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		provider := new(test.DynamoProviderMock)
		c := New("beep", "boop", provider)

		provider.On("UpdateItem", mock.MatchedBy(func(input *dynamodb.UpdateItemInput) bool {
			return true
		})).Return(&dynamodb.UpdateItemOutput{}, nil)

		err := c.UpdatePost("test", "test", "test", int64(42))
		assert.Nil(t, err)
	})
}

func TestGetCount(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		provider := new(test.DynamoProviderMock)
		c := New("beep", "boop", provider)

		provider.On("DescribeTable", mock.MatchedBy(func(input *dynamodb.DescribeTableInput) bool {
			return true
		})).Return(nil, errors.New(""))

		count := c.GetCount()

		assert.Equal(t, int64(0), count)
		provider.AssertExpectations(t)
	})

	t.Run("ok", func(t *testing.T) {
		provider := new(test.DynamoProviderMock)
		c := New("beep", "boop", provider)

		// table
		table := &dynamodb.TableDescription{}
		table.SetItemCount(31337)

		// output
		out := &dynamodb.DescribeTableOutput{}
		out.SetTable(table)

		provider.On("DescribeTable", mock.MatchedBy(func(input *dynamodb.DescribeTableInput) bool {
			return true
		})).Return(out, nil)

		count := c.GetCount()

		assert.Equal(t, int64(31337), count)
		provider.AssertExpectations(t)
	})
}

func TestCreate(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		provider := new(test.DynamoProviderMock)
		c := New("bee", "boop", provider)

		provider.On("PutItem", mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
			return true
		})).Return(nil, errors.New(""))

		p := c.CreatePost(map[string]interface{}{
			"title":    "hello world",
			"username": "hello",
		})

		assert.Nil(t, p)
		provider.AssertExpectations(t)
	})

	t.Run("ok", func(t *testing.T) {
		provider := new(test.DynamoProviderMock)
		c := New("bee", "boop", provider)

		out := &dynamodb.PutItemOutput{}
		attr, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
			"title": "hello world",
		})
		out.SetAttributes(attr)

		provider.On("PutItem", mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
			return true
		})).Return(out, nil)

		p := c.CreatePost(map[string]interface{}{
			"title":    "hello world",
			"username": "hello",
		})

		assert.NotNil(t, p)
		assert.Equal(t, "hello world", p.Title)
		provider.AssertExpectations(t)
	})
}

func TestQuery(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		p := new(test.DynamoProviderMock)
		c := New("bee", "boop", p)

		p.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(nil, errors.New(""))

		_, err := c.Query("x", "", "", map[string]interface{}{}, false, 0)
		assert.NotNil(t, err)
	})

	t.Run("0 items", func(t *testing.T) {
		p := new(test.DynamoProviderMock)
		c := New("bee", "boop", p)

		out := &dynamodb.QueryOutput{}
		out.SetItems([]map[string]*dynamodb.AttributeValue{})
		p.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(out, nil)

		out, err := c.Query("x", "", "", map[string]interface{}{}, false, 0)
		assert.Nil(t, err)
		assert.Empty(t, out.Items)
	})

	t.Run("found items", func(t *testing.T) {
		p := new(test.DynamoProviderMock)
		c := New("bee", "boop", p)

		out := &dynamodb.QueryOutput{}
		item, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
			"title":    "test",
			"id":       "testing",
			"username": "test",
		})
		out.SetItems([]map[string]*dynamodb.AttributeValue{
			item,
		})
		p.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(out, nil)

		out, err := c.Query("x", "", "", map[string]interface{}{}, false, 0)
		assert.Nil(t, err)
		assert.NotNil(t, out)
	})
}

func TestGetAll(t *testing.T) {
	t.Run("0 item", func(t *testing.T) {
		p := new(test.DynamoProviderMock)
		c := New("bee", "boop", p)

		p.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(&dynamodb.QueryOutput{}, nil)

		posts := c.GetPosts()
		assert.Nil(t, posts)
	})

	t.Run("ok", func(t *testing.T) {
		p := new(test.DynamoProviderMock)
		c := New("bee", "boop", p)

		out := &dynamodb.QueryOutput{}
		item, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
			"title":    "test",
			"id":       "testing",
			"username": "test",
		})
		out.SetItems([]map[string]*dynamodb.AttributeValue{
			item,
		})
		p.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(out, nil)

		posts := c.GetPosts()
		assert.NotNil(t, posts)
	})
}

func TestGet(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		p := new(test.DynamoProviderMock)
		c := New("bee", "boop", p)

		p.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(nil, errors.New(""))

		post := c.GetPost("test", "test")
		assert.Nil(t, post)
	})

	t.Run("0 items", func(t *testing.T) {
		p := new(test.DynamoProviderMock)
		c := New("bee", "boop", p)

		p.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(&dynamodb.QueryOutput{}, nil)

		post := c.GetPost("test", "test")
		assert.Nil(t, post)
	})

	t.Run("ok", func(t *testing.T) {
		p := new(test.DynamoProviderMock)
		c := New("bee", "boop", p)

		out := &dynamodb.QueryOutput{}
		item, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
			"title":    "test",
			"id":       "testing",
			"username": "test",
		})
		out.SetItems([]map[string]*dynamodb.AttributeValue{
			item,
		})
		p.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(out, nil)

		post := c.GetPost("test", "test")
		assert.NotNil(t, post)
	})
}

func TestGetUserPosts(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		p := new(test.DynamoProviderMock)
		c := New("bee", "boop", p)

		p.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(nil, errors.New(""))

		posts := c.GetUserPosts("test")
		assert.Nil(t, posts)
	})

	t.Run("ok", func(t *testing.T) {
		p := new(test.DynamoProviderMock)
		c := New("bee", "boop", p)

		out := &dynamodb.QueryOutput{}
		item, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
			"title":    "test",
			"id":       "testing",
			"username": "test",
		})
		out.SetItems([]map[string]*dynamodb.AttributeValue{
			item,
		})
		p.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(out, nil)

		posts := c.GetUserPosts("test")
		assert.NotNil(t, posts)
		assert.Equal(t, "testing", posts[0].ID)
		assert.Equal(t, "test", posts[0].Title)
		assert.Equal(t, "test", posts[0].Username)
	})
}
