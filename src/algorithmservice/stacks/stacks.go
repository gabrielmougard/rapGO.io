package stacks

import "rapGO.io/src/algorithmservice/containers"

//Stack interface
// Push : add a value to the stack (using LIFO principle)
// Pop : return the last element in the stack and deleting this element from the stack
// Peek : same as Pop but do not modify the stack
type Stack interface {
	Push(value interface{})
	Pop() (value interface{})
	Peek() (value interface{})
	containers.Container
}

