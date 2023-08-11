package bson_test

import "bson"

func Example_types() {
	arr := bson.A{"a", "b", "c", 12345}
	doc := bson.D{{"hello", "world"}, {"answer", 42}}
	doo := bson.M{"hello": "world", "pi": 3.14159}

	_, _, _ = arr, doc, doo

	// Output:
}
