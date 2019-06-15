package post

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/mock"
)

type providerMock struct {
	mock.Mock
}

func (p *providerMock) PutItem(input *dynamodb.PutItemInput) (
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
func (p *providerMock) UpdateItem(input *dynamodb.UpdateItemInput) (
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
func (p *providerMock) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	args := p.Called(input)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*dynamodb.DeleteItemOutput), args.Error(1)
}
func (p *providerMock) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	args := p.Called(input)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*dynamodb.QueryOutput), args.Error(1)
}
func (p *providerMock) DescribeTable(input *dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	args := p.Called(input)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}

	return resp.(*dynamodb.DescribeTableOutput), args.Error(1)
}
