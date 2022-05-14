package dbresolver

import "gorm.io/gorm/clause"

type Operation string

const writeName = "gorm:db_resolver:write"

func (op Operation) Name() string {
	if op == Write {
		return writeName
	} else {
		return "gorm:db_resolver:read"
	}
}

func (op Operation) Build(clause.Builder) {
}

func (op Operation) MergeClause(*clause.Clause) {
}

func Use(str string) clause.Interface {
	return using{Use: str}
}

type using struct {
	Use string
}

const usingName = "gorm:db_resolver:using"

func (u using) Name() string {
	return usingName
}

func (u using) Build(clause.Builder) {
}

func (u using) MergeClause(c *clause.Clause) {
	c.Expression = u
}
