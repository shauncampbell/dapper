package query

import (
	"github.com/nmcclain/ldap"
	"github.com/shauncampbell/dapper/pkg/query/errors"
	"strings"
)

// And conditions perform a binary AND operation on other conditions.
// In order for an AND condition to evaluate as true all conditions
// within it should evaluate as true also.
type And struct {
	Conditions []Evaluator	// An array of the conditions within the AND structure.
}

// Evaluate evaluates the query against the specified ldap entry.
func (a *And) Evaluate(entry *ldap.Entry) bool {
	for _, condition := range a.Conditions {
		if !condition.Evaluate(entry) {
			return false
		}
	}
	return true
}

// ToString produces a string version of this query condition
func (a *And) ToString() string {
	out := "(&"
	for _, condition := range a.Conditions {
		out = out + condition.ToString()
	}

	out = out + ")"
	return out
}

// ParseAnd takes an expression and attempts to parse it into an And condition.
func ParseAnd(expression string, offset int) (*And, int, error) {
	if offset > len(expression) || offset < 0 {
		return nil, -1, errors.InvalidOffset(offset, expression)
	}

	if !strings.HasPrefix(expression[offset:], "(&") {
		return nil, -1, errors.InvalidExpression("and")
	}

	and := &And{Conditions: make([]Evaluator, 0)}
	var cond Evaluator
	var err error
	offset = offset+2
	for offset < len(expression) -1 {
 		cond, offset, err = Parse(expression, offset)
		if err != nil {
			return nil, -1, err
		}
		and.Conditions = append(and.Conditions, cond)
	}

	if offset > len(expression)-1 || expression[offset] != ')' {
		return nil, -1, errors.InvalidExpression("and")
	}

	if len(and.Conditions) == 0 {
		return nil, -1, errors.InvalidExpression("and")
	}

	return and, offset+1, nil
}