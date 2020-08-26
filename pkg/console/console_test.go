package console

import (
	"github.com/onsi/gomega"
	"testing"
)

type Hello struct {
	Forename   string
	Surname    string
	Age        int
	IsAccurate bool
}

func Test(t *testing.T) {
	gomega.RegisterTestingT(t)
	// workaround as our naming clashes with gomega.
	Ω := gomega.Ω

	// Ensure that the string length matches. We are using a string match
	// here because the map field order is non-deterministic.
	b, err := Marshal(map[string]string{"Forename": "John", "Surname": "Doe"})
	Ω(err).Should(gomega.BeNil())
	Ω(len(string(b))).Should(gomega.Equal(len(`FORENAME     SURNAME     
John         Doe         
`)))

	// Ensure that the string length matches. We are using a string match
	// here because the map field order is non-deterministic.
	b, err = Marshal(map[string]int{"Forename": 123, "Surname": 234})
	Ω(err).Should(gomega.BeNil())
	Ω(len(string(b))).Should(gomega.Equal(len(`FORENAME     SURNAME     
123          234         
`)))

	// Ensure that the string length matches. We are using a string match
	// here because the map field order is non-deterministic.
	b, err = Marshal([]map[string]string{{"Forename": "John", "Surname": "Doe", "Age": "5"}, {"Forename": "Jane", "Surname": "Doe"}})
	Ω(err).Should(gomega.BeNil())
	Ω(len(string(b))).Should(gomega.Equal(len(`FORENAME     SURNAME     AGE     
John         Doe         5       
Jane         Doe                 
`)))

	b, err = Marshal(Hello{Forename: "John"})
	Ω(err).Should(gomega.BeNil())
	Ω(string(b)).Should(gomega.Equal(`FORENAME     SURNAME     AGE     ISACCURATE     
John                     0       false          
`))

	b, err = Marshal([]Hello{
		{Forename: "John"},
		{Surname: "Doe"},
		{Forename: "Jane", Surname: "Doe", Age: 55},
	})
	Ω(err).Should(gomega.BeNil())
	Ω(string(b)).Should(gomega.Equal(`FORENAME     SURNAME     AGE     ISACCURATE     
John                     0       false          
             Doe         0       false          
Jane         Doe         55      false          
`))
	b, err = Marshal(&Hello{Forename: "John"})
	Ω(err).Should(gomega.BeNil())
	Ω(string(b)).Should(gomega.Equal(`FORENAME     SURNAME     AGE     ISACCURATE     
John                     0       false          
`))
}
