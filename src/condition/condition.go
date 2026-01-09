package condition

import (
	_ "fmt"
	_ "log"
)

// Condition represents the condition state of a character
type Condition string

// Create a new Condition instance
func NewCondition(cond string) Condition {
	return Condition(cond)
}

// String returns the string representation of the Condition
func (c Condition) String() string {
	return string(c)
}
