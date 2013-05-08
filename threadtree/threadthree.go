package main

import (
	"fmt"
	"math"
	"math/rand"
)

type BinTree struct {
	root *Node
}

type pointerDescription int

const (
	noLinks       pointerDescription = 0
	leftIsLinked                     = 1
	rightIsLinked                    = 2
	bothAreLinked                    = leftIsLinked | rightIsLinked
)

type Node struct {
	value int
	links pointerDescription
	left  *Node
	right *Node
}

func New() *BinTree {
	return &BinTree{}
}

func (b *BinTree) add(val int) {

	if b.root == nil {
		b.root = &Node{val, bothAreLinked, nil, nil}
		return
	}

	b.root.add(val)
}

func (n *Node) add(val int) {

	//	println("adding ", val)

	for {

		if val <= n.value {
			if n.links&leftIsLinked == leftIsLinked {
				n.links ^= leftIsLinked
				n.left = &Node{val, leftIsLinked | rightIsLinked, n.left, n}
				return
			}

			n = n.left
		} else {
			if n.links&rightIsLinked == rightIsLinked {
				n.links ^= rightIsLinked
				n.right = &Node{val, leftIsLinked | rightIsLinked, n, n.right}
				return
			}

			n = n.right
		}
	}
}

func (b *BinTree) del(val int) {

	if b.root == nil {
		return
	}

	if b.root.isLeafNode() && b.root.value == val {
		b.root = nil
	}

	b.root = b.root.del(nil, val)
}

func (n *Node) isLeafNode() bool {
	return n.links == bothAreLinked
}

func (n *Node) del(parent *Node, val int) *Node {

	if n.value == val {

		p := n.pred()
		s := n.succ()

		if p != nil && p.links&rightIsLinked == rightIsLinked {
			p.right = s
		}

		if s != nil && s.links&leftIsLinked == leftIsLinked {
			s.left = p
		}

		// 1) leaf node
		if n.links&bothAreLinked == bothAreLinked {

			//			fmt.Println("both are linked: ", val)

			// we're about to free up a node for our parent
			// -> turn it into a linked-list pointer
			// these cases could probably be combined with the above
			var child *Node

			switch {
			case n == parent.left:
				child = n.left
			case n == parent.right:
				child = n.right
			}

			switch {
			case n == parent.left:
				parent.left = child
				parent.links |= leftIsLinked
			case n == parent.right:
				parent.right = child
				parent.links |= rightIsLinked
			}

			return child
		}

		// 2) left linked list, right is child(-tree)
		// 3) right is linked list, left is child(-tree)
		if n.links&leftIsLinked == leftIsLinked || n.links&rightIsLinked == rightIsLinked {

			//			fmt.Println("one is linked: ", val)
			var child *Node
			if n.links&leftIsLinked == leftIsLinked {
				child = n.right
			} else {
				child = n.left
			}

			// update the correct pointer from our parent
			if parent != nil {
				switch {
				case n == parent.left:
					parent.left = child
				case n == parent.right:
					parent.right = child
				}
			}

			return child
		}

		// 4) both left and right are children

		//		fmt.Println("none are linked: ", val)

		// add the right subtree to the 'end' of the left
		// the largest value on the left is n.pred(), so we already have it
		// we just need to set the child pointers
		if p != nil {
			p.right = n.right
			if p.links&rightIsLinked == 0 {
				panic("bad mojo")

			}
			p.links &= ^rightIsLinked
		} else {
			// I have no predecessor
			fmt.Println("no pred found for val=", val)
		}

		// update the correct pointer from our parent
		if parent != nil {
			switch {
			case n == parent.left:
				parent.left = n.left
			case n == parent.right:
				parent.right = n.left
			}
		}
		return n.left
	}

	if val < n.value && n.links&leftIsLinked == 0 {
		n.left.del(n, val)
	} else if val > n.value && n.links&rightIsLinked == 0 {
		n.right.del(n, val)
	}

	return n
}

func (b *BinTree) walkIn(f func(*Node)) {

	if b.root == nil {
		return
	}

	b.root.walkIn(f)
}

func (n *Node) walkIn(f func(*Node)) {
	if n.links&leftIsLinked == 0 {
		n.left.walkIn(f)
	}

	f(n)

	if n.links&rightIsLinked == 0 {
		n.right.walkIn(f)
	}
}

func (b *BinTree) walkLink(f func(*Node)) {

	if b.root == nil {
		return
	}

	b.root.walkLink(f)
}

func (b *BinTree) walkLinkReverse(f func(*Node)) {

	if b.root == nil {
		return
	}

	b.root.walkLinkReverse(f)
}

func (n *Node) pred() *Node {

	if n.links&leftIsLinked == leftIsLinked {
		return n.left
	}

	node := n.left
	for node.links&rightIsLinked == 0 {
		node = node.right
	}
	return node

}

func (n *Node) succ() *Node {

	if n.links&rightIsLinked == rightIsLinked {
		return n.right
	}

	node := n.right
	for node.links&leftIsLinked == 0 {
		node = node.left
	}
	return node

}

func (n *Node) walkLink(f func(*Node)) {

	var node *Node = nil

	// find the 'start' of the linked-list
	node = n
	for node.links&leftIsLinked == 0 {
		node = node.left
	}

	for node != nil {
		f(node)
		node = node.succ()
	}
}

func (n *Node) walkLinkReverse(f func(*Node)) {

	var node *Node = nil

	// find the 'end'' of the linked-list
	node = n
	for node.links&rightIsLinked == 0 {
		node = node.right
	}

	for node != nil {
		f(node)
		node = node.pred()
	}
}

// fisher-yates
func shuffle(array []int) {

	for i := len(array) - 1; i >= 1; i-- {
		j := rand.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
}

func check_tree(b *BinTree) {

	var i int

	i = 0

	b.walkIn(func(n *Node) {
		//		fmt.Printf("%p = %d, left = %x right = %x\n", n, n.value, n.left, n.right)
		if n.value < i {
			fmt.Println("in-order tree walk failed: i=", i, " value=", n.value)
		}
		i = n.value
	})
	//fmt.Println()

	i = 0

	b.walkLink(func(n *Node) {
		//		fmt.Printf("%p = %d, list=%d left = %x right = %x\n", n, n.value, n.links, n.left, n.right)

		if n.value < i {
			fmt.Println("reverse linked-list tree walk failed: i=", i, " value=", n.value)
		}
		i = n.value

	})

	//fmt.Println()
	i = math.MaxInt32
	b.walkLinkReverse(func(n *Node) {
		//		fmt.Printf("%p = %d, list=%d left = %x right = %x\n", n, n.value, n.links, n.left, n.right)

		if i < n.value {
			fmt.Println("reverse linked-list tree walk failed: i=", i, " value=", n.value)
		}
		i = n.value
	})

}

func main() {

	b := New()

	t := []int{2, 1, 7, 5, 4, 9, 3, 6, 8}

	for _, v := range t {
		b.add(v)
	}

	b.del(3)
	b.del(5)
	b.del(4)

	b.walkIn(func(n *Node) { fmt.Printf("%p = %d, left = %x right = %x\n", n, n.value, n.left, n.right) })
	fmt.Println()
	b.walkLink(func(n *Node) {
		fmt.Printf("%p = %d, list=%d left = %x right = %x\n", n, n.value, n.links, n.left, n.right)
	})

	fmt.Println()
	b.walkLinkReverse(func(n *Node) {
		fmt.Printf("%p = %d, list=%d left = %x right = %x\n", n, n.value, n.links, n.left, n.right)
	})

	t = []int{9, 8, 6, 1, 7}

	for _, v := range t {

		b.walkIn(func(n *Node) { fmt.Printf("%p = %d, left = %x right = %x\n", n, n.value, n.left, n.right) })

		b.del(v)

		fmt.Println()

	}

	b.walkIn(func(n *Node) { fmt.Printf("%p = %d, left = %x right = %x\n", n, n.value, n.left, n.right) })
	b.walkLink(func(n *Node) {
		fmt.Printf("%p = %d, list=%d left = %x right = %x\n", n, n.value, n.links, n.left, n.right)
	})

	b.walkLinkReverse(func(n *Node) {
		fmt.Printf("%p = %d, list=%d left = %x right = %x\n", n, n.value, n.links, n.left, n.right)
	})

	t = make([]int, 10000)

	for i := 0; i < 10000; i++ {
		t[i] = i
	}

	shuffle(t)

	for _, v := range t {
		b.add(v)
		check_tree(b)
	}

	check_tree(b)

	shuffle(t)

	for _, v := range t {
		b.del(v)
		check_tree(b)
	}

	fmt.Println("done adding")

	b.walkLinkReverse(func(n *Node) {
		fmt.Printf("%p = %d, list=%d left = %x right = %x\n", n, n.value, n.links, n.left, n.right)
	})

}
