package emissione_test

import (
	"github.com/kernle32dll/emissione-go"

	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestAcceptSlice_Len(t *testing.T) {
	tests := []struct {
		name string
		p    emissione.AcceptSlice
		want int
	}{
		{"empty", emissione.AcceptSlice{}, 0},
		{"one", emissione.AcceptSlice{"a"}, 1},
		{"two", emissione.AcceptSlice{"a", "b"}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAcceptSlice_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		p    emissione.AcceptSlice
		args args
		want bool
	}{
		{"itself-not-less", emissione.AcceptSlice{"image/webp"}, args{0, 0}, false},
		{"less-than-q-3", emissione.AcceptSlice{"image/webp", "image/webp;q=3"}, args{0, 1}, true},
		{"not-less-than-q-0.5", emissione.AcceptSlice{"image/webp", "image/webp;q=0.5"}, args{0, 1}, false},
		{"not-less-than-not-parsable", emissione.AcceptSlice{"image/webp", "image/webp;q=nope"}, args{0, 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAcceptSlice_Swap(t *testing.T) {
	slice := emissione.AcceptSlice{"image/webp", "image/png"}
	swapped := emissione.AcceptSlice{"image/png", "image/webp"}

	slice.Swap(0, 1)

	if fmt.Sprint(slice) != fmt.Sprint(swapped) {
		t.Errorf("Swap() = %v, want %v", swapped, slice)
	}
}

// The main purpose of AcceptSlice is to implement sorting by their quality.
// This method tests that.
func TestAcceptSlice_Sorting(t *testing.T) {
	tests := []struct {
		name string
		p    emissione.AcceptSlice
		want []string
	}{
		{"stable", emissione.AcceptSlice{"image/webp;q=1", "image/png;q=1"}, []string{"image/webp;q=1", "image/png;q=1"}},
		{"stable-with-default", emissione.AcceptSlice{"image/webp", "image/png;q=1"}, []string{"image/webp", "image/png;q=1"}},
		{"increasing-sort-int", emissione.AcceptSlice{"image/png;q=2", "image/webp;q=1"}, []string{"image/webp;q=1", "image/png;q=2"}},
		{"increasing-sort-float", emissione.AcceptSlice{"image/png;q=2.0", "image/webp;q=1.0"}, []string{"image/webp;q=1.0", "image/png;q=2.0"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(tt.p)

			if !reflect.DeepEqual([]string(tt.p), tt.want) {
				t.Errorf("Sorting = %v, want %v", tt.p, tt.want)
			}
		})
	}
}
