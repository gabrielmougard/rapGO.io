package bst
import (
	"bytes"
	"fmt"
	"rapGO.io/src/algorithmservice/trees"
	"rapGO.io/src/algorithmservice/utils"
	llstack "rapGO.io/src/algorithmservice/stacks/linkedliststack"
	"strings"
)

func assertTreeImplementation() {
	var _ trees.Tree = (*Tree)(nil)
}
//Tree hold the elements a BST (Binary Search Tree)
type Tree struct {
	Root 		*Node 				//Root node
	Comparator 	utils.Comparator 	//Key comparator
	size 		int					//Total number of keys in the Tree
}

//Node is the single element inside the Tree
type Node struct {
	Parent 	*Node 		
	Value   interface{}
	Left 	*Node		//Left child
	Right 	*Node		//Right child
}

func New(comparator utils.Comparator) *Tree {
	return &Tree{Comparator: comparator}
}

func NewWithIntComparator() *Tree {
	return &Tree{Comparator: utils.IntComparator}
}

func NewWithStringComparator() *Tree {
	return &Tree{Comparator: utils.StringComparator}
}

/**
 * 'Insert' inserts a value into the tree.
 * value should adhere to the comparator's type assertion, otherwise method panics.
 */
func (tree *Tree) Insert(newNode *Node) {
	y := &Node{}
	x := tree.Root
	for ; x != nil ; {
		y = x
		if tree.Comparator(newNode.Value, x.Value) > -1 { //newNode.Value < x.Value
			x = x.Left
		} else {
			x = x.Right
		}
	}
	newNode.Parent = y
	if y == nil {
		tree.Root = newNode //the tree was empty
	} else if tree.Comparator(newNode.Value, y.Value) > -1 {
		y.Left = newNode
	} else {
		y.Right = newNode
	}
	tree.size++
}

func (tree *Tree) Size() int {
	return tree.size
}
/**
 * Inorder traversal with no recursion. Run in O(n)
 */ 
func (tree *Tree) InorderNoRecursive() {
	stack := &llstack.Stack{}
	currentNode := tree.Root
	done := false
	
	for ; !done ; {
		if (currentNode != nil) {
			stack.Push(currentNode)
			currentNode = currentNode.Left
		} else {
			if (!stack.Empty()) {
				currentNode = stack.Pop().(*Node)
				fmt.Printf(utils.ToString(currentNode.Value))
				currentNode = currentNode.Right
			} else {
				done = true
			}
		}
	}
}

func (tree *Tree) Delete(nodeToDelete *Node) {
	if nodeToDelete.Left == nil {
		tree.transplant(nodeToDelete, nodeToDelete.Right)
	} else if nodeToDelete.Right == nil {
		tree.transplant(nodeToDelete, nodeToDelete.Left)
	} else {
		y := tree.Minimum(nodeToDelete.Right)
		if y.Parent != nodeToDelete {
			tree.transplant(y, y.Right)
			y.Right = nodeToDelete.Right
			y.Right.Parent = y
		}
		tree.transplant(nodeToDelete, y)
		y.Left = nodeToDelete.Left
		y.Left.Parent = y
	}
}

func (tree *Tree) transplant(u,v *Node) {
	if u.Parent == nil {
		tree.Root = v
	} else if u == u.Parent.Left {
		u.Parent.Left = v
	} else {
		u.Parent.Right = v
	}
	if v != nil {
		v.Parent = u.Parent
	}
}

func (tree *Tree) Minimum(startingNode *Node) *Node {
	y := startingNode
	for ; y.Left != nil ; {
		y = y.Left
	}
	return y
}

func (tree *Tree) Maximum(startingNode *Node) *Node {
	y := startingNode
	for ; y.Right != nil ; {
		y = y.Right
	}
	return y
}

func (tree *Tree) Successor(startingNode *Node) *Node {
	if startingNode.Right != nil {
		return tree.Minimum(startingNode.Right)
	}
	y := startingNode.Parent
	for ; y != nil && startingNode == y.Right ; {
		startingNode = y
		y = y.Parent
	}
	return y
}

func (tree *Tree) Predecessor(startingNode *Node) *Node {
	if startingNode.Left != nil {
		return tree.Maximum(startingNode.Left)
	}
	y := startingNode.Parent
	for ; y != nil && startingNode == y.Left ; {
		startingNode = y
		y = y.Parent
	}
	return y
}