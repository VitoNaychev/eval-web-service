package interp

import "github.com/VitoNaychev/eval-web-service/sm"

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
