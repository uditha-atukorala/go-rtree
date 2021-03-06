
package rtree

import (
	"math"
	"github.com/uditha-atukorala/go-rtree/geom"
)

type Rtree struct {
	root *node
}


func NewRtree( maxNodeEntries uint16 ) ( *Rtree ) {
	tree := &Rtree{
		root: NewNode( maxNodeEntries ),
	}
	return tree
}

func ( tree *Rtree ) Insert( item Item )  {
	leaf := tree.chooseLeaf( tree.root, item )
	leaf.insert( item )
}

func ( tree *Rtree ) Search( p *geom.Point ) ( []Item, int )  {
	n, nodes := tree.root, []*node{}
	results := []Item{}
	cost := 0

	for n != nil && n.Mbr().ContainsPoint( p ) {
		cost++
		if n.isLeaf() {
			for _, item := range n.items {
				if item.Mbr().ContainsPoint( p ) {
					results = append( results, item )
				}
			}
		} else {
			for _, child := range n.children {
				if child.Mbr().ContainsPoint( p ) {
					nodes = append( nodes, child )
				}
			}
		}

		if n = nil; len( nodes ) > 0 {
			n, nodes = nodes[0], nodes[1:]
		}
	}

	return results, cost
}

func ( tree *Rtree ) chooseLeaf( n *node, item Item ) ( *node ) {
	r := item.Mbr()

	for !n.isLeaf() {
		var chosen *node
		var minCost, minArea float64 = math.Inf( 0 ), math.Inf( 0 )

		for _, child := range n.children {
			cost, area := child.areaCost( r )
			if cost < minCost {
				minCost = cost
				chosen  = child
				if area < minArea {
					minArea = area
				}
			} else if cost == minCost && area < minArea {
				minArea = area
				chosen  = child
			}
		}

		n = chosen
	}

	return n
}

