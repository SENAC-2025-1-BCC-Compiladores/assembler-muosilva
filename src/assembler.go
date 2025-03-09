package src

import (
	"fmt"
	"strconv"
)

type Address uint8

type Instruction struct {
	Opcode  string
	Operand string
}

type Variable struct {
	Name  string
	Value Address
	Type  string
}

type Parser struct {
	tokens       []Token
	pos          int
	instructions []Instruction
	variables    map[string]Variable
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:       tokens,
		pos:          0,
		instructions: make([]Instruction, 0),
		variables:    make(map[string]Variable, 0),
	}
}

func (p *Parser) Parse() ([]Instruction, map[string]Variable, error) {
	for p.pos < len(p.tokens) {
		token := p.tokens[p.pos]
		err := p.parseToken(token)
		if err != nil {
			return nil, nil, err
		}
		p.pos++
	}
	return p.instructions, p.variables, nil
}

func (p *Parser) parseToken(token Token) error {
	switch token.Type {
	case INS:
		return p.parseInstruction()
	case VAR:
		return p.parseVariable()
	case SEC:
		return p.parseSection()
	case EOF:
		return nil
	default:
		return fmt.Errorf("token inesperado: %s", token.Value)
	}
}

func (p *Parser) parseInstruction() error {
	instruction := Instruction{Opcode: p.tokens[p.pos].Value}
	p.pos++

	if p.pos < len(p.tokens) && (p.tokens[p.pos].Type == NUM || p.tokens[p.pos].Type == VAR) {
		instruction.Operand = p.tokens[p.pos].Value
	} else {
		p.pos--
	}

	p.instructions = append(p.instructions, instruction)
	return nil
}

func (p *Parser) parseVariable() error {
	varName := p.tokens[p.pos].Value
	p.pos++
	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != DEF {
		return fmt.Errorf("esperado DB ou DS após o nome da variável")
	}

	directive := p.tokens[p.pos].Value
	p.pos++

	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != NUM {
		return fmt.Errorf("esperado valor após %s", directive)
	}

	value, err := strconv.ParseUint(p.tokens[p.pos].Value, 16, 8)
	if err != nil {
		return fmt.Errorf("valor inválido após %s: %s", directive, p.tokens[p.pos].Value)
	}

	p.variables[varName] = Variable{Name: varName, Value: Address(value), Type: directive}
	return nil
}

func (p *Parser) ResolveSymbols() error {
	address := Address(len(p.instructions))
	for name, variable := range p.variables {
		if variable.Type == "DS" {
			p.variables[name] = Variable{Name: name, Value: address, Type: "DS"}
			address++
		}
	}

	for i, inst := range p.instructions {
		if variable, exists := p.variables[inst.Operand]; exists && variable.Type == "DS" {
			p.instructions[i].Operand = fmt.Sprintf("%02X", variable.Value)
		}
	}

	return nil
}

func (p *Parser) parseSection() error {
	return nil
}
