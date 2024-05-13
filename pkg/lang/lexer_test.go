package lang

import (
	"testing"
)

func expectedNetwork(name, cidr4, cidr6 string) []Token {
	expect := []Token{
		{Type: Insert, Literal: "insert"},
		{Type: InsertInto, Literal: "into"},
		{Type: Identifier, Literal: "networks"},
		{Type: Values, Literal: "values"},
		{Type: LeftParen, Literal: "("},
		{Type: String, Literal: name},
		{Type: Comma, Literal: ","},
		{Type: CIDR, Literal: cidr4},
		{Type: Comma, Literal: ","},
		{Type: CIDR, Literal: cidr6},
		{Type: RightParen, Literal: ")"},
		{Type: Semicolon, Literal: ";"},
	}
	return expect
}

func TestNextTokenSingleLine(t *testing.T) {
	input_one_line := `insert into networks values ("VN_100", 8.8.8.8/24, 2001:10::/64);`
	expect := expectedNetwork("VN_100", "8.8.8.8/24", "2001:10::/64")
	l := NewLexer(input_one_line)
	for _, tok := range expect {
		actual := l.NextToken()
		if actual != tok {
			t.Errorf("Expected %v, got %v", tok, actual)
		}
	}
	if l.NextToken().Type != EOF {
		t.Errorf("Expected EOF, got %v", l.NextToken())
	}
}

func TestNextToekeMultiLine(t *testing.T) {
	input_multi_line := `
	   insert into networks values ("VN_100", 8.8.8.8/24, 2001:10:100:1::/64);
	   insert into networks values ("VN_101", 8.8.8.9/24, 2001:10:100:2::/64);
	   insert into networks values ("VN_102", 8.8.8.7/24, 2001:10:100:3::/64);
	 `
	expect1 := expectedNetwork("VN_100", "8.8.8.8/24", "2001:10:100:1::/64")
	expect2 := expectedNetwork("VN_101", "8.8.8.9/24", "2001:10:100:2::/64")
	expect3 := expectedNetwork("VN_102", "8.8.8.7/24", "2001:10:100:3::/64")

	l := NewLexer(input_multi_line)
	for _, tok := range expect1 {
		actual := l.NextToken()
		if actual != tok {
			t.Errorf("Expected %v, got %v", tok, actual)
		}
	}

	for _, tok := range expect2 {
		actual := l.NextToken()
		if actual != tok {
			t.Errorf("Expected %v, got %v", tok, actual)
		}
	}

	for _, tok := range expect3 {
		actual := l.NextToken()
		if actual != tok {
			t.Errorf("Expected %v, got %v", tok, actual)
		}
	}

	if l.NextToken().Type != EOF {
		t.Errorf("Expected EOF, got %v", l.NextToken())
	}
}

func TestAllTokens(t *testing.T) {
	all_token_input := "8.8.8.1/24 1992:10::/64 delete insert into from (), \"Hello\" ; where values abc"
  illegal_token := "1.2/23 1020:10:://64 1.2.3.4 10::1 1"

	expected_tokens := []Token{
		{Type: CIDR, Literal: "8.8.8.1/24"},
		{Type: CIDR, Literal: "1992:10::/64"},
		{Type: Delete, Literal: "delete"},
		{Type: Insert, Literal: "insert"},
		{Type: InsertInto, Literal: "into"},
		{Type: DeleteFrom, Literal: "from"},
		{Type: LeftParen, Literal: "("},
		{Type: RightParen, Literal: ")"},
		{Type: Comma, Literal: ","},
		{Type: String, Literal: "Hello"},
		{Type: Semicolon, Literal: ";"},
		{Type: Where, Literal: "where"},
		{Type: Values, Literal: "values"},
		{Type: Identifier, Literal: "abc"},
	}

	n := NewLexer(all_token_input)
	for _, tok := range expected_tokens {
		actual := n.NextToken()
		if actual != tok {
			t.Errorf("Expected %v, got %v", tok, actual)
		}
	}

	n = NewLexer(illegal_token)
	for tok := n.NextToken(); tok.Type != EOF; tok = n.NextToken() {
		if tok.Type != Illegal {
			t.Errorf("Expected illegal token, got %v", tok)
		}
	}
}
