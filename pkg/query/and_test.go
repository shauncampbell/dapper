package query

import (
	"github.com/onsi/gomega"
	"github.com/shauncampbell/dapper/pkg/query/errors"
	"testing"
)

func TestAndEvaluationHappyPath(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	not := And{Conditions: []Evaluator{
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
	Ω(not.Evaluate(&person2)).Should(gomega.Equal(false))
}

func TestAndEvaluationOneTrueOneFalse(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	and := And{Conditions: []Evaluator{
		&Equals{
			Attribute: "uid",
			Value:     "person1",
		},
		&Equals{
			Attribute: "objectClass",
			Value: "badPerson",
		},
	}}

	Ω(and.Evaluate(&person1)).Should(gomega.Equal(false))
	Ω(and.Evaluate(&person2)).Should(gomega.Equal(false))
}

func TestAndEvaluationBothFalse(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	and := And{Conditions: []Evaluator{
		&Equals{
			Attribute: "uid",
			Value:     "person1",
		},
		&Equals{
			Attribute: "objectClass",
			Value: "badPerson",
		},
	}}

	Ω(and.Evaluate(&person1)).Should(gomega.Equal(false))
	Ω(and.Evaluate(&person2)).Should(gomega.Equal(false))
}

func TestAndToString(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	and := And{Conditions: []Evaluator{
		&Equals{
			Attribute: "uid",
			Value:     "person1",
		},
		&Equals{
			Attribute: "objectClass",
			Value: "badPerson",
		},
	}}

	Ω(and.ToString()).Should(gomega.Equal("(&(uid=person1)(objectClass=badPerson))"))
}

func TestAndParseHappyPath(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	expression := "(&(uid=person1)(objectClass=inetOrgPerson))"
	and, offset, err := Parse(expression, 0)
	Ω(err).Should(gomega.BeNil())
	Ω(offset).Should(gomega.Equal(len(expression)))

	Ω(and.Evaluate(&person1)).Should(gomega.Equal(true))
	Ω(and.Evaluate(&person2)).Should(gomega.Equal(false))
}

func TestAndParseMissingBracket(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	expression := "(&(uid=person1)(objectClass=inetOrgPerson)"
	_, _, err := ParseAnd(expression, 0)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidExpression("and")))

	expression = "&(uid=person1)(objectClass=inetOrgPerson))"
	_, _, err = ParseAnd(expression, 0)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidExpression("and")))
}

func TestAndParseSillyOffsets(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	expression := "(&(uid=person1)(objectClass=inetOrgPerson))"
	_, _, err := ParseAnd(expression, -50)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidOffset(-50, expression)))

	expression = "(&(uid=person1)(objectClass=inetOrgPerson))"
	_, _, err = ParseAnd(expression, 500)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidOffset(500, expression)))
}

func TestAndParseNoEquals(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	expression := "(&)"
	_, _, err := ParseAnd(expression, 0)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidExpression("and")))
}

