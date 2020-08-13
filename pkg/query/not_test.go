package query

import (
	"github.com/onsi/gomega"
	"github.com/shauncampbell/dapper/pkg/query/errors"
	"testing"
)

func TestNotEvaluationHappyPath(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	not := Not{Parent: &Equals{
		Attribute: "uid",
		Value:     "person1",
	}}

	Ω(not.Evaluate(&person1)).Should(gomega.Equal(false))
	Ω(not.Evaluate(&person2)).Should(gomega.Equal(true))
}

func TestNotEvaluationWildcard(t *testing.T) {
	gomega.RegisterTestingT(t)

	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	not := Not{Parent: &Equals{
		Attribute: "uid",
		Value:     "person*",
	}}

	Ω(not.Evaluate(&person1)).Should(gomega.Equal(false))
	Ω(not.Evaluate(&person2)).Should(gomega.Equal(false))
}

func TestNotEvaluationNonEmpty(t *testing.T) {
	gomega.RegisterTestingT(t)

	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	not := Not{Parent: &Equals{
		Attribute: "userPassword",
		Value:     "*",
	}}

	Ω(not.Evaluate(&person1)).Should(gomega.Equal(false))
	Ω(not.Evaluate(&person2)).Should(gomega.Equal(true))
}

func TestNotToString(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	not := Not{ Parent: &Equals{
		Attribute: "uid",
		Value:     "person1",
	}}

	Ω(not.ToString()).Should(gomega.Equal("(!(uid=person1))"))
}

func TestNotParseHappyPath(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	expression := "(!(uid=person1))"
	not, offset, err := Parse(expression, 0)
	Ω(err).Should(gomega.BeNil())
	Ω(offset).Should(gomega.Equal(len(expression)))

	Ω(not.Evaluate(&person1)).Should(gomega.Equal(false))
	Ω(not.Evaluate(&person2)).Should(gomega.Equal(true))
}

func TestNotParseMissingBracket(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	expression := "(!(uid=person1)"
	_, _, err := ParseNot(expression, 0)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidExpression("not")))

	expression = "!(uid=person1))"
	_, _, err = ParseNot(expression, 0)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidExpression("not")))
}

func TestNotParseSillyOffsets(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	expression := "(!(uid=person1))"
	_, _, err := ParseNot(expression, -50)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidOffset(-50, expression)))

	expression = "(!(uid=person1))"
	_, _, err = ParseNot(expression, 500)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidOffset(500, expression)))
}

func TestNotParseNoEquals(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	expression := "(!)"
	_, _, err := ParseNot(expression, 0)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidExpression("")))
}
