package giocanvas

import (
	"testing"
)

func BenchmarkC0(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ColorLookup("red")
	}
}

func BenchmarkC1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ColorLookup("rgb(100)")
	}
}

func BenchmarkC2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ColorLookup("rgb(100,100)")
	}
}

func BenchmarkC3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ColorLookup("rgb(100,100,100)")
	}
}

func BenchmarkC4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ColorLookup("rgb(100,100,100,100)")
	}
}
