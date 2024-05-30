package service

type Token interface {
	GetToken() interface{}
}

type NewTokenFunc func(string) Token

var QuestionTokenPattern = `^What is`

type QuestionToken struct {
	Value string
}

func NewQuestionToken(value string) Token {
	return &QuestionToken{
		Value: value,
	}
}

func (q *QuestionToken) GetToken() interface{} {
	return q.Value
}

var NumberTokenPattern = `^\d+`

type NumberToken struct {
	Value string
}

func NewNumberToken(value string) Token {
	return &NumberToken{
		Value: value,
	}
}

func (n *NumberToken) GetToken() interface{} {
	return n.Value
}

var OperandTokenPattern = `^(plus|minus|multiplied by|divided by)`

type OperandToken struct {
	Value string
}

func NewOperandToken(value string) Token {
	return &OperandToken{
		Value: value,
	}
}

func (o *OperandToken) GetToken() interface{} {
	return o.Value
}

var PunctuationTokenPattern = `^\?`

type PunctuationToken struct {
	Value string
}

func NewPunctuationToken(value string) Token {
	return &PunctuationToken{
		Value: value,
	}
}

func (p *PunctuationToken) GetToken() interface{} {
	return p.Value
}
