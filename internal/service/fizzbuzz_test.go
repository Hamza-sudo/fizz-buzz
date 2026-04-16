package service

import (
	"reflect"
	"testing"
)

func TestFizzBuzz(t *testing.T) {
	tests := []struct {
		name     string
		params   FizzBuzzParams
		expected []string
	}{
		{
			name: "classic fizzbuzz 1-15",
			params: FizzBuzzParams{
				Int1:  3,
				Int2:  5,
				Limit: 15,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			expected: []string{
				"1", "2", "fizz", "4", "buzz",
				"fizz", "7", "8", "fizz", "buzz",
				"11", "fizz", "13", "14", "fizzbuzz",
			},
		},
		{
			name: "custom strings",
			params: FizzBuzzParams{
				Int1:  2,
				Int2:  7,
				Limit: 14,
				Str1:  "foo",
				Str2:  "bar",
			},
			expected: []string{
				"1", "foo", "3", "foo", "5", "foo", "bar",
				"foo", "9", "foo", "11", "foo", "13", "foobar",
			},
		},
		{
			name: "limit 1",
			params: FizzBuzzParams{
				Int1:  3,
				Int2:  5,
				Limit: 1,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			expected: []string{"1"},
		},
		{
			name: "int1 equals int2",
			params: FizzBuzzParams{
				Int1:  3,
				Int2:  3,
				Limit: 6,
				Str1:  "a",
				Str2:  "b",
			},
			expected: []string{"1", "2", "ab", "4", "5", "ab"},
		},
		{
			name: "int1 is zero",
			params: FizzBuzzParams{
				Int1:  0,
				Int2:  2,
				Limit: 4,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			expected: []string{"1", "buzz", "3", "buzz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FizzBuzz(tt.params)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FizzBuzz(%+v)\ngot:  %v\nwant: %v", tt.params, result, tt.expected)
			}
		})
	}
}
