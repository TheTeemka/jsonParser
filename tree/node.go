package tree

import "strings"

type Type string

const (
	JSON   Type = "JSON"
	Number Type = "NUMBER"
	String Type = "STRING"
	Array  Type = "ARRAY"
	Bool   Type = "BOOL"
	Null   Type = "NULL"
)

const defaultTAB = "  "

type Node struct {
	Key       string
	Value     any
	ValueType Type
	Next      *Node
}

func (n *Node) String() string {
	w := newNodeWriter()

	w.write("{\n")
	w.writeField(n, defaultTAB)
	w.write("}")

	return w.String()
}

func (n *Node) valueString() string {
	w := newNodeWriter()

	w.writeValue(n, "")

	return w.String()
}
func (n *Node) Get(path string) string {
	f := func(r rune) bool {
		return r == '.' || r == '[' || r == ']'
	}

	fields := strings.FieldsFunc(path, f)
	return n.get(fields)
}

func (n *Node) get(fields []string) string {
	if n.Key == fields[0] {
		if len(fields) == 1 {
			return n.valueString()
		}
		js := n.Value.(*Node)
		return js.get(fields[1:])
	} else {
		return n.Next.get(fields)
	}
}
