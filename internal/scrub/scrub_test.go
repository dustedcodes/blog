package scrub

import (
	"testing"
)

func Test_Middle(t *testing.T) {

	cases := map[string]string{
		"123":                            "********",
		"1234567":                        "********",
		"12345678":                       "********",
		"123456789":                      "12********89",
		"1234567890":                     "12********90",
		"1234567890a":                    "123********90a",
		"1234567890ab":                   "123********0ab",
		"1234567890abcdefghi":            "123********ghi",
		"1234567890abcdefghij":           "12345********fghij",
		"1234567890abcdefghijABCDEFGHI":  "12345********EFGHI",
		"1234567890abcdefghijABCDEFGHIJ": "1234567890********ABCDEFGHIJ",
	}

	for k, v := range cases {
		o := Middle(k)
		if o != v {
			t.Errorf("Expected %s, got %s", v, o)
		}
	}
}
