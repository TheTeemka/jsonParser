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

func (p *Parser) Parse() (*tree.Node, error) {
	t := p.lexer.NextToken()
	if t.Tp == token.ILLEGAL {
		return nil, fmt.Errorf("invalid token at %v from %v to %v", t.Line+1, t.Start, t.End)
	}
	if t.Tp != token.LeftBrace {
		return nil, fmt.Errorf("invalid json format (not starts with '{')")
	}
	return p.parseJSON()
}

func (p *Parser) parseJSON() (*tree.Node, error) {
	return p.parseField(false)
}

func (p *Parser) parseField(continuation bool) (*tree.Node, error) {
	var node tree.Node

	t := p.lexer.NextToken()
	if t.Tp == token.RightBrace {
		return nil, nil
	}

	if continuation {
		if t.Tp != token.Comma {
			return nil, fmt.Errorf("there is not comma in array at %v", t.Line)
		}
		t = p.lexer.NextToken()
	}

	if t.Tp != token.String {
		return nil, fmt.Errorf("field is not string at %v", t.Line)
	}
	node.Key = t.Literal

	t = p.lexer.NextToken()
	if t.Tp != token.Colon {
		return nil, fmt.Errorf("there is not colon after field at %v", t.Line)
	}

	t = p.lexer.NextToken()
	err := p.parseValue(&node, t)
	if err != nil {
		return nil, err
	}
	// fmt.Println(node)
	node.Next, err = p.parseField(true)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (p *Parser) parseArrayElements(continuation bool, index int) (*tree.Node, error) {
	var node tree.Node

	t := p.lexer.NextToken()
	if t.Tp == token.RightBracket {
		return nil, nil
	}

	if continuation {
		if t.Tp != token.Comma {
			return nil, fmt.Errorf("there is not comma in array at %v", t.Line)
		}
		t = p.lexer.NextToken()
	}

	node.Key = strconv.Itoa(index)
	err := p.parseValue(&node, t)
	if err != nil {
		return nil, err
	}
	node.Next, err = p.parseArrayElements(true, index+1)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (p *Parser) parseValue(node *tree.Node, t token.Token) error {
	var err error
	switch t.Tp {
	case token.LeftBrace:
		node.ValueType = tree.JSON
		node.Value, err = p.parseJSON()
	case token.LeftBracket:
		node.ValueType = tree.Array
		node.Value, err = p.parseArrayElements(false, 0)
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
		err = fmt.Errorf("invalid value type at %v", t.Line)
	}
	return err
}
