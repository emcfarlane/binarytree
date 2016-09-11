package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Node struct {
	left, right *Node
	value       int
}

func (root *Node) growTree(treeDepth int) {
	if treeDepth > 0 {
		if root.left == nil {
			root.left = &Node{}
			root.right = &Node{}
		}
		root.left.growTree(treeDepth - 1)
		root.right.growTree(treeDepth - 1)
	} else {
		root.left = nil
		root.right = nil
	}
}

func (root *Node) populateTree(nodeValue int) {
	root.value = nodeValue
	if root.left != nil {
		root.left.populateTree(2*nodeValue - 1)
		root.right.populateTree(2 * nodeValue)
	}
}

func (root *Node) computeTreeChecksum() int {
	if root.left != nil {
		return root.left.computeTreeChecksum() - root.right.computeTreeChecksum() + root.value
	}
	return root.value
}

func main() {
	flag.Parse()
	n, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		os.Exit(1)
	}

	minDepth := 4
	maxDepth := n
	if minDepth+2 > n {
		maxDepth = minDepth + 2
	}

	stretchTree := &Node{}
	stretchTree.growTree(maxDepth + 1)
	stretchTree.populateTree(0)
	fmt.Printf("stretch tree of depth %d\t check: %d\n", maxDepth+1, stretchTree.computeTreeChecksum())

	longLivedTree := stretchTree
	longLivedTree.growTree(maxDepth)
	stretchTree.populateTree(0)

	var wg = &sync.WaitGroup{}
	buf := make([]string, maxDepth+1)
	for d := minDepth; d <= maxDepth; d += 2 {
		wg.Add(1)
		go func(depth, iterations int) {
			defer wg.Done()

			treeRoot := &Node{}
			treeRoot.growTree(depth)

			var check int
			for i := 0; i < iterations; i++ {
				treeRoot.populateTree(i)
				check += treeRoot.computeTreeChecksum()
				treeRoot.populateTree(-i)
				check += treeRoot.computeTreeChecksum()
			}
			buf[depth] = fmt.Sprintf("%d\t trees of depth %d\t check: %d\n", iterations*2, depth, check)
		}(d, 1<<uint(maxDepth-d+minDepth))
	}
	wg.Wait()

	for depth := minDepth; depth <= maxDepth; depth += 2 {
		fmt.Print(buf[depth])
	}
	fmt.Printf("long lived tree of depth %d\t check: %d\n", maxDepth, longLivedTree.computeTreeChecksum())
}
