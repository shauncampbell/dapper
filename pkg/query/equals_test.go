package query

import (
	"github.com/nmcclain/ldap"
	"github.com/onsi/gomega"
	"github.com/shauncampbell/dapper/pkg/query/errors"
	"testing"
)

var person1 = ldap.Entry{DN: "cn=person1,dc=test,dc=lab",
	Attributes: []*ldap.EntryAttribute{
		{ Name:   "uid", Values: []string{ "person1" } },
		{ Name: "objectClass", Values: []string{ "inetOrgPerson", "jellyfinUser" }},
		{ Name: "userPassword", Values: []string{ "{SSHA}I8wq1+4gyJVJUtQW96JGcmCL46ADyPnW" }},	// password: test
	},
}

var person2 = ldap.Entry{DN: "cn=person2,dc=test,dc=lab",
	Attributes: []*ldap.EntryAttribute{
		{ Name:   "uid", Values: []string{ "person2" } },
		{ Name: "objectClass", Values: []string{ "inetOrgPerson" }},
	},
}

func TestEqualsEvaluationHappyPath(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	equals := Equals{
		Attribute: "uid",
		Value:     "person1",
	}

	Ω(equals.Evaluate(&person1)).Should(gomega.Equal(true))
	Ω(equals.Evaluate(&person2)).Should(gomega.Equal(false))
}

func TestEqualsEvaluationWildcard(t *testing.T) {
	gomega.RegisterTestingT(t)

	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	equals := Equals{
		Attribute: "uid",
		Value:     "person*",
	}

	Ω(equals.Evaluate(&person1)).Should(gomega.Equal(true))
	Ω(equals.Evaluate(&person2)).Should(gomega.Equal(true))
}

func TestEqualsEvaluationNonEmpty(t *testing.T) {
	gomega.RegisterTestingT(t)

	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	equals := Equals{
		Attribute: "userPassword",
		Value:     "*",
	}

	Ω(equals.Evaluate(&person1)).Should(gomega.Equal(true))
	Ω(equals.Evaluate(&person2)).Should(gomega.Equal(false))
}

func TestEqualsToString(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	equals := Equals{
		Attribute: "uid",
		Value:     "person1",
	}

	Ω(equals.ToString()).Should(gomega.Equal("(uid=person1)"))
}

func TestEqualsParseHappyPath(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	expression := "(uid=person1)"
	equals, offset, err := Parse(expression, 0)
	Ω(err).Should(gomega.BeNil())
	Ω(offset).Should(gomega.Equal(len(expression)))

	Ω(equals.Evaluate(&person1)).Should(gomega.Equal(true))
	Ω(equals.Evaluate(&person2)).Should(gomega.Equal(false))
}

func TestEqualsParseMissingBracket(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	expression := "uid=person1)"
	_, _, err := ParseEquals(expression, 0)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidExpression("equals")))

	expression = "(uid=person1"
	_, _, err = ParseEquals(expression, 0)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidExpression("equals")))
}

func TestEqualsParseSillyOffsets(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	expression := "(uid=person1)"
	_, _, err := ParseEquals(expression, -50)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidOffset(-50, expression)))

	expression = "(uid=person1)"
	_, _, err = ParseEquals(expression, 500)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidOffset(500, expression)))
}

func TestEqualsParseNoEquals(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	expression := "(uid)"
	_, _, err := ParseEquals(expression, 0)
	Ω(err).ShouldNot(gomega.BeNil())
	Ω(err).Should(gomega.Equal(errors.InvalidExpression("equals")))
}
