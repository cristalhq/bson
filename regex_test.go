package bson

import "testing"

func TestRegex(t *testing.T) {
	re := Regex{
		Pattern: "h.ll?ow(or)ld",
		Options: "i",
	}

	b, err := re.MarshalBSON()
	mustOk(t, err)
	wantBytes(t, b, "682e6c6c3f6f77286f72296c64006900")

	var re2 Regex
	err = re2.UnmarshalBSON(b)
	mustOk(t, err)
	mustEqual(t, re2.Pattern, re.Pattern)
	mustEqual(t, re2.Options, re.Options)
}

func TestRegexCompile(t *testing.T) {
	re := Regex{
		Pattern: "h.ll?ow(or)ld",
		Options: "i",
	}

	r, err := re.Compile()
	mustOk(t, err)
	mustEqual(t, r.String(), `(?i)h.ll?ow(or)ld`)
}
