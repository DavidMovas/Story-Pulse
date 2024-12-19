package dbx

import (
	"github.com/Masterminds/squirrel"
	"github.com/lann/builder"
)

type StatementBuilder struct {
	squirrel.StatementBuilderType
}

func NewStatementBuilder() *StatementBuilder {
	return &StatementBuilder{
		squirrel.StatementBuilderType(builder.EmptyBuilder).PlaceholderFormat(squirrel.Dollar),
	}
}

func (b *StatementBuilder) Build() (sql string, args []interface{}, err error) {
	return b.Build()
}
