package states

import (
	"testing"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	
)

func TestTree(t *testing.T) {
	fmt.Println("starting the test")
	tree := NewRbTree()
	fmt.Println("tree created")
	var save string
	for i := 0; i < 1000; i++ {
		fmt.Printf("iteration # %d\n", i)
		uuid := uuid.New().String()
		key := StringKey(uuid)
		if i == 100 {
			fmt.Println(uuid)
			save = uuid
		} 
		tree.Mu.Lock()
		tree.Insert(&key, "heartbeatDesc #"+strconv.Itoa(i))
		tree.Mu.Unlock()
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