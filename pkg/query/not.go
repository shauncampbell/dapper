package query

import (
	"github.com/nmcclain/ldap"
	"github.com/shauncampbell/dapper/pkg/query/errors"
	"strings"
)

type Not struct {
	Parent Evaluator
	Evaluator
}

// Evaluate evaluates the query against the specified ldap entry.
func (n *Not) Evaluate(entry *ldap.Entry) bool {
	return !n.Parent.Evaluate(entry)
}

// ToString produces a string version of this query condition
func (n *Not) ToString() string {
	return "(!" + n.Parent.ToString() + ")"
}

func ParseNot(expression string, offset int) (*Not, int, error) {
	if offset > len(expression) || offset < 0 {
		return nil, -1, errors.InvalidOffset(offset, expression)
	}

	if !strings.HasPrefix(expression[offset:], "(!") {
		return nil, -1, errors.InvalidExpression("not")
	}

	eq, pos, err := Parse(expression, offset+2)
	if err != nil {
		return nil, -1, err
	}

	if pos > len(expression)-1 || expression[pos] != ')' {
		return nil, -1, errors.InvalidExpression("not")
	}

	return &Not{Parent: eq}, pos+1, nil
}