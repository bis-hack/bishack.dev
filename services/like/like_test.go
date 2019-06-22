package like

import (
	"regexp"
	"testing"

	test "bishack.dev/testing"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetLike(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		m := new(test.DynamoProviderMock)
		c := New("a", "b", m)

		m.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(nil, errors.New("beep"))

		l, e := c.GetLike("test", "ing")
		assert.NotNil(t, e)
		assert.Nil(t, l)
		assert.Regexp(t, regexp.MustCompile(`(?i)getlike/query error: beep`), e.Error())
	})

	t.Run("not found", func(t *testing.T) {
		m := new(test.DynamoProviderMock)
		c := New("a", "b", m)

		m.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(&dynamodb.QueryOutput{}, nil)

		l, e := c.GetLike("test", "ing")
		assert.NotNil(t, e)
		assert.Nil(t, l)
		assert.Regexp(t, regexp.MustCompile(`(?i)not found`), e.Error())
	})

	t.Run("ok", func(t *testing.T) {
		m := new(test.DynamoProviderMock)
		c := New("a", "b", m)

		out := &dynamodb.QueryOutput{}
		item, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
			"ID":       "test",
			"Username": "ing",
		})
		out.SetItems([]map[string]*dynamodb.AttributeValue{
			item,
		})
		out.SetCount(1)
		m.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(out, nil)

		l, e := c.GetLike("test", "ing")

		assert.Nil(t, e)
		assert.NotNil(t, l)
		assert.Equal(t, "test", l.ID)
		assert.Equal(t, "ing", l.Username)
	})
}

func TestAddLike(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		m := new(test.DynamoProviderMock)
		c := New("a", "b", m)

		m.On("PutItem", mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
			return true
		})).Return(nil, errors.New("boop"))

		e := c.addLike("test", "ing")
		assert.NotNil(t, e)
	})

	t.Run("ok", func(t *testing.T) {
		m := new(test.DynamoProviderMock)
		c := New("a", "b", m)

		m.On("PutItem", mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
			return true
		})).Return(&dynamodb.PutItemOutput{}, nil)

		e := c.addLike("test", "ing")
		assert.Nil(t, e)
	})
}

func TestRemoveLike(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		m := new(test.DynamoProviderMock)
		c := New("a", "b", m)

		m.On("DeleteItem", mock.MatchedBy(func(input *dynamodb.DeleteItemInput) bool {
			return true
		})).Return(nil, errors.New("boop"))

		e := c.removeLike("test", 123)
		assert.NotNil(t, e)
	})

	t.Run("ok", func(t *testing.T) {
		m := new(test.DynamoProviderMock)
		c := New("a", "b", m)

		m.On("DeleteItem", mock.MatchedBy(func(input *dynamodb.DeleteItemInput) bool {
			return true
		})).Return(&dynamodb.DeleteItemOutput{}, nil)

		e := c.removeLike("test", 123)
		assert.Nil(t, e)
	})
}

func TestToggleLike(t *testing.T) {
	t.Run("addLike error", func(t *testing.T) {
		m := new(test.DynamoProviderMock)
		c := New("a", "b", m)

		m.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(nil, errors.New(""))
		m.On("PutItem", mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
			return true
		})).Return(nil, errors.New("addLike error"))

		e := c.ToggleLike("test", "ing")
		assert.NotNil(t, e)
	})

	t.Run("addLike ok", func(t *testing.T) {
		m := new(test.DynamoProviderMock)
		c := New("a", "b", m)

		m.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(nil, errors.New(""))
		m.On("PutItem", mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
			return true
		})).Return(&dynamodb.PutItemOutput{}, nil)

		e := c.ToggleLike("test", "ing")
		assert.Nil(t, e)
	})

	t.Run("removeLike error", func(t *testing.T) {
		m := new(test.DynamoProviderMock)
		c := New("a", "b", m)

		item, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
			"id":       "x",
			"username": "y",
		})
		m.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(&dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{
				item,
			},
		}, nil)

		m.On("DeleteItem", mock.MatchedBy(func(input *dynamodb.DeleteItemInput) bool {
			return true
		})).Return(nil, errors.New(""))

		e := c.ToggleLike("test", "ing")
		assert.NotNil(t, e)
	})

	t.Run("removeLike ok", func(t *testing.T) {
		m := new(test.DynamoProviderMock)
		c := New("a", "b", m)

		item, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
			"id":       "x",
			"username": "y",
		})
		m.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(&dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{
				item,
			},
		}, nil)

		m.On("DeleteItem", mock.MatchedBy(func(input *dynamodb.DeleteItemInput) bool {
			return true
		})).Return(&dynamodb.DeleteItemOutput{}, nil)

		e := c.ToggleLike("test", "ing")
		assert.Nil(t, e)
	})
}

func TestGetLikes(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		m := new(test.DynamoProviderMock)
		c := New("a", "b", m)

		m.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(nil, errors.New("beep"))

		l, e := c.GetLikes("test")
		assert.NotNil(t, e)
		assert.Nil(t, l)
		assert.Regexp(t, regexp.MustCompile(`(?i)getlikes/query error: beep`), e.Error())
	})

	t.Run("not found", func(t *testing.T) {
		m := new(test.DynamoProviderMock)
		c := New("a", "b", m)

		m.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(&dynamodb.QueryOutput{}, nil)

		l, e := c.GetLikes("test")
		assert.NotNil(t, e)
		assert.Nil(t, l)
		assert.Regexp(t, regexp.MustCompile(`(?i)notfound`), e.Error())
	})

	t.Run("ok", func(t *testing.T) {
		m := new(test.DynamoProviderMock)
		c := New("a", "b", m)

		out := &dynamodb.QueryOutput{}
		item, _ := dynamodbattribute.MarshalMap(map[string]interface{}{
			"ID":       "test",
			"Username": "ing",
		})
		out.SetItems([]map[string]*dynamodb.AttributeValue{
			item,
		})
		out.SetCount(1)
		m.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
			return true
		})).Return(out, nil)

		l, e := c.GetLikes("test")

		assert.Nil(t, e)
		assert.NotNil(t, l)
		assert.Equal(t, 1, len(l))
		assert.Equal(t, "test", l[0].ID)
		assert.Equal(t, "ing", l[0].Username)
	})
}
