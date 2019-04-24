package user

import (
	"testing"
)

type testpair struct {
	value string
	out   bool
}

func TestCheckDate(t *testing.T) {

	var tests = []testpair{
		{"1956-12-01", true},
		{"1596", false},
		{"", true},
		{"1956_01_01", false},
		{"I956-I1-II", false},
		{"1956-20-20", false},
	}

	for _, pair := range tests {
		v := CheckDate(pair.value)
		if (v == nil) != pair.out {
			t.Error(
				"For", pair.value,
				"expected", pair.out,
				"got", pair.out,
			)
		}
	}
}
