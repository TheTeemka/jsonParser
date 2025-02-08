package tree

import (
	"fmt"
	"strings"
)

type nodeWriter struct {
	strings.Builder
}

func newNodeWriter() *nodeWriter {
	return &nodeWriter{}
}

func (w *nodeWriter) write(ss ...any) {
	format := ss[0].(string)
	if len(ss) == 1 {
		w.WriteString(format)
	} else {
		w.WriteString(fmt.Sprintf(format, ss[1:]...))
	}
}

func (w *nodeWriter) writeField(n *Node, tab string) {
	w.write("%s%q: ", tab, n.Key)
	w.writeValue(n, tab)
	if n.Next != nil {
		w.write(",\n")
		w.writeField(n.Next, tab)
	} else {
		w.write("\n")
	}
}

func (w *nodeWriter) writeValue(n *Node, tab string) {
	switch n.ValueType {
	case JSON:
		w.write("{\n")
		js := n.Value.(*Node)
		w.writeField(js, tab+defaultTAB)
		w.write("%s}", tab)
	case String:
		w.write("%q", n.Value)
	case Number, Bool, Null:
		w.write(n.Value)
	case Array:
		w.write("[\n%s", tab+defaultTAB)
		arr := n.Value.(*Node)
		w.writeArray(arr, tab+defaultTAB)
		w.write("\n%s]", tab)
	}
}

func (w *nodeWriter) writeArray(n *Node, tab string) {
	w.writeValue(n, tab)
	if n.Next != nil {
		w.write(",\n%s", tab)
		w.writeArray(n.Next, tab)
	}
}
