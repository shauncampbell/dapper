package query

import (
	"fmt"
	"github.com/nmcclain/ldap"
	"github.com/shauncampbell/dapper/pkg/query/errors"
	regexp "regexp"
	"strings"
)

type Equals struct {
	Attribute string
	Value     string
	Evaluator
}

func (e *Equals) Evaluate(entry *ldap.Entry) bool {
	for _, a := range entry.Attributes {
		if strings.ToLower(a.Name) == e.Attribute {

			q := strings.ReplaceAll(e.Value, ".", "\\.")
			q = strings.ReplaceAll(q, "*", ".*")

			r, err := regexp.CompilePOSIX(q)
			if err != nil {
				return false
			}

			for _, v := range a.Values {

				if r.MatchString(strings.ToLower(v)) {
					return true
				}
			}
		}
	}
	return false
}

// ToString produces a string version of this query condition
func (e *Equals) ToString() string {
	return fmt.Sprintf("(%s=%s)", e.Attribute, e.Value)
}

func ParseEquals(expression string, offset int) (*Equals, int, error) {
	if offset > len(expression) || offset < 0 {
		return nil, -1, errors.InvalidOffset(offset, expression)
	}

	if !strings.HasPrefix(expression[offset:], "(") {
		return nil, -1, errors.InvalidExpression("equals")
	}

	firstBracket := strings.Index(expression[offset:], ")")+offset
	if firstBracket < offset+2 {
		return nil, -1, errors.InvalidExpression("equals")
	}

	equalsExpression := expression[offset+1:firstBracket]

	if eqSymbolIndex := strings.Index(equalsExpression, "="); eqSymbolIndex != -1 {
		attribute := equalsExpression[:eqSymbolIndex]
		expression := equalsExpression[eqSymbolIndex+1 : ]
		return &Equals{Attribute: strings.ToLower(attribute), Value: strings.ToLower(expression)}, firstBracket+1, nil
	}

	return nil, -1, errors.InvalidExpression("equals")
}
