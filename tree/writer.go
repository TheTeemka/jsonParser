package tree

import (
	"fmt"
	"strings"
)

type nodeWriter struct {
	strings.Builder
	tabSize string
}

func newNodeWriter(tab string) *nodeWriter {
	return &nodeWriter{tabSize: tab}
}

func (w *nodeWriter) write(ss ...any) {
	format := ss[0].(string)
	if len(ss) == 1 {
		w.WriteString(format)
	} else {
		w.WriteString(fmt.Sprintf(format, ss[1:]...))
	}
}

func (w *nodeWriter) writeField(n *Node) {
	w.write("%q: ", n.Key)
	w.writeValue(n)
	if n.Next != nil {
		w.write(",")
		w.writeField(n.Next)
	}
}

func (w *nodeWriter) writeValue(n *Node) {
	switch n.ValueType {
	case JSON:
		w.write("{")
		js := n.Value.(*Node)
		w.writeField(js)
		w.write("}")
	case String:
		w.write("%q", n.Value)
	case Number, Bool, Null:
		w.write(n.Value)
	case Array:
		w.write("[")
		arr := n.Value.(*Node)
		w.writeArray(arr)
		w.write("]")
	}
}

func (w *nodeWriter) writeArray(n *Node) {
	w.writeValue(n)
	if n.Next != nil {
		w.write(",")
		w.writeArray(n.Next)
	}
}
