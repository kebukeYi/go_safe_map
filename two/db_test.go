package two

import (
	"math/rand"
	"testing"
	"time"
	"unsafe"
)

func hash(key interface{}) uint {
	gc := &key
	start := uintptr(unsafe.Pointer(gc))
	offset := unsafe.Sizeof(key)
	sizeOfByte := unsafe.Sizeof(byte(0))
	hashSum := uint(0)
	for ptr := start; ptr < start+offset; ptr += sizeOfByte {
		b := *(*byte)(unsafe.Pointer(ptr))
		hashSum += uint(b)
		hashSum = uint(b) + (hashSum << 6) + (hashSum << 16) - hashSum
	}
	return hashSum
}

var (
	db = NewIDataBase(10, 1000, 1000, hash)
)

func init() {
	for i := 0; i < 10000; i++ {
		db.Set(rand.Int(), i, time.Second)
	}
}

func BenchmarkTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 33333; j++ {
			db.Set(rand.Int(), j, time.Second)
		}
		for f := 0; f < 33333; f++ {
			db.Delete(rand.Int())
		}
		for d := 0; d < 33333; d++ {
			db.Get(rand.Int())
		}
	}
}
