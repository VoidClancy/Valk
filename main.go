package main

import (
	"encoding/json"
	"fmt"
	"os"
	"valkyrie/schema"
)

func main() {

	fileBytes, _ := os.ReadFile("schema.prisma")
	rawString := string(fileBytes)

	// tokens := schema.ExtractTokens(rawString)
	// for i, token := range tokens {
	// 	fmt.Printf("%d: %+v\n", i, token)
	// }

	schema, errs := schema.ParseSchema(rawString)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}

	}

	b, _ := json.MarshalIndent(schema, "", "  ")

	os.WriteFile("result.json", b, 0644)
	fmt.Println(string(b))
}
