package post

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	t.Run("without default provider", func(t *testing.T) {
		c := New("beep", "boop", nil)
		assert.NotNil(t, c.Provider.DescribeTable)
	})

	t.Run("with default provider", func(t *testing.T) {
		provider := new(providerMock)
		c := New("beep", "boop", provider)
		assert.NotNil(t, c.Provider.DescribeTable)
	})
}

func TestGetCount(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		provider := new(providerMock)
		c := New("beep", "boop", provider)

		provider.On("DescribeTable", mock.MatchedBy(func(input *dynamodb.DescribeTableInput) bool {
			return true
		})).Return(nil, errors.New(""))

		count := c.GetCount()

		assert.Equal(t, int64(0), count)
		provider.AssertExpectations(t)
	})

	t.Run("ok", func(t *testing.T) {
		provider := new(providerMock)
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
