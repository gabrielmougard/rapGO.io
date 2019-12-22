package lists

import (
	"rapGO.io/src/algorithmservice/containers"
	"rapGO.io/src/algorithmservice/utils"
)

type List interface {
	Get(index int)(interface{})
	Remove(index int)
	Add(values ...interface{})
	Contains(values ...interface{}) bool
	Sort(comparator utils.Comparator)
	Swap(index1, index2 int)
	Insert(index int, values ...interface{})
	Set(index int, value interface{})

	containers.Container
}