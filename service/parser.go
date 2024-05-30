package service

import "github.com/VitoNaychev/eval-web-service/sm"

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

const (
	stateParserInitial = iota
	stateParserQuestion
	stateParserNumber
	stateParserOperand
	stateParserPunctuation
	stateParserFinal
	stateParserSyntaxError
)

type ParserEvent int

const (
	eventParserQuestion = iota
	eventParserNumber
	eventParserOperand
	eventParserPunctuation
	eventParserInvalid
)

func SignificantTokenCallback(delta sm.Delta, ctx sm.Context) error {
	parserCtx := ctx.(*ParserContext)

	currentToken := parserCtx.InputTokens[0]
	parserCtx.InputTokens = parserCtx.InputTokens[1:]

	parserCtx.OutputTokens = append(parserCtx.OutputTokens, currentToken)

	return nil
}

func NonsignificanTokenCallback(delta sm.Delta, ctx sm.Context) error {
	parserCtx := ctx.(*ParserContext)

	parserCtx.InputTokens = parserCtx.InputTokens[1:]

	return nil
}

func SyntaxErrorCallback(delta sm.Delta, ctx sm.Context) error {
	return ErrInvalidSyntax
}

var parserDeltas = []sm.Delta{
	{Current: sm.State(stateParserInitial), Event: sm.Event(eventParserQuestion), Next: sm.State(stateParserQuestion), Predicate: nil, Callback: NonsignificanTokenCallback},
	{Current: sm.State(stateParserInitial), Event: sm.Event(eventParserNumber), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},
	{Current: sm.State(stateParserInitial), Event: sm.Event(eventParserOperand), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},
	{Current: sm.State(stateParserInitial), Event: sm.Event(eventParserPunctuation), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},
	{Current: sm.State(stateParserInitial), Event: sm.Event(eventParserInvalid), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},

	{Current: sm.State(stateParserQuestion), Event: sm.Event(eventParserNumber), Next: sm.State(stateParserNumber), Predicate: nil, Callback: SignificantTokenCallback},
	{Current: sm.State(stateParserQuestion), Event: sm.Event(eventParserQuestion), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},
	{Current: sm.State(stateParserQuestion), Event: sm.Event(eventParserOperand), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},
	{Current: sm.State(stateParserQuestion), Event: sm.Event(eventParserPunctuation), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},
	{Current: sm.State(stateParserQuestion), Event: sm.Event(eventParserInvalid), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},

	{Current: sm.State(stateParserNumber), Event: sm.Event(eventParserOperand), Next: sm.State(stateParserOperand), Predicate: nil, Callback: SignificantTokenCallback},
	{Current: sm.State(stateParserNumber), Event: sm.Event(eventParserPunctuation), Next: sm.State(stateParserFinal), Predicate: nil, Callback: NonsignificanTokenCallback},
	{Current: sm.State(stateParserNumber), Event: sm.Event(eventParserQuestion), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},
	{Current: sm.State(stateParserNumber), Event: sm.Event(eventParserNumber), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},
	{Current: sm.State(stateParserNumber), Event: sm.Event(eventParserInvalid), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},

	{Current: sm.State(stateParserOperand), Event: sm.Event(eventParserNumber), Next: sm.State(stateParserNumber), Predicate: nil, Callback: SignificantTokenCallback},
	{Current: sm.State(stateParserOperand), Event: sm.Event(eventParserQuestion), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},
	{Current: sm.State(stateParserOperand), Event: sm.Event(eventParserOperand), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},
	{Current: sm.State(stateParserOperand), Event: sm.Event(eventParserPunctuation), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},
	{Current: sm.State(stateParserOperand), Event: sm.Event(eventParserInvalid), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},

	{Current: sm.State(stateParserFinal), Event: sm.Event(eventParserQuestion), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},
	{Current: sm.State(stateParserFinal), Event: sm.Event(eventParserNumber), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},
	{Current: sm.State(stateParserFinal), Event: sm.Event(eventParserOperand), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},
	{Current: sm.State(stateParserFinal), Event: sm.Event(eventParserPunctuation), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},
	{Current: sm.State(stateParserFinal), Event: sm.Event(eventParserInvalid), Next: sm.State(stateParserSyntaxError), Predicate: nil, Callback: SyntaxErrorCallback},
}

type ParserContext struct {
	InputTokens  []Token
	OutputTokens []Token
}

func Parse(tokens []Token) ([]Token, error) {
	ctx := ParserContext{
		InputTokens:  tokens,
		OutputTokens: []Token{},
	}
	parser := sm.New(stateParserInitial, parserDeltas, &ctx)

	for _, token := range tokens {
		var err error

		switch token.(type) {
		case *QuestionToken:
			err = parser.Exec(eventParserQuestion)
		case *NumberToken:
			err = parser.Exec(eventParserNumber)
		case *OperandToken:
			err = parser.Exec(eventParserOperand)
		case *PunctuationToken:
			err = parser.Exec(eventParserPunctuation)
		}

		if err != nil {
			return nil, err
		}
	}

	return ctx.OutputTokens, nil
}
