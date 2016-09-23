package lambda

import "testing"

func TestLambdaCreateFunction(t *testing.T) {
	lambdaCreateFunction := NewLambdaCreateFunction()
	lambdaCreateFunction.FunctionName = "irontesting/test"
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
	lambdaTestFunction.FunctionName = "irontesting/test"

	err := lambdaTestFunction.Action()
	if err != nil {
		t.Error(err)
	}
}

func TestLambdaPublishFunction(t *testing.T) {
	lambdaPublishFunction := NewLambdaPublishFunction()
	lambdaPublishFunction.FunctionName = "irontesting/test:latest"

	err := lambdaPublishFunction.Action()
	if err != nil {
		t.Error(err)
	}
}
