package bson

import (
	"reflect"
	"sort"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

type pairRefl struct {
	Key string
	Val reflect.Value
}

type docRefl []pairRefl

func (doc docRefl) Len() int           { return len(doc) }
func (doc docRefl) Less(i, j int) bool { return doc[j].Key < doc[j-1].Key }
func (doc docRefl) Swap(i, j int)      { doc[i], doc[j] = doc[j], doc[i] }

func sortPairRefl(doc docRefl) {
	// TODO(cristaloleg): 20 might be too high/low.
	if len(doc) <= 20 {
		// for i := 1; i < len(ss); i++ {
		// 	for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
		// 		ss[j], ss[j-1] = ss[j-1], ss[j]
		// 	}
		// }
		// below is the code above, but (almost) without bound checks

		for i := 1; i < len(doc); i++ {
			for j := i; j > 0; j-- {
				if doc[j].Key >= doc[j-1].Key {
					break
				}
				doc[j], doc[j-1] = doc[j-1], doc[j]
			}
		}
	} else {
		sort.Sort(doc)
	}
}
