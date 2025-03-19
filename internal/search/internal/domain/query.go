package domain

type QueryExpression struct {
}

type Equal struct {
	Field string
	Value []string
}

type NotEqual struct {
	Field string
	Value any
}
