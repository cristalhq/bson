package bson

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
)

// Regex represents BSON regular expression.
type Regex struct {
	Pattern string
	Options string
}

// String returns a hex string representation of the id.
func (re Regex) String() string {
	return fmt.Sprintf(`Regex('%s', '%s')`, re.Pattern, re.Options)
}

// MarshalBSON implements [bson.Marshaler].
func (re *Regex) MarshalBSON() ([]byte, error) {
	b := make([]byte, 0, len(re.Pattern)+1+len(re.Options)+1)
	b = append(b, re.Pattern...)
	b = append(b, 0)
	b = append(b, re.Options...)
	b = append(b, 0)
	return b, nil
}

// UnmarshalBSON implements [bson.Unmarshaler].
func (re *Regex) UnmarshalBSON(b []byte) error {
	idx := bytes.IndexByte(b, 0)
	if idx == -1 {
		return errors.New("malformed regex")
	}
	re.Pattern = string(b[0:idx])
	b = b[idx+1:]

	idx = bytes.IndexByte(b, 0)
	if idx == -1 {
		return errors.New("malformed regex")
	}
	re.Options = string(b[:idx])

	return nil
}

// Compile returns [regexp.Regexp].
func (r Regex) Compile() (*regexp.Regexp, error) {
	var opts string
	for _, o := range r.Options {
		switch o {
		case 'i', 'm', 's':
			opts += string(o)
		default:
		}
	}

	expr := r.Pattern
	if opts != "" {
		expr = "(?" + opts + ")" + expr
	}

	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return re, nil
}
