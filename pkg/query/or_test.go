package query

import (
"github.com/onsi/gomega"
"github.com/shauncampbell/dapper/pkg/query/errors"
"testing"
)

func TestOrEvaluationHappyPath(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	not := Or{Conditions: []Evaluator{
		&Equals{
			Attribute: "uid",
			Value:     "person1",
		},
		&Equals{
			Attribute: "objectClass",
			Value: "inetOrgPerson",
		},
	}}

	Ω(not.Evaluate(&person1)).Should(gomega.Equal(true))
	Ω(not.Evaluate(&person2)).Should(gomega.Equal(true))
}

func TestOrEvaluationOneTrueOneFalse(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	or := Or{Conditions: []Evaluator{
		&Equals{
			Attribute: "uid",
			Value:     "person1",
		},
		&Equals{
			Attribute: "objectClass",
			Value: "badPerson",
		},
	}}

	Ω(or.Evaluate(&person1)).Should(gomega.Equal(true))
	Ω(or.Evaluate(&person2)).Should(gomega.Equal(false))
}

func TestOrEvaluationBothFalse(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	or := Or{Conditions: []Evaluator{
		&Equals{
			Attribute: "uid",
			Value:     "person6",
		},
		&Equals{
			Attribute: "objectClass",
			Value: "badPerson",
		},
	}}

	Ω(or.Evaluate(&person1)).Should(gomega.Equal(false))
	Ω(or.Evaluate(&person2)).Should(gomega.Equal(false))
}

func TestOrToString(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	or := Or{Conditions: []Evaluator{
		&Equals{
			Attribute: "uid",
			Value:     "person1",
		},
		&Equals{
			Attribute: "objectClass",
			Value: "badPerson",
		},
	}}

	Ω(or.ToString()).Should(gomega.Equal("(|(uid=person1)(objectClass=badPerson))"))
}

func TestOrParseHappyPath(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	expression := "(|(uid=person1)(objectClass=inetOrgPerson))"
	or, offset, err := Parse(expression, 0)
	Ω(err).Should(gomega.BeNil())
	Ω(offset).Should(gomega.Equal(len(expression)))

	Ω(or.Evaluate(&person1)).Should(gomega.Equal(true))
	Ω(or.Evaluate(&person2)).Should(gomega.Equal(true))
}

func TestOrParseMissingBracket(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	expression := "(|(uid=person1)(objectClass=inetOrgPerson)"
	_, _, err := ParseOr(expression, 0)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidExpression("or")))

	expression = "|(uid=person1)(objectClass=inetOrgPerson))"
	_, _, err = ParseOr(expression, 0)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidExpression("or")))
}

func TestOrParseSillyOffsets(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	expression := "(|(uid=person1)(objectClass=inetOrgPerson))"
	_, _, err := ParseOr(expression, -50)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidOffset(-50, expression)))

	expression = "(|(uid=person1)(objectClass=inetOrgPerson))"
	_, _, err = ParseOr(expression, 500)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidOffset(500, expression)))
}

func TestOrParseNoEquals(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	expression := "(|)"
	_, _, err := ParseOr(expression, 0)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidExpression("or")))
}


