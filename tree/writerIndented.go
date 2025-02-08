package tree

func (w *nodeWriter) writeFieldIndented(n *Node, tab string) {
	w.write("%s%q: ", tab, n.Key)
	w.writeValueIndented(n, tab)
	if n.Next != nil {
		w.write(",\n")
		w.writeFieldIndented(n.Next, tab)
	} else {
		w.write("\n")
	}
}

func (w *nodeWriter) writeValueIndented(n *Node, tab string) {
	switch n.ValueType {
	case JSON:
		w.write("{\n")
		js := n.Value.(*Node)
		w.writeFieldIndented(js, tab+w.tabSize)
		w.write("%s}", tab)
	case String:
		w.write("%q", n.Value)
	case Number, Bool, Null:
		w.write(n.Value)
	case Array:
		w.write("[\n%s", tab+w.tabSize)
		arr := n.Value.(*Node)
		w.writeArrayIndented(arr, tab+w.tabSize)
		w.write("\n%s]", tab)
	}
}

func (w *nodeWriter) writeArrayIndented(n *Node, tab string) {
	w.writeValueIndented(n, tab)
	if n.Next != nil {
		w.write(",\n%s", tab)
		w.writeArrayIndented(n.Next, tab)
	}
}
