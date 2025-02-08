package tree

import (
	"fmt"
	"strings"
)

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
func (n *Node) Get(path string) (string, error) {
	f := func(r rune) bool {
		return r == '.' || r == '[' || r == ']'
	}

	fields := strings.FieldsFunc(path, f)
	return n.get(fields, 0)
}

func (n *Node) get(fields []string, ind int) (string, error) {
	if len(fields) == ind {
		return "", fmt.Errorf("there is not such field(%s)", strings.Join(fields[:ind], "."))
	}
	if n.Key == fields[ind] {
		if len(fields)-1 == ind {
			return n.valueString(), nil
		}
		js, ok := n.Value.(*Node)
		if !ok {
			return "", fmt.Errorf("there is not such field(%s)", strings.Join(fields[:ind+2], "."))
		}
		return js.get(fields, ind+1)
	} else {
		if n.Next == nil {
			return "", fmt.Errorf("there is not such field(%s)", strings.Join(fields[:ind+1], "."))
		}
		return n.Next.get(fields, ind)
	}
}
