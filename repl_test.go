package main

import (
	"reflect" 
	"testing"
)


func TestCleanInput(t *testing.T) {
	cases := []struct {
		name string
		input string
		expected []string
	}{
		{name: "excessive_spaces", input: "   hello  world   ", expected: []string{"hello", "world"}},
		{name: "zero_spaces", input: "helloworld", expected: []string{"helloworld"}},
		{name: "uppercase", input: "How Are YOU?", expected: []string{"how", "are", "you?"}},
		{name: "escaped", input: "  \tHello\nWORLD\t\n  ", expected: []string{"hello", "world"}},
	}

	for _, tc := range cases {
		// copies only the header btw
		want := tc.expected
		got := cleanInput(tc.input)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("%s(), expected: %#v, got: %#v", tc.name, want, got)
		}
	}
}

