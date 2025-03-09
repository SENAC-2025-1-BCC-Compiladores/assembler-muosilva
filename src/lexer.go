package src

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	SEC = "SECTION"
	EOF = "EOF"
	INS = "INSTRUCTION"
	NUM = "NUMBER"
	VAR = "VARIABLE"
	DEF = "DEFINE"
	UNK = "UNKNOWN"
)

var Codes = map[string]uint8{
	"NOP": 0x00,
	"STA": 0x10,
	"LDA": 0x20,
	"ADD": 0x30,
	"OR":  0x40,
	"AND": 0x50,
	"NOT": 0x60,
	"JMP": 0x80,
	"JN":  0x90,
	"JZ":  0xA0,
	"HLT": 0xF0,
}

var Types = map[string]bool{
	"DB":  true,
	"DS":  false,
	"EQU": true,
	"ORG": true,
}

var regex = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type Token struct {
	Type  string
	Value string
}

func isCode(piece string) bool {
	_, exist := Codes[piece]
	return exist
}

func isTypes(piece string) bool {
	_, exist := Types[piece]
	return exist
}

func isNumber(piece string) bool {
	_, err := strconv.ParseInt(piece, 16, 64)
	return err == nil
}

func isVariable(piece string) bool {
	return regex.MatchString(piece)
}

func defineToken(piece string) Token {
	switch {
	case strings.HasPrefix(piece, "."):
		return Token{Type: SEC, Value: strings.TrimPrefix(piece, ".")}
	case isTypes(piece):
		return Token{Type: DEF, Value: piece}
	case isCode(piece):
		return Token{Type: INS, Value: piece}
	case isNumber(piece):
		return Token{Type: NUM, Value: piece}
	case isVariable(piece):
		return Token{Type: VAR, Value: piece}
	default:
		log.Printf("Token desconhecido: %s", piece)
		return Token{Type: UNK, Value: piece}
	}
}

func Tokenization(filePath string) (tokens []Token) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo: %v", err)
	}

	lines := strings.Split(string(file), "\n")
	re := regexp.MustCompile(`\s+`)

	for _, line := range lines {
		line = strings.Split(line, ";")[0]
		pieces := re.Split(strings.TrimSpace(line), -1)

		for _, piece := range pieces {
			if piece != "" {
				tokens = append(tokens, defineToken(piece))
			}
		}
	}

	tokens = append(tokens, Token{Type: EOF, Value: ""})

	return tokens
}
