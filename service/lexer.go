package service

import (
	"errors"
	"regexp"
	"strings"

	"github.com/VitoNaychev/eval-web-service/sm"
)

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

type LexerState int

const (
	STATE_INITIAL LexerState = iota
	STATE_TOKENISE
	STATE_UNSUPPORTED_TOKEN
	STATE_NON_MATH_QUESTION
	STATE_UNSUPPORTED_OPERATION
	STATE_FINAL
)

type LexerEvent int

const (
	EVENT_SUPPORTED_TOKEN LexerEvent = iota
	EVENT_UNSUPPORTED_TOKEN
	EVENT_PUNCTUATION_MARK
	EVENT_QUESTION
	EVENT_NO_QUESTION
)

func HasMathQuestion(delta sm.Delta, ctx sm.Context) (bool, error) {
	lexerCtx := ctx.(*LexerContext)

	for _, token := range lexerCtx.Tokens {
		if _, ok := token.(*QuestionToken); !ok {
			return true, nil
		}
	}

	return false, nil
}

func UnsupportedOperationCallback(delta sm.Delta, ctx sm.Context) error {
	return ErrUnsupportedOperation
}

func HasNotMathQuestion(delta sm.Delta, ctx sm.Context) (bool, error) {
	hasMathQuestion, err := HasMathQuestion(delta, ctx)
	return !hasMathQuestion, err
}

func NonMathQuestionCallback(delta sm.Delta, ctx sm.Context) error {
	return ErrNonMathQuestion
}

var deltas = []sm.Delta{
	{Current: sm.State(STATE_TOKENISE), Event: sm.Event(EVENT_SUPPORTED_TOKEN), Next: sm.State(STATE_TOKENISE), Predicate: nil, Callback: nil},
	{Current: sm.State(STATE_TOKENISE), Event: sm.Event(EVENT_PUNCTUATION_MARK), Next: sm.State(STATE_FINAL), Predicate: HasMathQuestion, Callback: nil},
	{Current: sm.State(STATE_TOKENISE), Event: sm.Event(EVENT_PUNCTUATION_MARK), Next: sm.State(STATE_NON_MATH_QUESTION), Predicate: HasNotMathQuestion, Callback: NonMathQuestionCallback},
	{Current: sm.State(STATE_TOKENISE), Event: sm.Event(EVENT_UNSUPPORTED_TOKEN), Next: sm.State(STATE_NON_MATH_QUESTION), Predicate: HasNotMathQuestion, Callback: NonMathQuestionCallback},
	{Current: sm.State(STATE_TOKENISE), Event: sm.Event(EVENT_UNSUPPORTED_TOKEN), Next: sm.State(STATE_UNSUPPORTED_OPERATION), Predicate: HasMathQuestion, Callback: UnsupportedOperationCallback},
}

type LexerContext struct {
	Tokens []Token
}

func Lex(input string) ([]Token, error) {
	ctx := LexerContext{}
	lexer := sm.New(sm.State(STATE_TOKENISE), deltas, &ctx)

	qre := regexp.MustCompile(`^What is`)
	nre := regexp.MustCompile(`^\d+`)
	ore := regexp.MustCompile(`^(plus|minus|multiplied by|divided by)`)
	pre := regexp.MustCompile(`^\?`)

	input = strings.TrimSpace(input)

	for len(input) > 0 {
		var err error

		if qre.MatchString(input) {
			err = lexer.Exec(sm.Event(EVENT_SUPPORTED_TOKEN))

			match := qre.FindString(input)
			ctx.Tokens = append(ctx.Tokens, &QuestionToken{Value: match})
			input = input[len(match):]
		} else if nre.MatchString(input) {
			err = lexer.Exec(sm.Event(EVENT_SUPPORTED_TOKEN))

			match := nre.FindString(input)
			ctx.Tokens = append(ctx.Tokens, &NumberToken{Value: match})
			input = input[len(match):]
		} else if ore.MatchString(input) {
			err = lexer.Exec(sm.Event(EVENT_SUPPORTED_TOKEN))

			match := ore.FindString(input)
			ctx.Tokens = append(ctx.Tokens, &OperandToken{Value: match})
			input = input[len(match):]
		} else if pre.MatchString(input) {
			err = lexer.Exec(sm.Event(EVENT_PUNCTUATION_MARK))

			match := pre.FindString(input)
			input = input[len(match):]
		} else {
			err = lexer.Exec(sm.Event(EVENT_UNSUPPORTED_TOKEN))
		}

		if err != nil {
			return nil, err
		}

		input = strings.TrimSpace(input)
	}
	if lexer.Current != sm.State(STATE_FINAL) {
		return nil, ErrMissingPunctuationMark
	}
	return ctx.Tokens, nil
}
