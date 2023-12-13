package bson_test

import (
	"encoding/hex"
	"fmt"

	"github.com/cristalhq/bson"
)

func Example_types() {
	arr := bson.A{"a", "b", "c", 12345}
	doc := bson.D{{"hello", "world"}, {"answer", 42}}
	doo := bson.M{"hello": "world", "pi": 3.14159}

	fmt.Println(arr)
	fmt.Println(arr.AsD())
	fmt.Println()

	fmt.Println(doc)
	fmt.Println(doc.AsM())
	fmt.Println()

	fmt.Println(doo)
	fmt.Println(doo.AsD())

	_, _, _ = arr, doc, doo

	// Output:
	// [a b c 12345]
	// [{0 a} {1 b} {2 c} {3 12345}]
	//
	// [{hello world} {answer 42}]
	// map[answer:42 hello:world]
	//
	// map[hello:world pi:3.14159]
	// [{hello world} {pi 3.14159}]
}

func ExampleMarshal() {
	arr := bson.A{"a", "b", "c"}

	b, err := bson.Marshal(arr)
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(b))

	// Output:
	// 2000000002300002000000610002310002000000620002320002000000630000
}

func ExampleMarshalTo() {
	arr := bson.A{"a", "b", "c"}

	buf := make([]byte, 0, 128)

	buf, err := bson.MarshalTo(buf, arr)
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(buf))

	// Output:
	// 2000000002300002000000610002310002000000620002320002000000630000
}
