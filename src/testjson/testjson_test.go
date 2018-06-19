package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

type testmap struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type testmap2 struct {
	SName string `json:"name"` //접근제한자 말고 변수명 자체가 다르면 ""로 묶어 주어야 한다.
	SAge  int    `json:"age"`
}

func TestTestMap(test *testing.T) {
	doc := `	{		"name": "maria",		"age" : 18	}	`

	var data map[string]interface{}

	json.Unmarshal([]byte(doc), &data)
	fmt.Println("- Unmarshal - 1 - \n", data)
	fmt.Println()

	tm := testmap{}
	json.Unmarshal([]byte(doc), &tm)
	fmt.Printf("- Unmarshal - 2 - \n %+v \n", tm)
	fmt.Println()

	tm2 := testmap2{}
	json.Unmarshal([]byte(doc), &tm2)
	fmt.Printf("- Unmarshal - 3 - \n %+v \n", tm2)
	fmt.Println()

	//doc2, err := json.Marshal(data)
	doc2, err := json.MarshalIndent(data, "", "   ")
	if err == nil {
		fmt.Println("- Marshal - \n", string(doc2))
		fmt.Println()
	}

}
