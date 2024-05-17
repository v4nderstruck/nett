package lang

import "testing"

func contains(arr []Token, target Token) bool {
	for _, a := range arr {
		if a == target {
			return true
		}
	}
	return false
}

func TestInsertStatement(t *testing.T) {
	input_one_line := `insert into networks values ("VN_100", 9.0.0.0/8, 2001:10::/64);`

	n := NewLexer(input_one_line)
	i := NewInsertStatement()

	for {
		tok := n.NextToken()
		if tok.Type == EOF {
			break
		}
		next, err := i.ConsumeToken(tok)
		if err == nil {
			if next == CONSUME_DONE {
				break
			} else if next == CONSUME_MORE {
				continue
			}
		} else {
			t.Errorf("Unexpected error %v", err)
			break
		}
	}

	if i.State != CONSUME_FINAL {
		t.Errorf("Expected CONSUME_DONE state, got %v", i.State)
	}

	if i.TargetVal != "networks" {
		t.Errorf("Expected target networks, got %v", i.TargetVal)
	}

	if !contains(i.Values, Token{Type: String, Literal: "VN_100"}) {
		t.Errorf("Expected value VN_100, got %v", i.Values)
	}
	if !contains(i.Values, Token{Type: CIDR, Literal: "9.0.0.0/8"}) {
		t.Errorf("Expected value 9.0.0.0/8, got %v", i.Values)
	}

	if !contains(i.Values, Token{Type: CIDR, Literal: "2001:10::/64"}) {
		t.Errorf("Expected value 2001:10::/64, got %v", i.Values)
	}
}
