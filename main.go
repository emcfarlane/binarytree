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

func (root *Node) createTree(rootNodeValue, treeDepth int) {
	root.value = rootNodeValue

	if treeDepth > 0 {
		if root.left == nil {
			root.left = &Node{}
			root.right = &Node{}
		}
		root.left.createTree(2*rootNodeValue-1, treeDepth-1)
		root.right.createTree(2*rootNodeValue, treeDepth-1)

	} else if root.left != nil {
		root.left = nil
		root.right = nil
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

	var wg = &sync.WaitGroup{}
	buf := make([]string, maxDepth+2)

	wg.Add(1)
	go func() {
		defer wg.Done()

		stretchTree := &Node{}
		stretchTree.createTree(0, maxDepth+1)
		buf[maxDepth+1] = fmt.Sprintf("stretch tree of depth %d\t check: %d\n", maxDepth+1, stretchTree.computeTreeChecksum())
	}()

	longLivedTree := &Node{}
	longLivedTree.createTree(0, maxDepth)

	for d := minDepth; d <= maxDepth; d += 2 {
		wg.Add(1)
		go func(depth, iterations int) {
			defer wg.Done()

			treeRoot := &Node{}
			var check int
			for i := 0; i < iterations; i++ {
				treeRoot.createTree(i, depth)
				check += treeRoot.computeTreeChecksum()
				treeRoot.createTree(-i, depth)
				check += treeRoot.computeTreeChecksum()
			}
			buf[depth] = fmt.Sprintf("%d\t trees of depth %d\t check: %d\n", iterations*2, depth, check)
		}(d, 1<<uint(maxDepth-d+minDepth))
	}
	wg.Wait()

	fmt.Print(buf[maxDepth+1])
	for depth := minDepth; depth <= maxDepth; depth += 2 {
		fmt.Print(buf[depth])
	}
	fmt.Printf("long lived tree of depth %d\t check: %d\n", maxDepth, longLivedTree.computeTreeChecksum())
}
