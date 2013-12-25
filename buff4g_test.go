package buff4g

import (
	"sync"
	"testing"
)

func TestNewBlockBytes(t *testing.T) {
	bb := NewBlockBytes(1024)
	if len(bb.curBytes) != 1024 ||
		len(bb.newBytes) != 1024 {
		t.Fatal("TestNewBlockBytes failed")
	}
}

func TestInitBuffer(t *testing.T) {
	InitBuffer(1024, 4)
	if pool.poolMod != 3 {
		t.Fatal("TestInitBuffer failed")
	}
}

func TestAlloc(t *testing.T) {
	bb := NewBlockBytes(1024)
	for i := 0; i < 100; i++ {
		bytes := bb.Alloc(64)
		if len(bytes) != 64 && cap(bytes) != 64 {
			t.Fatal("Alloc failed")
		}
	}
}

// 多线程安全版本的malloc, 用之前必须先初始化: InitBuffer
func TestPAlloc(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		go func() {
			wg.Add(1)
			bytes := PAlloc(64)
			if len(bytes) != 64 && cap(bytes) != 64 {
				t.Fatal("PAlloc failed")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkAlloc(b *testing.B) {
	for j := 0; j < b.N; j++ {
		bytes := PAlloc(64)
		if len(bytes) != 64 && cap(bytes) != 64 {
			b.Fatal("PAlloc failed")
		}
	}
}

func BenchmarkPAlloc(b *testing.B) {
	wg := sync.WaitGroup{}
	for j := 0; j < b.N; j++ {
		go func() {
			for i := 0; i < 1000; i++ {
				wg.Add(1)
				bytes := PAlloc(64)
				if len(bytes) != 64 && cap(bytes) != 64 {
					b.Fatal("PAlloc failed")
				}
				wg.Done()
			}
		}()
	}
	wg.Wait()
}
