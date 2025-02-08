package parser

import (
	"fmt"
	"jsonParser/lexer"
	"jsonParser/token"
	"jsonParser/tree"
	"strconv"
)

type Parser struct {
	lexer lexer.Lexer
}

func New(data []byte) *Parser {
	return &Parser{
		lexer: *lexer.New(data),
	}
}

func (p *Parser) Parse() *tree.Node {
	t := p.lexer.NextToken()
	if t.Tp != token.LeftBrace {
		pan(t.Line, "Not JSON")
	}
	return p.parseJSON()
}

func (p *Parser) parseJSON() *tree.Node {
	return p.parseField(false)
}

func (p *Parser) parseField(continuation bool) *tree.Node {
	var node tree.Node

	t := p.lexer.NextToken()
	if t.Tp == token.RightBrace {
		return nil
	}

	if continuation {
		if t.Tp != token.Comma {
			pan(t.Line, "Not Comma")
		}
		t = p.lexer.NextToken()
	}

	if t.Tp != token.String {
		pan(t.Line, "Not Key")
	}
	node.Key = t.Literal

	t = p.lexer.NextToken()
	if t.Tp != token.Colon {
		pan(t.Line, "Not Colon")
	}

	t = p.lexer.NextToken()
	p.parseValue(&node, t)

	// fmt.Println(node)
	node.Next = p.parseField(true)
	return &node
}

func (p *Parser) parseArrayElements(continuation bool, index int) *tree.Node {
	var node tree.Node

	t := p.lexer.NextToken()
	if t.Tp == token.RightBracket {
		return nil
	}

	if continuation {
		if t.Tp != token.Comma {
			pan(t.Line, "Not Comma")
		}
		t = p.lexer.NextToken()
	}

	node.Key = strconv.Itoa(index)
	p.parseValue(&node, t)

	node.Next = p.parseArrayElements(true, index+1)
	return &node
}

func (p *Parser) parseValue(node *tree.Node, t token.Token) {
	switch t.Tp {
	case token.LeftBrace:
		node.ValueType = tree.JSON
		node.Value = p.parseJSON()
	case token.LeftBracket:
		node.ValueType = tree.Array
		node.Value = p.parseArrayElements(false, 0)
	case token.String:
		node.ValueType = tree.String
		node.Value = t.Literal
	case token.Number:
		node.ValueType = tree.Number
		node.Value = t.Literal
	case token.Bool:
		node.ValueType = tree.Bool
		node.Value = t.Literal
	case token.Null:
		node.ValueType = tree.Null
		node.Value = t.Literal
	default:
		pan(t.Line, "Not Value")
	}
}
func pan(line int, s string) {
	panic(fmt.Sprintf("%v: %s", line, s))
}
