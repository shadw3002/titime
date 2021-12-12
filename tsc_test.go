package titime

import (
	"fmt"
	"testing"
	"time"
)

func TestA(t *testing.T) {
	core := NewTscCore()
	fmt.Println(time.Now(), core.Now())
}

func BenchmarkTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now()
	}
}

func BenchmarkTiTime(b *testing.B) {
	core := NewTscCore()
	for i := 0; i < b.N; i++ {
		core.Now()
	}
}
