package states

import (
	"testing"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	
)

func TestTree(t *testing.T) {
	tree := NewRbTree()
	var save string
	for i := 0; i < 1000; i++ {
		uuid := uuid.New().String()
		key := StringKey(uuid)
		if i == 100 {
			fmt.Println(uuid)
			save = uuid
		} 
		tree.Insert(&key, "heartbeatDesc #"+strconv.Itoa(i))
	}
	fmt.Println("searching the value of the node : "+save)
	key := StringKey(save)
	res, ok := tree.Get(&key)
	if ok {
		fmt.Println("The result is : "+res)
	} else {
		fmt.Println("No result")
	}
	
	

}