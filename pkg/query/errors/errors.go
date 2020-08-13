package errors

import "fmt"

// InvalidExpression is an error message when the syntax of an expression is not valid.
func InvalidExpression(conditionType string) error {
	if conditionType == "" {
		return fmt.Errorf("the provided expression is not a valid condition")
	}
	return fmt.Errorf("the provided expression is not a valid '%s' condition", conditionType)
}

// InvalidOffset is an error message when the offset lies outside of the constraints of the expression.
func InvalidOffset(offset int, expression string) error {
	return fmt.Errorf("the provided offset (%d) is not valid for the length of this expression (%d)", offset, len(expression))
}