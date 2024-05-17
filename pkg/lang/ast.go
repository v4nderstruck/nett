package lang

import "fmt"

const (
	CONSUME_MORE = 1
	CONSUME_DONE = 0
  CONSUME_FINAL = -1000
)

type Expression interface {
	PrefixToken() Token
	ToString() string
	ConsumeToken(Token) error
}

func returnError(got, expected string) error {
	return fmt.Errorf("encoutered unpexpected character '%s', expected '%s'", got, expected)
}

type InsertStatement struct {
	TargetVal string
	Values    []Token // (....)
	State     int
}

func (i *InsertStatement) PrefixToken() Token {
	return Token{Type: Insert, Literal: "insert"}
}

func (i *InsertStatement) ToString() string {
	s := fmt.Sprintf("insert into %s values ( ", i.TargetVal)
	for k, t := range i.Values {
		if k == len(i.Values)-1 {
			s += t.Literal
		} else {
			s += t.Literal + ", "
		}
	}
	s += " )"
	return s
}

func (i *InsertStatement) ConsumeToken(t Token) (int, error) {
	// 0: Insert, 1: InsertInto, 2: Identifier, 3: LeftParen, 8: Values
	// 4: String, CIDR, 5: Comma, RightParen, 6: Semicolon, 7: EOS
	switch i.State {
	case 0:
		if t.Type != Insert {
			i.State = -1
			return CONSUME_DONE, returnError(t.Literal, Insert)
		} else {
			i.State = 1
		}
	case 1:
		if t.Type != InsertInto {
			i.State = -1
			return CONSUME_DONE, returnError(t.Literal, InsertInto)
		} else {
			i.State = 2
		}
	case 2:
		if t.Type != Identifier {
			i.State = -1
			return CONSUME_DONE, returnError(t.Literal, Identifier)
		} else {
			i.TargetVal = t.Literal
			i.State = 8
		}
	case 8:
		if t.Type != Values {
			i.State = -1
			return CONSUME_DONE, returnError(t.Literal, Values)
		} else {
			i.State = 3
		}
	case 3:
		if t.Type != LeftParen {
			i.State = -1
			return CONSUME_DONE, returnError(t.Literal, LeftParen)
		} else {
			i.State = 4
		}
	case 4:
		if t.Type != String && t.Type != CIDR {
			i.State = -1
			return CONSUME_DONE, returnError(t.Literal, "String or CIDR")
		} else {
			i.Values = append(i.Values, t)
			i.State = 5
		}
	case 5:
		if t.Type != Comma && t.Type != RightParen {
			i.State = -1
			return CONSUME_DONE, returnError(t.Literal, "',' or ')'")
		} else {
			switch t.Type {
			case RightParen:
				i.State = 6
			case Comma:
				i.State = 4
			}
		}
	case 6:
		if t.Type != Semicolon {
			i.State = -1
			return CONSUME_DONE, returnError(t.Literal, Semicolon)
		} else {
			i.State = CONSUME_FINAL 
		}
	case CONSUME_FINAL:
		return CONSUME_DONE, nil
	default:
		return CONSUME_DONE, fmt.Errorf("ast invalid state")
	}
	return CONSUME_MORE, nil
}

func NewInsertStatement() *InsertStatement {
	return &InsertStatement{
		TargetVal: "",
		Values:    []Token{},
		State:     0,
	}
}
