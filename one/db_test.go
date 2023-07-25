package one

import (
	"math/rand"
	"testing"
	"time"
)

var (
	db = NewDatabase(1, 1000)
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
