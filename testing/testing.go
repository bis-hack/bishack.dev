// Package testing ...
// Pretty neat little trick from:
// https://brandur.org/fragments/testing-go-project-root
package testing

import (
	"os"
	"path"
	"runtime"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/mock"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	_ = os.Chdir(dir)
}

// ProviderMock ...
type ProviderMock struct {
	mock.Mock
}

func (p *ProviderMock) PutItem(input *dynamodb.PutItemInput) (
	*dynamodb.PutItemOutput,
	error,
) {
	args := p.Called(input)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*dynamodb.PutItemOutput), args.Error(1)
}
func (p *ProviderMock) UpdateItem(input *dynamodb.UpdateItemInput) (
	*dynamodb.UpdateItemOutput,
	error,
) {
	args := p.Called(input)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*dynamodb.UpdateItemOutput), args.Error(1)
}
func (p *ProviderMock) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	args := p.Called(input)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*dynamodb.DeleteItemOutput), args.Error(1)
}
func (p *ProviderMock) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	args := p.Called(input)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*dynamodb.QueryOutput), args.Error(1)
}
func (p *ProviderMock) DescribeTable(input *dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	args := p.Called(input)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*dynamodb.DescribeTableOutput), args.Error(1)
}
