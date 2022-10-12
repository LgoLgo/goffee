package Lgoffee

import (
	"fmt"
	"strings"
	"testing"
)

func TestNestedGroup(t *testing.T) {
	r := New()
	v1 := r.Group("/v1")
	v2 := v1.Group("/v2")
	v3 := v2.Group("/v3")
	if v2.prefix != "/v1/v2" {
		t.Fatal("v2 prefix should be /v1/v2")
	}
	if v3.prefix != "/v1/v2/v3" {
		t.Fatal("v2 prefix should be /v1/v2")
	}
}

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
