package service

type Token interface {
	GetToken() interface{}
}

type QuestionToken struct {
	Value string
}

func (q *QuestionToken) GetToken() interface{} {
	return q.Value
}

type NumberToken struct {
	Value string
}

func (n *NumberToken) GetToken() interface{} {
	return n.Value
}

type OperandToken struct {
	Value string
}

func (o *OperandToken) GetToken() interface{} {
	return o.Value
}
