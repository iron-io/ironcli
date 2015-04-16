package common

import (
	"encoding/json"
	"fmt"
	"os"
)

// FailErr prints the error and calls os.Exit(1) if err != nil
func FailErr(err error) {
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}

// PrintJSON encodes i into easily readable JSON using json.MarshalIndent
// and uses fmt.Println to print it to the console. if there was an error
// encoding the json, prints a short error message and calls fmt.Println(i)
func PrintJSON(i interface{}) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		fmt.Println("error marshaling the JSON ", err)
		fmt.Println(i)
		return
	}
	fmt.Println(string(b))
}
