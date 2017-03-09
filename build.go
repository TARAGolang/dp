package dp

import (
	"math"
	"simplex/struct/bst"
	"simplex/struct/item"
	"simplex/struct/stack"
)

//Douglas Peucker BST
//  uses iteration to build tree, state managed with a array stack
//  Note: int fn_int must return Int type sorted by most interesting
//  process[Function]   - optional process node callback
//  should return node after process
//  signature : process(node) node{}
func (self *DP) Build(process ...func(item item.Item)) *DP {
	procFn := func(item item.Item) {}
	if len(process) > 0 && process[0] != nil {
		procFn = process[0]
	}
	var offset = self.Offset

	//use local if global is nil
	if offset == nil {
		offset = self.offset
	}

	var index int
	var val float64

	var range_ *item.Int2D
	var n, l, r *bst.Node
	var stk = stack.NewStack()

	var node *Node
	var root = bst.NewNode(
		NewNode(0, len(self.Pln)-1),
	)
	self.BST.Root = root

	//pre stk
	stk.Add(root)

	for !stk.IsEmpty() {
		n = self.AsBSTNode_Any(stk.Pop())
		procFn(n.Key)

		node = self.AsDPNode(n)
		range_ = node.Key

		node.Ints = offset.Offsets(node)

		vobj := node.Ints.Peek().(*Vertex)

		index = int(vobj.index)
		val = vobj.value

		if !math.IsNaN(val) && val <= self.Res {
			self.Simple.Add(range_[:]...)
		} else {
			//left and right branch
			l = bst.NewNode(
				NewNode(range_[0], index),
			)
			r = bst.NewNode(
				NewNode(index, range_[1]),
			)

			//update pointers
			bst.Ptr(n, l, bst.NewBranch().AsLeft())
			bst.Ptr(n, r, bst.NewBranch().AsRight())

			//node stk
			stk.Add(r)
			stk.Add(l)
		}
	}

	return self
}
