package states

import (
    "sync"
    "fmt"
    "rapGO.io/src/heartbeatservice/pkg/setting"
)
// KeyComparison structure used as result of comparing two keys 
type KeyComparison int8

const (
    // KeyIsLess is returned as result of key comparison if the first key is less than the second key
    KeyIsLess KeyComparison = iota - 1 
    // KeysAreEqual is returned as result of key comparison if the first key is equal to the second key
    KeysAreEqual
    // KeyIsGreater is returned as result of key comparison if the first key is greater than the second key
    KeyIsGreater
)

const (
    red byte = byte(0)
    black byte = byte(1)
    zeroOrEqual = int8(0)
)

func (tree KeyComparison) String() string {
    switch tree {
    case KeyIsLess:
        return "lessThan"
    case KeyIsGreater:
        return "greaterThan"
    default:
        return "equalTo"
    }
}

// RbKey interface
type RbKey interface {
    ComparedTo(key RbKey) KeyComparison
}

// rbNode structure used for storing key and value pairs
type rbNode struct {
    key RbKey
    heartbeatDesc string
    color byte
	left, right *rbNode

	DescChan chan string
}

// RbTree structure
type RbTree struct {
    root *rbNode
    count int
    version uint32

	Mu sync.Mutex
}

// NewRbTree creates a new RbTree and returns its address
func NewRbTree() *RbTree {
    return &RbTree{}
}

// newRbNode creates a new rbNode and returns its address
func newRbNode(key RbKey, heartbeatDesc string) *rbNode {
    result := &rbNode{
        key: key,
        heartbeatDesc: heartbeatDesc,
		color: red,
		DescChan: make(chan string, setting.TotalHeartbeatNumber()), //how many heartbeats do we have for a rap generation ? ==> I counted that 6 is our max.
    }

    result.DescChan <- heartbeatDesc //write to channel
    return result
}

// isRed checks if node exists and its color is red
func isRed(node *rbNode) bool {
    return node != nil && node.color == red
}

// isBlack checks if node exists and its color is black
func isBlack(node *rbNode) bool { 
    return node != nil && node.color == black 
}

// min finds the smallest node key including the given node
func min(node *rbNode) *rbNode {
    if node != nil {
        for node.left != nil {
            node = node.left
        }
    }
    return node
}

// max finds the greatest node key including the given node
func max(node *rbNode) *rbNode {
    if node != nil {
        for node.right != nil {
            node = node.right
        }
    }
    return node
}

// floor returns the largest key node in the subtree rooted at x less than or equal to the given key
func floor(node *rbNode, key RbKey) *rbNode {
    if node == nil {
        return nil
    }
    
    switch key.ComparedTo(node.key) {
    case KeysAreEqual:
        return node
    case KeyIsLess:
        return floor(node.left, key)
    default:
        fn := floor(node.right, key)
        if fn != nil {
            return fn
        }
        return node
    }
}

// ceilig returns the smallest key node in the subtree rooted at x greater than or equal to the given key
func ceiling(node *rbNode, key RbKey) *rbNode {  
    if node == nil {
        return nil
    }
    
    switch key.ComparedTo(node.key) {
    case KeysAreEqual:
        return node
    case KeyIsGreater:
        return ceiling(node.right, key)
    default:
        cn := ceiling(node.left, key)
        if cn != nil {
            return cn
        }
        return node
    }
}

// flipColor switchs the color of the node from red to black or black to red
func flipColor(node *rbNode) {
    if node.color == black {
        node.color = red
    } else {
        node.color = black
    }
}

// colorFlip switchs the color of the node and its children from red to black or black to red
func colorFlip(node *rbNode) {
    flipColor(node)
    flipColor(node.left)
    flipColor(node.right)
}

// rotateLeft makes a right-leaning link lean to the left
func rotateLeft(node *rbNode) *rbNode {
    child := node.right
    node.right = child.left
    child.left = node
    child.color = node.color
    node.color = red

    return child
}

// rotateRight makes a left-leaning link lean to the right
func rotateRight(node *rbNode) *rbNode {
    child := node.left
    node.left = child.right
    child.right = node
    child.color = node.color
    node.color = red

    return child
}

// moveRedLeft makes node.left or one of its children red,
// assuming that node is red and both children are black.
func moveRedLeft(node *rbNode) *rbNode {
    colorFlip(node)
    if isRed(node.right.left) {
        node.right = rotateRight(node.right)
        node = rotateLeft(node)
        colorFlip(node)
    }
    return node
}

// moveRedRight makes node.right or one of its children red,
// assuming that node is red and both children are black.
func moveRedRight(node *rbNode) *rbNode {
    colorFlip(node)
    if isRed(node.left.left) {
        node = rotateRight(node)
        colorFlip(node)
    }
    return node
}

// balance restores red-black tree invariant
func balance(node *rbNode) *rbNode {
    if isRed(node.right) {
        node = rotateLeft(node)
    }
    if isRed(node.left) && isRed(node.left.left) {
        node = rotateRight(node)
    }
    if isRed(node.left) && isRed(node.right) {
        colorFlip(node)
    }
    return node
}

// deleteMin removes the smallest key and associated value from the tree
func deleteMin(node *rbNode) *rbNode {
    if node.left == nil {
        return nil
    }    
    if isBlack(node.left) && !isRed(node.left.left) {
        node = moveRedLeft(node)
    }
    node.left = deleteMin(node.left)
    /* if node.left != nil {
        node.left.parent = node
    } */
    return balance(node)
}

func (node *rbNode) GetDescChan() chan string {
	if node.DescChan != nil {
		return node.DescChan
	} else {
		return nil
	}
}

// Count returns if count of the nodes stored.
func (tree *RbTree) Count() int {
    return tree.count
}

// IsEmpty returns if the tree has any node.
func (tree *RbTree) IsEmpty() bool {
    return tree.root == nil
}

// Min returns the smallest key in the tree.
func (tree *RbTree) Min() (RbKey, string) {
    if tree.root != nil {
        result := min(tree.root)
        return result.key, result.heartbeatDesc
    }
    return nil, ""
} 

// Max returns the largest key in the tree.
func (tree *RbTree) Max() (RbKey, string) {
    if tree.root != nil {
        result := max(tree.root)
        return result.key, result.heartbeatDesc
    }
    return nil, ""
} 

// Floor returns the largest key in the tree less than or equal to key
func (tree *RbTree) Floor(key RbKey) (RbKey, string) {
    if key != nil && tree.root != nil {
        node := floor(tree.root, key)
        if node == nil {
            return nil, ""
        }
        return node.key, node.heartbeatDesc
    }
    return nil, ""
}    

// Ceiling returns the smallest key in the tree greater than or equal to key
func (tree *RbTree) Ceiling(key RbKey) (RbKey, string) {
    if key != nil && tree.root != nil {
        node := ceiling(tree.root, key)
        if node == nil {
            return nil, ""
        }
        return node.key, node.heartbeatDesc
    }
    return nil, ""
}

// Get returns the stored value if key found and 'true', 
// otherwise returns 'false' with second return param if key not found 
func (tree *RbTree) Get(key RbKey) (string, bool) {
    if key != nil && tree.root != nil {
        node := tree.find(key)
        if node != nil {
            return node.heartbeatDesc, true
        }
    }
    return "", false
}

//GetNode returns the node with the associated key.
// if the node is found it returns a pointer to this node and 'true'
// else it returns nil and 'false'
func (tree *RbTree) GetNode(key RbKey) (*rbNode, bool) {
	if key != nil && tree.root != nil {
		node := tree.find(key)
		if node != nil {
			return node, true
		} else {
			return nil, false
		}
	}
	return nil, false
}

//Change the value of a node
func (tree *RbTree) EditDesc(key RbKey, newHeartbeatDesc string) {
    fmt.Println("EDITDESC : "+newHeartbeatDesc)
    if key != nil && tree.root != nil {
		node := tree.find(key)
		if node != nil {
			node.DescChan <- newHeartbeatDesc
			node.heartbeatDesc = newHeartbeatDesc
		}
	}
}

// find returns the node if key found, otherwise returns nil 
func (tree *RbTree) find(key RbKey) *rbNode {
    for node := tree.root; node != nil; { 
        switch key.ComparedTo(node.key) {
        case KeyIsLess:
            node = node.left
        case KeyIsGreater:
            node = node.right
        default:
            return node
        }    
    }
    return nil
}

// Exists returns the node if key found, otherwise returns nil 
func (tree *RbTree) Exists(key RbKey) bool {
    return tree.find(key) != nil
}

// Insert inserts the given key and value into the tree
func (tree *RbTree) Insert(key RbKey, heartbeatDesc string) {
    fmt.Println("INSERTNODE : "+heartbeatDesc)
    if key != nil {
        tree.version++
        tree.root = tree.insertNode(tree.root, key, heartbeatDesc);
        tree.root.color = black
        // tree.root.parent = nil
    }
}

// insertNode adds the given key and value into the node
func (tree *RbTree) insertNode(node *rbNode, key RbKey, heartbeatDesc string) *rbNode {
    if node == nil {
        tree.count++
        return newRbNode(key, heartbeatDesc)
    }

    switch key.ComparedTo(node.key) {
    case KeyIsLess:
        node.left  = tree.insertNode(node.left,  key, heartbeatDesc)
        // node.left.parent = node
    case KeyIsGreater:
        node.right = tree.insertNode(node.right, key, heartbeatDesc)
        // node.right.parent = node

    }
    return balance(node)
}

// Delete deletes the given key from the tree
func (tree *RbTree) Delete(key RbKey) {
    tree.version++
    tree.root = tree.deleteNode(tree.root, key)
    if tree.root != nil {
        tree.root.color = black
    }
}

// deleteNode deletes the given key from the node
func (tree *RbTree) deleteNode(node *rbNode, key RbKey) *rbNode {
    if node == nil {
        return nil
    }
    
    cmp := key.ComparedTo(node.key)
    if cmp == KeyIsLess {
        if isBlack(node.left) && !isRed(node.left.left) {
            node = moveRedLeft(node)
        }
        node.left = tree.deleteNode(node.left, key)
    } else {
        if cmp == KeysAreEqual {
            heartbeatDesc := node.heartbeatDesc
            if heartbeatDesc != "" {
                node.heartbeatDesc = heartbeatDesc
                return node
            }
        }
        
        if isRed(node.left) {
            node = rotateRight(node)
        }
        
        if isBlack(node.right) && !isRed(node.right.left) {
            node = moveRedRight(node)
        }
        
        if key.ComparedTo(node.key) != KeysAreEqual {
            node.right = tree.deleteNode(node.right, key)
        } else {
            if node.right == nil {
                return nil
            }

            rm := min(node.right)
            node.key   = rm.key
            node.heartbeatDesc = rm.heartbeatDesc
            node.right = deleteMin(node.right)

            rm.left = nil
            rm.right = nil
            
            tree.count--
        }
    }
    return balance(node)
}