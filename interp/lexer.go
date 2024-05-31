package interp

import (
	"fmt"
	"regexp"

	"github.com/VitoNaychev/eval-web-service/sm"
)

func Lex(input string) ([]Token, error) {
	ctx := LexerContext{
		Input:                  input,
		FinderToConstructorMap: newFinderToConstructorMap(),

		Tokens: []Token{},
	}
	lexer := sm.New(sm.State(stateLexerTokenise), lexerDeltas, &ctx)

	validPattern := fmt.Sprintf("%s|%s|%s|%s",
		QuestionTokenPattern, NumberTokenPattern, OperandTokenPattern, PunctuationTokenPattern)
	validRegex := regexp.MustCompile(validPattern)

	for len(ctx.Input) > 0 {
		var err error

		if validRegex.MatchString(ctx.Input) {
			err = lexer.Exec(sm.Event(eventLexerSupportedToken))
		} else {
			err = lexer.Exec(sm.Event(eventLexerUnsupportedToken))
		}

		if err != nil {
			return nil, err
		}
	}

	err := lexer.Exec(sm.Event(eventLexerEOF))
	if err != nil {
		return nil, err
	}

	return ctx.Tokens, nil
}
