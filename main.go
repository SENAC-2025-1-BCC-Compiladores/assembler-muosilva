package main

import (
	"fmt"
	"log"
)

func main() {
	filePath := "code.txt"

	tokens := src.Tokenization(filePath)
	fmt.Println("Tokens:")
	for _, token := range tokens {
		fmt.Printf("{Type: %s, Value: %s}\n", token.Type, token.Value)
	}
	fmt.Println("--------------------")

	parser := src.NewParser(tokens)
	instructions, variables, err := parser.Parse()
	if err != nil {
		log.Fatalf("Erro ao analisar o código: %v", err)
	}

	err = parser.ResolveSymbols()
	if err != nil {
		log.Fatalf("Erro ao resolver símbolos: %v", err)
	}

	fmt.Println("Instruções:")
	for _, inst := range instructions {
		fmt.Printf("{Opcode: %s, Operand: %s}\n", inst.Opcode, inst.Operand)
	}
	fmt.Println("--------------------")

	fmt.Println("Variáveis:")
	for name, variable := range variables {
		fmt.Printf("{Name: %s, Value: %d, Type: %s}\n", name, variable.Value, variable.Type)
	}
}
