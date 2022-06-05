//nolint:testpackage
package cron

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		input    string
		expected []int
	}{
		{"3.5", []int{3, 5}},
		{"3-5", []int{3, 4, 5}},
		{"*/15", []int{0, 15, 30, 45}},
		{"0-23/2", []int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22}},
		{"1-5/2,7-15/3", []int{1, 3, 5, 7, 10, 13}},
	}
	for _, v := range tests {
		if diff := cmp.Diff(
			decode(v.input, 0, 59), v.expected,
		); diff != "" {
			t.Errorf("Different parsing %s (-got +expected):\n%s\n", v.input, diff)
		}
	}
}

func TestTokenise(t *testing.T) {
	tests := []struct {
		input    string
		expected []token
	}{
		{"3,5", []token{
			{ttype: number, num: 3},
			{ttype: comma},
			{ttype: number, num: 5},
		}},
		{"3-5", []token{
			{ttype: number, num: 3},
			{ttype: dash},
			{ttype: number, num: 5},
		}},
		{"*/15", []token{
			{ttype: star},
			{ttype: slash},
			{ttype: number, num: 15},
		}},
		{"0-23/2", []token{
			{ttype: number, num: 0},
			{ttype: dash},
			{ttype: number, num: 23},
			{ttype: slash},
			{ttype: number, num: 2},
		}},
		{"1-5/2,7-15/3", []token{
			{ttype: number, num: 1},
			{ttype: dash},
			{ttype: number, num: 5},
			{ttype: slash},
			{ttype: number, num: 2},
			{ttype: comma},
			{ttype: number, num: 7},
			{ttype: dash},
			{ttype: number, num: 15},
			{ttype: slash},
			{ttype: number, num: 3},
		}},
	}
	for _, v := range tests {
		if diff := cmp.Diff(
			tokenise(v.input), v.expected,
			cmp.AllowUnexported(token{}),
		); diff != "" {
			t.Errorf("Different tokens %s (-got +expected):\n%s\n", v.input, diff)
		}
	}
}
