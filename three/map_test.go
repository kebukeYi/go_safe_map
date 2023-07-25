package three

import "testing"

// 单元测试内容

func TestCommonMap(t *testing.T) {
	t.Run("test-intIntCase", func(t *testing.T) {
		m := newCommonMap[int, int]()
		UnitTestFunc(m, &intIntCase, t)
	})
	t.Run("test-stringStringCase", func(t *testing.T) {
		m := newCommonMap[string, string]()
		UnitTestFunc(m, &stringStringCase, t)
	})
	t.Run("test-structStructCase", func(t *testing.T) {
		m := newCommonMap[testStruct, testStruct]()
		UnitTestFunc(m, &structStructCase, t)
	})
	t.Run("test-intStructCase", func(t *testing.T) {
		m := newCommonMap[int, testStruct]()
		UnitTestFunc(m, &intStructCase, t)
	})
}

func TestSimpleMap(t *testing.T) {
	t.Run("test-intIntCase", func(t *testing.T) {
		m := newSimpleMap[int, int]()
		UnitTestFunc(m, &intIntCase, t)
	})
	t.Run("test-stringStringCase", func(t *testing.T) {
		m := newSimpleMap[string, string]()
		UnitTestFunc(m, &stringStringCase, t)
	})
	t.Run("test-structStructCase", func(t *testing.T) {
		m := newSimpleMap[testStruct, testStruct]()
		UnitTestFunc(m, &structStructCase, t)
	})
	t.Run("test-intStructCase", func(t *testing.T) {
		m := newSimpleMap[int, testStruct]()
		UnitTestFunc(m, &intStructCase, t)
	})
}

func TestRWMap(t *testing.T) {
	t.Run("test-intIntCase", func(t *testing.T) {
		m := newRWMap[int, int](1)
		UnitTestFunc(m, &intIntCase, t)
	})
	t.Run("test-stringStringCase", func(t *testing.T) {
		m := newRWMap[string, string](1)
		UnitTestFunc(m, &stringStringCase, t)
	})
	t.Run("test-structStructCase", func(t *testing.T) {
		m := newRWMap[testStruct, testStruct](1)
		UnitTestFunc(m, &structStructCase, t)
	})
	t.Run("test-intStructCase", func(t *testing.T) {
		m := newRWMap[int, testStruct](1)
		UnitTestFunc(m, &intStructCase, t)
	})
}
