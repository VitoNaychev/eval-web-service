package interp

import "errors"

type LexerError struct {
	msg string
}

func NewLexerError(msg string) *LexerError {
	return &LexerError{
		msg: msg,
	}
}

func (l *LexerError) Error() string {
	return l.msg
}

var (
	ErrNonMathQuestion        = errors.New("non-math question")
	ErrUnsupportedOperation   = errors.New("unuspported operation")
	ErrMissingPunctuationMark = errors.New("missing punctuation mark")
)
