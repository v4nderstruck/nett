package lang

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	EOF     = "EOF"
	Illegal = "illegal"

	Number     = "number"
	Identifier = "identifier"
	String     = "string"

	// IP
	CIDR = "cidr"

	// Operations
	Insert     = "insert"
	InsertInto = "into"
	Delete     = "delete"
	DeleteFrom = "from"

	// Delimiters
	LeftParen   = "("
	RightParen  = ")"
	Comma       = ","
	SingleQuote = "'"
	Semicolon   = ";"

	// Selectors
	Where  = "where"
	Values = "values"
)

func LoopupIdent(ident string) TokenType {
	switch ident {
	case "insert":
		return Insert
	case "into":
		return InsertInto
	case "delete":
		return Delete
	case "from":
		return DeleteFrom
	case "where":
		return Where
	case "values":
		return Values
	default:
		return Identifier
	}
}
