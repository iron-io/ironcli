package common

import (
	"encoding/json"
	"fmt"
)

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
