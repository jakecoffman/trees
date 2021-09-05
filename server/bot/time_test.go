package main

import (
	"syscall"
	"testing"
	"time"
)

func BenchmarkTime(b *testing.B) {
	var tv syscall.Timeval

	for i := 0; i < b.N; i++ {
		_ = syscall.Gettimeofday(&tv)
	}
}

func BenchmarkTime2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now()
	}
}
