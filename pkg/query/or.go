package query

import (
	"github.com/nmcclain/ldap"
	"github.com/shauncampbell/dapper/pkg/query/errors"
	"strings"
)

type Or struct {
	Conditions []Evaluator
}

func (o *Or) Evaluate(entry *ldap.Entry) bool {
	for _, condition := range o.Conditions {
		if condition.Evaluate(entry) {
			return true
		}
	}
	return false
}

// ToString produces a string version of this query condition
func (o *Or) ToString() string {
	out := "(|"
	for _, condition := range o.Conditions {
		out = out + condition.ToString()
	}

	out = out + ")"
	return out
}

// ParseOr takes an expression and attempts to parse it into an Or condition.
func ParseOr(expression string, offset int) (*Or, int, error) {
	if offset > len(expression) || offset < 0 {
		return nil, -1, errors.InvalidOffset(offset, expression)
	}

	if !strings.HasPrefix(expression[offset:], "(|") {
		return nil, -1, errors.InvalidExpression("or")
	}

	and := &Or{Conditions: make([]Evaluator, 0)}
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
		return nil, -1, errors.InvalidExpression("or")
	}

	if len(and.Conditions) == 0 {
		return nil, -1, errors.InvalidExpression("or")
	}

	return and, offset+1, nil
}