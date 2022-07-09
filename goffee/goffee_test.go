package goffee

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkStringConcatenationByStringsBuilder(b *testing.B) {
	tests := []struct{ method, pattern string }{
		{"GET", "/hello"},
		{"POST", "/hello/other/:id"},
		{"PUT", "/"},
		{"DELETE", "/hello/delete/other/:id"},
	}
	for _, tt := range tests {
		var key strings.Builder
		key.WriteString(tt.method)
		key.WriteString("-")
		key.WriteString(tt.pattern)
		fmt.Println(key.String())
	}
}

func BenchmarkStringConcatenationByPlus(b *testing.B) {
	tests := []struct{ method, pattern string }{
		{"GET", "/hello"},
		{"POST", "/hello/other/:id"},
		{"PUT", "/"},
		{"DELETE", "/hello/delete/other/:id"},
	}
	for _, tt := range tests {
		key := tt.method + "-" + tt.pattern
		fmt.Println(key)
	}
}
