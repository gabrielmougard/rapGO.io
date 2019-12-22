package linkedliststack

import (
	"fmt"
	"rapGO.io/src/algorithmservice/stacks"
	"rapGO.io/src/algorithmservice/lists/singlylinkedlist"
	"strings"
)

func assertStackImplementation() {
	var _ stacks.Stack = (*Stack)(nil)
}

type Stack struct {
	list *singlylinkedlist.List
}

//Instantiate empty stack
func New() *Stack {
	return &Stack{list: &singlylinkedlist.List{}}
}

//Add element onto the top of the stack
func (stack *Stack) Push(value interface{}) {
	stack.list.Prepend(value)
}

func (stack *Stack) Pop() (value interface{}){
	value = stack.list.Get(0)
	stack.list.Remove(0)
	return
}

func (stack *Stack) Peek() (value interface{}) {
	return stack.list.Get(0)
}

func (stack *Stack) Empty() bool {
	return stack.list.Empty()
}

func (stack *Stack) Size() int {
	return stack.list.Size()
}

func (stack *Stack) Clear() {
	stack.list.Clear()
}

//Return the elements of the stack in LIFO order
func (stack *Stack) Values() []interface{} {
	return stack.list.Values()
}

//String returns a string representation of container
func (stack *Stack) String() string {
	str := "LinkedListStack\n"
	values := []string{}
	for _, value := range stack.list.Values() {
		values = append(values, fmt.Sprintf("%v",value))
	}
	str += strings.Join(values, ", ")
	return str
}

// Check that the index is within the bounds of the list
func (stack *Stack) withinRange(index int) bool {
	return index >= 0 && index < stack.list.Size()
}

