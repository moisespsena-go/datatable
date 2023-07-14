package datatable

import (
	"sort"
)

type Nodes []*Node

func (h Nodes) Sort() {
	var compare = func(x, y int) int {
		if x < y {
			return -1
		}
		if x == y {
			return 0
		}
		return 1
	}
	sort.Slice(h, func(i, j int) bool {
		var (
			a, b = h[i], h[j]
			pri  = compare(a.Row, b.Row)
			sec  = compare(a.Col, b.Col)
		)
		if pri != 0 {
			return pri < 0
		}
		return sec < 0
	})
}

type Node struct {
	val any
	Depth,
	Row,
	Col,
	rowspan,
	colspan int
	Children []*Node
}

var _ Cell = (*Node)(nil)

func (n *Node) Val() any {
	return n.val
}

func (n *Node) Colspan() int {
	return n.colspan
}

func (n *Node) Rowspan() int {
	return n.rowspan
}

func (n *Node) SetVal(v any) {
	n.val = v
}

func (n *Node) SetColspan(v int) {
	n.colspan = v
}

func (n *Node) SetRowspan(v int) {
	n.rowspan = v
}

func (n *Node) AddChild(child ...*Node) *Node {
	n.Children = append(n.Children, child...)
	return n
}

func (n *Node) Add(v any, cb ...func(n *Node)) *Node {
	c := &Node{val: v}
	n.Children = append(n.Children, c)
	for _, cb := range cb {
		cb(c)
	}
	return n
}

func (n *Node) BuildTree() (tree *Tree) {
	var (
		rowsToUse func(t *Node) int
		width     func(t *Node) int
		getCells  func(depth int, t *Node, row, col, rowsLeft int) Nodes

		lcm = func(a, b int) int {
			c := a * b
			for b > 0 {
				t := b
				b = a % b
				a = t
			}
			return c / a
		}
	)

	tree = &Tree{}

	rowsToUse = func(t *Node) int {
		var childrenRows int
		if len(t.Children) > 0 {
			childrenRows++
		}

		for _, child := range t.Children {
			childrenRows = lcm(childrenRows, rowsToUse(child))
		}
		return 1 + childrenRows
	}
	width = func(t *Node) int {
		if len(t.Children) == 0 {
			return 1
		}
		w := 0
		for _, child := range t.Children {
			w += width(child)
		}
		return w
	}
	getCells = func(depth int, t *Node, row, col, rowsLeft int) (cells Nodes) {
		// Add top-most cell corresponding to the root of the current tree.
		rootRows := rowsLeft / rowsToUse(t)
		cells = append(cells, &Node{t.val, depth, row, col, rootRows, width(t), nil})
		for _, child := range t.Children {
			cells = append(cells, getCells(depth+1, child, row+rootRows, col, rowsLeft-rootRows)...)
			col += width(child)
		}
		if (row + 1) > tree.NumRows {
			tree.NumRows = row + 1
		}
		return
	}

	cells := getCells(0, n, 0, 0, rowsToUse(n))
	cells.Sort()
	tree.Rows = make(Rows, tree.NumRows)

	for _, cell := range cells {
		tree.Rows[cell.Row] = append(tree.Rows[cell.Row], cell)
	}
	tree.Rows = tree.Rows[1:]
	tree.NumRows--

	return
}

type Tree struct {
	Rows    Rows
	NumRows int
}
