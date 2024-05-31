package interp

import "github.com/VitoNaychev/eval-web-service/sm"

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

	if parser.Current != stateParserFinal {
		return nil, ErrInvalidSyntax
	}

	return ctx.OutputTokens, nil
}
