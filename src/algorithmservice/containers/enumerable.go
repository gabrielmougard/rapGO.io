package containers

type EnumerableWithIndex interface {
	// Each calls the given function once for ordered containers whose values can be fetched by an index
	Each(func(index int, value interface{}))
	// Any passes each element of the container to the given function and
	// returns true if the function ever returns true for any element 
	Any(func(index int, value interface{}) bool) bool
	// All passes each element of the container to the given function and
	// returns true if the function returns true for all element
	All(func(index int, value interface{}) bool) bool
	// Find passes each element to the container to the given function and returns
	// the first (index,value) for which the function is true or -1, nil otherwise
	// if no element matches the criteria.
	Find(func(index int, value interface{}) bool) (int, interface{})
}

//EnumerableWithKey provides functions for ordered containers whose values whose elements are key/values pairs.
type EnumerableWithKey	interface {
	// Each calls the given function once for each element, passing that element's key and value
	Each(func(key interface{}, value interface{}))
	// Any passes each element of the container to the given function and
	// returns true if the function ever returns true for any element 
	Any(func(key interface{}, value interface{}) bool) bool
	// All passes each element of the container to the given function and
	// returns true if the function returns true for all element
	All(func(key interface{}, value interface{}) bool) bool
	// Find passes each element to the container to the given function and returns
	// the first (key,value) for which the function is true or -1, nil otherwise
	// if no element matches the criteria.
	Find(func(key interface{}, value interface{}) bool) (interface{}, interface{})

}