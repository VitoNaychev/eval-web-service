package interp

type ParserError struct {
	msg string
}

func NewParserError(msg string) *ParserError {
	return &ParserError{
		msg: msg,
	}
}

func (p *ParserError) Error() string {
	return p.msg
}

var (
	ErrInvalidSyntax = NewParserError("invalid syntax")
)
