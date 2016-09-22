package lambda

import "testing"

func TestLambdaCreateFunction(t *testing.T) {
	lambdaCreateFunction := NewLambdaCreateFunction()
	lambdaCreateFunction.FunctionName = "test"
	lambdaCreateFunction.Runtime = "nodejs"
	lambdaCreateFunction.Handler = "runtime.function"
	lambdaCreateFunction.FileNames = []string{"../../tests/test.js"}

	err := lambdaCreateFunction.Action()
	if err != nil {
		t.Error(err)
	}
}

func TestLambdaTestFunction(t *testing.T) {
	lambdaTestFunction := NewLambdaTestFunction()
	lambdaTestFunction.FunctionName = "test"
	lambdaTestFunction.ClientContext = "client context"
	lambdaTestFunction.Payload = "test payload"

	err := lambdaTestFunction.Action()
	if err != nil {
		t.Error(err)
	}
}

func TestLambdaPublishFunction(t *testing.T) {
	lambdaPublishFunction := NewLambdaPublishFunction()
	lambdaPublishFunction.FunctionName = "test:latest"

	err := lambdaPublishFunction.Action()
	if err != nil {
		t.Error(err)
	}
}
