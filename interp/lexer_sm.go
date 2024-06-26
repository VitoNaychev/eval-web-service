package interp

import (
	"errors"
	"regexp"
	"strings"

	"github.com/VitoNaychev/eval-web-service/sm"
)

type LexerState int

const (
	stateLexerInitial LexerState = iota
	stateLexerTokenise
	stateLexerNonMathQuestion
	stateLexerUnsupportedOperation
	stateLexerEOF
)

type LexerEvent int

const (
	eventLexerSupportedToken LexerEvent = iota
	eventLexerUnsupportedToken
	eventLexerEOF
)

func tokeniseCallback(delta sm.Delta, ctx sm.Context) error {
	lexerCtx := ctx.(*LexerContext)

	for finder, tokenConstructor := range lexerCtx.FinderToConstructorMap {
		if value := finder.FindString(lexerCtx.Input); value != "" {
			lexerCtx.Tokens = append(lexerCtx.Tokens, tokenConstructor(value))

			lexerCtx.Input = lexerCtx.Input[len(value):]
			lexerCtx.Input = strings.TrimSpace(lexerCtx.Input)

			return nil
		}
	}

	return errors.New("event cannot be executed, invalid context")
}

func hasMathQuestion(delta sm.Delta, ctx sm.Context) (bool, error) {
	lexerCtx := ctx.(*LexerContext)

	for _, token := range lexerCtx.Tokens {
		if _, ok := token.(*QuestionToken); ok {
			return true, nil
		}
	}

	return false, nil
}

func unsupportedOperationCallback(delta sm.Delta, ctx sm.Context) error {
	return ErrUnsupportedOperation
}

func hasNotMathQuestion(delta sm.Delta, ctx sm.Context) (bool, error) {
	hasMathQuestion, err := hasMathQuestion(delta, ctx)
	return !hasMathQuestion, err
}

func nonMathQuestionCallback(delta sm.Delta, ctx sm.Context) error {
	return ErrNonMathQuestion
}

var lexerDeltas = []sm.Delta{
	{Current: sm.State(stateLexerTokenise), Event: sm.Event(eventLexerSupportedToken), Next: sm.State(stateLexerTokenise), Predicate: nil, Callback: tokeniseCallback},
	{Current: sm.State(stateLexerTokenise), Event: sm.Event(eventLexerEOF), Next: sm.State(stateLexerEOF), Predicate: hasMathQuestion, Callback: nil},
	{Current: sm.State(stateLexerTokenise), Event: sm.Event(eventLexerEOF), Next: sm.State(stateLexerNonMathQuestion), Predicate: hasNotMathQuestion, Callback: nonMathQuestionCallback},
	{Current: sm.State(stateLexerTokenise), Event: sm.Event(eventLexerUnsupportedToken), Next: sm.State(stateLexerNonMathQuestion), Predicate: hasNotMathQuestion, Callback: nonMathQuestionCallback},
	{Current: sm.State(stateLexerTokenise), Event: sm.Event(eventLexerUnsupportedToken), Next: sm.State(stateLexerUnsupportedOperation), Predicate: hasMathQuestion, Callback: unsupportedOperationCallback},
}

type PatternFinder interface {
	FindString(string) string
}

type LexerContext struct {
	Input                  string
	FinderToConstructorMap map[PatternFinder]NewTokenFunc

	Tokens []Token
}

func NewLexerContext(input string) *LexerContext {
	return &LexerContext{
		Input:                  input,
		FinderToConstructorMap: newFinderToConstructorMap(),

		Tokens: []Token{},
	}
}

func newFinderToConstructorMap() map[PatternFinder]NewTokenFunc {
	ftcMap := map[PatternFinder]NewTokenFunc{}

	questionRegex := regexp.MustCompile(QuestionTokenPattern)
	ftcMap[questionRegex] = NewQuestionToken

	numberRegex := regexp.MustCompile(NumberTokenPattern)
	ftcMap[numberRegex] = NewNumberToken

	operandRegex := regexp.MustCompile(OperandTokenPattern)
	ftcMap[operandRegex] = NewOperandToken

	punctuationRegex := regexp.MustCompile(PunctuationTokenPattern)
	ftcMap[punctuationRegex] = NewPunctuationToken

	return ftcMap
}
