package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	testMap()
}

type testmap struct {
	Name string `json:name`
	Age  int    `json:age`
}

func testMap() {
	doc := `
	{
		"name": "maria",
		"age" : 18
	}
	`

	var data map[string]interface{}

	json.Unmarshal([]byte(doc), &data)
	fmt.Println("- Unmarshal - 1 - \n", data)
	fmt.Println()

	tm := testmap{}
	json.Unmarshal([]byte(doc), &tm)
	fmt.Printf("- Unmarshal - 2 - \n %+v \n", tm)
	fmt.Println()

	//doc2, err := json.Marshal(data)
	doc2, err := json.MarshalIndent(data, "", "   ")
	if err == nil {
		fmt.Println("- Marshal - \n", string(doc2))
		fmt.Println()
	}

}
