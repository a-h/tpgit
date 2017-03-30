package backend

import "github.com/aws/aws-sdk-go/service/dynamodb"

// Dynamo is a backend which stores a history of processed git commits in DynamoDB.
type Dynamo struct {
	DB *dynamodb.DynamoDB
}

// NewDynamo creates a new DynamoDB backend.
func NewDynamo() Dynamo {
	return Dynamo{}
}
