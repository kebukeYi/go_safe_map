package three

import (
	"fmt"
	"sync"
	"testing"
)

type TestCase[Key, Value any] struct {
	Indexes        []int
	StoreKeys      map[int]Key
	StoreValues    map[int]Value
	LoadKeys       map[int]Key
	WantLoadValues map[int]Value
	WantLoaded     map[int]bool
	DeleteKeys     map[int]Key
	RangeFunc      func(Key, Value) bool
	Equals         func(v1, v2 Value) bool
}

// 单元测试函数
func UnitTestFunc[Key, Value any](m IMap[Key, Value], testCase *TestCase[Key, Value], t *testing.T) {
	for _, index := range testCase.Indexes {
		m.Store(testCase.StoreKeys[index], testCase.StoreValues[index])
	}

	for _, index := range testCase.Indexes {
		m.Delete(testCase.DeleteKeys[index])
	}

	for _, index := range testCase.Indexes {
		got, gotten := m.Load(testCase.LoadKeys[index])
		if gotten != testCase.WantLoaded[index] {
			t.Errorf("key %v want loaded: %v, but load: %v", testCase.LoadKeys[index], testCase.WantLoaded[index], gotten)
		}
		if !testCase.Equals(got, testCase.WantLoadValues[index]) {
			t.Errorf("key %v want got: %v, but got: %v", testCase.LoadKeys[index], testCase.WantLoadValues[index], got)
		}
	}
	m.Range(testCase.RangeFunc)
}

type testStruct struct {
	val int
}

// 单元测试用例
var (
	intIntCase TestCase[int, int] = TestCase[int, int]{
		Indexes: []int{1, 2, 3},
		StoreKeys: map[int]int{
			1: 114514,
			2: 1919810,
			3: 2147,
		},
		StoreValues: map[int]int{
			1: 1919810,
			2: 114514,
			3: 65535,
		},
		DeleteKeys: map[int]int{
			1: 2147,
			2: 2147,
			3: 2147,
		},
		LoadKeys: map[int]int{
			1: 2147,
			2: 114514,
			3: 1919810,
		},
		WantLoaded: map[int]bool{
			1: false,
			2: true,
			3: true,
		},
		WantLoadValues: map[int]int{
			1: 0,
			2: 1919810,
			3: 114514,
		},
		RangeFunc: func(a, b int) bool {
			fmt.Printf("key: %d, value: %d\n", a, b)
			return true
		},
		Equals: func(a, b int) bool {
			return a == b
		},
	}

	stringStringCase TestCase[string, string] = TestCase[string, string]{
		Indexes: []int{1, 2, 3},
		StoreKeys: map[int]string{
			1: "114514",
			2: "1919810",
			3: "2147",
		},
		StoreValues: map[int]string{
			1: "1919810",
			2: "114514",
			3: "65535",
		},
		DeleteKeys: map[int]string{
			1: "2147",
			2: "2147",
			3: "2147",
		},
		LoadKeys: map[int]string{
			1: "2147",
			2: "114514",
			3: "1919810",
		},
		WantLoaded: map[int]bool{
			1: false,
			2: true,
			3: true,
		},
		WantLoadValues: map[int]string{
			1: "",
			2: "1919810",
			3: "114514",
		},
		RangeFunc: func(a, b string) bool {
			fmt.Printf("key: %v, value: %v\n", a, b)
			return true
		},
		Equals: func(a, b string) bool {
			return a == b
		},
	}

	intStructCase TestCase[int, testStruct] = TestCase[int, testStruct]{
		Indexes: []int{1, 2, 3},
		StoreKeys: map[int]int{
			1: 114514,
			2: 1919810,
			3: 2147,
		},
		StoreValues: map[int]testStruct{
			1: {val: 1919810},
			2: {val: 114514},
			3: {val: 65535},
		},
		DeleteKeys: map[int]int{
			1: 2147,
			2: 2147,
			3: 2147,
		},
		LoadKeys: map[int]int{
			1: 2147,
			2: 114514,
			3: 1919810,
		},
		WantLoaded: map[int]bool{
			1: false,
			2: true,
			3: true,
		},
		WantLoadValues: map[int]testStruct{
			1: {},
			2: {val: 1919810},
			3: {val: 114514},
		},
		RangeFunc: func(a int, b testStruct) bool {
			fmt.Printf("key: %d, value: %#v\n", a, b)
			return true
		},
		Equals: func(a, b testStruct) bool {
			return a.val == b.val
		},
	}

	structStructCase TestCase[testStruct, testStruct] = TestCase[testStruct, testStruct]{
		Indexes: []int{1, 2, 3},
		StoreKeys: map[int]testStruct{
			1: {val: 114514},
			2: {val: 1919810},
			3: {val: 2147},
		},
		StoreValues: map[int]testStruct{
			1: {val: 1919810},
			2: {val: 114514},
			3: {val: 65535},
		},
		DeleteKeys: map[int]testStruct{
			1: {val: 2147},
			2: {val: 2147},
			3: {val: 2147},
		},
		LoadKeys: map[int]testStruct{
			1: {val: 2147},
			2: {val: 114514},
			3: {val: 1919810},
		},
		WantLoaded: map[int]bool{
			1: false,
			2: true,
			3: true,
		},
		WantLoadValues: map[int]testStruct{
			1: {},
			2: {val: 1919810},
			3: {val: 114514},
		},
		RangeFunc: func(a, b testStruct) bool {
			fmt.Printf("key: %#v, value: %#v\n", a, b)
			return true
		},
		Equals: func(a, b testStruct) bool {
			return a.val == b.val
		},
	}
)

// 基准测试
// BenchmarkTestKv-4              5         217,674280 ns/op
func BenchmarkTestKv(b *testing.B) {
	goroutines, operations := 100, 10000
	for i := 0; i < b.N; i++ {
		m := newRWMap[int, int](10007)
		wg := &sync.WaitGroup{}
		wg.Add(goroutines * operations)
		for j := 0; j < goroutines; j++ {
			go func(index int) {
				for k := 0; k < operations; k++ {
					m.Store(index*operations+k, k)
					m.Load(index*operations + k)
					wg.Done()
				}
			}(j)
		}
		wg.Wait()
	}
}

// BenchmarkCommonMapTest 测试CommonMap进行100万次读和写的基准性能
// BenchmarkCommonMapTest-4               2         693,997350 ns/op
func BenchmarkCommonMapTest(b *testing.B) {
	// 只加一个全局锁
	m := newCommonMap[int, int]()
	MapBenchmarkTestFunc(m, 100, 10000, b)
}

// BenchmarkSimpleMapTest 测试SimpleMap进行100万次读和写的基准性能
// BenchmarkSimpleMapTest-4               2         530,467200 ns/op
func BenchmarkSimpleMapTest(b *testing.B) {
	// 将全局锁换成读写锁
	m := newSimpleMap[int, int]()
	MapBenchmarkTestFunc(m, 100, 10000, b)
}

// BenchmarkRWMapTest 测试RWMap进行100万次读和写的基准性能
// BenchmarkRWMapTest-4           7         157,554929 ns/op
func BenchmarkRWMapTest(b *testing.B) {
	// 分段map + 读写锁
	m := newRWMap[int, int](10007)
	MapBenchmarkTestFunc(m, 100, 10000, b)
}

// MapBenchmarkTestFunc 各种IMap的基准测试函数体
func MapBenchmarkTestFunc(m IMap[int, int], goroutines, operations int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := &sync.WaitGroup{}
		wg.Add(goroutines * operations)
		for j := 0; j < goroutines; j++ {
			go func(index int) {
				for k := 0; k < operations; k++ {
					m.Store(index*operations+k, k)
					m.Load(index*operations + k)
					wg.Done()
				}
			}(j)
		}
		wg.Wait()
	}
}
