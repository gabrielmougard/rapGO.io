package kafka

//Describes an interface for a class that emits events
type EventEmitter interface {
	Emit(e Event) error
}