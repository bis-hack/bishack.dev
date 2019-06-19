package dynamo

import (
	"errors"
	"testing"

	test "bishack.dev/testing"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	t.Run("with provider", func(t *testing.T) {
		cli := New("", "", nil)
		c := New("test", "", cli.Provider)
		assert.Equal(t, "test", c.TableName)
	})

	t.Run("with endpoint", func(t *testing.T) {
		c := New("test", "provider", nil)
		assert.Equal(t, "test", c.TableName)
	})
}

func TestQuery(t *testing.T) {
	c := New("", "", nil)
	p := new(test.ProviderMock)

	p.On("Query", mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
		return true
	})).Return(nil, errors.New(""))

	_, err := c.Query(
		"test",
		"test",
		"test",
		map[string]interface{}{},
		false,
		1,
	)

	assert.NotNil(t, err)
}
