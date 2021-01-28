package registry

import (
	"context"
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoMock struct {
	mockError error
}

func (m *mongoMock) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, m.mockError
}

func (m *mongoMock) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return nil, m.mockError
}

func TestSaveSchema(t *testing.T) {
	errorNotSave := errors.New("not able to save schema")
	testCases := []struct {
		name      string
		mockError error
		want      error
	}{
		{"Should save schema correctly", nil, nil},
		{"Should throw error when saving schema", errorNotSave, errorNotSave},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			MongoDB = mongoDB{collection: &mongoMock{mockError: tc.mockError}}
			schema := SchemaRegistry{
				ServiceName: "Cart",
				ServiceURL:  "http://cart-service",
				TypeDefs:    "type Query { placeHolder: String }",
			}
			if got := MongoDB.saveSchema(&schema); got != tc.want {
				t.Fail()
			}
		})
	}
}

func TestGetServiceSchemas(t *testing.T) {
	errorNotSave := errors.New("not able to get schema")
	testCases := []struct {
		name      string
		mockError error
		want      error
	}{
		{"Should throw error when getting service schema", errorNotSave, errorNotSave},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			MongoDB = mongoDB{collection: &mongoMock{mockError: tc.mockError}}
			results := []string{}
			if got := MongoDB.getServiceSchemas(&results); got != tc.want {
				t.Fail()
			}
		})
	}
}
