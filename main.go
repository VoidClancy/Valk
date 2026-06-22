package main

import (
	"encoding/json"
	"fmt"
	"os"
	"valkyrie/parser"
)

func main() {

	file, err := os.ReadFile("schema.prisma")
	if err != nil {
		panic(err)
	}
	tokens := parser.ExtractTokens(string(file))
	parser.LogTokens(tokens)
	parser := parser.Parser{
		Tokens: tokens,
		Pos:    0,
	}
	b, _ := json.MarshalIndent(parser.Parse(), "", "  ")
	fmt.Println(string(b))
}
