package query

import (
	"github.com/nmcclain/ldap"
	"github.com/shauncampbell/dapper/pkg/query/errors"
	"strings"
)

type Evaluator interface {
	Evaluate(entry *ldap.Entry) bool // Evaluate evaluates the query against the specified ldap entry.
	ToString() string                // ToString produces a string version of this query condition
}

// Parse takes an expression and an offset within that expression and turns it into
// an Evaluator chain. This can then be used to evaluate a query against a set of
// ldap entries.
func Parse(expression string, offset int) (Evaluator, int, error) {
	if strings.HasPrefix(expression[offset:], "(&") {
		return ParseAnd(expression, offset)
	} else if strings.HasPrefix(expression[offset:], "(|") {
		return ParseOr(expression, offset)
	} else if strings.HasPrefix(expression[offset:], "(!") {
		return ParseNot(expression, offset)
	} else if strings.HasPrefix(expression[offset:], "(") {
		return ParseEquals(expression, offset)
	}
	return nil, -1, errors.InvalidExpression("")
}
