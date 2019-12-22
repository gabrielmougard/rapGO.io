package linkedliststack

import "rapGO.io/src/algorithmservice/containers"

func assertSerializationImplementation() {
	var _ containers.JSONSerializer = (*Stack)(nil)
	var _ containers.JSONDeserializer = (*Stack)(nil)
}

// ToJSON outputs the JSON representation of the stack.
func (stack *Stack) ToJSON() ([]byte, error) {
	return stack.list.ToJSON()
}

// FromJSON populates the stack from the input JSON representation.
func (stack *Stack) FromJSON(data []byte) error {
	return stack.list.FromJSON(data)
}