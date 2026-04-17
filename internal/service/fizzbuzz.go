package service

import (
	"errors"
	"strconv"
	"strings"
)

// FizzBuzzParams contains the inputs used to build a FizzBuzz sequence.
type FizzBuzzParams struct {
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Limit int    `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
}

var (
	ErrInt1Required  = errors.New("int1 must be greater than 0")
	ErrInt2Required  = errors.New("int2 must be greater than 0")
	ErrLimitRequired = errors.New("limit must be greater than 0")
	ErrStr1Required  = errors.New("str1 must not be empty")
	ErrStr2Required  = errors.New("str2 must not be empty")
)

// Validate checks whether the parameters are usable by the API.
func (p FizzBuzzParams) Validate() error {
	switch {
	case p.Int1 <= 0:
		return ErrInt1Required
	case p.Int2 <= 0:
		return ErrInt2Required
	case p.Limit <= 0:
		return ErrLimitRequired
	case p.Str1 == "":
		return ErrStr1Required
	case p.Str2 == "":
		return ErrStr2Required
	default:
		return nil
	}
}

// FizzBuzz generates the sequence using the following rules:
// - multiples of Int1 -> Str1
// - multiples of Int2 -> Str2
// - multiples of both -> Str1+Str2
// - otherwise -> the number as a string
func FizzBuzz(p FizzBuzzParams) []string {
	result := make([]string, 0, p.Limit)

	for i := 1; i <= p.Limit; i++ {
		var sb strings.Builder

		// Concatenate both replacement strings when a number matches both divisors.
		if i%p.Int1 == 0 {
			sb.WriteString(p.Str1)
		}
		if i%p.Int2 == 0 {
			sb.WriteString(p.Str2)
		}

		// Keep the original number when no replacement rule applies.
		if sb.Len() == 0 {
			sb.WriteString(strconv.Itoa(i))
		}

		result = append(result, sb.String())
	}

	return result
}
